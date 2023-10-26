package main

import (
	"bytes"
	"fmt"
	"go/format"
	"go/importer"
	"go/types"
	"sort"
	"strings"
	"text/template"

	_ "github.com/ledgerwatch/erigon/cl/abstract"
	"github.com/pkg/errors"
	"golang.org/x/tools/go/packages"
)

// this code was adapted from
// "github.com/eiiches/go-gen-proxy/pkg/handler"
// it was broken (it uses Importer to import instead of packages, and importer doesnt work for some reason here
// also it was meant to be generic, but im not goint to keep it that way, and just assume this is only used here ever

func main() {
	x, err := Generate("github.com/ledgerwatch/erigon/cl/cltrace", "cltrace", "BeaconStateProxy", []string{
		"github.com/ledgerwatch/erigon/cl/abstract.BeaconStateUncopiable",
	})
	if err != nil {
		panic(err)
	}
	fmt.Print(x)
}

func Generate(pkgPath string, pkgName, structName string, ifacePathAndNames []string) (string, error) {
	imp := importer.Default()

	ifaces := []*types.Interface{}
	for _, ifacePathAndName := range ifacePathAndNames {
		iface, err := resolveInterface(imp, ifacePathAndName)
		if err != nil {
			return "", errors.Wrapf(err, "failed to resolve interface")
		}
		ifaces = append(ifaces, iface)
	}

	funcs, err := mergeInterfaceMethods(ifaces)
	if err != nil {
		return "", errors.Wrapf(err, "failed to merge interface methods")
	}

	in := buildTemplateInput(pkgPath, pkgName, structName, funcs)

	out, err := renderTemplate(in)
	if err != nil {
		return "", errors.Wrapf(err, "failed to render template")
	}

	return out, nil
}

func resolveInterface(imp types.Importer, ifacePathAndName string) (*types.Interface, error) {
	splitPos := strings.LastIndex(ifacePathAndName, ".")
	if splitPos < 0 {
		return nil, errors.Errorf("interface must be IMPORT_PATH.INTERFACE_NAME format")
	}
	pkgPath := ifacePathAndName[0:splitPos]
	ifaceName := ifacePathAndName[splitPos+1:]

	//pkg, err := imp.Import(pkgPath)
	//if err != nil {
	//	return nil, err
	//}
	conf := &packages.Config{Mode: packages.NeedTypes | packages.NeedTypesInfo}
	pkgs, err := packages.Load(conf, pkgPath)
	if err != nil {
		return nil, err
	}
	if len(pkgs) < 1 {
		return nil, errors.Errorf("pkg not found")
	}
	pkg := pkgs[0].Types

	obj := pkg.Scope().Lookup(ifaceName)
	if obj == nil {
		return nil, errors.Errorf("%s is not found", ifacePathAndName)
	}
	if !types.IsInterface(obj.Type()) {
		return nil, errors.Errorf("%s is not an interface", ifacePathAndName)
	}

	return obj.Type().Underlying().(*types.Interface), nil
}

func mergeInterfaceMethods(ifaces []*types.Interface) ([]*types.Func, error) {
	funcs := map[string]*types.Func{}
	for _, iface := range ifaces {
		iface = iface.Complete()
		for i := 0; i < iface.NumMethods(); i += 1 {
			f := iface.Method(i)

			if ff, ok := funcs[f.Name()]; ok {
				if types.Identical(ff.Type(), f.Type()) {
					continue
				} else {
					return nil, errors.Errorf("conflicting method")
				}
			}

			funcs[f.Name()] = f
		}
	}

	// sort funcs alphabetically
	names := []string{}
	for _, f := range funcs {
		names = append(names, f.Name())
	}
	sort.Strings(names)
	fs := []*types.Func{}
	for _, name := range names {
		fs = append(fs, funcs[name])
	}
	return fs, nil
}

