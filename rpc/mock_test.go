package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
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
	case "starknet_getTransactionReceipt":
		return mock_starknet_getTransactionReceipt(result, method, args...)
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
	case "starknet_protocolVersion":
		return mock_starknet_protocolVersion(result, method, args...)
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

func mock_starknet_getTransactionByHash(result interface{}, method string, args ...interface{}) error {
	r, ok := result.(*json.RawMessage)
	if !ok || r == nil {
		return errWrongType
	}
	if len(args) != 1 {
		return errWrongArgs
	}
	txHash, ok := args[0].(TxnHash)
	if !ok || !strings.HasPrefix(string(txHash), "0x") {
		return errWrongArgs
	}
	outputContent, _ := json.Marshal(InvokeTxnV00x705547f8f2f8f)
	json.Unmarshal(outputContent, r)
	return nil
}

func mock_starknet_getTransactionReceipt(result interface{}, method string, args ...interface{}) error {
	r, ok := result.(*json.RawMessage)
	if !ok || r == nil {
		return errWrongType
	}
	fmt.Printf("%T, %d", result, len(args))
	if len(args) != 1 {
		return errWrongArgs
	}
	txHash, ok := args[0].(TxnHash)
	if !ok || !strings.HasPrefix(string(txHash), "0x") {
		return errWrongArgs
	}
	transaction := InvokeTxnReceipt{
		CommonReceiptProperties{
			TransactionHash: txHash,
			Status:          TxnStatus("ACCEPTED_ON_L1"),
		},
		&InvokeTxnReceiptProperties{
			Events: []Event{{
				FromAddress: Address("0xdeadbeef"),
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
	contractHash, ok := args[1].(Address)
	if !ok || !strings.HasPrefix(string(contractHash), "0x") {
		return errWrongArgs
	}
	var class = ContractClass{
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
	contractHash, ok := args[1].(Address)
	if !ok || !strings.HasPrefix(string(contractHash), "0x") {
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
	if len(args) != 2 {
		return errWrongArgs
	}
	classHash, ok := args[1].(string)
	if !ok || !strings.HasPrefix(classHash, "0x") {
		fmt.Printf("%T\n", args[1])
		return errWrongArgs
	}
	var class = ContractClass{
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
	query, ok := args[0].(EventFilterParams)
	if !ok {
		return errWrongArgs
	}
	events := &EventsOutput{
		Events: []EmittedEvent{
			{BlockHash: BlockHash("0xdeadbeef"),
				Event: Event{
					FromAddress: query.EventFilter.Address,
				},
				BlockNumber:     BlockNumber(1),
				TransactionHash: TxnHash("deadbeef"),
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
	function, ok := args[0].(FunctionCall)
	if !ok || function.ContractAddress != "0xdeadbeef" {
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
	_, ok = args[0].(ContractClass)
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

	_, ok = args[2].(ContractClass)
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
	_, ok = args[0].(FunctionCall)
	if !ok {
		fmt.Printf("args[0] should be FunctionCall, got %T\n", args[0])
		return errWrongArgs
	}
	_, ok = args[1].(*blockID)
	if !ok {
		fmt.Printf("args[1] should be *blockID, got %T\n", args[1])
		return errWrongArgs
	}

	output := FeeEstimate{
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
	_, ok = args[0].(FunctionCall)
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
	for i := range []int{1, 2, 3} {
		_, ok = args[i].(string)
		if !ok {
			fmt.Printf("args[%d] should be string, got %T\n", i, args[i])
			return errWrongArgs
		}
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
	blockHash, ok := args[0].(string)
	if !ok {
		fmt.Printf("args[0] should be string, got %T\n", args[0])
		return errWrongArgs
	}
	output := &StateUpdateOutput{
		BlockHash: BlockHash(blockHash),
	}
	outputContent, _ := json.Marshal(output)
	json.Unmarshal(outputContent, r)
	return nil
}

func mock_starknet_protocolVersion(result interface{}, method string, args ...interface{}) error {
	r, ok := result.(*json.RawMessage)
	if !ok {
		return errWrongType
	}
	if len(args) != 0 {
		fmt.Printf("args: %d\n", len(args))
		return errWrongArgs
	}
	output := "0x312e30"
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
	if _, ok := args[0].(string); !ok {
		fmt.Printf("args[0] should be string, got %T\n", args[0])
		return errWrongArgs
	}
	output := big.NewInt(10)
	outputContent, _ := json.Marshal(output)
	json.Unmarshal(outputContent, r)
	return nil
}
