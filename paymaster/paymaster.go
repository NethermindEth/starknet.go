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

// callCloser is an interface that defines the methods for calling a remote procedure.
// It was created to match the Client struct from the 'client' package.
type callCloser interface {
	// CallContextWithSliceArgs call 'CallContext' with a slice of arguments.
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

// IsAvailable returns the status of the paymaster service.
// If the paymaster service is correctly functioning, return true. Else, return false
//
// Parameters:
//   - ctx: The context.Context object for controlling the function call
//
// Returns:
//   - bool: True if the paymaster service is correctly functioning, false otherwise
//   - error: An error if any
func (p *Paymaster) IsAvailable(ctx context.Context) (bool, error) {
	var response bool
	if err := p.c.CallContextWithSliceArgs(ctx, &response, "paymaster_isAvailable"); err != nil {
		return false, err
	}

	return response, nil
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

// BuildTransaction receives the transaction the user wants to execute. Returns the typed
// data along with the estimated gas cost and the maximum gas cost suggested to ensure execution
//
// Parameters:
//   - ctx: The context.Context object for controlling the function call
//   - request: The BuildTransactionRequest containing the transaction and parameters
//
// Returns:
//   - *BuildTransactionResponse: The response containing typed data and fee estimate
//   - error: An error if the request fails
func (p *Paymaster) BuildTransaction(ctx context.Context, request *BuildTransactionRequest) (*BuildTransactionResponse, error) {
	var response BuildTransactionResponse
	if err := p.c.CallContextWithSliceArgs(ctx, &response, "paymaster_buildTransaction", request); err != nil {
		return nil, err
	}

	return &response, nil
}

// ExecuteTransaction executes a signed transaction via the paymaster service.
// Returns an ExecuteTransactionResponse with tracking ID and transaction hash.
//
// Parameters:
//   - ctx: The context.Context object for controlling the function call
//   - request: The ExecuteTransactionRequest containing the transaction and parameters
//
// Returns:
//   - *ExecuteTransactionResponse: The response containing tracking ID and transaction hash
//   - error: An error if the execution fails
func (p *Paymaster) ExecuteTransaction(ctx context.Context, request *ExecuteTransactionRequest) (*ExecuteTransactionResponse, error) {
	var response ExecuteTransactionResponse
	if err := p.c.CallContextWithSliceArgs(ctx, &response, "paymaster_executeTransaction", request); err != nil {
		return nil, err
	}

	return &response, nil
}

// Constants for the SNIP-X specification
const (
	// OUTSIDE_EXECUTION_TYPED_DATA_V1 represents the typed data structure for version 1
	OUTSIDE_EXECUTION_TYPED_DATA_V1 = "OUTSIDE_EXECUTION_TYPED_DATA_V1"

	// OUTSIDE_EXECUTION_TYPED_DATA_V2 represents the typed data structure for version 2
	OUTSIDE_EXECUTION_TYPED_DATA_V2 = "OUTSIDE_EXECUTION_TYPED_DATA_V2"

	// OUTSIDE_EXECUTION_TYPED_DATA_V3_RC represents the typed data structure for version 3-rc
	OUTSIDE_EXECUTION_TYPED_DATA_V3_RC = "OUTSIDE_EXECUTION_TYPED_DATA_V3_RC"
)

// GetOutsideExecutionTypedDataV1 returns the typed data structure for version 1
func GetOutsideExecutionTypedDataV1(message OutsideExecutionMessageV1) OutsideExecutionTypedData {
	return OutsideExecutionTypedData{
		Types: map[string][]TypedDataField{
			"StarkNetDomain": {
				{Name: "name", Type: "felt"},
				{Name: "version", Type: "felt"},
				{Name: "chainId", Type: "felt"},
			},
			"OutsideExecution": {
				{Name: "caller", Type: "felt"},
				{Name: "nonce", Type: "felt"},
				{Name: "execute_after", Type: "felt"},
				{Name: "execute_before", Type: "felt"},
				{Name: "calls_len", Type: "felt"},
				{Name: "calls", Type: "OutsideCall*"},
			},
			"OutsideCall": {
				{Name: "to", Type: "felt"},
				{Name: "selector", Type: "felt"},
				{Name: "calldata_len", Type: "felt"},
				{Name: "calldata", Type: "felt*"},
			},
		},
		PrimaryType: "OutsideExecution",
		Domain: TypedDataDomain{
			Name:    "Account.execute_from_outside",
			Version: "1",
			ChainID: "0x534e5f4d41494e", // SN_MAINNET
		},
		Message: message,
	}
}

// GetOutsideExecutionTypedDataV2 returns the typed data structure for version 2
func GetOutsideExecutionTypedDataV2(message OutsideExecutionMessageV2) OutsideExecutionTypedData {
	return OutsideExecutionTypedData{
		Types: map[string][]TypedDataField{
			"StarknetDomain": {
				{Name: "name", Type: "shortstring"},
				{Name: "version", Type: "shortstring"},
				{Name: "chainId", Type: "shortstring"},
				{Name: "revision", Type: "shortstring"},
			},
			"OutsideExecution": {
				{Name: "Caller", Type: "ContractAddress"},
				{Name: "Nonce", Type: "felt"},
				{Name: "Execute After", Type: "u128"},
				{Name: "Execute Before", Type: "u128"},
				{Name: "Calls", Type: "Call*"},
			},
			"Call": {
				{Name: "To", Type: "ContractAddress"},
				{Name: "Selector", Type: "selector"},
				{Name: "Calldata", Type: "felt*"},
			},
		},
		PrimaryType: "OutsideExecution",
		Domain: TypedDataDomain{
			Name:    "Account.execute_from_outside",
			Version: "2",
			ChainID: "0x534e5f4d41494e", // SN_MAINNET
		},
		Message: message,
	}
}

// GetOutsideExecutionTypedDataV3RC returns the typed data structure for version 3-rc
func GetOutsideExecutionTypedDataV3RC(message *OutsideExecutionMessageV3) OutsideExecutionTypedData {
	return OutsideExecutionTypedData{
		Types: map[string][]TypedDataField{
			"StarknetDomain": {
				{Name: "name", Type: "shortstring"},
				{Name: "version", Type: "shortstring"},
				{Name: "chainId", Type: "shortstring"},
				{Name: "revision", Type: "shortstring"},
			},
			"OutsideExecution": {
				{Name: "Caller", Type: "ContractAddress"},
				{Name: "Nonce", Type: "felt"},
				{Name: "Execute After", Type: "u128"},
				{Name: "Execute Before", Type: "u128"},
				{Name: "Calls", Type: "Call*"},
				{Name: "Fee", Type: "Fee Mode"},
			},
			"Call": {
				{Name: "To", Type: "ContractAddress"},
				{Name: "Selector", Type: "selector"},
				{Name: "Calldata", Type: "felt*"},
			},
			"Fee Mode": {
				{Name: "No Fee", Type: "()"},
				{Name: "Pay Fee", Type: "(FeeTransfer)"},
			},
			"Fee Transfer": {
				{Name: "Fee Amount", Type: "TokenAmount"},
				{Name: "Fee Receiver", Type: "ContractAddress"},
			},
		},
		PrimaryType: "OutsideExecution",
		Domain: TypedDataDomain{
			Name:    "Account.execute_from_outside",
			Version: "3",
			ChainID: "0x534e5f4d41494e", // SN_MAINNET
		},
		Message: message,
	}
}
