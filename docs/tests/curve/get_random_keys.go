package main

import (
	"fmt"
	"log"

	"github.com/NethermindEth/starknet.go/curve"
)

func main() {
	// Generate random keys
	privKey, x, y, err := curve.GetRandomKeys()
	if err != nil {
		log.Fatal("Failed to generate random keys:", err)
	}

	fmt.Println("GetRandomKeys:")
	fmt.Printf("  Private Key: 0x%x\n", privKey)
	fmt.Printf("  Public Key X: 0x%x\n", x)
	fmt.Printf("  Public Key Y: 0x%x\n", y)
	fmt.Println("\nNote: Keys are randomly generated, values will differ on each run")
}
