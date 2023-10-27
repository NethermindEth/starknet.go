package rpc

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/NethermindEth/juno/core/felt"
)

// TraceTransaction returns the transaction trace for the given transaction hash.
//
// Parameters:
//   - ctx: the context.Context object for the request
//   - transactionHash: the transaction hash to trace
// Returns:
//   - TxnTrace: the transaction trace
//   - error: an error if the transaction trace cannot be retrieved
func (provider *Provider) TraceTransaction(ctx context.Context, transactionHash *felt.Felt) (TxnTrace, error) {
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

	switch rawTxnTrace["type"] {
	case string(TransactionType_Invoke):
		var trace InvokeTxnTrace
		err = json.Unmarshal(rawTraceByte, &trace)
		if err != nil {
			return nil, err
		}
		return trace, nil
	case string(TransactionType_Declare):
		var trace DeclareTxnTrace
		err = json.Unmarshal(rawTraceByte, &trace)
		if err != nil {
			return nil, err
		}
		return trace, nil
	case string(TransactionType_DeployAccount):
		var trace DeployAccountTxnTrace
		err = json.Unmarshal(rawTraceByte, &trace)
		if err != nil {
			return nil, err
		}
		return trace, nil
	case string(TransactionType_L1Handler):
		var trace L1HandlerTxnTrace
		err = json.Unmarshal(rawTraceByte, &trace)
		if err != nil {
			return nil, err
		}
		return trace, nil
	}
	return nil, errors.New("Unknown transaction type")

}

// TraceBlockTransactions retrieves the traces of transactions in a given block.
//
// Parameters:
// - ctx: the context.Context object for controlling the request
// - blockHash: the hash of the block to retrieve the traces from
// Returns:
// - []Trace: a slice of Trace objects representing the traces of transactions in the block
// - error: an error if there was a problem retrieving the traces.
func (provider *Provider) TraceBlockTransactions(ctx context.Context, blockID BlockID) ([]Trace, error) {
	var output []Trace
	if err := do(ctx, provider.c, "starknet_traceBlockTransactions", &output, blockID); err != nil {
		return nil, tryUnwrapToRPCErr(err, ErrBlockNotFound)
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
