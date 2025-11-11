package main

import (
	"context"
	"fmt"
	"math/big"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/account"
)

func main() {
	// Note: This example shows the structure without connecting to a real provider
	// In practice, you would use: provider, err := rpc.NewProvider(ctx, rpcURL)

	// Create account address (ETH token address for example)
	accountAddress, _ := new(felt.Felt).SetString("0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7")

	// Create keystore with test key
	publicKey := "0x03603a2692a2ae60abb343e832ee53b55d6b25f02a3ef1565ec691edc7a209b2"
	privateKey := new(big.Int).SetUint64(123456789)
	account.SetNewMemKeystore(publicKey, privateKey)

	fmt.Println("Account parameters prepared:")
	fmt.Printf("Account address: %s\n", accountAddress)
	fmt.Printf("Public key: %s\n", publicKey)
	fmt.Printf("Private key stored in keystore\n")

	// Show Cairo versions
	fmt.Println("\nAvailable Cairo versions:")
	fmt.Printf("Cairo v0: %d\n", account.CairoV0)
	fmt.Printf("Cairo v2: %d\n", account.CairoV2)

	// The actual account creation would be:
	// acc, err := account.NewAccount(provider, accountAddress, publicKey, ks, account.CairoV2)

	fmt.Println("\nNewAccount creates an Account instance with:")
	fmt.Println("- Provider for blockchain interaction")
	fmt.Println("- Account contract address")
	fmt.Println("- Public key for verification")
	fmt.Println("- Keystore for signing operations")
	fmt.Println("- Cairo version specification")
}