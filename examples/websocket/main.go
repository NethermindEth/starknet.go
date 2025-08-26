package main

import (
	"context"
	"fmt"
	"time"

	setup "github.com/NethermindEth/starknet.go/examples/internal"
	"github.com/NethermindEth/starknet.go/rpc"
)

func main() {
	fmt.Println("Starting websocket example")

	// Load variables from '.env' file
	wsProviderUrl := setup.GetWsProviderUrl()

	// Initialise connection to WS provider
	wsClient, err := rpc.NewWebsocketProvider(wsProviderUrl)
	if err != nil {
		panic(fmt.Sprintf("Error dialling the WS provider: %s", err))
	}
	defer wsClient.Close() // Close the WS client when the program finishes

	fmt.Println("Established connection with the client")

	// Let's now call the SubscribeNewHeads method. To do this, we need to create a channel to receive the new heads.
	//
	// Note: We'll need to do this for each of the methods we want to subscribe to, always creating a channel to receive the values from
	// the node. Check each method's description for the type required for the channel.
	newHeadsChan := make(chan *rpc.BlockHeader)

	// We then call the desired websocket method, passing in the channel and the parameters if needed.
	// For example, to subscribe to new block headers, we call the SubscribeNewHeads method, passing in the channel and the blockID.
	// As the description says it's optional, we pass an empty BlockID as value. That way, the latest block will be used by default.
	sub, err := wsClient.SubscribeNewHeads(context.Background(), newHeadsChan, rpc.SubscriptionBlockID{})
	if err != nil {
		panic(err)
	}
	fmt.Println()
	fmt.Println("Successfully subscribed to the node. Subscription ID:", sub.ID())

	var latestBlockNumber uint64

	// Now we'll create the loop to continuously read the new heads from the channel.
	// This will make the program wait indefinitely for new heads or errors if not interrupted.
loop1:
	for {
		select {
		case newHead := <-newHeadsChan:
			// This case will be triggered when a new block header is received.
			fmt.Println("New block header received:", newHead.Number)
			latestBlockNumber = newHead.Number

			break loop1 // Let's exit the loop after receiving the first block header
		case err = <-sub.Err():
			// This case will be triggered when an error occurs.
			panic(err)
		}
	}

	// We can also use the subscription returned by the WS methods to unsubscribe from the stream when we're done
	sub.Unsubscribe()

	fmt.Printf("Unsubscribed from the subscription %s successfully\n", sub.ID())

	// We'll now subscribe to the node again, but this time we'll pass in an older block number as the blockID.
	// This way, the node will send us block headers from that block number onwards.
	sub, err = wsClient.SubscribeNewHeads(
		context.Background(),
		newHeadsChan,
		new(rpc.SubscriptionBlockID).WithBlockNumber(latestBlockNumber-10),
	)
	if err != nil {
		panic(err)
	}
	fmt.Println()
	fmt.Println("Successfully subscribed to the node. Subscription ID:", sub.ID())

	go func() {
		time.Sleep(10 * time.Second)
		// Unsubscribe from the subscription after 10 seconds
		sub.Unsubscribe()
	}()

loop2:
	for {
		select {
		case newHead := <-newHeadsChan:
			fmt.Println("New block header received:", newHead.Number)
		case err := <-sub.Err():
			if err == nil { // when sub.Unsubscribe() is called a nil error is returned, so let's just break the loop if that's the case
				fmt.Printf("Unsubscribed from the subscription %s successfully\n", sub.ID())

				break loop2
			}
			panic(err)
		}
	}

	// This example can be used to understand how to use all the methods that return a subscription.
	// It's just a matter of creating a channel to receive the values from the node and calling the
	// desired method, passing in the channel and the parameters if needed. Remember to check the method's
	// description for the type required for the channel and whether there are any other parameters needed.
}
