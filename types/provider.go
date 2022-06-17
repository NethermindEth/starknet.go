package types

import (
	"context"
	"math/big"
)

type Provider interface {
	AccountNonce(context.Context, string) (*big.Int, error)
	BlockByHash(context.Context, string, string) (*Block, error)
	BlockByNumber(context.Context, *big.Int, string) (*Block, error)
	Call(context.Context, FunctionCall, string) ([]string, error)
	ChainID(context.Context) (string, error)
	Invoke(context.Context, Transaction) (*AddTxResponse, error)
	TransactionByHash(context.Context, string) (*Transaction, error)
	TransactionReceipt(context.Context, string) (*TransactionReceipt, error)
	Class(context.Context, string) (*ContractClass, error)
	ClassHashAt(context.Context, string) (*Felt, error)
	ClassAt(context.Context, string) (*ContractClass, error)
}
