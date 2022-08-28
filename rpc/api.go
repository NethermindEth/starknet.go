package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dontpanicdao/caigo"
)

var errInvalidBlockID = errors.New("invalid blockid")

// Call a starknet function without creating a StarkNet transaction.
func (sc *Client) Call(ctx context.Context, call FunctionCall, block BlockID) ([]string, error) {
	if !block.isValid() {
		return nil, errInvalidBlockID
	}
	call.EntryPointSelector = caigo.BigToHex(caigo.GetSelectorFromName(call.EntryPointSelector))
	if len(call.CallData) == 0 {
		call.CallData = make([]string, 0)
	}
	var result []string
	if tag, ok := block.tag(); ok {
		if err := sc.do(ctx, "starknet_call", &result, call, tag); err != nil {
			return nil, err
		}
		return result, nil
	}
	opt, err := block.getWithoutTag()
	if err != nil {
		return nil, err
	}
	if err := sc.do(ctx, "starknet_call", &result, call, opt); err != nil {
		return nil, err
	}
	return result, nil
}

// BlockNumber gets the most recent accepted block number.
func (sc *Client) BlockNumber(ctx context.Context) (uint64, error) {
	var blockNumber uint64
	if err := sc.c.CallContext(ctx, &blockNumber, "starknet_blockNumber"); err != nil {
		return 0, err
	}
	return blockNumber, nil
}

// BlockHashAndNumber gets block information given the block number or its hash.
func (sc *Client) BlockHashAndNumber(ctx context.Context) (*BlockHashAndNumberOutput, error) {
	var block BlockHashAndNumberOutput
	if err := sc.do(ctx, "starknet_blockHashAndNumber", &block); err != nil {
		return nil, err
	}
	return &block, nil
}

func WithBlockNumber(blockNumber BlockNumber) BlockID {
	return BlockID(func(b *blockID) error {
		b.BlockNumber = &blockNumber
		return nil
	})
}

func WithBlockHash(blockHash BlockHash) BlockID {
	return BlockID(func(b *blockID) error {
		b.BlockHash = &blockHash
		return nil
	})
}

func WithBlockTag(blockTag string) BlockID {
	return BlockID(func(b *blockID) error {
		if blockTag != "latest" && blockTag != "pending" {
			return errInvalidBlockTag
		}
		b.BlockTag = &blockTag
		return nil
	})
}

// PendingTransactions returns the list of pending transactions.
func (sc *Client) PendingTransactions(ctx context.Context) ([]Txn, error) {
	var pendingTransactions []Txn
	if err := sc.do(ctx, "starknet_pendingTransactions", &pendingTransactions); err != nil {
		return nil, err
	}
	pendingTransactionWithTypes, err := guessTxsWithType(pendingTransactions)
	if err != nil {
		return nil, err
	}
	return pendingTransactionWithTypes, nil
}

// BlockWithTxHashes gets block information given the block id.
func (sc *Client) BlockWithTxHashes(ctx context.Context, block BlockID) (interface{}, error) {
	if !block.isValid() {
		return nil, errInvalidBlockID
	}
	if block.isPending() {
		var result PendingBlockWithTxHashes
		if err := sc.do(ctx, "starknet_getBlockWithTxHashes", &result, "pending"); err != nil {
			return nil, err
		}
		return &result, nil
	}
	var result BlockWithTxHashes
	if block.isLatest() {
		if err := sc.do(ctx, "starknet_getBlockWithTxHashes", &result, "latest"); err != nil {
			return nil, err
		}
		return &result, nil
	}
	opt, err := block.getWithoutTag()
	if err != nil {
		return nil, err
	}
	if err := sc.do(ctx, "starknet_getBlockWithTxHashes", &result, *opt); err != nil {
		return nil, err
	}
	return &result, nil
}

// BlockTransactionCount gets the number of transactions in a block
func (sc *Client) BlockTransactionCount(ctx context.Context, block BlockID) (uint64, error) {
	if !block.isValid() {
		return 0, errInvalidBlockID
	}
	var result uint64
	if tag, ok := block.tag(); ok {
		if err := sc.do(ctx, "starknet_getBlockTransactionCount", &result, tag); err != nil {
			return 0, err
		}
		return result, nil
	}
	opt, err := block.getWithoutTag()
	if err != nil {
		return 0, err
	}
	if err := sc.do(ctx, "starknet_getBlockTransactionCount", &result, *opt); err != nil {
		return 0, err
	}
	return result, nil
}

