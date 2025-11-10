package main

import (
	"fmt"
	"math/big"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/account"
)

func main() {
	// Create a new empty MemKeystore
	ks := account.NewMemKeystore()

	fmt.Println("Created new empty MemKeystore")

	// Generate a test key pair (simplified for demo)
	privKeyFelt, _ := new(felt.Felt).SetRandom()
	bytes := privKeyFelt.Bytes()
	privKey := new(big.Int).SetBytes(bytes[:])

	// For this demo, we'll use a simple public key representation
	pubKey := "0x1234567890abcdef"

	fmt.Printf("\nTest key pair:")
	fmt.Printf("Public Key:  %s\n", pubKey)
	fmt.Printf("Private Key: %s\n", privKey)

	// Try to get a key that doesn't exist
	fmt.Println("\nAttempting to get non-existent key:")
	_, err := ks.Get("0x9999")
	if err != nil {
		fmt.Printf("Error (expected): %v\n", err)
	}

	// Add the key to the keystore
	fmt.Println("\nAdding key to keystore:")
	ks.Put(pubKey, privKey)
	fmt.Println("Key added successfully")

	// Retrieve the key
	fmt.Println("\nRetrieving key from keystore:")
	retrieved, err := ks.Get(pubKey)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Retrieved key: %s\n", retrieved)
		fmt.Printf("Key retrieved successfully: %v\n", retrieved != nil)
	}

	// Add another key
	fmt.Println("\nAdding another key:")
	pubKey2 := "0xfedcba0987654321"
	privKey2 := new(big.Int).SetUint64(999999)
	ks.Put(pubKey2, privKey2)
	fmt.Printf("Added key with public key: %s\n", pubKey2)

	// Retrieve the second key
	retrieved2, err := ks.Get(pubKey2)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Retrieved second key: %s\n", retrieved2)
	}
}