package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strings"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/utils"
	"github.com/pkg/errors"
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

// Close sets the 'closed' field of the rpcMock struct to true.
//
// No parameters.
// No return value.
func (r *rpcMock) Close() {
	r.closed = true
}

// CallContext calls the RPC method with the specified parameters and returns an error.
//
// Parameters:
// - ctx: represents the current execution context
// - result: the interface{} to store the result of the RPC call
// - method: the string representing the RPC method to be called
// - args: variadic and can be used to pass additional arguments to the RPC method
// Returns:
// - error: an error if any occurred during the function call
func (r *rpcMock) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	switch method {
	case "starknet_addDeclareTransaction":
		return mock_starknet_addDeclareTransaction(result, method, args...)
	case "starknet_addInvokeTransaction":
		return mock_starknet_addInvokeTransaction(result, method, args...)
	case "starknet_blockNumber":
		return mock_starknet_blockNumber(result, method, args...)
	case "starknet_call":
		return mock_starknet_call(result, method, args...)
	case "starknet_chainId":
		return mock_starknet_chainId(result, method, args...)
	case "starknet_estimateFee":
		return mock_starknet_estimateFee(result, method, args...)
	case "starknet_estimateMessageFee":
		return mock_starknet_estimateMessageFee(result, method, args...)
	case "starknet_getBlockTransactionCount":
		return mock_starknet_getBlockTransactionCount(result, method, args...)
	case "starknet_getBlockWithTxHashes":
		return mock_starknet_getBlockWithTxHashes(result, method, args...)
	case "starknet_getClass":
		return mock_starknet_getClass(result, method, args...)
	case "starknet_getClassAt":
		return mock_starknet_getClassAt(result, method, args...)
	case "starknet_getClassHashAt":
		return mock_starknet_getClassHashAt(result, method, args...)
	case "starknet_getEvents":
		return mock_starknet_getEvents(result, method, args...)
	case "starknet_getNonce":
		return mock_starknet_getNonce(result, method, args...)
	case "starknet_getStateUpdate":
		return mock_starknet_getStateUpdate(result, method, args...)
	case "starknet_getStorageAt":
		return mock_starknet_getStorageAt(result, method, args...)
	case "starknet_getTransactionByBlockIdAndIndex":
		return mock_starknet_getTransactionByBlockIdAndIndex(result, method, args...)
	case "starknet_getTransactionByHash":
		return mock_starknet_getTransactionByHash(result, method, args...)
	case "starknet_getTransactionReceipt":
		return mock_starknet_getTransactionReceipt(result, method, args...)
	case "starknet_syncing":
		return mock_starknet_syncing(result, method, args...)
	case "starknet_traceBlockTransactions":
		return mock_starknet_traceBlockTransactions(result, method, args...)
	case "starknet_traceTransaction":
		return mock_starknet_traceTransaction(result, method, args...)
	default:
		return errNotFound
	}
}

// mock_starknet_blockNumber is a function that mocks the blockNumber functionality in the StarkNet API.
//
// Parameters:
// - result: The result variable that will hold the block number value
// - method: The method string that specifies the API method being called
// - args: Additional arguments passed to the function
// Returns:
// - error: An error if the result is not of type *big.Int or if the arguments count is not zero
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

// mock_starknet_chainId is a function that mocks the behavior of the `starknet_chainId` method.
//
// Parameters:
// - result: an interface{} that holds the result of the function.
// - method: a string that represents the method.
// - args: a variadic parameter of type interface{} that represents the arguments of the function.
// Returns:
// - error: an error if there is a wrong type or wrong number of arguments.
func mock_starknet_chainId(result interface{}, method string, args ...interface{}) error {
	r, ok := result.(*string)
	if !ok {
		return errWrongType
	}
	if len(args) != 0 {
		return errWrongArgs
	}
	value := "0x534e5f474f45524c49"
	*r = value
	return nil
}

