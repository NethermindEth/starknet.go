package main

import (
	"context"
	"fmt"
	"os"

	"math/big"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
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
	base := os.Getenv("INTEGRATION_BASE")
	c, err := ethrpc.DialContext(context.Background(), base)
	if err != nil {
		fmt.Println("Failed to connect to the client, did you specify the url in the .env.mainnet?")
		panic(err)
	}
	clientv02 := rpc.NewProvider(c)
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
	callResp, err := clientv02.Call(context.Background(), tx, rpc.BlockID{Tag: "latest"})
	if err != nil {
		panic(err.Error())
	}

	// Get token's decimals
	getDecimalsTx := rpc.FunctionCall{
		ContractAddress:    tokenAddressInFelt,
		EntryPointSelector: utils.GetSelectorFromNameFelt(getDecimalsMethod),
	}
	getDecimalsResp, err := clientv02.Call(context.Background(), getDecimalsTx, rpc.BlockID{Tag: "latest"})
	if err != nil {
		panic(err)
	}

	floatValue := new(big.Float).SetInt(utils.FeltToBigInt(callResp[0]))
	floatValue.Quo(floatValue, new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), utils.FeltToBigInt(getDecimalsResp[0]), nil)))

	fmt.Printf("Token balance of %s is %f", accountAddress, floatValue)
}
