package rpc

import (
	"context"
	"net/http"
	"net/http/cookiejar"

	"github.com/NethermindEth/starknet.go/client"
	"golang.org/x/net/publicsuffix"
)

// @todo add description
// @new
type Caller interface {
	// CallContext calls the RPC method with the specified parameters and
	// returns an error.
	CallContext(ctx context.Context, result interface{}, method string, args interface{}) error
	// @todo remove quotes about Juno, same for the other interfaces
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

// @todo add description
// @new
type Subscriber interface {
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

// @new
func NewClient(
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
