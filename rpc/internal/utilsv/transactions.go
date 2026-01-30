package utilsv

import (
	"errors"
	"fmt"
	"math"
	"math/big"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/contracts"
	"github.com/NethermindEth/starknet.go/hash"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/NethermindEth/starknet.go/rpc/rpcv10"
	"github.com/NethermindEth/starknet.go/rpc/types"
	"github.com/NethermindEth/starknet.go/types/constraints"
)

const (
	// Ref: https://docs.starknet.io/learn/cheatsheets/chain-info#current-limits
	maxL2GasAmount = "0x3b9aca00" // = 10^9 = 1_000_000_000

	maxUint64  = "0xffffffffffffffff"
	maxUint128 = "0xffffffffffffffffffffffffffffffff"

	negativeResourceBoundsErr = "resource bounds cannot be negative, got '%#x'"
	invalidResourceBoundsErr  = "invalid resource bounds: '%v' is not a valid big.Int"
)

// Default fee limits for the Starknet network.
// Since there's no official limit for most resources, we use the max allowed values
// of the corresponding types defined by the Starknet specification¹ (uint64 or uint128).
// The L2 gas amount however is specified², so we use this value.
//
// ¹ Ref: https://github.com/starkware-libs/starknet-specs/blob/6485866d8b017f2dd615ee245275833028464419/api/starknet_api_openrpc.json#L3508
// ² Ref: https://docs.starknet.io/learn/cheatsheets/chain-info#current-limits
//
//nolint:lll // The link would be unclickable if we break the line.
var StarknetLimits = types.FeeLimits{
	L1GasPriceLimit:      maxUint128,
	L1GasAmountLimit:     maxUint64,
	L1DataGasPriceLimit:  maxUint128,
	L1DataGasAmountLimit: maxUint64,
	L2GasPriceLimit:      maxUint128,
	L2GasAmountLimit:     maxL2GasAmount,
}

type TxnOptions[TxVersion ~string] interface {
	TxnVersion() TxVersion
	SafeTip() types.U64
}

// BuildDeclareTxn creates a broadcast declare transaction (v3) by accepting a pointer
// to it and filling it with the given parameters and default values. It also
// returns the pointer to the filled transaction.
//
// Parameters:
//   - tx: A pointer to the desired broadcast declare transaction (e.g.
//     rpcv9.BroadcastDeclareTxnV3, rpcv10.BroadcastDeclareTxnV3, etc.) to be filled.
//     It needs to be signed before being sent
//   - senderAddress: The address of the account sending the transaction
//   - casmClass: The casm class of the contract to be declared
//   - contractClass: The contract class to be declared
//   - nonce: The account's nonce
//   - resourceBounds: Resource bounds for the transaction execution
//   - opts: optional settings for the transaction
//
// Returns:
//   - *BDeclareTxn: The pointer to the filled broadcast declare transaction. It needs to be signed
//     before being sent.
func BuildDeclareTxn[
	TransactionType, TransactionVersion ~string,
	u64 constraints.U64,
	u128 constraints.U128,
	RB constraints.ResourceBounds[u64, u128],
	RBM constraints.ResourceBoundsMapping[u64, u128, RB],
	DA constraints.DataAvailabilityMode,
	BDeclareTxn constraints.BroadcastDeclareTxnV3[
		TransactionType, TransactionVersion, u64, u128, RB, RBM, DA],
	Opts TxnOptions[TransactionVersion],
](
	tx *BDeclareTxn,
	senderAddress *felt.Felt,
	casmClass *contracts.CasmClass,
	contractClass *contracts.ContractClass,
	nonce *felt.Felt,
	resourceBounds *RBM,
	opts Opts,
) (*BDeclareTxn, error) {

	compiledClassHash, err := hash.CompiledClassHashV2(casmClass)
	if err != nil {
		return nil, err
	}

	*tx = BDeclareTxn{
		Type:                  TransactionType(rpcv10.TransactionTypeDeclare),
		SenderAddress:         senderAddress,
		CompiledClassHash:     compiledClassHash,
		Version:               TransactionVersion(opts.TxnVersion()),
		Signature:             []*felt.Felt{},
		Nonce:                 nonce,
		ContractClass:         contractClass,
		ResourceBounds:        resourceBounds,
		Tip:                   u64(opts.SafeTip()),
		PayMasterData:         []*felt.Felt{},
		AccountDeploymentData: []*felt.Felt{},
		NonceDataMode:         DA(types.DAModeL1),
		FeeMode:               DA(types.DAModeL1),
	}

	return tx, nil
}

