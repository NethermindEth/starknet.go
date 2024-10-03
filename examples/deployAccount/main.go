package main

import (
	"context"
	"fmt"
	"math"
	"strconv"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"

	setup "github.com/NethermindEth/starknet.go/examples/internal"
)

var (
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

	// Create transaction data
	tx := rpc.BroadcastDeployAccountTxn{
		DeployAccountTxn: rpc.DeployAccountTxn{
			Nonce:               &felt.Zero, // Contract accounts start with nonce zero.
			MaxFee:              new(felt.Felt).SetUint64(7268996239700),
			Type:                rpc.TransactionType_DeployAccount,
			Version:             rpc.TransactionV1,
			Signature:           []*felt.Felt{},
			ClassHash:           classHash,
			ContractAddressSalt: pub,
			ConstructorCalldata: []*felt.Felt{pub},
		},
	}

	precomputedAddress, err := accnt.PrecomputeAccountAddress(pub, classHash, tx.ConstructorCalldata)
	if err != nil {
		panic(err)
	}
	fmt.Println("PrecomputedAddress:", setup.PadZerosInFelt(precomputedAddress))
	// Sign the transaction
	err = accnt.SignDeployAccountTransaction(context.Background(), &tx.DeployAccountTxn, precomputedAddress)
	if err != nil {
		panic(err)
	}

	// Estimate the transaction fee
	feeRes, err := accnt.EstimateFee(context.Background(), []rpc.BroadcastTxn{tx}, []rpc.SimulationFlag{}, rpc.WithBlockTag("latest"))
	if err != nil {
		setup.PanicRPC(err)
	}
	estimatedFee := feeRes[0].OverallFee
	var feeInETH float64
	// If the estimated fee is higher than the current fee, let's override it and sign again
	if estimatedFee.Cmp(tx.MaxFee) == 1 {
		newFee, err := strconv.ParseUint(estimatedFee.String(), 0, 64)
		if err != nil {
			panic(err)
		}
		tx.MaxFee = new(felt.Felt).SetUint64(newFee + newFee/5) // fee + 20% to be sure
		// Signing the transaction again
		err = accnt.SignDeployAccountTransaction(context.Background(), &tx.DeployAccountTxn, precomputedAddress)
		if err != nil {
			panic(err)
		}
		feeInETH, _ = utils.FeltToBigInt(tx.MaxFee).Float64()
	} else {
		feeInETH, _ = utils.FeltToBigInt(estimatedFee).Float64()
		feeInETH += feeInETH / 5 // fee + 20% to be sure
	}
	//converts fee value from WEI to ETH
	feeInETH = feeInETH / (math.Pow(10, 18))

	// At this point you need to add funds to precomputed address to use it.
	var input string

	fmt.Println("\nThe `precomputedAddress` account needs to have enough ETH to perform a transaction.")
	fmt.Printf("Use the starknet faucet to send ETH to your `precomputedAddress`. You need aproximately %f ETH. \n", feeInETH)
	fmt.Println("When your account has been funded by the faucet, press any key, then `enter` to continue : ")
	fmt.Scan(&input)

	// Send transaction to the network
	resp, err := accnt.SendTransaction(context.Background(), tx)
	if err != nil {
		fmt.Println("Error returned from AddDeployAccountTransaction: ")
		setup.PanicRPC(err)
	}

	fmt.Println("AddDeployAccountTransaction successfully submitted! Wait a few minutes to see it in Voyager.")
	fmt.Printf("Transaction hash: %v \n", resp.TransactionHash)
	fmt.Printf("Contract address: %v \n", resp.ContractAddress)
}
