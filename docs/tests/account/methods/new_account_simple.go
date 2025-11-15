package main

import (
	"fmt"
	"math/big"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/account"
)

func main() {
	// Note: NewAccount requires a valid provider, but for documentation
	// we'll show the structure and parameters needed

	fmt.Println("NewAccount creates an Account instance with the following parameters:")
	fmt.Println()

	// Show the required parameters
	fmt.Println("Required parameters:")
	fmt.Println("1. provider - RPC provider for blockchain interaction")
	fmt.Println("2. accountAddress - The account contract address")
	fmt.Println("3. publicKey - Public key as hex string")
	fmt.Println("4. keystore - Keystore containing the private key")
	fmt.Println("5. cairoVersion - Cairo version (0, 1, or 2)")
	fmt.Println()

	// Show example values
	fmt.Println("Example values:")
	accountAddress, _ := new(felt.Felt).SetString("0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7")
	fmt.Printf("Account address: %s\n", accountAddress)

	publicKey := "0x03603a2692a2ae60abb343e832ee53b55d6b25f02a3ef1565ec691edc7a209b2"
	fmt.Printf("Public key: %s\n", publicKey)

	privateKey := new(big.Int).SetUint64(123456789)
	account.SetNewMemKeystore(publicKey, privateKey)
	fmt.Printf("Keystore created with private key: %s\n", privateKey)

	fmt.Println("\nCairo versions:")
	fmt.Printf("Cairo v0: %d\n", account.CairoV0)
	fmt.Printf("Cairo v2: %d\n", account.CairoV2)

	fmt.Println("\nAccount struct fields after creation:")
	fmt.Println("- CairoVersion: The Cairo version of the account")
	fmt.Println("- ChainId: The chain ID from the provider")
	fmt.Println("- provider: RPC provider instance")
	fmt.Println("- accountAddress: Account contract address")
	fmt.Println("- publicKey: Public key for signing")
	fmt.Println("- keystore: Keystore for signing operations")
}