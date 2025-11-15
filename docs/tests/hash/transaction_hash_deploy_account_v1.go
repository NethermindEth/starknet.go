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

	txn := &rpc.DeployAccountTxnV1{
		ClassHash:           new(felt.Felt).SetUint64(789),
		ContractAddressSalt: new(felt.Felt).SetUint64(999),
		ConstructorCalldata: []*felt.Felt{new(felt.Felt).SetUint64(1)},
		MaxFee:              new(felt.Felt).SetUint64(1000),
		Version:             rpc.TransactionV1,
		Signature:           []*felt.Felt{},
		Nonce:               new(felt.Felt).SetUint64(0),
	}

	// Calculate contract address
	contractAddress := new(felt.Felt).SetUint64(12345)

	txHash, _ := hash.TransactionHashDeployAccountV1(txn, contractAddress, chainID)
	fmt.Printf("TransactionHashDeployAccountV1: %s\n", txHash.String())
}
