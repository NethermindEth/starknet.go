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
	// Load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get RPC URL from environment
	rpcURL := os.Getenv("STARKNET_RPC_URL")
	if rpcURL == "" {
		log.Fatal("STARKNET_RPC_URL not set in .env")
	}

	// Create RPC provider
	ctx := context.Background()
	provider, err := rpc.NewProvider(ctx, rpcURL)
	if err != nil {
		log.Fatal(err)
	}

	// Create account address
	accountAddress, _ := new(felt.Felt).SetString("0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7")

	// Create keystore with test key
	publicKey := "0x03603a2692a2ae60abb343e832ee53b55d6b25f02a3ef1565ec691edc7a209b2"
	privateKey := new(big.Int).SetUint64(123456789)
	ks := account.SetNewMemKeystore(publicKey, privateKey)

	// Create Cairo v2 account
	fmt.Println("Testing Cairo v2 account:")
	acc2, err := account.NewAccount(provider, accountAddress, publicKey, ks, account.CairoV2)
	if err != nil {
		log.Fatal(err)
	}

	// Create function calls
	contractAddress, _ := new(felt.Felt).SetString("0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7")
	entryPointSelector, _ := new(felt.Felt).SetString("0x83afd3f4caedc6eebf44246fe54e38c95e3179a5ec9ea81740eca5b482d12e")
	recipient, _ := new(felt.Felt).SetString("0x1234567890abcdef")

	functionCalls := []rpc.FunctionCall{
		{
			ContractAddress:    contractAddress,
			EntryPointSelector: entryPointSelector,
			Calldata: []*felt.Felt{
				recipient,
				new(felt.Felt).SetUint64(100),
				new(felt.Felt).SetUint64(0),
			},
		},
	}

	// Format calldata for Cairo v2
	calldata2, err := acc2.FmtCalldata(functionCalls)
	if err != nil {
		fmt.Printf("Error formatting calldata: %v\n", err)
		return
	}

	fmt.Println("\nCairo v2 formatted calldata:")
	for i, data := range calldata2 {
		fmt.Printf("  [%d]: %s\n", i, data)
	}
	fmt.Printf("Total elements: %d\n", len(calldata2))

	// Create Cairo v0 account
	fmt.Println("\nTesting Cairo v0 account:")
	acc0, err := account.NewAccount(provider, accountAddress, publicKey, ks, account.CairoV0)
	if err != nil {
		log.Fatal(err)
	}

	// Format calldata for Cairo v0
	calldata0, err := acc0.FmtCalldata(functionCalls)
	if err != nil {
		fmt.Printf("Error formatting calldata: %v\n", err)
		return
	}

	fmt.Println("\nCairo v0 formatted calldata:")
	for i, data := range calldata0 {
		fmt.Printf("  [%d]: %s\n", i, data)
	}
	fmt.Printf("Total elements: %d\n", len(calldata0))
}