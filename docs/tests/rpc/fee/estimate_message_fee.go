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
 
	// L1 contract address (Ethereum address)
	fromAddress := "0x8453fc6cd1bcfe8d4dfc069c400b433054d47bdc"
 
	// L2 contract address (Starknet address)
	toAddress, err := utils.HexToFelt("0x04c5772d1914fe6ce891b64eb35bf3522aeae1315647314aac58b01137607f3f")
	if err != nil {
		log.Fatal(err)
	}
 
	// L1 handler selector
	selector, err := utils.HexToFelt("0x1b64b1b3b690b43b9b514fb81377518f4039cd3e4f4914d8a6bdf01d679fb19")
	if err != nil {
		log.Fatal(err)
	}
 
	// Message payload
	payload, err := utils.HexArrToFelt([]string{
		"0x455448",
		"0x2f14d277fc49e0e2d2967d019aea8d6bd9cb3998",
		"0x02000e6213e24b84012b1f4b1cbd2d7a723fb06950aeab37bedb6f098c7e051a",
		"0x01a055690d9db80000",
		"0x00",
	})
	if err != nil {
		log.Fatal(err)
	}
 
	// Create L1->L2 message
	l1Handler := rpc.MsgFromL1{
		FromAddress: fromAddress,
		ToAddress:   toAddress,
		Selector:    selector,
		Payload:     payload,
	}
 
	// Estimate message fee
	blockNumber := uint64(523066)
	result, err := provider.EstimateMessageFee(ctx, l1Handler, rpc.WithBlockNumber(blockNumber))
	if err != nil {
		log.Fatal(err)
	}
 
	resultJSON, _ := json.MarshalIndent(result, "", "  ")
	fmt.Printf("Estimate Message Fee:\n%s\n", resultJSON)
}