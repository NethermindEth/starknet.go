package main

import (
	"context"
	"fmt"
	"log"

	"github.com/NethermindEth/starknet.go/paymaster"
)

func main() {
	// Using AVNU's public SNIP-29 compliant paymaster on Sepolia testnet
	paymasterURL := "https://sepolia.paymaster.avnu.fi"
	ctx := context.Background()

	// Create a new paymaster client
	pm, err := paymaster.New(ctx, paymasterURL)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Paymaster client created successfully")

	// Check if the paymaster service is available
	available, err := pm.IsAvailable(ctx)
	if err != nil {
		log.Fatal(err)
	}

	if available {
		fmt.Println("Paymaster service is ready")
	}
}
