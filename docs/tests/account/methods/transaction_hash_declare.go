package main
 
import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
 
	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"
	"github.com/joho/godotenv"
)
 
func main() {
	ctx := context.Background()
 
	if err := godotenv.Load(); err != nil {
		log.Fatal("Failed to load .env file:", err)
	}
 
	rpcURL := os.Getenv("STARKNET_RPC_URL")
	if rpcURL == "" {
		log.Fatal("STARKNET_RPC_URL not set")
	}
 
	provider, err := rpc.NewProvider(ctx, rpcURL)
	if err != nil {
		log.Fatal("Failed to create provider:", err)
	}
 
	accountAddress := os.Getenv("ACCOUNT_ADDRESS")
	publicKey := os.Getenv("ACCOUNT_PUBLIC_KEY")
	privateKey := os.Getenv("ACCOUNT_PRIVATE_KEY")
 
	if accountAddress == "" || publicKey == "" || privateKey == "" {
		log.Fatal("Account credentials not set in .env")
	}
 
	ks := account.NewMemKeystore()
	privKeyBI, ok := new(big.Int).SetString(privateKey, 0)
	if !ok {
		log.Fatal("Failed to parse private key")
	}
	ks.Put(publicKey, privKeyBI)
 
	accountAddressFelt, err := utils.HexToFelt(accountAddress)
	if err != nil {
		log.Fatal("Failed to parse account address:", err)
	}
 
	accnt, err := account.NewAccount(provider, accountAddressFelt, publicKey, ks, account.CairoV2)
	if err != nil {
		log.Fatal("Failed to create account:", err)
	}
 
	nonce, _ := accnt.Nonce(ctx)
	classHash, _ := utils.HexToFelt("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
	compiledClassHash, _ := utils.HexToFelt("0xabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd")
 
	declareTxnV3 := &rpc.DeclareTxnV3{
		Type:              rpc.TransactionTypeDeclare,
		Version:           rpc.TransactionV3,
		SenderAddress:     accountAddressFelt,
		Nonce:             nonce,
		ClassHash:         classHash,
		CompiledClassHash: compiledClassHash,
		ResourceBounds: &rpc.ResourceBoundsMapping{
			L1Gas:     rpc.ResourceBounds{MaxAmount: "0x186a0", MaxPricePerUnit: "0x5f5e100"},
			L2Gas:     rpc.ResourceBounds{MaxAmount: "0x186a0", MaxPricePerUnit: "0x5f5e100"},
			L1DataGas: rpc.ResourceBounds{MaxAmount: "0x186a0", MaxPricePerUnit: "0x5f5e100"},
		},
		Tip:                   "0x0",
		PayMasterData:         []*felt.Felt{},
		AccountDeploymentData: []*felt.Felt{},
		NonceDataMode:         rpc.DAModeL1,
		FeeMode:               rpc.DAModeL1,
		Signature:             []*felt.Felt{},
	}
 
	txHashV3, err := accnt.TransactionHashDeclare(declareTxnV3)
	if err != nil {
		log.Fatal("Failed to compute V3 hash:", err)
	}
 
	fmt.Printf("V3 Transaction:\n")
	fmt.Printf("  Sender:            %s\n", declareTxnV3.SenderAddress.String())
	fmt.Printf("  Nonce:             %d\n", declareTxnV3.Nonce.Uint64())
	fmt.Printf("  ClassHash:         %s\n", classHash.String())
	fmt.Printf("  CompiledClassHash: %s\n", compiledClassHash.String())
	fmt.Printf("  Hash:              %s\n\n", txHashV3.String())
 
	declareTxnV2 := &rpc.DeclareTxnV2{
		Type:              rpc.TransactionTypeDeclare,
		Version:           rpc.TransactionV2,
		SenderAddress:     accountAddressFelt,
		Nonce:             nonce,
		ClassHash:         classHash,
		CompiledClassHash: compiledClassHash,
		MaxFee:            new(felt.Felt).SetUint64(1000000000000000),
		Signature:         []*felt.Felt{},
	}
 
	txHashV2, err := accnt.TransactionHashDeclare(declareTxnV2)
	if err != nil {
		log.Fatal("Failed to compute V2 hash:", err)
	}
 
	fmt.Printf("V2 Transaction:\n")
	fmt.Printf("  Sender:  %s\n", declareTxnV2.SenderAddress.String())
	fmt.Printf("  Nonce:   %d\n", declareTxnV2.Nonce.Uint64())
	fmt.Printf("  MaxFee:  %s\n", declareTxnV2.MaxFee.String())
	fmt.Printf("  Hash:    %s\n\n", txHashV2.String())
 
	declareTxnV1 := &rpc.DeclareTxnV1{
		Type:          rpc.TransactionTypeDeclare,
		Version:       rpc.TransactionV1,
		SenderAddress: accountAddressFelt,
		Nonce:         nonce,
		ClassHash:     classHash,
		MaxFee:        new(felt.Felt).SetUint64(1000000000000000),
		Signature:     []*felt.Felt{},
	}
 
	txHashV1, err := accnt.TransactionHashDeclare(declareTxnV1)
	if err != nil {
		log.Fatal("Failed to compute V1 hash:", err)
	}
 
	fmt.Printf("V1 Transaction:\n")
	fmt.Printf("  Sender: %s\n", declareTxnV1.SenderAddress.String())
	fmt.Printf("  Nonce:  %d\n", declareTxnV1.Nonce.Uint64())
	fmt.Printf("  Hash:   %s\n", txHashV1.String())
}