package main

import (
	"fmt"
	"log"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/contracts"
	"github.com/NethermindEth/starknet.go/hash"
)

func main() {
	// Create a minimal CASM class
	casmClass := &contracts.CasmClass{
		Prime:           "0x800000000000011000000000000000000000000000000000000000000000001",
		CompilerVersion: "2.1.0",
		Bytecode:        []*felt.Felt{new(felt.Felt).SetUint64(1), new(felt.Felt).SetUint64(2)},
	}

	// Calculate compiled class hash
	compiledHash, err := hash.CompiledClassHash(casmClass)
	if err != nil {
		log.Fatal("Failed to calculate compiled class hash:", err)
	}

	fmt.Println("CompiledClassHash:")
	fmt.Printf("  Compiler Version: %s\n", casmClass.CompilerVersion)
	fmt.Printf("  Bytecode Length: %d\n", len(casmClass.Bytecode))
	fmt.Printf("  Compiled Class Hash: %s\n", compiledHash.String())

	fmt.Println("\nNote: Use actual CASM JSON for real compiled class hash")
}
