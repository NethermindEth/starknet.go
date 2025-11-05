package account

import (
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"
)

const (
	defaultFeeMultiplier float64 = 1.5
	defaultTipMultiplier float64 = 1.0
)

// The `TxnOptions` struct is equal to the `utils.TxnOptions` struct + some new fields.
// Composition wasn't used to avoid the need to create a struct inside another struct
// when building the options.

// Optional settings for building/sending/estimating a transaction
// in the BuildAndSend* account methods.
type TxnOptions struct {
	// The multiplier to be used when estimating the tip when no custom tip is set.
	// If <= 0, it'll be set to 1.0 (no multiplier, just the estimated tip).
	TipMultiplier float64
	// A custom tip amount in FRI for the transaction in hexadecimal format.
	// If not set, the tip will be automatically estimated for the transaction.
	CustomTip rpc.U64

	// A boolean flag indicating whether the transaction version should have
	// the query bit when estimating fees. If true, the transaction version
	// will be `rpc.TransactionV3WithQueryBit` (0x100000000000000000000000000000003).
	// If false, the transaction version will be `rpc.TransactionV3` (0x3).
	// In case of doubt, set to `false`. Default: `false`.
	UseQueryBit bool

	// A safety factor for fee estimation that helps prevent transaction
	// failures due to fee fluctuations. It multiplies both the max amount
	// and max price per unit by this value.
	// A value of 1.5 (estimated fee + 50%) is recommended to balance between
	// transaction success rate and avoiding excessive fees. Higher values
	// provide more safety margin but may result in overpayment.
	// If FeeMultiplier <= 0, it'll be set to 1.5.
	FeeMultiplier float64

	// A boolean flag indicating whether to use the latest block tag
	// when estimating fees instead of the pre_confirmed block. Default: `false`.
	UseLatest bool
	// The simulation flag to be used when estimating fees. Default: none.
	SimulationFlag rpc.SimulationFlag
}

// BlockID returns the block ID for fee estimation based on the UseLatest flag.
// If UseLatest is `true`, returns the latest block ID.
// Otherwise, returns the pre_confirmed block ID.
func (opts *TxnOptions) BlockID() rpc.BlockID {
	if opts.UseLatest {
		return rpc.WithBlockTag(rpc.BlockTagLatest)
	}

	return rpc.WithBlockTag(rpc.BlockTagPreConfirmed)
}

// Returns a `[]rpc.SimulationFlag` containing the SimulationFlag.
// If the flag is not set, returns an empty slice.
func (opts *TxnOptions) SimulationFlags() []rpc.SimulationFlag {
	if opts.SimulationFlag == "" {
		return []rpc.SimulationFlag{}
	}

	return []rpc.SimulationFlag{opts.SimulationFlag}
}

// FmtFeeMultiplier returns the fee multiplier specified in the options.
// If not set/negative, returns the default fee multiplier.
func (opts *TxnOptions) FmtFeeMultiplier() float64 {
	if opts.FeeMultiplier <= 0 {
		return defaultFeeMultiplier
	}

	return opts.FeeMultiplier
}

// FmtTipMultiplier returns the tip multiplier specified in the options.
// If not set/negative, returns the default tip multiplier.
func (opts *TxnOptions) FmtTipMultiplier() float64 {
	if opts.TipMultiplier <= 0 {
		return defaultTipMultiplier
	}

	return opts.TipMultiplier
}

type UDCOptions = utils.UDCOptions
