package rpc

import (
	"context"
	"encoding/json"
	"errors"

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

	switch rawTxnTrace["type"] {
	case TransactionType_Invoke:
		var trace InvokeTxnTrace
		err = json.Unmarshal(rawTraceByte, &trace)
		if err != nil {
			return nil, err
		}
		return trace, nil
	case TransactionType_Declare:
		var trace DeclareTxnTrace
		err = json.Unmarshal(rawTraceByte, &trace)
		if err != nil {
			return nil, err
		}
		return trace, nil
	case TransactionType_DeployAccount:
		var trace DeployAccountTxnTrace
		err = json.Unmarshal(rawTraceByte, &trace)
		if err != nil {
			return nil, err
		}
		return trace, nil
	case TransactionType_L1Handler:
		var trace L1HandlerTxnTrace
		err = json.Unmarshal(rawTraceByte, &trace)
		if err != nil {
			return nil, err
		}
		return trace, nil
	}
	return nil, errors.New("Unknown transaction type")

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
