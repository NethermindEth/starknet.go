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

	// Transaction hash prefix for invoke
	txHashPrefix := new(felt.Felt).SetBytes([]byte("invoke"))
	version := new(felt.Felt).SetUint64(0)
	contractAddress := new(felt.Felt).SetUint64(123)
	entryPointSelector := new(felt.Felt).SetUint64(456)

	// Hash the calldata
	calldata := new(felt.Felt).SetUint64(789)
	maxFee := new(felt.Felt).SetUint64(1000)
	additionalData := []*felt.Felt{}

	txHash := hash.CalculateDeprecatedTransactionHashCommon(
		txHashPrefix,
		version,
		contractAddress,
		entryPointSelector,
		calldata,
		maxFee,
		chainID,
		additionalData,
	)

	fmt.Printf("CalculateDeprecatedTransactionHashCommon: %s\n", txHash.String())
}
