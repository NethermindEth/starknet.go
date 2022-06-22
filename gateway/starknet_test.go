package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/dontpanicdao/caigo"
	"github.com/dontpanicdao/caigo/types"
)

const (
	FEE_MARGIN         int64 = 115
	SEED               int     = 100000000
	ACCOUNT_CLASS_HASH string  = "0x3e327de1c40540b98d05cbcb13552008e36f0ec8d61d46956d2f9752c294328"
)

var (
	snTest         StarknetTest
	snTransactions map[string][]string
	accountClass   RawContractDefinition

	_, b, _, _             = runtime.Caller(0)
	projectRoot            = strings.TrimRight(filepath.Dir(b), "gateway")
	accountCompiled string = projectRoot + "gateway/contracts/account_class.json"
	proxyTest       string = projectRoot + "gateway/contracts/Proxy.cairo"
	proxyCompiled   string = projectRoot + "gateway/contracts/proxy_compiled.json"
)

type StarknetTest struct {
	Environments    []TestEnvironment `json:"environments"`
	AccountCalldata []string          `json:"accountCalldata"`
}

type TestEnvironment struct {
	Chain    string `json:"chain"`
	Accounts []struct {
		Address      string              `json:"address,omitempty"`
		Public       string              `json:"public_key,omitempty"`
		Private      string              `json:"private_key,omitempty"`
		Transactions []types.Transaction `json:"transactions"`
	} `json:"accounts"`
	ContractAddresses types.ContractAddresses `json:"contractAddresses"`
}

// requires starknet-devnet to be running and accessible and no seed:
// ex: starknet-devnet --port 5000 --seed 0
// (ref: https://github.com/Shard-Labs/starknet-devnet)
func init() {
	testFile, err := os.Open(projectRoot + "gateway/contracts/starknet_test.json")
	if err != nil {
		panic(err.Error())
	}

	defer testFile.Close()

	raw, _ := ioutil.ReadAll(testFile)
	json.Unmarshal(raw, &snTest)
	snTransactions = make(map[string][]string)
	for _, env := range snTest.Environments {
		snTransactions[env.Chain] = []string{}
	}

	if _, err := os.Stat(accountCompiled); os.IsNotExist(err) {
		accountClass, err := NewClient().ClassByHash(context.Background(), ACCOUNT_CLASS_HASH)
		if err != nil {
			panic(err.Error())
		}

		file, err := json.Marshal(accountClass)
		err = ioutil.WriteFile(accountCompiled, file, 0644)
		if err != nil {
			panic(err.Error())
		}
	}
	if _, err := os.Stat(proxyCompiled); os.IsNotExist(err) {
		home, err := os.UserHomeDir()
		if err != nil {
			panic(err.Error())
		}

		err = exec.Command(home+"/cairo_venv/bin/starknet-compile", "--cairo_path", projectRoot+"gateway", proxyTest, "--output", proxyCompiled).Run()
		if err != nil {
			panic(err.Error())
		}
	}
}

func TestExecute(t *testing.T) {
	for _, env := range snTest.Environments {
		// signature scheme for devnet is not curently compatible
		if env.Chain == "testnet" {
			for _, testAccount := range env.Accounts {
				if testAccount.Private != "" {
					curve, err := caigo.SC(caigo.WithConstants(projectRoot + "pedersen_params.json"))
					if err != nil {
						t.Errorf("%s: could not init with constant points: %v\n", env.Chain, err)
					}

					account, err := caigo.NewAccount(&curve, testAccount.Private, testAccount.Address, NewProvider(WithChain(env.Chain)))
					if err != nil {
						t.Errorf("%s: could not create account: %v\n", env.Chain, err)
					}

					feeEstimate, err := account.EstimateFee(context.Background(), testAccount.Transactions)
					if err != nil {
						t.Errorf("%s: could not estimate fee for transaction: %v\n", env.Chain, err)
					}
					fee := new(types.Felt)
					fee.Int = big.NewInt(feeEstimate.Amount * FEE_MARGIN/100)

					txResp, err := account.Execute(context.Background(), fee, testAccount.Transactions)
					if err != nil {
						t.Errorf("Could not execute test transaction: %v\n", err)
					}

					snTransactions[env.Chain] = append(snTransactions[env.Chain], txResp.TransactionHash)
				}
			}
		}
	}
}

