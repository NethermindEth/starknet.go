package main

import (
	"context"
	"fmt"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/client"
	"github.com/NethermindEth/starknet.go/curve"
	setup "github.com/NethermindEth/starknet.go/examples/internal"
	"github.com/NethermindEth/starknet.go/internal/utils"
	pm "github.com/NethermindEth/starknet.go/paymaster"
)

// An example of how to deploy an account and invoke a function in the same request using a paymaster.
func deployAndInvokeWithPaymaster() {
	fmt.Println("Starting paymaster example - deploy_and_invoke")

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
	fmt.Print("Step 1: Build the deploy_and_invoke transaction\n\n")

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
		ConstructorCalldata: constructorCalldata,
		SignatureData:       []*felt.Felt{}, // Optional. For the OZ account, we don't need to add anything in the signature data.
		Version:             2,              // The OZ account version is 2.
	}

	// The next step is to define what we want to execute.
	// The `deploy_and_invoke` transaction type requires both the deploy and invoke data in order to
	// deploy the account and invoke a function within the same request.

	// Here, we will execute a `mint` function in the `RAND_ERC20_CONTRACT_ADDRESS` contract, with the amount of `0xffffffff`.
	amount, _ := utils.HexToU256Felt("0xffffffff")
	invokeData := &pm.UserInvoke{
		UserAddress: precAddress, // The `user_address` is the address of the account that will be deployed.
		Calls: []pm.Call{
			{ // These fields were explained in the `main.go` file of this same example.
				To:       RandERC20ContractAddress,
				Selector: utils.GetSelectorFromNameFelt("mint"),
				Calldata: amount,
			},
		},
	}

	// With the deploy and invoke data, we can build the transaction by calling the `paymaster_buildTransaction` method.
	// REMEMBER: this will only work if you have a valid API key configured.
	//
	// A full explanation about the paymaster_buildTransaction method can be found in the `main.go` file of this same example.
	builtTxn, err := paymaster.BuildTransaction(context.Background(), &pm.BuildTransactionRequest{
		Transaction: &pm.UserTransaction{
			Type: pm.UserTxnDeployAndInvoke, // we are building an `deploy_and_invoke` transaction

			// Both the deploy and invoke data are required.
			Deployment: deployData,
			Invoke:     invokeData,
		},
		Parameters: &pm.UserParameters{
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
		panic(fmt.Sprintf("Error building the deploy_and_invoke transaction: %s", err))
	}
	fmt.Println("Transaction successfully built by the paymaster")
	PrettyPrint(builtTxn)

	fmt.Println("Step 2: Sign the transaction")

	// Now that we have the built transaction, we need to sign it.
	// Differently from the `deploy` transaction, where we just deploy a new account, in the `deploy_and_invoke`
	// we both deploy the account and invoke a function using it. This function request needs to be signed by the account.

	// The signing process consists of signing the SNIP-12 typed data contained in the built transaction.

	// Firstly, get the message hash of the typed data using our precomputed account address as input.
	messageHash, err := builtTxn.TypedData.GetMessageHash(precAddress.String())
	if err != nil {
		panic(fmt.Sprintf("Error getting the message hash of the typed data: %s", err))
	}
	fmt.Println("Message hash of the typed data:", messageHash)

	// Now, we sign the message hash using our account.
	r, s, err := curve.SignFelts(
		messageHash,
		privK,
	) // You can also use the `curve` package to sign the message hash.
	if err != nil {
		panic(fmt.Sprintf("Error signing the transaction: %s", err))
	}
	signature := []*felt.Felt{r, s}

	fmt.Println("Transaction successfully signed")
	PrettyPrint(signature)

	fmt.Println("Step 3: Send the signed transaction")

	// With our built deploy_and_invoke transaction, we can send it to the paymaster by calling the `paymaster_executeTransaction` method.
	response, err := paymaster.ExecuteTransaction(
		context.Background(),
		&pm.ExecuteTransactionRequest{
			Transaction: &pm.ExecutableUserTransaction{
				Type: pm.UserTxnDeployAndInvoke,

				Deployment: builtTxn.Deployment, // The deployment data is the same. We can use our `deployData` variable, or
				// the `builtTxn.Deployment` value.
				Invoke: &pm.ExecutableUserInvoke{
					UserAddress: precAddress,        // The `user_address` is the address of the account that will be deployed.
					TypedData:   builtTxn.TypedData, // The typed data returned by the `paymaster_buildTransaction` method.
					Signature:   signature,          // The signature of the message hash made in the previous step.
				},
			},
			Parameters: &pm.UserParameters{
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
		panic(
			fmt.Sprintf(
				"Error executing the deploy_and_invoke transaction with the paymaster: %s",
				err,
			),
		)
	}

	fmt.Println("Deploy_and_invoke transaction successfully executed by the paymaster")
	fmt.Println("Tracking ID:", response.TrackingID)
	fmt.Println("Transaction Hash:", response.TransactionHash)
}
