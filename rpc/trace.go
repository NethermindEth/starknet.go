package rpc

import (
	"context"
	"encoding/json"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/client/rpcerr"
)

// TraceTransaction returns the transaction trace for the given transaction hash.
//
// Parameters:
//   - ctx: the context.Context object for the request
//   - transactionHash: the transaction hash to trace
//
// Returns:
//   - TxnTrace: the transaction trace
//   - error: an error if the transaction trace cannot be retrieved
func (provider *Provider) TraceTransaction(
	ctx context.Context,
	transactionHash *felt.Felt,
) (TxnTrace, error) {
	var rawTxnTrace map[string]any
	if err := do(
		ctx, provider.c, "starknet_traceTransaction", &rawTxnTrace, transactionHash,
	); err != nil {
		return nil, rpcerr.UnwrapToRPCErr(err, ErrHashNotFound, ErrNoTraceAvailable)
	}

	rawTraceByte, err := json.Marshal(rawTxnTrace)
	if err != nil {
		return nil, rpcerr.Err(rpcerr.InternalError, StringErrData(err.Error()))
	}

	switch rawTxnTrace["type"] {
	case string(TransactionTypeInvoke):
		var trace InvokeTxnTrace
		err = json.Unmarshal(rawTraceByte, &trace)
		if err != nil {
			return nil, rpcerr.Err(rpcerr.InternalError, StringErrData(err.Error()))
		}

		return trace, nil
	case string(TransactionTypeDeclare):
		var trace DeclareTxnTrace
		err = json.Unmarshal(rawTraceByte, &trace)
		if err != nil {
			return nil, rpcerr.Err(rpcerr.InternalError, StringErrData(err.Error()))
		}

		return trace, nil
	case string(TransactionTypeDeployAccount):
		var trace DeployAccountTxnTrace
		err = json.Unmarshal(rawTraceByte, &trace)
		if err != nil {
			return nil, rpcerr.Err(rpcerr.InternalError, StringErrData(err.Error()))
		}

		return trace, nil
	case string(TransactionTypeL1Handler):
		var trace L1HandlerTxnTrace
		err = json.Unmarshal(rawTraceByte, &trace)
		if err != nil {
			return nil, rpcerr.Err(rpcerr.InternalError, StringErrData(err.Error()))
		}

		return trace, nil
	}

	return nil, rpcerr.Err(rpcerr.InternalError, StringErrData("Unknown transaction type"))
}

// TraceBlockTransactions retrieves the traces of transactions in a given block.
//
// Parameters:
//   - ctx: the context.Context object for controlling the request
//   - blockID: the block to retrieve the traces from. `pre_confirmed` tag is not allowed
//
// Returns:
//   - []Trace: a slice of Trace objects representing the traces of transactions in the block
//   - error: an error if there was a problem retrieving the traces.
func (provider *Provider) TraceBlockTransactions(
	ctx context.Context,
	blockID BlockID,
) ([]Trace, error) {
	err := checkForPreConfirmed(blockID)
	if err != nil {
		return nil, err
	}

	var output []Trace
	if err := do(
		ctx, provider.c, "starknet_traceBlockTransactions", &output, blockID,
	); err != nil {
		return nil, rpcerr.UnwrapToRPCErr(err, ErrBlockNotFound)
	}

	return output, nil
}

// SimulateTransactions returns the execution trace and consumed resources of the
// required transactions. When RETURN_INITIAL_READS is not present in simulation_flags,
// returns an array. When RETURN_INITIAL_READS is present in simulation_flags, returns an
// object with simulated_transactions and initial_reads fields.
//
// Parameters:
//   - ctx: The context of the function call
//   - blockID: The hash of the requested block, or number (height) of the requested
//     block, or a block tag, for
//     the block referencing the state or call the transaction on.
//   - txns: A sequence of transactions to simulate, running each transaction on the
//     state resulting from applying all the previous ones
//   - simulationFlags: Describes what parts of the transaction should be executed
//
// Returns:
//   - SimulateTxResult: The execution trace and consumed resources of the
//     required transactions + initial reads if the RETURN_INITIAL_READS flag was set.
//   - error: An error if any occurred during the execution
func (provider *Provider) SimulateTransactions(
	ctx context.Context,
	blockID BlockID,
	txns []BroadcastTxn,
	simulationFlags []SimulationFlag,
) (SimulateTxResult, error) {
	var output SimulateTxResult
	if err := do(
		ctx, provider.c, "starknet_simulateTransactions", &output, blockID, txns, simulationFlags,
	); err != nil {
		return output, rpcerr.UnwrapToRPCErr(err, ErrTxnExec, ErrBlockNotFound)
	}

	return output, nil
}
