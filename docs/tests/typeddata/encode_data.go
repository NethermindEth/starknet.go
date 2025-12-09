package main
 
import (
	"encoding/json"
	"fmt"
	"log"
 
	"github.com/NethermindEth/starknet.go/typeddata"
)
 
func main() {
	// Create TypedData from JSON (includes revision field)
	typedDataJSON := []byte(`{
		"types": {
			"StarknetDomain": [
				{"name": "name", "type": "shortstring"},
				{"name": "version", "type": "shortstring"},
				{"name": "chainId", "type": "shortstring"},
				{"name": "revision", "type": "shortstring"}
			],
			"Mail": [
				{"name": "from", "type": "felt"},
				{"name": "to", "type": "felt"},
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
			"from": "0x1234",
			"to": "0x5678",
			"contents": "Hello!"
		}
	}`)
 
	var td typeddata.TypedData
	err := json.Unmarshal(typedDataJSON, &td)
	if err != nil {
		log.Fatalf("Failed to unmarshal: %v", err)
	}
 
	// Get type definition for Mail
	typeDef := td.Types["Mail"]
 
	// Encode the message data
	encoded, err := typeddata.EncodeData(&typeDef, &td)
	if err != nil {
		log.Fatalf("Encoding failed: %v", err)
	}
 
	fmt.Printf("Encoded Mail message: %d felts\n", len(encoded))
	for i, felt := range encoded {
		fmt.Printf("  [%d]: %s\n", i, felt.String())
	}
}