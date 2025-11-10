package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/NethermindEth/starknet.go/typeddata"
)

func main() {
	typedDataJSON := []byte(`{
		"types": {
			"StarkNetDomain": [
				{ "name": "name", "type": "felt" },
				{ "name": "version", "type": "felt" },
				{ "name": "chainId", "type": "felt" }
			],
			"Person": [
				{ "name": "name", "type": "felt" },
				{ "name": "age", "type": "felt" }
			]
		},
		"primaryType": "Person",
		"domain": {
			"name": "Example",
			"version": "1",
			"chainId": "1"
		},
		"message": {
			"name": "Bob",
			"age": "25"
		}
	}`)

	var td typeddata.TypedData
	err := json.Unmarshal(typedDataJSON, &td)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("GetStructHash:")
	
	// Get struct hash for the Person message
	structHash, err := td.GetStructHash("Person")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("  Person struct hash: %s\n", structHash.String())
	
	// Get domain hash
	domainHash, err := td.GetStructHash("StarkNetDomain")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("  Domain struct hash: %s\n", domainHash.String())
}
