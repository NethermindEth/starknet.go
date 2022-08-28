package rpc

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// TestTransactionByHash tests transaction by hash
func TestTransactionByHash(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		TxHash      TxnHash
		ExpectedTxn Txn
	}
	testSet := map[string][]testSetType{
		"mock": {
			{
				TxHash:      TxnHash("0x705547f8f2f8fdfb10ed533d909f76482bb293c5a32648d476774516a0bebd0"),
				ExpectedTxn: InvokeTxnV00x705547f8f2f8f,
			},
		},
		"testnet": {
			{
				TxHash:      TxnHash("0x705547f8f2f8fdfb10ed533d909f76482bb293c5a32648d476774516a0bebd0"),
				ExpectedTxn: InvokeTxnV00x705547f8f2f8f,
			},
		},
		"mainnet": {},
	}[testEnv]
	for _, test := range testSet {
		spy := NewSpy(testConfig.client.c)
		testConfig.client.c = spy
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
		diff, err := spy.Compare(txTyped, false)
		if err != nil {
			t.Fatal("expecting to match", err)
		}
		if diff != "FullMatch" {
			spy.Compare(txTyped, true)
			t.Fatal("structure expecting to be FullMatch, instead", diff)
		}
		if !cmp.Equal(test.ExpectedTxn, txTyped) {
			t.Fatalf("the expected transaction blocks to match, instead: %s", cmp.Diff(test.ExpectedTxn, txTyped))
		}
	}
}

// TestTransactionReceipt tests transaction receipt
func TestTransactionReceipt(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		TxnHash            TxnHash
		ExpectedTxnReceipt TxnReceipt
	}
	testSet := map[string][]testSetType{
		"mock": {},
		"testnet": {
			{
				TxnHash:            TxnHash("0x40c82f79dd2bc1953fc9b347a3e7ab40fe218ed5740bf4e120f74e8a3c9ac99"),
				ExpectedTxnReceipt: receiptTxn310370_0,
			},
		},
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		spy := NewSpy(testConfig.client.c)
		testConfig.client.c = spy
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
		diff, err := spy.Compare(txnReceipt, false)
		if err != nil {
			t.Fatal("expecting to match", err)
		}
		if diff != "FullMatch" {
			spy.Compare(txnReceipt, true)
			t.Fatal("structure expecting to be FullMatch, instead", diff)
		}
		if !cmp.Equal(test.ExpectedTxnReceipt, txnReceipt) {
			t.Fatalf("the expected transaction blocks to match, instead: %s", cmp.Diff(test.ExpectedTxnReceipt, txnReceipt))
		}
	}
}

// TestDeployOrDeclareReceipt tests deploy or declare receipt
func TestDeployOrDeclareReceipt(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		TxnHash            TxnHash
		ExpectedTxnReceipt TxnReceipt
	}
	testSet := map[string][]testSetType{
		"mock": {},
		"testnet": {
			{
				TxnHash:            TxnHash("0x35bd2978d2061b3463498f83c09322ed6a82e4b2a188506525e272a7adcdf6a"),
				ExpectedTxnReceipt: receiptTxn310843_14,
			},
			{
				TxnHash:            TxnHash("0x46a9f52a96b2d226407929e04cb02507e531f7c78b9196fc8c910351d8c33f3"),
				ExpectedTxnReceipt: receiptTxn300114_3,
			},
		},
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		spy := NewSpy(testConfig.client.c)
		testConfig.client.c = spy
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
		diff, err := spy.Compare(txnReceipt, false)
		if err != nil {
			t.Fatal("expecting to match", err)
		}
		if diff != "FullMatch" {
			spy.Compare(txnReceipt, true)
			t.Fatal("structure expecting to be FullMatch, instead", diff)
		}
		if !cmp.Equal(test.ExpectedTxnReceipt, txnReceipt) {
			t.Fatalf("the expected transaction blocks to match, instead: %s", cmp.Diff(test.ExpectedTxnReceipt, txnReceipt))
		}
	}
}
