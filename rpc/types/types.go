package types

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/NethermindEth/juno/core/felt"
)

// An unsigned integer number in hex format (0x...)
type NumAsHex string

// A storage key, represented as a string of hex digits.
// Represented as up to 62 hex digits, 3 bits, and 5 leading zeroes.
type StorageKey string

// 64 bit unsigned integers, represented by hex string of length at most 16
type U64 string

// ToUint64 converts the U64 type to a uint64.
// If the value is greater than max uint64, returns an error.
func (u U64) ToUint64() (uint64, error) {
	hexStr := strings.TrimPrefix(string(u), "0x")

	val, err := strconv.ParseUint(hexStr, 16, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse hex string: %v", err)
	}

	return val, nil
}

// 128 bit unsigned integers, represented by hex string of length at most 32
type U128 string

// ToBigInt converts the U128 type to a *big.Int.
// If the value is greater than max uint128, returns an error.
//
//nolint:mnd // 16 means hex base
func (u U128) ToBigInt() (*big.Int, error) {
	hexStr := strings.TrimPrefix(string(u), "0x")

	result, ok := new(big.Int).SetString(hexStr, 16)
	if !ok {
		return nil, fmt.Errorf("failed to parse hex string: %v", hexStr)
	}

	maxUint128, _ := new(big.Int).SetString("ffffffffffffffffffffffffffffffff", 16)

	if result.Cmp(maxUint128) > 0 {
		return nil, fmt.Errorf("value is greater than max uint128: %v", u)
	}

	return result, nil
}

// @changed this entire pkg is new and contains migrated code

// @changed moved from the rpcvX pkg
// FunctionCall function call information
type FunctionCall struct {
	ContractAddress    *felt.Felt `json:"contract_address"`
	EntryPointSelector *felt.Felt `json:"entry_point_selector"`

	// Calldata The parameters passed to the function
	Calldata []*felt.Felt `json:"calldata"`
}

// @changed moved from the rpcvX pkg
// InvokeFunctionCall represents a function call to be invoked on a contract.
// It's a helper type used to build a FunctionCall for a v3 Invoke transaction.
type InvokeFunctionCall struct {
	// The address of the contract to invoke
	ContractAddress *felt.Felt
	// The name of the function to invoke
	FunctionName string
	// The parameters passed to the function
	CallData []*felt.Felt
}

// FeeLimits is a struct with custom limits for the fee values, used
// as a parameter for the `CustomFeeEstToResBoundsMap` function.
type FeeLimits struct {
	// Custom max value for L1 gas price
	L1GasPriceLimit U128
	// Custom max value for L1 gas amount
	L1GasAmountLimit U64

	// Custom max value for L2 gas price
	L2GasPriceLimit U128
	// Custom max value for L2 gas amount
	L2GasAmountLimit U64

	// Custom max value for L1 data gas price
	L1DataGasPriceLimit U128
	// Custom max value for L1 data gas amount
	L1DataGasAmountLimit U64
}
