package main

import (
	"fmt"
	"math/big"

	"github.com/NethermindEth/starknet.go/curve"
)

func main() {
	// Use a sample private key
	privKey := new(big.Int)
	privKey.SetString("1234567890123456789", 10)

	// Get the public key point
	x, y := curve.PrivateKeyToPoint(privKey)

	fmt.Println("PrivateKeyToPoint:")
	fmt.Printf("  Private Key: %s\n", privKey.String())
	fmt.Printf("  Public Key X: 0x%x\n", x)
	fmt.Printf("  Public Key Y: 0x%x\n", y)
}
