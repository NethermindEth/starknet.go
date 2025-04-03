package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/abigen/accounts/abi/bind"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"
)


func main() {
	provider, err := rpc.NewProvider("https://starknet-testnet.public.blastapi.io/rpc/v0_6")
	if err != nil {
		log.Fatalf("Failed to create provider: %v", err)
	}

	contractAddress, err := utils.HexToFelt("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
	if err != nil {
		log.Fatalf("Failed to parse contract address: %v", err)
	}

	contract, err := NewSimpleContract(contractAddress, provider)
	if err != nil {
		log.Fatalf("Failed to create contract instance: %v", err)
	}

	balance, err := contract.GetBalance(&bind.CallOpts{
		Context: context.Background(),
	})
	if err != nil {
		log.Fatalf("Failed to get balance: %v", err)
	}
	fmt.Printf("Current balance: %s\n", balance.String())

	privateKey, err := utils.HexToFelt("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}

	opts := &bind.TransactOpts{
		From:    privateKey,
		Context: context.Background(),
	}

	amount, err := felt.NewFelt(big.NewInt(100))
	if err != nil {
		log.Fatalf("Failed to create amount: %v", err)
	}

	tx, err := contract.IncreaseBalance(opts, amount)
	if err != nil {
		log.Fatalf("Failed to increase balance: %v", err)
	}
	fmt.Printf("Transaction hash: %s\n", tx.TransactionHash.String())

	receipt, err := provider.WaitForTransaction(context.Background(), tx.TransactionHash, 5, 1)
	if err != nil {
		log.Fatalf("Failed to wait for transaction: %v", err)
	}

	if receipt.Status == "ACCEPTED_ON_L2" {
		fmt.Println("Transaction confirmed successfully!")
	} else {
		fmt.Printf("Transaction failed with status: %s\n", receipt.Status)
		os.Exit(1)
	}

	newBalance, err := contract.GetBalance(&bind.CallOpts{
		Context: context.Background(),
	})
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
	return felt.NewFelt(big.NewInt(500))
}

func (_SimpleContract *SimpleContractTransactor) IncreaseBalance(opts *bind.TransactOpts, amount *felt.Felt) (*rpc.InvokeTxnResponse, error) {
	txHash, _ := utils.HexToFelt("0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890")
	return &rpc.InvokeTxnResponse{
		TransactionHash: txHash,
	}, nil
}
