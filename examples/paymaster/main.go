package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/NethermindEth/starknet.go/account"
	setup "github.com/NethermindEth/starknet.go/examples/internal"
	"github.com/NethermindEth/starknet.go/internal/utils"
	pm "github.com/NethermindEth/starknet.go/paymaster"
	"github.com/NethermindEth/starknet.go/rpc"
)

func main() {
	fmt.Println("Starting paymaster example")

	// ************* Set up things *************
	//
	// Load variables from '.env' file
	accountAddress := setup.GetAccountAddress()
	accountCairoVersion := setup.GetAccountCairoVersion()
	privateKey := setup.GetPrivateKey()
	publicKey := setup.GetPublicKey()
	rpcProviderUrl := setup.GetRpcProviderUrl()

	// Connect to a RPC provider to instantiate the account
	client, err := rpc.NewProvider(rpcProviderUrl)
	if err != nil {
		panic(fmt.Sprintf("Error dialling the RPC provider: %s", err))
	}

	// Instantiate the account to sign the transaction (we can also use the `curve` pkg for that, we'll see later)
	account := NewAccount(client, accountAddress, privateKey, publicKey, accountCairoVersion)

	// ************* done *************

	// Initialise connection to the paymaster provider - AVNU Sepolia in this case
	paymaster, err := pm.New("https://sepolia.paymaster.avnu.fi")
	if err != nil {
		panic(fmt.Sprintf("Error connecting to the paymaster provider: %s", err))
	}

	fmt.Println("Established connection with the paymaster provider")

	// Check if the paymaster provider is available by calling the `paymaster_isAvailable` method
	available, err := paymaster.IsAvailable(context.Background())
	if err != nil {
		panic(fmt.Sprintf("Error checking if the paymaster provider is available: %s", err))
	}
	fmt.Println("Is paymaster provider available?: ", available)

	// Get the supported tokens by calling the `paymaster_getSupportedTokens` method
	tokens, err := paymaster.GetSupportedTokens(context.Background())
	if err != nil {
		panic(fmt.Sprintf("Error getting the supported tokens: %s", err))
	}
	fmt.Println("\nSupported tokens:")
	PrettyPrint(tokens)

	// Now that we know the paymaster is available and we have the supported tokens list,
	// we can build and execute a transaction with the paymaster, paying the fees with any of
	// the supported tokens.
	// For the sake of simplicity, we will use the STRK token itself.

	// Sending an invoke transaction with a paymaster involves 3 steps:
	// 1. Build the transaction  by calling the `paymaster_buildTransaction` method
	// 2. Sign the transaction built by the paymaster
	// 3. Send the signed transaction by calling the `paymaster_executeTransaction` method

	fmt.Println("Step 1: Build the transaction")

	// a simple ERC20 contract with a public mint function
	simpleERC20, _ := utils.HexToFelt("0x0669e24364ce0ae7ec2864fb03eedbe60cfbc9d1c74438d10fa4b86552907d54")
	amount, _ := utils.HexToU256Felt("0xffffffff")

	// Here we are declaring the invoke data for the transaction.
	// It's a call to the `mint` function in the `simpleERC20` contract, with the amount of `0xffffffff`.
	invokeData := &pm.UserInvoke{
		UserAddress: account.Address,
		Calls: []pm.Call{
			{
				To:       simpleERC20,
				Selector: utils.GetSelectorFromNameFelt("mint"),
				Calldata: amount,
			},
			// we could add more calls to the transaction if we want. They would be executed in the
			// same paymaster transaction.
		},
	}

	STRKContractAddress, _ := utils.HexToFelt("0x04718f5a0Fc34cC1AF16A1cdee98fFB20C31f5cD61D6Ab07201858f4287c938D")

	// Now that we have the invoke data, we will build the transaction by calling the `paymaster_buildTransaction` method.
	builtTxn, err := paymaster.BuildTransaction(context.Background(), &pm.BuildTransactionRequest{
		Transaction: &pm.UserTransaction{
			Type:   pm.UserTxnInvoke, // we are building an `invoke` transaction
			Invoke: invokeData,
		},
		Parameters: &pm.UserParameters{
			Version: pm.UserParamV1, // Leave as is. This is the only version supported by the paymaster for now.

			// Here we specify the fee mode we want to use for the transaction.
			// We won't spend any value here; this step will just return a fee estimate based on our options.
			FeeMode: pm.FeeMode{
				// There are 2 fee modes supported by the paymaster: `sponsored` and `default`.
				// - `sponsored` fee mode is when an entity will cover your transaction fees. You need an API
				// key from an entity to use this mode.
				// - `default` fee mode is when you cover the fees yourself for the transaction using one of the supported tokens.
				Mode:     pm.FeeModeDefault,
				GasToken: STRKContractAddress, // For the `default` fee mode, use the `gas_token` field
				// to specify which token to use for the fees.

				// There's also the `tip` field to specify the tip for the transaction.
				// - `tip` field is used to specify a tip priority.
				// - `custom` field is used to specify a custom tip value.
				Tip: &pm.TipPriority{
					// Custom: 0, // You can use the `custom` field to specify a custom tip value.

					// Or, you can use the `priority` field to specify a tip priority mode.
					// There are 3 tip priority modes supported by the paymaster: `slow`, `normal` and `fast`.
					Priority: pm.TipPriorityNormal,

					// If you don't specify a tip priority or a custom tip value (`Tip: nil`),
					// the paymaster will use the `normal` tip priority by default.
				},
			},
		},
	})
	if err != nil {
		panic(fmt.Sprintf("Error building the transaction: %s", err))
	}
	fmt.Println("Transaction successfully built by the paymaster")

	// NOTE: Now that we have the built transaction, is up to you to check the fee estimate and
	// decide if you want to proceed with the transaction.
	// The fee estimate is contained in the `fee` JSON field of the built transaction, and looks like this:
	// `{
	// 	"gas_token_price_in_strk": "0xde0b6b3a7640000",
	// 	"estimated_fee_in_strk": "0xd0867e191fcc0",
	// 	"estimated_fee_in_gas_token": "0x83de54ac3228a",
	// 	"suggested_max_fee_in_strk": "0x4e326f496bec80",
	// 	"suggested_max_fee_in_gas_token": "0xab48f32cd750"
	// }`
	PrettyPrint(builtTxn)

	fmt.Println("Step 2: Sign the transaction")

	// Now that we have the built transaction, we need to sign it.
	// The signing process consists of signing the SNIP-12 typed data contained in the built transaction.

	// Firstly, get the message hash of the typed data using our account address as input.
	messageHash, err := builtTxn.TypedData.GetMessageHash(account.Address.String())
	if err != nil {
		panic(fmt.Sprintf("Error getting the message hash of the typed data: %s", err))
	}
	fmt.Println("Message hash of the typed data:", messageHash)

	// Now, we sign the message hash using our account.
	signature, err := account.Sign(context.Background(), messageHash)
	// r, s, err := curve.SignFelts(messageHash, privateKeyFelt) // You can also use the `curve` package to sign the message hash.
	if err != nil {
		panic(fmt.Sprintf("Error signing the transaction: %s", err))
	}
	fmt.Println("Transaction successfully signed")
	PrettyPrint(signature)

	fmt.Println("Step 3: Send the signed transaction")

	// Now that we have the signature, we can send our signed transaction to the paymaster by calling the `paymaster_executeTransaction` method.
	// NOTE: this is the final step, the transaction will be executed and the fees will be paid by us in the specified gas token.
	response, err := paymaster.ExecuteTransaction(context.Background(), &pm.ExecuteTransactionRequest{
		Transaction: &pm.ExecutableUserTransaction{
			Type: pm.UserTxnInvoke,
			Invoke: &pm.ExecutableUserInvoke{
				UserAddress: account.Address,    // Our account address
				TypedData:   builtTxn.TypedData, // The typed data returned by the `paymaster_buildTransaction` method
				Signature:   signature,          // The signature of the message hash made in the previous step
			},
		},
		Parameters: &pm.UserParameters{
			Version: pm.UserParamV1,

			// Using the same fee options as in the `paymaster_buildTransaction` method. A different fee mode here
			// will result in a different fee cost than the one we got in the build step.
			FeeMode: pm.FeeMode{
				Mode:     pm.FeeModeDefault,
				GasToken: STRKContractAddress,
			},
		},
	})
	if err != nil {
		panic(fmt.Sprintf("Error executing the txn with the paymaster: %s", err))
	}

	fmt.Println("Transaction successfully executed by the paymaster")
	fmt.Println("Tracking ID:", response.TrackingId)
	fmt.Println("Transaction Hash:", response.TransactionHash)
}

// PrettyPrint marshals the data with indentation and prints it.
func PrettyPrint(data interface{}) {
	prettyJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(prettyJSON))
	fmt.Println("--------------------------------")
}

// Just a helper function to instantiate the account for us.
func NewAccount(client *rpc.Provider, accountAddress, privateKey, publicKey string, accountCairoVersion account.CairoVersion) *account.Account {
	// Initialise the account memkeyStore (set public and private keys)
	ks := account.NewMemKeystore()
	privKeyBI, ok := new(big.Int).SetString(privateKey, 0)
	if !ok {
		panic("Failed to convert privKey to bigInt")
	}
	ks.Put(publicKey, privKeyBI)

	// Here we are converting the account address to felt
	accountAddressInFelt, err := utils.HexToFelt(accountAddress)
	if err != nil {
		fmt.Println("Failed to transform the account address, did you give the hex address?")
		panic(err)
	}
	// Initialise the account
	accnt, err := account.NewAccount(client, accountAddressInFelt, publicKey, ks, accountCairoVersion)
	if err != nil {
		panic(err)
	}

	return accnt
}
