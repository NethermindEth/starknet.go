package main

import (
	"fmt"

	"github.com/NethermindEth/starknet.go/utils"
)

func main() {
	// Convert 1 STRK to FRI
	fri := utils.STRKToFRI(1.0)
	fmt.Printf("1 STRK = %s FRI\n", fri.String())
	// Output: 0xde0b6b3a7640000

	// Convert 0.5 STRK to FRI
	fri2 := utils.STRKToFRI(0.5)
	fmt.Printf("0.5 STRK = %s FRI\n", fri2.String())
	// Output: 0x6f05b59d3b20000

	// Convert 10 STRK to FRI
	fri3 := utils.STRKToFRI(10.0)
	fmt.Printf("10 STRK = %s FRI\n", fri3.String())
	// Output: 0x8ac7230489e80000
}
