package rpc

import (
	"context"
	"errors"

	"github.com/NethermindEth/juno/core/felt"
)

var ErrNotImplemented = errors.New("not implemented")

// not implemented for testing yet
func (provider *Provider) TransactionTrace(ctx context.Context, hash string) error {
	return ErrNotImplemented
}

// Retrieve traces for all transactions in the given block
func (provider *Provider) TraceBlockTransactions(ctx context.Context, blockHash *felt.Felt) ([]Trace, error) {
	var output []Trace
	if err := do(ctx, provider.c, "starknet_traceBlockTransactions", &output, blockHash); err != nil {
		return nil, tryUnwrapToRPCErr(err, ErrInvalidBlockHash)
	}
	return output, nil

}

// simulate a given transaction on the requested state, and generate the execution trace
func (provider *Provider) SimulateTransactions(ctx context.Context, blockID BlockID, txns []Transaction, simulationFlags []SimulationFlag) ([]SimulatedTransaction, error) {

	var output []SimulatedTransaction
	if err := do(ctx, provider.c, "starknet_simulateTransactions", &output, blockID, txns, simulationFlags); err != nil {
		return nil, tryUnwrapToRPCErr(err, ErrContractNotFound, ErrContractError, ErrBlockNotFound)
	}

	return output, nil

}
