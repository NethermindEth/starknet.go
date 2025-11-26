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
	classHash, _ := new(felt.Felt).SetString("0x01234567890abcdef1234567890abcdef1234567890abcdef1234567890abcd")
	salt := new(felt.Felt).SetUint64(12345)
	maxFee := new(felt.Felt).SetUint64(2500000000000000)
	nonce := new(felt.Felt).SetUint64(0)

	publicKey, _ := new(felt.Felt).SetString("0x01234567890abcdef1234567890abcdef1234567890abcdef1234567890abcd")
	constructorCalldata := []*felt.Felt{publicKey}

	txn := &rpc.DeployAccountTxnV1{
		Type:                rpc.TransactionTypeDeployAccount,
		Version:             rpc.TransactionV1,
		ClassHash:           classHash,
		ContractAddressSalt: salt,
		ConstructorCalldata: constructorCalldata,
		MaxFee:              maxFee,
		Nonce:               nonce,
		Signature:           []*felt.Felt{},
	}

	contractAddress, _ := new(felt.Felt).SetString("0x04718f5a0fc34cc1af16a1cdee98ffb20c31f5cd61d6ab07201858f4287c938d")

	txHash, err := hash.TransactionHashDeployAccountV1(txn, contractAddress, chainID)
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
