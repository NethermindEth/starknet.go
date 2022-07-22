package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/dontpanicdao/caigo"
	"github.com/dontpanicdao/caigo/types"
)

const (
	FEE_MARGIN         uint64 = 115
	SEED               int    = 100000000
	ACCOUNT_CLASS_HASH string = "0x3e327de1c40540b98d05cbcb13552008e36f0ec8d61d46956d2f9752c294328"
)

var (
	devnetAccounts  []TestAccountType
	testnetAccounts []TestAccountType

	_, b, _, _             = runtime.Caller(0)
	projectRoot            = strings.TrimRight(filepath.Dir(b), "gateway")
	accountCompiled string = projectRoot + "gateway/contracts/account_class.json"
	proxyTest       string = projectRoot + "gateway/contracts/Proxy.cairo"
	proxyCompiled   string = projectRoot + "gateway/contracts/proxy.json"
)

type TestAccountType struct {
	PrivateKey   string              `json:"private_key"`
	PublicKey    string              `json:"public_key"`
	Address      string              `json:"address"`
	Transactions []types.Transaction `json:"transactions,omitempty"`
}

// requires starknet-devnet to be running and accessible and no seed:
// ex: starknet-devnet
// (ref: https://github.com/Shard-Labs/starknet-devnet)
func init() {
	if _, err := os.Stat(accountCompiled); os.IsNotExist(err) {
		accountClass, err := NewClient().ClassByHash(context.Background(), ACCOUNT_CLASS_HASH)
		if err != nil {
			panic(err.Error())
		}

		file, err := json.Marshal(accountClass)
		if err != nil {
			panic(err.Error())
		}

		if err = ioutil.WriteFile(accountCompiled, file, 0644); err != nil {
			panic(err.Error())
		}
	}

	var err error
	if devnetAccounts, err = DevnetAccounts(); err != nil {
		panic(err.Error())
	}

	testnetAccounts = []TestAccountType{
		{
			PrivateKey: "0x28a778906e0b5f4d240ad25c5993422e06769eb799483ae602cc3830e3f538",
			PublicKey:  "0x63f0f116c78146e1e4e193923fe3cad5f236c0ed61c2dc04487a733031359b8",
			Address:    "0x0254cfb85c43dee6f410867b9795b5309beb4a2640211c8f5b2c7681a47e5f3c",
			Transactions: []types.Transaction{
				{
					ContractAddress:    "0x22b0f298db2f1776f24cda70f431566d9ef1d0e54a52ee6d930b80ec8c55a62",
					EntryPointSelector: "update_single_store",
					Calldata:           []string{"3"},
				},
			},
		},
		{
			PrivateKey: "0x879d7dad7f9df54e1474ccf572266bba36d40e3202c799d6c477506647c126",
			PublicKey:  "0xb95246e1caeaf34672906d7b74bd6968231a2130f41e85aebb62d43b88068",
			Address:    "0x0126dd900b82c7fc95e8851f9c64d0600992e82657388a48d3c466553d4d9246",
			Transactions: []types.Transaction{
				{
					ContractAddress:    "0x22b0f298db2f1776f24cda70f431566d9ef1d0e54a52ee6d930b80ec8c55a62",
					EntryPointSelector: "update_multi_store",
					Calldata:           []string{"4", "7"},
				},
				{
					ContractAddress:    "0x22b0f298db2f1776f24cda70f431566d9ef1d0e54a52ee6d930b80ec8c55a62",
					EntryPointSelector: "update_struct_store",
					Calldata:           []string{"435921360636", "15000000000000000000", "0"},
				},
			},
		},
	}
}

func TestDeclare(t *testing.T) {
	for _, env := range []string{"devnet", "testnet"} {
		gw := NewClient(WithChain(env))
		declareTx, err := gw.Declare(context.Background(), accountCompiled, types.DeclareRequest{})
		if err != nil {
			t.Errorf("%s: could not 'DECLARE' contract: %v\n", env, err)
			return
		}

		tx, err := gw.Transaction(context.Background(), TransactionOptions{TransactionHash: declareTx.TransactionHash})
		if err != nil {
			t.Errorf("%s: could not get 'DECLARE' transaction: %v\n", env, err)
		}
		if tx.Transaction.Type != DECLARE {
			t.Errorf("%s: incorrect delcare transaction: %v\n", env, tx)
		}
	}
}

