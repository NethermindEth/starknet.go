package account

import "github.com/NethermindEth/starknet.go/rpc"

// Optional settings for building/sending/estimating a transaction
// in the BuildAndSend* account methods.
type TxnOptions struct {
	// Tip amount in FRI for the transaction. Default: `"0x0"`.
	Tip rpc.U64

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
	// If multiplier <= 0, it'll be set to 1.5.
	Multiplier float64

	// A boolean flag indicating whether to use the latest block tag
	// when estimating fees instead of the pending block. Default: `false`.
	UseLatest bool
	// The simulation flag to be used when estimating fees. Default: none.
	SimulationFlag rpc.SimulationFlag
}

// BlockID returns the block ID for fee estimation based on the UseLatest flag.
// If UseLatest is set, returns the latest block ID.
func (opts *TxnOptions) BlockID() rpc.BlockID {
	if opts.UseLatest {
		return rpc.WithBlockTag(rpc.BlockTagLatest)
	}

	return rpc.WithBlockTag(rpc.BlockTagPending)
}

// Returns a `[]rpc.SimulationFlag` containing the SimulationFlag.
// If the flag is not set, returns an empty slice.
func (opts *TxnOptions) SimulationFlags() []rpc.SimulationFlag {
	if opts.SimulationFlag == "" {
		return []rpc.SimulationFlag{}
	}

	return []rpc.SimulationFlag{opts.SimulationFlag}
}

// returns a copy of the default TxnOptions
func defaultTxnOptions() TxnOptions {
	return TxnOptions{
		Tip:            "0x0",
		UseQueryBit:    false,
		Multiplier:     1.5,
		UseLatest:      false,
		SimulationFlag: "",
	}
}

// takes a pointer to a TxnOptions struct and formats the tip and multiplier according to the default values
func fmtTipAndMultiplier(opts *TxnOptions) {
	if opts.Multiplier <= 0 {
		opts.Multiplier = 1.5
	}

	if opts.Tip == "" {
		opts.Tip = "0x0"
	}
}
