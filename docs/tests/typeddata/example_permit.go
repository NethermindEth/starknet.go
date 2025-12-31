package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/NethermindEth/starknet.go/typeddata"
)

func main() {
	// ERC20 Permit example - gasless token approval via signature
	typedDataJSON := []byte(`{
		"types": {
			"StarknetDomain": [
				{ "name": "name", "type": "string" },
				{ "name": "version", "type": "string" },
				{ "name": "chainId", "type": "felt" },
				{ "name": "revision", "type": "shortstring" }
			],
			"Permit": [
				{ "name": "owner", "type": "ContractAddress" },
				{ "name": "spender", "type": "ContractAddress" },
				{ "name": "value", "type": "u256" },
				{ "name": "nonce", "type": "felt" },
				{ "name": "deadline", "type": "felt" }
			]
		},
		"primaryType": "Permit",
		"domain": {
			"name": "MyToken",
			"version": "1",
			"chainId": "0x534e5f4d41494e",
			"revision": "1"
		},
		"message": {
			"owner": "0x7e00d496e324876bbc8531f2d9a82bf154d1a04a50218ee74cdd372f75a551a",
			"spender": "0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7",
			"value": {
				"low": "1000000000000000000",
				"high": "0"
			},
			"nonce": "0",
			"deadline": "1735689600"
		}
	}`)

	var td typeddata.TypedData
	err := json.Unmarshal(typedDataJSON, &td)
	if err != nil {
		log.Fatal(err)
	}

	// Owner account address
	ownerAddress := "0x7e00d496e324876bbc8531f2d9a82bf154d1a04a50218ee74cdd372f75a551a"

	// Calculate message hash for signing
	messageHash, err := td.GetMessageHash(ownerAddress)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("ERC20 Permit Example (Gasless Approval)")
	fmt.Printf("Token: %s\n", td.Domain.Name)
	fmt.Printf("Owner: %s\n", td.Message["owner"])
	fmt.Printf("Spender: %s\n", td.Message["spender"])

	value := td.Message["value"].(map[string]interface{})
	fmt.Printf("Amount: %s (low: %s, high: %s)\n", "1 ETH", value["low"], value["high"])
	fmt.Printf("Nonce: %s\n", td.Message["nonce"])
	fmt.Printf("Deadline: %s\n", td.Message["deadline"])

	fmt.Printf("\nMessage Hash: %s\n", messageHash.String())

	fmt.Println("\nPermit workflow:")
	fmt.Println("1. User signs permit message (no gas cost)")
	fmt.Println("2. Spender submits permit signature on-chain")
	fmt.Println("3. Contract verifies signature and grants approval")
	fmt.Println("4. Spender can now spend tokens on behalf of owner")
}
