package main
 
import (
	"fmt"
	"log"
 
	"github.com/NethermindEth/starknet.go/typeddata"
)
 
func main() {
	// Define types (must include StarknetDomain)
	types := []typeddata.TypeDefinition{
		{
			Name: "StarknetDomain",
			Parameters: []typeddata.TypeParameter{
				{Name: "name", Type: "felt"},
				{Name: "version", Type: "felt"},
				{Name: "chainId", Type: "felt"},
			},
		},
		{
			Name: "Mail",
			Parameters: []typeddata.TypeParameter{
				{Name: "from", Type: "felt"},
				{Name: "to", Type: "felt"},
				{Name: "contents", Type: "felt"},
			},
		},
	}
 
	// Define domain
	domain := typeddata.Domain{
		Name:     "StarkNet Mail",
		Version:  "1",
		ChainID:  "SN_SEPOLIA",
		Revision: 1,
	}
 
	// Message data
	message := []byte(`{
		"from": "0x1234",
		"to": "0x5678",
		"contents": "Hello!"
	}`)
 
	// Create TypedData
	td, err := typeddata.NewTypedData(types, "Mail", domain, message)
	if err != nil {
		log.Fatalf("Failed to create TypedData: %v", err)
	}
 
	fmt.Printf("Primary Type: %s\n", td.PrimaryType)
	fmt.Printf("Domain: %s v%s\n", td.Domain.Name, td.Domain.Version)
}