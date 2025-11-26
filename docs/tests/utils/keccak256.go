package main

import (
	"fmt"

	"github.com/NethermindEth/starknet.go/utils"
)

func main() {
	// Hash simple string
	hash := utils.Keccak256([]byte("hello"))
	fmt.Printf("keccak256('hello') = %x\n", hash)
	// Output: 1c8aff950685c2ed4bc3174f3472287b56d9517b9c948127319a09a7a36deac8

	// Hash function name (used for selectors)
	hash2 := utils.Keccak256([]byte("transfer"))
	fmt.Printf("keccak256('transfer') = %x\n", hash2)
	// Output: b483afd3f4caedc6eebf44246fe54e38c95e3179a5ec9ea81740eca5b482d12e
}
