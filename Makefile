GO = go
GOBIN = $(CURDIR)/build/bin
GOTEST = GODEBUG=cgocheck=0 $(GO) test ./... -p 2

GIT_COMMIT ?= $(shell git rev-list -1 HEAD)
GIT_BRANCH ?= $(shell git rev-parse --abbrev-ref HEAD)
GIT_TAG    ?= $(shell git describe --tags `git rev-list --tags="v*" --max-count=1`)
GOBUILD = $(GO) build -trimpath -ldflags "-X github.com/ledgerwatch/erigon/params.GitCommit=${GIT_COMMIT} -X github.com/ledgerwatch/erigon/params.GitBranch=${GIT_BRANCH} -X github.com/ledgerwatch/erigon/params.GitTag=${GIT_TAG}"
GO_DBG_BUILD = $(GO) build -trimpath -tags=debug -ldflags "-X github.com/ledgerwatch/erigon/params.GitCommit=${GIT_COMMIT} -X github.com/ledgerwatch/erigon/params.GitBranch=${GIT_BRANCH} -X github.com/ledgerwatch/erigon/params.GitTag=${GIT_TAG}" -gcflags=all="-N -l"  # see delve docs

GO_MAJOR_VERSION = $(shell $(GO) version | cut -c 14- | cut -d' ' -f1 | cut -d'.' -f1)
GO_MINOR_VERSION = $(shell $(GO) version | cut -c 14- | cut -d' ' -f1 | cut -d'.' -f2)

all: erigon hack rpctest state pics rpcdaemon integration db-tools sentry txpool

go-version:
	@if [ $(GO_MINOR_VERSION) -lt 16 ]; then \
		echo "minimum required Golang version is 1.16"; \
		exit 1 ;\
	fi

docker:
	DOCKER_BUILDKIT=1 docker build -t erigon:latest --build-arg git_commit='${GIT_COMMIT}' --build-arg git_branch='${GIT_BRANCH}' --build-arg git_tag='${GIT_TAG}' .

xdg_data_home :=  ~/.local/share
ifdef XDG_DATA_HOME
	xdg_data_home = $(XDG_DATA_HOME)
endif
docker-compose:
	mkdir -p $(xdg_data_home)/erigon $(xdg_data_home)/erigon-grafana $(xdg_data_home)/erigon-prometheus; \
	docker-compose up

# debug build allows see C stack traces, run it with GOTRACEBACK=crash. You don't need debug build for C pit for profiling. To profile C code use SETCGOTRCKEBACK=1
dbg:
	$(GO_DBG_BUILD) -o $(GOBIN)/ ./cmd/...

geth: erigon

erigon: go-version git-submodules
	@echo "Building Erigon"
	rm -f $(GOBIN)/tg # Remove old binary to prevent confusion where users still use it because of the scripts
	$(GOBUILD) -o $(GOBIN)/erigon ./cmd/erigon
	@echo "Run \"$(GOBIN)/erigon\" to launch Erigon."

hack:
	$(GOBUILD) -o $(GOBIN)/hack ./cmd/hack
	@echo "Run \"$(GOBIN)/hack\" to launch hack."

rpctest:
	$(GOBUILD) -o $(GOBIN)/rpctest ./cmd/rpctest
	@echo "Run \"$(GOBIN)/rpctest\" to launch rpctest."

state:
	$(GOBUILD) -o $(GOBIN)/state ./cmd/state
	@echo "Run \"$(GOBIN)/state\" to launch state."


pics:
	$(GOBUILD) -o $(GOBIN)/pics ./cmd/pics
	@echo "Run \"$(GOBIN)/pics\" to launch pics."

rpcdaemon:
	$(GOBUILD) -o $(GOBIN)/rpcdaemon ./cmd/rpcdaemon
	@echo "Run \"$(GOBIN)/rpcdaemon\" to launch rpcdaemon."

txpool:
	$(GOBUILD) -o $(GOBIN)/txpool ./cmd/txpool
	@echo "Run \"$(GOBIN)/txpool\" to launch txpool."

