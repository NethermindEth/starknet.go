package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	rpcURL := os.Getenv("STARKNET_RPC_URL")
	if rpcURL == "" {
		log.Fatal("STARKNET_RPC_URL not set in .env")
	}

	ctx := context.Background()
	client, err := rpc.NewProvider(ctx, rpcURL)
	if err != nil {
		log.Fatal("Failed to create client:", err)
	}

	fmt.Println("EstimateMessageFee RPC Method")
	fmt.Println("==============================")
	fmt.Println()
	fmt.Println("Method Signature:")
	fmt.Println("  func (provider *Provider) EstimateMessageFee(")
	fmt.Println("      ctx context.Context,")
	fmt.Println("      msg MsgFromL1,")
	fmt.Println("      blockID BlockID,")
	fmt.Println("  ) (MessageFeeEstimation, error)")
	fmt.Println()
	fmt.Println("Requirements:")
	fmt.Println("  - MsgFromL1 structure with:")
	fmt.Println("    - FromAddress: L1 Ethereum address")
	fmt.Println("    - ToAddress: L2 contract address")
	fmt.Println("    - EntryPointSelector: Function selector")
	fmt.Println("    - Payload: Message payload")
	fmt.Println()
	fmt.Println("Use Case:")
	fmt.Println("  Estimate the L2 fee for processing an L1->L2 message")
	fmt.Println("  sent from Ethereum to Starknet")
	fmt.Println()
	fmt.Println("Expected Response Structure:")
	fmt.Println("  MessageFeeEstimation{")
	fmt.Println("    L1GasConsumed:     *felt.Felt")
	fmt.Println("    L1GasPrice:        *felt.Felt")
	fmt.Println("    L2GasConsumed:     *felt.Felt")
	fmt.Println("    L2GasPrice:        *felt.Felt")
	fmt.Println("    L1DataGasConsumed: *felt.Felt")
	fmt.Println("    L1DataGasPrice:    *felt.Felt")
	fmt.Println("    OverallFee:        *felt.Felt")
	fmt.Println("    Unit:              \"WEI\"")
	fmt.Println("  }")
	fmt.Println()
	fmt.Println("Example:")
	fmt.Println("  msg := rpc.MsgFromL1{")
	fmt.Println("    FromAddress: ethAddress,")
	fmt.Println("    ToAddress:   starknetContractAddress,")
	fmt.Println("    EntryPointSelector: selector,")
	fmt.Println("    Payload: []*felt.Felt{...},")
	fmt.Println("  }")
	fmt.Println("  estimate, err := client.EstimateMessageFee(ctx, msg, blockID)")
	fmt.Println()
	fmt.Printf("RPC Provider connected: %T\n", client)
}
