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

	// Create an invoke transaction V1
	contractAddress, _ := new(felt.Felt).SetString("0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7")
	entryPointSelector, _ := new(felt.Felt).SetString("0x83afd3f4caedc6eebf44246fe54e38c95e3179a5ec9ea81740eca5b482d12e")
	recipient, _ := new(felt.Felt).SetString("0x1234567890abcdef")

	invokeTx := &rpc.InvokeTxnV1{
		Type:          rpc.TransactionTypeInvoke,
		SenderAddress: accountAddress,
		Nonce:         new(felt.Felt).SetUint64(0),
		MaxFee:        new(felt.Felt).SetUint64(1000000000000),
		Version:       rpc.TransactionV1,
		Signature:     []*felt.Felt{}, // Empty before signing
		Calldata: []*felt.Felt{
			new(felt.Felt).SetUint64(1), // num calls
			contractAddress,
			entryPointSelector,
			new(felt.Felt).SetUint64(3), // calldata len
			recipient,
			new(felt.Felt).SetUint64(100), // amount
			new(felt.Felt).SetUint64(0),   // amount high
		},
	}

	fmt.Println("Before signing:")
	fmt.Printf("Signature length: %d\n", len(invokeTx.Signature))

	// Sign the transaction
	err = acc.SignInvokeTransaction(ctx, invokeTx)
	if err != nil {
		fmt.Printf("Error signing transaction: %v\n", err)
		return
	}

	fmt.Println("\nAfter signing:")
	fmt.Printf("Signature length: %d\n", len(invokeTx.Signature))
	for i, sig := range invokeTx.Signature {
		fmt.Printf("Signature[%d]: %s\n", i, sig)
	}
}
