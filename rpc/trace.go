package rpc

import (
	"context"
	"encoding/json"

	"github.com/NethermindEth/juno/core/felt"
)

// For a given executed transaction, return the trace of its execution, including internal calls
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
