package paymaster

import (
	"context"
	"net/http"
	"net/http/cookiejar"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/client"
	"golang.org/x/net/publicsuffix"
)

// Paymaster is a client for interacting with a paymaster service via the SNIP-29 API.
// It provides methods to build and execute transactions, check service status, and track transaction status.
type Paymaster struct {
	// c is the underlying client for the paymaster service.
	c callCloser
}

type PaymasterInterface interface {
	IsAvailable(ctx context.Context) (bool, error)
	GetSupportedTokens(ctx context.Context) ([]TokenData, error)
	TrackingIdToLatestHash(ctx context.Context, trackingId *felt.Felt) (TrackingIdResponse, error)
	BuildTransaction(ctx context.Context, request *BuildTransactionRequest) (*BuildTransactionResponse, error)
	ExecuteTransaction(ctx context.Context, request *ExecuteTransactionRequest) (*ExecuteTransactionResponse, error)
}

var _ PaymasterInterface = &Paymaster{} //nolint:exhaustruct

// callCloser is an interface that defines the methods for calling a remote procedure.
// It was created to match the Client struct from the 'client' package.
type callCloser interface {
	CallContext(ctx context.Context, result interface{}, method string, args interface{}) error
	CallContextWithSliceArgs(ctx context.Context, result interface{}, method string, args ...interface{}) error
	Close()
}

// NewPaymasterClient creates a new paymaster client for the given service URL.
// The client will connect to the paymaster service at the specified URL.
//
// Parameters:
//   - url: The URL of the paymaster service
//
// Returns:
//   - *Paymaster: A new paymaster client instance
//   - error: An error if the client creation fails
//
// NewProvider creates a new HTTP rpc Provider instance.
func NewPaymasterClient(url string, options ...client.ClientOption) (*Paymaster, error) {
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

// Get a list of the tokens that the paymaster supports, together with their prices in STRK
//
// Parameters:
//   - ctx: The context.Context object for controlling the function call
//
// Returns:
//   - []TokenData: An array of token data
//   - error: An error if any
func (p *Paymaster) GetSupportedTokens(ctx context.Context) ([]TokenData, error) {
	var response []TokenData
	if err := p.c.CallContextWithSliceArgs(ctx, &response, "paymaster_getSupportedTokens"); err != nil {
		return nil, err
	}

	return response, nil
}

// TrackingIdToLatestHash gets the latest transaction hash and status for a given tracking ID.
// Returns a TrackingIdResponse.
//
// Parameters:
//   - ctx: The context.Context object for controlling the function call
//   - trackingId: A unique identifier used to track an execution request of a user.
//     This identitifier is returned by the paymaster after a successful call to `execute`.
//     Its purpose is to track the possibly different transaction hashes in the mempool which
//     are associated with a same user request.
//
// Returns:
//   - *TrackingIdResponse: The hash of the latest transaction broadcasted by the paymaster
//     corresponding to the requested ID and the status of the ID.
//   - error: An error if any
func (p *Paymaster) TrackingIdToLatestHash(ctx context.Context, trackingId *felt.Felt) (TrackingIdResponse, error) {
	var response TrackingIdResponse
	if err := p.c.CallContextWithSliceArgs(ctx, &response, "paymaster_trackingIdToLatestHash", trackingId); err != nil {
		return TrackingIdResponse{}, err
	}

	return response, nil
}
