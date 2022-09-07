package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"github.com/dontpanicdao/caigo/rpc/types"
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
	case "starknet_getTransactionByHash":
		return mock_starknet_getTransactionByHash(result, method, args...)
	case "starknet_getTransactionByBlockIdAndIndex":
		return mock_starknet_getTransactionByBlockIdAndIndex(result, method, args...)
	case "starknet_getBlockTransactionCount":
		return mock_starknet_getBlockTransactionCount(result, method, args...)
	case "starknet_getTxnReceipt":
		return mock_starknet_getTxnReceipt(result, method, args...)
	case "starknet_getClassAt":
		return mock_starknet_getClassAt(result, method, args...)
	case "starknet_getClassHashAt":
		return mock_starknet_getClassHashAt(result, method, args...)
	case "starknet_getClass":
		return mock_starknet_getClass(result, method, args...)
	case "starknet_getEvents":
		return mock_starknet_getEvents(result, method, args...)
	case "starknet_getNonce":
		return mock_starknet_getNonce(result, method, args...)
	case "starknet_getStorageAt":
		return mock_starknet_getStorageAt(result, method, args...)
	case "starknet_getStateUpdate":
		return mock_starknet_getStateUpdate(result, method, args...)
	case "starknet_call":
		return mock_starknet_call(result, method, args...)
	case "starknet_addDeclareTransaction":
		return mock_starknet_addDeclareTransaction(result, method, args...)
	case "starknet_addDeployTransaction":
		return mock_starknet_addDeployTransaction(result, method, args...)
	case "starknet_addInvokeTransaction":
		return mock_starknet_addInvokeTransaction(result, method, args...)
	case "starknet_estimateFee":
		return mock_starknet_estimateFee(result, method, args...)
	default:
		return errNotFound
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
		fmt.Printf("%T\n", result)
		return errWrongType
	}
	if len(args) != 0 {
		fmt.Println(args...)
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
		StartingBlockHash: "0x4b238e99c40d448b85dfc69e4414c2dbeb4d21d5c670b1662b5ad2ad2fcb061",
		StartingBlockNum:  "0x4c602",
		CurrentBlockHash:  "0x9cee6f457637180c36532bb0bfc5a091bb410b70f0489bcbbb0f1eca6650be",
		CurrentBlockNum:   "0x4c727",
		HighestBlockHash:  "0x9cee6f457637180c36532bb0bfc5a091bb410b70f0489bcbbb0f1eca6650be",
		HighestBlockNum:   "0x4c727",
	}
	*r = value
	return nil
}

func mock_starknet_getTransactionByBlockIdAndIndex(result interface{}, method string, args ...interface{}) error {
	r, ok := result.(*json.RawMessage)
	if !ok || r == nil {
		return errWrongType
	}
	if len(args) != 2 {
		return errWrongArgs
	}
	_, ok = args[1].(uint64)
	if !ok {
		return errWrongArgs
	}
	outputContent, _ := json.Marshal(InvokeTxnV0_300000_0)
	json.Unmarshal(outputContent, r)
	return nil
}

func mock_starknet_getBlockTransactionCount(result interface{}, method string, args ...interface{}) error {
	r, ok := result.(*json.RawMessage)
	if !ok || r == nil {
		return errWrongType
	}
	if len(args) != 1 {
		return errWrongArgs
	}
	outputContent, _ := json.Marshal(uint64(10))
	json.Unmarshal(outputContent, r)
	return nil
}

func mock_starknet_getTransactionByHash(result interface{}, method string, args ...interface{}) error {
	r, ok := result.(*json.RawMessage)
	if !ok || r == nil {
		return errWrongType
	}
	if len(args) != 1 {
		return errWrongArgs
	}

	_, ok = args[0].(types.Hash)
	if !ok {
		return errWrongArgs
	}
	outputContent, _ := json.Marshal(InvokeTxnV00x705547f8f2f8f)
	json.Unmarshal(outputContent, r)
	return nil
}

func mock_starknet_getTxnReceipt(result interface{}, method string, args ...interface{}) error {
	r, ok := result.(*json.RawMessage)
	if !ok || r == nil {
		return errWrongType
	}
	fmt.Printf("%T, %d", result, len(args))
	if len(args) != 1 {
		return errWrongArgs
	}

	transaction := types.InvokeTransactionReceipt{
		CommonTransactionReceipt: types.CommonTransactionReceipt{
			TransactionHash: types.HexToHash(args[0].(string)),
			Status:          types.TransactionStatus("ACCEPTED_ON_L1"),
		},
		InvokeTransactionReceiptProperties: types.InvokeTransactionReceiptProperties{
			Events: []types.Event{{
				FromAddress: types.HexToHash("0xdeadbeef"),
			}},
		},
	}
	outputContent, _ := json.Marshal(transaction)
	json.Unmarshal(outputContent, r)
	return nil
}

