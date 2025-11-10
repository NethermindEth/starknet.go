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
	
	fmt.Println("=== NEW CAIRO V2 ACCOUNT DETAILS ===")
	fmt.Println()
	fmt.Printf("Private Key: %s\n", privateKey)
	fmt.Printf("Public Key:  %s\n", publicKey)
	fmt.Println()

	// Use OpenZeppelin Cairo 1.0 account class hash on Sepolia
	// This is the OZ Cairo v1 account (latest version)
	// Reference: https://docs.openzeppelin.com/contracts-cairo/0.10.0/accounts
	classHash, _ := new(felt.Felt).SetString("0x00e2eb8f5672af4e6a4e8a8f1b44989685e668489b0a25437733756c5a34a1d6")
	
	// Salt for address derivation
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
	fmt.Println("2. Go to: https://starknet-faucet.vercel.app/")
	fmt.Println("3. Request ~0.002 ETH (need more for Cairo v2 deployment)")
	fmt.Println()
	fmt.Println("=== SAVE THESE VALUES ===")
	fmt.Println()
	fmt.Println("Replace in your .env file:")
	fmt.Printf("TEST_ACCOUNT_ADDRESS=%s\n", precomputedAddress)
	fmt.Printf("TEST_ACCOUNT_PRIVATE_KEY=%s\n", privateKey)
	fmt.Printf("TEST_ACCOUNT_PUBLIC_KEY=%s\n", publicKey)
	fmt.Printf("TEST_ACCOUNT_CLASS_HASH=%s\n", classHash)
	fmt.Println()
	fmt.Println("This will be a Cairo v2 account (latest version).")
}
