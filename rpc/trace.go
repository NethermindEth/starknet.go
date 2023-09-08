package rpc

import (
	"context"
	"errors"
)

var ErrNotImplemented = errors.New("not implemented")

// not implemented for testing yet
func (provider *Provider) TransactionTrace(ctx context.Context, hash string) error {
	return ErrNotImplemented
}

// not implemented for testing yet
func (provider *Provider) TraceBlockTransactions(ctx context.Context, hash string) error {
	return ErrNotImplemented
}

// simulate a given transaction on the requested state, and generate the execution trace
func (provider *Provider) SimulateTransactions(ctx context.Context, blockID BlockID, txns []BroadcastedTransaction, simulationFlags []SimulationFlag) ([]SimulatedTransaction, error) {

	var output []SimulatedTransaction
	if err := do(ctx, provider.c, "starknet_simulateTransactions", &output, blockID, txns, simulationFlags); err != nil {
		return nil, tryUnwrapToRPCErr(err, ErrContractNotFound, ErrContractError, ErrBlockNotFound)
	}

	return output, nil

}