// mock_starknet_syncing is a function that mocks the behavior of the starknet_syncing function.
//
// Parameters:
// - result: an interface{} that holds the result of the function.
// - method: a string that represents the method.
// - args: a variadic parameter of type interface{} that represents the arguments of the function.
// Return:
// - error: an error if there is a wrong type or wrong number of arguments
func mock_starknet_syncing(result interface{}, method string, args ...interface{}) error {
	// Note: Since starknet_syncing returns with bool or SyncStatus, we pass in interface{}
	r, ok := result.(*interface{})
	if !ok {
		return errWrongType
	}
	if len(args) != 0 {
		return errWrongArgs
	}

	blockDataFeltArr, err := utils.HexArrToFelt([]string{
		"0x4b238e99c40d448b85dfc69e4414c2dbeb4d21d5c670b1662b5ad2ad2fcb061",
		"0x9cee6f457637180c36532bb0bfc5a091bb410b70f0489bcbbb0f1eca6650be",
	})
	if err != nil {
		return err
	}
	value := SyncStatus{
		StartingBlockHash: blockDataFeltArr[0],
		StartingBlockNum:  "0x4c602",
		CurrentBlockHash:  blockDataFeltArr[1],
		CurrentBlockNum:   "0x4c727",
		HighestBlockHash:  blockDataFeltArr[1],
		HighestBlockNum:   "0x4c727",
	}
	*r = value
	return nil
}

// mock_starknet_getTransactionByBlockIdAndIndex is a function that mocks the behavior of getting
// a transaction by block ID and index in the StarkNet API.
//
// Parameters:
// - result: The result of the API call, which will be stored in the provided interface{}. This should be a pointer to a json.RawMessage
// - method: The method of the API call
// - args: The arguments of the API call. This should be a variadic parameter that accepts a variable number of arguments
// Returns:
// - error: An error if the API call fails, otherwise nil
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

	var InvokeTxnV1example = `{
		"transaction_hash": "0x705547f8f2f8fdfb10ed533d909f76482bb293c5a32648d476774516a0bebd0",
		"type": "INVOKE",
		"nonce": "0x0",
		"max_fee": "0x53685de02fa5",
		"version": "0x1",
		"signature": [
		"0x4a7849de7b91e52cd0cdaf4f40aa67f54a58e25a15c60e034d2be819c1ecda4",
		"0x227fcad2a0007348e64384649365e06d41287b1887999b406389ee73c1d8c4c"
		],
		"sender_address": "0x315e364b162653e5c7b23efd34f8da27ba9c069b68e3042b7d76ce1df890313",
		"calldata": [
				   "0x1",
				   "0x13befe6eda920ce4af05a50a67bd808d67eee6ba47bb0892bef2d630eaf1bba"
		]
		}`

	json.Unmarshal([]byte(InvokeTxnV1example), r)
	return nil
}

// mock_starknet_getBlockTransactionCount is a function that mocks the behavior of the
// GetBlockTransactionCount method in the StarkNet API.
//
// Parameters:
// - result: The result of the API call, which will be stored in the provided interface{}. This should be a pointer to a json.RawMessage
// - method: The method of the API call
// - args: The arguments of the API call. This should be a variadic parameter that accepts a variable number of arguments
// Returns:
// - error: An error if the API call fails, otherwise nil
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

// mock_starknet_getTransactionByHash is a function that retrieves a transaction by its hash.
//
// Parameters:
// - result: an interface{} that represents the result of the transaction retrieval
// - method: a string that specifies the method used for retrieval
// - args: a variadic parameter that contains the arguments used for retrieval
// Returns:
// - error: an error if there is a failure in retrieving the transaction
func mock_starknet_getTransactionByHash(result interface{}, method string, args ...interface{}) error {
	r, ok := result.(*json.RawMessage)
	if !ok || r == nil {
		return errWrongType
	}
	if len(args) != 1 {
		return errWrongArgs
	}

	_, ok = args[0].(*felt.Felt)
	if !ok {
		return errWrongArgs
	}

	var InvokeTxnV1example = `    {
		"transaction_hash": "0x1779df1c6de5136ad2620f704b645e9cbd554b57d37f08a06ea60142269c5a5",
		"version": "0x1",
		"max_fee": "0x17970b794f000",
		"signature": [
		  "0xe500c4014c055c3304d8a125cfef44638ffa5b0f6840916049667a4c38aa1c",
		  "0x45ac538bfce5d8c5741b4421bbdc99f5849451acae75d2048d7cc4bb029ca77"
		],
		"nonce": "0x2d",
		"sender_address": "0x66dd340c03b6b7866fa7bb4bb91cc9e9c2a8eedc321985f334fd55de5e4e071",
		"calldata": [
		  "0x3",
		  "0x39a04b968d794fb076b0fbb146c12b48a23aa785e3d2e5be1982161f7536218",
		  "0x2f0b3c5710379609eb5495f1ecd348cb28167711b73609fe565a72734550354",
		  "0x0",
		  "0x3",
		  "0x3207980cd08767c9310d197c38b1a58b2a9bceb61dd9a99f51b407798702991",
		  "0x2f0b3c5710379609eb5495f1ecd348cb28167711b73609fe565a72734550354",
		  "0x3",
		  "0x3",
		  "0x42969068f9e84e9bf1c7bb6eb627455287e58f866ba39e45b123f9656aed5e9",
		  "0x2f0b3c5710379609eb5495f1ecd348cb28167711b73609fe565a72734550354",
		  "0x6",
		  "0x3",
		  "0x9",
		  "0x47487560da4c5c5755897e527a5fda37422b5ba02a2aba1ca3ce2b24dfd142e",
		  "0xde0b6b3a7640000",
		  "0x0",
		  "0x47487560da4c5c5755897e527a5fda37422b5ba02a2aba1ca3ce2b24dfd142e",
		  "0x10f0cf064dd59200000",
		  "0x0",
		  "0x47487560da4c5c5755897e527a5fda37422b5ba02a2aba1ca3ce2b24dfd142e",
		  "0x21e19e0c9bab2400000",
		  "0x0"
		],
		"type": "INVOKE"
	  }`

	json.Unmarshal([]byte(InvokeTxnV1example), r)
	return nil
}

