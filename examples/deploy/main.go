package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	starknetgo "github.com/NethermindEth/starknet.go"
	"github.com/NethermindEth/starknet.go/artifacts"
	"github.com/NethermindEth/starknet.go/gateway"
	"github.com/NethermindEth/starknet.go/types"
)

const (
	env                   = "testnet"
	compiledOZAccount     = "./contracts/account/OZAccount_compiled.json"
	compiledERC20Contract = "./contracts/erc20/erc20_custom_compiled.json"
	predeployedContract   = "0x0024e9f35c5d6a14dcbb3f08be5fb7703e76611767266068c521fe8cba27983c"
	maxPoll               = 15
	pollInterval          = 5
)

// main is the entry point of the deploy example.
//
// It initializes a gateway client and generates a random private key. It then deploys an account to the testnet,
// waits for the deployment transaction to complete, fetches the transaction data, creates a gateway account,
// saves the private key, and prompts the user to fund the deployed account. After that, it deploys an ERC20 contract,
// waits for the deployment transaction to complete, fetches the transaction data, and mints tokens to the gateway account.
// Finally, it transfers tokens from the gateway account to a predeployed contract and prints the account balances.
func main() {
	gw := gateway.NewClient(gateway.WithChain(env))

	privateKey, err := starknetgo.Curve.GetRandomPrivateKey()
	if err != nil {
		fmt.Println("can't get random private key:", err)
		os.Exit(1)
	}
	pubX, _, err := starknetgo.Curve.PrivateToPoint(privateKey)
	if err != nil {
		fmt.Println("can't generate public key:", err)
		os.Exit(1)
	}

	contractClass := types.ContractClass{}
	err = json.Unmarshal(artifacts.AccountCompiled, &contractClass)
	if err != nil {
		fmt.Println("could not log file", err)
		os.Exit(1)
	}
	fmt.Println("Deploying account to testnet. It may take a while.")
	accountResponse, err := gw.Deploy(context.Background(), contractClass, types.DeployRequest{
		Type:                gateway.DEPLOY,
		ContractAddressSalt: types.BigToHex(pubX),     // salt to hex
		ConstructorCalldata: []string{pubX.String()}}) // public key
	if err != nil {
		fmt.Println("can't deploy account:", err)
		os.Exit(1)
	}

	if err := waitForTransaction(gw, accountResponse.TransactionHash); err != nil {
		fmt.Println("Account deployement transaction failure:", err)
		os.Exit(1)
	}

	tx, err := gw.Transaction(context.Background(), gateway.TransactionOptions{TransactionHash: accountResponse.TransactionHash})
	if err != nil {
		fmt.Println("can't fetch transaction data:", err)
		os.Exit(1)
	}

	account, err := starknetgo.NewGatewayAccount(privateKey.String(), types.StrToFelt(tx.Transaction.ContractAddress), gw)
	if err != nil {
		fmt.Println("can't create account:", err)
		os.Exit(1)
	}

	fmt.Println("Account deployed. Contract address: ", account.Address)
	if err := savePrivateKey(types.BigToHex(privateKey)); err != nil {
		fmt.Println("can't save private key:", err)
		os.Exit(1)
	}

	// At this point you need to add funds to the deployed account in order to use it.
	var input string
	fmt.Println("The deployed account has to be feeded with ETH to perform transaction.")
	fmt.Print("When your account has been funded with the faucet, press any key and enter to continue : ")
	fmt.Scan(&input)

	fmt.Println("Deploying erc20 contract. It may take a while")
	erc20Response, err := gw.Deploy(context.Background(), compiledERC20Contract, types.DeployRequest{
		Type:                gateway.DEPLOY,
		ContractAddressSalt: types.BigToHex(pubX), // salt to hex
		ConstructorCalldata: []string{
			account.Address.String(), // owner
			"2000",                   // initial supply
			"0",                      // Uint256 additional parameter
		},
	})
	if err != nil {
		fmt.Println("can't deploy erc20 contract:", err)
		os.Exit(1)
	}

	if err := waitForTransaction(gw, erc20Response.TransactionHash); err != nil {
		fmt.Println("ERC20 deployment transaction failure:", err)
		os.Exit(1)
	}

	txERC20, err := gw.Transaction(context.Background(), gateway.TransactionOptions{TransactionHash: erc20Response.TransactionHash})
	if err != nil {
		fmt.Println("can't fetch transaction data:", err)
		os.Exit(1)
	}
	fmt.Println("ERC20 contract deployed.",
		"Contract address: ", txERC20.Transaction.ContractAddress,
		"Transaction hash: ", txERC20.Transaction.TransactionHash,
	)

	erc20ContractAddr := txERC20.Transaction.ContractAddress

	fmt.Println("Minting 10 tokens to your account...")
	if err := mint(gw, account, erc20ContractAddr); err != nil {
		fmt.Println("can't mint erc20 contract:", err)
		os.Exit(1)
	}

	balance, err := balanceOf(gw, erc20ContractAddr, account.Address.Hex())
	if err != nil {
		fmt.Println("can't get balance of:", account.Address, err)
		os.Exit(1)
	}
	fmt.Println("Your account has ", balance, " tokens.")

	fmt.Println("Transferring 5 tokens from", account.Address, "to", predeployedContract)
	if err := transferFrom(gw, account, erc20ContractAddr, predeployedContract); err != nil {
		fmt.Println("can't transfer tokens:", account.Address, err)
		os.Exit(1)
	}

	balanceAccount, err := balanceOf(gw, erc20ContractAddr, account.Address.Hex())
	if err != nil {
		fmt.Println("can't get balance of:", account.Address, err)
		os.Exit(1)
	}
	balancePredeployed, err := balanceOf(gw, erc20ContractAddr, account.Address.Hex())
	if err != nil {
		fmt.Println("can't get balance of:", predeployedContract, err)
		os.Exit(1)
	}

	fmt.Println("Transfer done.")
	fmt.Println("Account balance: ", balanceAccount, ". Predeployed account balance: ", balancePredeployed)
}

