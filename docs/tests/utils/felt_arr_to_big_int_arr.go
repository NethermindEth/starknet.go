package main

import (
	"fmt"
	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/utils"
)

func main() {
	// Create an array of Felt values
	feltArr := []*felt.Felt{
		new(felt.Felt).SetUint64(100),
		new(felt.Felt).SetUint64(200),
		new(felt.Felt).SetUint64(300),
	}

	// Convert to big.Int array
	bigIntArr := utils.FeltArrToBigIntArr(feltArr)

	fmt.Println("Felt Array:")
	for i, f := range feltArr {
		fmt.Printf("  [%d]: %s\n", i, f.String())
	}

	fmt.Println("\nBig Int Array:")
	for i, b := range bigIntArr {
		fmt.Printf("  [%d]: %s\n", i, b.String())
	}
	// Output:
	// Felt Array:
	//   [0]: 100
	//   [1]: 200
	//   [2]: 300
	//
	// Big Int Array:
	//   [0]: 100
	//   [1]: 200
	//   [2]: 300
}