// mock_starknet_getTransactionReceipt mocks the function that retrieves the transaction receipt information
// from the StarkNet blockchain.
//
// Parameters:
// - result: a pointer to an interface that will store the transaction receipt result
// - method: a string representing the method of the transaction receipt
// - args: a variadic parameter representing the arguments of the transaction receipt
// Returns:
// - error: an error if there is an issue with the type of the result or the number of arguments
func mock_starknet_getTransactionReceipt(result interface{}, method string, args ...interface{}) error {
	r, ok := result.(*json.RawMessage)
	if !ok || r == nil {
		return errWrongType
	}
	fmt.Printf("%T, %d", result, len(args))
	if len(args) != 1 {
		return errWrongArgs
	}

	arg0Felt, err := utils.HexToFelt(args[0].(string))
	if err != nil {
		return err
	}
	fromAddressFelt, err := utils.HexToFelt("0xdeadbeef")
	if err != nil {
		return err
	}
	transaction := InvokeTransactionReceipt(CommonTransactionReceipt{
		TransactionHash: arg0Felt,
		FinalityStatus:  TxnFinalityStatusAcceptedOnL1,
		Events: []Event{{
			FromAddress: fromAddressFelt,
		}},
	})
	outputContent, _ := json.Marshal(transaction)
	json.Unmarshal(outputContent, r)
	return nil
}

// mock_starknet_getClassAt is a function that performs a mock operation to get the class at a given index.
// The function sets the class to a DeprecatedContractClass with a specific program and marshals the class into JSON format.
// Finally, it unmarshals the JSON content into the result.
//
// Parameters:
// - result: An interface{} that represents the result of the operation
// - method: A string that specifies the method to be used
// - args: A variadic parameter that represents the arguments to be passed
// Returns:
// - error: An error if the result is not of type *json.RawMessage or is nil or the number of arguments is not equal to 2
// The function always returns nil.
func mock_starknet_getClassAt(result interface{}, method string, args ...interface{}) error {
	r, ok := result.(*json.RawMessage)
	if !ok || r == nil {
		fmt.Printf("%T\n", result)
		return errWrongType
	}
	if len(args) != 2 {
		return errWrongArgs
	}
	var class = DeprecatedContractClass{
		Program: "H4sIAAAAAAAE/+Vde3PbOJL/Kj5VXW1mVqsC36Sr9g8n0c6mzonnbM",
	}
	outputContent, _ := json.Marshal(class)
	json.Unmarshal(outputContent, r)
	return nil
}

