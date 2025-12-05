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
		log.Fatal("STARKNET_RPC_URL not set in .env file")
	}
 
	provider, err := rpc.NewProvider(ctx, rpcURL)
	if err != nil {
		log.Fatal("Failed to create provider:", err)
	}
 
	accountAddress := os.Getenv("ACCOUNT_ADDRESS")
	publicKey := os.Getenv("ACCOUNT_PUBLIC_KEY")
	privateKey := os.Getenv("ACCOUNT_PRIVATE_KEY")
 
	if accountAddress == "" || publicKey == "" || privateKey == "" {
		log.Fatal("ACCOUNT_ADDRESS, ACCOUNT_PUBLIC_KEY, or ACCOUNT_PRIVATE_KEY not set")
	}
 
	ks := account.NewMemKeystore()
	privKeyBI, ok := new(big.Int).SetString(privateKey, 0)
	if !ok {
		log.Fatal("Failed to parse private key")
	}
	ks.Put(publicKey, privKeyBI)
 
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
 
	fmt.Printf("Class Hash:       %s\n", classHash.String())
	fmt.Printf("Salt:             %s\n", salt.String())
	fmt.Printf("Public Key:       %s\n", newAccountPubKey.String())
	fmt.Printf("Precomputed Addr: %s\n\n", precomputedAddress.String())
 
	deployTxn := &rpc.DeployAccountTxnV3{
		Type:                rpc.TransactionTypeDeployAccount,
		Version:             rpc.TransactionV3,
		Nonce:               new(felt.Felt).SetUint64(0),
		ClassHash:           classHash,
		ContractAddressSalt: salt,
		ConstructorCalldata: constructorCalldata,
		ResourceBounds: &rpc.ResourceBoundsMapping{
			L1Gas: rpc.ResourceBounds{
				MaxAmount:       "0x0",
				MaxPricePerUnit: "0x0",
			},
			L1DataGas: rpc.ResourceBounds{
				MaxAmount:       "0x0",
				MaxPricePerUnit: "0x0",
			},
			L2Gas: rpc.ResourceBounds{
				MaxAmount:       "0x0",
				MaxPricePerUnit: "0x0",
			},
		},
		Tip:           "0x0",
		PayMasterData: []*felt.Felt{},
		NonceDataMode: rpc.DAModeL1,
		FeeMode:       rpc.DAModeL1,
		Signature:     []*felt.Felt{},
	}
 
	fmt.Printf("Nonce:            %d\n", deployTxn.Nonce.Uint64())
	fmt.Printf("Signature (before): %v\n\n", deployTxn.Signature)
 
	newAcctKs := account.NewMemKeystore()
	newPrivKey := new(big.Int).SetBytes([]byte("temp_private_key_for_demo"))
	newAcctKs.Put(newAccountPubKey.String(), newPrivKey)
 
	newAcct, err := account.NewAccount(
		provider,
		precomputedAddress,
		newAccountPubKey.String(),
		newAcctKs,
		account.CairoV2,
	)
	if err != nil {
		log.Fatal("Failed to create new account instance:", err)
	}
 
	err = newAcct.SignDeployAccountTransaction(ctx, deployTxn, precomputedAddress)
	if err != nil {
		log.Fatal("Failed to sign deploy account transaction:", err)
	}
 
	fmt.Printf("Signature components:\n")
	fmt.Printf("  r: %s\n", deployTxn.Signature[0].String())
	fmt.Printf("  s: %s\n", deployTxn.Signature[1].String())
}