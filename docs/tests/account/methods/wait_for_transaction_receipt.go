package main
 
import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"
 
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
 
	strkContract, _ := utils.HexToFelt("0x04718f5a0fc34cc1af16a1cdee98ffb20c31f5cd61d6ab07201858f4287c938d")
	recipient := accountAddressFelt
	amount := new(felt.Felt).SetUint64(1000000000000000)
	u256Amount, _ := utils.HexToU256Felt(amount.String())
 
	invokeFnCall := rpc.InvokeFunctionCall{
		ContractAddress: strkContract,
		FunctionName:    "transfer",
		CallData:        append([]*felt.Felt{recipient}, u256Amount...),
	}
 
	tx, err := accnt.BuildAndSendInvokeTxn(ctx, []rpc.InvokeFunctionCall{invokeFnCall}, nil)
	if err != nil {
		log.Fatal("Failed to send transaction:", err)
	}
 
	fmt.Printf("Transaction sent: %s\n", tx.Hash.String())
 
	receipt, err := accnt.WaitForTransactionReceipt(ctx, tx.Hash, 3*time.Second)
	if err != nil {
		log.Fatal("Failed to get receipt:", err)
	}
 
	fmt.Printf("Block Hash: %s\n", receipt.BlockHash.String())
	fmt.Printf("Block Number: %d\n", receipt.BlockNumber)
	fmt.Printf("Status: %s\n", receipt.FinalityStatus)
}