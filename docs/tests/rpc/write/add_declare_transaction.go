package main
 
import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"
 
	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/contracts"
	"github.com/NethermindEth/starknet.go/hash"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/joho/godotenv"
)
 
func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
 
	// Get configuration from environment
	rpcURL := os.Getenv("STARKNET_RPC_URL")
	privateKeyStr := os.Getenv("ACCOUNT_PRIVATE_KEY")
	publicKeyStr := os.Getenv("ACCOUNT_PUBLIC_KEY")
	accountAddressStr := os.Getenv("ACCOUNT_ADDRESS")
 
	if rpcURL == "" || privateKeyStr == "" || publicKeyStr == "" || accountAddressStr == "" {
		log.Fatal("Missing required environment variables")
	}
 
	// Parse credentials
	privateKey, err := new(felt.Felt).SetString(privateKeyStr)
	if err != nil {
		log.Fatal(err)
	}
 
	publicKey, err := new(felt.Felt).SetString(publicKeyStr)
	if err != nil {
		log.Fatal(err)
	}
 
	accountAddress, err := new(felt.Felt).SetString(accountAddressStr)
	if err != nil {
		log.Fatal(err)
	}
 
	// Initialize provider
	ctx := context.Background()
	client, err := rpc.NewProvider(ctx, rpcURL)
	if err != nil {
		log.Fatal(err)
	}
 
	// Create keystore
	ks := account.NewMemKeystore()
	ks.Put(publicKey.String(), privateKey.BigInt(new(big.Int)))
 
	// Create account controller
	acct, err := account.NewAccount(client, accountAddress, publicKey.String(), ks, 2)
	if err != nil {
		log.Fatal(err)
	}
 
	// Get current nonce
	nonce, err := acct.Nonce(ctx)
	if err != nil {
		log.Fatal(err)
	}
 
	// Load the compiled contract files
	// Replace with your actual contract paths
	sierraPath := "counter_contract/target/dev/counter_Counter.contract_class.json"
	casmPath := "counter_contract/target/dev/counter_Counter.compiled_contract_class.json"
 
	// Check if files exist
	if _, err := os.Stat(sierraPath); os.IsNotExist(err) {
		log.Fatal("Contract files not found! Run: cd counter_contract && scarb build")
	}
 
	// Load Sierra contract
	sierraContent, err := os.ReadFile(sierraPath)
	if err != nil {
		log.Fatal(err)
	}
 
	var contractClass contracts.ContractClass
	if err := json.Unmarshal(sierraContent, &contractClass); err != nil {
		log.Fatal(err)
	}
 
	// Load CASM contract
	casmContent, err := os.ReadFile(casmPath)
	if err != nil {
		log.Fatal(err)
	}
 
	var casmClass contracts.CasmClass
	if err := json.Unmarshal(casmContent, &casmClass); err != nil {
		log.Fatal(err)
	}
 
	// Calculate the compiled class hash first
	// Note: There's sometimes a mismatch between local calculation and what the sequencer expects
	// If you get a mismatch error, use the "Expected" hash from the error message
	compiledClassHash, err := hash.CompiledClassHash(&casmClass)
	if err != nil {
		log.Fatal(err)
	}
 
	// Override with the correct hash if there's a known mismatch
	// IMPORTANT: Replace this with the "Expected" hash from your error message if you get a mismatch
	// Example: compiledClassHash, _ = new(felt.Felt).SetString("0x4fe67cc3cb8e0e3c06161a5b8ccaed841e5e3116138bef832b19fa298c74f6c")
 
	// Calculate class hash from the contract
	classHash := hash.ClassHash(&contractClass)
 
	fmt.Printf("Class Hash: %s\n", classHash.String())
	fmt.Printf("Compiled Class Hash: %s\n", compiledClassHash.String())
 
	// Check if contract is already declared
	fmt.Println("Checking if contract is already declared...")
	_, err = client.ClassAt(ctx, rpc.WithBlockTag(rpc.BlockTagLatest), classHash)
	if err == nil {
		fmt.Println("Contract is already declared!")
		fmt.Printf("Class Hash: %s\n", classHash.String())
		return
	}
 
	// Try alternative check with Class method
	_, err = client.Class(ctx, rpc.WithBlockTag(rpc.BlockTagLatest), classHash)
	if err == nil {
		fmt.Println("Contract is already declared!")
		fmt.Printf("Class Hash: %s\n", classHash.String())
		return
	}
 
	fmt.Println("Contract not found on-chain, proceeding with declaration...")
 
	// Create the declare transaction
	declareTx := &rpc.BroadcastDeclareTxnV3{
		Type:              rpc.TransactionTypeDeclare,
		Version:           rpc.TransactionV3,
		SenderAddress:     accountAddress,
		Nonce:             nonce,
		ContractClass:     &contractClass,
		CompiledClassHash: compiledClassHash,
		Signature:         []*felt.Felt{}, // Will be filled by signing
		ResourceBounds: &rpc.ResourceBoundsMapping{
			L1Gas: rpc.ResourceBounds{
				MaxAmount:       rpc.U64("0x186a0"),         // 100000
				MaxPricePerUnit: rpc.U128("0x33a937098d80"), // High price for estimation
			},
			L1DataGas: rpc.ResourceBounds{
				MaxAmount:       rpc.U64("0x186a0"),         // 100000
				MaxPricePerUnit: rpc.U128("0x33a937098d80"), // High price for estimation
			},
			L2Gas: rpc.ResourceBounds{
				MaxAmount:       rpc.U64("0x6b5c0"),      // 440000 (minimum required)
				MaxPricePerUnit: rpc.U128("0x10c388d00"), // High price for estimation
			},
		},
		Tip:                   rpc.U64("0x0"),
		PayMasterData:         []*felt.Felt{},
		AccountDeploymentData: []*felt.Felt{},
		NonceDataMode:         rpc.DAModeL1,
		FeeMode:               rpc.DAModeL1,
	}
 
	// Estimate fee
	fmt.Println("Estimating transaction fee...")
	simFlags := []rpc.SimulationFlag{rpc.SkipValidate}
	feeEstimate, err := client.EstimateFee(
		ctx,
		[]rpc.BroadcastTxn{declareTx},
		simFlags,
		rpc.WithBlockTag(rpc.BlockTagLatest),
	)
 
	if err != nil {
		fmt.Printf("Fee estimation failed: %v\n", err)
		fmt.Println("Continuing with default resource bounds...")
	} else if len(feeEstimate) > 0 {
		// Update resource bounds based on estimate with 20% buffer
		estimatedL1DataGas := feeEstimate[0].L1DataGasConsumed.Uint64()
		estimatedL2Gas := feeEstimate[0].L2GasConsumed.Uint64()
 
		declareTx.ResourceBounds = &rpc.ResourceBoundsMapping{
			L1Gas: rpc.ResourceBounds{
				MaxAmount:       rpc.U64(fmt.Sprintf("0x%x", estimatedL1DataGas*12/10)),
				MaxPricePerUnit: rpc.U128("0x33a937098d80"),
			},
			L1DataGas: rpc.ResourceBounds{
				MaxAmount:       rpc.U64(fmt.Sprintf("0x%x", estimatedL1DataGas*12/10)),
				MaxPricePerUnit: rpc.U128("0x33a937098d80"),
			},
			L2Gas: rpc.ResourceBounds{
				MaxAmount:       rpc.U64(fmt.Sprintf("0x%x", estimatedL2Gas*12/10)),
				MaxPricePerUnit: rpc.U128("0x10c388d00"),
			},
		}
 
		fmt.Printf("Estimated L1 Data Gas: %d\n", estimatedL1DataGas)
		fmt.Printf("Estimated L2 Gas: %d\n", estimatedL2Gas)
		fmt.Printf("Overall Fee: %s STRK\n", feeEstimate[0].OverallFee.String())
	}
 
	// Sign the transaction
	// Create a DeclareTxnV3 for signing (required for proper hash calculation)
	declareTxnV3 := rpc.DeclareTxnV3{
		Type:                  declareTx.Type,
		SenderAddress:         declareTx.SenderAddress,
		CompiledClassHash:     declareTx.CompiledClassHash,
		Version:               declareTx.Version,
		Signature:             declareTx.Signature,
		Nonce:                 declareTx.Nonce,
		ClassHash:             classHash, // Important: ClassHash is required for signing
		ResourceBounds:        declareTx.ResourceBounds,
		Tip:                   declareTx.Tip,
		PayMasterData:         declareTx.PayMasterData,
		AccountDeploymentData: declareTx.AccountDeploymentData,
		NonceDataMode:         declareTx.NonceDataMode,
		FeeMode:               declareTx.FeeMode,
	}
 
	err = acct.SignDeclareTransaction(ctx, &declareTxnV3)
	if err != nil {
		log.Fatal(err)
	}
 
	// Copy signature back to broadcast transaction
	declareTx.Signature = declareTxnV3.Signature
 
	// Submit the transaction
	resp, err := client.AddDeclareTransaction(ctx, declareTx)
	if err != nil {
		log.Fatal(err)
	}
 
	fmt.Printf("Contract declared successfully!\n")
	fmt.Printf("Transaction Hash: %s\n", resp.Hash.String())
	fmt.Printf("Class Hash: %s\n", resp.ClassHash.String())
 
	// Wait for transaction confirmation
	receipt, err := waitForTransaction(ctx, client, resp.Hash)
	if err != nil {
		log.Fatal(err)
	}
 
	fmt.Printf("Transaction confirmed!\n")
	fmt.Printf("Block Number: %d\n", receipt.BlockNumber)
	fmt.Printf("Status: %s\n", receipt.FinalityStatus)
	fmt.Printf("Actual Fee: %s\n", receipt.ActualFee.Amount.String())
}
 
// waitForTransaction waits for a transaction to be confirmed on the network
func waitForTransaction(ctx context.Context, client *rpc.Provider, txHash *felt.Felt) (*rpc.TransactionReceiptWithBlockInfo, error) {
	for i := 0; i < 60; i++ { // Wait up to 5 minutes
		receipt, err := client.TransactionReceipt(ctx, txHash)
		if err == nil {
			if receipt.FinalityStatus == rpc.TxnFinalityStatusAcceptedOnL2 ||
				receipt.FinalityStatus == rpc.TxnFinalityStatusAcceptedOnL1 {
				return receipt, nil
			}
		}
 
		time.Sleep(5 * time.Second)
		fmt.Print(".")
	}
 
	return nil, fmt.Errorf("transaction confirmation timeout")
}