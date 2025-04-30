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

// main initializes the client, sets up the temporary account, precomputes the address of the new account,
// estimates deployment fees, and prepares for the account deployment transaction.
//
// It loads environment variables, initializes a Starknet RPC client, generates random cryptographic keys,
// sets up an account with the generated keys, precomputes the address of a new account, estimates deployment fees,
// and prepares for the account deployment transaction.
func main() {
	// Load variables from '.env' file
	rpcProviderUrl := setup.GetRpcProviderUrl()

	// Initialise the client.
	client, err := rpc.NewProvider(rpcProviderUrl)
	if err != nil {
		panic(err)
	}

	// Get random keys for being able to sign the deploy transaction.
	// These keys will always be used to sign transactions in the new account.
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

	// Build and estimate fees for the deploy account transaction, and precompute the address of the new account.
	// In our case, the OZ account constructor requires the public key of the account as calldata, so we pass it as calldata.
	// The multiplier for the fee estimation is 1.5, as we want to be sure that the transaction will be accepted.
	deployAccountTxn, precomputedAddress, err := accnt.BuildAndEstimateDeployAccountTxn(context.Background(), pub, classHash, []*felt.Felt{pub}, 1.5, false)
	if err != nil {
		panic(err)
	}

	fmt.Println("PrecomputedAddress:", setup.PadZerosInFelt(precomputedAddress))

	// Convert the estimated fee to STRK. The multiplier is 1, as we already estimated the fee in BuildAndEstimateDeployAccountTxn multiplying by 1.5.
	overallFee, err := utils.ResBoundsMapToOverallFee(deployAccountTxn.ResourceBounds, 1)
	if err != nil {
		panic(err)
	}
	feeInSTRK := utils.FRIToSTRK(overallFee)

	// At this point you need to add funds to precomputed address to use it.
	var input string

	fmt.Println("\nThe `precomputedAddress` account needs to have enough STRK to perform a transaction.")
	fmt.Printf("You can use the starknet faucet or send STRK to your `precomputedAddress`. You need approximately %f STRK. \n", feeInSTRK)
	fmt.Println("When your account has been funded, press any key, then `enter` to continue: ")
	fmt.Scan(&input)

	// Send transaction to the network
	resp, err := accnt.SendTransaction(context.Background(), deployAccountTxn)
	if err != nil {
		fmt.Println("Error returned from SendTransaction: ")
		panic(err)
	}

	fmt.Println("BroadcastDeployAccountTxn successfully submitted! Wait a few minutes to see it in Voyager.")
	fmt.Printf("Transaction hash: %v \n", resp.TransactionHash)
	fmt.Printf("Contract address: %v \n", setup.PadZerosInFelt(resp.ContractAddress))
}