// waitForTransaction waits for a transaction to be accepted on Layer 2 of the gateway nd print tx status.
//
// It takes the gateway instance and the transaction hash as parameters.
// It returns an error if there is a transaction failure.
func waitForTransaction(gw *gateway.Gateway, transactionHash string) error {
	acceptedOnL2 := false
	var receipt *gateway.TransactionReceipt
	var err error
	fmt.Println("Polling until transaction is accepted on L2...")
	for !acceptedOnL2 {
		_, receipt, err = gw.PollTx(context.Background(), transactionHash, types.ACCEPTED_ON_L2, pollInterval, maxPoll)
		if err != nil {
			fmt.Println(receipt.Status)
			return fmt.Errorf("Transaction Failure (%s): can't poll to desired status: %s", transactionHash, err.Error())
		}
		fmt.Println("Current status : ", receipt.Status)
		if receipt.Status == types.ACCEPTED_ON_L2.String() {
			acceptedOnL2 = true
		}
	}
	return nil
}

// mint mints a specified amount of ERC20 tokens for a given account on the StarkNet gateway.
//
// Parameters:
// - gw: The instance of the StarkNet gateway.
// - account: The account to mint tokens for.
// - erc20address: The address of the ERC20 contract to interact with.
//
// Returns:
// An error if there was a problem executing the transaction or waiting for it to be confirmed.
func mint(gw *gateway.Gateway, account *starknetgo.Account, erc20address string) error {
	// Transaction that will be executed by the account contract.
	tx := []types.FunctionCall{
		{
			ContractAddress:    types.StrToFelt(erc20address),
			EntryPointSelector: "mint",
			Calldata: []string{
				account.Address.String(), // to
				"10",                     // amount to mint
				"0",                      // UInt256 additional parameter
			},
		},
	}

	execResp, err := account.Execute(context.Background(), tx, types.ExecuteDetails{})
	if err != nil {
		return fmt.Errorf("can't execute transaction: %w", err)
	}

	if err := waitForTransaction(gw, execResp.TransactionHash); err != nil {
		return fmt.Errorf("a problem occured with the transaction: %w", err)
	}
	return nil
}

// transferFrom executes a transaction to transfer a specified amount (5) of tokens from one account to another using the ERC20 transferFrom method.
//
// Parameters:
// - gw: the Gateway instance used to interact with the StarkNet network.
// - account: the account from which the tokens will be transferred.
// - erc20address: the address of the ERC20 contract.
// - otherAccount: the address of the recipient account.
//
// Returns:
// - error: an error if the transaction execution or waiting for the transaction fails.
func transferFrom(gw *gateway.Gateway, account *starknetgo.Account, erc20address, otherAccount string) error {
	// Transaction that will be executed by the account contract.
	tx := []types.FunctionCall{
		{
			ContractAddress:    types.StrToFelt(erc20address),
			EntryPointSelector: "transferFrom",
			Calldata: []string{
				account.Address.String(),             // sender
				types.HexToBN(otherAccount).String(), // recipient
				"5",                                  // amount to transfer
				"0",                                  // UInt256 additional parameter
			},
		},
	}

	execResp, err := account.Execute(context.Background(), tx, types.ExecuteDetails{})
	if err != nil {
		return fmt.Errorf("can't execute transaction: %w", err)
	}

	if err := waitForTransaction(gw, execResp.TransactionHash); err != nil {
		return fmt.Errorf("a problem occured with transaction: %w", err)
	}
	return nil
}

// balanceOf returns the balance of a specific ERC20 token for a given accountAddress.
//
// Parameters:
// - gw: The gateway used to make the function call.
// - erc20address: The address of the ERC20 token contract.
// - accountAddress: The address of the account to check the balance for.
//
// Returns:
// - string: The balance of the ERC20 token for the given account address.
// - error: An error if the function call fails.
func balanceOf(gw *gateway.Gateway, erc20address, accountAddress string) (string, error) {
	res, err := gw.Call(context.Background(), types.FunctionCall{
		ContractAddress:    types.StrToFelt(erc20address),
		EntryPointSelector: "balanceOf",
		Calldata: []string{
			types.HexToBN(accountAddress).String(),
		},
	}, "")
	if err != nil {
		return "", fmt.Errorf("can't call erc20: %s. Error: %w", accountAddress, err)
	}
	low := types.StrToFelt(res[0])
	hi := types.StrToFelt(res[1])

	balance, err := types.NewUint256(low, hi)
	if err != nil {
		return "", nil
	}
	return balance.String(), nil
}

// savePrivateKey saves the private key to a file (private_key.txt).
//
// It takes a string parameter 'privKey' which represents the private key to be saved.
// Returns an error if there is any issue creating or writing to the file.
func savePrivateKey(privKey string) error {
	file, err := os.Create("private_key.txt")
	if err != nil {
		return fmt.Errorf("can't create private_key.txt")
	}
	defer file.Close()
	if _, err := file.WriteString(privKey); err != nil {
		return fmt.Errorf("can't write private_key.txt")
	}
	return nil
}
