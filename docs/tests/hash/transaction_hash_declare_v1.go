package main

import (
	"context"
	"fmt"
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
	chainIDStr, _ := client.ChainID(ctx)
	chainID := new(felt.Felt).SetBytes([]byte(chainIDStr))

	txn := &rpc.DeclareTxnV1{
		SenderAddress: new(felt.Felt).SetUint64(123),
		ClassHash:     new(felt.Felt).SetUint64(789),
		MaxFee:        new(felt.Felt).SetUint64(1000),
		Version:       rpc.TransactionV1,
		Signature:     []*felt.Felt{},
		Nonce:         new(felt.Felt).SetUint64(0),
	}

	txHash, _ := hash.TransactionHashDeclareV1(txn, chainID)
	fmt.Printf("TransactionHashDeclareV1: %s\n", txHash.String())
}
