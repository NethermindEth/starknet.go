package main
 
import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
 
	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/contracts"
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
 
	accountAddress := os.Getenv("ACCOUNT_ADDRESS")
	publicKey := os.Getenv("ACCOUNT_PUBLIC_KEY")
	privateKey := os.Getenv("ACCOUNT_PRIVATE_KEY")
 
	if accountAddress == "" || publicKey == "" || privateKey == "" {
		log.Fatal("ACCOUNT_ADDRESS, ACCOUNT_PUBLIC_KEY, or ACCOUNT_PRIVATE_KEY not set in .env")
	}
 
	ks := account.NewMemKeystore()
	privKeyBI, ok := new(big.Int).SetString(privateKey, 0)
	if !ok {
		log.Fatal("Failed to parse private key")
	}
	ks.Put(publicKey, privKeyBI)
 
	accountAddressFelt, err := utils.HexToFelt(accountAddress)
	if err != nil {
		log.Fatal("Failed to parse account address:", err)
	}
 
	accnt, err := account.NewAccount(
		provider,
		accountAddressFelt,
		publicKey,
		ks,
		account.CairoV2,
	)
	if err != nil {
		log.Fatal("Failed to create account:", err)
	}
 
	sierraFile, err := os.ReadFile("contract.sierra.json")
	if err != nil {
		log.Fatal("Failed to read sierra file:", err)
	}
 
	var contractClass contracts.ContractClass
	if err := json.Unmarshal(sierraFile, &contractClass); err != nil {
		log.Fatal("Failed to unmarshal contract class:", err)
	}
 
	casmClass, err := contracts.UnmarshalCasmClass("contract.casm.json")
	if err != nil {
		log.Fatal("Failed to unmarshal casm class:", err)
	}
 
	opts := &account.TxnOptions{
		FeeMultiplier: 1.5,
		TipMultiplier: 1.0,
	}
 
	response, err := accnt.BuildAndSendDeclareTxn(
		ctx,
		casmClass,
		&contractClass,
		opts,
	)
	if err != nil {
		log.Fatal("Failed to declare contract:", err)
	}
 
	fmt.Printf("Declare Transaction Successful:\n")
	fmt.Printf("Transaction Hash: %s\n", response.Hash.String())
	fmt.Printf("Class Hash:       %s\n", response.ClassHash.String())
	fmt.Printf("\nUse this class hash to deploy contract instances\n")
}