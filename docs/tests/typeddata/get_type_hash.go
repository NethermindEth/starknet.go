package main
 
import (
	"encoding/json"
	"fmt"
	"log"
 
	"github.com/NethermindEth/starknet.go/typeddata"
)
 
func main() {
	// Create TypedData
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
			"from": {"name": "Alice", "wallet": "0x1234"},
			"to": {"name": "Bob", "wallet": "0x5678"},
			"contents": "Hello!"
		}
	}`)
 
	var td typeddata.TypedData
	if err := json.Unmarshal(typedDataJSON, &td); err != nil {
		log.Fatalf("Failed to unmarshal: %v", err)
	}
 
	// Get type hash for Mail
	mailTypeHash, err := td.GetTypeHash("Mail")
	if err != nil {
		log.Fatalf("Failed to get type hash: %v", err)
	}
	fmt.Printf("Mail type hash: %s\n", mailTypeHash.String())
 
	// Get type hash for Person
	personTypeHash, err := td.GetTypeHash("Person")
	if err != nil {
		log.Fatalf("Failed to get type hash: %v", err)
	}
	fmt.Printf("Person type hash: %s\n", personTypeHash.String())
 
	// Try to get hash for preset type (u128)
	u128TypeHash, err := td.GetTypeHash("u128")
	if err != nil {
		fmt.Printf("Note: u128 is a preset type: %v\n", err)
	} else {
		fmt.Printf("u128 type hash: %s\n", u128TypeHash.String())
	}
}