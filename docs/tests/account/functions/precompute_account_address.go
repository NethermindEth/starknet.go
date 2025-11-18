package main

import (
	"fmt"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/account"
)

func main() {
	// Create a salt for the account address
	salt, _ := new(felt.Felt).SetString("0x12345678")

	// Use a standard account class hash (OpenZeppelin account contract)
	classHash, _ := new(felt.Felt).SetString("0x061dac032f228abef9c6626f995015233097ae253a7f72d68552db02f2971b8f")

	// Constructor calldata - typically the public key for an account
	publicKey, _ := new(felt.Felt).SetString("0x01234567890abcdef01234567890abcdef01234567890abcdef01234567890ab")
	constructorCalldata := []*felt.Felt{publicKey}

	// Precompute the account address
	address := account.PrecomputeAccountAddress(salt, classHash, constructorCalldata)
	fmt.Printf("Precomputed account address: %s\n", address)

	// Test with empty constructor calldata
	fmt.Println("\nWith empty constructor calldata:")
	addressEmpty := account.PrecomputeAccountAddress(salt, classHash, []*felt.Felt{})
	fmt.Printf("Precomputed account address: %s\n", addressEmpty)

	// Test with different salt
	salt2, _ := new(felt.Felt).SetString("0x87654321")
	address2 := account.PrecomputeAccountAddress(salt2, classHash, constructorCalldata)
	fmt.Printf("\nWith different salt (0x87654321):\nPrecomputed account address: %s\n", address2)
}