package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"math/big"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"
	"github.com/joho/godotenv"
)

var (
	name               string = "mainnet"
	ethMainnetContract string = "0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"
	// usdcMainnetContract string="0x053c91253bc9682c04929ca02ed00b3e423f6710d2ee7e0d5ebb06f3ecf368a8"
	// daiMainnetContract string="0x00da114221cb83fa859dbdb4c44beeaa0bb37c7537ad5ae66fe5e0efd20e6eb3"
	getBalanceMethod  string = "balanceOf"
	getDecimalsMethod string = "decimals"
)

// main entry point of the program.
//
// It initializes the environment and establishes a connection with the client.
// It then calls a token contract to get the account balance and decimals.
// In the end, it transfomers the response to readable number.
//
// Parameters:
//
//	none
//
// Returns:
//
//	none
func main() {
	fmt.Println("Starting getTokenBalance example")
	godotenv.Load(fmt.Sprintf(".env.%s", name))
	url := os.Getenv("INTEGRATION_BASE")

	clientv02, err := rpc.NewProvider(url)
	if err != nil {
		log.Fatal(fmt.Sprintf("Error dialing the RPC provider: %s", err))
	}

	fmt.Println("Established connection with the client")

	tokenAddressInFelt, err := utils.HexToFelt(ethMainnetContract)
	if err != nil {
		fmt.Println("Failed to transform the token contract address, did you give the hex address?")
		panic(err)
	}

	accountAddress := os.Getenv("ACCOUNT_ADDR")
	accountAddressInFelt, err := utils.HexToFelt(accountAddress)
	if err != nil {
		fmt.Println("Failed to transform the account address, did you give the hex address?")
		panic(err)
	}

	// Make read contract call
	tx := rpc.FunctionCall{
		ContractAddress:    tokenAddressInFelt,
		EntryPointSelector: utils.GetSelectorFromNameFelt(getBalanceMethod),
		Calldata:           []*felt.Felt{accountAddressInFelt},
	}

	fmt.Println("Making balanceOf() request")
	callResp, rpcErr := clientv02.Call(context.Background(), tx, rpc.BlockID{Tag: "latest"})
	if rpcErr != nil {
		panic(rpcErr)
	}

	// Get token's decimals
	getDecimalsTx := rpc.FunctionCall{
		ContractAddress:    tokenAddressInFelt,
		EntryPointSelector: utils.GetSelectorFromNameFelt(getDecimalsMethod),
	}
	getDecimalsResp, rpcErr := clientv02.Call(context.Background(), getDecimalsTx, rpc.BlockID{Tag: "latest"})
	if rpcErr != nil {
		panic(rpcErr)
	}

	floatValue := new(big.Float).SetInt(utils.FeltToBigInt(callResp[0]))
	floatValue.Quo(floatValue, new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), utils.FeltToBigInt(getDecimalsResp[0]), nil)))

	fmt.Printf("Token balance of %s is %f", accountAddress, floatValue)
}
