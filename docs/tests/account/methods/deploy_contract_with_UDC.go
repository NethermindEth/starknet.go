package main
 
import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"
 
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
		log.Fatal("ACCOUNT_ADDRESS, ACCOUNT_PUBLIC_KEY, or ACCOUNT_PRIVATE_KEY not set in .env")
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
 
	accnt, err := account.NewAccount(
		provider,
		accountAddressFelt,
		publicKey,
		ks,
		account.CairoV2,
	)
	if err != nil {
		log.Fatal("Failed to create account:", err)
	}
 
	erc20ClassHash, err := utils.HexToFelt("0x073d71c37e20c569186445d2c497d2195b4c0be9a255d72dbad86662fcc63ae6")
	if err != nil {
		log.Fatal("Failed to parse ERC20 class hash:", err)
	}
 
	name, err := utils.StringToByteArrFelt("MyToken")
	if err != nil {
		log.Fatal("Failed to convert name:", err)
	}
 
	symbol, err := utils.StringToByteArrFelt("MTK")
	if err != nil {
		log.Fatal("Failed to convert symbol:", err)
	}
 
	initialSupply, err := utils.HexToU256Felt("0xd3c21bcecceda1000000")
	if err != nil {
		log.Fatal("Failed to convert supply:", err)
	}
 
	recipient := accnt.Address
	owner := accnt.Address
 
	constructorCalldata := make([]*felt.Felt, 0, 20)
	constructorCalldata = append(constructorCalldata, name...)
	constructorCalldata = append(constructorCalldata, symbol...)
	constructorCalldata = append(constructorCalldata, initialSupply...)
	constructorCalldata = append(constructorCalldata, recipient, owner)
 
	txnOpts := &account.TxnOptions{
		FeeMultiplier: 1.5,
		TipMultiplier: 1.0,
	}
 
	udcOpts := &utils.UDCOptions{
		Salt:              nil,
		OriginIndependent: false,
		UDCVersion:        utils.UDCCairoV2,
	}
 
	response, salt, err := accnt.DeployContractWithUDC(
		ctx,
		erc20ClassHash,
		constructorCalldata,
		txnOpts,
		udcOpts,
	)
	if err != nil {
		log.Fatal("Failed to deploy contract:", err)
	}
 
	fmt.Printf("Deploy Transaction Successful:\n")
	fmt.Printf("Transaction Hash: %s\n", response.Hash.String())
	fmt.Printf("Salt Used:        %s\n", salt.String())
 
	deployedAddress := utils.PrecomputeAddressForUDC(
		erc20ClassHash,
		salt,
		constructorCalldata,
		utils.UDCCairoV2,
		accnt.Address,
	)
 
	fmt.Printf("Deployed Address: %s\n", deployedAddress.String())
 
	txReceipt, err := accnt.WaitForTransactionReceipt(
		ctx,
		response.Hash,
		3*time.Second,
	)
	if err != nil {
		log.Fatal("Failed to get transaction receipt:", err)
	}
 
	fmt.Printf("Block Number:     %d\n", txReceipt.BlockNumber)
	fmt.Printf("Execution Status: %s\n", txReceipt.ExecutionStatus)
	fmt.Printf("Finality Status:  %s\n", txReceipt.FinalityStatus)
}