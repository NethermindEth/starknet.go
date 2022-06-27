package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

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
	case "starknet_getTransactionByBlockHashAndIndex":
		return mock_starknet_getTransactionByBlockHashAndIndex(result, method, args...)
	case "starknet_getTransactionByBlockNumberAndIndex":
		return mock_starknet_getTransactionByBlockNumberAndIndex(result, method, args...)
	case "starknet_getBlockTransactionCountByNumber":
		return mock_starknet_getBlockTransactionCountByNumber(result, method, args...)
	case "starknet_getBlockTransactionCountByHash":
		return mock_starknet_getBlockTransactionCountByHash(result, method, args...)
	case "starknet_getTransactionByHash":
		return mock_starknet_getTransactionByHash(result, method, args...)
	case "starknet_getTransactionReceipt":
		return mock_starknet_getTransactionReceipt(result, method, args...)
	case "starknet_getCode":
		return mock_starknet_getCode(result, method, args...)
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
	case "starknet_getStateUpdateByHash":
		return mock_starknet_getStateUpdateByHash(result, method, args...)
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

func mock_starknet_getTransactionByHash(result interface{}, method string, args ...interface{}) error {
	r, ok := result.(*json.RawMessage)
	if !ok || r == nil {
		return errWrongType
	}
	if len(args) != 1 {
		return errWrongArgs
	}
	txHash, ok := args[0].(string)
	if !ok || !strings.HasPrefix(txHash, "0x") {
		return errWrongArgs
	}
	transaction := types.Transaction{
		TransactionHash:    txHash,
		ContractAddress:    "0xdeadbeef",
		EntryPointSelector: "0xdeadbeef",
	}
	outputContent, _ := json.Marshal(transaction)
	json.Unmarshal(outputContent, r)
	return nil
}

func mock_starknet_getTransactionByBlockHashAndIndex(result interface{}, method string, args ...interface{}) error {
	r, ok := result.(*json.RawMessage)
	if !ok || r == nil {
		return errWrongType
	}
	if len(args) != 2 {
		return errWrongArgs
	}
	blockHash, ok := args[0].(string)
	if !ok || !strings.HasPrefix(blockHash, "0x") {
		return errWrongArgs
	}
	_, ok = args[1].(int)
	if !ok {
		fmt.Printf("args[1] expecting int, got %T\n", args[1])
		return errWrongArgs
	}
	transaction := types.Transaction{
		TransactionHash:    "0xdeadbeef",
		ContractAddress:    "0xdeadbeef",
		EntryPointSelector: "0xdeadbeef",
	}
	outputContent, _ := json.Marshal(transaction)
	json.Unmarshal(outputContent, r)
	return nil
}

func mock_starknet_getTransactionByBlockNumberAndIndex(result interface{}, method string, args ...interface{}) error {
	r, ok := result.(*json.RawMessage)
	if !ok || r == nil {
		return errWrongType
	}
	if len(args) != 2 {
		return errWrongArgs
	}
	_, ok1 := args[0].(int)
	_, ok2 := args[0].(string)
	if !ok1 && !ok2 {
		fmt.Printf("args[0] expecting int or string, got %T\n", args[0])
		return errWrongArgs
	}
	_, ok = args[1].(int)
	if !ok {
		fmt.Printf("args[1] expecting int, got %T\n", args[1])
		return errWrongArgs
	}
	transaction := types.Transaction{
		TransactionHash:    "0xdeadbeef",
		ContractAddress:    "0xdeadbeef",
		EntryPointSelector: "0xdeadbeef",
	}
	outputContent, _ := json.Marshal(transaction)
	json.Unmarshal(outputContent, r)
	return nil
}

func mock_starknet_getTransactionReceipt(result interface{}, method string, args ...interface{}) error {
	r, ok := result.(*json.RawMessage)
	if !ok || r == nil {
		return errWrongType
	}
	if len(args) != 1 {
		return errWrongArgs
	}
	txHash, ok := args[0].(string)
	if !ok || !strings.HasPrefix(txHash, "0x") {
		return errWrongArgs
	}
	transaction := types.TransactionReceipt{
		TransactionHash: txHash,
		Status:          "ACCEPTED_ON_L1",
	}
	outputContent, _ := json.Marshal(transaction)
	json.Unmarshal(outputContent, r)
	return nil
}

