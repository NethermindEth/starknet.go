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

// NewSpy creates a new spy object.
//
// It takes a client callCloser as the first parameter and an optional debug parameter.
// The client callCloser is the interface that the spy will be based on.
// The debug parameter is a variadic parameter that specifies whether debug mode is enabled.
//
// Parameters:
// - client: the interface that the spy will be based on
// - debug: a boolean flag indicating whether to print debug information
// Returns:
// - spy: a new spy object
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

// CallContext calls the spy function with the given context, result, method, and arguments.
//
// Parameters:
// - ctx: the context.Context to be used.
// - result: the interface{} to store the result of the function call.
// - method: the string representing the method to be called.
// - args: variadic arguments to be passed to the function call.
// Returns:
// - error: an error if any occurred during the function call
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

// Compare compares the spy object with the given object and returns the difference between them.
//
// Parameters:
// - o: the object to compare with the spy object
// - debug: a boolean flag indicating whether to print debug information
// Returns:
// - string: the difference between the spy object and the given object
// - error: an error if any occurred during the comparison
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
		return "", Err(InternalError, err)
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
