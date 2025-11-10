package main

import (
	"fmt"
	"log"
	"math/big"

	"github.com/NethermindEth/starknet.go/curve"
)

func main() {
	// Generate a key pair
	privKey, pubX, _, err := curve.GetRandomKeys()
	if err != nil {
		log.Fatal("Failed to generate keys:", err)
	}

	// Message hash to sign
	msgHash := new(big.Int)
	msgHash.SetString("1234567890abcdef", 16)

	// Sign the message
	r, s, err := curve.Sign(msgHash, privKey)
	if err != nil {
		log.Fatal("Failed to sign message:", err)
	}

	// Verify the signature
	valid, err := curve.Verify(msgHash, r, s, pubX)
	if err != nil {
		log.Fatal("Failed to verify signature:", err)
	}

	fmt.Println("Verify:")
	fmt.Printf("  Message Hash: 0x%x\n", msgHash)
	fmt.Printf("  Public Key X: 0x%x\n", pubX)
	fmt.Printf("  Signature R: 0x%x\n", r)
	fmt.Printf("  Signature S: 0x%x\n", s)
	fmt.Printf("  Valid: %v\n", valid)
}