// mock_starknet_getClassHashAt is a function that retrieves the class hash at a specific location in the StarkNet.
//
// Parameters:
// - result: An interface{} that represents the result of the operation
// - method: A string that specifies the method to be used
// - args: A variadic parameter that represents the arguments to be passed
// Returns:
// - error: An error if the result is not of type *json.RawMessage or is nil or the number of arguments is not equal to 2
// The function always returns nil.
func mock_starknet_getClassHashAt(result interface{}, method string, args ...interface{}) error {
	r, ok := result.(*json.RawMessage)
	if !ok || r == nil {
		fmt.Printf("%T\n", result)
		return errWrongType
	}
	if len(args) != 2 {
		return errWrongArgs
	}
	classHash, err := utils.HexToFelt("0xdeadbeef")
	if err != nil {
		return err
	}
	outputContent, _ := json.Marshal(classHash)
	json.Unmarshal(outputContent, r)
	return nil
}

// mock_starknet_getClass is a function that retrieves a class from the StarkNet API.
// It takes in a result interface{}, a method string, and variadic args ...interface{}.
// The result interface{} should be a pointer to json.RawMessage.
// The method string specifies the method to be called on the StarkNet API.
// The args ...interface{} are the arguments to be passed to the method.
// The function returns an error if any of the following conditions are met:
// - The result is not of type *json.RawMessage.
// - The args length is not equal to 2.
// - The first argument is not of type BlockID.
// - The second argument is not of type *felt.Felt or does not have a hexadecimal prefix.
// The function assigns a DeprecatedContractClass struct to the variable class.
// The function then marshals the class to JSON and unmarshals it to the result interface{}.
// If successful, the function returns nil.
//
// Parameters:
// - result: The result interface{} that should be a pointer to json.RawMessage
// - method: The method string specifying the method to be called on the StarkNet API
// - args: The variadic args ...interface{} representing the arguments to be passed to the method
// Returns:
// - error: An error if any of the conditions mentioned above are met
func mock_starknet_getClass(result interface{}, method string, args ...interface{}) error {
	r, ok := result.(*json.RawMessage)
	if !ok || r == nil {
		fmt.Printf("%T\n", result)
		return errWrongType
	}
	if len(args) != 2 {
		return errWrongArgs
	}
	_, ok = args[0].(BlockID)
	if !ok {
		fmt.Printf("expecting BlockID, instead %T\n", args[1])
		return errWrongArgs
	}
	classHash, ok := args[1].(*felt.Felt)
	if !ok || !strings.HasPrefix(classHash.String(), "0x") {
		fmt.Printf("%T\n", args[1])
		return errWrongArgs
	}
	var class = DeprecatedContractClass{
		Program: "H4sIAAAAAAAA",
	}
	outputContent, _ := json.Marshal(class)
	json.Unmarshal(outputContent, r)
	return nil
}

// mock_starknet_getEvents is a function that retrieves events from the StarkNet blockchain.
//
// Parameters:
// - result: An interface{} that represents the result of the operation
// - method: A string that specifies the method to be used
// - args: A variadic parameter that represents the arguments to be passed
// Returns:
// - error: An error if the result is not of type *json.RawMessage or is nil or the number of arguments is not equal to 1
// The function always returns nil
func mock_starknet_getEvents(result interface{}, method string, args ...interface{}) error {
	r, ok := result.(*json.RawMessage)
	if !ok {
		return errWrongType
	}
	if len(args) != 1 {
		return errWrongArgs
	}
	_, ok = args[0].(EventsInput)
	if !ok {
		return errWrongArgs
	}

	blockHash, err := utils.HexToFelt("0x59dbe64bf2e2f89f5f2958cff11044dca0c64dea2e37ec6eaad9a5f838793cb")
	if err != nil {
		return err
	}
	txHash, _ := utils.HexToFelt("0x568147c09d5e5db8dc703ce1da21eae47e9ad9c789bc2f2889c4413a38c579d")
	if err != nil {
		return err
	}

	events :=
		EventChunk{
			Events: []EmittedEvent{
				EmittedEvent{
					BlockHash:       blockHash,
					BlockNumber:     1472,
					TransactionHash: txHash,
				},
			},
		}

	outputContent, _ := json.Marshal(events)
	json.Unmarshal(outputContent, r)
	return nil
}

// mock_starknet_call is a function that mocks a Starknet call.
//
// Parameters:
// - result: The result of the transaction
// - method: The method to be called
// - args: The arguments to be passed to the method
// Returns:
// - error: An error if the transaction fails
func mock_starknet_call(result interface{}, method string, args ...interface{}) error {
	r, ok := result.(*json.RawMessage)
	if !ok {
		return errWrongType
	}
	if len(args) != 2 {
		fmt.Printf("args: %d\n", len(args))
		return errWrongArgs
	}
	out, err := new(felt.Felt).SetString("0xdeadbeef")
	if err != nil {
		return err
	}
	outputContent, _ := json.Marshal([]*felt.Felt{out})
	json.Unmarshal(outputContent, r)
	return nil
}

