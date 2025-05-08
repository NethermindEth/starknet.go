package main

import (
	"context"
	"fmt"
	"math"
	"math/big"

	"github.com/NethermindEth/juno/core/felt"
	setup "github.com/NethermindEth/starknet.go/examples/internal"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"
)

var (
	someContract               string = "0x04718f5a0fc34cc1af16a1cdee98ffb20c31f5cd61d6ab07201858f4287c938d" // Sepolia STRK contract address
	contractMethod             string = "decimals"
	contractMethodWithCalldata string = "balance_of"
)

// main entry point of the program.
//
// It initialises the environment and establishes a connection with the client.
// It then makes two contract calls and prints the responses.
func main() {
	fmt.Println("Starting simpleCall example")

	// Load variables from '.env' file
	rpcProviderUrl := setup.GetRpcProviderUrl()
	accountAddress := setup.GetAccountAddress()

	// Initialise connection to RPC provider
	client, err := rpc.NewProvider(rpcProviderUrl)
	if err != nil {
		panic(fmt.Sprintf("Error dialling the RPC provider: %s", err))
	}

	fmt.Println("Established connection with the client")

	contractAddress, err := utils.HexToFelt(someContract)
	if err != nil {
		fmt.Println("Failed to transform the token contract address, did you give the hex address?")
		panic(err)
	}

	// Here we are converting the account address to felt
	accountAddressInFelt, err := utils.HexToFelt(accountAddress)
	if err != nil {
		fmt.Println("Failed to transform the account address, did you give the hex address?")
		panic(err)
	}

	// Get token's decimals. As the contract method doesn't require any parameters, we can omit the calldata field.
	getDecimalsTx := rpc.FunctionCall{
		ContractAddress:    contractAddress,
		EntryPointSelector: utils.GetSelectorFromNameFelt(contractMethod),
	}
	decimalsResp, rpcErr := client.Call(context.Background(), getDecimalsTx, rpc.WithBlockTag("latest"))
	if rpcErr != nil {
		panic(rpcErr)
	}
	decimals, _ := utils.FeltToBigInt(decimalsResp[0]).Float64()
	fmt.Printf("Decimals: %v \n", decimals)

	// Get balance from specified account address. As the contract method requires a parameter, we need to pass it in the calldata field.
	tx := rpc.FunctionCall{
		ContractAddress:    contractAddress,
		EntryPointSelector: utils.GetSelectorFromNameFelt(contractMethodWithCalldata),
		Calldata:           []*felt.Felt{accountAddressInFelt},
	}
	balanceResp, rpcErr := client.Call(context.Background(), tx, rpc.WithBlockTag("latest"))
	if rpcErr != nil {
		panic(rpcErr)
	}
	balance := balanceResp[0].BigInt(new(big.Int))
	fmt.Printf("Balance: %v \n", balance)

	// Getting result
	balance = balance.Div(balance, big.NewInt(int64(math.Pow(10, decimals))))
	fmt.Printf("Token balance of %s is %v STRK \n", accountAddress, balance)
}
