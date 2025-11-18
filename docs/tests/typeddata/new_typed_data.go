package main

import (
	"fmt"
	"log"

	"github.com/NethermindEth/starknet.go/typeddata"
)

func main() {
	// Define types
	types := []typeddata.TypeDefinition{
		{
			Parameters: []typeddata.TypeParameter{
				{Name: "name", Type: "felt"},
				{Name: "version", Type: "felt"},
				{Name: "chainId", Type: "felt"},
			},
		},
		{
			Parameters: []typeddata.TypeParameter{
				{Name: "username", Type: "felt"},
				{Name: "level", Type: "felt"},
			},
		},
	}
	types[0].Name = "StarkNetDomain"
	types[1].Name = "User"

	// Define domain
	domain := typeddata.Domain{
		Name:    "GameApp",
		Version: "1",
		ChainID: "1",
	}

	// Define message
	messageJSON := []byte(`{
		"username": "player1",
		"level": "10"
	}`)

	// Create TypedData
	td, err := typeddata.NewTypedData(types, "User", domain, messageJSON)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("NewTypedData:")
	fmt.Printf("  Primary Type: %s\n", td.PrimaryType)
	fmt.Printf("  Domain Name: %s\n", td.Domain.Name)
	fmt.Printf("  Domain ChainID: %s\n", td.Domain.ChainID)
	fmt.Printf("  Number of Types: %d\n", len(td.Types))
	
	// Get message hash
	accountAddr := "0x123"
	hash, err := td.GetMessageHash(accountAddr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("  Message Hash: %s\n", hash.String())
}
