package tests

import (
	"context"
	"encoding/json"
	"fmt"
)

// The purpose of the Spy type is to spy on the JSON-RPC calls made by the client.
// It's used in the tests to mock the JSON-RPC calls and to check if the client is
// making the correct calls.
type Spy struct {
	callCloser
	buff  []byte
	mock  bool
	debug bool
}

// The callCloser interface used in `rpc` and `paymaster` tests.
// It's implemented by the `client.Client` type.
type callCloser interface {
	CallContext(ctx context.Context, result interface{}, method string, args interface{}) error
	CallContextWithSliceArgs(
		ctx context.Context,
		result interface{},
		method string,
		args ...interface{},
	) error
	Close()
}

// The Spyer interface implemented by the Spy type.
type Spyer interface {
	CallContext(ctx context.Context, result interface{}, method string, args interface{}) error
	CallContextWithSliceArgs(
		ctx context.Context,
		result interface{},
		method string,
		args ...interface{},
	) error
	Close()
	LastResponse() json.RawMessage
}

// Assert that the Spy type implements the callCloser and Spyer interfaces.
var (
	_ callCloser = (*Spy)(nil)
	_ Spyer      = (*Spy)(nil)
)

// NewJSONRPCSpy creates a new spy object.
//
// It takes a client callCloser as the first parameter and an optional debug parameter.
// The client callCloser is the interface that the spy will be based on.
// The debug parameter is a variadic parameter that specifies whether debug mode is enabled.
//
// Parameters:
//   - client: the interface that the spy will be based on
//   - debug: a boolean flag indicating whether to print debug information
//
// Returns:
//   - spy: a new spy object
func NewJSONRPCSpy(client callCloser, debug ...bool) Spyer {
	d := false
	if len(debug) > 0 {
		d = debug[0]
	}
	if TEST_ENV == MockEnv {
		return &Spy{
			callCloser: client,
			buff:       []byte{},
			mock:       true,
			debug:      d,
		}
	}

	return &Spy{
		callCloser: client,
		buff:       []byte{},
		mock:       false,
		debug:      d,
	}
}

// CallContext calls the spy function with the given context, result, method, and arguments.
//
// Parameters:
//   - ctx: the context.Context to be used.
//   - result: the interface{} to store the result of the function call.
//   - method: the string representing the method to be called.
//   - arg: argument to be passed to the function call.
//
// Returns:
//   - error: an error if any occurred during the function call
func (s *Spy) CallContext(
	ctx context.Context,
	result interface{},
	method string,
	arg interface{},
) error {
	if s.mock {
		return s.callCloser.CallContext(ctx, result, method, arg)
	}

	if s.debug {
		fmt.Printf("### Spy Debug mode: in parameters\n")
		fmt.Printf("   arg.(%T): %+v\n", arg, arg)
		PrettyPrint(arg)
		fmt.Println("--------------------------------------------")
	}

	raw := json.RawMessage{}
	err := s.callCloser.CallContext(ctx, &raw, method, arg)
	if err != nil {
		return err
	}

	if s.debug {
		fmt.Printf("### Spy Debug mode: output\n")
		PrettyPrint(raw)
	}

	err = json.Unmarshal(raw, result)
	s.buff = raw

	return err
}

// CallContextWithSliceArgs calls the spy CallContext function with args as a slice.
//
// Parameters:
//   - ctx: the context.Context to be used.
//   - result: the interface{} to store the result of the function call.
//   - method: the string representing the method to be called.
//   - args: variadic arguments to be passed to the function call.
//
// Returns:
//   - error: an error if any occurred during the function call
func (s *Spy) CallContextWithSliceArgs(
	ctx context.Context,
	result interface{},
	method string,
	args ...interface{},
) error {
	if s.mock {
		return s.callCloser.CallContextWithSliceArgs(ctx, result, method, args...)
	}

	if s.debug {
		fmt.Printf("### Spy Debug mode: in parameters\n")
		for i, v := range args {
			fmt.Printf("   Arg[%d].(%T): %+v\n", i, v, v)
			PrettyPrint(v)
			fmt.Println("--------------------------------------------")
		}
	}

	raw := json.RawMessage{}
	err := s.callCloser.CallContextWithSliceArgs(ctx, &raw, method, args...)
	if err != nil {
		return err
	}

	if s.debug {
		fmt.Printf("### Spy Debug mode: output\n")
		PrettyPrint(raw)
	}

	err = json.Unmarshal(raw, result)
	s.buff = raw

	return err
}

// LastResponse returns the last response captured by the spy.
// In other words, it returns the raw JSON response received from the server when
// calling a `callCloser` method.
func (s *Spy) LastResponse() json.RawMessage {
	return s.buff
}

// PrettyPrint pretty marshals the data with indentation and prints it.
func PrettyPrint(data interface{}) {
	prettyJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println("Raw data:")
	fmt.Println(string(prettyJSON))
}
