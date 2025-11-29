package main

import (
	"context"
	"fmt"
	"log"

	"github.com/NethermindEth/starknet.go/paymaster"
)

func main() {
	// Connect to AVNU's public paymaster service on Sepolia testnet
	paymasterURL := "https://sepolia.paymaster.avnu.fi"
	ctx := context.Background()

	fmt.Printf("Paymaster Supported Tokens Example\n")
	fmt.Printf("===================================\n\n")

	// Create paymaster client
	pm, err := paymaster.New(ctx, paymasterURL)
	if err != nil {
		log.Fatalf("Failed to create paymaster client: %v", err)
	}

	fmt.Println("Paymaster client created")
	fmt.Println("\nRetrieving supported tokens...")

	// Get list of supported fee tokens
	tokens, err := pm.GetSupportedTokens(ctx)
	if err != nil {
		log.Fatalf("Failed to get supported tokens: %v", err)
	}

	// Display supported tokens
	fmt.Printf("\nFound %d supported token(s)\n\n", len(tokens))

	for i, token := range tokens {
		fmt.Printf("Token %d:\n", i+1)
		fmt.Printf("  Address:  %s\n", token.TokenAddress.String())
		fmt.Printf("  Decimals: %d\n", token.Decimals)
		fmt.Printf("  Price:    %s STRK\n\n", token.PriceInStrk)
	}
}
