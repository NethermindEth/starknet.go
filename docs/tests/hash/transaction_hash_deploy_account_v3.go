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

	txn := &rpc.DeployAccountTxnV3{
		ClassHash:           new(felt.Felt).SetUint64(789),
		ContractAddressSalt: new(felt.Felt).SetUint64(999),
		ConstructorCalldata: []*felt.Felt{new(felt.Felt).SetUint64(1)},
		Version:             rpc.TransactionV3,
		Signature:           []*felt.Felt{},
		Nonce:               new(felt.Felt).SetUint64(0),
		ResourceBounds: &rpc.ResourceBoundsMapping{
			L1Gas: rpc.ResourceBounds{
				MaxAmount:       rpc.U64("0x186a0"),
				MaxPricePerUnit: rpc.U128("0x3e8"),
			},
			L1DataGas: rpc.ResourceBounds{
				MaxAmount:       rpc.U64("0x0"),
				MaxPricePerUnit: rpc.U128("0x0"),
			},
			L2Gas: rpc.ResourceBounds{
				MaxAmount:       rpc.U64("0x0"),
				MaxPricePerUnit: rpc.U128("0x0"),
			},
		},
		Tip:           rpc.U64("0x0"),
		PayMasterData: []*felt.Felt{},
		NonceDataMode: rpc.DAModeL1,
		FeeMode:       rpc.DAModeL1,
	}

	// Calculate contract address
	contractAddress := new(felt.Felt).SetUint64(12345)

	txHash, _ := hash.TransactionHashDeployAccountV3(txn, contractAddress, chainID)
	fmt.Printf("TransactionHashDeployAccountV3: %s\n", txHash.String())
}
