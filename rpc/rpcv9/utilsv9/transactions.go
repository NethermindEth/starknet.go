package utilsv9

import (
	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/contracts"
	"github.com/NethermindEth/starknet.go/rpc/internal/utilsv"
	"github.com/NethermindEth/starknet.go/rpc/rpcv9"
	"github.com/NethermindEth/starknet.go/rpc/types"
)

// Optional settings when building a transaction.
type TxnOptions struct {
	// Tip amount in FRI for the transaction. Default: `"0x0"`.
	Tip types.U64
	// A boolean flag indicating whether the transaction version should have
	// the query bit when estimating fees. If true, the transaction version
	// will be `rpcv9.TransactionV3WithQueryBit` (0x100000000000000000000000000000003).
	// If false, the transaction version will be `rpcv9.TransactionV3` (0x3).
	// In case of doubt, set to `false`. Default: `false`.
	UseQueryBit bool

	// ONLY FOR DECLARE TXN: A boolean flag indicating whether to use the Blake2s hash
	// function to calculate the compiled class hash. This must be set to true after
	// the Starknet v0.14.1 upgrade, when the Poseidon hash function will be deprecated
	// for the compiled class hash.
	UseBlake2sHash bool
	// TODO: remove this field after the Starknet v0.14.1 upgrade
}

// TxnVersion returns `rpcv9.TransactionV3WithQueryBit` when UseQueryBit is true, and
// `rpcv9.TransactionV3` if false.
func (opts *TxnOptions) TxnVersion() rpcv9.TransactionVersion {
	if opts.UseQueryBit {
		return rpcv9.TransactionV3WithQueryBit
	}

	return rpcv9.TransactionV3
}

// SafeTip returns the tip amount in FRI for the transaction. If the tip is not set
// or invalid, returns "0x0".
func (opts *TxnOptions) SafeTip() types.U64 {
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
	resourceBounds *rpcv9.ResourceBoundsMapping,
	opts *TxnOptions,
) *rpcv9.BroadcastInvokeTxnV3 {
	if opts == nil {
		opts = new(TxnOptions)
	}

	return &rpcv9.BroadcastInvokeTxnV3{
		Type:                  rpcv9.TransactionTypeInvoke,
		SenderAddress:         senderAddress,
		Calldata:              calldata,
		Version:               opts.TxnVersion(),
		Signature:             []*felt.Felt{},
		Nonce:                 nonce,
		ResourceBounds:        resourceBounds,
		Tip:                   opts.SafeTip(),
		PayMasterData:         []*felt.Felt{},
		AccountDeploymentData: []*felt.Felt{},
		NonceDataMode:         rpcv9.DAModeL1,
		FeeMode:               rpcv9.DAModeL1,
	}
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
	resourceBounds *rpcv9.ResourceBoundsMapping,
	opts *TxnOptions,
) (*rpcv9.BroadcastDeclareTxnV3, error) {
	if opts == nil {
		opts = new(TxnOptions)
	}

	var tx rpcv9.BroadcastDeclareTxnV3
	return utilsv.BuildDeclareTxn(
		&tx,
		senderAddress,
		casmClass,
		contractClass,
		nonce,
		resourceBounds,
		opts,
	)
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
	resourceBounds *rpcv9.ResourceBoundsMapping,
	opts *TxnOptions,
) *rpcv9.BroadcastDeployAccountTxnV3 {
	if opts == nil {
		opts = new(TxnOptions)
	}

	var tx rpcv9.BroadcastDeployAccountTxnV3
	return utilsv.BuildDeployAccountTxn(
		&tx,
		nonce,
		contractAddressSalt,
		constructorCalldata,
		classHash,
		resourceBounds,
		opts,
	)
}

// FeeEstToResBoundsMap converts a FeeEstimation to ResourceBoundsMapping with applied multipliers.
// Parameters:
//   - feeEstimation: The fee estimation to convert
//   - multiplier: Multiplier for max amount and max price per unit. Recommended to be 1.5,
//     but at least greater than 0.
//     If multiplier <= 0, all resources bounds will be set to 0.
//     If resource bounds overflow, they will be set to the max allowed value (U64 or U128).
//
// Returns:
//   - types.ResourceBoundsMapping: Resource bounds with applied multipliers
func FeeEstToResBoundsMap(
	feeEstimation *rpcv9.FeeEstimation,
	multiplier float64,
) *rpcv9.ResourceBoundsMapping {
	var resources rpcv9.ResourceBoundsMapping

	utilsv.FeeEstToResBoundsMap(
		feeEstimation,
		&resources,
		multiplier,
	)

	// TODO: return by value instead of pointer
	return &resources
}

// CustomFeeEstToResBoundsMap converts a FeeEstimation to ResourceBoundsMapping with applied
// multipliers and limits.
// Parameters:
//   - feeEstimation: The fee estimation to convert
//   - multiplier: Multiplier for max amount and max price per unit. Recommended to be 1.5,
//     but at least greater than 0.
//     If multiplier <= 0, all resources bounds will be set to 0.
//     If resource bounds overflow, they will be set to the max allowed value (U64 or U128).
//   - limits: Limits for the resource bounds. They are still capped to the max allowed
//     values (U64 or U128).
//
// Returns:
//   - types.ResourceBoundsMapping: Resource bounds with applied multipliers and limits
func CustomFeeEstToResBoundsMap(
	feeEstimation *rpcv9.FeeEstimation,
	multiplier float64,
	limits *types.FeeLimits,
) rpcv9.ResourceBoundsMapping {
	var resources rpcv9.ResourceBoundsMapping

	utilsv.CustomFeeEstToResBoundsMap(
		feeEstimation,
		&resources,
		multiplier,
		limits,
	)

	return resources
}

// ResBoundsMapToOverallFee calculates the overall fee for a ResourceBoundsMapping with
// applied multipliers.
// Parameters:
//   - resBounds: The resource bounds to calculate the fee for
//   - multiplier: Multiplier for max amount and max price per unit. Recommended to be 1.5,
//     but at least greater than 0
//   - tip: The tip amount in FRI in hexadecimal string format
//
// Returns:
//   - *felt.Felt: The overall fee in FRI
//   - error: An error if any
func ResBoundsMapToOverallFee(
	resBounds *rpcv9.ResourceBoundsMapping,
	multiplier float64,
	tip types.U64,
) (*felt.Felt, error) {
	return utilsv.ResBoundsMapToOverallFee(
		resBounds,
		multiplier,
		tip,
	)
}