func buildTemplateInput(pkgPath string, pkgName, name string, funcs []*types.Func) *TemplateInput {
	importAliaser := NewImportAliaser(pkgPath)

	methods := []*Method{}
	for _, f := range funcs {
		signature := f.Type().(*types.Signature)
		toParams := func(tuple *types.Tuple, variadic bool) []*Param {
			params := []*Param{}
			for i := 0; i < tuple.Len(); i += 1 {
				v := tuple.At(i)

				var typ types.Type
				if variadic && i == tuple.Len()-1 {
					typ = v.Type().(*types.Slice).Elem()
				} else {
					typ = v.Type()
				}

				importAliaser.AssignAliasFromType(typ)

				param := &Param{
					Name: v.Name(),
					Type: types.TypeString(typ, importAliaser.GetQualifier),
				}

				params = append(params, param)
			}
			return params
		}
		methods = append(methods, &Method{
			Name:     f.Name(),
			Params:   toParams(signature.Params(), signature.Variadic()),
			Variadic: signature.Variadic(),
			Results:  toParams(signature.Results(), false),
		})
	}

	in := &TemplateInput{
		PackageName: pkgName,
		Name:        name,
		Methods:     methods,
		Imports:     importAliaser.GetAllAliases(),
	}

	return in
}

func renderTemplate(in *TemplateInput) (string, error) {
	buf := &bytes.Buffer{}

	templateFunctions := template.FuncMap{
		"add": add,
	}

	tmpl := template.Must(template.New("template.go").Funcs(templateFunctions).Parse(t))
	if err := tmpl.Execute(buf, in); err != nil {
		return "", err
	}

	bytes, err := format.Source(buf.Bytes())
	if err != nil {
		return "", errors.Wrapf(err, "failed to format source; body = %s", buf.String())
	}

	return string(bytes), nil
}

type TemplateInput struct {
	PackageName string
	Name        string
	Methods     []*Method
	Imports     []*ImportAlias
}

type Method struct {
	Name string

	Params   []*Param
	Variadic bool

	Results []*Param
}

type Param struct {
	Name string
	Type string
}

func add(a, b int) int {
	return a + b
}

const (
	t = `// Code generated by go-gen-proxy; DO NOT EDIT.
package {{ .PackageName }}

import (
	"fmt"
	"github.com/ledgerwatch/erigon/cl/abstract"

{{- range $import := .Imports }}
	{{ if not $import.Implicit }}{{ $import.Alias }} {{ end }}"{{ $import.Path }}"
{{- end }}
)

type {{ .Name }} struct {
	Handler InvocationHandler
	Underlying abstract.BeaconState
}

{{ range $method := .Methods }}

func (this *{{ $.Name }}) {{ $method.Name }}(
	{{- range $index, $param := $method.Params -}}
		{{- if ne $index 0 }}, {{ end -}}
		arg{{ $index }} {{ if and $method.Variadic (eq (add $index 1) (len $method.Params)) }} ... {{ end }} {{ $param.Type }}
	{{- end -}}
) (
	{{- range $index, $res := .Results -}}
		{{- if ne $index 0 }}, {{ end -}}
		ret{{$index}} {{ $res.Type }}
	{{- end -}}
) {
	args := []any{
		{{- range $index, $param := $method.Params }}
			{{- if or (not $method.Variadic) (lt (add $index 1) (len $method.Params)) }}
				arg{{ $index }},
			{{- end }}
		{{- end}}
	}
	{{- if $method.Variadic }}
	for _, arg := range arg{{ add (len $method.Params) -1 }} {
		args = append(args, arg)
	}
	{{- end }}
	rets, intercept := this.Handler.Invoke("{{ .Name }}", args)
	_ = rets
	if !intercept {
	{{- if $method.Results -}}{{range $index, $res := $method.Results -}}
		{{- if ne $index 0 }}, {{ end -}}
		ret{{$index}}
	{{- end}} = {{- end -}}this.Underlying.{{$method.Name}}(
		{{- range $index, $param := $method.Params }}
			{{- if or (not $method.Variadic) (lt (add $index 1) (len $method.Params)) }}
				arg{{ $index }},
			{{- end }}
		{{- end}}
	)
		return
	}
  {{ if $method.Results }}
	var ok bool
    {{ range $index, $res := $method.Results }}
  	ret{{ $index }}, ok = rets[{{ $index }}].({{ $res.Type }})
  	if rets[{{ $index }}] != nil && !ok {
	  	panic(fmt.Sprintf("%+v is not a valid {{ $res.Type }} value", rets[{{ $index }}]))
  	}
	  {{- end }}
		return {{ range $index, $res := $method.Results -}}
			{{- if ne $index 0 }}, {{ end -}}
			ret{{ $index }}
		{{- end -}}
	{{ end }}
}

{{end}}
`
)
