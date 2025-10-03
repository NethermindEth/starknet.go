package paymaster

import (
	"context"
	"net/http"
	"net/http/cookiejar"

	"github.com/NethermindEth/starknet.go/client"
	"golang.org/x/net/publicsuffix"
)

// Paymaster is a client for interacting with a paymaster service via the SNIP-29 API.
// It provides methods to build and execute transactions, check service status, and track transaction status.
type Paymaster struct {
	// c is the underlying client for the paymaster service.
	c callCloser
}

// Used to assert that the Paymaster struct implements all the paymaster methods.
// Ref: https://github.com/starknet-io/SNIPs/blob/ea46a8777d8c8d53a43f45b7beb1abcc301a1a69/assets/snip-29/paymaster_api.json
type paymasterInterface interface {
	IsAvailable(ctx context.Context) (bool, error)
	// More methods coming...
}

var _ paymasterInterface = &Paymaster{} //nolint:exhaustruct

// callCloser is an interface that defines the methods for calling a remote procedure.
// It was created to match the Client struct from the 'client' package.
type callCloser interface {
	CallContext(ctx context.Context, result interface{}, method string, args interface{}) error
	CallContextWithSliceArgs(ctx context.Context, result interface{}, method string, args ...interface{}) error
	Close()
}

// Creates a new paymaster client for the given service URL.
// Additional options can be passed to the client to configure the connection.
//
// Parameters:
//   - url: The URL of the paymaster service
//   - options: Additional options to configure the client
//
// Returns:
//   - *Paymaster: A new paymaster client instance
//   - error: An error if the client creation fails
func New(url string, options ...client.ClientOption) (*Paymaster, error) {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, err
	}
	httpClient := &http.Client{Jar: jar} //nolint:exhaustruct
	// prepend the custom client to allow users to override
	options = append([]client.ClientOption{client.WithHTTPClient(httpClient)}, options...)
	c, err := client.DialOptions(context.Background(), url, options...)
	if err != nil {
		return nil, err
	}

	paymaster := &Paymaster{c: c}

	return paymaster, nil
}
