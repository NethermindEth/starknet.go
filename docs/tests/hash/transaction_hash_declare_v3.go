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
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	chainID, _ := new(felt.Felt).SetString("0x534e5f5345504f4c4941")
	senderAddress, _ := new(felt.Felt).SetString("0x04718f5a0fc34cc1af16a1cdee98ffb20c31f5cd61d6ab07201858f4287c938d")
	classHash, _ := new(felt.Felt).SetString("0x01234567890abcdef1234567890abcdef1234567890abcdef1234567890abcd")
	compiledClassHash, _ := new(felt.Felt).SetString("0x0fedcba9876543210fedcba9876543210fedcba9876543210fedcba98765432")
	nonce := new(felt.Felt).SetUint64(15)

	resourceBounds := &rpc.ResourceBoundsMapping{
		L1Gas: rpc.ResourceBounds{
			MaxAmount:       "0x5000",
			MaxPricePerUnit: "0x33a937098d80",
		},
		L1DataGas: rpc.ResourceBounds{
			MaxAmount:       "0x5000",
			MaxPricePerUnit: "0x33a937098d80",
		},
		L2Gas: rpc.ResourceBounds{
			MaxAmount:       "0x5000",
			MaxPricePerUnit: "0x10c388d00",
		},
	}

	txn := &rpc.DeclareTxnV3{
		Type:                  rpc.TransactionTypeDeclare,
		Version:               rpc.TransactionV3,
		SenderAddress:         senderAddress,
		ClassHash:             classHash,
		CompiledClassHash:     compiledClassHash,
		Nonce:                 nonce,
		ResourceBounds:        resourceBounds,
		Tip:                   "0x0",
		PayMasterData:         []*felt.Felt{},
		AccountDeploymentData: []*felt.Felt{},
		NonceDataMode:         rpc.DAModeL1,
		FeeMode:               rpc.DAModeL1,
		Signature:             []*felt.Felt{},
	}

	txHash, err := hash.TransactionHashDeclareV3(txn, chainID)
	if err != nil {
		log.Fatalf("Failed to calculate transaction hash: %v", err)
	}

	fmt.Printf("Transaction Hash: %s\n", txHash.String())

	if rpcURL := os.Getenv("STARKNET_RPC_URL"); rpcURL != "" {
		verifyTransaction(txHash, rpcURL)
	}
}

func verifyTransaction(txHash *felt.Felt, rpcURL string) {
	client, err := rpc.NewProvider(context.Background(), rpcURL)
	if err != nil {
		log.Printf("Warning: Could not connect to RPC: %v", err)
		return
	}

	ctx := context.Background()
	tx, err := client.TransactionByHash(ctx, txHash)

	if err == nil {
		fmt.Printf("\nVerification: FOUND on-chain\n")
		fmt.Printf("Type: %T\n", tx)
	} else {
		fmt.Printf("\nVerification: NOT FOUND\n")
	}
}
