package rpc

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/nsf/jsondiff"
)

type spy struct {
	callCloser
	s     []byte
	mock  bool
	debug bool
}

func NewSpy(client callCloser, debug ...bool) *spy {
	d := false
	if len(debug) > 0 {
		d = debug[0]
	}
	if _, ok := client.(*rpcMock); ok {
		return &spy{
			callCloser: client,
			s:          []byte{},
			mock:       true,
			debug:      d,
		}
	}
	return &spy{
		callCloser: client,
		s:          []byte{},
		debug:      d,
	}
}

func (s *spy) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	if s.mock {
		return s.callCloser.CallContext(ctx, result, method, args...)
	}
	raw := json.RawMessage{}
	if s.debug {
		fmt.Printf("... in parameters\n")
		for k, v := range args {
			fmt.Printf("   arg[%d].(%T): %+v\n", k, v, v)
		}
	}
	err := s.callCloser.CallContext(ctx, &raw, method, args...)
	if err != nil {
		return err
	}
	if s.debug {
		fmt.Printf("... output\n")
		data, err := raw.MarshalJSON()
		if err != nil {
			return err
		}
		fmt.Println("output:", string(data))
	}

	err = json.Unmarshal(raw, result)
	s.s = []byte(raw)
	return err
}

func (s *spy) Compare(o interface{}, debug bool) (string, error) {
	if s.mock {
		if debug {
			fmt.Println("**************************")
			fmt.Println("This is a mock")
			fmt.Println("**************************")
		}
		return "FullMatch", nil
	}
	b, err := json.Marshal(o)
	if err != nil {
		return "", err
	}
	diff, _ := jsondiff.Compare(s.s, b, &jsondiff.Options{})
	if debug {
		fmt.Println("**************************")
		fmt.Println(string(s.s))
		fmt.Println("**************************")
		fmt.Println(string(b))
		fmt.Println("**************************")
	}
	return diff.String(), nil
}
