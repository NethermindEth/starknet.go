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
		log.Fatal("STARKNET_RPC_URL not set")
	}
 
	provider, err := rpc.NewProvider(ctx, rpcURL)
	if err != nil {
		log.Fatal("Failed to create provider:", err)
	}
 
	accountAddress := os.Getenv("ACCOUNT_ADDRESS")
	publicKey := os.Getenv("ACCOUNT_PUBLIC_KEY")
	privateKey := os.Getenv("ACCOUNT_PRIVATE_KEY")
 
	if accountAddress == "" || publicKey == "" || privateKey == "" {
		log.Fatal("Account credentials not set in .env")
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
 
	accnt, err := account.NewAccount(provider, accountAddressFelt, publicKey, ks, account.CairoV2)
	if err != nil {
		log.Fatal("Failed to create account:", err)
	}
 
	message := new(felt.Felt).SetUint64(12345)
	signature, err := accnt.Sign(ctx, message)
	if err != nil {
		log.Fatal("Failed to sign message:", err)
	}
 
	isValid, err := accnt.Verify(message, signature)
	if err != nil {
		log.Fatal("Failed to verify signature:", err)
	}
 
	fmt.Printf("Message: %s\n", message.String())
	fmt.Printf("Signature: [%s, %s]\n", signature[0].String(), signature[1].String())
	fmt.Printf("Verification result: %v\n", isValid)
}