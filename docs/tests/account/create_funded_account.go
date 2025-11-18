package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	rpcURL := os.Getenv("STARKNET_RPC_URL")
	if rpcURL == "" {
		log.Fatal("STARKNET_RPC_URL not set in .env")
	}

	ctx := context.Background()
	_, err = rpc.NewProvider(ctx, rpcURL)
	if err != nil {
		log.Fatal(err)
	}

	// Generate random keys
	_, publicKey, privateKey := account.GetRandomKeys()
	
	fmt.Println("=== NEW TEST ACCOUNT DETAILS ===")
	fmt.Println()
	fmt.Printf("Private Key: %s\n", privateKey)
	fmt.Printf("Public Key:  %s\n", publicKey)
	fmt.Println()

	// Use OpenZeppelin account class hash on Sepolia
	// This is the standard OZ account v0.8.1 on Sepolia
	classHash, _ := new(felt.Felt).SetString("0x061dac032f228abef9c6626f995015233097ae253a7f72d68552db02f2971b8f")
	
	// Salt for address derivation (using 0 for simplicity)
	salt := new(felt.Felt).SetUint64(0)
	
	// Constructor calldata is just the public key for OZ accounts
	constructorCalldata := []*felt.Felt{publicKey}

	// Precompute the account address
	precomputedAddress := account.PrecomputeAccountAddress(
		salt,
		classHash,
		constructorCalldata,
	)

	fmt.Printf("Precomputed Account Address: %s\n", precomputedAddress)
	fmt.Println()
	fmt.Println("=== FUNDING INSTRUCTIONS ===")
	fmt.Println()
	fmt.Println("1. Copy the account address above")
	fmt.Println("2. Go to the Starknet Sepolia Faucet:")
	fmt.Println("   https://starknet-faucet.vercel.app/")
	fmt.Println()
	fmt.Println("3. Paste the address and request funds (you'll get 0.001 ETH)")
	fmt.Println()
	fmt.Println("4. Alternatively, use Alchemy faucet:")
	fmt.Println("   https://www.alchemy.com/faucets/starknet-sepolia")
	fmt.Println()
	fmt.Println("5. Wait for the transaction to complete (usually 30-60 seconds)")
	fmt.Println()
	fmt.Println("6. Once funded, we'll deploy the account")
	fmt.Println()
	fmt.Println("=== SAVE THESE VALUES ===")
	fmt.Println()
	fmt.Println("Add these to your .env file:")
	fmt.Printf("TEST_ACCOUNT_ADDRESS=%s\n", precomputedAddress)
	fmt.Printf("TEST_ACCOUNT_PRIVATE_KEY=%s\n", privateKey)
	fmt.Printf("TEST_ACCOUNT_PUBLIC_KEY=%s\n", publicKey)
	fmt.Println()
	fmt.Println("Amount needed: ~0.001 ETH (minimum for account deployment)")
}