func mock_starknet_getClassAt(result interface{}, method string, args ...interface{}) error {
	r, ok := result.(*json.RawMessage)
	if !ok || r == nil {
		fmt.Printf("%T\n", result)
		return errWrongType
	}
	if len(args) != 2 {
		return errWrongArgs
	}
	var class = types.ContractClass{
		Program: "H4sIAAAAAAAE/+Vde3PbOJL/Kj5VXW1mVqsC36Sr9g8n0c6mzonnbM",
	}
	outputContent, _ := json.Marshal(class)
	json.Unmarshal(outputContent, r)
	return nil
}

func mock_starknet_getClassHashAt(result interface{}, method string, args ...interface{}) error {
	r, ok := result.(*json.RawMessage)
	if !ok || r == nil {
		fmt.Printf("%T\n", result)
		return errWrongType
	}
	if len(args) != 2 {
		return errWrongArgs
	}
	classHash := "0xdeadbeef"
	outputContent, _ := json.Marshal(classHash)
	json.Unmarshal(outputContent, r)
	return nil
}

func mock_starknet_getClass(result interface{}, method string, args ...interface{}) error {
	r, ok := result.(*json.RawMessage)
	if !ok || r == nil {
		fmt.Printf("%T\n", result)
		return errWrongType
	}
	if len(args) != 1 {
		return errWrongArgs
	}
	classHash, ok := args[0].(string)
	if !ok || !strings.HasPrefix(classHash, "0x") {
		fmt.Printf("%T\n", args[1])
		return errWrongArgs
	}
	var class = types.ContractClass{
		Program: "H4sIAAAAAAAE/+Vde3PbOJL/Kj5VXW1mVqsC36Sr9g8n0c6mzonnbM",
	}
	outputContent, _ := json.Marshal(class)
	json.Unmarshal(outputContent, r)
	return nil
}

func mock_starknet_getEvents(result interface{}, method string, args ...interface{}) error {
	r, ok := result.(*json.RawMessage)
	if !ok {
		return errWrongType
	}
	if len(args) != 1 {
		return errWrongArgs
	}
	query, ok := args[0].(types.EventFilter)
	if !ok {
		return errWrongArgs
	}
	events := &types.EventsOutput{
		Events: []types.EmittedEvent{
			{BlockHash: types.HexToHash("0xdeadbeef"),
				Event: types.Event{
					FromAddress: query.Address,
				},
				BlockNumber:     1,
				TransactionHash: types.HexToHash("0xdeadbeef"),
			},
		},
	}
	outputContent, _ := json.Marshal(events)
	json.Unmarshal(outputContent, r)
	return nil
}

func mock_starknet_call(result interface{}, method string, args ...interface{}) error {
	r, ok := result.(*json.RawMessage)
	if !ok {
		return errWrongType
	}
	if len(args) != 2 {
		fmt.Printf("args: %d\n", len(args))
		return errWrongArgs
	}
	output := []string{"0x12"}
	outputContent, _ := json.Marshal(output)
	json.Unmarshal(outputContent, r)
	return nil
}

func mock_starknet_addDeclareTransaction(result interface{}, method string, args ...interface{}) error {
	r, ok := result.(*json.RawMessage)
	if !ok {
		return errWrongType
	}
	if len(args) != 2 {
		fmt.Printf("args: %d\n", len(args))
		return errWrongArgs
	}
	_, ok = args[0].(types.ContractClass)
	if !ok {
		fmt.Printf("args[2] should be ContractClass, got %T\n", args[0])
		return errWrongArgs
	}
	_, ok = args[1].(string)
	if !ok {
		fmt.Printf("args[1] should be string, got %T\n", args[1])
		return errWrongArgs
	}
	output := AddDeclareTransactionOutput{
		TransactionHash: "0xdeadbeef",
		ClassHash:       "0xdeadbeef",
	}
	outputContent, _ := json.Marshal(output)
	json.Unmarshal(outputContent, r)
	return nil
}

func mock_starknet_addDeployTransaction(result interface{}, method string, args ...interface{}) error {
	r, ok := result.(*json.RawMessage)
	if !ok {
		return errWrongType
	}
	if len(args) != 3 {
		fmt.Printf("args: %d\n", len(args))
		return errWrongArgs
	}
	_, ok = args[0].(string)
	if !ok {
		fmt.Printf("args[0] should be string, got %T\n", args[0])
		return errWrongArgs
	}
	_, ok = args[1].([]string)
	if !ok {
		fmt.Printf("args[1] should be ConstructorCallData, got %T\n", args[1])
		return errWrongArgs
	}

	_, ok = args[2].(types.ContractClass)
	if !ok {
		fmt.Printf("args[2] should be ContractClass, got %T\n", args[2])
		return errWrongArgs
	}

	output := AddDeployTransactionOutput{
		TransactionHash: "0xdeadbeef",
		ContractAddress: "0xdeadbeef",
	}
	outputContent, _ := json.Marshal(output)
	json.Unmarshal(outputContent, r)
	return nil
}

