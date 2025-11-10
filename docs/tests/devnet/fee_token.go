package main

import (
	"fmt"
	"log"
	"os"

	"github.com/NethermindEth/starknet.go/devnet"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")

	// Get DevNet URL from environment variable
	devnetURL := os.Getenv("DEVNET_URL")
	if devnetURL == "" {
		devnetURL = "http://localhost:5050"
	}

	devNet := devnet.NewDevNet(devnetURL)

	// Get fee token information
	feeToken, err := devNet.FeeToken()
	if err != nil {
		log.Fatal("Failed to get fee token:", err)
	}

	fmt.Println("Fee Token Information:")
	fmt.Printf("  Symbol:  %s\n", feeToken.Symbol)
	fmt.Printf("  Address: %s\n", feeToken.Address.String())
}
