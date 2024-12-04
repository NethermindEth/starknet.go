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
// https://sepolia.voyager.online/tx/0x9bc6f6352663aafd71a9ebe1bde9c042590d8f3c8c265e5826274708cf0133

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
		panic("Fail to convert privKey to bitInt")
	}
	ks.Put(publicKey, privKeyBI)

	fmt.Println("Established connection with the client")

	// Initialize the account
	accnt, err := account.NewAccount(client, accountAddressInFelt, publicKey, ks, accountCairoVersion)
	if err != nil {
		panic(err)
	}

	// Get the accounts nonce
	nonce, err := accnt.Nonce(context.Background(), rpc.BlockID{Tag: "latest"}, accnt.AccountAddress)
	if err != nil {
		setup.PanicRPC(err)
	}

	// Build the InvokeTx struct
	InvokeTx := rpc.BroadcastInvokev1Txn{
		InvokeTxnV1: rpc.InvokeTxnV1{
			MaxFee:        new(felt.Felt).SetUint64(100000000000000),
			Version:       rpc.TransactionV1,
			Nonce:         nonce,
			Type:          rpc.TransactionType_Invoke,
			SenderAddress: accnt.AccountAddress,
		}}

	// Convert the contractAddress from hex to felt
	contractAddress, err := utils.HexToFelt(UDCAddress)
	if err != nil {
		panic(err)
	}

	// Build the functionCall struct, where :
	FnCall := rpc.FunctionCall{
		ContractAddress:    contractAddress,                               //contractAddress is the contract that we want to call
		EntryPointSelector: utils.GetSelectorFromNameFelt(contractMethod), //this is the function that we want to call
		Calldata:           getUDCCalldata(accountAddress),                //change this function content to your use case
	}

	// Building the Calldata with the help of FmtCalldata where we pass in the FnCall struct along with the Cairo version
	InvokeTx.Calldata, err = accnt.FmtCalldata([]rpc.FunctionCall{FnCall})
	if err != nil {
		panic(err)
	}

	// Sign the transaction
	err = accnt.SignInvokeTransaction(context.Background(), &InvokeTx.InvokeTxnV1)
	if err != nil {
		panic(err)
	}

	// Estimate the transaction fee
	feeRes, err := accnt.EstimateFee(context.Background(), []rpc.BroadcastTxn{InvokeTx}, []rpc.SimulationFlag{}, rpc.WithBlockTag("latest"))
	if err != nil {
		setup.PanicRPC(err)
	}
	estimatedFee := feeRes[0].OverallFee
	// If the estimated fee is higher than the current fee, let's override it and sign again
	if estimatedFee.Cmp(InvokeTx.MaxFee) == 1 {
		newFee, err := strconv.ParseUint(estimatedFee.String(), 0, 64)
		if err != nil {
			panic(err)
		}
		InvokeTx.MaxFee = new(felt.Felt).SetUint64(newFee + newFee/5) // fee + 20% to be sure
		// Signing the transaction again
		err = accnt.SignInvokeTransaction(context.Background(), &InvokeTx.InvokeTxnV1)
		if err != nil {
			panic(err)
		}
	}

	// After the signing we finally call the AddInvokeTransaction in order to invoke the contract function
	resp, err := accnt.SendTransaction(context.Background(), InvokeTx)
	if err != nil {
		setup.PanicRPC(err)
	}

	fmt.Println("Waiting for the transaction status...")
	time.Sleep(time.Second * 3) // Waiting 3 seconds

	//Getting the transaction status
	txStatus, err := client.GetTransactionStatus(context.Background(), resp.TransactionHash)
	if err != nil {
		setup.PanicRPC(err)
	}

	// This returns us with the transaction hash and status
	fmt.Printf("Transaction hash response: %v\n", resp.TransactionHash)
	fmt.Printf("Transaction execution status: %s\n", txStatus.ExecutionStatus)
	fmt.Printf("Transaction status: %s\n", txStatus.FinalityStatus)
}

// getUDCCalldata is a simple helper to set the call data required by the UDCs deployContract function. Update as needed.
func getUDCCalldata(data ...string) []*felt.Felt {

	classHash, err := new(felt.Felt).SetString(someContractHash)
	if err != nil {
		panic(err)
	}

	randomInt := rand.Uint64()
	salt := new(felt.Felt).SetUint64(randomInt) // to prevent address clashes

	unique, err := new(felt.Felt).SetString("0x0") // see https://docs.starknet.io/architecture-and-concepts/accounts/universal-deployer/#deployment_types
	if err != nil {
		panic(err)
	}

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

	length := int64(len(calldata))
	calldataLen, err := new(felt.Felt).SetString(strconv.FormatInt(length, 16))
	if err != nil {
		panic(err)
	}

	return append([]*felt.Felt{classHash, salt, unique, calldataLen}, calldata...)
}
