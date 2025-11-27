package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/NethermindEth/starknet.go/typeddata"
)

func main() {
	// Authentication flow example - login message signing
	timestamp := time.Now().Unix()

	typedDataJSON := []byte(fmt.Sprintf(`{
		"types": {
			"StarknetDomain": [
				{ "name": "name", "type": "shortstring" },
				{ "name": "version", "type": "shortstring" },
				{ "name": "chainId", "type": "shortstring" },
				{ "name": "revision", "type": "shortstring" }
			],
			"Login": [
				{ "name": "action", "type": "string" },
				{ "name": "timestamp", "type": "timestamp" },
				{ "name": "nonce", "type": "felt" }
			]
		},
		"primaryType": "Login",
		"domain": {
			"name": "MyDapp Auth",
			"version": "1.0",
			"chainId": "SN_MAIN",
			"revision": "1"
		},
		"message": {
			"action": "Sign in to MyDapp",
			"timestamp": "%d",
			"nonce": "123456789"
		}
	}`, timestamp))

	var td typeddata.TypedData
	err := json.Unmarshal(typedDataJSON, &td)
	if err != nil {
		log.Fatal(err)
	}

	// User's account address
	accountAddress := "0x7e00d496e324876bbc8531f2d9a82bf154d1a04a50218ee74cdd372f75a551a"

	// Calculate message hash for signing
	messageHash, err := td.GetMessageHash(accountAddress)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Authentication Flow Example")
	fmt.Printf("Action: %s\n", td.Message["action"])
	fmt.Printf("Timestamp: %s\n", td.Message["timestamp"])
	fmt.Printf("Nonce: %s\n", td.Message["nonce"])
	fmt.Printf("Account: %s\n", accountAddress)
	fmt.Printf("\nMessage Hash: %s\n", messageHash.String())
	fmt.Println("\nAuthentication flow:")
	fmt.Println("1. User requests to sign in")
	fmt.Println("2. Frontend creates typed data with current timestamp and nonce")
	fmt.Println("3. User signs the message hash")
	fmt.Println("4. Backend verifies signature and establishes session")
}
