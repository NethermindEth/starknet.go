package rpc

import (
	"context"
	"fmt"
	"math/big"
)

var (
	errWrongType = fmt.Errorf("wrong type")
	errWrongArgs = fmt.Errorf("wrong number of args")
)

// rpcMock is a mock of the go-ethereum Client that can be used for local tests
// when no integration environment exists.
type rpcMock struct{}

func (r *rpcMock) Close() {}

func (r *rpcMock) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	switch method {
	case "starknet_blockNumber":
		return mock_starknet_blockNumber(result, method, args...)
	case "starknet_chainId":
		return mock_starknet_chainId(result, method, args...)
	case "starknet_syncing":
		return mock_starknet_syncing(result, method, args...)
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
