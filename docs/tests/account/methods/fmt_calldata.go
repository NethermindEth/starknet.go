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
 
	strkContract, _ := utils.HexToFelt("0x04718f5a0fc34cc1af16a1cdee98ffb20c31f5cd61d6ab07201858f4287c938d")
	transferSelector := utils.GetSelectorFromNameFelt("transfer")
	recipient := accountAddressFelt
	amount := new(felt.Felt).SetUint64(1000000000000000)
	u256Amount, _ := utils.HexToU256Felt(amount.String())
 
	fnCall := rpc.FunctionCall{
		ContractAddress:    strkContract,
		EntryPointSelector: transferSelector,
		Calldata:           append([]*felt.Felt{recipient}, u256Amount...),
	}
 
	accntV2, err := account.NewAccount(
		provider,
		accountAddressFelt,
		publicKey,
		ks,
		account.CairoV2,
	)
	if err != nil {
		log.Fatal("Failed to create Cairo v2 account:", err)
	}
 
	calldataV2, err := accntV2.FmtCalldata([]rpc.FunctionCall{fnCall})
	if err != nil {
		log.Fatal("Failed to format Cairo v2 calldata:", err)
	}
 
	fmt.Println("Cairo v2 Formatted Calldata:")
	fmt.Printf("Total elements: %d\n", len(calldataV2))
	fmt.Printf("  [0] Num calls:         %s\n", calldataV2[0].String())
	fmt.Printf("  [1] Contract address:  %s\n", calldataV2[1].String())
	fmt.Printf("  [2] Entry point sel:   %s\n", calldataV2[2].String())
	fmt.Printf("  [3] Calldata length:   %s\n", calldataV2[3].String())
 
	accntV0, err := account.NewAccount(
		provider,
		accountAddressFelt,
		publicKey,
		ks,
		account.CairoV0,
	)
	if err != nil {
		log.Fatal("Failed to create Cairo v0 account:", err)
	}
 
	calldataV0, err := accntV0.FmtCalldata([]rpc.FunctionCall{fnCall})
	if err != nil {
		log.Fatal("Failed to format Cairo v0 calldata:", err)
	}
 
	fmt.Println("Cairo v0 Formatted Calldata:")
	fmt.Printf("Total elements: %d\n", len(calldataV0))
	fmt.Printf("  [0] Num calls:             %s\n", calldataV0[0].String())
	fmt.Printf("  [1] Contract address:      %s\n", calldataV0[1].String())
	fmt.Printf("  [2] Entry point selector:  %s\n", calldataV0[2].String())
	fmt.Printf("  [3] Calldata offset:       %s\n", calldataV0[3].String())
}