package rpc

import (
	"context"
	"regexp"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/utils"
	"github.com/google/go-cmp/cmp"
	"github.com/test-go/testify/require"
)

// TestTransactionByHash tests transaction by hash
func TestTransactionByHash(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		TxHash      *felt.Felt
		ExpectedTxn Transaction
	}

	var InvokeTxnV1example = InvokeTxnV1{
		CommonTransaction: CommonTransaction{
			TransactionHash: utils.TestHexToFelt(t, "0x1779df1c6de5136ad2620f704b645e9cbd554b57d37f08a06ea60142269c5a5"),

			BroadcastedTxnCommonProperties: BroadcastedTxnCommonProperties{

				Type:    TransactionType_Invoke,
				MaxFee:  utils.TestHexToFelt(t, "0x17970b794f000"),
				Version: TransactionV1,
				Nonce:   utils.TestHexToFelt(t, "0x2d"),
				Signature: []*felt.Felt{
					utils.TestHexToFelt(t, "0xe500c4014c055c3304d8a125cfef44638ffa5b0f6840916049667a4c38aa1c"),
					utils.TestHexToFelt(t, "0x45ac538bfce5d8c5741b4421bbdc99f5849451acae75d2048d7cc4bb029ca77"),
				},
			},
		},
		SenderAddress: utils.TestHexToFelt(t, "0x66dd340c03b6b7866fa7bb4bb91cc9e9c2a8eedc321985f334fd55de5e4e071"),
		Calldata: []*felt.Felt{
			utils.TestHexToFelt(t, "0x3"),
			utils.TestHexToFelt(t, "0x39a04b968d794fb076b0fbb146c12b48a23aa785e3d2e5be1982161f7536218"),
			utils.TestHexToFelt(t, "0x2f0b3c5710379609eb5495f1ecd348cb28167711b73609fe565a72734550354"),
			utils.TestHexToFelt(t, "0x0"),
			utils.TestHexToFelt(t, "0x3"),
			utils.TestHexToFelt(t, "0x3207980cd08767c9310d197c38b1a58b2a9bceb61dd9a99f51b407798702991"),
			utils.TestHexToFelt(t, "0x2f0b3c5710379609eb5495f1ecd348cb28167711b73609fe565a72734550354"),
			utils.TestHexToFelt(t, "0x3"),
			utils.TestHexToFelt(t, "0x3"),
			utils.TestHexToFelt(t, "0x42969068f9e84e9bf1c7bb6eb627455287e58f866ba39e45b123f9656aed5e9"),
			utils.TestHexToFelt(t, "0x2f0b3c5710379609eb5495f1ecd348cb28167711b73609fe565a72734550354"),
			utils.TestHexToFelt(t, "0x6"),
			utils.TestHexToFelt(t, "0x3"),
			utils.TestHexToFelt(t, "0x9"),
			utils.TestHexToFelt(t, "0x47487560da4c5c5755897e527a5fda37422b5ba02a2aba1ca3ce2b24dfd142e"),
			utils.TestHexToFelt(t, "0xde0b6b3a7640000"),
			utils.TestHexToFelt(t, "0x0"),
			utils.TestHexToFelt(t, "0x47487560da4c5c5755897e527a5fda37422b5ba02a2aba1ca3ce2b24dfd142e"),
			utils.TestHexToFelt(t, "0x10f0cf064dd59200000"),
			utils.TestHexToFelt(t, "0x0"),
			utils.TestHexToFelt(t, "0x47487560da4c5c5755897e527a5fda37422b5ba02a2aba1ca3ce2b24dfd142e"),
			utils.TestHexToFelt(t, "0x21e19e0c9bab2400000"),
			utils.TestHexToFelt(t, "0x0"),
		},
	}

	testSet := map[string][]testSetType{
		"mock": {
			{
				TxHash:      utils.TestHexToFelt(t, "0x1779df1c6de5136ad2620f704b645e9cbd554b57d37f08a06ea60142269c5a5"),
				ExpectedTxn: InvokeTxnV1example,
			},
		},
		"testnet": {
			{
				TxHash:      utils.TestHexToFelt(t, "0x1779df1c6de5136ad2620f704b645e9cbd554b57d37f08a06ea60142269c5a5"),
				ExpectedTxn: InvokeTxnV1example,
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

		txCasted, ok := (tx).(InvokeTxnV1)
		if !ok {
			t.Fatalf("transaction should be InvokeTxnV1, instead %T", tx)
		}
		require.Equal(t, txCasted.Type, TransactionType_Invoke)
		require.Equal(t, txCasted.TransactionHash, test.TxHash)
		require.Equal(t, txCasted, test.ExpectedTxn)
	}
}

func TestTransactionByBlockIdAndIndex(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		BlockID     BlockID
		Index       uint64
		ExpectedTxn Transaction
	}

	var InvokeTxnV1example = InvokeTxnV1{
		CommonTransaction: CommonTransaction{
			TransactionHash: utils.TestHexToFelt(t, "0x705547f8f2f8fdfb10ed533d909f76482bb293c5a32648d476774516a0bebd0"),
			BroadcastedTxnCommonProperties: BroadcastedTxnCommonProperties{
				Type:    TransactionType_Invoke,
				MaxFee:  utils.TestHexToFelt(t, "0x53685de02fa5"),
				Version: TransactionV1,
				Nonce:   &felt.Zero,
				Signature: []*felt.Felt{
					utils.TestHexToFelt(t, "0x4a7849de7b91e52cd0cdaf4f40aa67f54a58e25a15c60e034d2be819c1ecda4"),
					utils.TestHexToFelt(t, "0x227fcad2a0007348e64384649365e06d41287b1887999b406389ee73c1d8c4c"),
				},
			},
		},
		SenderAddress: utils.TestHexToFelt(t, "0x315e364b162653e5c7b23efd34f8da27ba9c069b68e3042b7d76ce1df890313"),
		Calldata: []*felt.Felt{
			utils.TestHexToFelt(t, "0x1"),
			utils.TestHexToFelt(t, "0x13befe6eda920ce4af05a50a67bd808d67eee6ba47bb0892bef2d630eaf1bba"),
		},
	}

	testSet := map[string][]testSetType{
		"mock": {
			{
				BlockID:     WithBlockNumber(300000),
				Index:       0,
				ExpectedTxn: InvokeTxnV1example,
			},
		},
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
		txCasted, ok := (tx).(InvokeTxnV1)
		if !ok {
			t.Fatalf("transaction should be InvokeTxnV0, instead %T", tx)
		}
		require.Equal(t, txCasted.Type, TransactionType_Invoke)
		require.Equal(t, txCasted, test.ExpectedTxn)
	}
}

