package main
 
import (
	"encoding/json"
	"fmt"
	"log"
 
	"github.com/NethermindEth/starknet.go/typeddata"
)
 
func main() {
	// Create TypedData from JSON
	typedDataJSON := []byte(`{
		"types": {
			"StarkNetDomain": [
				{"name": "name", "type": "felt"},
				{"name": "version", "type": "felt"},
				{"name": "chainId", "type": "felt"}
			],
			"Mail": [
				{"name": "from", "type": "felt"},
				{"name": "to", "type": "felt"},
				{"name": "contents", "type": "felt"}
			]
		},
		"primaryType": "Mail",
		"domain": {
			"name": "StarkNet Mail",
			"version": "1",
			"chainId": "SN_SEPOLIA"
		},
		"message": {
			"from": "0x1234",
			"to": "0x5678",
			"contents": "Hello!"
		}
	}`)
 
	var td typeddata.TypedData
	if err := json.Unmarshal(typedDataJSON, &td); err != nil {
		log.Fatalf("Failed to unmarshal: %v", err)
	}
 
	// Account address that will sign
	accountAddress := "0x02cdAb749380950e7a7c0deFf5ea8eDD716fEb3a2952aDd4E5659655077B8510"
 
	// Get message hash for signing
	messageHash, err := td.GetMessageHash(accountAddress)
	if err != nil {
		log.Fatalf("Failed to get message hash: %v", err)
	}
 
	fmt.Printf("Message Hash: %s\n", messageHash.String())
	fmt.Printf("This hash should be signed with the account's private key\n")
}