package main

import (
	"fmt"
	"math/big"

	"github.com/NethermindEth/juno/core/felt"
)

func main() {
	// Create account address (ETH token address for example)
	accountAddress, _ := new(felt.Felt).SetString("0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7")

	// Public key for the account
	publicKey := "0x03603a2692a2ae60abb343e832ee53b55d6b25f02a3ef1565ec691edc7a209b2"

	// Private key for keystore
	privateKey := new(big.Int).SetUint64(123456789)

	fmt.Println("NewAccount parameters:")
	fmt.Printf("Account address: %s\n", accountAddress)
	fmt.Printf("Public key: %s\n", publicKey)
	fmt.Printf("Private key: %d\n", privateKey)

	fmt.Println("\nCairo versions:")
	fmt.Println("CairoV0 = 0")
	fmt.Println("CairoV2 = 2")
}