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

	// Create account
	acc, err := account.NewAccount(provider, accountAddress, publicKey, ks, account.CairoV2)
	if err != nil {
		log.Fatal(err)
	}

	// Create a message to sign
	message, _ := new(felt.Felt).SetString("0x4d7920746573742073747277742073747277742073747277742073747277742053696d706c65206d657373616765")

	fmt.Println("Signing message:")
	fmt.Printf("Message: %s\n", message)

	// Sign the message
	signature, err := acc.Sign(ctx, message)
	if err != nil {
		fmt.Printf("Error signing: %v\n", err)
		return
	}

	fmt.Println("\nSignature components:")
	for i, sig := range signature {
		fmt.Printf("Signature[%d]: %s\n", i, sig)
	}
	fmt.Printf("\nTotal signature components: %d\n", len(signature))
}