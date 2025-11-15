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

	txn := &rpc.InvokeTxnV0{
		FunctionCall: rpc.FunctionCall{
			ContractAddress:    new(felt.Felt).SetUint64(123),
			EntryPointSelector: new(felt.Felt).SetUint64(456),
			Calldata:           []*felt.Felt{new(felt.Felt).SetUint64(1), new(felt.Felt).SetUint64(2)},
		},
		MaxFee:    new(felt.Felt).SetUint64(1000),
		Version:   rpc.TransactionV0,
		Signature: []*felt.Felt{},
	}

	txHash, _ := hash.TransactionHashInvokeV0(txn, chainID)
	fmt.Printf("TransactionHashInvokeV0: %s\n", txHash.String())
}
