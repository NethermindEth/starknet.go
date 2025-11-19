package main
 
import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
 
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"
	"github.com/joho/godotenv"
)
 
func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
 
	// Get RPC URL from environment variable
	rpcURL := os.Getenv("STARKNET_RPC_URL")
	if rpcURL == "" {
		log.Fatal("STARKNET_RPC_URL not found in .env file")
	}
 
	// Initialize provider
	provider, err := rpc.NewProvider(context.Background(), rpcURL)
	if err != nil {
		log.Fatal(err)
	}
 
	ctx := context.Background()
 
	// Account address
	senderAddress, err := utils.HexToFelt("0x36d67ab362562a97f9fba8a1051cf8e37ff1a1449530fb9f1f0e32ac2da7d06")
	if err != nil {
		log.Fatal(err)
	}
 
	// Get current nonce
	nonce, err := provider.Nonce(ctx, rpc.WithBlockTag("latest"), senderAddress)
	if err != nil {
		log.Fatal(err)
	}
 
	fmt.Printf("Current nonce: %s\n", nonce)
 
	// Build transaction JSON with current nonce
	txData := fmt.Sprintf(`{
		"type": "INVOKE",
		"version": "0x3",
		"nonce": "%s",
		"sender_address": "0x36d67ab362562a97f9fba8a1051cf8e37ff1a1449530fb9f1f0e32ac2da7d06",
		"signature": [
		  "0x33a831e9428920f71c1df9248d4dbf9101fb5ee2bd100c0ad0d10c94c28dfe3",
		  "0x3fa865114ae29b2a49469401e11eb0db953a7d854916512c2ed400320405c8a"
		],
		"calldata": [
		  "0x1",
		  "0x669e24364ce0ae7ec2864fb03eedbe60cfbc9d1c74438d10fa4b86552907d54",
		  "0x2f0b3c5710379609eb5495f1ecd348cb28167711b73609fe565a72734550354",
		  "0x2",
		  "0xffffffff",
		  "0x0"
		],
		"resource_bounds": {
		  "l1_data_gas": {
			"max_amount": "0x1e0",
			"max_price_per_unit": "0x922"
		  },
		  "l1_gas": {
			"max_amount": "0x0",
			"max_price_per_unit": "0xfbfdefe2186"
		  },
		  "l2_gas": {
			"max_amount": "0x16eea0",
			"max_price_per_unit": "0x1830e58f7"
		  }
		},
		"tip": "0x0",
		"paymaster_data": [],
		"account_deployment_data": [],
		"nonce_data_availability_mode": "L1",
		"fee_data_availability_mode": "L1"
	}`, nonce)
 
	var invokeTx rpc.BroadcastInvokeTxnV3
	if err := json.Unmarshal([]byte(txData), &invokeTx); err != nil {
		log.Fatal(err)
	}
 
	// Estimate fee with SKIP_VALIDATE flag
	result, err := provider.EstimateFee(
		ctx,
		[]rpc.BroadcastTxn{invokeTx},
		[]rpc.SimulationFlag{rpc.SkipValidate},
		rpc.WithBlockTag("latest"),
	)
	if err != nil {
		log.Fatal(err)
	}
 
	resultJSON, _ := json.MarshalIndent(result, "", "  ")
	fmt.Printf("Fee estimate:\n%s\n", resultJSON)
}