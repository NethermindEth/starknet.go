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
			"name": "Example",
			"version": "1",
			"chainId": "1"
		},
		"message": {}
	}`)

	var td typeddata.TypedData
	err := json.Unmarshal(typedDataJSON, &td)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("GetTypeHash:")
	
	// Get type hash for Person
	personHash, err := td.GetTypeHash("Person")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("  Person type hash: %s\n", personHash.String())

	// Get type hash for Mail
	mailHash, err := td.GetTypeHash("Mail")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("  Mail type hash: %s\n", mailHash.String())
}
