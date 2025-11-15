package main

import (
	"fmt"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"
)

func main() {
	// Define contract address and entry point
	contractAddress, _ := utils.HexToFelt("0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7")
	entryPointSelector := utils.GetSelectorFromNameFelt("transfer")

	// Create recipient address
	recipient, _ := new(felt.Felt).SetString("0x1234567890abcdef")

	// Create function calls
	functionCalls := []rpc.FunctionCall{
		{
			ContractAddress:    contractAddress,
			EntryPointSelector: entryPointSelector,
			Calldata: []*felt.Felt{
				recipient,                      // recipient
				new(felt.Felt).SetUint64(100), // amount low
				new(felt.Felt).SetUint64(0),   // amount high
			},
		},
	}

	// Format calldata for Cairo 0
	calldata := account.FmtCallDataCairo0(functionCalls)
	fmt.Printf("Formatted calldata length: %d\n", len(calldata))

	// Display the formatted calldata structure
	fmt.Println("\nCalldata structure:")
	fmt.Printf("  Number of calls: %s\n", calldata[0])
	if len(calldata) > 1 {
		fmt.Printf("  First call contract: %s\n", calldata[1])
	}
	if len(calldata) > 2 {
		fmt.Printf("  First call selector: %s\n", calldata[2])
	}

	// Display all calldata elements
	fmt.Println("\nFull calldata array:")
	for i, data := range calldata {
		fmt.Printf("  [%d]: %s\n", i, data)
	}
}