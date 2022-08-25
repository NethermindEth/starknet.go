package main

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/dontpanicdao/caigo"
	"github.com/dontpanicdao/caigo/gateway"
	"github.com/dontpanicdao/caigo/types"
)

const (
	env                   = "testnet"
	compiledOZAccount     = "./contracts/account/OZAccount_compiled.json"
	compiledERC20Contract = "./contracts/erc20/erc20_custom_compiled.json"
	predeployedContract   = "0x0024e9f35c5d6a14dcbb3f08be5fb7703e76611767266068c521fe8cba27983c"
	maxPoll               = 15
	pollInterval          = 5
)

func main() {
	// init starknet gateway client
	gw := gateway.NewProvider(gateway.WithChain(env))

	// Generating priv / pub key pair
	privateKey, err := caigo.Curve.GetRandomPrivateKey()
	if err != nil {
		fmt.Println("can't get random private key :", err)
		os.Exit(1)
	}
	pubX, _, err := caigo.Curve.PrivateToPoint(privateKey)
	if err != nil {
		fmt.Println("can't generate public key :", err)
		os.Exit(1)
	}

	// Deploying OpenZeppelin account
	fmt.Println("Deploying account to testnet. It may take a while.")
	accoutResponse, err := gw.Deploy(context.Background(), compiledOZAccount, types.DeployRequest{
		Type:                gateway.DEPLOY,
		ConstructorCalldata: []string{pubX.String()},                            // public key
		ContractAddressSalt: caigo.BigToHex(big.NewInt(time.Now().UnixNano()))}) // salt to hex
	if err != nil {
		fmt.Println("can't deploy account :", err)
		os.Exit(1)
	}

	// poll until the transaction is accepted on L2
	if err := waitForTransaction(gw, accoutResponse.TransactionHash); err != nil {
		fmt.Println("Account deployement transaction failure :", err)
		os.Exit(1)
	}

	// fetch transaction data
	tx, err := gw.Transaction(context.Background(), gateway.TransactionOptions{TransactionHash: accoutResponse.TransactionHash})
	if err != nil {
		fmt.Println("can't fetch transaction data :", err)
		os.Exit(1)
	}

	// Instantiate account object
	account, err := caigo.NewAccount(privateKey.String(), tx.Transaction.ContractAddress, gw)
	if err != nil {
		fmt.Println("can't create account:", err)
		os.Exit(1)
	}

	fmt.Println("Account deployed. contract address : ", account.Address)

	// At this point you need to add fund to the deployed account in order to use it.
	var input string
	fmt.Println("The deployed account has to be feeded with ETH to perform transaction.")
	fmt.Print("When your account has been funded with the faucet, press any key and enter to continue : ")
	fmt.Scan(&input)

	fmt.Println("Deploying erc20 contract. It may take a while")
	erc20Response, err := gw.Deploy(context.Background(), compiledERC20Contract, types.DeployRequest{
		Type:                gateway.DEPLOY,
		ContractAddressSalt: caigo.BigToHex(big.NewInt(time.Now().UnixNano())), // salt to hex
		ConstructorCalldata: []string{
			caigo.HexToBN(account.Address).String(), // owner
			"2000",                                  // initial supply
			"0",                                     // Uint256 additionnal parameter
		},
	})
	if err != nil {
		fmt.Println("can't deploy erc20 contract :", err)
		os.Exit(1)
	}

	// Poll until the transaction is accepted on L2
	if err := waitForTransaction(gw, erc20Response.TransactionHash); err != nil {
		fmt.Println("ERC20 deployement transaction failure :", err)
		os.Exit(1)
	}

	// fetch transaction data
	txERC20, err := gw.Transaction(context.Background(), gateway.TransactionOptions{TransactionHash: erc20Response.TransactionHash})
	if err != nil {
		fmt.Println("can't fetch transaction data :", err)
		os.Exit(1)
	}
	fmt.Println("ERC20 contract deployed.",
		"Contract address : ", txERC20.Transaction.ContractAddress,
		"Transaction hash : ", txERC20.Transaction.TransactionHash,
	)

	erc20ContractAddr := txERC20.Transaction.ContractAddress

	// minting the erc20 contract
	fmt.Println("Minting 10 tokens to your account...")
	if err := mint(gw, account, erc20ContractAddr); err != nil {
		fmt.Println("can't mint erc20 contract :", err)
		os.Exit(1)
	}

	balance, err := getBalanceOf(gw, erc20ContractAddr, account.Address)
	if err != nil {
		fmt.Println("can't get balance of :", account.Address, err)
		os.Exit(1)
	}
	fmt.Println("Your account has ", balance, " tokens.")

	// Make transfer from transacation
	fmt.Println("Transfering 5 tokens from", account.Address, "to", predeployedContract)
	if err := transferFrom(gw, account, erc20ContractAddr, predeployedContract); err != nil {
		fmt.Println("can't transfer tokens  :", account.Address, err)
		os.Exit(1)
	}

	balanceAccount, err := getBalanceOf(gw, erc20ContractAddr, account.Address)
	if err != nil {
		fmt.Println("can't get balance of :", account.Address, err)
		os.Exit(1)
	}
	balancePredeployed, err := getBalanceOf(gw, erc20ContractAddr, account.Address)
	if err != nil {
		fmt.Println("can't get balance of :", predeployedContract, err)
		os.Exit(1)
	}

	fmt.Println("Transfer done.")
	fmt.Println("Account balance : ", balanceAccount, ". Predeployed account balance : ", balancePredeployed)
}

