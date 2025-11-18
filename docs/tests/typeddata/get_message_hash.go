package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/NethermindEth/starknet.go/typeddata"
)

func main() {
	// Example typed data following SNIP-12
	typedDataJSON := []byte(`{
		"types": {
			"StarkNetDomain": [
				{ "name": "name", "type": "felt" },
				{ "name": "version", "type": "felt" },
				{ "name": "chainId", "type": "felt" }
			],
			"Person": [
				{ "name": "name", "type": "felt" },
				{ "name": "wallet", "type": "felt" },
				{ "name": "age", "type": "felt" }
			]
		},
		"primaryType": "Person",
		"domain": {
			"name": "MyDapp",
			"version": "1",
			"chainId": "SN_SEPOLIA"
		},
		"message": {
			"name": "Alice",
			"wallet": "0x1234567890abcdef",
			"age": "30"
		}
	}`)

	var td typeddata.TypedData
	err := json.Unmarshal(typedDataJSON, &td)
	if err != nil {
		log.Fatal(err)
	}

	// Calculate message hash for signing (using valid Starknet address)
	accountAddress := "0x7e00d496e324876bbc8531f2d9a82bf154d1a04a50218ee74cdd372f75a551a"
	messageHash, err := td.GetMessageHash(accountAddress)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("GetMessageHash:")
	fmt.Printf("  Account: %s\n", accountAddress)
	fmt.Printf("  Message Hash: %s\n", messageHash.String())
}
