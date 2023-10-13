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

func (r *rpcMock) Close() {
	r.closed = true
}

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
	value := "0x534e5f474f45524c49"
	*r = value
	return nil
}

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

	output := FeeEstimate{
		GasConsumed: NumAsHex("0x01a4"),
		GasPrice:    NumAsHex("0x45"),
		OverallFee:  NumAsHex("0x7134"),
	}
	outputContent, _ := json.Marshal(output)
	json.Unmarshal(outputContent, r)
	return nil
}

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
		GasConsumed: NumAsHex("0x1"),
		GasPrice:    NumAsHex("0x2"),
		OverallFee:  NumAsHex("0x3"),
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
	if len(args) != 1 {
		return errors.Wrap(errWrongArgs, fmt.Sprint("wrong number of args ", len(args)))
	}
	invokeTx, ok := args[0].(InvokeTxnV1)
	if !ok {
		return errors.Wrap(errWrongArgs, fmt.Sprintf("args[0] should be InvokeTxnV1, got %T\n", args[0]))
	}
	if invokeTx.SenderAddress != nil {

		if invokeTx.SenderAddress.Equal(new(felt.Felt).SetUint64(123)) {
			unexpErr := ErrUnexpectedError
			unexpErr.data = "Something crazy happened"
			return unexpErr
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
	output := "0x0"
	outputContent, _ := json.Marshal(output)
	json.Unmarshal(outputContent, &r)
	return nil
}

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

func mock_starknet_traceBlockTransactions(result interface{}, method string, args ...interface{}) error {
	r, ok := result.(*json.RawMessage)
	if !ok || r == nil {
		return errWrongType
	}
	if len(args) != 1 {
		return errWrongArgs
	}
	blockHash, ok := args[0].(*felt.Felt)
	if !ok {
		return errors.Wrap(errWrongArgs, fmt.Sprintf("args[0] should be felt, got %T\n", args[0]))
	}
	if blockHash.String() == "0x3ddc3a8aaac071ecdc5d8d0cfbb1dc4fc6a88272bc6c67523c9baaee52a5ea2" {

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