func mock_starknet_getCode(result interface{}, method string, args ...interface{}) error {
	r, ok := result.(*json.RawMessage)
	if !ok || r == nil {
		return errWrongType
	}
	if len(args) != 1 {
		return errWrongArgs
	}
	contractHash, ok := args[0].(string)
	if !ok || !strings.HasPrefix(contractHash, "0x") {
		return errWrongArgs
	}
	var codeRaw = struct {
		Bytecode []string `json:"bytecode"`
		AbiRaw   string   `json:"abi"`
	}{
		Bytecode: types.Bytecode{"0xdeadbeef"},
		AbiRaw:   `[{"name": "Uint256", "type": "Struct"}]`,
	}
	outputContent, _ := json.Marshal(codeRaw)
	json.Unmarshal(outputContent, r)
	return nil
}

func mock_starknet_getClassAt(result interface{}, method string, args ...interface{}) error {
	r, ok := result.(*json.RawMessage)
	if !ok || r == nil {
		fmt.Printf("%T\n", result)
		return errWrongType
	}
	if len(args) != 1 {
		return errWrongArgs
	}
	contractHash, ok := args[0].(string)
	if !ok || !strings.HasPrefix(contractHash, "0x") {
		return errWrongArgs
	}
	var class = types.ContractClass{
		Program: []string{"0xdeadbeef"},
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
	if len(args) != 1 {
		return errWrongArgs
	}
	contractHash, ok := args[0].(string)
	if !ok || !strings.HasPrefix(contractHash, "0x") {
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
		return errWrongArgs
	}
	var class = types.ContractClass{
		Program: []string{"0xdeadbeef"},
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
	query, ok := args[0].(EventParams)
	if !ok {
		return errWrongArgs
	}
	events := Events{
		Events: []Event{
			{
				BlockNumber: int(query.FromBlock),
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
	function, ok := args[0].(types.FunctionCall)
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
	_, ok = args[0].(FunctionCall)
	if !ok {
		fmt.Printf("args[0] should be FunctionCall, got %T\n", args[0])
		return errWrongArgs
	}
	_, ok = args[1].(string)
	if !ok {
		fmt.Printf("args[1] should be string, got %T\n", args[1])
		return errWrongArgs
	}

	output := EstimateFeeOutput{
		GasConsumed: "0xdeadbeef",
		GasPrice:    "0xdeadbeef",
		OverallFee:  "0xdeadbeef",
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
		fmt.Printf("args[0] should be types.FunctionCall, got %T\n", args[0])
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

func mock_starknet_getStateUpdateByHash(result interface{}, method string, args ...interface{}) error {
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
		BlockHash: blockHash,
	}
	outputContent, _ := json.Marshal(output)
	json.Unmarshal(outputContent, r)
	return nil
}

func mock_starknet_getBlockTransactionCountByHash(result interface{}, method string, args ...interface{}) error {
	r, ok := result.(*json.RawMessage)
	if !ok {
		return errWrongType
	}
	if len(args) != 1 {
		fmt.Printf("args: %d\n", len(args))
		return errWrongArgs
	}
	_, ok = args[0].(string)
	if !ok {
		fmt.Printf("args[0] should be string, got %T\n", args[0])
		return errWrongArgs
	}
	output := 7
	outputContent, _ := json.Marshal(output)
	json.Unmarshal(outputContent, r)
	return nil
}

func mock_starknet_getBlockTransactionCountByNumber(result interface{}, method string, args ...interface{}) error {
	r, ok := result.(*json.RawMessage)
	if !ok {
		return errWrongType
	}
	if len(args) != 1 {
		fmt.Printf("args: %d\n", len(args))
		return errWrongArgs
	}
	_, ok1 := args[0].(string)
	_, ok2 := args[0].(int)
	if !ok1 && !ok2 {
		fmt.Printf("args[0] should be int or string, got %T\n", args[0])
		return errWrongArgs
	}
	output := 7
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
