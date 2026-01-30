package rpcv10

import (
	"context"
	"encoding/json"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/client/rpcerr"
	"github.com/NethermindEth/starknet.go/rpc/callers"
	"github.com/NethermindEth/starknet.go/rpc/internal"
	"github.com/NethermindEth/starknet.go/rpc/types"
)

// SimulateTransactions simulates transactions on the blockchain.
// Simulate a given sequence of transactions on the requested state, and generate
// the execution traces.
// Note that some of the transactions may revert, in which case no error is thrown,
// but revert details can be seen on the returned trace object.
// Note that some of the transactions may revert, this will be reflected by the
// revert_error property in the trace. Other types of failures (e.g. unexpected error
// or failure in the validation phase) will result in TRANSACTION_EXECUTION_ERROR.
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
//   - []SimulatedTransaction: The execution trace and consumed resources of the
//     required transactions
//   - error: An error if any occurred during the execution
func SimulateTransactions(
	ctx context.Context,
	c callers.Caller,
	blockID BlockID,
	txns []BroadcastTxn,
	simulationFlags []types.SimulationFlag,
) ([]SimulatedTransaction, error) {
	var output []SimulatedTransaction
	if err := internal.Do(
		ctx, c, "starknet_simulateTransactions", &output, blockID, txns, simulationFlags,
	); err != nil {
		return nil, rpcerr.UnwrapToRPCErr(err, ErrTxnExec, ErrBlockNotFound)
	}

	return output, nil
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
func TraceBlockTransactions(
	ctx context.Context,
	c callers.Caller,
	blockID BlockID,
) ([]Trace, error) {
	err := checkForPreConfirmed(blockID)
	if err != nil {
		return nil, err
	}

	var output []Trace
	if err := internal.Do(
		ctx, c, "starknet_traceBlockTransactions", &output, blockID,
	); err != nil {
		return nil, rpcerr.UnwrapToRPCErr(err, ErrBlockNotFound)
	}

	return output, nil
}

// TraceTransaction returns the transaction trace for the given transaction hash.
//
// Parameters:
//   - ctx: the context.Context object for the request
//   - transactionHash: the transaction hash to trace
//
// Returns:
//   - TxnTrace: the transaction trace
//   - error: an error if the transaction trace cannot be retrieved
func TraceTransaction(
	ctx context.Context,
	c callers.Caller,
	transactionHash *felt.Felt,
) (TxnTrace, error) {
	var rawTxnTrace map[string]any
	if err := internal.Do(
		ctx, c, "starknet_traceTransaction", &rawTxnTrace, transactionHash,
	); err != nil {
		return nil, rpcerr.UnwrapToRPCErr(err, ErrHashNotFound, ErrNoTraceAvailable)
	}

	rawTraceByte, err := json.Marshal(rawTxnTrace)
	if err != nil {
		return nil, rpcerr.Err(rpcerr.InternalError, StringErrData(err.Error()))
	}

	switch rawTxnTrace["type"] {
	case string(types.TransactionTypeInvoke):
		var trace InvokeTxnTrace
		err = json.Unmarshal(rawTraceByte, &trace)
		if err != nil {
			return nil, rpcerr.Err(rpcerr.InternalError, StringErrData(err.Error()))
		}

		return trace, nil
	case string(types.TransactionTypeDeclare):
		var trace DeclareTxnTrace
		err = json.Unmarshal(rawTraceByte, &trace)
		if err != nil {
			return nil, rpcerr.Err(rpcerr.InternalError, StringErrData(err.Error()))
		}

		return trace, nil
	case string(types.TransactionTypeDeployAccount):
		var trace DeployAccountTxnTrace
		err = json.Unmarshal(rawTraceByte, &trace)
		if err != nil {
			return nil, rpcerr.Err(rpcerr.InternalError, StringErrData(err.Error()))
		}

		return trace, nil
	case string(types.TransactionTypeL1Handler):
		var trace L1HandlerTxnTrace
		err = json.Unmarshal(rawTraceByte, &trace)
		if err != nil {
			return nil, rpcerr.Err(rpcerr.InternalError, StringErrData(err.Error()))
		}

		return trace, nil
	}

	return nil, rpcerr.Err(rpcerr.InternalError, StringErrData("Unknown transaction type"))
}