// Nonce returns the Nnce of a contract
func (sc *Client) Nonce(ctx context.Context, block BlockID, contractAddress Address) (*string, error) {
	nonce := ""
	if !block.isValid() {
		return &nonce, errInvalidBlockID
	}
	if tag, ok := block.tag(); ok {

		if err := sc.do(ctx, "starknet_getNonce", &nonce, tag, contractAddress); err != nil {
			return nil, err
		}
		return &nonce, nil
	}
	opt, err := block.getWithoutTag()
	if err != nil {
		return &nonce, err
	}
	if err := sc.do(ctx, "starknet_getNonce", &nonce, opt, contractAddress); err != nil {
		return nil, err
	}
	return &nonce, nil
}

func (i BroadcastedInvokeTxnV0) Version() uint64 {
	return 0
}

func (i BroadcastedInvokeTxnV1) Version() uint64 {
	return 1
}

func (s *StructABIEntry) IsType() string {
	return string(s.Type)
}

func (e *EventABIEntry) IsType() string {
	return string(e.Type)
}

func (f *FunctionABIEntry) IsType() string {
	return string(f.Type)
}

func guessTxWithType(i interface{}) (interface{}, error) {
	switch local := i.(type) {
	case map[string]interface{}:
		typeValue, ok := local["type"]
		if !ok {
			return nil, errBadTxType
		}
		value, ok := typeValue.(string)
		if !ok {
			return nil, errBadTxType
		}
		switch value {
		case "DECLARE":
			data, err := json.Marshal(i)
			if err != nil {
				return nil, err
			}
			tx := DeclareTxn{}
			err = json.Unmarshal(data, &tx)
			return tx, err
		case "DEPLOY":
			data, err := json.Marshal(i)
			if err != nil {
				return nil, err
			}
			tx := DeployTxn{}
			err = json.Unmarshal(data, &tx)
			return tx, err
		case "L1_HANDLER":
			data, err := json.Marshal(i)
			if err != nil {
				return nil, err
			}
			tx := L1HandlerTxn{}
			err = json.Unmarshal(data, &tx)
			return tx, err
		case "INVOKE":
			data, err := json.Marshal(i)
			if err != nil {
				return nil, err
			}
			tx := InvokeTxnDuck{}
			err = json.Unmarshal(data, &tx)
			if err != nil {
				return nil, err
			}
			if tx.AccountAddress != "" {
				txv1 := InvokeTxnV1{}
				err = json.Unmarshal(data, &txv1)
				return txv1, err
			}
			if tx.ContractAddress != "" && tx.EntryPointSelector != "" {
				txv0 := InvokeTxnV0{}
				err = json.Unmarshal(data, &txv0)
				return txv0, err
			}
			return nil, errBadTxType
		}
		return nil, errBadTxType
	}
	return nil, errBadTxType
}

func guessTxnReceiptWithType(i interface{}) (TxnReceipt, error) {
	switch local := i.(type) {
	case map[string]interface{}:
		txnType := "INVOKE"
		typeValue, ok := local["type"]
		if ok {
			txnType, ok = typeValue.(string)
			if !ok {
				return nil, errBadTxType
			}
		}
		switch txnType {
		case "DECLARE":
			data, err := json.Marshal(i)
			if err != nil {
				return nil, err
			}
			fmt.Printf("%s\n", string(data))
			receipt := DeclareTxnReceipt{}
			err = json.Unmarshal(data, &receipt)
			return receipt, err
		case "DEPLOY":
			data, err := json.Marshal(i)
			if err != nil {
				return nil, err
			}
			receipt := DeployTxnReceipt{}
			err = json.Unmarshal(data, &receipt)
			return receipt, err
		case "L1_HANDLER":
			data, err := json.Marshal(i)
			if err != nil {
				return nil, err
			}
			receipt := L1HandlerTxnReceipt{}
			err = json.Unmarshal(data, &receipt)
			return receipt, err
		case "INVOKE":
			data, err := json.Marshal(i)
			if err != nil {
				return nil, err
			}
			receipt := InvokeTxnReceipt{}
			err = json.Unmarshal(data, &receipt)
			return receipt, err
		}
		return nil, errBadTxType
	}
	return nil, errBadTxType
}

