package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	rpcURL := os.Getenv("STARKNET_RPC_URL")
	if rpcURL == "" {
		log.Fatal("STARKNET_RPC_URL not set in .env")
	}

	ctx := context.Background()
	provider, err := rpc.NewProvider(ctx, rpcURL)
	if err != nil {
		log.Fatal(err)
	}

	accountAddress, _ := new(felt.Felt).SetString("0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7")

	publicKey := "0x03603a2692a2ae60abb343e832ee53b55d6b25f02a3ef1565ec691edc7a209b2"
	privateKey := new(big.Int).SetUint64(123456789)
	ks := account.SetNewMemKeystore(publicKey, privateKey)

	acc, err := account.NewAccount(provider, accountAddress, publicKey, ks, account.CairoV2)
	if err != nil {
		log.Fatal(err)
	}

	// Create a deploy account transaction V1
	classHash, _ := new(felt.Felt).SetString("0x04c6d6cf894f8bc96bb9c525e6853e5483177841f7388f74a46cfda6f028c755")
	contractAddressSalt, _ := new(felt.Felt).SetString("0x04a7a67901e7f64e7d4f46fa17a0c57aefb0e91b3ec31e83feb758ae56b6e29e")
	pubKey, _ := new(felt.Felt).SetString("0x03603a2692a2ae60abb343e832ee53b55d6b25f02a3ef1565ec691edc7a209b2")
	
	deployAccountTx := rpc.DeployAccountTxnV1{
		Type:                rpc.TransactionTypeDeployAccount,
		MaxFee:              new(felt.Felt).SetUint64(1000000000000),
		Version:             rpc.TransactionV1,
		Signature:           []*felt.Felt{},
		Nonce:               new(felt.Felt).SetUint64(0),
		ClassHash:           classHash,
		ContractAddressSalt: contractAddressSalt,
		ConstructorCalldata: []*felt.Felt{pubKey},
	}

	// Precompute the contract address (required for deploy account hash)
	precomputedAddress, _ := new(felt.Felt).SetString("0x04a7a67901e7f64e7d4f46fa17a0c57aefb0e91b3ec31e83feb758ae56b6e29e")

	fmt.Println("Calculating transaction hash for DeployAccountTxnV1:")
	txHash, err := acc.TransactionHashDeployAccount(&deployAccountTx, precomputedAddress)
	if err != nil {
		fmt.Printf("Error calculating hash: %v\n", err)
		return
	}

	fmt.Printf("Transaction hash: %s\n", txHash)
}