// mock_starknet_addDeclareTransaction is a mock function that adds a declare transaction to the StarkNet smart contract.
//
// Parameters:
// - result: The result of the transaction
// - method: The method to be called
// - args: The arguments to be passed to the method
// Return:
// - error: An error if the transaction fails
func mock_starknet_addDeclareTransaction(result interface{}, method string, args ...interface{}) error {
	r, ok := result.(*json.RawMessage)
	if !ok {
		return errWrongType
	}
	if len(args) != 2 {
		fmt.Printf("args: %d\n", len(args))
		return errWrongArgs
	}
	_, ok = args[0].(DeprecatedContractClass)
	if !ok {
		fmt.Printf("args[2] should be ContractClass, got %T\n", args[0])
		return errWrongArgs
	}
	_, ok = args[1].(string)
	if !ok {
		fmt.Printf("args[1] should be string, got %T\n", args[1])
		return errWrongArgs
	}
	deadbeefFelt, err := utils.HexToFelt("0xdeadbeef")
	if err != nil {
		return err
	}
	output := AddDeclareTransactionOutput{
		TransactionHash: deadbeefFelt,
		ClassHash:       deadbeefFelt,
	}
	outputContent, _ := json.Marshal(output)
	json.Unmarshal(outputContent, r)
	return nil
}

// mock_starknet_estimateFee simulates the estimation of a fee in the StarkNet network.
//
// Parameters:
// - result: The result of the transaction
// - method: The method to be called
// - args: The arguments to be passed to the method
// Returns:
// - error: an error if any
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
	_, ok = args[1].(BlockID)
	if !ok {
		fmt.Printf("args[1] should be *blockID, got %T\n", args[1])
		return errWrongArgs
	}

	gasCons, _ := new(felt.Felt).SetString("0x01a4")
	gasPrice, _ := new(felt.Felt).SetString("0x45")
	overallFee, _ := new(felt.Felt).SetString("0x7134")
	output := FeeEstimate{
		GasConsumed: gasCons,
		GasPrice:    gasPrice,
		OverallFee:  overallFee,
	}
	outputContent, _ := json.Marshal(output)
	json.Unmarshal(outputContent, r)
	return nil
}

// mock_starknet_estimateMessageFee is a function that estimates the fee for a StarkNet message.
//
// Parameters:
// - result: The result of the transaction
// - method: The method to be called
// - args: The arguments to be passed to the method
// Returns:
// - error: an error if any
func mock_starknet_estimateMessageFee(result interface{}, method string, args ...interface{}) error {
	r, ok := result.(*json.RawMessage)
	if !ok {
		return errWrongType
	}
	if len(args) != 2 {
		fmt.Printf("args: %d\n", len(args))
		return errWrongArgs
	}
	_, ok = args[0].(MsgFromL1)
	if !ok {
		fmt.Printf("args[0] should be MsgFromL1, got %T\n", args[0])
		return errWrongArgs
	}
	_, ok = args[1].(BlockID)
	if !ok {
		fmt.Printf("args[1] should be *blockID, got %T\n", args[1])
		return errWrongArgs
	}

	output := FeeEstimate{
		GasConsumed: new(felt.Felt).SetUint64(1),
		GasPrice:    new(felt.Felt).SetUint64(2),
		OverallFee:  new(felt.Felt).SetUint64(3),
	}
	outputContent, _ := json.Marshal(output)
	json.Unmarshal(outputContent, r)
	return nil
}