func guessTxsWithType(txs []Txn) ([]Txn, error) {
	for k, v := range txs {
		tv, err := guessTxWithType(v)
		if err != nil {
			return nil, errBadTxType
		}
		txs[k] = tv
	}
	return txs, nil
}

// BlockWithTxs get block information with full transactions given the block id.
func (sc *Client) BlockWithTxs(ctx context.Context, block BlockID) (interface{}, error) {
	if !block.isValid() {
		return nil, errInvalidBlockID
	}
	if block.isPending() {
		var result PendingBlockWithTxs
		if err := sc.do(ctx, "starknet_getBlockWithTxs", &result, "pending"); err != nil {
			return nil, err
		}
		txns, err := guessTxsWithType(result.Transactions)
		if err != nil {
			return nil, err
		}
		result.Transactions = txns
		return &result, nil
	}

	var result BlockWithTxs
	if block.isLatest() {
		if err := sc.do(ctx, "starknet_getBlockWithTxs", &result, "latest"); err != nil {
			return nil, err
		}
		txns, err := guessTxsWithType(result.Transactions)
		if err != nil {
			return nil, err
		}
		result.Transactions = txns
		return &result, nil
	}
	opt, err := block.getWithoutTag()
	if err != nil {
		return nil, err
	}
	if err := sc.do(ctx, "starknet_getBlockWithTxs", &result, *opt); err != nil {
		return nil, err
	}
	txns, err := guessTxsWithType(result.Transactions)
	if err != nil {
		return nil, err
	}
	result.Transactions = txns
	return &result, nil
}

// Class gets the contract class definition associated with the given hash.
func (sc *Client) Class(ctx context.Context, block BlockID, classHash string) (*ContractClass, error) {
	if !block.isValid() {
		return nil, errInvalidBlockID
	}
	var rawClass ContractClass
	if tag, ok := block.tag(); ok {
		if err := sc.do(ctx, "starknet_getClass", &rawClass, tag, classHash); err != nil {
			return nil, err
		}
		return &rawClass, nil
	}
	opt, err := block.getWithoutTag()
	if err != nil {
		return nil, err
	}
	if err := sc.do(ctx, "starknet_getClass", &rawClass, *opt, classHash); err != nil {
		return nil, err
	}
	return &rawClass, nil
}

// ClassAt get the contract class definition at the given address.
func (sc *Client) ClassAt(ctx context.Context, block BlockID, contractAddress Address) (*ContractClass, error) {
	if !block.isValid() {
		return nil, errInvalidBlockID
	}
	var rawClass ContractClass
	if tag, ok := block.tag(); ok {
		if err := sc.do(ctx, "starknet_getClassAt", &rawClass, tag, contractAddress); err != nil {
			return nil, err
		}
		return &rawClass, nil
	}
	opt, err := block.getWithoutTag()
	if err != nil {
		return nil, err
	}
	if err := sc.do(ctx, "starknet_getClassAt", &rawClass, *opt, contractAddress); err != nil {
		return nil, err
	}
	return &rawClass, nil
}

// ClassHashAt gets the contract class hash for the contract deployed at the given address.
func (sc *Client) ClassHashAt(ctx context.Context, block BlockID, contractAddress Address) (*string, error) {
	if !block.isValid() {
		return nil, errInvalidBlockID
	}
	var result string
	if tag, ok := block.tag(); ok {
		if err := sc.do(ctx, "starknet_getClassHashAt", &result, tag, contractAddress); err != nil {
			return nil, err
		}
		return &result, nil
	}
	opt, err := block.getWithoutTag()
	if err != nil {
		return nil, err
	}
	if err := sc.do(ctx, "starknet_getStateUpdate", &result, *opt, contractAddress); err != nil {
		return nil, err
	}
	return &result, nil
}

// StorageAt gets the value of the storage at the given address and key.
func (sc *Client) StorageAt(ctx context.Context, contractAddress Address, key string, block BlockID) (string, error) {
	if !block.isValid() {
		return "", errInvalidBlockID
	}
	var value string
	hashKey := fmt.Sprintf("0x%s", caigo.GetSelectorFromName(key).Text(16))
	if tag, ok := block.tag(); ok {
		if err := sc.do(ctx, "starknet_getStorageAt", &value, string(contractAddress), hashKey, tag); err != nil {
			return "", err
		}
		return value, nil
	}
	opt, err := block.getWithoutTag()
	if err != nil {
		return "", err
	}
	if err := sc.do(ctx, "starknet_getStorageAt", &value, string(contractAddress), hashKey, *opt); err != nil {
		return "", err
	}
	return value, nil
}