// BuildDeployAccountTxn creates a broadcast deploy account transaction (v3) by
// accepting a pointer to it and filling it with the given parameters and default
// values. It also returns the pointer to the filled transaction.
//
// Parameters:
//   - tx: A pointer to the desired broadcast deploy account transaction (e.g.
//     rpcv9.BroadcastDeployAccountTxnV3, rpcv10.BroadcastDeployAccountTxnV3, etc.) to
//     be filled. It needs to be signed before being sent
//   - nonce: The account's nonce
//   - contractAddressSalt: A value used to randomise the deployed contract address
//   - constructorCalldata: The parameters for the constructor function
//   - classHash: The hash of the contract class to deploy
//   - resourceBounds: Resource bounds for the transaction execution
//   - opts: optional settings for the transaction
//
// Returns:
//   - *BDeployAccountTxn: The pointer to the filled broadcast deploy account transaction.
//     It needs to be signed before being sent.
func BuildDeployAccountTxn[
	TransactionType, TransactionVersion ~string,
	u64 constraints.U64,
	u128 constraints.U128,
	RB constraints.ResourceBounds[u64, u128],
	RBM constraints.ResourceBoundsMapping[u64, u128, RB],
	DA constraints.DataAvailabilityMode,
	BDeployAccountTxn constraints.DeployAccountTxnV3[
		TransactionType, TransactionVersion, u64, u128, RB, RBM, DA],
	Opts TxnOptions[TransactionVersion],
](
	tx *BDeployAccountTxn,
	nonce *felt.Felt,
	contractAddressSalt *felt.Felt,
	constructorCalldata []*felt.Felt,
	classHash *felt.Felt,
	resourceBounds *RBM,
	opts Opts,
) *BDeployAccountTxn {
	*tx = BDeployAccountTxn{
		Type:                TransactionType(rpcv10.TransactionTypeDeployAccount),
		Version:             TransactionVersion(opts.TxnVersion()),
		Signature:           []*felt.Felt{},
		Nonce:               nonce,
		ContractAddressSalt: contractAddressSalt,
		ConstructorCalldata: constructorCalldata,
		ClassHash:           classHash,
		ResourceBounds:      resourceBounds,
		Tip:                 u64(opts.SafeTip()),
		PayMasterData:       []*felt.Felt{},
		NonceDataMode:       DA(types.DAModeL1),
		FeeMode:             DA(types.DAModeL1),
	}

	return tx
}

// @changed
// FeeEstToResBoundsMap converts a FeeEstimation to a ResourceBoundsMapping with
// applied multipliers.
// Parameters:
//   - feeEstimation: The fee estimation to convert
//   - resources: a pointer to the resource bounds mapping to fill
//     (e.g. *rpcv9.ResourceBoundsMapping, *rpcv10.ResourceBoundsMapping, etc.)
//   - multiplier: Multiplier for max amount and max price per unit. Recommended to be 1.5,
//     but at least greater than 0.
//     If multiplier <= 0, all resources bounds will be set to 0.
//     If resource bounds overflow, they will be set to the max allowed value (U64 or U128).
//
// Returns:
//   - *RBM: The pointer to the filled resource bounds mapping.
func FeeEstToResBoundsMap[
	PriceUnit ~string,
	FE constraints.FeeEstimation[PriceUnit],
	u64 constraints.U64,
	u128 constraints.U128,
	RB constraints.ResourceBounds[u64, u128],
	RBM constraints.ResourceBoundsMapping[u64, u128, RB],
](
	feeEstimation *FE,
	resources *RBM,
	multiplier float64,
) *RBM {
	bounds := CustomFeeEstToResBoundsMap(
		feeEstimation,
		resources,
		multiplier,
		&StarknetLimits)

	// TODO: return by value instead of pointer
	return bounds
}

