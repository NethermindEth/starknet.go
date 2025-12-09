package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/NethermindEth/starknet.go/typeddata"
)

func main() {
	// Define a simple mail message structure
	typedDataJSON := []byte(`{
		"types": {
			"StarkNetDomain": [
				{ "name": "name", "type": "felt" },
				{ "name": "version", "type": "felt" },
				{ "name": "chainId", "type": "felt" }
			],
			"Person": [
				{ "name": "name", "type": "felt" },
				{ "name": "wallet", "type": "felt" }
			],
			"Mail": [
				{ "name": "from", "type": "Person" },
				{ "name": "to", "type": "Person" },
				{ "name": "contents", "type": "felt" }
			]
		},
		"primaryType": "Mail",
		"domain": {
			"name": "StarkNet Mail",
			"version": "1",
			"chainId": 1
		},
		"message": {
			"from": {
				"name": "Alice",
				"wallet": "0x7e00d496e324876bbc8531f2d9a82bf154d1a04a50218ee74cdd372f75a551a"
			},
			"to": {
				"name": "Bob",
				"wallet": "0x69b49c2cc8b16e80e86bfc5b0614a59aa8c9b601569c7b80dde04d3f3151b79"
			},
			"contents": "Hello Bob!"
		}
	}`)

	// Parse the typed data
	var td typeddata.TypedData
	err := json.Unmarshal(typedDataJSON, &td)
	if err != nil {
		log.Fatal(err)
	}

	// Calculate the message hash for signing
	accountAddress := "0x7e00d496e324876bbc8531f2d9a82bf154d1a04a50218ee74cdd372f75a551a"
	messageHash, err := td.GetMessageHash(accountAddress)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Basic Mail Signing Example")
	fmt.Printf("From: Alice (%s)\n", td.Message["from"].(map[string]interface{})["wallet"])
	fmt.Printf("To: Bob (%s)\n", td.Message["to"].(map[string]interface{})["wallet"])
	fmt.Printf("Message Hash: %s\n", messageHash.String())
	fmt.Println("\nThis hash can now be signed with your private key")
}
