package main

import (
	"fmt"
	"math/big"

	"github.com/NethermindEth/starknet.go/curve"
)

func main() {
	// Known private key
	privKey := big.NewInt(1234567890123456789)

	// Derive public key
	x, y := curve.PrivateKeyToPoint(privKey)

	fmt.Println("PrivateKeyToPoint:")
	fmt.Printf("  Private Key: 0x%x\n", privKey)
	fmt.Printf("  Public Key X: 0x%x\n", x)
	fmt.Printf("  Public Key Y: 0x%x\n", y)
}
