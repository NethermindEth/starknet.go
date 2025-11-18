package main

import (
	"fmt"
	"github.com/NethermindEth/juno/core/felt"
)

func main() {
	// Test if we can use felt without utils
	f := new(felt.Felt).SetUint64(123)
	fmt.Printf("Felt: %s\n", f.String())
}
