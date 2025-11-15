package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/rpc"
)

func main() {
	rpcURL := os.Getenv("STARKNET_RPC_URL")
	if rpcURL == "" {
		rpcURL = "http://localhost:5050/rpc"
	}

	client, err := rpc.NewProvider(context.Background(), rpcURL)
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}

	// Use a known class hash (this one is from Starknet)
	classHashStr := "0x05400e90f7e0ae78bd02c77cd75527280470e2fe19c54970dd79dc37a9d3645c"
	classHash, err := new(felt.Felt).SetString(classHashStr)
	if err != nil {
		log.Fatal("Invalid class hash:", err)
	}

	fmt.Println("CompiledCasm:")
	fmt.Printf("  Class Hash: %s\n", classHash.String())

	casmClass, err := client.CompiledCasm(context.Background(), classHash)
	if err != nil {
		log.Fatal("Failed to get compiled CASM:", err)
	}

	fmt.Printf("  Prime: %s\n", casmClass.Prime)
	fmt.Printf("  Compiler Version: %s\n", casmClass.CompilerVersion)
	fmt.Printf("  Bytecode Length: %d\n", len(casmClass.ByteCode))
	fmt.Printf("  Entry Points External: %d\n", len(casmClass.EntryPointsByType.External))
	fmt.Printf("  Entry Points L1Handler: %d\n", len(casmClass.EntryPointsByType.L1Handler))
	fmt.Printf("  Entry Points Constructor: %d\n", len(casmClass.EntryPointsByType.Constructor))
}
