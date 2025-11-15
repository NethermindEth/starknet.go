package main

import (
	"context"
	"fmt"
	"log"

	"github.com/NethermindEth/starknet.go/paymaster"
)

func main() {
	// Paymaster service URL (example)
	paymasterURL := "https://paymaster.example.com/rpc"

	// Create a new paymaster client
	ctx := context.Background()
	pm, err := paymaster.New(ctx, paymasterURL)
	if err != nil {
		log.Fatal("Failed to create paymaster client:", err)
	}

	fmt.Println("New Paymaster Client:")
	fmt.Printf("  Service URL: %s\n", paymasterURL)
	fmt.Printf("  Client created: %v\n", pm != nil)
	
	// Note: Actual usage requires a running paymaster service
	fmt.Println("\nNote: This requires a SNIP-29 compliant paymaster service")
}
