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
	
	ctx := context.Background()
	client, _ := rpc.NewProvider(ctx, os.Getenv("STARKNET_RPC_URL"))
	chainID, _ := client.ChainID(ctx)

	txn := &rpc.InvokeTxnV1{
		SenderAddress: new(felt.Felt).SetUint64(123),
		Calldata:      []*felt.Felt{new(felt.Felt).SetUint64(1)},
		MaxFee:        new(felt.Felt).SetUint64(1000),
		Version:       rpc.TransactionV1,
		Signature:     []*felt.Felt{},
		Nonce:         new(felt.Felt).SetUint64(0),
	}

	txHash, _ := hash.TransactionHashInvokeV1(txn, chainID)
	fmt.Printf("TransactionHashInvokeV1: %s\n", txHash.String())
}