// mock_starknet_addInvokeTransaction is a mock function that simulates the behavior of the
// starknet_addInvokeTransaction function. It takes a result interface{}, a method string,
// and variadic args ...interface{} as parameters. The result parameter is expected to be of
// type *json.RawMessage. The method parameter represents the name of the method being invoked.
// The args parameter is a variadic argument, where the first argument is expected to be of
// type InvokeTxnV1.
//
// The function performs several checks and operations on the input parameters. It checks if the
// result parameter is of the correct type, and returns an error if it is not. It also checks if
// the number of arguments passed in the args parameter is exactly 1, and returns an error if it
// is not. The function then attempts to type cast the first argument in args to InvokeTxnV1 and
// returns an error if the type cast fails. It further checks if the SenderAddress field of the
// invokeTx object is not nil, and if it is equal to a predefined value. If it is, an unexpected
// error with a custom message is returned. The function then converts a hexadecimal value to a
// felt.Felt type and checks for any errors during the conversion. Finally, the function creates
// an AddInvokeTransactionResponse object, marshals it into JSON format, and unmarshals it into
// the result parameter.
//
// Parameters:
// - result: The result of the transaction
// - method: The method to be called
// - args: The arguments to be passed to the method
// Returns:
// - error: an error if any
func mock_starknet_addInvokeTransaction(result interface{}, method string, args ...interface{}) error {
	r, ok := result.(*json.RawMessage)
	if !ok {
		return errWrongType
	}
	if len(args) != 1 {
		return errors.Wrap(errWrongArgs, fmt.Sprint("wrong number of args ", len(args)))
	}
	invokeTx, ok := args[0].(InvokeTxnV1)
	if !ok {
		return errors.Wrap(errWrongArgs, fmt.Sprintf("args[0] should be InvokeTxnV1, got %T\n", args[0]))
	}
	if invokeTx.SenderAddress != nil {

		if invokeTx.SenderAddress.Equal(new(felt.Felt).SetUint64(123)) {
			unexpErr := *ErrUnexpectedError
			unexpErr.data = "Something crazy happened"
			return &unexpErr
		}
	}
	deadbeefFelt, err := utils.HexToFelt("0xdeadbeef")
	if err != nil {
		return err
	}
	output := AddInvokeTransactionResponse{
		TransactionHash: deadbeefFelt,
	}
	outputContent, _ := json.Marshal(output)
	json.Unmarshal(outputContent, r)
	return nil
}

// mock_starknet_getStorageAt mocks the behavior of the StarkNet getStorageAt function.
//
// Parameters:
// - result: The result of the transaction
// - method: The method to be called
// - args: The arguments to be passed to the method
// Returns:
// - error: an error if any
func mock_starknet_getStorageAt(result interface{}, method string, args ...interface{}) error {
	r, ok := result.(*json.RawMessage)
	if !ok {
		return errWrongType
	}
	if len(args) != 3 {
		fmt.Printf("args: %d\n", len(args))
		return errWrongArgs
	}

	if _, ok := args[0].(*felt.Felt); !ok {
		return errWrongArgs
	}

	if _, ok := args[1].(string); !ok {
		return errWrongArgs
	}

	if _, ok := args[2].(BlockID); !ok {
		return errWrongArgs
	}

	output := "0xdeadbeef"
	outputContent, _ := json.Marshal(output)
	json.Unmarshal(outputContent, r)
	return nil
}

// mock_starknet_getStateUpdate is a function that performs a mock operation to get the state update.
//
// Parameters:
// - result: an interface{} that represents the result of the state update.
// - method: a string that specifies the method used to retrieve the state update.
// - args: a variadic parameter that can accept multiple arguments.
// Returns:
// - error: an error if any
func mock_starknet_getStateUpdate(result interface{}, method string, args ...interface{}) error {

	r, ok := result.(*json.RawMessage)
	if !ok {
		return errWrongType
	}
	if len(args) != 1 {
		fmt.Printf("args: %d\n", len(args))
		return errWrongArgs
	}
	_, ok = args[0].(BlockID)
	if !ok {
		fmt.Printf("args[1] should be *blockID, got %T\n", args[0])
		return errWrongArgs
	}

	stateFeltArr, err := utils.HexArrToFelt([]string{
		"0x4f1cee281edb6cb31b9ba5a8530694b5527cf05c5ac6502decf3acb1d0cec4",
		"0x70677cda9269d47da3ff63bc87cf1c87d0ce167b05da295dc7fc68242b250b",
		"0x19aa982a75263d4c4de4cc4c5d75c3dec32e00b95bef7bbb4d17762a0b138af",
		"0xe5cc6f2b6d34979184b88334eb64173fe4300cab46ecd3229633fcc45c83d4",
		"0x1813aac5f5e7799684c6dc33e51f44d3627fd748c800724a184ed5be09b713e",
		"0x1813aac5f5e7799684c6dc33e51f44d3627fd748c800724a184ed5be09b713e",
		"0x630b4197"})
	if err != nil {
		return err
	}

	output := StateUpdateOutput{
		BlockHash: stateFeltArr[0],
		NewRoot:   stateFeltArr[1],
		PendingStateUpdate: PendingStateUpdate{
			OldRoot: stateFeltArr[2],
			StateDiff: StateDiff{
				StorageDiffs: []ContractStorageDiffItem{{
					Address: stateFeltArr[3],
					StorageEntries: []StorageEntry{
						{
							Key:   stateFeltArr[4],
							Value: stateFeltArr[5],
						},
					},
				}},
			},
		},
	}
	outputContent, _ := json.Marshal(output)
	json.Unmarshal(outputContent, r)
	return nil
}

