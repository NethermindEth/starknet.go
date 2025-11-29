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

	fmt.Printf("Paymaster Health Check Example\n")

	// Create paymaster client
	pm, err := paymaster.New(ctx, paymasterURL)
	if err != nil {
		log.Fatalf("Failed to create paymaster client: %v", err)
	}

	fmt.Println("Paymaster client created")
	fmt.Println("\nPerforming health check...")

	// Check if paymaster service is available and operational
	available, err := pm.IsAvailable(ctx)
	if err != nil {
		log.Fatalf("Health check failed: %v", err)
	}

	// Display results
	fmt.Println()
	if available {
		fmt.Println("Paymaster service is AVAILABLE")
		fmt.Println("\nService Status: OPERATIONAL")
		fmt.Println("Ready to process transactions")
	} else {
		fmt.Println("Paymaster service is UNAVAILABLE")
		fmt.Println("\nService Status: DOWN")
		fmt.Println("The service is temporarily unavailable")
	}
}
