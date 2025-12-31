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
	classHash, _ := new(felt.Felt).SetString("0x01234567890abcdef1234567890abcdef1234567890abcdef1234567890abcd")
	salt := new(felt.Felt).SetUint64(12345)
	nonce := new(felt.Felt).SetUint64(0)

	publicKey, _ := new(felt.Felt).SetString("0x01234567890abcdef1234567890abcdef1234567890abcdef1234567890abcd")
	constructorCalldata := []*felt.Felt{publicKey}

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

	txn := &rpc.DeployAccountTxnV3{
		Type:                rpc.TransactionTypeDeployAccount,
		Version:             rpc.TransactionV3,
		ClassHash:           classHash,
		ContractAddressSalt: salt,
		ConstructorCalldata: constructorCalldata,
		Nonce:               nonce,
		ResourceBounds:      resourceBounds,
		Tip:                 "0x0",
		PayMasterData:       []*felt.Felt{},
		NonceDataMode:       rpc.DAModeL1,
		FeeMode:             rpc.DAModeL1,
		Signature:           []*felt.Felt{},
	}

	contractAddress, _ := new(felt.Felt).SetString("0x04718f5a0fc34cc1af16a1cdee98ffb20c31f5cd61d6ab07201858f4287c938d")

	txHash, err := hash.TransactionHashDeployAccountV3(txn, contractAddress, chainID)
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
