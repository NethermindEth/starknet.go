package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	rpcURL := os.Getenv("STARKNET_RPC_URL")
	if rpcURL == "" {
		log.Fatal("STARKNET_RPC_URL not set in .env")
	}

	ctx := context.Background()
	provider, err := rpc.NewProvider(ctx, rpcURL)
	if err != nil {
		log.Fatal(err)
	}

	accountAddress, _ := new(felt.Felt).SetString("0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7")

	publicKey := "0x03603a2692a2ae60abb343e832ee53b55d6b25f02a3ef1565ec691edc7a209b2"
	privateKey := new(big.Int).SetUint64(123456789)
	ks := account.SetNewMemKeystore(publicKey, privateKey)

	acc, err := account.NewAccount(provider, accountAddress, publicKey, ks, account.CairoV2)
	if err != nil {
		log.Fatal(err)
	}

	// Create a declare transaction V2
	classHash, _ := new(felt.Felt).SetString("0x01f372292df22d28f2d4c5798734421afe9596e6a566b8bc9b7b50e26521b855")
	compiledClassHash, _ := new(felt.Felt).SetString("0x017f655f7a639a49ea1d8d56172e99cff8b51f4123b733f0378dfd6378a2cd37")

	declareTx := &rpc.DeclareTxnV2{
		Type:              rpc.TransactionTypeDeclare,
		SenderAddress:     accountAddress,
		CompiledClassHash: compiledClassHash,
		ClassHash:         classHash,
		MaxFee:            new(felt.Felt).SetUint64(1000000000000),
		Version:           rpc.TransactionV2,
		Signature:         []*felt.Felt{}, // Empty before signing
		Nonce:             new(felt.Felt).SetUint64(0),
	}

	fmt.Printf("Before signing - Signature length: %d\n", len(declareTx.Signature))

	// Sign the transaction
	err = acc.SignDeclareTransaction(ctx, declareTx)
	if err != nil {
		fmt.Printf("Error signing transaction: %v\n", err)
		return
	}

	fmt.Printf("After signing - Signature length: %d\n", len(declareTx.Signature))
	for i, sig := range declareTx.Signature {
		fmt.Printf("Signature[%d]: %s\n", i, sig)
	}
}
