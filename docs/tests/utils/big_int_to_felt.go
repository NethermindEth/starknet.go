package main

import (
	"fmt"
	"math/big"

	"github.com/NethermindEth/starknet.go/utils"
)

func main() {
	bigIntValues := []*big.Int{
		big.NewInt(123),
		big.NewInt(0),
		new(big.Int).SetBytes([]byte{0x7e, 0x00, 0xd4, 0x96, 0xe3, 0x24, 0x87, 0x6b}),
	}

	fmt.Println("BigIntToFelt:")
	for _, bigInt := range bigIntValues {
		felt := utils.BigIntToFelt(bigInt)
		fmt.Printf("  Input BigInt: %s\n", bigInt.String())
		fmt.Printf("  Output Felt: %s\n\n", felt.String())
	}
}
