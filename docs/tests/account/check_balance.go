package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/NethermindEth/juno/core/felt"
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

	testAccountAddress := os.Getenv("TEST_ACCOUNT_ADDRESS")
	if testAccountAddress == "" {
		log.Fatal("TEST_ACCOUNT_ADDRESS not set in .env")
	}

	ctx := context.Background()
	provider, err := rpc.NewProvider(ctx, rpcURL)
	if err != nil {
		log.Fatal(err)
	}

	// Parse the account address
	accountAddress, err := new(felt.Felt).SetString(testAccountAddress)
	if err != nil {
		log.Fatalf("Invalid account address: %v", err)
	}

	// ETH token contract address on Starknet Sepolia
	ethContractAddress, _ := new(felt.Felt).SetString("0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7")

	// balanceOf selector
	balanceOfSelector, _ := new(felt.Felt).SetString("0x2e4263afad30923c891518314c3c95dbe830a16874e8abc5777a9a20b54c76e")

	// Call balanceOf
	callData := rpc.FunctionCall{
		ContractAddress:    ethContractAddress,
		EntryPointSelector: balanceOfSelector,
		Calldata:           []*felt.Felt{accountAddress},
	}

	result, err := provider.Call(ctx, callData, rpc.BlockID{Tag: "latest"})
	if err != nil {
		log.Fatalf("Error calling balanceOf: %v", err)
	}

	fmt.Println("=== ACCOUNT BALANCE CHECK ===")
	fmt.Println()
	fmt.Printf("Account Address: %s\n", testAccountAddress)
	fmt.Println()

	if len(result) >= 2 {
		// Balance is returned as [low, high] uint256
		low := result[0].BigInt(new(big.Int))
		high := result[1].BigInt(new(big.Int))
		
		// Combine low and high to get full uint256
		balance := new(big.Int).Lsh(high, 128)
		balance.Add(balance, low)

		// Convert to ETH (divide by 10^18)
		divisor := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
		ethBalance := new(big.Float).Quo(new(big.Float).SetInt(balance), new(big.Float).SetInt(divisor))

		fmt.Printf("Raw Balance (wei): %s\n", balance.String())
		fmt.Printf("Balance (ETH):     %s ETH\n", ethBalance.Text('f', 6))
		fmt.Println()

		if balance.Cmp(big.NewInt(0)) > 0 {
			fmt.Println("✅ Account is FUNDED and ready to use!")
		} else {
			fmt.Println("❌ Account has NO FUNDS. Please fund it first.")
		}
	} else {
		fmt.Println("❌ Could not retrieve balance")
	}
	fmt.Println()
}
