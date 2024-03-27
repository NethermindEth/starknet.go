package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"
	"github.com/joho/godotenv"
)

var (
	network                string = "testnet"
	predeployedClassHash          = "0x2794ce20e5f2ff0d40e632cb53845b9f4e526ebd8471983f7dbd355b721d5a"
	accountAddress                = "0xdeadbeef"
	accountContractVersion        = 0 //Replace with the cairo version of your account contract
)

// main initializes the client, sets up the account, deploys a contract, and sends a transaction to the network.
//
// It loads environment variables, dials the Ethereum RPC, creates a new account, casts the account address to a felt type,
// sets up the account using the client, converts the predeployed class hash to a felt type, creates transaction data,
// precomputes an address, prompts the user to add funds to the precomputed address, signs the transaction,
// and finally sends the transaction to the network.
//
// Parameters:
//
//	none
//
// Returns:
//
//	none
func main() {
	// Initialise the client.
	godotenv.Load(fmt.Sprintf(".env.%s", network))
	url := os.Getenv("INTEGRATION_BASE")

	clientv02, err := rpc.NewProvider(url)
	if err != nil {
		log.Fatal(fmt.Sprintf("Error dialing the RPC provider: %s", err))
	}

	// Get random keys for test purposes
	ks, pub, _ := account.GetRandomKeys()

	accountAddressFelt, err := new(felt.Felt).SetString(accountAddress)
	if err != nil {
		panic("Error casting accountAddress to felt")
	}

	// Set up account
	acnt, err := account.NewAccount(clientv02, accountAddressFelt, pub.String(), ks, accountContractVersion)
	if err != nil {
		panic(err)
	}

	classHash, err := utils.HexToFelt(predeployedClassHash)
	if err != nil {
		panic(err)
	}

	// Create transaction data
	tx := rpc.DeployAccountTxn{
		Nonce:               &felt.Zero, // Contract accounts start with nonce zero.
		MaxFee:              new(felt.Felt).SetUint64(4724395326064),
		Type:                rpc.TransactionType_DeployAccount,
		Version:             rpc.TransactionV1,
		Signature:           []*felt.Felt{},
		ClassHash:           classHash,
		ContractAddressSalt: pub,
		ConstructorCalldata: []*felt.Felt{pub},
	}

	precomputedAddress, err := acnt.PrecomputeAddress(&felt.Zero, pub, classHash, tx.ConstructorCalldata)
	fmt.Println("precomputedAddress:", precomputedAddress)

	// At this point you need to add funds to precomputed address to use it.
	var input string

	fmt.Println("The `precomputedAddress` account needs to have enough ETH to perform a transaction.")
	fmt.Println("Use the starknet faucet to send ETH to your `precomputedAddress`")
	fmt.Println("When your account has been funded by the faucet, press any key, then `enter` to continue : ")
	fmt.Scan(&input)

	// Sign the transaction
	err = acnt.SignDeployAccountTransaction(context.Background(), &tx, precomputedAddress)
	if err != nil {
		panic(err)
	}

	// Send transaction to the network
	resp, err := acnt.AddDeployAccountTransaction(context.Background(), rpc.BroadcastDeployAccountTxn{DeployAccountTxn: tx})
	if err != nil {
		panic(fmt.Sprintf("Error returned from AddDeployAccountTransaction: %s", err))
	}
	fmt.Println("AddDeployAccountTransaction response:", resp)
}
