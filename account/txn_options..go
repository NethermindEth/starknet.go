package account

import "github.com/NethermindEth/starknet.go/rpc"

// Optional settings when building/sending/estimating a transaction.
type TxnOptions struct {
	// Tip amount in FRI for the transaction. Default: `"0x0"`.
	Tip rpc.U64

	// A boolean flag indicating whether the transaction version should have
	// the query bit when estimating fees. If true, the transaction version
	// will be `rpc.TransactionV3WithQueryBit` (0x100000000000000000000000000000003).
	// If false, the transaction version will be `rpc.TransactionV3` (0x3).
	// In case of doubt, set to `false`. Default: `false`.
	WithQueryBitVersion bool

	// A safety factor for fee estimation that helps prevent transaction
	// failures due to fee fluctuations. It multiplies both the max amount
	// and max price per unit by this value.
	// A value of 1.5 (estimated fee + 50%) is recommended to balance between
	// transaction success rate and avoiding excessive fees. Higher values
	// provide more safety margin but may result in overpayment.
	// If multiplier <= 0, it'll be set to 1.5. Default: `1.5`.
	Multiplier float64

	// The block tag to be used for fee estimation. Default: `"pending"`.
	EstimationBlockTag rpc.BlockTag
	// The flag to be used when estimating fees. Default: none.
	SimulationFlag rpc.SimulationFlag
}

// SafeMultiplier returns the multiplier for the transaction. If the multiplier is not set or negative, returns 1.5.
func (opts *TxnOptions) SafeMultiplier() float64 {
	if opts == nil || opts.Multiplier <= 0 {
		return 1.5
	}

	return opts.Multiplier
}

// BlockID returns the block ID for fee estimation based on the EstimationBlockTag.
// If EstimationBlockTag is not set, returns the pending block ID.
func (opts *TxnOptions) BlockID() rpc.BlockID {
	if opts == nil || opts.EstimationBlockTag == "" {
		return rpc.WithBlockTag(rpc.BlockTagPending)
	}

	return rpc.WithBlockTag(opts.EstimationBlockTag)
}

// Returns a `[]rpc.SimulationFlag` containing the SimulationFlag.
// If the flag is not set, returns an empty slice.
func (opts *TxnOptions) SimulationFlags() []rpc.SimulationFlag {
	if opts == nil || opts.SimulationFlag == "" {
		return []rpc.SimulationFlag{}
	}

	return []rpc.SimulationFlag{opts.SimulationFlag}
}
