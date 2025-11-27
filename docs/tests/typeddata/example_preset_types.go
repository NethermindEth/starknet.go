package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/NethermindEth/starknet.go/typeddata"
)

func main() {
	// Example using preset types: u256, TokenAmount, NftId
	typedDataJSON := []byte(`{
		"types": {
			"StarknetDomain": [
				{ "name": "name", "type": "shortstring" },
				{ "name": "version", "type": "shortstring" },
				{ "name": "chainId", "type": "shortstring" },
				{ "name": "revision", "type": "shortstring" }
			],
			"Transfer": [
				{ "name": "token", "type": "TokenAmount" },
				{ "name": "nft", "type": "NftId" }
			]
		},
		"primaryType": "Transfer",
		"domain": {
			"name": "NFT Marketplace",
			"version": "1",
			"chainId": "SN_MAIN",
			"revision": "1"
		},
		"message": {
			"token": {
				"token_address": "0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7",
				"amount": {
					"low": "1000",
					"high": "0"
				}
			},
			"nft": {
				"collection_address": "0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7",
				"token_id": {
					"low": "42",
					"high": "0"
				}
			}
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

	// Get type hashes for preset types
	u256Hash, _ := td.GetTypeHash("u256")
	tokenAmountHash, _ := td.GetTypeHash("TokenAmount")
	nftIdHash, _ := td.GetTypeHash("NftId")

	fmt.Println("Preset Types Example")
	fmt.Println("\nPreset Type Hashes:")
	fmt.Printf("  u256: %s\n", u256Hash.String())
	fmt.Printf("  TokenAmount: %s\n", tokenAmountHash.String())
	fmt.Printf("  NftId: %s\n", nftIdHash.String())
	fmt.Printf("\nMessage Hash: %s\n", messageHash.String())
	fmt.Println("\nPreset types: u256, TokenAmount, NftId are built-in")
}
