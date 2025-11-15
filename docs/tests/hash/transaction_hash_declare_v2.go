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

	txn := &rpc.DeclareTxnV2{
		SenderAddress:     new(felt.Felt).SetUint64(123),
		CompiledClassHash: new(felt.Felt).SetUint64(456),
		MaxFee:            new(felt.Felt).SetUint64(1000),
		Version:           rpc.TransactionV2,
		Signature:         []*felt.Felt{},
		Nonce:             new(felt.Felt).SetUint64(0),
		ClassHash:         new(felt.Felt).SetUint64(789),
	}

	txHash, _ := hash.TransactionHashDeclareV2(txn, chainID)
	fmt.Printf("TransactionHashDeclareV2: %s\n", txHash.String())
}
