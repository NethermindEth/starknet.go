package utils

import (
	"errors"
	"fmt"
	"math"
	"math/big"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/contracts"
	"github.com/NethermindEth/starknet.go/hash"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/NethermindEth/starknet.go/rpc"
)

const (
	maxUint64  uint64 = math.MaxUint64
	maxUint128        = "0xffffffffffffffffffffffffffffffff"

	negativeResourceBoundsErr = "resource bounds cannot be negative, got '%#x'"
	invalidResourceBoundsErr  = "invalid resource bounds: '%v' is not a valid big.Int"
)

// Optional settings when building a transaction.
type TxnOptions struct {
	// Tip amount in FRI for the transaction. Default: `"0x0"`.
	// Note: only ready to be used after Starknet v0.14.0 upgrade.
	Tip rpc.U64
	// A boolean flag indicating whether the transaction version should have
	// the query bit when estimating fees. If true, the transaction version
	// will be `rpc.TransactionV3WithQueryBit` (0x100000000000000000000000000000003).
	// If false, the transaction version will be `rpc.TransactionV3` (0x3).
	// In case of doubt, set to `false`. Default: `false`.
	UseQueryBit bool

	// ONLY FOR DECLARE TXN: A boolean flag indicating whether to use the Blake2s hash
	// function to calculate the compiled class hash. This must be set to true after
	// the Starknet v0.14.1 upgrade, when the Poseidon hash function will be deprecated
	// for the compiled class hash.
	UseBlake2sHash bool
}

// TxnVersion returns `rpc.TransactionV3WithQueryBit` when UseQueryBit is true, and
// `rpc.TransactionV3` if false.
func (opts *TxnOptions) TxnVersion() rpc.TransactionVersion {
	if opts.UseQueryBit {
		return rpc.TransactionV3WithQueryBit
	}

	return rpc.TransactionV3
}

// SafeTip returns the tip amount in FRI for the transaction. If the tip is not set or invalid, returns "0x0".
func (opts *TxnOptions) SafeTip() rpc.U64 {
	if opts.Tip == "" {
		return "0x0"
	}

	if _, err := opts.Tip.ToUint64(); err != nil {
		return "0x0"
	}

	return opts.Tip
}

// BuildInvokeTxn creates a new invoke transaction (v3) for the StarkNet network.
//
// Parameters:
//   - senderAddress: The address of the account sending the transaction
//   - nonce: The account's nonce
//   - calldata: The data expected by the account's `execute` function (in most usecases,
//     this includes the called contract address and a function selector)
//   - resourceBounds: Resource bounds for the transaction execution
//   - opts: optional settings for the transaction
//
// Returns:
//   - rpc.BroadcastInvokev3Txn: A broadcast invoke transaction with default values
//     for signature, paymaster data, etc. Needs to be signed before being sent.
func BuildInvokeTxn(
	senderAddress *felt.Felt,
	nonce *felt.Felt,
	calldata []*felt.Felt,
	resourceBounds *rpc.ResourceBoundsMapping,
	opts *TxnOptions,
) *rpc.BroadcastInvokeTxnV3 {
	if opts == nil {
		opts = new(TxnOptions)
	}

	invokeTxn := rpc.BroadcastInvokeTxnV3{
		Type:                  rpc.TransactionType_Invoke,
		SenderAddress:         senderAddress,
		Calldata:              calldata,
		Version:               opts.TxnVersion(),
		Signature:             []*felt.Felt{},
		Nonce:                 nonce,
		ResourceBounds:        resourceBounds,
		Tip:                   opts.SafeTip(),
		PayMasterData:         []*felt.Felt{},
		AccountDeploymentData: []*felt.Felt{},
		NonceDataMode:         rpc.DAModeL1,
		FeeMode:               rpc.DAModeL1,
	}

	return &invokeTxn
}

