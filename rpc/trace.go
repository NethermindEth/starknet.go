package rpc

import (
	"context"
	"errors"

	"github.com/NethermindEth/juno/core/felt"
)

var ErrNotImplemented = errors.New("not implemented")

// For a given executed transaction, return the trace of its execution, including internal calls
func (provider *Provider) TransactionTrace(ctx context.Context, transactionHash *felt.Felt) (TxnTrace, error) {
	var output TxnTrace
	if err := do(ctx, provider.c, "starknet_traceTransaction", &output, transactionHash); err != nil {
		return nil, tryUnwrapToRPCErr(err, ErrInvalidTxnHash, ErrNoTraceAvailable)
	}

	return output, nil

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
func (provider *Provider) SimulateTransactions(ctx context.Context, blockID BlockID, txns []BroadcastedTransaction, simulationFlags []SimulationFlag) ([]SimulatedTransaction, error) {

	var output []SimulatedTransaction
	if err := do(ctx, provider.c, "starknet_simulateTransactions", &output, blockID, txns, simulationFlags); err != nil {
		return nil, tryUnwrapToRPCErr(err, ErrContractNotFound, ErrContractError, ErrBlockNotFound)
	}

	return output, nil

}
