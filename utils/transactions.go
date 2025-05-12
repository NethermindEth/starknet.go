package utils

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"strings"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/contracts"
	"github.com/NethermindEth/starknet.go/hash"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/NethermindEth/starknet.go/rpc"
)

var (
	maxUint64                 uint64 = math.MaxUint64
	maxUint128                       = "0xffffffffffffffffffffffffffffffff"
	negativeResourceBoundsErr        = "resource bounds cannot be negative, got '%#x'"
	invalidResourceBoundsErr         = "invalid resource bounds: '%v' is not a valid big.Int"
)

// TransactionOptions holds options for building transactions
// Multiplier: safety factor for fee estimation
// WithQueryBitVersion: whether to use the query bit version
// Tip: tip amount for the transaction
//
type TransactionOptions struct {
	Multiplier          float64
	WithQueryBitVersion bool
	Tip                 rpc.U64
}

// validate sets defaults and checks for edge cases
func (opts *TransactionOptions) validate() {
	opts.validateMultiplier()
	opts.validateTip()
}

func (opts *TransactionOptions) validateMultiplier() {
	if opts.Multiplier <= 0 {
		opts.Multiplier = 1.5
	}
}

func (opts *TransactionOptions) validateTip() {
	if opts.Tip == "" {
		opts.Tip = "0x0"
	}
	if !strings.HasPrefix(string(opts.Tip), "0x") {
		opts.Tip = "0x0"
	}
	if val, err := opts.Tip.ToUint64(); err != nil {
		opts.Tip = "0x0"
	} else if val > 0xFFFFFFFFFFFFFFFF {
		opts.Tip = "0xFFFFFFFFFFFFFFFF" // max U64
	}
}

// BuildInvokeTxn creates a new invoke transaction (v3) for the StarkNet network.
//
// The default version of the returned transaction is rpc.TransactionV3 (0x3). If a version with
// rpc.TransactionV3WithQueryBit ('0x100000000000000000000000000000003') is required, it should be set manually
// in the returned transaction.
//
// Parameters:
//   - senderAddress: The address of the account sending the transaction
//   - nonce: The account's nonce
//   - calldata: The data expected by the account's `execute` function (in most usecases,
//     this includes the called contract address and a function selector)
//   - resourceBounds: Resource bounds for the transaction execution
//   - opts: TransactionOptions pointer for tip, multiplier, etc.
//
// Returns:
//   - rpc.BroadcastInvokev3Txn: A broadcast invoke transaction with default values
//     for signature, paymaster data, etc. Needs to be signed before being sent.
func BuildInvokeTxn(
	senderAddress *felt.Felt,
	nonce *felt.Felt,
	calldata []*felt.Felt,
	resourceBounds rpc.ResourceBoundsMapping,
	opts *TransactionOptions,
) *rpc.BroadcastInvokeTxnV3 {
	opts.validate()

	invokeTxn := rpc.BroadcastInvokeTxnV3{
		Type:                  rpc.TransactionType_Invoke,
		SenderAddress:         senderAddress,
		Calldata:              calldata,
		Version:               rpc.TransactionV3,
		Signature:             []*felt.Felt{},
		Nonce:                 nonce,
		ResourceBounds:        resourceBounds,
		Tip:                   opts.Tip,
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
// The default version of the returned transaction is rpc.TransactionV3 (0x3). If a version with
// rpc.TransactionV3WithQueryBit ('0x100000000000000000000000000000003') is required, it should be set manually
// in the returned transaction.
//
// Parameters:
//   - senderAddress: The address of the account sending the transaction
//   - casmClass: The casm class of the contract to be declared
//   - contractClass: The contract class to be declared
//   - nonce: The account's nonce
//   - resourceBounds: Resource bounds for the transaction execution
//   - opts: TransactionOptions pointer for tip, multiplier, etc.
//
// Returns:
//   - rpc.BroadcastDeclareTxnV3: A broadcast declare transaction with default values
//     for signature, paymaster data, etc. Needs to be signed before being sent.
func BuildDeclareTxn(
	senderAddress *felt.Felt,
	casmClass *contracts.CasmClass,
	contractClass *contracts.ContractClass,
	nonce *felt.Felt,
	resourceBounds rpc.ResourceBoundsMapping,
	opts *TransactionOptions,
) (*rpc.BroadcastDeclareTxnV3, error) {
	opts.validate()

	compiledClassHash, err := hash.CompiledClassHash(casmClass)
	if err != nil {
		return nil, err
	}

	declareTxn := rpc.BroadcastDeclareTxnV3{
		Type:                  rpc.TransactionType_Declare,
		SenderAddress:         senderAddress,
		CompiledClassHash:     compiledClassHash,
		Version:               rpc.TransactionV3,
		Signature:             []*felt.Felt{},
		Nonce:                 nonce,
		ContractClass:         contractClass,
		ResourceBounds:        resourceBounds,
		Tip:                   opts.Tip,
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
// The default version of the returned transaction is rpc.TransactionV3 (0x3). If a version with
// rpc.TransactionV3WithQueryBit ('0x100000000000000000000000000000003') is required, it should be set manually
// in the returned transaction.
//
// Parameters:
//   - nonce: The account's nonce
//   - contractAddressSalt: A value used to randomize the deployed contract address
//   - constructorCalldata: The parameters for the constructor function
//   - classHash: The hash of the contract class to deploy
//   - resourceBounds: Resource bounds for the transaction execution
//   - opts: TransactionOptions pointer for tip, multiplier, etc.
//
// Returns:
//   - rpc.BroadcastDeployAccountTxnV3: A broadcast deploy account transaction with default values
//     for signature, paymaster data, etc. Needs to be signed before being sent.
func BuildDeployAccountTxn(
	nonce *felt.Felt,
	contractAddressSalt *felt.Felt,
	constructorCalldata []*felt.Felt,
	classHash *felt.Felt,
	resourceBounds rpc.ResourceBoundsMapping,
	opts *TransactionOptions,
) *rpc.BroadcastDeployAccountTxnV3 {
	opts.validate()

	deployAccountTxn := rpc.BroadcastDeployAccountTxnV3{
		Type:                rpc.TransactionType_DeployAccount,
		Version:             rpc.TransactionV3,
		Signature:           []*felt.Felt{},
		Nonce:               nonce,
		ContractAddressSalt: contractAddressSalt,
		ConstructorCalldata: constructorCalldata,
		ClassHash:           classHash,
		ResourceBounds:      resourceBounds,
		Tip:                 opts.Tip,
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
) rpc.ResourceBoundsMapping {

	// Create L1 resources bounds
	l1Gas := toResourceBounds(feeEstimation.L1GasPrice, feeEstimation.L1GasConsumed, multiplier)
	l1DataGas := toResourceBounds(feeEstimation.L1DataGasPrice, feeEstimation.L1DataGasConsumed, multiplier)
	// Create L2 resource bounds
	l2Gas := toResourceBounds(feeEstimation.L2GasPrice, feeEstimation.L2GasConsumed, multiplier)

	return rpc.ResourceBoundsMapping{
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
	resBounds rpc.ResourceBoundsMapping,
	multiplier float64,
) (*felt.Felt, error) {
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
