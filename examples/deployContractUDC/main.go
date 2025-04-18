package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"math/rand"
	"strconv"
	"time"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"

	setup "github.com/NethermindEth/starknet.go/examples/internal"
)

// More info: https://docs.starknet.io/architecture-and-concepts/accounts/universal-deployer/
// NOTE : Please add in your keys only for testing purposes, in case of a leak you would potentially lose your funds.
var (
	someContractHash string = "0x046ded64ae2dead6448e247234bab192a9c483644395b66f2155f2614e5804b0" // The contract hash to be deployed (in this example, it's an ERC20 contract)
	UDCAddress       string = "0x041a78e741e5af2fec34b695679bc6891742439f7afb8484ecd7766661ad02bf" // UDC contract address
	contractMethod   string = "deployContract"                                                     // UDC method to deploy a contract (from pre-declared contracts)
)

// Example succesful transaction created from this example on Sepolia
// https://sepolia.voyager.online/tx/0xa9a67a7cd8d218bd225335ea2ad4ea4d4c906a5806f14603c7036bfe49ca92

func main() {
	fmt.Println("Starting deployContractUDC example")

	// Load variables from '.env' file
	rpcProviderUrl := setup.GetRpcProviderUrl()
	accountAddress := setup.GetAccountAddress()
	accountCairoVersion := setup.GetAccountCairoVersion()
	privateKey := setup.GetPrivateKey()
	publicKey := setup.GetPublicKey()

	// Initialize connection to RPC provider
	client, err := rpc.NewProvider(rpcProviderUrl)
	if err != nil {
		panic(fmt.Sprintf("Error dialing the RPC provider: %s", err))
	}

	// Here we are converting the account address to felt
	accountAddressInFelt, err := utils.HexToFelt(accountAddress)
	if err != nil {
		panic(err)
	}

	// Initialize the account memkeyStore (set public and private keys)
	ks := account.NewMemKeystore()
	privKeyBI, ok := new(big.Int).SetString(privateKey, 0)
	if !ok {
		panic("Failed to convert privKey to bigInt")
	}
	ks.Put(publicKey, privKeyBI)

	fmt.Println("Established connection with the client")

	// Initialize the account
	accnt, err := account.NewAccount(client, accountAddressInFelt, publicKey, ks, accountCairoVersion)
	if err != nil {
		panic(err)
	}

	// Convert the contractAddress from hex to felt
	contractAddress, err := utils.HexToFelt(UDCAddress)
	if err != nil {
		panic(err)
	}

	// Build the functionCall struct, where :
	FnCall := rpc.InvokeFunctionCall{
		ContractAddress: contractAddress,                //contractAddress is the contract that we want to call
		FunctionName:    contractMethod,                 //this is the function that we want to call
		CallData:        getUDCCalldata(accountAddress), //change this function content to your use case
	}

	// After the signing we finally call the AddInvokeTransaction in order to invoke the contract function
	resp, err := accnt.BuildAndSendInvokeTxn(context.Background(), []rpc.InvokeFunctionCall{FnCall}, 1.5)
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

// getUDCCalldata is a simple helper to set the call data required by the UDCs deployContract function. Update as needed.
func getUDCCalldata(data ...string) []*felt.Felt {

	classHash, err := utils.HexToFelt(someContractHash)
	if err != nil {
		panic(err)
	}

	salt := new(felt.Felt).SetUint64(rand.Uint64()) // to prevent address clashes

	unique := felt.Zero // see https://docs.starknet.io/architecture-and-concepts/accounts/universal-deployer/#deployment_types

	// As we are using an ERC20 token in this example, the calldata needs to have the ERC20 constructor required parameters.
	// You must adjust these fields to match the constructor's parameters of your desired contract.
	// https://docs.openzeppelin.com/contracts-cairo/0.8.1/api/erc20#ERC20-constructor-section
	calldata, err := utils.HexArrToFelt([]string{
		hex.EncodeToString([]byte("MyERC20Token")), //name
		hex.EncodeToString([]byte("MET")),          //symbol
		strconv.FormatInt(200000000000000000, 16),  //fixed_supply (u128 low). See https://book.cairo-lang.org/ch02-02-data-types.html#integer-types
		strconv.FormatInt(0, 16),                   //fixed_supply (u128 high)
		data[0],                                    //recipient
	})
	if err != nil {
		panic(err)
	}

	calldataLen := new(felt.Felt).SetUint64(uint64(len(calldata)))

	return append([]*felt.Felt{classHash, salt, &unique, calldataLen}, calldata...)
}
