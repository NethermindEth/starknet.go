package main

import (
	"fmt"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"
)

func main() {
	// Build deploy account transaction parameters
	nonce := new(felt.Felt).SetUint64(0) // First transaction from account
	contractAddressSalt := new(felt.Felt).SetUint64(12345)
	publicKey := new(felt.Felt).SetUint64(987654321)
	constructorCalldata := []*felt.Felt{publicKey}
	classHash := new(felt.Felt).SetUint64(123456789)
	
	// Resource bounds
	resourceBounds := &rpc.ResourceBoundsMapping{
		L1Gas: rpc.ResourceBounds{
			MaxAmount:       rpc.U64(10000),
			MaxPricePerUnit: rpc.U128(100),
		},
	}
	
	// Build transaction
	txn := utils.BuildDeployAccountTxn(
		nonce,
		contractAddressSalt,
		constructorCalldata,
		classHash,
		resourceBounds,
		nil,
	)
	
	fmt.Printf("Built deploy account transaction\n")
	fmt.Printf("  Class Hash: %s\n", txn.ClassHash.String())
	fmt.Printf("  Contract Address Salt: %s\n", txn.ContractAddressSalt.String())
	fmt.Println("✅ Transaction built successfully")
}