// BuildDeclareTxn creates a new declare transaction (v3) for the StarkNet network.
// A declare transaction is used to declare a new contract class on the network.
//
// Parameters:
//   - senderAddress: The address of the account sending the transaction
//   - casmClass: The casm class of the contract to be declared
//   - contractClass: The contract class to be declared
//   - nonce: The account's nonce
//   - resourceBounds: Resource bounds for the transaction execution
//   - opts: optional settings for the transaction
//
// Returns:
//   - rpc.BroadcastDeclareTxnV3: A broadcast declare transaction with default values
//     for signature, paymaster data, etc. Needs to be signed before being sent.
func BuildDeclareTxn(
	senderAddress *felt.Felt,
	casmClass *contracts.CasmClass,
	contractClass *contracts.ContractClass,
	nonce *felt.Felt,
	resourceBounds *rpc.ResourceBoundsMapping,
	opts *TxnOptions,
) (*rpc.BroadcastDeclareTxnV3, error) {
	if opts == nil {
		opts = new(TxnOptions)
	}

	var compiledClassHash *felt.Felt
	var err error

	if opts.UseBlake2sHash {
		compiledClassHash, err = hash.CompiledClassHashV2(casmClass)
		if err != nil {
			return nil, err
		}
	} else {
		compiledClassHash, err = hash.CompiledClassHash(casmClass)
		if err != nil {
			return nil, err
		}
	}

	declareTxn := rpc.BroadcastDeclareTxnV3{
		Type:                  rpc.TransactionType_Declare,
		SenderAddress:         senderAddress,
		CompiledClassHash:     compiledClassHash,
		Version:               opts.TxnVersion(),
		Signature:             []*felt.Felt{},
		Nonce:                 nonce,
		ContractClass:         contractClass,
		ResourceBounds:        resourceBounds,
		Tip:                   opts.SafeTip(),
		PayMasterData:         []*felt.Felt{},
		AccountDeploymentData: []*felt.Felt{},
		NonceDataMode:         rpc.DAModeL1,
		FeeMode:               rpc.DAModeL1,
	}

	return &declareTxn, nil
}

// BuildDeployAccountTxn creates a new deploy account transaction (v3) for the StarkNet network.
// A deploy account transaction is used to deploy a new account contract on the network.
//
// Parameters:
//   - nonce: The account's nonce
//   - contractAddressSalt: A value used to randomise the deployed contract address
//   - constructorCalldata: The parameters for the constructor function
//   - classHash: The hash of the contract class to deploy
//   - resourceBounds: Resource bounds for the transaction execution
//   - opts: optional settings for the transaction
//
// Returns:
//   - rpc.BroadcastDeployAccountTxnV3: A broadcast deploy account transaction with default values
//     for signature, paymaster data, etc. Needs to be signed before being sent.
func BuildDeployAccountTxn(
	nonce *felt.Felt,
	contractAddressSalt *felt.Felt,
	constructorCalldata []*felt.Felt,
	classHash *felt.Felt,
	resourceBounds *rpc.ResourceBoundsMapping,
	opts *TxnOptions,
) *rpc.BroadcastDeployAccountTxnV3 {
	if opts == nil {
		opts = new(TxnOptions)
	}

	deployAccountTxn := rpc.BroadcastDeployAccountTxnV3{
		Type:                rpc.TransactionType_DeployAccount,
		Version:             opts.TxnVersion(),
		Signature:           []*felt.Felt{},
		Nonce:               nonce,
		ContractAddressSalt: contractAddressSalt,
		ConstructorCalldata: constructorCalldata,
		ClassHash:           classHash,
		ResourceBounds:      resourceBounds,
		Tip:                 opts.SafeTip(),
		PayMasterData:       []*felt.Felt{},
		NonceDataMode:       rpc.DAModeL1,
		FeeMode:             rpc.DAModeL1,
	}

	return &deployAccountTxn
}

// InvokeFuncCallsToFunctionCalls converts a slice of InvokeFunctionCall to a slice of FunctionCall.
//
// Parameters:
//   - invokeFuncCalls: The invoke function calls to convert
//
// Returns:
//   - []*rpc.FunctionCall: A new function calls
func InvokeFuncCallsToFunctionCalls(invokeFuncCalls []rpc.InvokeFunctionCall) []rpc.FunctionCall {
	functionCalls := make([]rpc.FunctionCall, len(invokeFuncCalls))
	for i, call := range invokeFuncCalls {
		functionCalls[i] = rpc.FunctionCall{
			ContractAddress:    call.ContractAddress,
			EntryPointSelector: GetSelectorFromNameFelt(call.FunctionName),
			Calldata:           call.CallData,
		}
	}

	return functionCalls
}