func TestDeclareAndDeploy(t *testing.T) {
	for _, env := range snTest.Environments {
		if env.Chain != "mainnet" {
			gw := NewClient(WithChain(env.Chain))
			declareTx, err := gw.Declare(context.Background(), accountCompiled, types.DeclareRequest{})
			if err != nil {
				t.Errorf("%s: could not 'DECLARE' contract: %v\n", env.Chain, err)
				return
			}

			tx, err := gw.Transaction(context.Background(), TransactionOptions{TransactionHash: declareTx.TransactionHash})
			if err != nil {
				t.Errorf("%s: could not get 'DECLARE' transaction: %v\n", env.Chain, err)
			}
			if tx.Transaction.Type != DECLARE {
				t.Errorf("%s: incorrect delcare transaction: %v\n", env.Chain, tx)
			}

			for _, testAccount := range env.Accounts {
				salt := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(SEED)
				calldata := append(snTest.AccountCalldata, caigo.HexToBN(testAccount.Public).String(), "0")

				deployTx, err := gw.Deploy(context.Background(), proxyCompiled, types.DeployRequest{
					ContractAddressSalt: fmt.Sprintf("0x%x", salt),
					ConstructorCalldata: calldata,
				})
				if err != nil {
					t.Errorf("%s: could not deploy contract: %v\n", env.Chain, err)
				}

				snTransactions[env.Chain] = append(snTransactions[env.Chain], deployTx.TransactionHash)
			}
		}
	}
}

func TestContractAddresses(t *testing.T) {
	for _, env := range snTest.Environments {
		if env.Chain != "devnet" {
			gw := NewClient(WithChain(env.Chain))
			addresses, err := gw.ContractAddresses(context.Background())
			if err != nil {
				t.Errorf("%s: could not get starknet addresses - \n%v\n", env.Chain, err)
			}

			if strings.ToLower(addresses.Starknet) != env.ContractAddresses.Starknet {
				t.Errorf("%s: fetched incorrect addresses - \n%s %s\n", env.Chain, strings.ToLower(addresses.Starknet), env.ContractAddresses.Starknet)
			}
		}
	}
}

func TestStateDiff(t *testing.T) {
	for _, env := range snTest.Environments {
		if env.Chain != "devnet" {
			gw := NewClient(WithChain(env.Chain))
			diff, err := gw.StateUpdate(context.Background(), nil)
			if err != nil {
				t.Errorf("%s: could not get starknet addresses - \n%v\n", env.Chain, err)
			}

			if diff.OldRoot == "" || diff.NewRoot == "" {
				t.Errorf("%s: could not fetch accurate state update - \n%s %s\n", env.Chain, diff.OldRoot, diff.NewRoot)
			}
		}
	}
}

func TestCall(t *testing.T) {
	for _, env := range snTest.Environments {
		gw := NewClient(WithChain(env.Chain))
		for _, testAccount := range env.Accounts {
			call := types.FunctionCall{
				ContractAddress: testAccount.Address,
			}
			if env.Chain == "devnet" {
				call.EntryPointSelector = "get_public_key"
			} else {
				call.EntryPointSelector = "get_signer"
			}
			resp, err := gw.Call(context.Background(), call, "")
			if err != nil {
				t.Errorf("%s: could 'Call' deployed contract: %v\n", env.Chain, err)
			}
			if len(resp) == 0 {
				t.Errorf("%s: could get signing key for account: %v\n", env.Chain, err)
			}

			if resp[0] != testAccount.Public {
				t.Errorf("%s: signing key is incorrect: \n%s %v\n", env.Chain, resp[0], testAccount.Public)
			}
		}
	}
}

// func TestTransactions(t *testing.T) {
// 	for k, txs := range snTransactions {
// 		fmt.Println("ENV - ", k)

// 		gw := NewClient(WithChain(k))
// 		for _, tx := range txs {
// 			_, status, err := gw.PollTx(context.Background(), tx, types.ACCEPTED_ON_L2, 5, 150)
// 			if err != nil {
// 				t.Errorf("Bad transaction poll: %v\n", err)
// 			}

// 			transaction, err := gw.Transaction(context.Background(), TransactionOptions{TransactionHash: tx})
// 			if err != nil {
// 				t.Errorf("Bad transaction poll: %v\n", err)
// 			}
// 			fmt.Printf("\tTxHash(%s): %s-%s\n", transaction.Transaction.Type, status, tx)
// 		}
// 	}
// }
