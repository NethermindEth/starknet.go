package main

import (
	"context"
	"fmt"

	"github.com/NethermindEth/starknet.go/rpc"
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

	// Initialize connection to RPC provider
	client, err := rpc.NewWebsocketProvider("ws://localhost:6061") //local juno node for testing
	if err != nil {
		panic(fmt.Sprintf("Error dialing the RPC provider: %s", err))
	}

	fmt.Println("Established connection with the client")

	ch := make(chan *rpc.BlockHeader)
	sub, err := client.SubscribeNewHeads(context.Background(), ch)
	if err != nil {
		rpcErr := err.(*rpc.RPCError)
		panic(fmt.Sprintf("Error subscribing: %s", rpcErr.Error()))
	}

	for {
		select {
		case resp := <-ch:
			fmt.Printf("New block: %d \n", resp.BlockNumber)
		case err := <-sub.Err():
			panic(fmt.Sprintf("Error subscribing to new heads: %s", err))
		}
	}
}
