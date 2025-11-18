package main

import (
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/NethermindEth/juno/core/felt"
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

	// Get test accounts to mint to
	accounts, err := devNet.Accounts()
	if err != nil {
		log.Fatal("Failed to get accounts:", err)
	}

	if len(accounts) == 0 {
		log.Fatal("No accounts available")
	}

	// Use first account address
	address, err := new(felt.Felt).SetString(accounts[0].Address)
	if err != nil {
		log.Fatal("Failed to parse address:", err)
	}

	// Amount: 1 ETH = 1e18 wei
	amount := new(big.Int)
	amount.SetString("1000000000000000000", 10)

	// Mint tokens
	response, err := devNet.Mint(address, amount)
	if err != nil {
		log.Fatal("Failed to mint tokens:", err)
	}

	fmt.Println("Tokens minted successfully!")
	fmt.Printf("New Balance: %s %s\n", response.NewBalance, response.Unit)
	fmt.Printf("Transaction Hash: %s\n", response.TransactionHash)
}
