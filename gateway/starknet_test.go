package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"math/rand"
	"net/http"
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

func TestDeclare(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct{}
	testSet := map[string][]testSetType{
		"devnet":  {{}},
		"mainnet": {},
		"mock":    {},
		"testnet": {{}},
	}[testEnv]

	for _, env := range testSet {
		gw := testConfig.client
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
	gw := NewClient()
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
		fee.Int = new(big.Int).SetUint64(feeEstimate.OverallFee * FEE_MARGIN / 100)

		nonce, err := gw.AccountNonce(context.Background(), testAccount.Address)
		if err != nil {
			t.Errorf("testnet: could not get account nonce: %v", err)
		}

		_, err = account.Execute(context.Background(), 
			&caigo.ExecuteDetails{
				Calls: &testAccount.Transactions,
				MaxFee: fee,
				Nonce: nonce,
			})
		if err != nil {
			t.Errorf("Could not execute test transaction: %v\n", err)
		}
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
	testConfig := beforeEach(t)

	type testSetType struct{}
	testSet := map[string][]testSetType{
		"devnet":  {{}},
		"mainnet": {},
		"mock":    {},
		"testnet": {},
	}[testEnv]

	for _, env := range testSet {
		gw := testConfig.client

		deployTx, err := gw.Deploy(context.Background(), "../rpc/tests/counter.json", types.DeployRequest{})
		if err != nil {
			t.Errorf("%s: could not deploy devnet counter: %v", env, err)
		}

		_, _, err = gw.PollTx(context.Background(), deployTx.TransactionHash, types.ACCEPTED_ON_L2, 1, 10)
		if err != nil {
			t.Errorf("%s: could not deploy devnet counter: %v", env, err)
		}

		txDetails, err := gw.Transaction(context.Background(), TransactionOptions{TransactionHash: deployTx.TransactionHash})
		if err != nil {
			t.Errorf("%s: fetching transaction: %v", env, err)
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
			fee.Int = new(big.Int).SetUint64(feeEstimate.OverallFee * FEE_MARGIN / 100)

			nonce, err := gw.AccountNonce(context.Background(), account.Address)
			if err != nil {
				t.Errorf("testnet: could not get account nonce: %v", err)
			}

			execResp, err := account.Execute(context.Background(), 
				&caigo.ExecuteDetails{
					Calls: &tx,
					MaxFee: fee,
					Nonce: nonce,
				})
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
