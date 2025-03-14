package main

import (
	"context"
	"fmt"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"

	setup "github.com/NethermindEth/starknet.go/examples/internal"
)

var (
	// OpenZeppelin Account Class Hash in Sepolia
	predeployedClassHash = "0x61dac032f228abef9c6626f995015233097ae253a7f72d68552db02f2971b8f"
)

// main initializes the client, sets up the account, deploys a contract, and sends a transaction to the network.
//
// It loads environment variables, dials the Starknet Sepolia RPC, creates a new account, casts the account address to a felt type,
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
	// Load variables from '.env' file
	rpcProviderUrl := setup.GetRpcProviderUrl()

	// Initialise the client.
	client, err := rpc.NewProvider(rpcProviderUrl)
	if err != nil {
		panic(err)
	}

	// Get random keys for test purposes
	ks, pub, privKey := account.GetRandomKeys()
	fmt.Printf("Generated public key: %v\n", pub)
	fmt.Printf("Generated private key: %v\n", privKey)

	// Set up the account passing random values to 'accountAddress' and 'cairoVersion' variables,
	// as for this case we only need the 'ks' to sign the deploy transaction.
	accnt, err := account.NewAccount(client, pub, pub.String(), ks, 2)
	if err != nil {
		panic(err)
	}

	classHash, err := utils.HexToFelt(predeployedClassHash)
	if err != nil {
		panic(err)
	}

	// Precompute the address of the new account
	// In our case, the OZ account constructor requires the public key of the account as calldata, so we pass it as a parameter
	deployAccountTxn, precomputedAddress, err := accnt.BuildAndEstimateDeployAccountTxn(context.Background(), pub, classHash, []*felt.Felt{pub}, 1.5)
	if err != nil {
		panic(err)
	}

	fmt.Println("PrecomputedAddress:", setup.PadZerosInFelt(precomputedAddress))

	overallFee, err := utils.ResBoundsMapToOverallFee(deployAccountTxn.ResourceBounds, 1.5)
	if err != nil {
		panic(err)
	}
	feeInSTRK := utils.FRIToSTRK(overallFee)

	// At this point you need to add funds to precomputed address to use it.
	var input string

	fmt.Println("\nThe `precomputedAddress` account needs to have enough STRK to perform a transaction.")
	fmt.Printf("You can use the starknet faucet or send STRK to your `precomputedAddress`. You need aproximately %f STRK. \n", feeInSTRK)
	fmt.Println("When your account has been funded, press any key, then `enter` to continue : ")
	fmt.Scan(&input)

	// Send transaction to the network
	resp, err := accnt.SendTransaction(context.Background(), deployAccountTxn)
	if err != nil {
		fmt.Println("Error returned from AddDeployAccountTransaction: ")
		setup.PanicRPC(err)
	}

	fmt.Println("AddDeployAccountTransaction successfully submitted! Wait a few minutes to see it in Voyager.")
	fmt.Printf("Transaction hash: %v \n", resp.TransactionHash)
	fmt.Printf("Contract address: %v \n", setup.PadZerosInFelt(resp.ContractAddress))
}
