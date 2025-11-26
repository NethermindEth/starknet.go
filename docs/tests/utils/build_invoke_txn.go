package main

import (
	"fmt"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"
)

func main() {
	// Define example values
	senderAddress := new(felt.Felt).SetUint64(12345)
	nonce := new(felt.Felt).SetUint64(1)
	
	// Build calldata (example: transfer call)
	recipientAddr := new(felt.Felt).SetUint64(67890)
	amount := new(felt.Felt).SetUint64(1000)
	transferSelector := utils.GetSelectorFromNameFelt("transfer")
	calldata := []*felt.Felt{
		new(felt.Felt).SetUint64(1), // num calls
		new(felt.Felt).SetUint64(12345), // contract address
		transferSelector,
		recipientAddr,
		amount,
	}
	
	// Resource bounds
	resourceBounds := &rpc.ResourceBoundsMapping{
		L1Gas: rpc.ResourceBounds{
			MaxAmount:       rpc.U64(10000),
			MaxPricePerUnit: rpc.U128(100),
		},
	}

	// Build transaction
	txn := utils.BuildInvokeTxn(senderAddress, nonce, calldata, resourceBounds, nil)

	fmt.Printf("Built invoke transaction\n")
	fmt.Printf("  Sender: %s\n", txn.SenderAddress.String())
	fmt.Printf("  Nonce: %s\n", txn.Nonce.String())
	fmt.Println("✅ Transaction built successfully")
}
