package main

import (
	"fmt"

	"github.com/NethermindEth/starknet.go/hash"
	"github.com/NethermindEth/starknet.go/rpc"
)

func main() {
	feeDAMode := rpc.DAModeL1
	nonceDAMode := rpc.DAModeL1

	result, err := hash.DataAvailabilityModeConc(feeDAMode, nonceDAMode)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("DataAvailabilityModeConc (L1, L1): %d\n", result)

	// Try with L2 modes
	feeDAMode = rpc.DAModeL2
	nonceDAMode = rpc.DAModeL2

	result, err = hash.DataAvailabilityModeConc(feeDAMode, nonceDAMode)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("DataAvailabilityModeConc (L2, L2): %d\n", result)
}