// Utils function to wait for transaction to be accepted on L2 and print tx status
func waitForTransaction(gw *gateway.GatewayProvider, transactionHash string) error {
	acceptedOnL2 := false
	var receipt *types.TransactionReceipt
	var err error
	fmt.Println("Polling until transaction is accepted on L2...")
	for !acceptedOnL2 {
		_, receipt, err = gw.PollTx(context.Background(), transactionHash, types.ACCEPTED_ON_L2, pollInterval, maxPoll)
		if err != nil {
			fmt.Println(receipt.Status, receipt.StatusData)
			return fmt.Errorf("Transaction Failure (%s) : can't poll to desired status : %s", transactionHash, err.Error())
		}
		fmt.Println("Current status : ", receipt.Status)
		if receipt.Status == types.ACCEPTED_ON_L2.String() {
			acceptedOnL2 = true
		}
	}
	return nil
}

// mint will mint the erc20 contracts through the account.
func mint(gw *gateway.GatewayProvider, account *caigo.Account, erc20address string) error {
	// Transaction that will be executed by the account contract.
	tx := []types.Transaction{
		{
			ContractAddress:    erc20address,
			EntryPointSelector: "mint",
			Calldata: []string{
				caigo.HexToBN(account.Address).String(), // owner
				"10",                                    // amount to mint
				"0",                                     // UInt256 additional parameter
			},
		},
	}

	execResp, err := account.Execute(context.Background(), tx, caigo.ExecuteDetails{})
	if err != nil {
		return fmt.Errorf("can't execute transacation : %w", err)
	}

	if err := waitForTransaction(gw, execResp.TransactionHash); err != nil {
		return fmt.Errorf("a problem occured with transacation : %w", err)
	}
	return nil
}

// transferFrom will transfer 5 tokens from account balance to the otherAccount by
// calling the transferFrom function of the erc20 contract.
func transferFrom(gw *gateway.GatewayProvider, account *caigo.Account, erc20address, otherAccount string) error {
	// Transaction that will be executed by the account contract.
	tx := []types.Transaction{
		{
			ContractAddress:    erc20address,
			EntryPointSelector: "transferFrom",
			Calldata: []string{
				caigo.HexToBN(account.Address).String(), // sender
				caigo.HexToBN(otherAccount).String(),    // recipient
				"5",                                     // amount to transfer
				"0",                                     // UInt256 additional parameter
			},
		},
	}

	execResp, err := account.Execute(context.Background(), tx, caigo.ExecuteDetails{})
	if err != nil {
		return fmt.Errorf("can't execute transacation : %w", err)
	}

	if err := waitForTransaction(gw, execResp.TransactionHash); err != nil {
		return fmt.Errorf("a problem occured with transacation : %w", err)
	}
	return nil
}

// getBalanceOf returns the balance of the account at the accountAddress address.
func getBalanceOf(gw *gateway.GatewayProvider, erc20address, accountAddress string) (string, error) {
	res, err := gw.Call(context.Background(), types.FunctionCall{
		ContractAddress:    erc20address,
		EntryPointSelector: "balanceOf",
		Calldata: []string{
			caigo.HexToBN(accountAddress).String(),
		},
	}, "")
	if err != nil {
		return "", fmt.Errorf("can't call erc20 : %s. Error : %w", accountAddress, err)
	}
	return res[0], nil
}
