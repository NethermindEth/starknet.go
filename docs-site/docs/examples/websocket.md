---
sidebar_position: 7
---

# WebSocket Example

This example demonstrates how to use WebSocket connections with StarkNet.go to receive real-time updates from the StarkNet network.

## Overview

The WebSocket example shows how to:

1. Initialize a connection to a StarkNet WebSocket provider
2. Subscribe to new block headers
3. Process real-time updates
4. Unsubscribe from the subscription
5. Subscribe with specific parameters

## Prerequisites

Before running this example, you need to:

1. Rename the `.env.template` file located at the root of the "examples" folder to `.env`
2. Uncomment and assign your WebSocket provider URL to the `WS_PROVIDER_URL` variable in the `.env` file

## Code Example

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/NethermindEth/starknet.go/rpc"

	setup "github.com/NethermindEth/starknet.go/examples/internal"
)

func main() {
	fmt.Println("Starting websocket example")

	// Load variables from '.env' file
	wsProviderUrl := setup.GetWsProviderUrl()

	// Initialize connection to WS provider
	wsClient, err := rpc.NewWebsocketProvider(wsProviderUrl)
	if err != nil {
		panic(fmt.Sprintf("Error dialing the WS provider: %s", err))
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
	sub, err := wsClient.SubscribeNewHeads(context.Background(), newHeadsChan, rpc.BlockID{})
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
			fmt.Println("New block header received:", newHead.BlockNumber)
			latestBlockNumber = newHead.BlockNumber
			break loop1 // Let's exit the loop after receiving the first block header
		case err := <-sub.Err():
			// This case will be triggered when an error occurs.
			panic(err)
		}
	}

	// We can also use the subscription returned by the WS methods to unsubscribe from the stream when we're done
	sub.Unsubscribe()

	fmt.Printf("Unsubscribed from the subscription %s successfully\n", sub.ID())

	// We'll now subscribe to the node again, but this time we'll pass in an older block number as the blockID.
	// This way, the node will send us block headers from that block number onwards.
	sub, err = wsClient.SubscribeNewHeads(context.Background(), newHeadsChan, rpc.WithBlockNumber(latestBlockNumber-10))
	if err != nil {
		panic(err)
	}
	fmt.Println()
	fmt.Println("Successfully subscribed to the node. Subscription ID:", sub.ID())

	go func() {
		time.Sleep(20 * time.Second)
		// Unsubscribe from the subscription after 20 seconds
		sub.Unsubscribe()
	}()

loop2:
	for {
		select {
		case newHead := <-newHeadsChan:
			fmt.Println("New block header received:", newHead.BlockNumber)
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
```

## Running the Example

To run this example:

1. Make sure you are in the "websocket" directory
2. Execute `go run main.go`

## Expected Output

```
Starting websocket example
Established connection with the client

Successfully subscribed to the node. Subscription ID: 0x1
New block header received: 123456
Unsubscribed from the subscription 0x1 successfully

Successfully subscribed to the node. Subscription ID: 0x2
New block header received: 123447
New block header received: 123448
New block header received: 123449
New block header received: 123450
Unsubscribed from the subscription 0x2 successfully
```

## Key Concepts

### WebSocket Provider

The WebSocket provider allows for real-time updates from the StarkNet network:

```go
wsClient, err := rpc.NewWebsocketProvider(wsProviderUrl)
```

### Subscriptions

StarkNet.go supports several subscription types:

1. **New Block Headers**: Subscribe to new block headers
2. **Events**: Subscribe to contract events
3. **Pending Transactions**: Subscribe to pending transactions
4. **Transaction Status**: Subscribe to transaction status updates

### Creating Channels

For each subscription, you need to create a channel to receive the updates:

```go
newHeadsChan := make(chan *rpc.BlockHeader)
```

### Subscribing to Updates

To subscribe to updates:

```go
sub, err := wsClient.SubscribeNewHeads(context.Background(), newHeadsChan, rpc.BlockID{})
```

### Processing Updates

To process updates, you need to read from the channel:

```go
select {
case newHead := <-newHeadsChan:
    fmt.Println("New block header received:", newHead.BlockNumber)
case err := <-sub.Err():
    panic(err)
}
```

### Unsubscribing

To unsubscribe from a subscription:

```go
sub.Unsubscribe()
```

## Subscription Types

### Subscribe to New Block Headers

```go
newHeadsChan := make(chan *rpc.BlockHeader)
sub, err := wsClient.SubscribeNewHeads(context.Background(), newHeadsChan, rpc.BlockID{})
```

### Subscribe to Events

```go
eventsChan := make(chan *rpc.Event)
sub, err := wsClient.SubscribeEvents(
    context.Background(),
    eventsChan,
    rpc.EventsInput{
        Address: contractAddress,
        Keys:    [][]string{{"0x..."}}, // Event key to filter by
    },
)
```

### Subscribe to Pending Transactions

```go
pendingTxsChan := make(chan *string)
sub, err := wsClient.SubscribePendingTransactions(context.Background(), pendingTxsChan)
```

### Subscribe to Transaction Status

```go
txStatusChan := make(chan *rpc.TxnStatus)
sub, err := wsClient.SubscribeTransactionStatus(context.Background(), txStatusChan, txHash)
```

## Next Steps

After understanding how to use WebSocket connections, you can:

- Implement real-time updates in your application
- Create a block explorer with live updates
- Build a notification system for contract events
- Monitor transaction status in real-time
