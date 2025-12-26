package main
 
import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
 
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
 
	accountAddress := os.Getenv("ACCOUNT_ADDRESS")
	publicKey := os.Getenv("ACCOUNT_PUBLIC_KEY")
	privateKey := os.Getenv("ACCOUNT_PRIVATE_KEY")
 
	if accountAddress == "" || publicKey == "" || privateKey == "" {
		log.Fatal("ACCOUNT_ADDRESS, ACCOUNT_PUBLIC_KEY, or ACCOUNT_PRIVATE_KEY not set")
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
 
	nonce, err := accnt.Nonce(ctx)
	if err != nil {
		log.Fatal("Failed to get nonce:", err)
	}
 
	fmt.Printf("Account Address: %s\n", accountAddressFelt.String())
	fmt.Printf("Current Nonce:   %s\n", nonce.String())
	fmt.Printf("Nonce (uint64):  %d\n", nonce.Uint64())
	fmt.Printf("Nonce (hex):     0x%x\n", nonce.Uint64())
 
	noncePreConfirmed, err := provider.Nonce(
		ctx,
		rpc.WithBlockTag("pre_confirmed"),
		accountAddressFelt,
	)
	if err != nil {
		log.Fatal("Failed to get pre_confirmed nonce:", err)
	}
	fmt.Printf("Nonce (pre_confirmed): %s (%d)\n", noncePreConfirmed.String(), noncePreConfirmed.Uint64())
 
	nonceLatest, err := provider.Nonce(
		ctx,
		rpc.WithBlockTag("latest"),
		accountAddressFelt,
	)
	if err != nil {
		log.Fatal("Failed to get latest nonce:", err)
	}
	fmt.Printf("Nonce (latest):        %s (%d)\n", nonceLatest.String(), nonceLatest.Uint64())
 
	nonceL1Accepted, err := provider.Nonce(
		ctx,
		rpc.WithBlockTag("l1_accepted"),
		accountAddressFelt,
	)
	if err != nil {
		log.Fatal("Failed to get l1_accepted nonce:", err)
	}
	fmt.Printf("Nonce (l1_accepted):   %s (%d)\n", nonceL1Accepted.String(), nonceL1Accepted.Uint64())
}