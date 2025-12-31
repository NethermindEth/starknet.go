package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/hash"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	chainID, _ := new(felt.Felt).SetString("0x534e5f5345504f4c4941")
	senderAddress, _ := new(felt.Felt).SetString("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
	maxFee := new(felt.Felt).SetUint64(2000000000000000)
	nonce := new(felt.Felt).SetUint64(5)

	ethContract, _ := new(felt.Felt).SetString("0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7")
	transferSelector, _ := new(felt.Felt).SetString("0x83afd3f4caedc6eebf44246fe54e38c95e3179a5ec9ea81740eca5b482d12e")

	calldata := []*felt.Felt{
		new(felt.Felt).SetUint64(1),
		ethContract,
		transferSelector,
		new(felt.Felt).SetUint64(2),
		new(felt.Felt).SetUint64(1),
		new(felt.Felt).SetUint64(2),
	}

	txn := &rpc.InvokeTxnV1{
		Type:          rpc.TransactionTypeInvoke,
		Version:       rpc.TransactionV1,
		SenderAddress: senderAddress,
		Nonce:         nonce,
		MaxFee:        maxFee,
		Calldata:      calldata,
		Signature:     []*felt.Felt{},
	}

	txHash, err := hash.TransactionHashInvokeV1(txn, chainID)
	if err != nil {
		log.Fatalf("Failed to calculate transaction hash: %v", err)
	}

	fmt.Printf("Transaction Hash: %s\n", txHash.String())

	if rpcURL := os.Getenv("STARKNET_RPC_URL"); rpcURL != "" {
		verifyTransaction(txHash, rpcURL)
	}
}

func verifyTransaction(txHash *felt.Felt, rpcURL string) {
	client, err := rpc.NewProvider(context.Background(), rpcURL)
	if err != nil {
		log.Printf("Warning: Could not connect to RPC: %v", err)
		return
	}

	ctx := context.Background()
	tx, err := client.TransactionByHash(ctx, txHash)

	if err == nil {
		fmt.Printf("\nVerification: FOUND on-chain\n")
		fmt.Printf("Type: %T\n", tx)
	} else {
		fmt.Printf("\nVerification: NOT FOUND\n")
	}
}