// StateUpdate gets the information about the result of executing the requested block.
func (sc *Client) StateUpdate(ctx context.Context, block BlockID) (*StateUpdateOutput, error) {
	if !block.isValid() {
		return nil, errInvalidBlockID
	}
	var state StateUpdateOutput
	if block.isLatest() {
		if err := sc.do(ctx, "starknet_getStateUpdate", &state, "latest"); err != nil {
			return nil, err
		}
		return &state, nil
	}
	opt, err := block.getWithoutTag()
	if err != nil {
		return nil, err
	}
	if err := sc.do(ctx, "starknet_getStateUpdate", &state, *opt); err != nil {
		return nil, err
	}
	return &state, nil
}

// TransactionByHash gets the details and status of a submitted transaction.
func (sc *Client) TransactionByHash(ctx context.Context, hash TxnHash) (*Txn, error) {
	var tx interface{}
	if err := sc.do(ctx, "starknet_getTransactionByHash", &tx, hash); err != nil {
		return nil, err
	}
	txWithType, err := guessTxWithType(tx)
	if err != nil {
		return nil, err
	}
	txTxn := Txn(txWithType)
	return &txTxn, nil
}

// TransactionByBlockIdAndIndex Get the details of the transaction given by the identified block and index in that block. If no transaction is found, null is returned.
func (sc *Client) TransactionByBlockIdAndIndex(ctx context.Context, block BlockID, index uint64) (*Txn, error) {
	if !block.isValid() {
		return nil, errInvalidBlockID
	}
	var tx interface{}
	if tag, ok := block.tag(); ok {
		if err := sc.do(ctx, "starknet_getTransactionByBlockIdAndIndex", &tx, tag, index); err != nil {
			return nil, err
		}
		txWithType, err := guessTxWithType(tx)
		if err != nil {
			return nil, err
		}
		txTxn := Txn(txWithType)
		return &txTxn, nil
	}
	opt, err := block.getWithoutTag()
	if err != nil {
		return nil, err
	}
	if err := sc.do(ctx, "starknet_getTransactionByBlockIdAndIndex", &tx, *opt, index); err != nil {
		return nil, err
	}
	txWithType, err := guessTxWithType(tx)
	if err != nil {
		return nil, err
	}
	txTxn := Txn(txWithType)
	return &txTxn, nil
}

// TransactionReceipt gets the transaction receipt by the transaction hash.
func (sc *Client) TransactionReceipt(ctx context.Context, transactionHash TxnHash) (TxnReceipt, error) {
	var receipt interface{}
	err := sc.do(ctx, "starknet_getTransactionReceipt", &receipt, transactionHash)
	if err != nil {
		return nil, err
	}
	return guessTxnReceiptWithType(receipt)
}

// Events returns all events matching the given filter
func (sc *Client) Events(ctx context.Context, filter EventFilterParams) (*EventsOutput, error) {
	var result EventsOutput
	if err := sc.do(ctx, "starknet_getEvents", &result, filter); err != nil {
		return nil, err
	}

	return &result, nil
}

// EstimateFee estimates the fee for a given StarkNet transaction.
func (sc *Client) EstimateFee(ctx context.Context, request BroadcastedTxn, block BlockID) (*FeeEstimate, error) {
	if !block.isValid() {
		return nil, errInvalidBlockID
	}
	var raw FeeEstimate
	if tag, ok := block.tag(); ok {
		if err := sc.do(ctx, "starknet_estimateFee", &raw, request, tag); err != nil {
			return nil, err
		}
		fmt.Printf("%+v\n", raw)
		return &raw, nil
	}
	opt, err := block.getWithoutTag()
	if err != nil {
		return nil, err
	}
	if err := sc.do(ctx, "starknet_estimateFee", &raw, request, *opt); err != nil {
		return nil, err
	}
	fmt.Printf("%+v\n", raw)
	return &raw, nil
}
