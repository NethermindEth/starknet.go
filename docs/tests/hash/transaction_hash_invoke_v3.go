package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/hash"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load("../.env")

	// Get chain ID
	rpcURL := os.Getenv("STARKNET_RPC_URL")
	if rpcURL == "" {
		log.Fatal("STARKNET_RPC_URL not set in environment")
	}

	ctx := context.Background()
	client, err := rpc.NewProvider(ctx, rpcURL)
	if err != nil {
		log.Fatal("Failed to create RPC provider:", err)
	}

	chainID, err := client.ChainID(ctx)
	if err != nil {
		log.Fatal("Failed to get chain ID:", err)
	}

	// Create a sample InvokeTxnV3
	senderAddress, _ := new(felt.Felt).SetString("0x1234567890abcdef")

	txn := &rpc.InvokeTxnV3{
		SenderAddress: senderAddress,
		Calldata:      []*felt.Felt{new(felt.Felt).SetUint64(1)},
		Version:       rpc.TransactionV3,
		Signature:     []*felt.Felt{},
		Nonce:         new(felt.Felt).SetUint64(0),
		ResourceBounds: rpc.ResourceBoundsMapping{
			L1Gas: rpc.ResourceBounds{
				MaxAmount:       new(felt.Felt).SetUint64(10000),
				MaxPricePerUnit: new(felt.Felt).SetUint64(100),
			},
			L2Gas: rpc.ResourceBounds{
				MaxAmount:       new(felt.Felt).SetUint64(10000),
				MaxPricePerUnit: new(felt.Felt).SetUint64(100),
			},
		},
		Tip:                   0,
		PayMasterData:         []*felt.Felt{},
		AccountDeploymentData: []*felt.Felt{},
		NonceDataMode:         rpc.DAModeL1,
		FeeMode:               rpc.DAModeL1,
	}

	// Calculate transaction hash
	txHash, err := hash.TransactionHashInvokeV3(txn, chainID)
	if err != nil {
		log.Fatal("Failed to calculate transaction hash:", err)
	}

	fmt.Println("TransactionHashInvokeV3:")
	fmt.Printf("  Chain ID: %s\n", chainID)
	fmt.Printf("  Sender: %s\n", txn.SenderAddress.String())
	fmt.Printf("  Nonce: %s\n", txn.Nonce.String())
	fmt.Printf("  Transaction Hash: %s\n", txHash.String())
}
