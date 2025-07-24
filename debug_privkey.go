package main

import (
	"fmt"
	"math/big"
)

func main() {
	// Test the original key
	privKey1 := "0xbfee5bcea8219ac607547e60fb26210472d97d793dd17b0d8e4f16429ff38b"
	fmt.Printf("Original key: %s\n", privKey1)
	fmt.Printf("Length after 0x: %d\n", len(privKey1)-2)
	
	privKeyBI1, ok1 := new(big.Int).SetString(privKey1, 0)
	fmt.Printf("SetString result: %t\n", ok1)
	if ok1 {
		fmt.Printf("BigInt value: %s\n", privKeyBI1.String())
	}
	
	// Test with leading zeros (64 chars total)
	privKey2 := "0x00bfee5bcea8219ac607547e60fb26210472d97d793dd17b0d8e4f16429ff38b"
	fmt.Printf("\nPadded key: %s\n", privKey2)
	fmt.Printf("Length after 0x: %d\n", len(privKey2)-2)
	
	privKeyBI2, ok2 := new(big.Int).SetString(privKey2, 0)
	fmt.Printf("SetString result: %t\n", ok2)
	if ok2 {
		fmt.Printf("BigInt value: %s\n", privKeyBI2.String())
	}
	
	// Check for hidden characters
	fmt.Printf("\nOriginal key bytes: %v\n", []byte(privKey1))
}
