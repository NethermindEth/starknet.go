package main

import (
	"fmt"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/contracts"
)

func main() {
	// Example parameters for contract deployment
	deployerAddress, _ := new(felt.Felt).SetString("0x1234567890abcdef")
	salt, _ := new(felt.Felt).SetString("0x12345")
	classHash, _ := new(felt.Felt).SetString("0x07e2e0ba00c247e6c2c7e38e8cadfcc59f828bb94c182e69bd8ea667bcbb65e7")

	// Constructor calldata
	constructorCalldata := []*felt.Felt{
		new(felt.Felt).SetUint64(100),
		new(felt.Felt).SetUint64(200),
	}

	// Compute the contract address
	contractAddress := contracts.PrecomputeAddress(
		deployerAddress,
		salt,
		classHash,
		constructorCalldata,
	)

	fmt.Println("PrecomputeAddress:")
	fmt.Printf("  Deployer Address: %s\n", deployerAddress.String())
	fmt.Printf("  Salt: %s\n", salt.String())
	fmt.Printf("  Class Hash: %s\n", classHash.String())
	fmt.Printf("  Constructor Calldata: [%s, %s]\n",
		constructorCalldata[0].String(),
		constructorCalldata[1].String())
	fmt.Printf("  Computed Contract Address: %s\n", contractAddress.String())
}
