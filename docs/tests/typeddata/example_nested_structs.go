package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/NethermindEth/starknet.go/typeddata"
)

func main() {
	// Complex nested structure example
	typedDataJSON := []byte(`{
		"types": {
			"StarknetDomain": [
				{ "name": "name", "type": "string" },
				{ "name": "chainId", "type": "felt" },
				{ "name": "version", "type": "string" }
			],
			"Address": [
				{ "name": "street", "type": "string" },
				{ "name": "city", "type": "string" },
				{ "name": "country", "type": "string" }
			],
			"Person": [
				{ "name": "name", "type": "string" },
				{ "name": "age", "type": "felt" },
				{ "name": "address", "type": "Address" }
			],
			"Transfer": [
				{ "name": "from", "type": "Person" },
				{ "name": "to", "type": "Person" },
				{ "name": "amount", "type": "felt" }
			]
		},
		"primaryType": "Transfer",
		"domain": {
			"name": "Payment System",
			"chainId": "0x534e5f5345504f4c4941",
			"version": "1.0.0",
			"revision": "1"
		},
		"message": {
			"from": {
				"name": "Alice",
				"age": "30",
				"address": {
					"street": "123 Main St",
					"city": "New York",
					"country": "USA"
				}
			},
			"to": {
				"name": "Bob",
				"age": "25",
				"address": {
					"street": "456 Oak Ave",
					"city": "Los Angeles",
					"country": "USA"
				}
			},
			"amount": "1000"
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

	// Get struct hash for nested Address type
	addressHash, err := td.GetStructHash("Address", "from", "address")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Nested Structs Example")
	fmt.Printf("Primary Type: %s\n", td.PrimaryType)
	fmt.Printf("Number of Custom Types: %d\n", len(td.Types))
	fmt.Printf("\nNested Address Hash: %s\n", addressHash.String())
	fmt.Printf("Full Message Hash: %s\n", messageHash.String())
	fmt.Println("\nNested structures allow complex data modeling")
}
