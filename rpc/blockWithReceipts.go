package rpc

import (
	"context"
	"encoding/json"

	"github.com/NethermindEth/starknet.go/client/rpcerr"
)

// Get block information with full transactions and receipts given the block id
func (provider *Provider) BlockWithReceipts(
	ctx context.Context,
	blockID BlockID,
) (interface{}, error) {
	var result json.RawMessage
	if err := do(ctx, provider.c, "starknet_getBlockWithReceipts", &result, blockID); err != nil {
		return nil, rpcerr.UnwrapToRPCErr(err, ErrBlockNotFound)
	}

	var m map[string]interface{}
	if err := json.Unmarshal(result, &m); err != nil {
		return nil, rpcerr.Err(rpcerr.InternalError, StringErrData(err.Error()))
	}

	// Pre_confirmedBlockWithReceipts doesn't contain a "status" field
	if _, ok := m["status"]; ok {
		var block BlockWithReceipts
		if err := json.Unmarshal(result, &block); err != nil {
			return nil, rpcerr.Err(rpcerr.InternalError, StringErrData(err.Error()))
		}

		return &block, nil
	} else {
		var preConfirmedBlock PreConfirmedBlockWithReceipts
		if err := json.Unmarshal(result, &preConfirmedBlock); err != nil {
			return nil, rpcerr.Err(rpcerr.InternalError, StringErrData(err.Error()))
		}

		return &preConfirmedBlock, nil
	}
}
