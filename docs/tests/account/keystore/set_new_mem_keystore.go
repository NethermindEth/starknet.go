package main

import (
	"fmt"
	"math/big"

	"github.com/NethermindEth/starknet.go/account"
)

func main() {
	// Create a keystore with an initial key pair
	pubKey1 := "0xabc123def456789"
	privKey1 := new(big.Int).SetUint64(111111)

	fmt.Println("Creating MemKeystore with initial key pair:")
	fmt.Printf("Public Key:  %s\n", pubKey1)
	fmt.Printf("Private Key: %s\n", privKey1)

	ks := account.SetNewMemKeystore(pubKey1, privKey1)
	fmt.Println("\nKeystore created with initial key")

	// Retrieve the initial key
	fmt.Println("\nRetrieving the initial key:")
	retrieved1, err := ks.Get(pubKey1)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Retrieved key: %s\n", retrieved1)
		fmt.Printf("Keys match: %v\n", retrieved1.String() == privKey1.String())
	}

	// Add more keys to the keystore
	fmt.Println("\nAdding additional keys to the keystore:")

	pubKey2 := "0xdef456789abc123"
	privKey2 := new(big.Int).SetUint64(222222)
	ks.Put(pubKey2, privKey2)
	fmt.Printf("Added key 2: Public=%s, Private=%s\n", pubKey2, privKey2)

	pubKey3 := "0x789abcdef123456"
	privKey3 := new(big.Int).SetUint64(333333)
	ks.Put(pubKey3, privKey3)
	fmt.Printf("Added key 3: Public=%s, Private=%s\n", pubKey3, privKey3)

	// Retrieve all keys to verify they're stored
	fmt.Println("\nRetrieving all keys:")

	retrieved2, err := ks.Get(pubKey2)
	if err != nil {
		fmt.Printf("Error retrieving key 2: %v\n", err)
	} else {
		fmt.Printf("Retrieved key 2: %s\n", retrieved2)
	}

	retrieved3, err := ks.Get(pubKey3)
	if err != nil {
		fmt.Printf("Error retrieving key 3: %v\n", err)
	} else {
		fmt.Printf("Retrieved key 3: %s\n", retrieved3)
	}

	// Try to retrieve a non-existent key
	fmt.Println("\nAttempting to retrieve non-existent key:")
	_, err = ks.Get("0xnonexistent")
	if err != nil {
		fmt.Printf("Error (expected): %v\n", err)
	}

	// Create another keystore with different initial key
	fmt.Println("\nCreating another keystore with different initial key:")
	pubKey4 := "0xfedcba9876543210"
	privKey4 := new(big.Int).SetUint64(444444)
	ks2 := account.SetNewMemKeystore(pubKey4, privKey4)

	fmt.Printf("Created second keystore with key: Public=%s, Private=%s\n", pubKey4, privKey4)

	// Retrieve from second keystore
	retrieved4, err := ks2.Get(pubKey4)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Retrieved from keystore 2: %s\n", retrieved4)
	}

	// Verify keystores are independent
	fmt.Println("\nVerifying keystores are independent:")
	_, err1 := ks.Get(pubKey4)
	_, err2 := ks2.Get(pubKey1)
	fmt.Printf("Keystore 1 has key 4: %v\n", err1 == nil)
	fmt.Printf("Keystore 2 has key 1: %v\n", err2 == nil)
}