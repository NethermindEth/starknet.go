package main
 
import (
	"context"
	"fmt"
	"log"
	"os"
 
	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"
	"github.com/joho/godotenv"
)
 
func main() {
	ctx := context.Background()
 
	if err := godotenv.Load(); err != nil {
		log.Fatal("Failed to load .env file:", err)
	}
 
	rpcURL := os.Getenv("STARKNET_RPC_URL")
	if rpcURL == "" {
		log.Fatal("STARKNET_RPC_URL not set in .env file")
	}
 
	provider, err := rpc.NewProvider(ctx, rpcURL)
	if err != nil {
		log.Fatal("Failed to create provider:", err)
	}
 
	ks, publicKey, privateKey := account.GetRandomKeys()
	fmt.Printf("Generated keys:\n")
	fmt.Printf("  Public Key:  %s\n", publicKey.String())
	fmt.Printf("  Private Key: %s\n\n", privateKey)
 
	accnt, err := account.NewAccount(
		provider,
		publicKey,
		publicKey.String(),
		ks,
		account.CairoV2,
	)
	if err != nil {
		log.Fatal("Failed to create account:", err)
	}
 
	classHash, err := utils.HexToFelt("0x061dac032f228abef9c6626f995015233097ae253a7f72d68552db02f2971b8f")
	if err != nil {
		log.Fatal("Failed to parse class hash:", err)
	}
 
	constructorCalldata := []*felt.Felt{publicKey}
 
	salt := publicKey
 
	opts := &account.TxnOptions{
		FeeMultiplier: 1.5,
		TipMultiplier: 1.0,
	}
 
	deployTxn, precomputedAddress, err := accnt.BuildAndEstimateDeployAccountTxn(
		ctx,
		salt,
		classHash,
		constructorCalldata,
		opts,
	)
	if err != nil {
		log.Fatal("Failed to build and estimate deploy account transaction:", err)
	}
 
	fmt.Printf("\nDeploy Account Transaction:\n")
	fmt.Printf("Type:                %s\n", deployTxn.Type)
	fmt.Printf("Version:             %s\n", deployTxn.Version)
	fmt.Printf("Precomputed Address: %s\n", precomputedAddress.String())
	fmt.Printf("Class Hash:          %s\n", deployTxn.ClassHash.String())
	fmt.Printf("Salt:                %s\n", deployTxn.ContractAddressSalt.String())
	fmt.Printf("Nonce:               %s\n", deployTxn.Nonce.String())
	fmt.Printf("Signature Length:    %d\n", len(deployTxn.Signature))
 
	fmt.Printf("\nResource Bounds:\n")
	fmt.Printf("L1 Gas:\n")
	fmt.Printf("  Max Amount:         %s\n", deployTxn.ResourceBounds.L1Gas.MaxAmount)
	fmt.Printf("  Max Price Per Unit: %s\n", deployTxn.ResourceBounds.L1Gas.MaxPricePerUnit)
	fmt.Printf("L2 Gas:\n")
	fmt.Printf("  Max Amount:         %s\n", deployTxn.ResourceBounds.L2Gas.MaxAmount)
	fmt.Printf("  Max Price Per Unit: %s\n", deployTxn.ResourceBounds.L2Gas.MaxPricePerUnit)
	fmt.Printf("Tip:                  %s\n\n", deployTxn.Tip)
 
	overallFee, err := utils.ResBoundsMapToOverallFee(
		deployTxn.ResourceBounds,
		1,
		deployTxn.Tip,
	)
	if err != nil {
		log.Fatal("Failed to calculate overall fee:", err)
	}
 
	fmt.Printf("Estimated Fee: %s STRK\n\n", overallFee.String())
 
	fmt.Printf("Fund the precomputed address with STRK tokens:\n")
	fmt.Printf("   Address: %s\n", precomputedAddress.String())
	fmt.Printf("   Amount:  At least %s STRK\n\n", overallFee.String())
}