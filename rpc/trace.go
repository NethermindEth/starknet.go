package rpc

import (
	"context"
	"encoding/json"

	"github.com/NethermindEth/juno/core/felt"
)

// TransactionTrace returns the transaction trace for the given transaction hash.
//
// Parameters:
//   - ctx: the context.Context object for the request
//   - transactionHash: the transaction hash to trace
// Returns:
//   - TxnTrace: the transaction trace
//   - error: an error if the transaction trace cannot be retrieved
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

// TraceBlockTransactions retrieves the traces of transactions in a given block.
//
// Parameters:
// - ctx: the context.Context object for controlling the request
// - blockHash: the hash of the block to retrieve the traces from
// Returns:
// - []Trace: a slice of Trace objects representing the traces of transactions in the block
// - error: an error if there was a problem retrieving the traces.
func (provider *Provider) TraceBlockTransactions(ctx context.Context, blockHash *felt.Felt) ([]Trace, error) {
	var output []Trace
	if err := do(ctx, provider.c, "starknet_traceBlockTransactions", &output, blockHash); err != nil {
		return nil, tryUnwrapToRPCErr(err, ErrInvalidBlockHash)
	}
	return output, nil

}

// SimulateTransactions simulates transactions on the blockchain.
//
// Parameters:
// - ctx: the context.Context object for controlling the request
// - blockID: the ID of the block on which the transactions should be simulated
// - txns: a slice of Transaction objects representing the transactions to be simulated
// - simulationFlags: a slice of SimulationFlag objects representing additional simulation flags
// Returns:
// - []SimulatedTransaction: a slice of SimulatedTransaction objects representing the simulated transactions
// - error: an error if any error occurs during the simulation process
func (provider *Provider) SimulateTransactions(ctx context.Context, blockID BlockID, txns []Transaction, simulationFlags []SimulationFlag) ([]SimulatedTransaction, error) {

	var output []SimulatedTransaction
	if err := do(ctx, provider.c, "starknet_simulateTransactions", &output, blockID, txns, simulationFlags); err != nil {
		return nil, tryUnwrapToRPCErr(err, ErrContractNotFound, ErrContractError, ErrBlockNotFound)
	}

	return output, nil

}
