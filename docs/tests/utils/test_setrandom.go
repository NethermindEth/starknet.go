package main

import (
	"fmt"
	"github.com/NethermindEth/juno/core/felt"
)

func main() {
	f := new(felt.Felt)
	result := f.SetRandom()
	fmt.Printf("Type: %T, Value: %v\n", result, result)
}
