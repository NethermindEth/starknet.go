package xsessions

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	_ "embed"

	"github.com/NethermindEth/juno/core/felt"
	starknetgo "github.com/NethermindEth/starknet.go"
	"github.com/NethermindEth/starknet.go/artifacts"
	"github.com/NethermindEth/starknet.go/types"
	"github.com/NethermindEth/starknet.go/utils"
)

var sessionPluginCompiled = artifacts.PluginV0Compiled

// sessionToken generates a session key token using the provided private key, account address, and session public key.
//
// Parameters:
// - privateKey: The private key used to sign the token.
// - accountAddress: The address of the account associated with the session.
// - sessionPublicKey: The public key of the session.
//
// Returns:
// - *SessionKeyToken: The generated session key token.
func sessionToken(privateKey, accountAddress, sessionPublicKey string) *SessionKeyToken {
	token, _ := SignToken(
		privateKey,
		types.UTF8StrToBig("SN_GOERLI").Text(16),
		sessionPublicKey,
		accountAddress,
		2*time.Hour,
		[]Policy{{ContractAddress: counterAddress, Selector: "increment"}},
	)
	return token
}

// TestSessionKey_RegisterPlugin is a test function that registers a plugin for the session key.
//
// It generates a plugin hash using the RegisterClass function and creates a new instance of the accountPlugin struct.
// The plugin hash is then set in the accountPlugin struct.
//
// The accountPlugin instance is then written to the ".sessionkey.json" file.
// If there is an error while writing to the file, the function will fail the test.
func TestSessionKey_RegisterPlugin(t *testing.T) {
	pluginHash := RegisterClass(t, sessionPluginCompiled)
	v := &accountPlugin{
		PluginHash: pluginHash,
	}
	err := v.Write(".sessionkey.json")
	if err != nil {
		t.Fatal("should be able to save pluginHash, instead:", err)
	}
}

// TestSessionKey_DeployAccount is a test function that deploys an account and updates the session key.
//
// The function takes no parameters.
// It does not return anything.
func TestSessionKey_DeployAccount(t *testing.T) {
	pk, ok := big.NewInt(0).SetString(privateKey, 0)
	if !ok {
		t.Fatal("could not match *big.Int private key with current value")
	}
	publicKey, _, err := starknetgo.Curve.PrivateToPoint(pk)
	if err != nil {
		t.Fatal(err)
	}
	publicKeyString := fmt.Sprintf("0x%s", publicKey.Text(16))
	v := &accountPlugin{}
	err = v.Read(".sessionkey.json")
	if err != nil {
		t.Fatal(err)
	}
	inputs := []string{
		publicKeyString,
		v.PluginHash,
	}
	accountAddress := DeployContract(t, accountCompiled, inputs)
	v.AccountAddress = accountAddress
	err = v.Write(".sessionkey.json")
	if err != nil {
		t.Fatal(err)
	}
}

// TestSessionKey_MintEth is a test function that verifies the MintEth method of the SessionKey type.
//
// This function reads the ".sessionkey.json" file using the accountPlugin type, and then calls the MintEth method with the given testing.T object and the AccountAddress field of the accountPlugin instance.
func TestSessionKey_MintEth(t *testing.T) {
	v := &accountPlugin{}
	err := v.Read(".sessionkey.json")
	if err != nil {
		t.Fatal(err)
	}
	MintEth(t, v.AccountAddress)
}

// TestSessionKey_CheckEth is a test function that checks the Ethereum session key.
//
// It reads the session key from the ".sessionkey.json" file and checks if there are any errors.
// If there is an error, it fails the test.
// Finally, it calls the CheckEth function with the testing.T object and the account address.
func TestSessionKey_CheckEth(t *testing.T) {
	v := &accountPlugin{}
	err := v.Read(".sessionkey.json")
	if err != nil {
		t.Fatal(err)
	}
	CheckEth(t, v.AccountAddress)
}

