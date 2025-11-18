package main

import (
	"fmt"

	"github.com/NethermindEth/starknet.go/utils"
)

func main() {
	functionNames := []string{
		"transfer",
		"balanceOf",
		"approve",
		"increase_balance",
	}

	fmt.Println("GetSelectorFromNameFelt:")
	for _, name := range functionNames {
		selector := utils.GetSelectorFromNameFelt(name)
		fmt.Printf("  Function: %s\n", name)
		fmt.Printf("  Selector: %s\n\n", selector.String())
	}
}
