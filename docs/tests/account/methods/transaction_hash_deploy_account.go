package main
 
import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
 
	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/contracts"
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
 
	classHash, _ := utils.HexToFelt("0x036078334509b514626504edc9fb252328d1a240e4e948bef8d0c08dff45927f")
	salt := new(felt.Felt).SetUint64(12345)
	newAccountPubKey, _ := utils.HexToFelt("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
	constructorCalldata := []*felt.Felt{newAccountPubKey}
 
	precomputedAddress := contracts.PrecomputeAddress(
		new(felt.Felt).SetUint64(0),
		salt,
		classHash,
		constructorCalldata,
	)
 
	fmt.Printf("ClassHash:           %s\n", classHash.String())
	fmt.Printf("Salt:                %s\n", salt.String())
	fmt.Printf("PublicKey:           %s\n", newAccountPubKey.String())
	fmt.Printf("Precomputed Address: %s\n\n", precomputedAddress.String())
 
	deployTxnV3 := &rpc.DeployAccountTxnV3{
		Type:                rpc.TransactionTypeDeployAccount,
		Version:             rpc.TransactionV3,
		Nonce:               new(felt.Felt).SetUint64(0),
		ClassHash:           classHash,
		ContractAddressSalt: salt,
		ConstructorCalldata: constructorCalldata,
		ResourceBounds: &rpc.ResourceBoundsMapping{
			L1Gas:     rpc.ResourceBounds{MaxAmount: "0x186a0", MaxPricePerUnit: "0x5f5e100"},
			L2Gas:     rpc.ResourceBounds{MaxAmount: "0x186a0", MaxPricePerUnit: "0x5f5e100"},
			L1DataGas: rpc.ResourceBounds{MaxAmount: "0x186a0", MaxPricePerUnit: "0x5f5e100"},
		},
		Tip:           "0x0",
		PayMasterData: []*felt.Felt{},
		NonceDataMode: rpc.DAModeL1,
		FeeMode:       rpc.DAModeL1,
		Signature:     []*felt.Felt{},
	}
 
	txHashV3, err := accnt.TransactionHashDeployAccount(deployTxnV3, precomputedAddress)
	if err != nil {
		log.Fatal("Failed to compute V3 hash:", err)
	}
 
	fmt.Printf("V3 Transaction:\n")
	fmt.Printf("  Nonce:     %d\n", deployTxnV3.Nonce.Uint64())
	fmt.Printf("  ClassHash: %s\n", deployTxnV3.ClassHash.String())
	fmt.Printf("  Salt:      %s\n", deployTxnV3.ContractAddressSalt.String())
	fmt.Printf("  Hash:      %s\n\n", txHashV3.String())
 
	deployTxnV1 := &rpc.DeployAccountTxnV1{
		Type:                rpc.TransactionTypeDeployAccount,
		Version:             rpc.TransactionV1,
		Nonce:               new(felt.Felt).SetUint64(0),
		ClassHash:           classHash,
		ContractAddressSalt: salt,
		ConstructorCalldata: constructorCalldata,
		MaxFee:              new(felt.Felt).SetUint64(1000000000000000),
		Signature:           []*felt.Felt{},
	}
 
	txHashV1, err := accnt.TransactionHashDeployAccount(deployTxnV1, precomputedAddress)
	if err != nil {
		log.Fatal("Failed to compute V1 hash:", err)
	}
 
	fmt.Printf("V1 Transaction:\n")
	fmt.Printf("  Nonce:  %d\n", deployTxnV1.Nonce.Uint64())
	fmt.Printf("  MaxFee: %s\n", deployTxnV1.MaxFee.String())
	fmt.Printf("  Hash:   %s\n", txHashV1.String())
}