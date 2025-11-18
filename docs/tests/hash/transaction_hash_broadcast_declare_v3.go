package main

import (
	"context"
	"fmt"
	"os"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/contracts"
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

	senderAddress, _ := new(felt.Felt).SetString("0x123")

	txn := &rpc.BroadcastDeclareTxnV3{
		SenderAddress:     senderAddress,
		CompiledClassHash: new(felt.Felt).SetUint64(456),
		Version:           rpc.TransactionV3,
		Signature:         []*felt.Felt{},
		Nonce:             new(felt.Felt).SetUint64(0),
		ContractClass:     &contracts.ContractClass{}, // Empty contract class for testing
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
		Tip:                   rpc.U64("0x0"),
		PayMasterData:         []*felt.Felt{},
		AccountDeploymentData: []*felt.Felt{},
		NonceDataMode:         rpc.DAModeL1,
		FeeMode:               rpc.DAModeL1,
	}

	txHash, _ := hash.TransactionHashBroadcastDeclareV3(txn, chainID)
	fmt.Printf("TransactionHashBroadcastDeclareV3: %s\n", txHash.String())
}