// @changed
// CustomFeeEstToResBoundsMap converts a FeeEstimation to ResourceBoundsMapping with applied
// multipliers and limits.
// Parameters:
//   - feeEstimation: The fee estimation to convert
//   - resources: a pointer to the resource bounds mapping to fill
//   - multiplier: Multiplier for max amount and max price per unit. Recommended to be 1.5,
//     but at least greater than 0.
//     If multiplier <= 0, all resources bounds will be set to 0.
//     If resource bounds overflow, they will be set to the max allowed value (U64 or U128).
//   - limits: Limits for the resource bounds. They are still capped to the max allowed
//     values (U64 or U128).
//
// Returns:
//   - *RBM: The pointer to the filled resource bounds mapping.
func CustomFeeEstToResBoundsMap[
	PriceUnit ~string,
	FE constraints.FeeEstimation[PriceUnit],
	u64 constraints.U64,
	u128 constraints.U128,
	RB constraints.ResourceBounds[u64, u128],
	RBM constraints.ResourceBoundsMapping[u64, u128, RB],
](
	feeEstimation *FE,
	resources *RBM,
	multiplier float64,
	limits *types.FeeLimits,
) *RBM {
	innerFE := constraints.FeeEstimationImpl[PriceUnit](*feeEstimation)

	// Create L1 resources bounds
	l1Gas := toResourceBounds[u64, u128, RB](
		innerFE.L1GasPrice,
		limits.L1GasPriceLimit,
		innerFE.L1GasConsumed,
		limits.L1GasAmountLimit,
		multiplier,
	)
	l1DataGas := toResourceBounds[u64, u128, RB](
		innerFE.L1DataGasPrice,
		limits.L1DataGasPriceLimit,
		innerFE.L1DataGasConsumed,
		limits.L1DataGasAmountLimit,
		multiplier,
	)

	// Create L2 resource bounds
	l2Gas := toResourceBounds[u64, u128, RB](
		innerFE.L2GasPrice,
		limits.L2GasPriceLimit,
		innerFE.L2GasConsumed,
		limits.L2GasAmountLimit,
		multiplier,
	)

	*resources = RBM{
		L1Gas:     l1Gas,
		L1DataGas: l1DataGas,
		L2Gas:     l2Gas,
	}

	return resources
}

// toResourceBounds converts a gas price and gas consumed to a ResourceBounds with
// applied multiplier.
//
// Parameters:
//   - gasPrice: The gas price
//   - gasPriceLimit: The limit for the gas price. If invalid, a default value
//     will be used.
//   - gasConsumed: The gas consumed
//   - gasAmountLimit: The limit for the gas amount. If invalid, a default value
//     will be used.
//   - multiplier: Multiplier for max amount and max price per unit
//
// Returns:
//   - RB: Resource bounds with applied multiplier
func toResourceBounds[
	u64 constraints.U64,
	u128 constraints.U128,
	RB constraints.ResourceBounds[u64, u128],
](
	gasPrice *felt.Felt,
	gasPriceLimit interface{ ToBigInt() (*big.Int, error) },
	gasConsumed *felt.Felt,
	gasAmountLimit interface{ ToUint64() (uint64, error) },
	multiplier float64,
) RB {
	// multiplier must be greater than 0. Default to 0 if not
	if multiplier <= 0 {
		return RB{
			MaxAmount:       u64("0x0"),
			MaxPricePerUnit: u128("0x0"),
		}
	}

	// Convert felt to big.Int
	gasPriceInt := gasPrice.BigInt(new(big.Int))
	gasConsumedInt := gasConsumed.BigInt(new(big.Int))

	// multiply values by the multiplier
	maxAmount := new(big.Float)
	maxPricePerUnit := new(big.Float)

	maxAmount.Mul(new(big.Float).SetInt(gasConsumedInt), big.NewFloat(multiplier))
	maxPricePerUnit.Mul(new(big.Float).SetInt(gasPriceInt), big.NewFloat(multiplier))
	// Convert big.Float to big.Int for proper hex formatting. The result is a truncated int
	maxAmountInt, _ := maxAmount.Int(new(big.Int))
	maxPricePerUnitInt, _ := maxPricePerUnit.Int(new(big.Int))

	// Get the limits OR set default values if invalid
	gasPL, err := gasPriceLimit.ToBigInt()
	if err != nil {
		gasPL = internalUtils.HexToBN(maxUint128)
	}
	tempGasAL, err := gasAmountLimit.ToUint64()
	if err != nil {
		tempGasAL = uint64(math.MaxUint64)
	}
	gasAL := new(big.Int).SetUint64(tempGasAL)

	// Check for overflow comparing with the limits
	if maxAmountInt.Cmp(gasAL) > 0 {
		maxAmountInt = gasAL
	}
	if maxPricePerUnitInt.Cmp(gasPL) > 0 {
		maxPricePerUnitInt = gasPL
	}

	return RB{
		MaxAmount:       u64(fmt.Sprintf("%#x", maxAmountInt)),
		MaxPricePerUnit: u128(fmt.Sprintf("%#x", maxPricePerUnitInt)),
	}
}

