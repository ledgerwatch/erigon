package types

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	// "io"

	// rlp2 "github.com/ledgerwatch/erigon-lib/rlp"
	"github.com/ledgerwatch/erigon/rlp"
)

// type Request struct {
// 	inner RequestData
// }

type Request interface {
	EncodeRLP(io.Writer) error
	DecodeRLP([]byte) error
	RequestType() byte
	copy() Request
	EncodingSize() int
}

// func (r *Request) Type() byte {
//     return r.inner.requestType()
// }

// func NewRequest(inner RequestData) *Request {
// 	req := new(Request)
// 	req.inner = inner.copy()
// 	return req
// }

// func (r *Request) EncodingSize() int {
// 	switch r.Type() {
// 	case DepositRequestType:
// 		total := r.inner.encodingSize() + 1 // +1 byte for requset type
// 		return rlp2.ListPrefixLen(total) + total
// 	case WithdrawalRequestType:
// 		total := r.inner.EncodingSize() + 1
// 		return rlp2.ListPrefixLen(total) + total
// 	default:
// 		panic(fmt.Sprintf("Unknown request type: %d", r.Type()))
// 	}
// }

// func (r *Request) EncodeRLP(w io.Writer) error {
// 	var buf bytes.Buffer    // TODO(racytech): find a solution to reuse the same buffer instead of recreating it
// 	buf.WriteByte(r.Type()) // first write type of request then encode inner data
// 	r.inner.encodeRLP(&buf)
// 	return rlp.Encode(w, buf.Bytes())
// }

// func (r *Request) DecodeRLP(s *rlp.Stream) error {
// 	kind, _, err := s.Kind()
// 	switch {
// 	case err != nil:
// 		return err
// 	case kind == rlp.List:
// 		return fmt.Errorf("error: untyped request (unexpected lit)")
// 	case kind == rlp.Byte:
// 		return fmt.Errorf("error: too short request")
// 	default:
// 		var buf []byte
// 		if buf, err = s.Bytes(); err != nil {
// 			return err
// 		}
// 		return r.decode(buf)
// 	}
// }

func decode(data []byte) (Request, error) {
	if len(data) <= 1 {
		return nil, fmt.Errorf("error: too short type request")
	}
	var req Request
	switch data[0] {
	case DepositRequestType:
		req = new(Deposit)
	case WithdrawalRequestType:
		req = new(WithdrawalRequest)
	default:
		return nil, fmt.Errorf("unknown request type - %d", data[0])
	}

	if err := req.DecodeRLP(data); err != nil {
		return nil, err
	}
	return req, nil
}

type Requests []Request

func (r *Requests) DecodeRLP(s *rlp.Stream) (err error) {
	if _, err = s.List(); err != nil {
		if errors.Is(err, rlp.EOL) {
			r = nil
			return nil
		}
		return fmt.Errorf("read requests: %v", err)
	}
	for err == nil {
		var req Request
		kind, _, err := s.Kind()
		switch {
		case err != nil:
			return err
		case kind == rlp.List:
			return fmt.Errorf("error: untyped request (unexpected lit)")
		case kind == rlp.Byte:
			return fmt.Errorf("error: too short request")
		default:
			var buf []byte
			if buf, err = s.Bytes(); err != nil {
				return err
			}
			if req, err = decode(buf); err != nil {
				break
			}
			*r = append(*r, req)
		}
	}
	return err
}

func (r Requests) Deposits() Deposits {
	deposits := make(Deposits, 0, len(r))
	for _, req := range r {
		if req.RequestType() == DepositRequestType {
			deposits = append(deposits, req.(*Deposit))
		}
	}
	return deposits
}

func MarshalRequestsBinary(requests Requests) ([][]byte, error) {
	var ret [][]byte
	for _, req := range requests {
		buf := new(bytes.Buffer)
		if err := req.EncodeRLP(buf); err != nil {
			return nil, err
		}
		ret = append(ret, buf.Bytes())
	}
	return ret, nil
}

func UnmarshalRequestsFromBinary(requests [][]byte) (reqs Requests, err error) {
	for _, b := range requests {
		switch b[0] {
		case DepositRequestType:
			d := new(Deposit)
			if err = d.DecodeRLP(b); err != nil {
				return nil, err
			}
			reqs = append(reqs, d)
		case WithdrawalRequestType:
			w := new(WithdrawalRequest)
			if err = w.DecodeRLP(b); err != nil {
				return nil, err
			}
			reqs = append(reqs, w)
		default:
			continue
		}	
	}
	return
}

func (r Requests) Len() int { return len(r) }

// EncodeIndex encodes the i'th request to w. Note that this does not check for errors
// because we assume that *request will only ever contain valid requests that were either
// constructed by decoding or via public API in this package.
func (r Requests) EncodeIndex(i int, w *bytes.Buffer) {
	r[i].EncodeRLP(w)
}