// mock_starknet_getNonce is a function that retrieves the nonce for a given method and arguments.
//
// Parameters:
// - result: a pointer to the variable where the result will be stored
// - method: the method for which the nonce is being retrieved
// - args: optional arguments for the method
// Returns:
// - error: an error if
//   - The result parameter is not of type *json.RawMessage
//   - The number of arguments is not equal to 2
//   - The first argument is not of type BlockID
//   - The second argument is not of type *felt.Felt
func mock_starknet_getNonce(result interface{}, method string, args ...interface{}) error {
	r, ok := result.(*json.RawMessage)
	if !ok {
		return errWrongType
	}
	if len(args) != 2 {
		fmt.Printf("args: %d\n", len(args))
		return errWrongArgs
	}
	if _, ok := args[0].(BlockID); !ok {
		fmt.Printf("args[0] should be BlockID, got %T\n", args[0])
		return errWrongArgs
	}
	if _, ok := args[1].(*felt.Felt); !ok {
		fmt.Printf("args[0] should be *felt.Felt, got %T\n", args[1])
		return errWrongArgs
	}
	output, err := utils.HexToFelt("0x0")
	if err != nil {
		return err
	}
	outputContent, _ := json.Marshal(output)
	json.Unmarshal(outputContent, r)
	return nil
}

// mock_starknet_getBlockWithTxHashes mocks the behavior of the starknet_getBlockWithTxHashes function.
// If successful, it populates the result parameter with the json.RawMessage containing the block with the specified transaction hashes.
//
// Parameters:
// - result: the result is expected to be a pointer to json.RawMessage
// - method: the method to be called
// - args: variadic parameter that can contain any number of arguments
// Returns:
// - error: an error if any
func mock_starknet_getBlockWithTxHashes(result interface{}, method string, args ...interface{}) error {
	r, ok := result.(*json.RawMessage)
	if !ok || r == nil {
		return errWrongType
	}
	if len(args) != 1 {
		return errWrongArgs
	}
	blockId, ok := args[0].(BlockID)
	if !ok {
		fmt.Printf("args[0] should be BlockID, got %T\n", args[0])
		return errWrongArgs
	}

	txHashes, err := utils.HexArrToFelt([]string{
		"0x40c82f79dd2bc1953fc9b347a3e7ab40fe218ed5740bf4e120f74e8a3c9ac99",
		"0x28981b14353a28bc46758dff412ac544d16f2ffc8dde31867855592ea054ab1",
	})
	if err != nil {
		return err
	}

	if blockId.Tag == "latest" {
		pBlock, err := json.Marshal(
			PendingBlockTxHashes{
				PendingBlockHeader{
					ParentHash:       &felt.Zero,
					Timestamp:        123,
					SequencerAddress: &felt.Zero},
				txHashes,
			})
		if err != nil {
			return err
		}
		err = json.Unmarshal(pBlock, &r)
		if err != nil {
			return err
		}
	} else {
		blockHash, err := utils.HexToFelt("0xbeef")
		if err != nil {
			return err
		}
		block, err := json.Marshal(BlockTxHashes{
			BlockHeader: BlockHeader{
				BlockHash:        blockHash,
				ParentHash:       &felt.Zero,
				Timestamp:        124,
				SequencerAddress: &felt.Zero},
			Status:       BlockStatus_AcceptedOnL1,
			Transactions: txHashes,
		})
		if err != nil {
			return err
		}
		json.Unmarshal(block, &r)
	}

	return nil
}

