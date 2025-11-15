package main

import (
	"fmt"
	"math/big"

	"github.com/NethermindEth/starknet.go/account"
)

func main() {
	// Create an empty keystore
	ks := account.NewMemKeystore()
	fmt.Println("Created empty MemKeystore")

	// Put a single key
	fmt.Println("\nAdding first key:")
	pubKey1 := "0x1234567890abcdef"
	privKey1 := new(big.Int).SetUint64(99999)
	ks.Put(pubKey1, privKey1)
	fmt.Printf("Added key: public=%s, private=%s\n", pubKey1, privKey1)

	// Verify key was added
	retrieved1, err := ks.Get(pubKey1)
	if err != nil {
		fmt.Printf("Error retrieving key: %v\n", err)
	} else {
		fmt.Printf("Verified key was stored: %s\n", retrieved1)
	}

	// Put multiple keys
	fmt.Println("\nAdding multiple keys:")
	keys := map[string]*big.Int{
		"0xabc123":           new(big.Int).SetUint64(111),
		"0xdef456":           new(big.Int).SetUint64(222),
		"0x789xyz":           new(big.Int).SetUint64(333),
		"0xfedcba9876543210": new(big.Int).SetUint64(444),
	}

	for pub, priv := range keys {
		ks.Put(pub, priv)
		fmt.Printf("Added key: public=%s, private=%s\n", pub, priv)
	}

	// Verify all keys are stored
	fmt.Println("\nVerifying all keys are stored:")
	for pub, expectedPriv := range keys {
		retrieved, err := ks.Get(pub)
		if err != nil {
			fmt.Printf("Error getting %s: %v\n", pub, err)
		} else {
			matches := retrieved.String() == expectedPriv.String()
			fmt.Printf("Key %s: retrieved=%s, matches=%v\n", pub, retrieved, matches)
		}
	}

	// Overwrite an existing key
	fmt.Println("\nOverwriting existing key:")
	newPrivKey1 := new(big.Int).SetUint64(88888)
	fmt.Printf("Original key %s had value: %s\n", pubKey1, privKey1)
	ks.Put(pubKey1, newPrivKey1)
	fmt.Printf("Overwrote with new value: %s\n", newPrivKey1)

	// Verify the key was overwritten
	retrievedNew, err := ks.Get(pubKey1)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Retrieved updated key: %s\n", retrievedNew)
		fmt.Printf("Key was successfully overwritten: %v\n", retrievedNew.String() == newPrivKey1.String())
	}

	// Put with empty string key (edge case)
	fmt.Println("\nTesting edge cases:")
	emptyKey := ""
	emptyPriv := new(big.Int).SetUint64(555)
	ks.Put(emptyKey, emptyPriv)
	fmt.Printf("Added key with empty string: public='%s', private=%s\n", emptyKey, emptyPriv)

	retrievedEmpty, err := ks.Get(emptyKey)
	if err != nil {
		fmt.Printf("Error retrieving empty key: %v\n", err)
	} else {
		fmt.Printf("Retrieved empty string key: %s\n", retrievedEmpty)
	}

	// Put with very large number
	fmt.Println("\nTesting with large number:")
	largePub := "0xlarge"
	largePriv := new(big.Int)
	largePriv.SetString("123456789012345678901234567890123456789012345678901234567890", 10)
	ks.Put(largePub, largePriv)
	fmt.Printf("Added large key: public=%s\n", largePub)
	fmt.Printf("Private key: %s\n", largePriv)

	retrievedLarge, err := ks.Get(largePub)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Retrieved large key successfully: %v\n", retrievedLarge.String() == largePriv.String())
	}
}