package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/contracts"
	"github.com/NethermindEth/starknet.go/hash"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	sierraPath := "counter_contract/target/dev/counter_Counter.contract_class.json"
	casmPath := "counter_contract/target/dev/counter_Counter.compiled_contract_class.json"

	sierraData, err := os.ReadFile(sierraPath)
	if err != nil {
		log.Fatalf("Failed to read Sierra file: %v", err)
	}

	var contractClass contracts.ContractClass
	if err := json.Unmarshal(sierraData, &contractClass); err != nil {
		log.Fatalf("Failed to parse Sierra JSON: %v", err)
	}

	casmData, err := os.ReadFile(casmPath)
	if err != nil {
		log.Fatalf("Failed to read CASM file: %v", err)
	}

	var casmClass contracts.CasmClass
	if err := json.Unmarshal(casmData, &casmClass); err != nil {
		log.Fatalf("Failed to parse CASM JSON: %v", err)
	}

	compiledClassHash, err := hash.CompiledClassHash(&casmClass)
	if err != nil {
		log.Fatalf("Failed to calculate compiled class hash: %v", err)
	}

	chainID, _ := new(felt.Felt).SetString("0x534e5f5345504f4c4941")
	senderAddress, _ := new(felt.Felt).SetString("0x04718f5a0fc34cc1af16a1cdee98ffb20c31f5cd61d6ab07201858f4287c938d")
	nonce := new(felt.Felt).SetUint64(18)

	resourceBounds := &rpc.ResourceBoundsMapping{
		L1Gas: rpc.ResourceBounds{
			MaxAmount:       "0x5000",
			MaxPricePerUnit: "0x33a937098d80",
		},
		L1DataGas: rpc.ResourceBounds{
			MaxAmount:       "0x5000",
			MaxPricePerUnit: "0x33a937098d80",
		},
		L2Gas: rpc.ResourceBounds{
			MaxAmount:       "0x5000",
			MaxPricePerUnit: "0x10c388d00",
		},
	}

	txn := &rpc.BroadcastDeclareTxnV3{
		Type:                  rpc.TransactionTypeDeclare,
		Version:               rpc.TransactionV3,
		SenderAddress:         senderAddress,
		ContractClass:         &contractClass,
		CompiledClassHash:     compiledClassHash,
		Nonce:                 nonce,
		ResourceBounds:        resourceBounds,
		Tip:                   "0x0",
		PayMasterData:         []*felt.Felt{},
		AccountDeploymentData: []*felt.Felt{},
		NonceDataMode:         rpc.DAModeL1,
		FeeMode:               rpc.DAModeL1,
		Signature:             []*felt.Felt{},
	}

	txHash, err := hash.TransactionHashBroadcastDeclareV3(txn, chainID)
	if err != nil {
		log.Fatalf("Failed to calculate transaction hash: %v", err)
	}

	fmt.Printf("Transaction Hash: %s\n", txHash.String())

	if rpcURL := os.Getenv("STARKNET_RPC_URL"); rpcURL != "" {
		verifyTransaction(txHash, rpcURL)
	}
}

func verifyTransaction(txHash *felt.Felt, rpcURL string) {
	client, err := rpc.NewProvider(context.Background(), rpcURL)
	if err != nil {
		log.Printf("Warning: Could not connect to RPC: %v", err)
		return
	}

	ctx := context.Background()
	tx, err := client.TransactionByHash(ctx, txHash)

	if err == nil {
		fmt.Printf("\nVerification: FOUND on-chain\n")
		fmt.Printf("Type: %T\n", tx)
	} else {
		fmt.Printf("\nVerification: NOT FOUND\n")
	}
}