func TestExecuteGoerli(t *testing.T) {
	for _, testAccount := range testnetAccounts {
		account, err := caigo.NewAccount(testAccount.PrivateKey, testAccount.Address, NewProvider())
		if err != nil {
			t.Errorf("testnet: could not create account: %v\n", err)
		}

		feeEstimate, err := account.EstimateFee(context.Background(), testAccount.Transactions)
		if err != nil {
			t.Errorf("testnet: could not estimate fee for transaction: %v\n", err)
		}
		fee := new(types.Felt)
		fee.Int = new(big.Int).SetUint64(feeEstimate.Amount * FEE_MARGIN / 100)

		// TODO: fix estimate_fee call
		// _, err = account.Execute(context.Background(), fee, testAccount.Transactions)
		// if err = nil {
		// 	t.Errorf("Could not execute test transaction: %v\n", err)
		// }
	}
}

func TestDeployGoerli(t *testing.T) {
	accountCalldata := []string{
		"622947212727630016888676747031248582302218969620261661169319263633873449798",
		"215307247182100370520050591091822763712463273430149262739280891880522753123",
		"2",
	}

	gw := NewClient()

	for _, testAccount := range testnetAccounts {
		salt := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(SEED)
		calldata := append(accountCalldata, caigo.HexToBN(testAccount.PublicKey).String(), "0")

		_, err := gw.Deploy(context.Background(), proxyCompiled, types.DeployRequest{
			ContractAddressSalt: fmt.Sprintf("0x%x", salt),
			ConstructorCalldata: calldata,
		})
		if err != nil {
			t.Errorf("testnet: could not deploy contract: %v\n", err)
		}
	}
}

func TestCallGoerli(t *testing.T) {
	gw := NewClient()
	for _, testAccount := range testnetAccounts {
		call := types.FunctionCall{
			ContractAddress:    testAccount.Address,
			EntryPointSelector: "get_signer",
		}

		resp, err := gw.Call(context.Background(), call, "")
		if err != nil {
			t.Errorf("testnet: could 'Call' deployed contract: %v\n", err)
		}
		if len(resp) == 0 {
			t.Errorf("testnet: could get signing key for account: %v\n", err)
		}

		if resp[0] != testAccount.PublicKey {
			t.Errorf("testnet: signing key is incorrect: \n%s %v\n", resp[0], testAccount.PublicKey)
		}
	}
}

func TestE2EDevnet(t *testing.T) {
	gw := NewClient(WithChain("devnet"))

	deployTx, err := gw.Deploy(context.Background(), "../rpc/tests/counter.json", types.DeployRequest{})
	if err != nil {
		t.Errorf("could not deploy devnet counter: %v\n", err)
	}

	_, _, err = gw.PollTx(context.Background(), deployTx.TransactionHash, types.ACCEPTED_ON_L2, 1, 10)
	if err != nil {
		t.Errorf("could not deploy devnet counter: %v\n", err)
	}

	txDetails, err := gw.Transaction(context.Background(), TransactionOptions{TransactionHash: deployTx.TransactionHash})
	if err != nil {
		t.Errorf("fetching transaction: %v", err)
	}

	for i := 0; i < 3; i++ {
		rand := fmt.Sprintf("0x%x", rand.New(rand.NewSource(time.Now().UnixNano())).Intn(SEED))

		tx := []types.Transaction{
			{
				ContractAddress:    txDetails.Transaction.ContractAddress,
				EntryPointSelector: "set_rand",
				Calldata:           []string{rand},
			},
		}

		account, err := caigo.NewAccount(devnetAccounts[i].PrivateKey, devnetAccounts[i].Address, gw)
		if err != nil {
			t.Errorf("testnet: could not create account: %v\n", err)
		}

		feeEstimate, err := account.EstimateFee(context.Background(), tx)
		if err != nil {
			t.Errorf("testnet: could not estimate fee for transaction: %v\n", err)
		}
		fee := new(types.Felt)
		fee.Int = new(big.Int).SetUint64(feeEstimate.Amount * FEE_MARGIN / 100)

		execResp, err := account.Execute(context.Background(), fee, tx)
		if err != nil {
			t.Errorf("Could not execute test transaction: %v\n", err)
		}

		_, _, err = gw.PollTx(context.Background(), execResp.TransactionHash, types.ACCEPTED_ON_L2, 1, 10)
		if err != nil {
			t.Errorf("could not deploy devnet counter: %v\n", err)
		}

		call := types.FunctionCall{
			ContractAddress:    txDetails.Transaction.ContractAddress,
			EntryPointSelector: "get_rand",
		}
		callResp, err := gw.Call(context.Background(), call, "")
		if err != nil {
			t.Errorf("could not call counter contract: %v\n", err)
		}

		if rand != callResp[0] {
			t.Errorf("could not set value on counter contract: %v\n", err)
		}
	}
}

func DevnetAccounts() ([]TestAccountType, error) {
	req, err := http.NewRequest("GET", "http://localhost:5050/predeployed_accounts", nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var accounts []TestAccountType
	err = json.NewDecoder(resp.Body).Decode(&accounts)
	return accounts, err
}