// ResBoundsMapToOverallFee calculates the overall fee for a ResourceBoundsMapping with
// applied multipliers.
// Parameters:
//   - resBounds: The resource bounds to calculate the fee for
//     (e.g. *rpcv9.ResourceBoundsMapping, *rpcv10.ResourceBoundsMapping, etc.)
//   - multiplier: Multiplier for max amount and max price per unit. Recommended to be 1.5,
//     but at least greater than 0
//   - tip: The tip amount in FRI in hexadecimal string format
//
// Returns:
//   - *felt.Felt: The overall fee in FRI
//   - error: An error if any
func ResBoundsMapToOverallFee[
	u64, Tip constraints.U64,
	u128 constraints.U128,
	RB constraints.ResourceBounds[u64, u128],
	RBM constraints.ResourceBoundsMapping[u64, u128, RB],
](
	resBounds *RBM,
	multiplier float64,
	tip Tip,
) (*felt.Felt, error) {
	if resBounds == nil {
		return nil, errors.New("resource bounds are nil")
	}
	innerResBounds := constraints.ResourceBoundsMappingImpl[u64, u128, RB](*resBounds)

	// negative multiplier is not allowed
	if multiplier <= 0 {
		return nil, errors.New("multiplier must be greater than 0")
	}

	tipInt, err := tip.ToUint64()
	if err != nil {
		return nil, fmt.Errorf("invalid tip: %w", err)
	}
	tipBigInt := new(big.Int).SetUint64(tipInt)

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

	innerL1Gas := constraints.ResourceBoundsImpl[u64, u128](innerResBounds.L1Gas)
	l1GasAmount, err := parseBound(string(innerL1Gas.MaxAmount))
	if err != nil {
		return nil, err
	}
	l1GasPrice, err := parseBound(string(innerL1Gas.MaxPricePerUnit))
	if err != nil {
		return nil, err
	}

	innerL1DataGas := constraints.ResourceBoundsImpl[u64, u128](innerResBounds.L1DataGas)
	l1DataGasAmount, err := parseBound(string(innerL1DataGas.MaxAmount))
	if err != nil {
		return nil, err
	}
	l1DataGasPrice, err := parseBound(string(innerL1DataGas.MaxPricePerUnit))
	if err != nil {
		return nil, err
	}

	innerL2Gas := constraints.ResourceBoundsImpl[u64, u128](innerResBounds.L2Gas)
	l2GasAmount, err := parseBound(string(innerL2Gas.MaxAmount))
	if err != nil {
		return nil, err
	}
	l2GasPrice, err := parseBound(string(innerL2Gas.MaxPricePerUnit))
	if err != nil {
		return nil, err
	}

	// calculate fee
	// Ref: https://docs.starknet.io/learn/protocol/fees#overall-fee
	l1GasFee := new(big.Int).Mul(l1GasAmount, l1GasPrice)
	l1DataGasFee := new(big.Int).Mul(l1DataGasAmount, l1DataGasPrice)
	l2GasFee := l2GasPrice.Add(l2GasPrice, tipBigInt).Mul(l2GasPrice, l2GasAmount)
	overallFee := l1GasFee.Add(l1GasFee, l1DataGasFee).Add(l1GasFee, l2GasFee)

	// multiply fee by multiplier
	multipliedOverallFee := new(
		big.Float,
	).Mul(new(big.Float).SetInt(overallFee), big.NewFloat(multiplier))
	overallFeeInt, _ := multipliedOverallFee.Int(nil) // truncated int

	// Convert big.Int to felt. SetString() validates if it's a valid felt
	return new(felt.Felt).SetString(fmt.Sprintf("%#x", overallFeeInt))
}
