package methods

import (
	"context"

	"github.com/NethermindEth/starknet.go/client/rpcerr"
	"github.com/NethermindEth/starknet.go/rpc/callers"
	"github.com/NethermindEth/starknet.go/rpc/internal"
	"github.com/NethermindEth/starknet.go/rpc/rpcv10"
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
	blockID rpcv10.BlockID,
	txns []types.BroadcastTxn,
	simulationFlags []rpcv10.SimulationFlag,
) ([]rpcv10.SimulatedTransaction, error) {
	var output []rpcv10.SimulatedTransaction
	if err := internal.Do(
		ctx, c, "starknet_simulateTransactions", &output, blockID, txns, simulationFlags,
	); err != nil {
		return nil, rpcerr.UnwrapToRPCErr(err, rpcv10.ErrTxnExec, rpcv10.ErrBlockNotFound)
	}

	return output, nil
}
