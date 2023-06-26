package gateway_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/smartcontractkit/caigo/gateway"
	"github.com/smartcontractkit/caigo/types"
)

func TestTransaction(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		TransactionIndex int
		BlockNumber      int
		Transaction      gateway.Transaction
		BlockHash        string
		Status           string
		opts             gateway.TransactionOptions
	}

	testSet := map[string][]testSetType{
		"devnet": {},
		"testnet": {{TransactionIndex: 3,
			Status:      "ACCEPTED_ON_L1",
			BlockNumber: 0,
			BlockHash:   "0x7d328a71faf48c5c3857e99f20a77b18522480956d1cd5bff1ff2df3c8b427b",
			opts:        gateway.TransactionOptions{TransactionHash: "0x1822471b7751cbaf98a5cce0003181af95d588e38c958739213af59f389fdc5"}}},
	}[testEnv]

	for _, test := range testSet {
		tx, err := testConfig.client.Transaction(context.Background(), test.opts)

		if err != nil {
			t.Fatal(err)
		}
		if tx.TransactionIndex != test.TransactionIndex {
			t.Fatalf("expecting %d, instead: %d", test.TransactionIndex, tx.TransactionIndex)
		}
	}
}

func TestTransactionStatus(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		TxStatus        string
		BlockHash       string
		TxFailureReason struct {
			ErrorMessage string
		}
		opts gateway.TransactionStatusOptions
	}
	testSet := map[string][]testSetType{
		"devnet": {},
		"testnet": {{
			TxStatus:  "ACCEPTED_ON_L1",
			BlockHash: "0x7d328a71faf48c5c3857e99f20a77b18522480956d1cd5bff1ff2df3c8b427b",
			opts:      gateway.TransactionStatusOptions{TransactionHash: "0x1822471b7751cbaf98a5cce0003181af95d588e38c958739213af59f389fdc5"}}},
	}[testEnv]

	for _, test := range testSet {
		tx, err := testConfig.client.TransactionStatus(context.Background(), test.opts)

		if err != nil {
			t.Fatal(err)
		}
		if tx.TxStatus != test.TxStatus {
			t.Fatalf("expecting %s, instead: %s", test.TxStatus, tx.TxStatus)
		}
	}
}

func TestTransactionID(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		TxHash string
		TxId   *big.Int
	}
	testSet := map[string][]testSetType{
		"devnet": {},
		"testnet": {{
			TxHash: "0x1822471b7751cbaf98a5cce0003181af95d588e38c958739213af59f389fdc5",
			TxId:   big.NewInt(3)}},
	}[testEnv]

	for _, test := range testSet {
		id, err := testConfig.client.TransactionID(context.Background(), test.TxHash)

		if err != nil {
			t.Fatal(err)
		}
		if id.String() != test.TxId.String() {
			t.Fatalf("expecting %d, instead: %d", id, test.TxId)
		}
	}
}

func TestTransactionHash(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		TxHash string
		TxId   *big.Int
	}
	testSet := map[string][]testSetType{
		"devnet": {},
		"testnet": {{
			TxHash: "0x1822471b7751cbaf98a5cce0003181af95d588e38c958739213af59f389fdc5",
			TxId:   big.NewInt(3)}},
	}[testEnv]

	for _, test := range testSet {
		hash, err := testConfig.client.TransactionHash(context.Background(), test.TxId)

		if err != nil {
			t.Fatal(err)
		}
		if hash != test.TxHash {
			t.Fatalf("expecting %s, instead: %s", hash, test.TxHash)
		}
	}
}

func TestTransactionReceipt(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		TransactionHash       string
		Status                types.TransactionState
		BlockHash             string
		BlockNumber           int
		TransactionIndex      int
		L1ToL2ConsumedMessage struct {
			FromAddress string
			ToAddress   string
			Selector    string
			Payload     []string
		}
		L2ToL1Messages []interface{}
		Events         []interface{}
		//ExecutionResources types.ExecutionResources
	}
	testSet := map[string][]testSetType{
		"devnet": {},
		"testnet": {{
			TransactionHash: "0x1822471b7751cbaf98a5cce0003181af95d588e38c958739213af59f389fdc5",
			Status:          "ACCEPTED_ON_L1"}},
	}[testEnv]

	for _, test := range testSet {
		receipt, err := testConfig.client.TransactionReceipt(context.Background(), test.TransactionHash)

		if err != nil {
			t.Fatal(err)
		}
		if receipt.Status != test.Status {
			t.Fatalf("expecting %s, instead: %s", receipt.Status, test.Status)
		}
	}
}

func TestTransactionTrace(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		TransactionHash string
		ContractAddress string
		//FunctionInvocation FunctionInvocation
		//Signature          []*Felt
	}
	testSet := map[string][]testSetType{
		"devnet": {},
		"testnet": {{
			TransactionHash: "0x1822471b7751cbaf98a5cce0003181af95d588e38c958739213af59f389fdc5",
			ContractAddress: "0x2adb4393384c09f049c06bc0070b7a2f72c9cbdcbe841fa7e109a520466cd66"}},
	}[testEnv]

	for _, test := range testSet {
		trace, err := testConfig.client.TransactionTrace(context.Background(), test.TransactionHash)

		if err != nil {
			t.Fatal(err)
		}
		if trace.FunctionInvocation.ContractAddress != test.ContractAddress {
			t.Fatalf("expecting %s, instead: %s", trace.FunctionInvocation.ContractAddress, test.ContractAddress)
		}
	}
}
