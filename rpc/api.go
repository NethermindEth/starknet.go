package rpc

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dontpanicdao/caigo"
	"github.com/dontpanicdao/caigo/types"
)

// Call a starknet function without creating a StarkNet transaction.
func (sc *Client) Call(ctx context.Context, call types.FunctionCall, blockIDOption BlockIDOption) ([]string, error) {
	opt := &blockID{}
	err := blockIDOption(opt)
	if err != nil {
		return nil, err
	}
	call.EntryPointSelector = caigo.BigToHex(caigo.GetSelectorFromName(call.EntryPointSelector))
	if len(call.Calldata) == 0 {
		call.Calldata = make([]string, 0)
	}
	var result []string
	if opt.BlockTag != nil {
		if err := sc.do(ctx, "starknet_call", &result, call, *opt.BlockTag); err != nil {
			return nil, err
		}
		return result, nil
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

func WithBlockIDNumber(blockNumber BlockNumber) BlockIDOption {
	return BlockIDOption(func(b *blockID) error {
		b.BlockNumber = &blockNumber
		return nil
	})
}

func WithBlockIDHash(blockHash BlockHash) BlockIDOption {
	return BlockIDOption(func(b *blockID) error {
		b.BlockHash = &blockHash
		return nil
	})
}

func WithBlockIDTag(blockTag string) BlockIDOption {
	return BlockIDOption(func(b *blockID) error {
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
func (sc *Client) BlockWithTxHashes(ctx context.Context, blockIDOption BlockIDOption) (interface{}, error) {
	opt := &blockID{}
	err := blockIDOption(opt)
	if err != nil {
		return nil, err
	}
	if opt.BlockTag != nil && *opt.BlockTag == "pending" {
		var block PendingBlockWithTxHashes
		if err := sc.do(ctx, "starknet_getBlockWithTxHashes", &block, "pending"); err != nil {
			return nil, err
		}
		return &block, nil
	}
	var block BlockWithTxHashes
	if opt.BlockTag != nil && *opt.BlockTag == "latest" {
		if err := sc.do(ctx, "starknet_getBlockWithTxHashes", &block, "latest"); err != nil {
			return nil, err
		}
		return &block, nil
	}
	if err := sc.do(ctx, "starknet_getBlockWithTxHashes", &block, opt); err != nil {
		return nil, err
	}
	return &block, nil
}

// BlockTransactionCount gets the number of transactions in a block
func (sc *Client) BlockTransactionCount(ctx context.Context, blockIDOption BlockIDOption) (uint64, error) {
	opt := &blockID{}
	err := blockIDOption(opt)
	if err != nil {
		return 0, err
	}
	var result uint64
	if opt.BlockTag != nil {
		if err := sc.do(ctx, "starknet_getBlockTransactionCount", &result, *opt.BlockTag); err != nil {
			return 0, err
		}
		return result, nil
	}

	if err := sc.do(ctx, "starknet_getBlockTransactionCount", &result, opt); err != nil {
		return 0, err
	}
	return result, nil
}

// Nonce returns the Nnce of a contract
func (sc *Client) Nonce(ctx context.Context, blockIDOption BlockIDOption, contractAddress Address) (*string, error) {
	opt := &blockID{}
	err := blockIDOption(opt)
	if err != nil {
		return nil, err
	}
	var result string
	if opt.BlockTag != nil {
		if err := sc.do(ctx, "starknet_getNonce", &result, *opt.BlockTag, contractAddress); err != nil {
			return nil, err
		}
		return &result, nil
	}

	if err := sc.do(ctx, "starknet_getNonce", &result, opt, contractAddress); err != nil {
		return nil, err
	}
	return &result, nil
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
func (sc *Client) BlockWithTxs(ctx context.Context, blockIDOption BlockIDOption) (interface{}, error) {
	opt := &blockID{}
	err := blockIDOption(opt)
	if err != nil {
		return nil, err
	}
	if opt.BlockTag != nil && *opt.BlockTag == "pending" {
		var block PendingBlockWithTxs
		if err := sc.do(ctx, "starknet_getBlockWithTxs", &block, "pending"); err != nil {
			return nil, err
		}
		txns, err := guessTxsWithType(block.Transactions)
		if err != nil {
			return nil, err
		}
		block.Transactions = txns
		return &block, nil
	}
	var block BlockWithTxs
	if opt.BlockTag != nil && *opt.BlockTag == "latest" {
		if err := sc.do(ctx, "starknet_getBlockWithTxs", &block, "latest"); err != nil {
			return nil, err
		}
		txns, err := guessTxsWithType(block.Transactions)
		if err != nil {
			return nil, err
		}
		block.Transactions = txns
		return &block, nil
	}
	if err := sc.do(ctx, "starknet_getBlockWithTxs", &block, opt); err != nil {
		return nil, err
	}
	txns, err := guessTxsWithType(block.Transactions)
	if err != nil {
		return nil, err
	}
	block.Transactions = txns
	return &block, nil
}

// Class gets the contract class definition associated with the given hash.
func (sc *Client) Class(ctx context.Context, blockIDOption BlockIDOption, classHash string) (*ContractClass, error) {
	opt := &blockID{}
	err := blockIDOption(opt)
	if err != nil {
		return nil, err
	}
	var rawClass ContractClass
	if opt.BlockTag != nil {
		if err := sc.do(ctx, "starknet_getClass", &rawClass, *opt.BlockTag, classHash); err != nil {
			return nil, err
		}
		return &rawClass, nil
	}
	if err := sc.do(ctx, "starknet_getClass", &rawClass, opt, classHash); err != nil {
		return nil, err
	}
	return &rawClass, nil
}

// ClassAt get the contract class definition at the given address.
func (sc *Client) ClassAt(ctx context.Context, blockIDOption BlockIDOption, contractAddress Address) (*ContractClass, error) {
	opt := &blockID{}
	err := blockIDOption(opt)
	if err != nil {
		return nil, err
	}
	var rawClass ContractClass
	if opt.BlockTag != nil {
		if err := sc.do(ctx, "starknet_getClassAt", &rawClass, *opt.BlockTag, contractAddress); err != nil {
			return nil, err
		}
		return &rawClass, nil
	}
	if err := sc.do(ctx, "starknet_getClassAt", &rawClass, opt, contractAddress); err != nil {
		return nil, err
	}
	return &rawClass, nil
}

// ClassHashAt gets the contract class hash for the contract deployed at the given address.
func (sc *Client) ClassHashAt(ctx context.Context, blockIDOption BlockIDOption, contractAddress Address) (*string, error) {
	opt := &blockID{}
	err := blockIDOption(opt)
	if err != nil {
		return nil, err
	}
	if opt.BlockTag != nil && *opt.BlockTag != "pending" && *opt.BlockTag != "latest" {
		return nil, errInvalidBlockTag
	}
	var result string
	if opt.BlockTag != nil {
		if err := sc.do(ctx, "starknet_getClassHashAt", &result, *opt.BlockTag, contractAddress); err != nil {
			return nil, err
		}
		return &result, nil
	}
	if err := sc.do(ctx, "starknet_getStateUpdate", &result, opt, contractAddress); err != nil {
		return nil, err
	}
	return &result, nil
}

// StorageAt gets the value of the storage at the given address and key.
func (sc *Client) StorageAt(ctx context.Context, contractAddress Address, key string, blockIDOption BlockIDOption) (string, error) {
	opt := &blockID{}
	err := blockIDOption(opt)
	if err != nil {
		return "", err
	}
	var value string
	hashKey := fmt.Sprintf("0x%s", caigo.GetSelectorFromName(key).Text(16))
	if opt.BlockTag != nil {
		if err := sc.do(ctx, "starknet_getStorageAt", &value, string(contractAddress), hashKey, *opt.BlockTag); err != nil {
			return "", err
		}
		return value, nil
	}
	if err := sc.do(ctx, "starknet_getStorageAt", &value, string(contractAddress), hashKey, opt); err != nil {
		return "", err
	}
	return value, nil
}

// StateUpdate gets the information about the result of executing the requested block.
func (sc *Client) StateUpdate(ctx context.Context, blockIDOption BlockIDOption) (*StateUpdateOutput, error) {
	opt := &blockID{}
	err := blockIDOption(opt)
	if err != nil {
		return nil, err
	}
	if opt.BlockTag != nil && *opt.BlockTag != "latest" {
		return nil, errInvalidBlockTag
	}
	var state StateUpdateOutput
	if opt.BlockTag != nil {
		if err := sc.do(ctx, "starknet_getStateUpdate", &state, "latest"); err != nil {
			return nil, err
		}
		return &state, nil
	}
	if err := sc.do(ctx, "starknet_getStateUpdate", &state, opt); err != nil {
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
func (sc *Client) TransactionByBlockIdAndIndex(ctx context.Context, blockIDOption BlockIDOption, index uint64) (*Txn, error) {
	opt := &blockID{}
	err := blockIDOption(opt)
	if err != nil {
		return nil, err
	}
	if opt.BlockTag != nil && *opt.BlockTag != "latest" {
		return nil, errInvalidBlockTag
	}
	var tx interface{}
	if opt.BlockTag != nil {
		if err := sc.do(ctx, "starknet_getTransactionByBlockIdAndIndex", &tx, *opt.BlockTag, index); err != nil {
			return nil, err
		}
		txWithType, err := guessTxWithType(tx)
		if err != nil {
			return nil, err
		}
		txTxn := Txn(txWithType)
		return &txTxn, nil
	}
	if err := sc.do(ctx, "starknet_getTransactionByBlockIdAndIndex", &tx, opt, index); err != nil {
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
	var receipt types.TransactionReceipt
	err := sc.do(ctx, "starknet_getTransactionReceipt", &receipt)
	if err != nil {
		return nil, err
	} else if receipt.TransactionHash == "" {
		return nil, errNotFound
	}

	return &receipt, nil
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
func (sc *Client) EstimateFee(ctx context.Context, request BroadcastedTxn, blockIDOption BlockIDOption) (*FeeEstimate, error) {
	opt := &blockID{}
	err := blockIDOption(opt)
	if err != nil {
		return nil, err
	}
	var raw FeeEstimate
	if opt.BlockTag != nil {
		if err := sc.do(ctx, "starknet_estimateFee", &raw, request, *opt.BlockTag); err != nil {
			return nil, err
		}
		return &raw, nil
	}
	if err := sc.do(ctx, "starknet_estimateFee", &raw, request, opt); err != nil {
		return nil, err
	}
	return &raw, nil
}

func (sc *Client) Invoke(context.Context, types.FunctionInvoke) (*types.AddTxResponse, error) {
	panic("not implemented")
}
