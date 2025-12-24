package tests

import (
	"context"
	"encoding/json"

	"github.com/NethermindEth/starknet.go/client"
)

// The purpose of the Spy type is to spy on the JSON-RPC calls made by the client.
// It's used in the tests to mock the JSON-RPC calls and to check if the client is
// making the correct calls.
type WsSpy struct {
	wsConn
	buff  []byte
	debug bool
}

// Toggles the debug mode of the spy to the opposite of the current value.
func (s *WsSpy) ToggleDebug() {
	s.debug = !s.debug
}

// The wsConn interface used in `rpc.websocket` tests.
// It's implemented by the `client.Client` type.
type wsConn interface {
	Subscribe(
		ctx context.Context,
		namespace string,
		methodSuffix string,
		channel interface{},
		args interface{},
	) (*client.ClientSubscription, error)
	SubscribeWithSliceArgs(
		ctx context.Context,
		namespace string,
		methodSuffix string,
		channel interface{},
		args ...interface{},
	) (*client.ClientSubscription, error)
	Close()
}

// The Spyer interface implemented by the Spy type.
type WsSpyer interface {
	callCloser
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
func NewWsSpy(client wsConn, debug ...bool) WsSpyer {
	d := false
	if len(debug) > 0 {
		d = debug[0]
	}

	return &WsSpy{
		wsConn: client,
		buff:   []byte{},
		debug:  d,
	}
}

// CallContext calls the original CallContext function with the given parameters
// and captures the response.
func (s *WsSpy) CallContext(
	ctx context.Context,
	result interface{},
	method string,
	arg interface{},
) error {
	// if s.debug {
	// 	fmt.Printf("### Spy Debug mode: in parameters\n")
	// 	fmt.Printf("   arg.(%T): %+v\n", arg, arg)
	// 	PrettyPrint(arg)
	// 	fmt.Println("--------------------------------------------")
	// }

	// raw := json.RawMessage{}
	// err := s.callCloser.CallContext(ctx, &raw, method, arg)
	// if err != nil {
	// 	return err
	// }

	// if s.debug {
	// 	fmt.Printf("### Spy Debug mode: output\n")
	// 	PrettyPrint(raw)
	// }

	// err = json.Unmarshal(raw, result)
	// s.buff = raw

	// return err
	return nil
}

// Subscribe calls the original Subscribe function with the given parameters
// and captures the response.
func (s *WsSpy) Subscribe(
	ctx context.Context,
	namespace string,
	methodSuffix string,
	channel interface{},
	args interface{},
) error {
	// if s.debug {
	// 	fmt.Printf("### Spy Debug mode: in parameters\n")
	// 	fmt.Printf("   args.(%T): %+v\n", args, args)
	// 	PrettyPrint(args)
	// 	fmt.Println("--------------------------------------------")
	// }

	// raw := json.RawMessage{}
	// err := s.wsConn.Subscribe(ctx, namespace, methodSuffix, channel, args)
	// if err != nil {
	// 	return err
	// }

	// if s.debug {
	// 	fmt.Printf("### Spy Debug mode: output\n")
	// 	PrettyPrint(raw)
	// }

	// err = json.Unmarshal(raw, result)
	// s.buff = raw

	// return err
	return nil
}

// CallContextWithSliceArgs calls the original CallContextWithSliceArgs function with the given parameters
// and captures the response.
func (s *WsSpy) CallContextWithSliceArgs(
	ctx context.Context,
	result interface{},
	method string,
	args ...interface{},
) error {
	// if s.debug {
	// 	fmt.Printf("### Spy Debug mode: in parameters\n")
	// 	for i, v := range args {
	// 		fmt.Printf("   Arg[%d].(%T): %+v\n", i, v, v)
	// 		PrettyPrint(v)
	// 		fmt.Println("--------------------------------------------")
	// 	}
	// }

	// raw := json.RawMessage{}
	// err := s.callCloser.CallContextWithSliceArgs(ctx, &raw, method, args...)
	// if err != nil {
	// 	return err
	// }

	// if s.debug {
	// 	fmt.Printf("### Spy Debug mode: output\n")
	// 	PrettyPrint(raw)
	// }

	// err = json.Unmarshal(raw, result)
	// s.buff = raw

	// return err
	return nil
}

// LastResponse returns the last response captured by the spy.
// In other words, it returns the raw JSON response received from the server when
// calling a `callCloser` method.
func (s *WsSpy) LastResponse() json.RawMessage {
	return s.buff
}
