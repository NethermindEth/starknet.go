package methods

import (
	"context"
	"encoding/json"

	"github.com/NethermindEth/starknet.go/client/rpcerr"
	"github.com/NethermindEth/starknet.go/rpc/rpcv10"
)

// Get block information with full transactions and receipts given the block id
func GetBlockWithReceipts(
	ctx context.Context,
	c callCloser,
	blockID rpcv10.BlockID,
) (interface{}, error) {
	var result json.RawMessage
	if err := do(ctx, c, "starknet_getBlockWithReceipts", &result, blockID); err != nil {
		return nil, rpcerr.UnwrapToRPCErr(err, ErrBlockNotFound)
	}

	var m map[string]interface{}
	if err := json.Unmarshal(result, &m); err != nil {
		return nil, rpcerr.Err(rpcerr.InternalError, StringErrData(err.Error()))
	}

	// Pre_confirmedBlockWithReceipts doesn't contain a "status" field
	if _, ok := m["status"]; ok {
		var block rpcv10.BlockWithReceipts
		if err := json.Unmarshal(result, &block); err != nil {
			return nil, rpcerr.Err(rpcerr.InternalError, StringErrData(err.Error()))
		}

		return &block, nil
	} else {
		var preConfirmedBlock rpcv10.PreConfirmedBlockWithReceipts
		if err := json.Unmarshal(result, &preConfirmedBlock); err != nil {
			return nil, rpcerr.Err(rpcerr.InternalError, StringErrData(err.Error()))
		}

		return &preConfirmedBlock, nil
	}
}
