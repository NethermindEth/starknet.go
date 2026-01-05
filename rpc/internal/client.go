package internal

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/cookiejar"

	"github.com/NethermindEth/starknet.go/client"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/gorilla/websocket"
	"golang.org/x/net/publicsuffix"
)

// @todo see docs

// Do is a function that performs a remote procedure call (RPC) using the
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
func Do(
	ctx context.Context,
	c rpc.Caller,
	method string,
	data interface{},
	args ...interface{},
) error {
	var raw json.RawMessage
	err := c.CallContextWithSliceArgs(ctx, &raw, method, args...)
	if err != nil {
		return err
	}
	if len(raw) == 0 {
		return errors.New("not found")
	}
	if err := json.Unmarshal(raw, &data); err != nil {
		return err
	}

	return nil
}

// DoAsObject is a function that performs a remote procedure call (RPC) using
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
func DoAsObject(ctx context.Context, c rpc.Caller, method string, data, arg interface{}) error {
	var raw json.RawMessage
	err := c.CallContext(ctx, &raw, method, arg)
	if err != nil {
		return err
	}
	if len(raw) == 0 {
		return errors.New("not found")
	}
	if err := json.Unmarshal(raw, &data); err != nil {
		return err
	}

	return nil
}

// @new
func NewHTTPClient(
	ctx context.Context,
	url string,
	options ...client.ClientOption,
) (*client.Client, error) {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, err
	}
	httpClient := &http.Client{Jar: jar} //nolint:exhaustruct // Only the Jar field is used.
	// prepend the custom client to allow users to override
	options = append([]client.ClientOption{client.WithHTTPClient(httpClient)}, options...)
	c, err := client.DialOptions(ctx, url, options...)
	if err != nil {
		return nil, err
	}

	return c, nil
}

// @new
func NewWSClient(
	ctx context.Context,
	url string,
	options ...client.ClientOption,
) (*client.Client, error) {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, err
	}
	dialer := websocket.Dialer{Jar: jar} //nolint:exhaustruct // Only the Jar field is used.

	// prepend the custom client to allow users to override
	options = append([]client.ClientOption{client.WithWebsocketDialer(dialer)}, options...)
	c, err := client.DialOptions(ctx, url, options...)
	if err != nil {
		return nil, err
	}

	return c, nil
}
