package rpc

import (
	"context"
	"testing"
)

// TestTransactionByHash tests transaction by hash
func TestTransactionByHash(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		TxHash                     TxnHash
		ExpectedContractAddress    string
		ExpectedEntrypointSelector string
	}
	testSet := map[string][]testSetType{
		"mock": {
			{
				TxHash:                     TxnHash("0xdeadbeef"),
				ExpectedContractAddress:    "0xdeadbeef",
				ExpectedEntrypointSelector: "0xdeadbeef",
			},
		},
		"testnet": {
			{
				TxHash:                     TxnHash("0x705547f8f2f8fdfb10ed533d909f76482bb293c5a32648d476774516a0bebd0"),
				ExpectedContractAddress:    "0x315e364b162653e5c7b23efd34f8da27ba9c069b68e3042b7d76ce1df890313",
				ExpectedEntrypointSelector: "0x15d40a3d6ca2ac30f4031e42be28da9b056fef9bb7357ac5e85627ee876e5ad",
			},
		},
		"mainnet": {
			{
				TxHash:                     TxnHash("0x5f904b9185d4ed442846ac7e26bc4c60249a2a7f0bb85376c0bc7459665bae6"),
				ExpectedContractAddress:    "0x3b4be7def2fc08589348966255e101824928659ebb724855223ff3a8c831efa",
				ExpectedEntrypointSelector: "0x2913ee03e5e3308c41e308bd391ea4faac9b9cb5062c76a6b3ab4f65397e106",
			},
		},
	}[testEnv]

	for _, test := range testSet {
		tx, err := testConfig.client.TransactionByHash(context.Background(), test.TxHash)
		if err != nil {
			t.Fatal(err)
		}
		if tx == nil {
			t.Fatal("transaction should exist")
		}
		txTyped, ok := (*tx).(InvokeTxnV0)
		if !ok {
			t.Fatalf("transaction should be InvokeTxnV0, instead %T", tx)

		}
		if txTyped.FunctionCall.ContractAddress != Address(test.ExpectedContractAddress) {
			t.Fatalf("expecting contract %s, got %s", test.ExpectedContractAddress, txTyped.FunctionCall.ContractAddress)
		}
		if txTyped.FunctionCall.EntryPointSelector != test.ExpectedEntrypointSelector {
			t.Fatalf("expecting entrypoint %s, got %s", test.ExpectedEntrypointSelector, txTyped.FunctionCall.EntryPointSelector)
		}
	}
}

// TestTransactionReceipt tests transaction receipt
func TestTransactionReceipt(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		TxnHash        TxnHash
		ExpectedStatus TxnStatus
	}
	testSet := map[string][]testSetType{
		"mock": {
			{
				TxnHash:        TxnHash("0xdeadbeef"),
				ExpectedStatus: TxnStatus("ACCEPTED_ON_L1"),
			},
		},
		"testnet": {
			{
				TxnHash:        TxnHash("0x705547f8f2f8fdfb10ed533d909f76482bb293c5a32648d476774516a0bebd0"),
				ExpectedStatus: TxnStatus("ACCEPTED_ON_L1"),
			},
		},
		"mainnet": {
			{
				TxnHash:        TxnHash("0x5f904b9185d4ed442846ac7e26bc4c60249a2a7f0bb85376c0bc7459665bae6"),
				ExpectedStatus: TxnStatus("ACCEPTED_ON_L1"),
			},
		},
	}[testEnv]

	for _, test := range testSet {
		txReceiptInterface, err := testConfig.client.TransactionReceipt(context.Background(), test.TxnHash)
		if err != nil {
			t.Fatal(err)
		}
		if txReceiptInterface == nil {
			t.Fatal("transaction receipt should exist")
		}
		txnReceipt, ok := txReceiptInterface.(InvokeTxnReceipt)
		if !ok {
			t.Fatalf("transaction receipt should be InvokeTxnReceipt, instead %T", txReceiptInterface)

		}
		if txnReceipt.Status != test.ExpectedStatus {
			t.Fatalf("expecting status %s, got %s", test.ExpectedStatus, txnReceipt.Status)
		}
	}
}