integration:
	$(GOBUILD) -o $(GOBIN)/integration ./cmd/integration
	@echo "Run \"$(GOBIN)/integration\" to launch integration tests."

sentry:
	$(GOBUILD) -o $(GOBIN)/sentry ./cmd/sentry
	rm -f $(GOBIN)/headers # Remove old binary to prevent confusion where users still use it because of the scripts
	@echo "Run \"$(GOBIN)/sentry\" to run sentry"

cons:
	$(GOBUILD) -o $(GOBIN)/cons ./cmd/cons
	@echo "Run \"$(GOBIN)/cons\" to run consensus engine PoC."

evm:
	$(GOBUILD) -o $(GOBIN)/evm ./cmd/evm
	@echo "Run \"$(GOBIN)/evm\" to run EVM"

downloader: git-submodules
	$(GOBUILD) -o $(GOBIN)/downloader ./cmd/downloader
	@echo "Run \"$(GOBIN)/downloader\" to download and seed snapshots."

devnettest:
	$(GOBUILD) -o $(GOBIN)/devnettest ./cmd/devnettest
	@echo "Run \"$(GOBIN)/devnettest\" to launch devnettest."

db-tools: git-submodules
	@echo "Building db-tools"

	cd libmdbx && MDBX_BUILD_TIMESTAMP=unknown make tools
	cp libmdbx/mdbx_chk $(GOBIN)
	cp libmdbx/mdbx_copy $(GOBIN)
	cp libmdbx/mdbx_dump $(GOBIN)
	cp libmdbx/mdbx_drop $(GOBIN)
	cp libmdbx/mdbx_load $(GOBIN)
	cp libmdbx/mdbx_stat $(GOBIN)
	@echo "Run \"$(GOBIN)/mdbx_stat -h\" to get info about mdbx db file."

test:
	$(GOTEST) --timeout 30m

lint:
	@./build/bin/golangci-lint run --config ./.golangci.yml

lintci:
	@echo "--> Running linter for code"
	@./build/bin/golangci-lint run --config ./.golangci.yml

lintci-deps:
	rm -f ./build/bin/golangci-lint
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./build/bin v1.42.1

clean:
	go clean -cache
	rm -fr build/*
	cd libmdbx/ && make clean
	./build/bin/golangci-lint cache clean

# The devtools target installs tools required for 'go generate'.
# You need to put $GOBIN (or $GOPATH/bin) in your PATH to use 'go generate'.

devtools:
	# Notice! If you adding new binary - add it also to cmd/hack/binary-deps/main.go file
	$(GOBUILD) -o $(GOBIN)/stringer golang.org/x/tools/cmd/stringer
	$(GOBUILD) -o $(GOBIN)/go-bindata github.com/kevinburke/go-bindata/go-bindata
	$(GOBUILD) -o $(GOBIN)/gencodec github.com/fjl/gencodec
	$(GOBUILD) -o $(GOBIN)/codecgen github.com/ugorji/go/codec/codecgen
	$(GOBUILD) -o $(GOBIN)/abigen ./cmd/abigen
	PATH=$(GOBIN):$(PATH) go generate ./common
	PATH=$(GOBIN):$(PATH) go generate ./core/types
	PATH=$(GOBIN):$(PATH) go generate ./consensus/aura/...
	@type "npm" 2> /dev/null || echo 'Please install node.js and npm'
	@type "solc" 2> /dev/null || echo 'Please install solc'
	@type "protoc" 2> /dev/null || echo 'Please install protoc'

bindings:
	PATH=$(GOBIN):$(PATH) go generate ./tests/contracts/
	PATH=$(GOBIN):$(PATH) go generate ./core/state/contracts/

prometheus:
	docker-compose up prometheus grafana


escape:
	cd $(path) && go test -gcflags "-m -m" -run none -bench=BenchmarkJumpdest* -benchmem -memprofile mem.out

git-submodules:
	@git submodule update --init --recursive --force || true # Dockerhub using ./hooks/post-checkout to set submodules, so this line will fail on Dockerhub