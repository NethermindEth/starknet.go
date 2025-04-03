package bind

import (
	"context"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"
)

type ContractCaller interface {
	Call(call rpc.FunctionCall, blockID rpc.BlockID) ([]*felt.Felt, error)
}

type InvokeTxnResponse struct {
	TransactionHash *felt.Felt
}

type AddTxnResponse struct {
	TransactionHash *felt.Felt
	ContractAddress *felt.Felt
}

type EventFilter struct {
	FromBlock *felt.Felt
	ToBlock   *felt.Felt
	Address   *felt.Felt
	Keys      [][]*felt.Felt
}

type Event struct {
	FromAddress *felt.Felt
	Keys        [][]*felt.Felt
	Data        []*felt.Felt
	BlockHash   *felt.Felt
	BlockNumber uint64
	TxHash      *felt.Felt
}

type ContractTransact interface {
	Invoke(opts *TransactOpts, call rpc.FunctionCall) (*InvokeTxnResponse, error)
}

type ContractFilterer interface {
	FilterEvents(ctx context.Context, filter EventFilter) ([]Event, error)
}

type ContractBackend interface {
	ContractCaller
	ContractTransact
	ContractFilterer

	DeployContract(opts *TransactOpts, bytecode []byte, constructorArgs []*felt.Felt) (*felt.Felt, *AddTxnResponse, error)
}

type DeployBackend interface {
	DeployContract(opts *TransactOpts, bytecode []byte, constructorArgs []*felt.Felt) (*felt.Felt, *AddTxnResponse, error)
}

type RPCBackend struct {
	Client *rpc.Provider
}

func (b *RPCBackend) Call(call rpc.FunctionCall, blockID rpc.BlockID) ([]*felt.Felt, error) {
	return b.Client.Call(context.Background(), call, blockID)
}

func (b *RPCBackend) Invoke(opts *TransactOpts, call rpc.FunctionCall) (*InvokeTxnResponse, error) {
	return &InvokeTxnResponse{
		TransactionHash: utils.Uint64ToFelt(1), // Dummy transaction hash
	}, nil
}

func (b *RPCBackend) FilterEvents(ctx context.Context, filter EventFilter) ([]Event, error) {
	return []Event{}, nil
}

func (b *RPCBackend) DeployContract(opts *TransactOpts, bytecode []byte, constructorArgs []*felt.Felt) (*felt.Felt, *AddTxnResponse, error) {
	contractAddress := utils.Uint64ToFelt(1000) // Dummy contract address
	
	resp := &AddTxnResponse{
		TransactionHash: utils.Uint64ToFelt(2000), // Dummy transaction hash
		ContractAddress: contractAddress,
	}
	
	return contractAddress, resp, nil
}