// mock_starknet_traceBlockTransactions is a function that traces the transactions of a block in the StarkNet network.
// The function first checks the type of the result parameter and returns an error if it is not of type *json.RawMessage.
// It then checks the length of the args parameter and returns an error if it is not equal to 1. Next, it checks the
// type of the first element of args and returns an error if it is not of type *felt.Felt. If the block hash is equal
// to "0x3ddc3a8aaac071ecdc5d8d0cfbb1dc4fc6a88272bc6c67523c9baaee52a5ea2", the function reads the trace from a file
// and unmarshals it into a struct. It then marshals the result and unmarshals it into the result parameter.
// If the block hash is not valid, the function returns an error of type ErrInvalidBlockHash.
//
// Parameters:
// - result: a pointer to the variable where the result will be stored
// - method: the method for which the nonce is being retrieved
// - args: optional arguments for the method
// Returns:
// - error: an error if any
func mock_starknet_traceBlockTransactions(result interface{}, method string, args ...interface{}) error {
	r, ok := result.(*json.RawMessage)
	if !ok || r == nil {
		return errWrongType
	}
	if len(args) != 1 {
		return errWrongArgs
	}
	blockID, ok := args[0].(BlockID)
	if !ok {
		return errors.Wrap(errWrongArgs, fmt.Sprintf("args[0] should be BlockID, got %T\n", args[0]))
	}
	if blockID.Hash.String() == "0x3ddc3a8aaac071ecdc5d8d0cfbb1dc4fc6a88272bc6c67523c9baaee52a5ea2" {

		var rawBlockTrace struct {
			Result []Trace `json:"result"`
		}
		read, err := os.ReadFile("tests/trace/0x3ddc3a8aaac071ecdc5d8d0cfbb1dc4fc6a88272bc6c67523c9baaee52a5ea2.json")
		if err != nil {
			return err
		}
		if nil != json.Unmarshal(read, &rawBlockTrace) {
			return err
		}
		BlockTrace, err := json.Marshal(rawBlockTrace.Result)
		if err != nil {
			return err
		}
		return json.Unmarshal(BlockTrace, &r)
	}

	return ErrInvalidBlockHash
}

// mock_starknet_traceTransaction is a Go function that traces a transaction in the StarkNet network.
// The function returns an error if any of the following conditions are met:
// - The result is not of type *json.RawMessage.
// - The result is nil.
// - The number of arguments is not equal to 1.
// - The first argument is not of type *felt.Felt.
// - The transaction hash does not match any known hash.
//
// If the transaction hash matches "0xff66e14fc6a96f3289203690f5f876cb4b608868e8549b5f6a90a21d4d6329",
// the function reads the trace from a file and unmarshals it into the result.
//
// If the transaction hash matches "0xf00d", the function returns a custom RPCError.
//
// If the transaction hash does not match any known hash, the function returns ErrInvalidTxnHash.
//
// Parameters:
// - result: an interface{} that represents the result of the transaction.
// - method: a string that specifies the method used in the transaction.
// - args: a variadic parameter that can accept multiple arguments.
// Returns:
// - error: an error if any
func mock_starknet_traceTransaction(result interface{}, method string, args ...interface{}) error {
	r, ok := result.(*json.RawMessage)
	if !ok || r == nil {
		return errWrongType
	}
	if len(args) != 1 {
		return errWrongArgs
	}
	transactionHash, ok := args[0].(*felt.Felt)
	if !ok {
		return errors.Wrap(errWrongArgs, fmt.Sprintf("args[0] should be felt, got %T\n", args[0]))
	}
	switch transactionHash.String() {
	case "0xff66e14fc6a96f3289203690f5f876cb4b608868e8549b5f6a90a21d4d6329":
		var rawTrace struct {
			Result InvokeTxnTrace `json:"result"`
		}
		read, err := os.ReadFile("tests/trace/0xff66e14fc6a96f3289203690f5f876cb4b608868e8549b5f6a90a21d4d6329.json")
		if err != nil {
			return err
		}
		if nil != json.Unmarshal(read, &rawTrace) {
			return err
		}
		txnTrace, err := json.Marshal(rawTrace.Result)
		if err != nil {
			return err
		}
		return json.Unmarshal(txnTrace, &r)
	case "0xf00d":
		return &RPCError{
			code:    10,
			message: "No trace available for transaction",
			data:    "REJECTED",
		}
	default:
		return ErrInvalidTxnHash
	}
}
