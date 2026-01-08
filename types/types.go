package types

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

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
