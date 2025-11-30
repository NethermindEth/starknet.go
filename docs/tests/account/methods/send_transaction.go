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
 
	strkContract, _ := utils.HexToFelt("0x04718f5a0fc34cc1af16a1cdee98ffb20c31f5cd61d6ab07201858f4287c938d")
	amount := new(felt.Felt).SetUint64(1)
	u256Amount, _ := utils.HexToU256Felt(amount.String())
 
	fnCall := rpc.FunctionCall{
		ContractAddress:    strkContract,
		EntryPointSelector: utils.GetSelectorFromNameFelt("transfer"),
		Calldata:           append([]*felt.Felt{accountAddressFelt}, u256Amount...),
	}
 
	calldata, err := accnt.FmtCalldata([]rpc.FunctionCall{fnCall})
	if err != nil {
		log.Fatal("Failed to format calldata:", err)
	}
 
	invokeTx := utils.BuildInvokeTxn(
		accnt.Address,
		nonce,
		calldata,
		&rpc.ResourceBoundsMapping{
			L1Gas: rpc.ResourceBounds{
				MaxAmount:       "0x0",
				MaxPricePerUnit: "0x0",
			},
			L1DataGas: rpc.ResourceBounds{
				MaxAmount:       "0x0",
				MaxPricePerUnit: "0x0",
			},
			L2Gas: rpc.ResourceBounds{
				MaxAmount:       "0x0",
				MaxPricePerUnit: "0x0",
			},
		},
		nil,
	)
 
	err = accnt.SignInvokeTransaction(ctx, invokeTx)
	if err != nil {
		log.Fatal("Failed to sign for estimation:", err)
	}
 
	feeRes, err := accnt.Provider.EstimateFee(
		ctx,
		[]rpc.BroadcastTxn{invokeTx},
		[]rpc.SimulationFlag{},
		rpc.WithBlockTag("pre_confirmed"),
	)
	if err != nil {
		log.Fatal("Failed to estimate fee:", err)
	}
 
	invokeTx.ResourceBounds = utils.FeeEstToResBoundsMap(feeRes[0], 1.5)
 
	err = accnt.SignInvokeTransaction(ctx, invokeTx)
	if err != nil {
		log.Fatal("Failed to sign final transaction:", err)
	}
 
	response, err := accnt.SendTransaction(ctx, invokeTx)
	if err != nil {
		log.Fatal("Failed to send transaction:", err)
	}
 
	fmt.Printf("Transaction Hash: %s\n", response.Hash.String())
}