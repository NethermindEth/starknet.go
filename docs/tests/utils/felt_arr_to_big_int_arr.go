package main

import (
	"fmt"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/utils"
)

func main() {
	fmt.Println("FeltArrToBigIntArr:")

	// Create an array of Felt values
	feltArr := []*felt.Felt{
		new(felt.Felt).SetUint64(123),
		new(felt.Felt).SetUint64(456),
		new(felt.Felt).SetUint64(789),
	}

	// Convert to big.Int array
	bigIntArr := utils.FeltArrToBigIntArr(feltArr)

	fmt.Print("  Input Felt array: [")
	for i, f := range feltArr {
		if i > 0 {
			fmt.Print(", ")
		}
		fmt.Print(f.String())
	}
	fmt.Println("]")

	fmt.Print("  Output big.Int array: [")
	for i, b := range bigIntArr {
		if i > 0 {
			fmt.Print(", ")
		}
		fmt.Print(b.String())
	}
	fmt.Println("]")
}
