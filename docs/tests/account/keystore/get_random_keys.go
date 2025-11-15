package main

import (
	"fmt"

	"github.com/NethermindEth/starknet.go/account"
)

func main() {
	// Get random keys (keystore, public key, private key)
	ks, pubKey, privKey := account.GetRandomKeys()

	fmt.Println("Random keys generated:")
	fmt.Printf("Public Key:  %s\n", pubKey)
	fmt.Printf("Private Key: %s\n", privKey)

	// Show that keystore contains the key
	fmt.Println("\nVerifying keystore contains the key:")
	storedKey, err := ks.Get(pubKey.String())
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Retrieved key from keystore: %s\n", storedKey)
		// The keystore stores the private key associated with the public key
		fmt.Printf("Keys stored successfully: %v\n", storedKey != nil)
	}

	// Generate another set of random keys
	fmt.Println("\nGenerating another set of random keys:")
	ks2, pubKey2, privKey2 := account.GetRandomKeys()
	fmt.Printf("Public Key 2:  %s\n", pubKey2)
	fmt.Printf("Private Key 2: %s\n", privKey2)

	// Verify keys are different
	fmt.Printf("\nKeys are different: %v\n", pubKey.String() != pubKey2.String())

	// Show both keystores are independent
	_, err1 := ks.Get(pubKey2.String())
	_, err2 := ks2.Get(pubKey.String())
	fmt.Printf("Keystore 1 has key 2: %v\n", err1 == nil)
	fmt.Printf("Keystore 2 has key 1: %v\n", err2 == nil)
}