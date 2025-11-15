package main

import (
	"fmt"

	"github.com/NethermindEth/starknet.go/utils"
)

func main() {
	// Generate selectors for common function names
	functionNames := []string{
		"transfer",
		"balanceOf",
		"approve",
		"transferFrom",
		"get_balance",
	}

	fmt.Println("GetSelectorFromNameFelt:")
	for _, name := range functionNames {
		selector := utils.GetSelectorFromNameFelt(name)
		fmt.Printf("  Function: %s\n", name)
		fmt.Printf("  Selector: %s\n", selector.String())
		fmt.Println()
	}
}
