package main

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/account"
	setup "github.com/NethermindEth/starknet.go/examples/internal"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"
)

// More info: https://docs.starknet.io/architecture-and-concepts/accounts/universal-deployer/
// NOTE : Please add in your keys only for testing purposes, in case of a leak you would potentially lose your funds.
var (
	// This is the class hash of a modern (Sierra) OpenZeppelin ERC20 contract.
	erc20ContractHash, _ = utils.HexToFelt(
		"0x073d71c37e20c569186445d2c497d2195b4c0be9a255d72dbad86662fcc63ae6",
	)
)

// Example successful transaction created from this example on Sepolia
// https://sepolia.starkscan.co/tx/0x6f70bc3756087f02fb3c281f7895520ba87c38152f87f43e0afa595d026469b
func main() {
	fmt.Println("Starting deployContractUDC example")

	// Load variables from '.env' file
	rpcProviderURL := setup.GetRPCProviderURL()
	accountAddress := setup.GetAccountAddress()
	accountCairoVersion := setup.GetAccountCairoVersion()
	privateKey := setup.GetPrivateKey()
	publicKey := setup.GetPublicKey()

	// Initialise connection to RPC provider
	client, err := rpc.NewProvider(context.Background(), rpcProviderURL)
	if err != nil {
		panic(fmt.Sprintf("Error dialling the RPC provider: %s", err))
	}

	// Here we are converting the account address to felt
	accountAddressInFelt, err := utils.HexToFelt(accountAddress)
	if err != nil {
		panic(err)
	}

	// Initialise the account memkeyStore (set public and private keys)
	ks := account.NewMemKeystore()
	privKeyBI, ok := new(big.Int).SetString(privateKey, 0)
	if !ok {
		panic("Failed to convert privKey to bigInt")
	}
	ks.Put(publicKey, privKeyBI)

	fmt.Println("Established connection with the client")

	// Initialise the account
	accnt, err := account.NewAccount(
		client,
		accountAddressInFelt,
		publicKey,
		ks,
		accountCairoVersion,
	)
	if err != nil {
		panic(err)
	}

	// Build the constructor calldata for the ERC20 contract
	// ref: https://docs.openzeppelin.com/contracts-cairo/1.0.0/api/erc20#ERC20Upgradeable-constructor
	name, err := utils.StringToByteArrFelt("My Test Token")
	if err != nil {
		panic(err)
	}
	symbol, err := utils.StringToByteArrFelt("MTT")
	if err != nil {
		panic(err)
	}
	supply, err := utils.HexToU256Felt("0x200000000000000000")
	if err != nil {
		panic(err)
	}
	recipient := accnt.Address
	owner := accnt.Address

	constructorCalldata := make([]*felt.Felt, 0, 10)
	constructorCalldata = append(constructorCalldata, name...)
	constructorCalldata = append(constructorCalldata, symbol...)
	constructorCalldata = append(constructorCalldata, supply...)
	constructorCalldata = append(constructorCalldata, recipient, owner)

	// Deploy the contract with UDC
	resp, salt, err := accnt.DeployContractWithUDC(
		context.Background(),
		erc20ContractHash,
		constructorCalldata,
		nil,
		nil,
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("Waiting for the transaction status...")

	txReceipt, err := accnt.WaitForTransactionReceipt(context.Background(), resp.Hash, time.Second)
	if err != nil {
		panic(err)
	}

	// This returns us with the transaction hash and status
	fmt.Printf("Transaction hash response: %v\n", resp.Hash)
	fmt.Printf("Transaction execution status: %s\n", txReceipt.ExecutionStatus)
	fmt.Printf("Transaction status: %s\n", txReceipt.FinalityStatus)

	// Compute the contract address
	contractAddress := utils.PrecomputeAddressForUDC(
		erc20ContractHash,
		salt,
		constructorCalldata,
		utils.UDCCairoV0,
		accnt.Address,
	)

	fmt.Printf("Contract deployed address: %s\n", contractAddress)
}
