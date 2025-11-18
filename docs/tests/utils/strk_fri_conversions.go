package main

import (
	"fmt"

	"github.com/NethermindEth/starknet.go/utils"
)

func main() {
	fmt.Println("STRK/FRI Conversions:")
	
	// STRK to FRI
	strkValues := []float64{1.0, 0.5, 10.0}
	
	fmt.Println("\nSTRKToFRI:")
	for _, strk := range strkValues {
		fri := utils.STRKToFRI(strk)
		fmt.Printf("  Input STRK: %v\n", strk)
		fmt.Printf("  Output FRI: %s\n\n", fri.String())
	}
	
	// FRI to STRK
	friValues := []string{
		"0xde0b6b3a7640000",    // 1 STRK
		"0x6f05b59d3b20000",    // 0.5 STRK
		"0x8ac7230489e80000",   // 10 STRK
	}
	
	fmt.Println("FRIToSTRK:")
	for _, friHex := range friValues {
		fri, _ := utils.HexToFelt(friHex)
		strk := utils.FRIToSTRK(fri)
		fmt.Printf("  Input FRI: %s\n", fri.String())
		fmt.Printf("  Output STRK: %v\n\n", strk)
	}
}
