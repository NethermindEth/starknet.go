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

	// Get pre-funded accounts
	accounts, err := devNet.Accounts()
	if err != nil {
		log.Fatal("Failed to get accounts:", err)
	}

	fmt.Printf("Found %d pre-funded accounts\n\n", len(accounts))

	// Display first account
	if len(accounts) > 0 {
		fmt.Println("First Account:")
		fmt.Printf("  Address:     %s\n", accounts[0].Address)
		fmt.Printf("  Public Key:  %s\n", accounts[0].PublicKey)
		fmt.Printf("  Private Key: %s\n", accounts[0].PrivateKey)
	}
}
