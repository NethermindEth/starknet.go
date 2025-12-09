package main
 
import (
	"context"
	"fmt"
	"log"
	"os"
 
	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/joho/godotenv"
)
 
func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
 
	// Get RPC URL from environment variable
	rpcURL := os.Getenv("STARKNET_RPC_URL")
	if rpcURL == "" {
		log.Fatal("STARKNET_RPC_URL not found in .env file")
	}
 
	// Initialize provider
	ctx := context.Background()
	client, err := rpc.NewProvider(ctx, rpcURL)
	if err != nil {
		log.Fatal(err)
	}
 
	// Class hash of a declared contract
	classHash, _ := new(felt.Felt).SetString("0x05400e90f7e0ae78bd02c77cd75527280470e2fe19c54970dd79dc37a9d3645c")
 
	// Get compiled CASM
	casmClass, err := client.CompiledCasm(ctx, classHash)
	if err != nil {
		log.Fatal(err)
	}
 
	fmt.Printf("Prime: %s\n", casmClass.Prime)
	fmt.Printf("Compiler Version: %s\n", casmClass.CompilerVersion)
	fmt.Printf("Bytecode Length: %d\n", len(casmClass.ByteCode))
	fmt.Printf("External Entry Points: %d\n", len(casmClass.EntryPointsByType.External))
}