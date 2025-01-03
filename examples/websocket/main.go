package main

import (
	"context"
	"fmt"

	"github.com/NethermindEth/starknet.go/rpc"

	setup "github.com/NethermindEth/starknet.go/examples/internal"
)

// main entry point of the program.
//
// It initializes the environment and establishes a connection with the client.
// It then makes two contract calls and prints the responses.
//
// Parameters:
//
//	none
//
// Returns:
//
//	none
func main() {
	fmt.Println("Starting simpleCall example")

	// Load variables from '.env' file
	rpcProviderUrl := setup.GetRpcProviderUrl()

	// Initialize connection to RPC provider
	client, err := rpc.NewWebsocketProvider(rpcProviderUrl)
	if err != nil {
		panic(fmt.Sprintf("Error dialing the RPC provider: %s", err))
	}

	fmt.Println("Established connection with the client")

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		panic(fmt.Sprintf("Error getting chain ID: %s", err))
	}
	fmt.Printf("Chain ID: %s\n", chainID)

}
