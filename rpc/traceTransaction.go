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
