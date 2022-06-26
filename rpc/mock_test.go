package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/dontpanicdao/caigo/types"
)

var (
	errWrongType = fmt.Errorf("wrong type")
	errWrongArgs = fmt.Errorf("wrong number of args")
)

// rpcMock is a mock of the go-ethereum Client that can be used for local tests
// when no integration environment exists.
type rpcMock struct {
	closed bool
}

func (r *rpcMock) Close() {
	r.closed = true
}

func (r *rpcMock) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	switch method {
	case "starknet_blockNumber":
		return mock_starknet_blockNumber(result, method, args...)
	case "starknet_chainId":
		return mock_starknet_chainId(result, method, args...)
	case "starknet_syncing":
		return mock_starknet_syncing(result, method, args...)
	case "starknet_getBlockByHash":
		return mock_starknet_getBlockByHash(result, method, args...)
	case "starknet_getBlockByNumber":
		return mock_starknet_getBlockByNumber(result, method, args...)
	default:
		return ErrNotFound
	}
}

func mock_starknet_blockNumber(result interface{}, method string, args ...interface{}) error {
	r, ok := result.(*big.Int)
	if !ok || r == nil {
		return errWrongType
	}
	if len(args) != 0 {
		return errWrongArgs
	}
	value1 := big.NewInt(1)
	*r = *value1
	return nil
}

func mock_starknet_getBlockByHash(result interface{}, method string, args ...interface{}) error {
	r, ok := result.(*json.RawMessage)
	if !ok || r == nil {
		return errWrongType
	}
	if len(args) != 2 {
		return errWrongArgs
	}
	blockHash, ok := args[0].(string)
	if !ok || blockHash != "0xdeadbeef" {
		return errWrongArgs
	}
	blockFormat, ok := args[1].(string)
	if !ok || (blockFormat != "FULL_TXN_AND_RECEIPTS" && blockFormat != "FULL_TXNS") {
		return errWrongArgs
	}
	transactionReceipt := types.TransactionReceipt{}
	if blockFormat == "FULL_TXN_AND_RECEIPTS" {
		transactionReceipt = types.TransactionReceipt{
			Status: "ACCEPTED_ON_L1",
		}
	}
	transaction := types.Transaction{
		TransactionReceipt: transactionReceipt,
		TransactionHash:    "0xdeadbeef",
	}
	output := types.Block{
		BlockNumber:  1000,
		BlockHash:    "0xdeadbeef",
		Transactions: []*types.Transaction{&transaction},
	}
	outputContent, _ := json.Marshal(output)
	json.Unmarshal(outputContent, r)
	return nil
}

func mock_starknet_getBlockByNumber(result interface{}, method string, args ...interface{}) error {
	r, ok := result.(*json.RawMessage)
	if !ok || r == nil {
		return errWrongType
	}
	if len(args) != 2 {
		return errWrongArgs
	}
	blockNumber, ok := args[0].(uint64)
	if !ok || blockNumber != 1000 {
		return errWrongArgs
	}
	blockFormat, ok := args[1].(string)
	if !ok || (blockFormat != "FULL_TXN_AND_RECEIPTS" && blockFormat != "FULL_TXNS") {
		return errWrongArgs
	}
	transactionReceipt := types.TransactionReceipt{}
	if blockFormat == "FULL_TXN_AND_RECEIPTS" {
		transactionReceipt = types.TransactionReceipt{
			Status: "ACCEPTED_ON_L1",
		}
	}
	transaction := types.Transaction{
		TransactionReceipt: transactionReceipt,
		TransactionHash:    "0xdeadbeef",
	}
	output := types.Block{
		BlockHash:    "0xdeadbeef",
		Transactions: []*types.Transaction{&transaction},
	}
	outputContent, _ := json.Marshal(output)
	json.Unmarshal(outputContent, r)
	return nil
}

func mock_starknet_chainId(result interface{}, method string, args ...interface{}) error {
	r, ok := result.(*string)
	if !ok {
		return errWrongType
	}
	if len(args) != 0 {
		return errWrongArgs
	}
	value := "0x4d4f434b"
	*r = value
	return nil
}

func mock_starknet_syncing(result interface{}, method string, args ...interface{}) error {
	r, ok := result.(*SyncResponse)
	if !ok {
		return errWrongType
	}
	if len(args) != 0 {
		return errWrongArgs
	}
	value := SyncResponse{
		StartingBlockNum: "0x1",
		CurrentBlockNum:  "0x1",
		HighestBlockNum:  "0x1",
	}
	*r = value
	return nil
}
