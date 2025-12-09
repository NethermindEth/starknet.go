package main

import (
	"fmt"
	"log"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/curve"
)

func main() {
	// Create message hash (typically from transaction hash)
	msgHash, _ := new(felt.Felt).SetString("0x1234567890abcdef")

	// Private key (never hardcode in production!)
	privKey, _ := new(felt.Felt).SetString("0x1234567890abcdef1234567890abcdef")

	// Sign the message
	r, s, err := curve.SignFelts(msgHash, privKey)
	if err != nil {
		log.Fatal("Signing failed:", err)
	}

	fmt.Println("SignFelts:")
	fmt.Printf("  Message Hash: %s\n", msgHash.String())
	fmt.Printf("  Signature R: %s\n", r.String())
	fmt.Printf("  Signature S: %s\n", s.String())
}
