package rpc

import (
	"context"
	"encoding/json"

	"github.com/NethermindEth/juno/core/felt"
)

// TransactionTrace retrieves the transaction trace for a given transaction hash, including internal calls.
//
// ctx: The context.Context object for cancellation.
// transactionHash: The hash of the transaction to retrieve the trace for.
// Returns: The transaction trace object and an error if any.
func (provider *Provider) TransactionTrace(ctx context.Context, transactionHash *felt.Felt) (TxnTrace, error) {
	var rawTxnTrace map[string]any
	if err := do(ctx, provider.c, "starknet_traceTransaction", &rawTxnTrace, transactionHash); err != nil {
		if noTraceAvailableError, ok := isErrNoTraceAvailableError(err); ok {
			return nil, noTraceAvailableError
		}
		return nil, tryUnwrapToRPCErr(err, ErrInvalidTxnHash)
	}

	rawTraceByte, err := json.Marshal(rawTxnTrace)
	if err != nil {
		return nil, err
	}

	// if execute_invocation exists, then it's an InvokeTxnTrace type
	if _, exists := rawTxnTrace["execute_invocation"]; exists {
		var trace InvokeTxnTrace
		err = json.Unmarshal(rawTraceByte, &trace)
		if err != nil {
			return nil, err
		}
		return trace, nil
	}

	// if constructor_invocation exists, then it's a DeployAccountTxnTrace type
	if _, exists := rawTxnTrace["constructor_invocation"]; exists {
		var trace DeployAccountTxnTrace
		err = json.Unmarshal(rawTraceByte, &trace)
		if err != nil {
			return nil, err
		}
		return trace, nil
	}

	// if function_invocation exists, then it's an L1HandlerTxnTrace type
	if _, exists := rawTxnTrace["function_invocation"]; exists {
		var trace L1HandlerTxnTrace
		err = json.Unmarshal(rawTraceByte, &trace)
		if err != nil {
			return nil, err
		}
		return trace, nil
	}

	// the other possible choice is for it to be a DeclareTxnTrace type
	var trace DeclareTxnTrace
	err = json.Unmarshal(rawTraceByte, &trace)
	if err != nil {
		return nil, err
	}
	return trace, nil
}

// TraceBlockTransactions retrieves the trace of all transactions for a given block.
//
// ctx: The context.Context used for the function call.
// blockHash: The hash of the block for which to retrieve the trace.
// []Trace: The trace of transactions for the given block.
// error: An error if the function call fails.
func (provider *Provider) TraceBlockTransactions(ctx context.Context, blockHash *felt.Felt) ([]Trace, error) {
	var output []Trace
	if err := do(ctx, provider.c, "starknet_traceBlockTransactions", &output, blockHash); err != nil {
		return nil, tryUnwrapToRPCErr(err, ErrInvalidBlockHash)
	}
	return output, nil

}

// SimulateTransactions simulates transactions on the StarkNet blockchain, and generate the execution trace.
//
// ctx - The context of the function call.
// blockID - The ID of the block to simulate transactions on.
// txns - The list of transactions to simulate.
// simulationFlags - The flags to control the simulation process.
// []SimulatedTransaction - The list of simulated transactions.
// error - An error, if any occurred.
func (provider *Provider) SimulateTransactions(ctx context.Context, blockID BlockID, txns []Transaction, simulationFlags []SimulationFlag) ([]SimulatedTransaction, error) {

	var output []SimulatedTransaction
	if err := do(ctx, provider.c, "starknet_simulateTransactions", &output, blockID, txns, simulationFlags); err != nil {
		return nil, tryUnwrapToRPCErr(err, ErrContractNotFound, ErrContractError, ErrBlockNotFound)
	}

	return output, nil

}
