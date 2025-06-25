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
	someContractHash string = "0x073d71c37e20c569186445d2c497d2195b4c0be9a255d72dbad86662fcc63ae6"
)

// Example successful transaction created from this example on Sepolia
// https://sepolia.starkscan.co/tx/0x04d646a167c25530a0e4d1be296885dd14c1a1bf32a39cd6a4ddc4cb5ce1c5b2
func main() {
	fmt.Println("Starting deployContractUDC example")

	// Load variables from '.env' file
	rpcProviderUrl := setup.GetRpcProviderUrl()
	accountAddress := setup.GetAccountAddress()
	accountCairoVersion := setup.GetAccountCairoVersion()
	privateKey := setup.GetPrivateKey()
	publicKey := setup.GetPublicKey()

	// Initialise connection to RPC provider
	client, err := rpc.NewProvider(rpcProviderUrl)
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
	accnt, err := account.NewAccount(client, accountAddressInFelt, publicKey, ks, accountCairoVersion)
	if err != nil {
		panic(err)
	}

	classHash, err := utils.HexToFelt(someContractHash)
	if err != nil {
		panic(err)
	}

	// For modern contracts, strings are represented as `ByteArray` which serialize into multiple felts.
	nameAsFelts, err := utils.StringToByteArrFelt("My Test Token")
	if err != nil {
		panic(err)
	}
	symbolAsFelts, err := utils.StringToByteArrFelt("MTT")
	if err != nil {
		panic(err)
	}
	recipient, err := utils.HexToFelt(accountAddress)
	if err != nil {
		panic(err)
	}

	// Assemble the constructor calldata.
	constructorCalldata := nameAsFelts
	constructorCalldata = append(constructorCalldata, symbolAsFelts...)
	constructorCalldata = append(constructorCalldata,
		// u256 supply: 1000 tokens with 18 decimals (10^21)
		new(felt.Felt).SetBigInt(new(big.Int).And(new(big.Int).Exp(big.NewInt(10), big.NewInt(21), nil), new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 128), big.NewInt(1)))), // low
		new(felt.Felt).SetBigInt(new(big.Int).Rsh(new(big.Int).Exp(big.NewInt(10), big.NewInt(21), nil), 128)),                                                                   // high
		recipient, // recipient
		recipient, // owner
	)

	resp, err := accnt.DeployContractWithUDC(context.Background(), classHash, constructorCalldata, nil, nil)
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
}