// IncrementWithSessionKeyPlugin is a function that increments a counter on a contract using a session key plugin.
//
// It takes the following parameters:
// - t: a testing object for running tests and reporting failures.
// - accountAddress: the address of the account to be used for the transaction.
// - pluginClass: the class of the session key plugin to be used.
// - token: a session key token.
// - counterAddress: the address of the counter contract.
func IncrementWithSessionKeyPlugin(t *testing.T, accountAddress string, pluginClass string, token *SessionKeyToken, counterAddress string) {
	provider := beforeEachRPC(t)
	// shim a keystore into existing tests.
	// use a string representation of the PK as a fake sender address for the keystore
	ks := starknetgo.NewMemKeystore()

	fakeSenderAddress := sessionPrivateKey
	k := types.SNValToBN(sessionPrivateKey)
	ks.Put(fakeSenderAddress, k)
	account, err := starknetgo.NewRPCAccount(
		utils.TestHexToFelt(t, fakeSenderAddress),
		utils.TestHexToFelt(t, accountAddress),
		ks,
		provider,
		WithSessionKeyPlugin(
			pluginClass,
			token,
		))
	if err != nil {
		t.Fatal("deploy should succeed, instead:", err)
	}
	calls := []types.FunctionCall{
		{
			ContractAddress:    utils.TestHexToFelt(t, counterAddress),
			EntryPointSelector: types.GetSelectorFromNameFelt("increment"),
			Calldata:           []*felt.Felt{},
		},
	}
	ctx := context.Background()
	tx, err := account.Execute(ctx, calls, types.ExecuteDetails{})
	if err != nil {
		t.Fatal("execute should succeed, instead:", err)
	}
	if !strings.HasPrefix(tx.TransactionHash.String(), "0x") {
		t.Fatal("execute should return transaction hash, instead:", tx.TransactionHash)
	}
	status, err := provider.WaitForTransaction(ctx, tx.TransactionHash, 8*time.Second)
	if err != nil {
		t.Fatal("declare should succeed, instead:", err)
	}
	if status != types.TransactionAcceptedOnL2 {
		t.Log("unexpected status transaction status, check:", status)
		t.Log("...")
		t.Log("   verify transaction")
		t.Log("...")
		t.Log("export STARKNET_WALLET=starkware.starknet.wallets.open_zeppelin.OpenZeppelinAccount")
		t.Log("export STARKNET_NETWORK=alpha-goerli")
		t.Logf("export HASH=%s\n", tx.TransactionHash)
		t.Log("starknet get_transaction --hash $HASH --feeder_gateway http://localhost:5050/feeder_gateway")
		t.Log("...")
		t.Fail()
	}
	fmt.Printf("tx hash: %s\n", tx.TransactionHash)
}

// TestCounter_IncrementWithSessionKeyPlugin tests the IncrementWithSessionKeyPlugin function in the Counter package.
//
// This function verifies the functionality of the IncrementWithSessionKeyPlugin function by performing a series of steps:
// 1. Reads the account plugin using the provided session key file path.
// 2. Checks for any errors during the reading process. If an error occurs, the test fails.
// 3. Converts the session private key to a *big.Int value.
// 4. Checks if the conversion was successful. If not, the test fails.
// 5. Converts the session private key to a session public key using starknetgo.Curve.PrivateToPoint function.
// 6. Checks for any errors during the conversion process. If an error occurs, the test fails.
// 7. Formats the session public key as a hexadecimal string.
// 8. Generates a session token using the private key, account address, and session public key.
// 9. Calls the IncrementWithSessionKeyPlugin function with the provided parameters.
//
// Parameters:
// - t: A testing.T object used for reporting test failures and logging.
//
// Return Type: void.
func TestCounter_IncrementWithSessionKeyPlugin(t *testing.T) {
	v := &accountPlugin{}
	err := v.Read(".sessionkey.json")
	if err != nil {
		t.Fatal(err)
	}
	sessionPrivateKeyInt, ok := big.NewInt(0).SetString(sessionPrivateKey, 0)
	if !ok {
		t.Fatal("could not match *big.Int private key with current value")
	}
	sessionPublicKeyInt, _, err := starknetgo.Curve.PrivateToPoint(sessionPrivateKeyInt)
	if err != nil {
		t.Fatal(err)
	}
	sessionPublicKey := fmt.Sprintf("0x%s", sessionPublicKeyInt.Text(16))
	token := sessionToken(privateKey, v.AccountAddress, sessionPublicKey)
	IncrementWithSessionKeyPlugin(t, v.AccountAddress, v.PluginHash, token, counterAddress)
}
