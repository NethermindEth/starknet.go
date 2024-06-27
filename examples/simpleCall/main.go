package main

import (
	"context"
	"fmt"
	"math"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"

	setup "github.com/NethermindEth/starknet.go/examples/internal"
)

var (
	someContract               string = "0x049D36570D4e46f48e99674bd3fcc84644DdD6b96F7C741B1562B82f9e004dC7" // Sepolia ETH contract address
	contractMethod             string = "decimals"
	contractMethodWithCalldata string = "balance_of"
)

// main entry point of the program.
//
// It initializes the environment and establishes a connection with the client.
// It then makes two contract calls and prints the responses.
//
// Parameters:
//
//	none
//
// Returns:
//
//	none
func main() {
	fmt.Println("Starting simpleCall example")

	// Load variables from '.env' file
	rpcProviderUrl := setup.GetRpcProviderUrl()
	accountAddress := setup.GetAccountAddress()

	// Initialize connection to RPC provider
	client, err := rpc.NewProvider(rpcProviderUrl)
	if err != nil {
		panic(fmt.Sprintf("Error dialing the RPC provider: %s", err))
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

	// Get token's decimals. Make read contract call without calldata
	getDecimalsTx := rpc.FunctionCall{
		ContractAddress:    contractAddress,
		EntryPointSelector: utils.GetSelectorFromNameFelt(contractMethod),
	}
	decimalsResp, rpcErr := client.Call(context.Background(), getDecimalsTx, rpc.BlockID{Tag: "latest"})
	if rpcErr != nil {
		panic(rpcErr)
	}
	decimals, _ := utils.FeltToBigInt(decimalsResp[0]).Float64()
	fmt.Printf("Decimals: %v \n", decimals)

	// Get balance from specified account address. Make read contract call with calldata
	tx := rpc.FunctionCall{
		ContractAddress:    contractAddress,
		EntryPointSelector: utils.GetSelectorFromNameFelt(contractMethodWithCalldata),
		Calldata:           []*felt.Felt{accountAddressInFelt},
	}
	balanceResp, rpcErr := client.Call(context.Background(), tx, rpc.BlockID{Tag: "latest"})
	if rpcErr != nil {
		panic(rpcErr)
	}
	balance, _ := utils.FeltToBigInt(balanceResp[0]).Float64()
	fmt.Printf("Balance: %d \n", int(balance))

	// Getting result
	balance = balance / (math.Pow(10, decimals))
	fmt.Printf("Token balance of %s is %f ETH \n", accountAddress, balance)
}
