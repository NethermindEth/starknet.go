package main

import (
	"fmt"
	"log"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/curve"
)

func main() {
	// Generate a key pair
	privKeyBig, pubXBig, _, err := curve.GetRandomKeys()
	if err != nil {
		log.Fatal("Failed to generate keys:", err)
	}

	// Convert to felt
	privKey := new(felt.Felt).SetBigInt(privKeyBig)
	pubX := new(felt.Felt).SetBigInt(pubXBig)

	// Message hash to sign
	msgHash, _ := new(felt.Felt).SetString("0x1234567890abcdef")

	// Sign the message using felts
	r, s, err := curve.SignFelts(msgHash, privKey)
	if err != nil {
		log.Fatal("Failed to sign message:", err)
	}

	fmt.Println("SignFelts:")
	fmt.Printf("  Message Hash: %s\n", msgHash.String())
	fmt.Printf("  Private Key: %s\n", privKey.String())
	fmt.Printf("  Public Key X: %s\n", pubX.String())
	fmt.Printf("  Signature R: %s\n", r.String())
	fmt.Printf("  Signature S: %s\n", s.String())
}
