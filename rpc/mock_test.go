package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/contracts"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/pkg/errors"
)

var (
	errWrongType = errors.New("wrong type")
	errWrongArgs = errors.New("wrong number of args")
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

// Call 'CallContext' with a slice of arguments.
//
// For RPC-Calls with optional arguments, use 'CallContext' instead and pass a struct containing
// the arguments, because Juno doesn't support optional arguments being passed in an array, only within an object.
//
// Parameters:
//   - ctx: represents the current execution context
//   - result: the interface{} to store the result of the RPC call
//   - method: the name of the method to call
//   - args: variadic and can be used to pass additional arguments to the RPC method
//
// Returns:
//   - error: an error if any occurred during the function call
func (r *rpcMock) CallContextWithSliceArgs(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	return methodsSwitchList(result, method, args...)
}

// CallContext calls the RPC method with the specified parameters and returns an error.
//
// Parameters:
//   - ctx: represents the current execution context
//   - result: the interface{} to store the result of the RPC call
//   - args: variadic and can be used to pass additional arguments to the RPC method
//
// Returns:
//   - error: an error if any occurred during the function call
func (r *rpcMock) CallContext(ctx context.Context, result interface{}, method string, args interface{}) error {
	return methodsSwitchList(result, method, args)
}

// methodsSwitchList is a function that switches on the method name and calls the corresponding mock function.
//
// Parameters:
//   - result: The result of the RPC call
//   - method: The name of the method to call
//   - args: The arguments to pass to the method
//
// Returns:
//   - error: An error if the method is not found or if the arguments are invalid
//
//nolint:gocyclo // Inevitable due to many switch cases
func methodsSwitchList(result interface{}, method string, args ...interface{}) error {
	switch method {
	case "starknet_addDeclareTransaction":
		return mock_starknet_addDeclareTransaction(result, args...)
	case "starknet_addDeployAccountTransaction":
		return mock_starknet_addDeployAccountTransaction(result, args...)
	case "starknet_addInvokeTransaction":
		return mock_starknet_addInvokeTransaction(result, args...)
	case "starknet_blockHashAndNumber":
		return mock_starknet_blockHashAndNumber(result, args...)
	case "starknet_blockNumber":
		return mock_starknet_blockNumber(result, args...)
	case "starknet_call":
		return mock_starknet_call(result, args...)
	case "starknet_chainId":
		return mock_starknet_chainId(result, args...)
	case "starknet_estimateFee":
		return mock_starknet_estimateFee(result, args...)
	case "starknet_estimateMessageFee":
		return mock_starknet_estimateMessageFee(result, args...)
	case "starknet_getBlockTransactionCount":
		return mock_starknet_getBlockTransactionCount(result, args...)
	case "starknet_getBlockWithReceipts":
		return mock_starknet_getBlockWithReceipts(result, args...)
	case "starknet_getBlockWithTxHashes":
		return mock_starknet_getBlockWithTxHashes(result, args...)
	case "starknet_getBlockWithTxs":
		return mock_starknet_getBlockWithTxs(result, args...)
	case "starknet_getClass":
		return mock_starknet_getClass(result, args...)
	case "starknet_getClassAt":
		return mock_starknet_getClassAt(result, args...)
	case "starknet_getClassHashAt":
		return mock_starknet_getClassHashAt(result, args...)
	case "starknet_getCompiledCasm":
		return mock_starknet_getCompiledCasm(result, args...)
	case "starknet_getEvents":
		return mock_starknet_getEvents(result, args...)
	case "starknet_getMessagesStatus":
		return mock_starknet_getMessagesStatus(result, args...)
	case "starknet_getNonce":
		return mock_starknet_getNonce(result, args...)
	case "starknet_getStateUpdate":
		return mock_starknet_getStateUpdate(result, args...)
	case "starknet_getStorageAt":
		return mock_starknet_getStorageAt(result, args...)
	case "starknet_getTransactionByBlockIdAndIndex":
		return mock_starknet_getTransactionByBlockIdAndIndex(result, args...)
	case "starknet_getTransactionByHash":
		return mock_starknet_getTransactionByHash(result, args...)
	case "starknet_getTransactionReceipt":
		return mock_starknet_getTransactionReceipt(result, args...)
	case "starknet_simulateTransactions":
		return mock_starknet_simulateTransactions(result, args...)
	case "starknet_syncing":
		return mock_starknet_syncing(result, args...)
	case "starknet_traceBlockTransactions":
		return mock_starknet_traceBlockTransactions(result, args...)
	case "starknet_traceTransaction":
		return mock_starknet_traceTransaction(result, args...)
	default:
		return errNotFound
	}
}

// mock_starknet_blockNumber is a function that mocks the blockNumber functionality in the StarkNet API.
//
// Parameters:
//   - result: The result variable that will hold the block number value
//   - args: Additional arguments passed to the function
//
// Returns:
//   - error: An error if the result is not of type uint64 or if the arguments count is not zero
func mock_starknet_blockNumber(result interface{}, args ...interface{}) error {
	r, ok := result.(*json.RawMessage)
	if !ok || r == nil {
		return errWrongType
	}
	if len(args) != 0 {
		return errWrongArgs
	}

	resp, err := json.Marshal(uint64(1234))
	if err != nil {
		return err
	}

	return json.Unmarshal(resp, r)
}

// mock_starknet_blockHashAndNumber is a function that mocks the behaviour of the `blockHashAndNumber` method in the StarkNet API.
//
// Parameters:
//   - result: an interface{} that holds the result of the function.
//   - args: a variadic parameter of type interface{} that represents the arguments of the function.
//
// Returns:
//   - error: an error if there is a wrong type or wrong number of arguments.
func mock_starknet_blockHashAndNumber(result interface{}, args ...interface{}) error {
	r, ok := result.(*json.RawMessage)
	if !ok || r == nil {
		return errWrongType
	}
	if len(args) != 0 {
		return errWrongArgs
	}

	blockData := BlockHashAndNumberOutput{
		Number: 1234,
		Hash:   internalUtils.RANDOM_FELT,
	}

	resp, err := json.Marshal(blockData)
	if err != nil {
		return err
	}

	return json.Unmarshal(resp, r)
}

// mock_starknet_chainId is a function that mocks the behaviour of the `starknet_chainId` method.
//
// Parameters:
//   - result: an interface{} that holds the result of the function.
//   - args: a variadic parameter of type interface{} that represents the arguments of the function.
//
// Returns:
//   - error: an error if there is a wrong type or wrong number of arguments.
func mock_starknet_chainId(result interface{}, args ...interface{}) error {
	r, ok := result.(*json.RawMessage)
	if !ok {
		return errWrongType
	}
	if len(args) != 0 {
		return errWrongArgs
	}
	resp, err := json.Marshal("0x534e5f5345504f4c4941")
	if err != nil {
		return err
	}

	return json.Unmarshal(resp, r)
}

// mock_starknet_syncing is a function that mocks the behaviour of the starknet_syncing function.
//
// Parameters:
//   - result: an interface{} that holds the result of the function.
//   - args: a variadic parameter of type interface{} that represents the arguments of the function.
//
// Return:
//   - error: an error if there is a wrong type or wrong number of arguments
func mock_starknet_syncing(result interface{}, args ...interface{}) error {
	// Note: Since starknet_syncing returns with bool or SyncStatus, we pass in interface{}
	r, ok := result.(*interface{})
	if !ok {
		return errWrongType
	}
	if len(args) != 0 {
		return errWrongArgs
	}

	value := SyncStatus{
		StartingBlockHash: internalUtils.RANDOM_FELT,
		StartingBlockNum:  "0x4c602",
		CurrentBlockHash:  internalUtils.RANDOM_FELT,
		CurrentBlockNum:   "0x4c727",
		HighestBlockHash:  internalUtils.RANDOM_FELT,
		HighestBlockNum:   "0x4c727",
	}
	*r = value

	return nil
}

// mock_starknet_getTransactionByBlockIdAndIndex is a function that mocks the behaviour of getting
// a transaction by block ID and index in the StarkNet API.
//
// Parameters:
//   - result: The result of the API call, which will be stored in the provided interface{}. This should be a pointer to a json.RawMessage
//   - args: The arguments of the API call. This should be a variadic parameter that accepts a variable number of arguments
//
// Returns:
//   - error: An error if the API call fails, otherwise nil
func mock_starknet_getTransactionByBlockIdAndIndex(result interface{}, args ...interface{}) error {
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

	blockTxnexample, err := internalUtils.UnmarshalJSONFileToType[BlockTransaction](
		"testData/transactions/sepoliaBlockInvokeTxV3_0x265f6a59e7840a4d52cec7db37be5abd724fdfd72db9bf684f416927a88bc89.json",
		"",
	)
	if err != nil {
		return err
	}

	txBytes, err := json.Marshal(blockTxnexample)
	if err != nil {
		return err
	}

	return json.Unmarshal(txBytes, &r)
}

// mock_starknet_getBlockTransactionCount is a function that mocks the behaviour of the
// GetBlockTransactionCount method in the StarkNet API.
//
// Parameters:
//   - result: The result of the API call, which will be stored in the provided interface{}. This should be a pointer to a json.RawMessage
//   - args: The arguments of the API call. This should be a variadic parameter that accepts a variable number of arguments
//
// Returns:
//   - error: An error if the API call fails, otherwise nil
func mock_starknet_getBlockTransactionCount(result interface{}, args ...interface{}) error {
	r, ok := result.(*json.RawMessage)
	if !ok || r == nil {
		return errWrongType
	}
	if len(args) != 1 {
		return errWrongArgs
	}
	outputContent, err := json.Marshal(uint64(10))
	if err != nil {
		return err
	}
	if err := json.Unmarshal(outputContent, r); err != nil {
		return err
	}

	return nil
}

// mock_starknet_getTransactionByHash is a function that retrieves a transaction by its hash.
//
// Parameters:
//   - result: an interface{} that represents the result of the transaction retrieval
//   - args: a variadic parameter that contains the arguments used for retrieval
//
// Returns:
//   - error: an error if there is a failure in retrieving the transaction
func mock_starknet_getTransactionByHash(result interface{}, args ...interface{}) error {
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

	BlockDeclareTxnV2Example := `{
		"transaction_hash": "0xd109474cd037bad60a87ba0ccf3023d5f2d1cd45220c62091d41a614d38eda",
		"type": "DECLARE",
		"sender_address": "0x5fd4befee268bf6880f955875cbed3ade8346b1f1e149cc87b317e62b6db569",
		"compiled_class_hash": "0x7130f75fc2f1400813d1e96ea7ebee334b568a87b645a62aade0eb2fa2cf252",
		"max_fee": "0x4a0fbb2d7a43",
		"version": "0x2",
		"signature": [
		   "0x5569787df42fece1184537b0d480900a403386355b9d6a59e7c7a7e758287f0",
		   "0x2acaeea2e0817da33ed5dbeec295b0177819b5a5a50b0a669e6eecd88e42e92"
		],
		"nonce": "0x16e",
		"class_hash": "0x79b7ec8fdf40a4ff6ed47123049dfe36b5c02db93aa77832682344775ef70c6"
	}`

	if err := json.Unmarshal([]byte(BlockDeclareTxnV2Example), r); err != nil {
		return err
	}

	return nil
}

// mock_starknet_getTransactionReceipt mocks the function that retrieves the transaction receipt information
// from the StarkNet blockchain.
//
// Parameters:
//   - result: a pointer to an interface that will store the transaction receipt result
//   - args: a variadic parameter representing the arguments of the transaction receipt
//
// Returns:
//   - error: an error if there is an issue with the type of the result or the number of arguments
func mock_starknet_getTransactionReceipt(result interface{}, args ...interface{}) error {
	r, ok := result.(*json.RawMessage)
	if !ok || r == nil {
		return errWrongType
	}
	if len(args) != 1 {
		return errWrongArgs
	}

	arg0Felt := args[0].(*felt.Felt)
	l1BlockHash, err := new(felt.Felt).SetString("0x74011377f326265f5a54e27a27968355e7033ad1de11b77b225374875aff519")
	if err != nil {
		return err
	}
	testTxnHash, err := internalUtils.HexToFelt("0xf2f3d50192637e8d5e817363460c39d3a668fe12f117ecedb9749466d8352b")
	if err != nil {
		return err
	}
	if arg0Felt.Equal(testTxnHash) {
		var txnRec *TransactionReceiptWithBlockInfo
		txnRec, err = internalUtils.UnmarshalJSONFileToType[TransactionReceiptWithBlockInfo](
			"testData/receipt/sepoliaRec_0xf2f3d50192637e8d5e817363460c39d3a668fe12f117ecedb9749466d8352b.json",
			"",
		)
		if err != nil {
			return err
		}

		var txnReceipt []byte
		txnReceipt, err = json.Marshal(txnRec)
		if err != nil {
			return err
		}

		return json.Unmarshal(txnReceipt, &r)
	} else if arg0Felt.Equal(l1BlockHash) {
		var txnRec *TransactionReceiptWithBlockInfo
		txnRec, err = internalUtils.UnmarshalJSONFileToType[TransactionReceiptWithBlockInfo]("testData/receipt/mainnetRc_0x74011377f326265f5a54e27a27968355e7033ad1de11b77b225374875aff519.json", "")
		if err != nil {
			return err
		}

		var txnReceipt []byte
		txnReceipt, err = json.Marshal(txnRec)
		if err != nil {
			return err
		}

		return json.Unmarshal(txnReceipt, &r)
	}

	fromAddressFelt, err := internalUtils.HexToFelt("0xdeadbeef")
	if err != nil {
		return err
	}

	transaction := TransactionReceipt{
		Hash:           arg0Felt,
		FinalityStatus: TxnFinalityStatusAcceptedOnL1,
		Events: []Event{{
			FromAddress: fromAddressFelt,
		}},
	}
	outputContent, err := json.Marshal(transaction)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(outputContent, r); err != nil {
		return err
	}

	return nil
}

// mock_starknet_getClassAt is a function that performs a mock operation to get the class at a given index.
// The function sets the class to a DeprecatedContractClass with a specific program and marshals the class into JSON format.
// Finally, it unmarshals the JSON content into the result.
//
// Parameters:
//   - result: An interface{} that represents the result of the operation
//   - args: A variadic parameter that represents the arguments to be passed
//
// Returns:
//   - error: An error if the result is not of type *json.RawMessage or is nil or the number of arguments is not equal to 2
//
// The function always returns nil.
func mock_starknet_getClassAt(result interface{}, args ...interface{}) error {
	r, ok := result.(*json.RawMessage)
	if !ok || r == nil {
		fmt.Printf("%T\n", result)

		return errWrongType
	}
	if len(args) != 2 {
		return errWrongArgs
	}
	fakeSelector, err := internalUtils.HexToFelt("0xdeadbeef")
	if err != nil {
		return err
	}
	class := contracts.DeprecatedContractClass{
		Program: "H4sIAAAAAAAE/+Vde3PbOJL/Kj5VXW1mVqsC36Sr9g8n0c6mzonnbM",
		DeprecatedEntryPointsByType: contracts.DeprecatedEntryPointsByType{
			Constructor: []contracts.DeprecatedCairoEntryPoint{},
			External: []contracts.DeprecatedCairoEntryPoint{
				{
					Offset:   "0x0xdeadbeef",
					Selector: fakeSelector,
				},
			},
			L1Handler: []contracts.DeprecatedCairoEntryPoint{},
		},
	}
	outputContent, err := json.Marshal(class)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(outputContent, r); err != nil {
		return err
	}

	return nil
}

// mock_starknet_getClassHashAt is a function that retrieves the class hash at a specific location in the StarkNet.
//
// Parameters:
//   - result: An interface{} that represents the result of the operation
//   - args: A variadic parameter that represents the arguments to be passed
//
// Returns:
//   - error: An error if the result is not of type *json.RawMessage or is nil or the number of arguments is not equal to 2
//
// The function always returns nil.
func mock_starknet_getClassHashAt(result interface{}, args ...interface{}) error {
	r, ok := result.(*json.RawMessage)
	if !ok || r == nil {
		fmt.Printf("%T\n", result)

		return errWrongType
	}
	if len(args) != 2 {
		return errWrongArgs
	}
	classHash, err := internalUtils.HexToFelt("0xdeadbeef")
	if err != nil {
		return err
	}
	outputContent, err := json.Marshal(classHash)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(outputContent, r); err != nil {
		return err
	}

	return nil
}

// mock_starknet_getClass is a function that retrieves a class from the StarkNet API.
// It takes in a result interface{} and variadic args ...interface{}.
// The result interface{} should be a pointer to json.RawMessage.
// The method string specifies the method to be called on the StarkNet API.
// The args ...interface{} are the arguments to be passed to the method.
// The function returns an error if any of the following conditions are met:
//   - The result is not of type *json.RawMessage.
//   - The args length is not equal to 2.
//   - The first argument is not of type BlockID.
//   - The second argument is not of type *felt.Felt or does not have a hexadecimal prefix.
//
// The function assigns a DeprecatedContractClass struct to the variable class.
// The function then marshals the class to JSON and unmarshals it to the result interface{}.
// If successful, the function returns nil.
//
// Parameters:
//   - result: The result interface{} that should be a pointer to json.RawMessage
//   - args: The variadic args ...interface{} representing the arguments to be passed to the method
//
// Returns:
//   - error: An error if any of the conditions mentioned above are met
func mock_starknet_getClass(result interface{}, args ...interface{}) error {
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
	class := contracts.DeprecatedContractClass{
		Program: "H4sIAAAAAAAA",
	}
	outputContent, err := json.Marshal(class)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(outputContent, r); err != nil {
		return err
	}

	return nil
}

// mock_starknet_getEvents is a function that retrieves events from the StarkNet blockchain.
//
// Parameters:
//   - result: An interface{} that represents the result of the operation
//   - args: A variadic parameter that represents the arguments to be passed
//
// Returns:
//   - error: An error if the result is not of type *json.RawMessage or is nil or the number of arguments is not equal to 1
//
// The function always returns nil
func mock_starknet_getEvents(result interface{}, args ...interface{}) error {
	r, ok := result.(*json.RawMessage)
	if !ok {
		return errWrongType
	}
	if len(args) != 1 {
		return errWrongArgs
	}
	ei, ok := args[0].(EventsInput)
	if !ok {
		return errWrongArgs
	}
	if ei.ChunkSize == 0 {
		return errors.New("-ChuckSize error message-")
	}

	blockHash, err := internalUtils.HexToFelt("0x59dbe64bf2e2f89f5f2958cff11044dca0c64dea2e37ec6eaad9a5f838793cb")
	if err != nil {
		return err
	}
	txHash, err := internalUtils.HexToFelt("0x568147c09d5e5db8dc703ce1da21eae47e9ad9c789bc2f2889c4413a38c579d")
	if err != nil {
		return err
	}

	events := EventChunk{
		Events: []EmittedEvent{
			{
				BlockHash:       blockHash,
				BlockNumber:     1472,
				TransactionHash: txHash,
			},
		},
	}

	outputContent, err := json.Marshal(events)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(outputContent, r); err != nil {
		return err
	}

	return nil
}

// mock_starknet_call is a function that mocks a Starknet call.
//
// Parameters:
//   - result: The result of the transaction
//   - args: The arguments to be passed to the method
//
// Returns:
//   - error: An error if the transaction fails
func mock_starknet_call(result interface{}, args ...interface{}) error {
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
	outputContent, err := json.Marshal([]*felt.Felt{out})
	if err != nil {
		return err
	}
	if err := json.Unmarshal(outputContent, r); err != nil {
		return err
	}

	return nil
}

// mock_starknet_addDeclareTransaction is a mock function that adds a declare transaction to the StarkNet smart contract.
//
// Parameters:
//   - result: The result of the transaction
//   - args: The arguments to be passed to the method
//
// Return:
//   - error: An error if the transaction fails
func mock_starknet_addDeclareTransaction(result interface{}, args ...interface{}) error {
	r, ok := result.(*json.RawMessage)
	if !ok {
		return errWrongType
	}

	switch args[0].(type) {
	case *BroadcastDeclareTxnV3, BroadcastDeclareTxnV3:
		deadbeefFelt, err := internalUtils.HexToFelt("0x41d1f5206ef58a443e7d3d1ca073171ec25fa75313394318fc83a074a6631c3")
		if err != nil {
			return err
		}
		output := AddDeclareTransactionOutput{
			TransactionHash: deadbeefFelt,
			ClassHash:       deadbeefFelt,
		}
		outputContent, err := json.Marshal(output)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(outputContent, r); err != nil {
			return err
		}

		return nil
	}

	return errors.Wrap(errWrongArgs, fmt.Sprintf("args[0] should be BroadcastDeclareTxnV3, got %T\n", args[0]))
}

// mock_starknet_estimateFee simulates the estimation of a fee in the StarkNet network.
//
// Parameters:
//   - result: The result of the transaction
//   - args: The arguments to be passed to the method
//
// Returns:
//   - error: an error if any
func mock_starknet_estimateFee(result interface{}, args ...interface{}) error {
	r, ok := result.(*json.RawMessage)
	if !ok {
		return errWrongType
	}
	if len(args) != 3 {
		fmt.Printf("args: %d\n", len(args))

		return errWrongArgs
	}
	_, ok = args[0].([]BroadcastTxn)
	if !ok {
		fmt.Printf("args[0] should be BroadcastTxn, got %T\n", args[0])

		return errWrongArgs
	}
	flags, ok := args[1].([]SimulationFlag)
	if !ok {
		fmt.Printf("args[1] should be SimulationFlag, got %T\n", args[1])

		return errWrongArgs
	}
	_, ok = args[2].(BlockID)
	if !ok {
		fmt.Printf("args[2] should be *blockID, got %T\n", args[2])

		return errWrongArgs
	}

	var output FeeEstimation

	if len(flags) > 0 {
		output = FeeEstimation{
			FeeEstimationCommon: FeeEstimationCommon{
				L1GasConsumed:     new(felt.Felt).SetUint64(1234),
				L1GasPrice:        new(felt.Felt).SetUint64(1234),
				L2GasConsumed:     new(felt.Felt).SetUint64(1234),
				L2GasPrice:        new(felt.Felt).SetUint64(1234),
				L1DataGasConsumed: new(felt.Felt).SetUint64(1234),
				L1DataGasPrice:    new(felt.Felt).SetUint64(1234),
				OverallFee:        new(felt.Felt).SetUint64(1234),
			},
			Unit: FriUnit,
		}
	} else {
		output = FeeEstimation{
			FeeEstimationCommon: FeeEstimationCommon{
				L1GasConsumed:     internalUtils.RANDOM_FELT,
				L1GasPrice:        internalUtils.RANDOM_FELT,
				L2GasConsumed:     internalUtils.RANDOM_FELT,
				L2GasPrice:        internalUtils.RANDOM_FELT,
				L1DataGasConsumed: internalUtils.RANDOM_FELT,
				L1DataGasPrice:    internalUtils.RANDOM_FELT,
				OverallFee:        internalUtils.RANDOM_FELT,
			},
			Unit: FriUnit,
		}
	}

	outputContent, err := json.Marshal([]FeeEstimation{output})
	if err != nil {
		return err
	}
	if err := json.Unmarshal(outputContent, r); err != nil {
		return err
	}

	return nil
}

// mock_starknet_estimateMessageFee is a function that estimates the fee for a StarkNet message.
//
// Parameters:
//   - result: The result of the transaction
//   - args: The arguments to be passed to the method
//
// Returns:
//   - error: an error if any
func mock_starknet_estimateMessageFee(result interface{}, args ...interface{}) error {
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

	output := MessageFeeEstimation{
		FeeEstimationCommon: FeeEstimationCommon{
			L1DataGasConsumed: internalUtils.RANDOM_FELT,
			L1DataGasPrice:    internalUtils.RANDOM_FELT,
			L1GasConsumed:     internalUtils.RANDOM_FELT,
			L1GasPrice:        internalUtils.RANDOM_FELT,
			L2GasConsumed:     internalUtils.RANDOM_FELT,
			L2GasPrice:        internalUtils.RANDOM_FELT,
			OverallFee:        internalUtils.RANDOM_FELT,
		},
		Unit: WeiUnit,
	}
	outputContent, err := json.Marshal(output)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(outputContent, r); err != nil {
		return err
	}

	return nil
}

// mock_starknet_simulateTransactions is a function that simulates transactions on the StarkNet network.
// If any error occurs during the process, it is returned.
//
// Parameters:
//   - result: The result of the method
//   - args: The arguments to be passed to the method
//
// Returns:
//   - error: an error if any
func mock_starknet_simulateTransactions(result interface{}, args ...interface{}) error {
	r, ok := result.(*json.RawMessage)
	if !ok {
		return errWrongType
	}
	if len(args) != 3 {
		fmt.Printf("args: %d\n", len(args))

		return errWrongArgs
	}
	_, ok = args[0].(BlockID)
	if !ok {
		fmt.Printf("args[0] should be *blockID, got %T\n", args[0])

		return errWrongArgs
	}
	_, ok = args[1].([]BroadcastTxn)
	if !ok {
		fmt.Printf("args[1] should be BroadcastTxn, got %T\n", args[1])

		return errWrongArgs
	}
	_, ok = args[2].([]SimulationFlag)
	if !ok {
		fmt.Printf("args[2] should be SimulationFlag, got %T\n", args[2])

		return errWrongArgs
	}

	output, err := internalUtils.UnmarshalJSONFileToType[[]SimulatedTransaction]("./testData/trace/sepoliaSimulateInvokeTxResp.json", "")
	if err != nil {
		return err
	}

	outputContent, err := json.Marshal(output)
	if err != nil {
		return err
	}

	return json.Unmarshal(outputContent, r)
}

// mock_starknet_addInvokeTransaction is a mock function that simulates the behaviour of the
// starknet_addInvokeTransaction function.
// Parameters:
//   - result: The result of the transaction
//   - args: The arguments to be passed to the method
//
// Returns:
//   - error: an error if any
func mock_starknet_addInvokeTransaction(result interface{}, args ...interface{}) error {
	r, ok := result.(*json.RawMessage)
	if !ok {
		return errWrongType
	}
	if len(args) != 1 {
		return errors.Wrap(errWrongArgs, fmt.Sprint("wrong number of args ", len(args)))
	}
	switch args[0].(type) {
	case *BroadcastInvokeTxnV3, BroadcastInvokeTxnV3:
		deadbeefFelt, err := internalUtils.HexToFelt("0x49728601e0bb2f48ce506b0cbd9c0e2a9e50d95858aa41463f46386dca489fd")
		if err != nil {
			return err
		}
		output := AddInvokeTransactionResponse{
			Hash: deadbeefFelt,
		}
		outputContent, err := json.Marshal(output)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(outputContent, r); err != nil {
			return err
		}

		return nil
	default:
		return errors.Wrap(errWrongArgs, fmt.Sprintf("args[0] should be InvokeTxnV3, got %T\n", args[0]))
	}
}

func mock_starknet_addDeployAccountTransaction(result interface{}, args ...interface{}) error {
	r, ok := result.(*json.RawMessage)
	if !ok {
		return errWrongType
	}
	if len(args) != 1 {
		return errors.Wrap(errWrongArgs, fmt.Sprint("wrong number of args ", len(args)))
	}
	switch args[0].(type) {
	case *BroadcastDeployAccountTxnV3, BroadcastDeployAccountTxnV3:

		deadbeefFelt, err := internalUtils.HexToFelt("0x32b272b6d0d584305a460197aa849b5c7a9a85903b66e9d3e1afa2427ef093e")
		if err != nil {
			return err
		}
		output := AddDeployAccountTransactionResponse{
			Hash:            deadbeefFelt,
			ContractAddress: new(felt.Felt).SetUint64(0),
		}
		outputContent, err := json.Marshal(output)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(outputContent, r); err != nil {
			return err
		}

		return nil
	default:
		return errors.Wrap(errWrongArgs, fmt.Sprintf("args[0] should be DeployAccountTxnV3, got %T\n", args[0]))
	}
}

// mock_starknet_getStorageAt mocks the behaviour of the StarkNet getStorageAt function.
//
// Parameters:
//   - result: The result of the transaction
//   - args: The arguments to be passed to the method
//
// Returns:
//   - error: an error if any
func mock_starknet_getStorageAt(result interface{}, args ...interface{}) error {
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
	outputContent, err := json.Marshal(output)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(outputContent, r); err != nil {
		return err
	}

	return nil
}

// mock_starknet_getStateUpdate is a function that performs a mock operation to get the state update.
//
// Parameters:
//   - result: an interface{} that represents the result of the state update.
//   - args: a variadic parameter that can accept multiple arguments.
//
// Returns:
//   - error: an error if any
func mock_starknet_getStateUpdate(result interface{}, args ...interface{}) error {
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

	stateFeltArr, err := internalUtils.HexArrToFelt([]string{
		"0x62ab7b3ade3e7c26d0f50cb539c621b679e07440685d639904663213f906938",
		"0x491250c959067f21177f50cfdfede2bd9c8f2597f4ed071dbdba4a7ee3dabec",
		"0x19aa982a75263d4c4de4cc4c5d75c3dec32e00b95bef7bbb4d17762a0b138af",
		"0xe5cc6f2b6d34979184b88334eb64173fe4300cab46ecd3229633fcc45c83d4",
		"0x1813aac5f5e7799684c6dc33e51f44d3627fd748c800724a184ed5be09b713e",
		"0x630b4197",
	})
	if err != nil {
		return err
	}

	output := StateUpdateOutput{
		BlockHash: stateFeltArr[0],
		NewRoot:   stateFeltArr[1],
		Pre_confirmedStateUpdate: Pre_confirmedStateUpdate{
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
	outputContent, err := json.Marshal(output)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(outputContent, r); err != nil {
		return err
	}

	return nil
}

// mock_starknet_getNonce is a function that retrieves the nonce for a given method and arguments.
//
// Parameters:
//   - result: a pointer to the variable where the result will be stored
//   - args: optional arguments for the method
//
// Returns:
//   - error: an error if
//   - The result parameter is not of type *json.RawMessage
//   - The number of arguments is not equal to 2
//   - The first argument is not of type BlockID
//   - The second argument is not of type *felt.Felt
func mock_starknet_getNonce(result interface{}, args ...interface{}) error {
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
	output, err := internalUtils.HexToFelt("0xdeadbeef")
	if err != nil {
		return err
	}
	outputContent, err := json.Marshal(output)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(outputContent, r); err != nil {
		return err
	}

	return nil
}

// mock_starknet_getBlockWithTxs mocks the behaviour of the starknet_getBlockWithTxs function.
// If successful, it populates the result parameter with the json.RawMessage containing the block with the transactions.
//
// Parameters:
//   - result: the result is expected to be a pointer to json.RawMessage
//   - args: variadic parameter that can contain any number of arguments
//
// Returns:
//   - error: an error if any
func mock_starknet_getBlockWithTxs(result interface{}, args ...interface{}) error {
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

	if blockId.Tag == BlockTagPre_confirmed {
		pBlock, err := json.Marshal(
			Pre_confirmedBlock{
				Pre_confirmedBlockHeader{
					Number:           1234,
					Timestamp:        123,
					SequencerAddress: internalUtils.RANDOM_FELT,
				},
				[]BlockTransaction{},
			},
		)
		if err != nil {
			return err
		}

		return json.Unmarshal(pBlock, &r)
	} else {
		fullBlockSepolia64159, err := internalUtils.UnmarshalJSONFileToType[Block]("testData/block/sepoliaBlockTxs65083.json", "result")
		if err != nil {
			return err
		}

		blockBytes, err := json.Marshal(fullBlockSepolia64159)
		if err != nil {
			return err
		}

		return json.Unmarshal(blockBytes, &r)
	}
}

// mock_starknet_getBlockWithTxHashes mocks the behaviour of the starknet_getBlockWithTxHashes function.
// If successful, it populates the result parameter with the json.RawMessage containing the block with the specified transaction hashes.
//
// Parameters:
//   - result: the result is expected to be a pointer to json.RawMessage
//   - args: variadic parameter that can contain any number of arguments
//
// Returns:
//   - error: an error if any
func mock_starknet_getBlockWithTxHashes(result interface{}, args ...interface{}) error {
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

	txHashes, err := internalUtils.HexArrToFelt([]string{
		"0x5754961d70d6f39d0e2c71a1a4ff5df0a26b1ceda4881ca82898994379e1e73",
		"0x692381bba0e8505a8e0b92d0f046c8272de9e65f050850df678a0c10d8781d",
	})
	if err != nil {
		return err
	}

	if blockId.Tag == BlockTagPre_confirmed {
		pBlock, innerErr := json.Marshal(
			Pre_confirmedBlockTxHashes{
				Pre_confirmedBlockHeader{
					Number:           1234,
					Timestamp:        123,
					SequencerAddress: internalUtils.RANDOM_FELT,
				},
				txHashes,
			})
		if innerErr != nil {
			return innerErr
		}
		err = json.Unmarshal(pBlock, &r)
		if err != nil {
			return err
		}
	} else {
		block, innerErr := json.Marshal(
			BlockTxHashes{
				BlockHeader: BlockHeader{
					Hash:             internalUtils.RANDOM_FELT,
					ParentHash:       internalUtils.RANDOM_FELT,
					Timestamp:        124,
					SequencerAddress: internalUtils.RANDOM_FELT,
				},
				Status:       BlockStatus_AcceptedOnL1,
				Transactions: txHashes,
			})
		if innerErr != nil {
			return innerErr
		}
		err = json.Unmarshal(block, &r)
		if err != nil {
			return err
		}
	}

	return nil
}

func mock_starknet_getBlockWithReceipts(result interface{}, args ...interface{}) error {
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

	if blockId.Tag == BlockTagPre_confirmed {
		pBlock, innerErr := json.Marshal(
			Pre_confirmedBlockWithReceipts{
				Pre_confirmedBlockHeader{
					Number: 1234,
				},
				BlockBodyWithReceipts{
					Transactions: []TransactionWithReceipt{
						{
							Transaction: BlockTransaction{
								Hash: internalUtils.RANDOM_FELT,
								Transaction: InvokeTxnV1{
									Type:          "INVOKE",
									Version:       TransactionV1,
									SenderAddress: internalUtils.RANDOM_FELT,
								},
							},
							Receipt: TransactionReceipt{
								Type:            "INVOKE",
								Hash:            internalUtils.RANDOM_FELT,
								ExecutionStatus: TxnExecutionStatusSUCCEEDED,
								FinalityStatus:  TxnFinalityStatusAcceptedOnL1,
							},
						},
					},
				},
			},
		)
		if innerErr != nil {
			return innerErr
		}
		err := json.Unmarshal(pBlock, &r)
		if err != nil {
			return err
		}
	} else {
		block, innerErr := json.Marshal(
			BlockWithReceipts{
				BlockHeader{
					Hash: internalUtils.RANDOM_FELT,
				},
				"ACCEPTED_ON_L1",
				BlockBodyWithReceipts{
					Transactions: []TransactionWithReceipt{
						{
							Transaction: BlockTransaction{
								Hash: internalUtils.RANDOM_FELT,
								Transaction: InvokeTxnV1{
									Type:          "INVOKE",
									Version:       TransactionV1,
									SenderAddress: internalUtils.RANDOM_FELT,
								},
							},
							Receipt: TransactionReceipt{
								Type:            "INVOKE",
								Hash:            internalUtils.RANDOM_FELT,
								ExecutionStatus: TxnExecutionStatusSUCCEEDED,
								FinalityStatus:  TxnFinalityStatusAcceptedOnL1,
							},
						},
					},
				},
			},
		)
		if innerErr != nil {
			return innerErr
		}
		err := json.Unmarshal(block, &r)
		if err != nil {
			return err
		}
	}

	return nil
}

// mock_starknet_traceBlockTransactions is a function that traces the transactions of a block in the StarkNet network.
// The function first checks the type of the result parameter and returns an error if it is not of type *json.RawMessage.
// It then checks the length of the args parameter and returns an error if it is not equal to 1. Next, it checks the
// type of the first element of args and returns an error if it is not of type *felt.Felt. If the block hash is equal
// to "0x42a4c6a4c3dffee2cce78f04259b499437049b0084c3296da9fbbec7eda79b2", the function reads the trace from a file
// and unmarshals it into a struct. It then marshals the result and unmarshals it into the result parameter.
// If the block hash is not valid, the function returns an error of type ErrInvalidBlockHash.
//
// Parameters:
//   - result: a pointer to the variable where the result will be stored
//   - args: optional arguments for the method
//
// Returns:
//   - error: an error if any
func mock_starknet_traceBlockTransactions(result interface{}, args ...interface{}) error {
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
	if blockID.Hash != nil && blockID.Hash.String() == "0x42a4c6a4c3dffee2cce78f04259b499437049b0084c3296da9fbbec7eda79b2" {
		rawBlockTrace, err := internalUtils.UnmarshalJSONFileToType[[]Trace](
			"testData/trace/sepoliaBlockTrace_0x42a4c6a4c3dffee2cce78f04259b499437049b0084c3296da9fbbec7eda79b2.json",
			"",
		)
		if err != nil {
			return err
		}

		BlockTrace, err := json.Marshal(rawBlockTrace)
		if err != nil {
			return err
		}

		return json.Unmarshal(BlockTrace, &r)
	}

	return ErrBlockNotFound
}

// mock_starknet_traceTransaction is a Go function that traces a transaction in the StarkNet network.
// The function returns an error if any of the following conditions are met:
//   - The result is not of type *json.RawMessage.
//   - The result is nil.
//   - The number of arguments is not equal to 1.
//   - The first argument is not of type *felt.Felt.
//   - The transaction hash does not match any known hash.
//
// If the transaction hash matches "0xff66e14fc6a96f3289203690f5f876cb4b608868e8549b5f6a90a21d4d6329",
// the function reads the trace from a file and unmarshals it into the result.
//
// If the transaction hash matches "0xf00d", the function returns a custom RPCError.
//
// If the transaction hash does not match any known hash, the function returns ErrInvalidTxnHash.
//
// Parameters:
//   - result: an interface{} that represents the result of the transaction.
//   - args: a variadic parameter that can accept multiple arguments.
//
// Returns:
//   - error: an error if any
func mock_starknet_traceTransaction(result interface{}, args ...interface{}) error {
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
	case "0x6a4a9c4f1a530f7d6dd7bba9b71f090a70d1e3bbde80998fde11a08aab8b282":
		rawTrace, err := internalUtils.UnmarshalJSONFileToType[InvokeTxnTrace](
			"testData/trace/sepoliaInvokeTrace_0x6a4a9c4f1a530f7d6dd7bba9b71f090a70d1e3bbde80998fde11a08aab8b282.json",
			"",
		)
		if err != nil {
			return err
		}

		txnTrace, err := json.Marshal(rawTrace)
		if err != nil {
			return err
		}

		return json.Unmarshal(txnTrace, &r)
	case "0xf00d":
		return &RPCError{
			Code:    10,
			Message: "No trace available for transaction",
			Data:    &TraceStatusErrData{Status: TraceStatusRejected},
		}
	default:
		return ErrHashNotFound
	}
}

// mock_starknet_getCompiledCasm mocks the behaviour of getting compiled CASM for a contract class.
//
// Parameters:
//   - result: The result of the operation
//   - args: The arguments to be passed to the method
//
// Returns:
//   - error: an error if any
func mock_starknet_getCompiledCasm(result interface{}, args ...interface{}) error {
	r, ok := result.(*json.RawMessage)
	if !ok || r == nil {
		return errWrongType
	}
	if len(args) != 1 {
		return errWrongArgs
	}
	classHash, ok := args[0].(*felt.Felt)
	if !ok {
		return errWrongArgs
	}

	if classHash.String() == "0xbad" {
		return &RPCError{
			Code:    100,
			Message: "Failed to compile the contract",
			Data:    &CompilationErrData{CompilationError: "compilation error: invalid sierra class"},
		}
	}

	// Return error for specific test case
	if classHash != internalUtils.RANDOM_FELT {
		return ErrClassHashNotFound
	}

	// Read the test data from file
	resp, err := internalUtils.UnmarshalJSONFileToType[json.RawMessage]("testData/compiledCasm.json", "result")
	if err != nil {
		return err
	}

	return json.Unmarshal(*resp, r)
}

// mock_starknet_getMessagesStatus mocks the behaviour of getting the status of L1->L2 messages.
//
// Parameters:
//   - result: The result of the operation
//   - args: The arguments to be passed to the method
//
// Returns:
//   - error: an error if any
func mock_starknet_getMessagesStatus(result interface{}, args ...interface{}) error {
	r, ok := result.(*json.RawMessage)
	if !ok || r == nil {
		return errWrongType
	}
	if len(args) != 1 {
		return errWrongArgs
	}
	txHash, ok := args[0].(NumAsHex)
	if !ok {
		return errWrongArgs
	}

	if txHash == "0xdededededededededededededededededededededededededededededededede" {
		return ErrHashNotFound
	}

	// Return mock response for successful case
	response := []MessageStatus{
		{
			Hash:            internalUtils.RANDOM_FELT,
			FinalityStatus:  TxnFinalityStatusAcceptedOnL2,
			ExecutionStatus: TxnExecutionStatusSUCCEEDED,
		},
		{
			Hash:            internalUtils.RANDOM_FELT,
			FinalityStatus:  TxnFinalityStatusAcceptedOnL2,
			ExecutionStatus: TxnExecutionStatusSUCCEEDED,
		},
	}

	outputContent, err := json.Marshal(response)
	if err != nil {
		return err
	}

	return json.Unmarshal(outputContent, r)
}
