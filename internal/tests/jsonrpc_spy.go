package tests

import (
	"context"
	"encoding/json"
	"fmt"
)

// The purpose of the RPCSpy type is to spy on the JSON-RPC calls made by the client.
// It's used in the tests to observe and store the responses from the JSON-RPC calls.
type RPCSpy struct {
	callCloser
	buff  []byte
	debug bool
}

// Toggles the debug mode of the spy to the opposite of the current value.
func (s *RPCSpy) ToggleDebug() {
	s.debug = !s.debug
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

// The RPCSpyer interface implemented by the Spy type.
type RPCSpyer interface {
	CallContext(ctx context.Context, result interface{}, method string, args interface{}) error
	CallContextWithSliceArgs(
		ctx context.Context,
		result interface{},
		method string,
		args ...interface{},
	) error
	Close()
	LastResponse() json.RawMessage
	ToggleDebug()
}

// Assert that the RPCSpy type implements the RPCSpyer interface.
var (
	_ RPCSpyer = (*RPCSpy)(nil)
)

// NewRPCSpy creates a new RPCSpy object.
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
//   - RPCSpyer: a new RPCSpy object that implements the RPCSpyer interface
func NewRPCSpy(client callCloser, debug ...bool) RPCSpyer {
	d := false
	if len(debug) > 0 {
		d = debug[0]
	}

	return &RPCSpy{
		callCloser: client,
		buff:       []byte{},
		debug:      d,
	}
}

// CallContext calls the original CallContext function with the given parameters
// and captures the response.
func (s *RPCSpy) CallContext(
	ctx context.Context,
	result interface{},
	method string,
	arg interface{},
) error {
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

// CallContextWithSliceArgs calls the original CallContextWithSliceArgs function with the given parameters
// and captures the response.
func (s *RPCSpy) CallContextWithSliceArgs(
	ctx context.Context,
	result interface{},
	method string,
	args ...interface{},
) error {
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
func (s *RPCSpy) LastResponse() json.RawMessage {
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
