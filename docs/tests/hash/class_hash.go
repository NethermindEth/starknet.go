package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/contracts"
	"github.com/NethermindEth/starknet.go/hash"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/joho/godotenv"
)

// main demonstrates how to calculate and verify contract class hashes
// using starknet.go v0.17.0
func main() {
	// Define paths to compiled contract artifacts
	sierraPath := "counter_contract/target/dev/counter_Counter.contract_class.json"
	casmPath := "counter_contract/target/dev/counter_Counter.compiled_contract_class.json"

	// Load and parse Sierra contract class
	sierraData, err := os.ReadFile(sierraPath)
	if err != nil {
		log.Fatalf("Failed to read Sierra file: %v", err)
	}

	var contractClass contracts.ContractClass
	if err := json.Unmarshal(sierraData, &contractClass); err != nil {
		log.Fatalf("Failed to parse Sierra JSON: %v", err)
	}

	// Calculate class hash from Sierra representation
	// This hash uniquely identifies the contract code on Starknet
	classHash := hash.ClassHash(&contractClass)
	fmt.Printf("Class Hash: %s\n", classHash.String())

	// Load and parse CASM (Cairo Assembly) compiled contract
	casmData, err := os.ReadFile(casmPath)
	if err != nil {
		log.Fatalf("Failed to read CASM file: %v", err)
	}

	var casmClass contracts.CasmClass
	if err := json.Unmarshal(casmData, &casmClass); err != nil {
		log.Fatalf("Failed to parse CASM JSON: %v", err)
	}

	// Calculate compiled class hash from CASM bytecode
	// This hash is required for declare transactions
	compiledHash, err := hash.CompiledClassHash(&casmClass)
	if err != nil {
		log.Fatalf("Failed to calculate compiled class hash: %v", err)
	}
	fmt.Printf("Compiled Class Hash: %s\n", compiledHash.String())

	// Optional: Verify contract declaration on-chain
	if err := godotenv.Load(); err == nil {
		if rpcURL := os.Getenv("STARKNET_RPC_URL"); rpcURL != "" {
			verifyOnChain(classHash, rpcURL)
		}
	}
}

// verifyOnChain checks if the contract is declared on the Starknet network
func verifyOnChain(classHash *felt.Felt, rpcURL string) {
	client, err := rpc.NewProvider(context.Background(), rpcURL)
	if err != nil {
		log.Printf("Warning: Could not connect to RPC: %v", err)
		return
	}

	ctx := context.Background()
	_, err = client.Class(ctx, rpc.WithBlockTag(rpc.BlockTagLatest), classHash)

	if err == nil {
		fmt.Printf("\nContract Status: DECLARED on-chain\n")
	} else {
		fmt.Printf("\nContract Status: NOT DECLARED on-chain\n")
	}
}