// FeeEstToResBoundsMap converts a FeeEstimation to ResourceBoundsMapping with applied multipliers.
// Parameters:
//   - feeEstimation: The fee estimation to convert
//   - multiplier: Multiplier for max amount and max price per unit. Recommended to be 1.5, but at least greater than 0.
//     If multiplier < 0, all resources bounds will be set to 0.
//     If resource bounds overflow, they will be set to the max allowed value (U64 or U128).
//
// Returns:
//   - rpc.ResourceBoundsMapping: Resource bounds with applied multipliers
func FeeEstToResBoundsMap(
	feeEstimation rpc.FeeEstimation,
	multiplier float64,
) *rpc.ResourceBoundsMapping {
	// Create L1 resources bounds
	l1Gas := toResourceBounds(feeEstimation.L1GasPrice, feeEstimation.L1GasConsumed, multiplier)
	l1DataGas := toResourceBounds(feeEstimation.L1DataGasPrice, feeEstimation.L1DataGasConsumed, multiplier)
	// Create L2 resource bounds
	l2Gas := toResourceBounds(feeEstimation.L2GasPrice, feeEstimation.L2GasConsumed, multiplier)

	return &rpc.ResourceBoundsMapping{
		L1Gas:     l1Gas,
		L1DataGas: l1DataGas,
		L2Gas:     l2Gas,
	}
}

// toResourceBounds converts a gas price and gas consumed to a ResourceBounds with applied multiplier.
//
// Parameters:
//   - gasPrice: The gas price
//   - gasConsumed: The gas consumed
//   - multiplier: Multiplier for max amount and max price per unit
//
// Returns:
//   - rpc.ResourceBounds: Resource bounds with applied multiplier
func toResourceBounds(
	gasPrice *felt.Felt,
	gasConsumed *felt.Felt,
	multiplier float64,
) rpc.ResourceBounds {
	// negative multiplier is not allowed, default to 0
	if multiplier < 0 {
		return rpc.ResourceBounds{
			MaxAmount:       rpc.U64("0x0"),
			MaxPricePerUnit: rpc.U128("0x0"),
		}
	}

	// Convert felt to big.Int
	gasPriceInt := gasPrice.BigInt(new(big.Int))
	gasConsumedInt := gasConsumed.BigInt(new(big.Int))

	// Check for overflow
	maxUint64 := new(big.Int).SetUint64(maxUint64)
	maxUint128, _ := new(big.Int).SetString(maxUint128, 0)
	// max_price_per_unit is U128 by the spec
	if gasPriceInt.Cmp(maxUint128) > 0 {
		gasPriceInt = maxUint128
	}
	// max_amount is U64 by the spec
	if gasConsumedInt.Cmp(maxUint64) > 0 {
		gasConsumedInt = maxUint64
	}

	maxAmount := new(big.Float)
	maxPricePerUnit := new(big.Float)

	maxAmount.Mul(new(big.Float).SetInt(gasConsumedInt), big.NewFloat(multiplier))
	maxPricePerUnit.Mul(new(big.Float).SetInt(gasPriceInt), big.NewFloat(multiplier))

	// Convert big.Float to big.Int for proper hex formatting. The result is a truncated int
	maxAmountInt, _ := maxAmount.Int(new(big.Int))
	maxPricePerUnitInt, _ := maxPricePerUnit.Int(new(big.Int))

	// Check for overflow after mul operation
	if maxAmountInt.Cmp(maxUint64) > 0 {
		maxAmountInt = maxUint64
	}
	if maxPricePerUnitInt.Cmp(maxUint128) > 0 {
		maxPricePerUnitInt = maxUint128
	}

	return rpc.ResourceBounds{
		MaxAmount:       rpc.U64(fmt.Sprintf("%#x", maxAmountInt)),
		MaxPricePerUnit: rpc.U128(fmt.Sprintf("%#x", maxPricePerUnitInt)),
	}
}

