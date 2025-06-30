package rpc

import (
	"context"

	"github.com/NethermindEth/juno/core/felt"
)

// Package rpc provides the RPC client implementation for Starknet.
// This file contains the paymaster API implementation based on SNIP-29 specification.

// PaymasterClient is a client for interacting with a paymaster service via the SNIP-29 API.
// It provides methods to build and execute transactions, check service status, and track transaction status.
type PaymasterClient struct {
	c callCloser
}

// NewPaymasterClient creates a new paymaster client for the given service URL.
// The client will connect to the paymaster service at the specified URL.
//
// Parameters:
//   - url: The URL of the paymaster service
//
// Returns:
//   - *PaymasterClient: A new paymaster client instance
//   - error: An error if the client creation fails
func NewPaymasterClient(url string) (*PaymasterClient, error) {
	// For now, we'll use the same client creation pattern as Provider
	// In a real implementation, this would connect to a paymaster service
	provider, err := NewProvider(url)
	if err != nil {
		return nil, err
	}

	return &PaymasterClient{
		c: provider.c,
	}, nil
}

// IsAvailable checks if the paymaster service is up and running.
// Returns true if available, false otherwise.
//
// Parameters:
//   - ctx: The context.Context object for controlling the function call
//
// Returns:
//   - bool: True if the paymaster service is correctly functioning, false otherwise
//   - error: An error if the request fails
func (pc *PaymasterClient) IsAvailable(ctx context.Context) (bool, error) {
	var result bool
	if err := do(ctx, pc.c, "paymaster_isAvailable", &result); err != nil {
		return false, err
	}
	return result, nil
}

// GetSupportedTokens gets a list of tokens supported by the paymaster service.
// Returns an array of TokenData.
//
// Parameters:
//   - ctx: The context.Context object for controlling the function call
//
// Returns:
//   - []TokenData: An array of token data
//   - error: An error if the request fails
func (pc *PaymasterClient) GetSupportedTokens(ctx context.Context) ([]TokenData, error) {
	var result []TokenData
	if err := do(ctx, pc.c, "paymaster_getSupportedTokens", &result); err != nil {
		return nil, err
	}
	return result, nil
}

// TrackingIdToLatestHash gets the latest transaction hash and status for a given tracking ID.
// Returns a TrackingIdResponse.
//
// Parameters:
//   - ctx: The context.Context object for controlling the function call
//   - trackingId: A unique identifier used to track an execution request
//
// Returns:
//   - *TrackingIdResponse: The response containing transaction hash and status
//   - error: An error if the request fails
func (pc *PaymasterClient) TrackingIdToLatestHash(ctx context.Context, trackingId *felt.Felt) (*TrackingIdResponse, error) {
	var result TrackingIdResponse
	if err := do(ctx, pc.c, "paymaster_trackingIdToLatestHash", &result, trackingId); err != nil {
		return nil, err
	}
	return &result, nil
}

// BuildTransaction builds a transaction, returning typed data for signature and a fee estimate.
// Returns a BuildTransactionResponse.
//
// Parameters:
//   - ctx: The context.Context object for controlling the function call
//   - request: The BuildTransactionRequest containing the transaction and parameters
//
// Returns:
//   - *BuildTransactionResponse: The response containing typed data and fee estimate
//   - error: An error if the request fails
func (pc *PaymasterClient) BuildTransaction(ctx context.Context, request BuildTransactionRequest) (*BuildTransactionResponse, error) {
	var response BuildTransactionResponse
	if err := do(ctx, pc.c, "paymaster_buildTransaction", &response, request); err != nil {
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
func (pc *PaymasterClient) ExecuteTransaction(ctx context.Context, request ExecuteTransactionRequest) (*ExecuteTransactionResponse, error) {
	var response ExecuteTransactionResponse
	if err := do(ctx, pc.c, "paymaster_executeTransaction", &response, request); err != nil {
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
func GetOutsideExecutionTypedDataV3RC(message OutsideExecutionMessageV3) OutsideExecutionTypedData {
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
