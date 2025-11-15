package main

import (
	"fmt"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/rpc"
)

func main() {
	// Create contract address manually to avoid utils import
	contractAddress, _ := new(felt.Felt).SetString("0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7")

	// Create entry point selector manually (this is the selector for "transfer")
	entryPointSelector, _ := new(felt.Felt).SetString("0x83afd3f4caedc6eebf44246fe54e38c95e3179a5ec9ea81740eca5b482d12e")

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

	// Format calldata for Cairo 2
	calldata := account.FmtCallDataCairo2(functionCalls)
	fmt.Printf("Formatted calldata length: %d\n", len(calldata))

	// Display the formatted calldata structure
	fmt.Println("\nCalldata structure (Cairo 2):")
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