// ResBoundsMapToOverallFee calculates the overall fee for a ResourceBoundsMapping with applied multipliers.
// Parameters:
//   - resBounds: The resource bounds to calculate the fee for
//   - multiplier: Multiplier for max amount and max price per unit. Recommended to be 1.5, but at least greater than 0
//
// Returns:
//   - *felt.Felt: The overall fee in FRI
//   - error: An error if any
func ResBoundsMapToOverallFee(
	resBounds *rpc.ResourceBoundsMapping,
	multiplier float64,
) (*felt.Felt, error) {
	if resBounds == nil {
		return nil, errors.New("resource bounds are nil")
	}

	// negative multiplier is not allowed
	if multiplier < 0 {
		return nil, errors.New("multiplier cannot be negative")
	}

	parseBound := func(value string) (*big.Int, error) {
		// get big int values
		val, ok := new(big.Int).SetString(value, 0)
		if !ok {
			return nil, fmt.Errorf(invalidResourceBoundsErr, value)
		}
		// Check for negative values
		if val.Sign() < 0 {
			return nil, fmt.Errorf(negativeResourceBoundsErr, val)
		}

		return val, nil
	}

	l1GasAmount, err := parseBound(string(resBounds.L1Gas.MaxAmount))
	if err != nil {
		return nil, err
	}

	l1GasPrice, err := parseBound(string(resBounds.L1Gas.MaxPricePerUnit))
	if err != nil {
		return nil, err
	}

	l1DataGasAmount, err := parseBound(string(resBounds.L1DataGas.MaxAmount))
	if err != nil {
		return nil, err
	}

	l1DataGasPrice, err := parseBound(string(resBounds.L1DataGas.MaxPricePerUnit))
	if err != nil {
		return nil, err
	}

	l2GasAmount, err := parseBound(string(resBounds.L2Gas.MaxAmount))
	if err != nil {
		return nil, err
	}

	l2GasPrice, err := parseBound(string(resBounds.L2Gas.MaxPricePerUnit))
	if err != nil {
		return nil, err
	}

	// calculate fee
	l1GasFee := new(big.Int).Mul(l1GasAmount, l1GasPrice)
	l1DataGasFee := new(big.Int).Mul(l1DataGasAmount, l1DataGasPrice)
	l2GasFee := new(big.Int).Mul(l2GasAmount, l2GasPrice)
	overallFee := l1GasFee.Add(l1GasFee, l1DataGasFee).Add(l1GasFee, l2GasFee)

	// multiply fee by multiplier
	multipliedOverallFee := new(big.Float).Mul(new(big.Float).SetInt(overallFee), big.NewFloat(multiplier))
	overallFeeInt, _ := multipliedOverallFee.Int(nil) // truncated int

	// Convert big.Int to felt. SetString() validates if it's a valid felt
	return new(felt.Felt).SetString(fmt.Sprintf("%#x", overallFeeInt))
}

// FillHexWithZeroes normalises a hex string to have a '0x' prefix and pads it with leading zeros
// to a total length of 66 characters (including the '0x' prefix).
func FillHexWithZeroes(hex string) string {
	return internalUtils.FillHexWithZeroes(hex)
}

// WeiToETH converts a Wei amount to ETH
// Returns the ETH value as a float64
func WeiToETH(wei *felt.Felt) float64 {
	return internalUtils.WeiToETH(wei)
}

// ETHToWei converts an ETH amount to Wei
// Returns the Wei value as a *felt.Felt
func ETHToWei(eth float64) *felt.Felt {
	return internalUtils.ETHToWei(eth)
}

// FRIToSTRK converts a FRI amount to STRK
// Returns the STRK value as a float64
func FRIToSTRK(fri *felt.Felt) float64 {
	return internalUtils.WeiToETH(fri)
}

// STRKToFRI converts a STRK amount to FRI
// Returns the FRI value as a *felt.Felt
func STRKToFRI(strk float64) *felt.Felt {
	return internalUtils.ETHToWei(strk)
}