// TestTransactionReceipt tests transaction receipt
func TestTransactionReceipt_MatchesCapturedTransaction(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		TxnHash            *felt.Felt
		ExpectedTxnReceipt TransactionReceipt
	}
	var receiptTxn310370_0 = InvokeTransactionReceipt(CommonTransactionReceipt{
		TransactionHash: utils.TestHexToFelt(t, "0x40c82f79dd2bc1953fc9b347a3e7ab40fe218ed5740bf4e120f74e8a3c9ac99"),
		ActualFee:       utils.TestHexToFelt(t, "0x1709a2f3a2"),
		Type:            "INVOKE",
		ExecutionStatus: TxnExecutionStatusSUCCEEDED,
		FinalityStatus:  TxnFinalityStatusAcceptedOnL1,
		BlockHash:       utils.TestHexToFelt(t, "0x6c2fe3db009a2e008c2d65fca14204f3405cb74742fcf685f02473acaf70c72"),
		BlockNumber:     310370,
		MessagesSent:    []MsgToL1{},
		Events: []Event{
			{
				FromAddress: utils.TestHexToFelt(t, "0x37de00fb1416936b3074fc78bcc811d83046009b162c4a822ce84dabedd0ea9"),
				Data: []*felt.Felt{
					utils.TestHexToFelt(t, "0x0"),
					utils.TestHexToFelt(t, "0x35b32bb4a1969175fb14b6c09838d1b3200724cc4d2b0891be319764021f5ac"),
					utils.TestHexToFelt(t, "0xe9"),
					utils.TestHexToFelt(t, "0x0"),
				},
				Keys: []*felt.Felt{utils.TestHexToFelt(t, "0x99cd8bde557814842a3121e8ddfd433a539b8c9f14bf31ebf108d12e6196e9")},
			},
			{
				FromAddress: utils.TestHexToFelt(t, "0x33830ce413e4c096eef81b5e6ffa9b9f5d963f57b8cd63c9ae4c839c383c1a6"),
				Data: []*felt.Felt{
					utils.TestHexToFelt(t, "0x61c6e7484657e5dc8b21677ffa33e4406c0600bba06d12cf1048fdaa55bdbc3"),
					utils.TestHexToFelt(t, "0x2e28403d7ee5e337b7d456327433f003aa875c29631906908900058c83d8cb6"),
				},
				Keys: []*felt.Felt{utils.TestHexToFelt(t, "0xf806f71b19e4744968b37e3fb288e61309ab33a782ea9d11e18f67a1fbb110")},
			},
		},
	})

	testSet := map[string][]testSetType{
		"mock": {},
		"testnet": {
			{
				TxnHash:            utils.TestHexToFelt(t, "0x40c82f79dd2bc1953fc9b347a3e7ab40fe218ed5740bf4e120f74e8a3c9ac99"),
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
		TxnHash         *felt.Felt
		ExecutionStatus string
	}
	testSet := map[string][]testSetType{
		"mock": {},
		"testnet": {
			{
				TxnHash:         utils.TestHexToFelt(t, "0x650667fb0f17e63e1c9d1040e750d160f3dbfebcab990e7d4382f33468b1b59"),
				ExecutionStatus: "(SUCCEEDED|REVERTED)",
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
		if ok, err := regexp.MatchString(test.ExecutionStatus, string(txnReceipt.ExecutionStatus)); err != nil || !ok {
			t.Fatal("error checking transaction status", ok, err, txnReceipt.ExecutionStatus)
		}
	}
}

// TestDeployOrDeclareReceipt tests deploy or declare receipt
func TestDeployOrDeclareReceipt(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		TxnHash            *felt.Felt
		ExpectedTxnReceipt TransactionReceipt
	}

	var receiptTxn310843_14 = DeployTransactionReceipt{
		CommonTransactionReceipt: CommonTransactionReceipt{
			TransactionHash: utils.TestHexToFelt(t, "0x035bd2978d2061b3463498f83c09322ed6a82e4b2a188506525e272a7adcdf6a"),
			ActualFee:       &felt.Zero,
			FinalityStatus:  TxnFinalityStatusAcceptedOnL1,
			ExecutionStatus: TxnExecutionStatusSUCCEEDED,
			BlockHash:       utils.TestHexToFelt(t, "0x0424fba26a7760b63895abe0c366c2d254cb47090c6f9e91ba2b3fa0824d4fc9"),
			BlockNumber:     310843,
			Type:            "DEPLOY",
			MessagesSent:    []MsgToL1{},
			Events:          []Event{},
		},
		ContractAddress: "0x21c40b1377353924e185c9536469787dbe0cdb77b6877fa3a9946b795c71ec7",
	}

	var receiptTxn300114_3 = DeclareTransactionReceipt(
		CommonTransactionReceipt{
			TransactionHash: utils.TestHexToFelt(t, "0x46a9f52a96b2d226407929e04cb02507e531f7c78b9196fc8c910351d8c33f3"),
			ActualFee:       utils.TestHexToFelt(t, "0x0"),
			FinalityStatus:  TxnFinalityStatusAcceptedOnL1,
			ExecutionStatus: TxnExecutionStatusSUCCEEDED,
			BlockHash:       utils.TestHexToFelt(t, "0x184268bfbce24766fa53b65c9c8b30b295e145e8281d543a015b46308e27fdf"),
			BlockNumber:     300114,
			Type:            "DECLARE",
			MessagesSent:    []MsgToL1{},
			Events:          []Event{},
		})

	testSet := map[string][]testSetType{
		"mock": {},
		"testnet": {
			{
				TxnHash:            utils.TestHexToFelt(t, "0x35bd2978d2061b3463498f83c09322ed6a82e4b2a188506525e272a7adcdf6a"),
				ExpectedTxnReceipt: receiptTxn310843_14,
			},
			{
				TxnHash:            utils.TestHexToFelt(t, "0x46a9f52a96b2d226407929e04cb02507e531f7c78b9196fc8c910351d8c33f3"),
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
