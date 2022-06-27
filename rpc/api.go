package rpc

import (
	"context"
	"encoding/json"
	"math/big"

	"github.com/dontpanicdao/caigo"
	"github.com/dontpanicdao/caigo/types"
)

type Events struct {
	Events []Event `json:"events"`
}

type Event struct {
	*types.Event
	FromAddress     string `json:"from_address"`
	BlockHash       string `json:"block_hash"`
	BlockNumber     int    `json:"block_number"`
	TransactionHash string `json:"transaction_hash"`
}

type EventParams struct {
	FromBlock  uint64 `json:"fromBlock"`
	ToBlock    uint64 `json:"toBlock"`
	PageSize   uint64 `json:"page_size"`
	PageNumber uint64 `json:"page_number"`
}

func (sc *Client) Call(ctx context.Context, call types.FunctionCall, hash string) ([]string, error) {
	call.EntryPointSelector = caigo.BigToHex(caigo.GetSelectorFromName(call.EntryPointSelector))
	if len(call.Calldata) == 0 {
		call.Calldata = make([]string, 0)
	}

	var result []string
	if err := sc.do(ctx, "starknet_call", &result, call, hash); err != nil {
		return nil, err
	}

	return result, nil
}

// BlockNumber returns the current block managed by the API.
func (sc *Client) BlockNumber(ctx context.Context) (*big.Int, error) {
	var blockNumber big.Int
	if err := sc.c.CallContext(ctx, &blockNumber, "starknet_blockNumber"); err != nil {
		return nil, err
	}

	return &blockNumber, nil
}

func (sc *Client) BlockByHash(ctx context.Context, hash string, scope string) (*types.Block, error) {
	var block types.Block
	if err := sc.do(ctx, "starknet_getBlockByHash", &block, hash, scope); err != nil {
		return nil, err
	}

	return &block, nil
}

func (sc *Client) BlockByNumber(ctx context.Context, number *big.Int, scope string) (*types.Block, error) {
	var block types.Block
	if err := sc.do(ctx, "starknet_getBlockByNumber", &block, toBlockNumArg(number), scope); err != nil {
		return nil, err
	}

	return &block, nil
}

func (sc *Client) CodeAt(ctx context.Context, address string) (*types.Code, error) {
	var contractRaw struct {
		Bytecode []string `json:"bytecode"`
		AbiRaw   string   `json:"abi"`
		Abi      types.ABI
	}
	if err := sc.do(ctx, "starknet_getCode", &contractRaw, address); err != nil {
		return nil, err
	}

	contract := types.Code{
		Bytecode: contractRaw.Bytecode,
	}
	if err := json.Unmarshal([]byte(contractRaw.AbiRaw), &contract.Abi); err != nil {
		return nil, err
	}

	return &contract, nil
}

func (sc *Client) Class(ctx context.Context, hash string) (*types.ContractClass, error) {
	var rawClass types.ContractClass
	if err := sc.do(ctx, "starknet_getClass", &rawClass, hash); err != nil {
		return nil, err
	}

	return &rawClass, nil
}

func (sc *Client) ClassAt(ctx context.Context, address string) (*types.ContractClass, error) {
	var rawClass types.ContractClass
	if err := sc.do(ctx, "starknet_getClassAt", &rawClass, address); err != nil {
		return nil, err
	}

	return &rawClass, nil
}

func (sc *Client) ClassHashAt(ctx context.Context, address string) (string, error) {
	result := new(string)
	if err := sc.do(ctx, "starknet_getClassHashAt", &result, address); err != nil {
		return "", err
	}

	return *result, nil
}

func (sc *Client) TransactionByHash(ctx context.Context, hash string) (*types.Transaction, error) {
	var tx types.Transaction
	if err := sc.do(ctx, "starknet_getTransactionByHash", &tx, hash); err != nil {
		return nil, err
	} else if tx.TransactionHash == "" {
		return nil, ErrNotFound
	}

	return &tx, nil
}

func (sc *Client) TransactionReceipt(ctx context.Context, hash string) (*types.TransactionReceipt, error) {
	var receipt types.TransactionReceipt
	err := sc.do(ctx, "starknet_getTransactionReceipt", &receipt, hash)
	if err != nil {
		return nil, err
	} else if receipt.TransactionHash == "" {
		return nil, ErrNotFound
	}

	return &receipt, nil
}

func (sc *Client) Events(ctx context.Context, evParams EventParams) (*Events, error) {
	var result Events
	if err := sc.do(ctx, "starknet_getEvents", &result, evParams); err != nil {
		return nil, err
	}

	return &result, nil
}

func (sc *Client) EstimateFee(context.Context, types.Transaction) (*types.FeeEstimate, error) {
	panic("not implemented")
}

func (sc *Client) AccountNonce(context.Context, string) (*big.Int, error) {
	panic("not implemented")
}

func toBlockNumArg(number *big.Int) interface{} {
	var numOrTag interface{}

	if number == nil {
		numOrTag = "latest"
	} else if number.Cmp(big.NewInt(-1)) == 0 {
		numOrTag = "pending"
	} else {
		numOrTag = number.Uint64()
	}

	return numOrTag
}
