package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/NethermindEth/starknet.go/typeddata"
)

func main() {
	// Example with array fields
	typedDataJSON := []byte(`{
		"types": {
			"StarknetDomain": [
				{ "name": "name", "type": "shortstring" },
				{ "name": "version", "type": "shortstring" },
				{ "name": "chainId", "type": "shortstring" },
				{ "name": "revision", "type": "shortstring" }
			],
			"Batch": [
				{ "name": "description", "type": "string" },
				{ "name": "recipients", "type": "ContractAddress*" },
				{ "name": "amounts", "type": "u128*" }
			]
		},
		"primaryType": "Batch",
		"domain": {
			"name": "BatchTransfer",
			"version": "1",
			"chainId": "SN_SEPOLIA",
			"revision": "1"
		},
		"message": {
			"description": "Monthly payments",
			"recipients": [
				"0x7e00d496e324876bbc8531f2d9a82bf154d1a04a50218ee74cdd372f75a551a",
				"0x69b49c2cc8b16e80e86bfc5b0614a59aa8c9b601569c7b80dde04d3f3151b79",
				"0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"
			],
			"amounts": [100, 200, 300]
		}
	}`)

	var td typeddata.TypedData
	err := json.Unmarshal(typedDataJSON, &td)
	if err != nil {
		log.Fatal(err)
	}

	// Calculate message hash
	accountAddress := "0x7e00d496e324876bbc8531f2d9a82bf154d1a04a50218ee74cdd372f75a551a"
	messageHash, err := td.GetMessageHash(accountAddress)
	if err != nil {
		log.Fatal(err)
	}

	recipients := td.Message["recipients"].([]interface{})
	amounts := td.Message["amounts"].([]interface{})

	fmt.Println("Array Types Example")
	fmt.Printf("Description: %s\n", td.Message["description"])
	fmt.Printf("Number of recipients: %d\n", len(recipients))
	fmt.Printf("Number of amounts: %d\n", len(amounts))
	fmt.Printf("\nMessage Hash: %s\n", messageHash.String())
	fmt.Println("\nArrays are denoted with * suffix (e.g., u128*, ContractAddress*)")
}
