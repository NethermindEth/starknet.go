package rpc

import (
	"context"
	"encoding/json"

	"github.com/NethermindEth/starknet.go/client"
)

type callCloser interface {
	// CallContext calls the RPC method with the specified parameters and
	// returns an error.
	CallContext(ctx context.Context, result interface{}, method string, args interface{}) error
	// CallContextWithSliceArgs call 'CallContext' with a slice of arguments.
	// For RPC-Calls with optional arguments, use 'CallContext' instead and
	// pass a struct containing the arguments, because Juno doesn't support
	// optional arguments being passed in an array, only within an object.
	CallContextWithSliceArgs(
		ctx context.Context,
		result interface{},
		method string,
		args ...interface{},
	) error
	Close()
}

type wsConn interface {
	// Subscribe calls the "<namespace>_subscribe" method with the given arguments,
	// registering a subscription. Server notifications for the subscription are
	// sent to the given channel. The element type of the channel must match the
	// expected type of content returned by the subscription.
	Subscribe(
		ctx context.Context,
		namespace string,
		methodSuffix string,
		channel interface{},
		args interface{},
	) (*client.ClientSubscription, error)
	// SubscribeWithSliceArgs call 'Subscribe' with a slice of arguments.
	// For RPC-Subscriptions with optional arguments, use 'Subscribe' instead and pass
	// a struct containing the arguments, because Juno doesn't support optional arguments
	// being passed in an array, only within an object.
	SubscribeWithSliceArgs(
		ctx context.Context,
		namespace string,
		methodSuffix string,
		channel interface{},
		args ...interface{},
	) (*client.ClientSubscription, error)
	Close()
}

// do is a function that performs a remote procedure call (RPC) using the
// provided callCloser.
// It passes the parameters as an array in the JSON-RPC call.
//
// Parameters:
//   - ctx: represents the current execution context
//   - call: the callCloser object
//   - method: the string representing the RPC method to be called
//   - data: the interface{} to store the result of the RPC call
//   - args: variadic and can be used to pass additional arguments to the RPC method
//
// Returns:
//   - error: an error if any occurred during the function call
func do(
	ctx context.Context,
	call callCloser,
	method string,
	data interface{},
	args ...interface{},
) error {
	var raw json.RawMessage
	err := call.CallContextWithSliceArgs(ctx, &raw, method, args...)
	if err != nil {
		return err
	}
	if len(raw) == 0 {
		return errNotFound
	}
	if err := json.Unmarshal(raw, &data); err != nil {
		return err
	}

	return nil
}

// doAsObject is a function that performs a remote procedure call (RPC) using
// the provided callCloser. It passes the parameter as an object in the JSON-RPC
// call, used for RPC-Calls with optional arguments since Juno doesn't support
// optional arguments being passed in an array, only within an object.
//
// Parameters:
//   - ctx: represents the current execution context
//   - call: the callCloser object
//   - method: the string representing the RPC method to be called
//   - data: the interface{} to store the result of the RPC call
//   - arg: the interface{} to pass as an object to the RPC method
//
// Returns:
//   - error: an error if any occurred during the function call
func doAsObject(ctx context.Context, call callCloser, method string, data, arg interface{}) error {
	var raw json.RawMessage
	err := call.CallContext(ctx, &raw, method, arg)
	if err != nil {
		return err
	}
	if len(raw) == 0 {
		return errNotFound
	}
	if err := json.Unmarshal(raw, &data); err != nil {
		return err
	}

	return nil
}
