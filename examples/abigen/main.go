package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/abigen/accounts/abi/bind"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"
)

type MockContractBackend struct {
	*rpc.Provider
}

func (m *MockContractBackend) Call(call rpc.FunctionCall, blockID rpc.BlockID) ([]*felt.Felt, error) {
	return m.Provider.Call(context.Background(), call, blockID)
}

func (m *MockContractBackend) Invoke(opts *bind.TransactOpts, call rpc.FunctionCall) (*bind.InvokeTxnResponse, error) {
	return &bind.InvokeTxnResponse{
		TransactionHash: utils.Uint64ToFelt(1), // Dummy transaction hash
	}, nil
}

func (m *MockContractBackend) FilterEvents(ctx context.Context, filter bind.EventFilter) ([]bind.Event, error) {
	return []bind.Event{}, nil
}

func (m *MockContractBackend) DeployContract(opts *bind.TransactOpts, bytecode []byte, constructorArgs []*felt.Felt) (*felt.Felt, *bind.AddTxnResponse, error) {
	contractAddress := utils.Uint64ToFelt(1000) // Dummy contract address
	
	resp := &bind.AddTxnResponse{
		TransactionHash: utils.Uint64ToFelt(2000), // Dummy transaction hash
		ContractAddress: contractAddress,
	}
	
	return contractAddress, resp, nil
}

func main() {
	provider, err := rpc.NewProvider("https://starknet-testnet.public.blastapi.io/rpc/v0_6")
	if err != nil {
		log.Fatalf("Failed to create provider: %v", err)
	}

	backend := &MockContractBackend{Provider: provider}

	contractAddress, err := utils.HexToFelt("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
	if err != nil {
		log.Fatalf("Failed to parse contract address: %v", err)
	}

	contract, err := NewSimpleContract(contractAddress, backend)
	if err != nil {
		log.Fatalf("Failed to create contract instance: %v", err)
	}

	balance, err := contract.GetBalance(&bind.CallOpts{})
	if err != nil {
		log.Fatalf("Failed to get balance: %v", err)
	}
	fmt.Printf("Current balance: %s\n", balance.String())

	privateKey, err := utils.HexToFelt("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}

	opts := &bind.TransactOpts{
		From: privateKey,
	}

	amount := new(felt.Felt).SetUint64(100)

	tx, err := contract.IncreaseBalance(opts, amount)
	if err != nil {
		log.Fatalf("Failed to increase balance: %v", err)
	}
	fmt.Printf("Transaction hash: %s\n", tx.TransactionHash.String())

	receipt, err := provider.TransactionReceipt(context.Background(), tx.TransactionHash)
	if err != nil {
		log.Fatalf("Failed to get transaction receipt: %v", err)
	}

	if receipt.FinalityStatus == "ACCEPTED_ON_L2" {
		fmt.Println("Transaction confirmed successfully!")
	} else {
		fmt.Printf("Transaction failed with status: %s\n", receipt.FinalityStatus)
		os.Exit(1)
	}

	newBalance, err := contract.GetBalance(&bind.CallOpts{})
	if err != nil {
		log.Fatalf("Failed to get updated balance: %v", err)
	}
	fmt.Printf("New balance: %s\n", newBalance.String())
}

type SimpleContract struct {
	SimpleContractCaller
	SimpleContractTransactor
	SimpleContractFilterer
}

type SimpleContractCaller struct {
	contract *bind.BoundContract
}

type SimpleContractTransactor struct {
	contract *bind.BoundContract
}

type SimpleContractFilterer struct {
	contract *bind.BoundContract
}

func NewSimpleContract(address *felt.Felt, backend bind.ContractBackend) (*SimpleContract, error) {
	boundContract := &bind.BoundContract{}
	return &SimpleContract{
		SimpleContractCaller:     SimpleContractCaller{contract: boundContract},
		SimpleContractTransactor: SimpleContractTransactor{contract: boundContract},
		SimpleContractFilterer:   SimpleContractFilterer{contract: boundContract},
	}, nil
}

func (_SimpleContract *SimpleContractCaller) GetBalance(opts *bind.CallOpts) (*felt.Felt, error) {
	return new(felt.Felt).SetUint64(500), nil
}

func (_SimpleContract *SimpleContractTransactor) IncreaseBalance(opts *bind.TransactOpts, amount *felt.Felt) (*rpc.AddInvokeTransactionResponse, error) {
	txHash, _ := utils.HexToFelt("0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890")
	return &rpc.AddInvokeTransactionResponse{
		TransactionHash: txHash,
	}, nil
}
