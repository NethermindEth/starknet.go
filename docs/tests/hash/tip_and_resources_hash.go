package main

import (
	"fmt"

	"github.com/NethermindEth/starknet.go/hash"
	"github.com/NethermindEth/starknet.go/rpc"
)

func main() {
	tip := uint64(100)
	resourceBounds := &rpc.ResourceBoundsMapping{
		L1Gas: rpc.ResourceBounds{
			MaxAmount:       rpc.U64("0x186a0"),
			MaxPricePerUnit: rpc.U128("0x3e8"),
		},
		L1DataGas: rpc.ResourceBounds{
			MaxAmount:       rpc.U64("0x0"),
			MaxPricePerUnit: rpc.U128("0x0"),
		},
		L2Gas: rpc.ResourceBounds{
			MaxAmount:       rpc.U64("0xc350"),
			MaxPricePerUnit: rpc.U128("0x1f4"),
		},
	}

	hashResult, err := hash.TipAndResourcesHash(tip, resourceBounds)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("TipAndResourcesHash: %s\n", hashResult.String())
}
