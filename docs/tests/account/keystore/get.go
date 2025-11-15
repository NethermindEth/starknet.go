package main

import (
	"fmt"
	"math/big"

	"github.com/NethermindEth/starknet.go/account"
)

func main() {
	// Create a keystore and add some keys
	ks := account.NewMemKeystore()

	// Add keys to the keystore
	pubKey1 := "0x0123456789abcdef"
	privKey1 := new(big.Int).SetUint64(11111)
	ks.Put(pubKey1, privKey1)

	pubKey2 := "0xfedcba9876543210"
	privKey2 := new(big.Int).SetUint64(22222)
	ks.Put(pubKey2, privKey2)

	pubKey3 := "0xaabbccddeeff0011"
	privKey3 := new(big.Int).SetUint64(33333)
	ks.Put(pubKey3, privKey3)

	fmt.Println("Keystore populated with 3 keys")

	// Get existing keys
	fmt.Println("\nRetrieving existing keys:")

	key1, err := ks.Get(pubKey1)
	if err != nil {
		fmt.Printf("Error getting key 1: %v\n", err)
	} else {
		fmt.Printf("Key 1 (public: %s): %s\n", pubKey1, key1)
	}

	key2, err := ks.Get(pubKey2)
	if err != nil {
		fmt.Printf("Error getting key 2: %v\n", err)
	} else {
		fmt.Printf("Key 2 (public: %s): %s\n", pubKey2, key2)
	}

	key3, err := ks.Get(pubKey3)
	if err != nil {
		fmt.Printf("Error getting key 3: %v\n", err)
	} else {
		fmt.Printf("Key 3 (public: %s): %s\n", pubKey3, key3)
	}

	// Try to get non-existent keys
	fmt.Println("\nAttempting to retrieve non-existent keys:")

	_, err = ks.Get("0xnonexistent")
	if err != nil {
		fmt.Printf("Error for '0xnonexistent': %v\n", err)
	}

	_, err = ks.Get("")
	if err != nil {
		fmt.Printf("Error for empty string: %v\n", err)
	}

	_, err = ks.Get("0x0000000000000000")
	if err != nil {
		fmt.Printf("Error for '0x0000000000000000': %v\n", err)
	}

	// Verify keys are retrieved correctly
	fmt.Println("\nVerifying retrieved keys match original:")
	fmt.Printf("Key 1 matches: %v\n", key1.String() == privKey1.String())
	fmt.Printf("Key 2 matches: %v\n", key2.String() == privKey2.String())
	fmt.Printf("Key 3 matches: %v\n", key3.String() == privKey3.String())

	// Test case sensitivity
	fmt.Println("\nTesting case sensitivity:")

	// Try uppercase version of existing key
	upperKey := "0X0123456789ABCDEF"
	_, err = ks.Get(upperKey)
	if err != nil {
		fmt.Printf("Uppercase key not found (keys are case-sensitive): %v\n", err)
	}

	// Mixed case
	mixedKey := "0x0123456789AbCdEf"
	_, err = ks.Get(mixedKey)
	if err != nil {
		fmt.Printf("Mixed case key not found (keys are case-sensitive): %v\n", err)
	}
}