package methods

import (
	"context"
	"encoding/json"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/client/rpcerr"
	"github.com/NethermindEth/starknet.go/rpc/callers"
	"github.com/NethermindEth/starknet.go/rpc/internal"
	"github.com/NethermindEth/starknet.go/rpc/rpcv10"
	"github.com/NethermindEth/starknet.go/rpc/types"
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
func TraceTransaction(
	ctx context.Context,
	c callers.Caller,
	transactionHash *felt.Felt,
) (rpcv10.TxnTrace, error) {
	var rawTxnTrace map[string]any
	if err := internal.Do(
		ctx, c, "starknet_traceTransaction", &rawTxnTrace, transactionHash,
	); err != nil {
		return nil, rpcerr.UnwrapToRPCErr(err, rpcv10.ErrHashNotFound, rpcv10.ErrNoTraceAvailable)
	}

	rawTraceByte, err := json.Marshal(rawTxnTrace)
	if err != nil {
		return nil, rpcerr.Err(rpcerr.InternalError, rpcv10.StringErrData(err.Error()))
	}

	switch rawTxnTrace["type"] {
	case string(types.TransactionTypeInvoke):
		var trace rpcv10.InvokeTxnTrace
		err = json.Unmarshal(rawTraceByte, &trace)
		if err != nil {
			return nil, rpcerr.Err(rpcerr.InternalError, rpcv10.StringErrData(err.Error()))
		}

		return trace, nil
	case string(types.TransactionTypeDeclare):
		var trace rpcv10.DeclareTxnTrace
		err = json.Unmarshal(rawTraceByte, &trace)
		if err != nil {
			return nil, rpcerr.Err(rpcerr.InternalError, rpcv10.StringErrData(err.Error()))
		}

		return trace, nil
	case string(types.TransactionTypeDeployAccount):
		var trace rpcv10.DeployAccountTxnTrace
		err = json.Unmarshal(rawTraceByte, &trace)
		if err != nil {
			return nil, rpcerr.Err(rpcerr.InternalError, rpcv10.StringErrData(err.Error()))
		}

		return trace, nil
	case string(types.TransactionTypeL1Handler):
		var trace rpcv10.L1HandlerTxnTrace
		err = json.Unmarshal(rawTraceByte, &trace)
		if err != nil {
			return nil, rpcerr.Err(rpcerr.InternalError, rpcv10.StringErrData(err.Error()))
		}

		return trace, nil
	}

	return nil, rpcerr.Err(rpcerr.InternalError, rpcv10.StringErrData("Unknown transaction type"))
}
