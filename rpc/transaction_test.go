package rpc

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	ctypes "github.com/NethermindEth/starknet.go/types"
	"github.com/google/go-cmp/cmp"
)

// TestTransactionByHash tests transaction by hash
func TestTransactionByHash(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		TxHash      ctypes.Felt
		ExpectedTxn Transaction
	}
	testSet := map[string][]testSetType{
		"mock": {
			{
				TxHash:      ctypes.StrToFelt("0x705547f8f2f8fdfb10ed533d909f76482bb293c5a32648d476774516a0bebd0"),
				ExpectedTxn: InvokeTxnV00x705547f8f2f8f,
			},
		},
		"testnet": {
			{
				TxHash:      ctypes.StrToFelt("0x705547f8f2f8fdfb10ed533d909f76482bb293c5a32648d476774516a0bebd0"),
				ExpectedTxn: InvokeTxnV00x705547f8f2f8f,
			},
		},
		"mainnet": {},
	}[testEnv]
	for _, test := range testSet {
		spy := NewSpy(testConfig.provider.c)
		testConfig.provider.c = spy
		tx, err := testConfig.provider.TransactionByHash(context.Background(), test.TxHash)
		if err != nil {
			t.Fatal(err)
		}
		if tx == nil {
			t.Fatal("transaction should exist")
		}

		txTyped, ok := (tx).(InvokeTxnV0)
		if !ok {
			t.Fatalf("transaction should be InvokeTxnV0, instead %T", tx)
		}
		if !cmp.Equal(test.ExpectedTxn, txTyped) {
			t.Fatalf("the expected transaction blocks to match, instead: %s", cmp.Diff(test.ExpectedTxn, txTyped))
		}
	}
}

func TestTransactionByBlockIdAndIndex(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		BlockID     BlockID
		Index       uint64
		ExpectedTxn Transaction
	}
	testSet := map[string][]testSetType{
		"mock": {
			{
				BlockID:     WithBlockNumber(300000),
				Index:       0,
				ExpectedTxn: InvokeTxnV0_300000_0,
			},
		},
		// "testnet": {
		// 	{
		// 		BlockID:     WithBlockNumber(300000),
		// 		Index:       0,
		// 		ExpectedTxn: InvokeTxnV0_300000_0,
		// 	},
		// },
		"mainnet": {},
	}[testEnv]
	for _, test := range testSet {
		spy := NewSpy(testConfig.provider.c)
		testConfig.provider.c = spy
		tx, err := testConfig.provider.TransactionByBlockIdAndIndex(context.Background(), test.BlockID, test.Index)
		if err != nil {
			t.Fatal(err)
		}
		if tx == nil {
			t.Fatal("transaction should exist")
		}
		txTyped, ok := (tx).(InvokeTxnV0)
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
func TestTransactionReceipt_MatchesCapturedTransaction(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		TxnHash            ctypes.Felt
		ExpectedTxnReceipt TransactionReceipt
	}
	testSet := map[string][]testSetType{
		"mock": {},
		"testnet": {
			{
				TxnHash:            ctypes.StrToFelt("0x40c82f79dd2bc1953fc9b347a3e7ab40fe218ed5740bf4e120f74e8a3c9ac99"),
				ExpectedTxnReceipt: receiptTxn310370_0,
			},
		},
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		spy := NewSpy(testConfig.provider.c)
		testConfig.provider.c = spy
		txReceiptInterface, err := testConfig.provider.TransactionReceipt(context.Background(), test.TxnHash)
		if err != nil {
			t.Fatal(err)
		}
		if txReceiptInterface == nil {
			t.Fatal("transaction receipt should exist")
		}
		txnReceipt, ok := txReceiptInterface.(InvokeTransactionReceipt)
		if !ok {
			t.Fatalf("transaction receipt should be InvokeTransactionReceipt, instead %T", txReceiptInterface)
		}
		if !cmp.Equal(test.ExpectedTxnReceipt, txnReceipt) {
			t.Fatalf("the expected transaction blocks to match, instead: %s", cmp.Diff(test.ExpectedTxnReceipt, txnReceipt))
		}
	}
}

// TestTransactionReceipt tests transaction receipt
func TestTransactionReceipt_MatchesStatus(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		TxnHash     ctypes.Felt
		StatusMatch string
	}
	testSet := map[string][]testSetType{
		"mock": {},
		"testnet": {
			{
				TxnHash:     ctypes.StrToFelt("0x650667fb0f17e63e1c9d1040e750d160f3dbfebcab990e7d4382f33468b1b59"),
				StatusMatch: "(ACCEPTED_ON_L1|ACCEPTED_ON_L2|PENDING)",
			},
		},
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		spy := NewSpy(testConfig.provider.c, false)
		testConfig.provider.c = spy
		txReceiptInterface, err := testConfig.provider.TransactionReceipt(context.Background(), test.TxnHash)
		if err != nil {
			t.Fatal(err)
		}
		if txReceiptInterface == nil {
			t.Fatal("transaction receipt should exist")
		}
		txnReceipt, ok := txReceiptInterface.(InvokeTransactionReceipt)
		if !ok {
			t.Fatalf("transaction receipt should be InvokeTransactionReceipt, instead %T", txReceiptInterface)
		}
		if ok, err := regexp.MatchString(test.StatusMatch, string(txnReceipt.Status)); err != nil || !ok {
			t.Fatal("error checking transaction status", ok, err, txnReceipt.Status)
		}
		fmt.Println("transaction status", txnReceipt.Status)
	}
}

// TestDeployOrDeclareReceipt tests deploy or declare receipt
func TestDeployOrDeclareReceipt(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		TxnHash            ctypes.Felt
		ExpectedTxnReceipt TransactionReceipt
	}
	testSet := map[string][]testSetType{
		"mock": {},
		"testnet": {
			{
				TxnHash:            ctypes.StrToFelt("0x35bd2978d2061b3463498f83c09322ed6a82e4b2a188506525e272a7adcdf6a"),
				ExpectedTxnReceipt: receiptTxn310843_14,
			},
			{
				TxnHash:            ctypes.StrToFelt("0x46a9f52a96b2d226407929e04cb02507e531f7c78b9196fc8c910351d8c33f3"),
				ExpectedTxnReceipt: receiptTxn300114_3,
			},
		},
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		spy := NewSpy(testConfig.provider.c)
		testConfig.provider.c = spy
		txReceiptInterface, err := testConfig.provider.TransactionReceipt(context.Background(), test.TxnHash)
		if err != nil {
			t.Fatal(err)
		}
		if txReceiptInterface == nil {
			t.Fatal("transaction receipt should exist")
		}
		txnDeployReceipt, ok1 := txReceiptInterface.(DeployTransactionReceipt)
		txnDeclareReceipt, ok2 := txReceiptInterface.(DeclareTransactionReceipt)
		if !ok1 && !ok2 {
			t.Fatalf("transaction receipt should be Deploy or Declare, instead %T", txReceiptInterface)
		}
		switch {
		case ok1:
			if !cmp.Equal(test.ExpectedTxnReceipt, txnDeployReceipt) {
				t.Fatalf("the expected transaction blocks to match, instead: %s", cmp.Diff(test.ExpectedTxnReceipt, txnDeployReceipt))
			}
		case ok2:
			if !cmp.Equal(test.ExpectedTxnReceipt, txnDeclareReceipt) {
				t.Fatalf("the expected transaction blocks to match, instead: %s", cmp.Diff(test.ExpectedTxnReceipt, txnDeclareReceipt))
			}
		}
	}
}
