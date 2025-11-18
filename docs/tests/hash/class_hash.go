package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/NethermindEth/starknet.go/contracts"
	"github.com/NethermindEth/starknet.go/hash"
)

func main() {
	// Read a sample Sierra contract class
	// For demonstration, we'll use a minimal contract structure
	contractJSON := `{
		"sierra_program": ["0x1", "0x2", "0x3"],
		"contract_class_version": "0.1.0",
		"entry_points_by_type": {
			"CONSTRUCTOR": [],
			"EXTERNAL": [],
			"L1_HANDLER": []
		},
		"abi": []
	}`

	var contractClass contracts.ContractClass
	err := json.Unmarshal([]byte(contractJSON), &contractClass)
	if err != nil {
		log.Fatal("Failed to parse contract class:", err)
	}

	// Calculate class hash
	classHash := hash.ClassHash(&contractClass)
	if classHash == nil {
		log.Fatal("Failed to calculate class hash")
	}

	fmt.Println("ClassHash:")
	fmt.Printf("  Contract Class Version: %s\n", contractClass.ContractClassVersion)
	fmt.Printf("  Class Hash: %s\n", classHash.String())

	// Note: In real usage, load from actual Sierra JSON file
	fmt.Println("\nNote: Use actual Sierra contract JSON for real class hash calculation")
}
