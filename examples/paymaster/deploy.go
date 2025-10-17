package main

import (
	"context"
	"fmt"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/client"
	setup "github.com/NethermindEth/starknet.go/examples/internal"
	"github.com/NethermindEth/starknet.go/internal/utils"
	pm "github.com/NethermindEth/starknet.go/paymaster"
)

// OpenZeppelin account class hash that supports outside executions
const OZAccountClassHash = "0x05b4b537eaa2399e3aa99c4e2e0208ebd6c71bc1467938cd52c798c601e43564"

// An example of how to deploy a contract with a paymaster.
func deployWithPaymaster() {
	fmt.Println("Starting paymaster example - deploying an account")

	// Load variables from '.env' file
	AVNUApiKey := setup.GetAVNUApiKey()

	// Since all accounts in Starknet are smart contracts, we need to deploy them first before we can use them.
	// And to do so, we need to calculate the address of the new account and fund it with
	// enough STRK tokens before deploying it. This tokens will be used to pay the fees for the `deploy` txn.
	//
	// Deploy an account with a paymaster using the `default` fee mode doesn't make much sense, as we will
	// need to send some tokens for the account anyway. So, we will use the `sponsored` fee mode now,
	// which will allow the paymaster to fully cover the fees for the `deploy` txn. This mode requires
	// an API key from an entity. You can only run this example with it.

	// Let's initialise the paymaster client, but now, we will also pass our API key to the client.
	// In the AVNU paymaster, the API key is a http header called `x-paymaster-api-key`.
	// In the current Starknet.go client, you can set a custom http header using the `client.WithHeader` option.
	paymaster, err := pm.New(
		AVNUPaymasterURL,
		client.WithHeader("x-paymaster-api-key", AVNUApiKey),
	)
	if err != nil {
		panic(fmt.Sprintf("Error connecting to the paymaster provider with the API key: %s", err))
	}

	fmt.Println("Established connection with the paymaster provider")
	fmt.Print("Step 1: Build the deploy transaction\n\n")

	// First, let's get all the data we need for deploy an account.
	_, pubKey, privK := account.GetRandomKeys() // Get random keys for the account
	fmt.Println("Public key:", pubKey)
	fmt.Println("Private key:", privK)
	classHash, _ := utils.HexToFelt(
		OZAccountClassHash,
	) // It needs to be an SNIP-9 compatible account
	constructorCalldata := []*felt.Felt{
		pubKey,
	} // The OZ account constructor requires the public key
	salt, _ := utils.HexToFelt("0xdeadbeef") // Just a random salt
	// Precompute the address of the new account based on the salt, class hash and constructor calldata
	precAddress := account.PrecomputeAccountAddress(salt, classHash, constructorCalldata)

	fmt.Println("Precomputed address:", precAddress)

	// Now we can create the deploy data for the transaction.
	deployData := &pm.AccDeploymentData{
		Address:             precAddress, // The precomputed address of the new account
		ClassHash:           classHash,
		Salt:                salt,
		Calldata: constructorCalldata,
		SignatureData:       []*felt.Felt{}, // Optional. For the OZ account, we don't need to add anything in the signature data.
		Version:             2,              // The OZ account version is 2.
	}

	// With the deploy data, we can build the transaction by calling the `paymaster_buildTransaction` method.
	// REMEMBER: this will only work if you have a valid API key configured.
	//
	// A full explanation about the paymaster_buildTransaction method can be found in the `main.go` file of this same example.
	builtTxn, err := paymaster.BuildTransaction(context.Background(), pm.BuildTransactionRequest{
		Transaction: pm.UserTransaction{
			Type:       pm.UserTxnDeploy, // we are building an `deploy` transaction
			Deployment: deployData,
		},
		Parameters: pm.UserParameters{
			Version: pm.UserParamV1,
			FeeMode: pm.FeeMode{
				Mode: pm.FeeModeSponsored, // We then set the fee mode to `sponsored`
				Tip: &pm.TipPriority{
					Priority: pm.TipPriorityNormal,
				},
			},
		},
	})
	if err != nil {
		panic(fmt.Sprintf("Error building the deploy transaction: %s", err))
	}
	fmt.Println("Transaction successfully built by the paymaster")
	PrettyPrint(builtTxn)

	// Since we are deploying an account, we don't need to sign the transaction, just execute it.

	fmt.Println("Step 2: Send the signed transaction")

	// With our built deploy transaction, we can send it to the paymaster by calling the `paymaster_executeTransaction` method.
	response, err := paymaster.ExecuteTransaction(
		context.Background(),
		pm.ExecuteTransactionRequest{
			Transaction: pm.ExecutableUserTransaction{
				Type:       pm.UserTxnDeploy,
				Deployment: builtTxn.Deployment, // The deployment data is the same. We can use our `deployData` variable, or
				// the `builtTxn.Deployment` value.
			},
			Parameters: pm.UserParameters{
				Version: pm.UserParamV1,

				// Using the same fee options as in the `paymaster_buildTransaction` method.
				FeeMode: pm.FeeMode{
					Mode: pm.FeeModeSponsored,
					Tip: &pm.TipPriority{
						Priority: pm.TipPriorityNormal,
					},
				},
			},
		},
	)
	if err != nil {
		panic(fmt.Sprintf("Error executing the deploy transaction with the paymaster: %s", err))
	}

	fmt.Println("Deploy transaction successfully executed by the paymaster")
	fmt.Println("Tracking ID:", response.TrackingID)
	fmt.Println("Transaction Hash:", response.TransactionHash)
}
