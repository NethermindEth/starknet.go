package main

import (
	"fmt"
	"log"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/curve"
)

func main() {
	// Message hash
	msgHash, _ := new(felt.Felt).SetString("0x1234567890abcdef")

	// Signature components (from SignFelts)
	r, _ := new(felt.Felt).SetString("0x4d2b6e6e88e01af828b0f68237ffde7e6742ae86169a89c32185141ad1c6e7e")
	s, _ := new(felt.Felt).SetString("0x4d6f5fe7927e73ffddd5bca9ea3be17e0ae62d12e34c772691fc7b829904f92")

	// Public key X coordinate
	pubX, _ := new(felt.Felt).SetString("0x1ef15c18599971b7beced415a40f0c7deacfd9b0d1819e03d723d8bc943cfca")

	// Verify signature
	valid, err := curve.VerifyFelts(msgHash, r, s, pubX)
	if err != nil {
		log.Fatal("Verification failed:", err)
	}

	fmt.Println("VerifyFelts:")
	fmt.Printf("  Message Hash: %s\n", msgHash.String())
	fmt.Printf("  Signature R: %s\n", r.String())
	fmt.Printf("  Signature S: %s\n", s.String())
	fmt.Printf("  Public Key X: %s\n", pubX.String())
	fmt.Printf("  Valid: %v\n", valid)
}
