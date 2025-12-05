package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/NethermindEth/starknet.go/typeddata"
)

func main() {
	// Same message structure, different domains produce different hashes
	accountAddress := "0x7e00d496e324876bbc8531f2d9a82bf154d1a04a50218ee74cdd372f75a551a"

	// Domain 1: Mainnet
	mainnetData := []byte(`{
		"types": {
			"StarkNetDomain": [
				{ "name": "name", "type": "felt" },
				{ "name": "version", "type": "felt" },
				{ "name": "chainId", "type": "felt" }
			],
			"Action": [
				{ "name": "action", "type": "felt" },
				{ "name": "timestamp", "type": "felt" }
			]
		},
		"primaryType": "Action",
		"domain": {
			"name": "MyDapp",
			"version": "1",
			"chainId": "SN_MAIN"
		},
		"message": {
			"action": "approve",
			"timestamp": "1234567890"
		}
	}`)

	// Domain 2: Sepolia (testnet)
	sepoliaData := []byte(`{
		"types": {
			"StarkNetDomain": [
				{ "name": "name", "type": "felt" },
				{ "name": "version", "type": "felt" },
				{ "name": "chainId", "type": "felt" }
			],
			"Action": [
				{ "name": "action", "type": "felt" },
				{ "name": "timestamp", "type": "felt" }
			]
		},
		"primaryType": "Action",
		"domain": {
			"name": "MyDapp",
			"version": "1",
			"chainId": "SN_SEPOLIA"
		},
		"message": {
			"action": "approve",
			"timestamp": "1234567890"
		}
	}`)

	// Parse and hash mainnet message
	var mainnetTD typeddata.TypedData
	if err := json.Unmarshal(mainnetData, &mainnetTD); err != nil {
		log.Fatal(err)
	}
	mainnetHash, err := mainnetTD.GetMessageHash(accountAddress)
	if err != nil {
		log.Fatal(err)
	}

	// Parse and hash sepolia message
	var sepoliaTD typeddata.TypedData
	if err := json.Unmarshal(sepoliaData, &sepoliaTD); err != nil {
		log.Fatal(err)
	}
	sepoliaHash, err := sepoliaTD.GetMessageHash(accountAddress)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Domain Separation Example")
	fmt.Println("\nMainnet Domain:")
	fmt.Printf("  Chain ID: %s\n", mainnetTD.Domain.ChainID)
	fmt.Printf("  Message Hash: %s\n", mainnetHash.String())

	fmt.Println("\nSepolia Domain:")
	fmt.Printf("  Chain ID: %s\n", sepoliaTD.Domain.ChainID)
	fmt.Printf("  Message Hash: %s\n", sepoliaHash.String())

	fmt.Println("\nNote: Same message, different domains produce different hashes")
	fmt.Println("This prevents replay attacks across networks")
}
