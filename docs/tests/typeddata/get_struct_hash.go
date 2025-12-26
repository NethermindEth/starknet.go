package main
 
import (
	"encoding/json"
	"fmt"
	"log"
 
	"github.com/NethermindEth/starknet.go/typeddata"
)
 
func main() {
	// Create TypedData with nested types
	typedDataJSON := []byte(`{
		"types": {
			"StarknetDomain": [
				{"name": "name", "type": "shortstring"},
				{"name": "version", "type": "shortstring"},
				{"name": "chainId", "type": "shortstring"},
				{"name": "revision", "type": "shortstring"}
			],
			"Person": [
				{"name": "name", "type": "shortstring"},
				{"name": "wallet", "type": "felt"}
			],
			"Mail": [
				{"name": "from", "type": "Person"},
				{"name": "to", "type": "Person"},
				{"name": "contents", "type": "shortstring"}
			]
		},
		"primaryType": "Mail",
		"domain": {
			"name": "StarkNet Mail",
			"version": "1",
			"chainId": "SN_SEPOLIA",
			"revision": "1"
		},
		"message": {
			"from": {
				"name": "Alice",
				"wallet": "0x1234"
			},
			"to": {
				"name": "Bob",
				"wallet": "0x5678"
			},
			"contents": "Hello!"
		}
	}`)
 
	var td typeddata.TypedData
	if err := json.Unmarshal(typedDataJSON, &td); err != nil {
		log.Fatalf("Failed to unmarshal: %v", err)
	}
 
	// Get hash of the primary type
	primaryHash, err := td.GetStructHash("Mail")
	if err != nil {
		log.Fatalf("Failed to get struct hash: %v", err)
	}
	fmt.Printf("Mail struct hash: %s\n", primaryHash.String())
 
	// Get hash of nested Person type for "from" field
	fromHash, err := td.GetStructHash("Person", "from")
	if err != nil {
		log.Fatalf("Failed to get nested hash: %v", err)
	}
	fmt.Printf("From Person hash: %s\n", fromHash.String())
 
	// Get hash of the domain
	domainHash, err := td.GetStructHash(td.Revision.Domain())
	if err != nil {
		log.Fatalf("Failed to get domain hash: %v", err)
	}
	fmt.Printf("Domain hash: %s\n", domainHash.String())
}