func mock_starknet_estimateFee(result interface{}, method string, args ...interface{}) error {
	r, ok := result.(*json.RawMessage)
	if !ok {
		return errWrongType
	}
	if len(args) != 2 {
		fmt.Printf("args: %d\n", len(args))
		return errWrongArgs
	}
	_, ok = args[0].(types.FunctionCall)
	if !ok {
		fmt.Printf("args[0] should be FunctionCall, got %T\n", args[0])
		return errWrongArgs
	}
	_, ok = args[1].(types.BlockID)
	if !ok {
		fmt.Printf("args[1] should be *blockID, got %T\n", args[1])
		return errWrongArgs
	}

	output := types.FeeEstimate{
		GasConsumed: "0x01a4",
		GasPrice:    "0x45",
		OverallFee:  "0x7134",
	}
	outputContent, _ := json.Marshal(output)
	json.Unmarshal(outputContent, r)
	return nil
}

func mock_starknet_addInvokeTransaction(result interface{}, method string, args ...interface{}) error {
	r, ok := result.(*json.RawMessage)
	if !ok {
		return errWrongType
	}
	if len(args) != 4 {
		fmt.Printf("args: %d\n", len(args))
		return errWrongArgs
	}
	_, ok = args[0].(types.FunctionCall)
	if !ok {
		fmt.Printf("args[0] should be FunctionCall, got %T\n", args[0])
		return errWrongArgs
	}
	_, ok = args[1].([]string)
	if !ok {
		fmt.Printf("args[1] should be []string, got %T\n", args[1])
		return errWrongArgs
	}
	_, ok = args[2].(string)
	if !ok {
		fmt.Printf("args[2] should be []string, got %T\n", args[2])
		return errWrongArgs
	}
	_, ok = args[3].(string)
	if !ok {
		fmt.Printf("args[3] should be []string, got %T\n", args[3])
		return errWrongArgs
	}

	output := AddInvokeTransactionOutput{
		TransactionHash: "0xdeadbeef",
	}
	outputContent, _ := json.Marshal(output)
	json.Unmarshal(outputContent, r)
	return nil
}

func mock_starknet_getStorageAt(result interface{}, method string, args ...interface{}) error {
	r, ok := result.(*json.RawMessage)
	if !ok {
		return errWrongType
	}
	if len(args) != 3 {
		fmt.Printf("args: %d\n", len(args))
		return errWrongArgs
	}

	if _, ok := args[0].(types.Hash); !ok {
		return errWrongArgs
	}

	if _, ok := args[1].(string); !ok {
		return errWrongArgs
	}

	if _, ok := args[2].(types.BlockID); !ok {
		return errWrongArgs
	}

	output := "0xdeadbeef"
	outputContent, _ := json.Marshal(output)
	json.Unmarshal(outputContent, r)
	return nil
}

func mock_starknet_getStateUpdate(result interface{}, method string, args ...interface{}) error {

	r, ok := result.(*json.RawMessage)
	if !ok {
		return errWrongType
	}
	if len(args) != 1 {
		fmt.Printf("args: %d\n", len(args))
		return errWrongArgs
	}
	_, ok = args[0].(types.BlockID)
	if !ok {
		fmt.Printf("args[1] should be *blockID, got %T\n", args[0])
		return errWrongArgs
	}

	output := types.StateUpdateOutput{
		BlockHash:    types.HexToHash("0x4f1cee281edb6cb31b9ba5a8530694b5527cf05c5ac6502decf3acb1d0cec4"),
		NewRoot:      "0x70677cda9269d47da3ff63bc87cf1c87d0ce167b05da295dc7fc68242b250b",
		OldRoot:      "0x19aa982a75263d4c4de4cc4c5d75c3dec32e00b95bef7bbb4d17762a0b138af",
		AcceptedTime: 0,
		StateDiff: types.StateDiff{
			StorageDiffs: []types.ContractStorageDiffItem{{
				Address: "0xe5cc6f2b6d34979184b88334eb64173fe4300cab46ecd3229633fcc45c83d4",
				Entries: []types.StorageEntry{{
					Key:   "0x1813aac5f5e7799684c6dc33e51f44d3627fd748c800724a184ed5be09b713e",
					Value: "0x630b4197",
				}},
			}},
		},
	}
	outputContent, _ := json.Marshal(output)
	json.Unmarshal(outputContent, r)
	return nil
}

func mock_starknet_getNonce(result interface{}, method string, args ...interface{}) error {
	r, ok := result.(*json.RawMessage)
	if !ok {
		return errWrongType
	}
	if len(args) != 1 {
		fmt.Printf("args: %d\n", len(args))
		return errWrongArgs
	}
	if _, ok := args[0].(types.Hash); !ok {
		fmt.Printf("args[0] should be string, got %T\n", args[0])
		return errWrongArgs
	}
	output := "0x0"
	outputContent, _ := json.Marshal(output)
	json.Unmarshal(outputContent, &r)
	return nil
}
