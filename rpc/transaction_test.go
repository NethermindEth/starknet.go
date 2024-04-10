package rpc

import (
	"context"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/utils"
	"github.com/test-go/testify/require"
)

// TestTransactionByHash tests transaction by hash
//
// Parameters:
// - t: the testing object for running the test cases
// Returns:
// none
func TestTransactionByHash(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		TxHash      *felt.Felt
		ExpectedTxn Transaction
	}

	var DeclareTxnV2Example = DeclareTxnV2{
		Type:              TransactionType_Declare,
		Version:           TransactionV2,
		MaxFee:            utils.TestHexToFelt(t, "0x4a0fbb2d7a43"),
		ClassHash:         utils.TestHexToFelt(t, "0x79b7ec8fdf40a4ff6ed47123049dfe36b5c02db93aa77832682344775ef70c6"),
		CompiledClassHash: utils.TestHexToFelt(t, "0x7130f75fc2f1400813d1e96ea7ebee334b568a87b645a62aade0eb2fa2cf252"),
		Nonce:             utils.TestHexToFelt(t, "0x16e"),
		Signature: []*felt.Felt{
			utils.TestHexToFelt(t, "0x5569787df42fece1184537b0d480900a403386355b9d6a59e7c7a7e758287f0"),
			utils.TestHexToFelt(t, "0x2acaeea2e0817da33ed5dbeec295b0177819b5a5a50b0a669e6eecd88e42e92"),
		},
		SenderAddress: utils.TestHexToFelt(t, "0x5fd4befee268bf6880f955875cbed3ade8346b1f1e149cc87b317e62b6db569"),
	}

	testSet := map[string][]testSetType{
		"mock": {
			{
				TxHash:      utils.TestHexToFelt(t, "0xd109474cd037bad60a87ba0ccf3023d5f2d1cd45220c62091d41a614d38eda"),
				ExpectedTxn: DeclareTxnV2Example,
			},
		},
		"testnet": {
			{
				TxHash:      utils.TestHexToFelt(t, "0xd109474cd037bad60a87ba0ccf3023d5f2d1cd45220c62091d41a614d38eda"),
				ExpectedTxn: DeclareTxnV2Example,
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

		txCasted, ok := (tx).(DeclareTxnV2)
		if !ok {
			t.Fatalf("transaction should be DeclareTnxV2, instead %T", tx)
		}
		require.Equal(t, txCasted.Type, TransactionType_Declare)
		require.Equal(t, txCasted, test.ExpectedTxn)
	}
}

// TestTransactionByBlockIdAndIndex tests the TransactionByBlockIdAndIndex function.
//
// It sets up a test environment and defines a test set. For each test in the set,
// it creates a spy object and assigns it to the provider's c field. It then calls
// the TransactionByBlockIdAndIndex function with the specified block ID and index.
// If there is an error, it fails the test. If the transaction is nil, it fails the test.
// If the transaction is not of type InvokeTxnV1, it fails the test. Finally, it asserts
// that the transaction type is TransactionType_Invoke and that the transaction is equal to the expected transaction.
//
// Parameters:
// - t: the testing object for running the test cases
// Returns:
//
//	none
func TestTransactionByBlockIdAndIndex(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		BlockID     BlockID
		Index       uint64
		ExpectedTxn Transaction
	}

	var InvokeTxnV1example = InvokeTxnV1{
		Type:    TransactionType_Invoke,
		MaxFee:  utils.TestHexToFelt(t, "0x53685de02fa5"),
		Version: TransactionV1,
		Nonce:   &felt.Zero,
		Signature: []*felt.Felt{
			utils.TestHexToFelt(t, "0x4a7849de7b91e52cd0cdaf4f40aa67f54a58e25a15c60e034d2be819c1ecda4"),
			utils.TestHexToFelt(t, "0x227fcad2a0007348e64384649365e06d41287b1887999b406389ee73c1d8c4c"),
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
			t.Fatalf("transaction should be InvokeTxnV1, instead %T", tx)
		}
		require.Equal(t, txCasted.Type, TransactionType_Invoke)
		require.Equal(t, txCasted, test.ExpectedTxn)
	}
}

func TestTransactionReceipt(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		TxnHash      *felt.Felt
		ExpectedResp TransactionReceiptWithBlockInfo
	}
	var receiptTxn52767_16 = InvokeTransactionReceipt(CommonTransactionReceipt{
		TransactionHash: utils.TestHexToFelt(t, "0xf2f3d50192637e8d5e817363460c39d3a668fe12f117ecedb9749466d8352b"),
		BlockHash:       utils.TestHexToFelt(t, "0x4ae5d52c75e4dea5694f456069f830cfbc7bec70427eee170c3385f751b8564"),
		BlockNumber:     52767,
		ActualFee: FeePayment{
			Amount: utils.TestHexToFelt(t, "0x16409a78a10b00"),
			Unit:   UnitStrk,
		},
		Type:            "INVOKE",
		ExecutionStatus: TxnExecutionStatusSUCCEEDED,
		FinalityStatus:  TxnFinalityStatusAcceptedOnL1,
		MessagesSent:    []MsgToL1{},
		Events: []Event{
			{
				FromAddress: utils.TestHexToFelt(t, "0x243d436e1f7cea085aaa42834975488029b1ebf67cea1d2e86f7de58e7d34a3"),
				Data: []*felt.Felt{
					utils.TestHexToFelt(t, "0x3028044a4c4df95c0b0a907307c6feffa76b9c38e83088ade29b186a250eb13"),
					utils.TestHexToFelt(t, "0x3"),
					utils.TestHexToFelt(t, "0x17a393a5e943cec833c8a8f4cbbf7c58361fb2fdd9caa0c36d901eedec4938e"),
					utils.TestHexToFelt(t, "0x776731d30bd922ac0390edfc664ed31b232aa7c7ce389c333e34c6b32957532"),
					utils.TestHexToFelt(t, "0x40851db0ebaebb9f8a18eda25005c050793f2a69e9e7d1f44bc133752898918"),
				},
				Keys: []*felt.Felt{
					utils.TestHexToFelt(t, "0x243d436e1f7cea085aaa42834975488029b1ebf67cea1d2e86f7de58e7d34a3"),
				},
			},
			{
				FromAddress: utils.TestHexToFelt(t, "0x243d436e1f7cea085aaa42834975488029b1ebf67cea1d2e86f7de58e7d34a3"),
				Data: []*felt.Felt{
					utils.TestHexToFelt(t, "0x6016d919abf2ddefe03dacc2ff5c8f42eb80cf65add1e90dd73c5c5e06ef3e2"),
					utils.TestHexToFelt(t, "0x1176a1bd84444c89232ec27754698e5d2e7e1a7f1539f12027f28b23ec9f3d8"),
					utils.TestHexToFelt(t, "0x16409a78a10b00"),
					utils.TestHexToFelt(t, "0x0"),
				},
				Keys: []*felt.Felt{
					utils.TestHexToFelt(t, "0x15bd0500dc9d7e69ab9577f73a8d753e8761bed10f25ba0f124254dc4edb8b4"),
				},
			},
		},
		ExecutionResources: ExecutionResources{
			ComputationResources: ComputationResources{
				Steps:          5774,
				PedersenApps:   24,
				ECOPApps:       3,
				RangeCheckApps: 152,
			},
		},
	})

	var receiptTxnIntegration = InvokeTransactionReceipt(CommonTransactionReceipt{
		TransactionHash: utils.TestHexToFelt(t, "0x49728601e0bb2f48ce506b0cbd9c0e2a9e50d95858aa41463f46386dca489fd"),
		ActualFee:       FeePayment{Amount: utils.TestHexToFelt(t, "0x16d8b4ad4000"), Unit: UnitStrk},
		Type:            "INVOKE",
		ExecutionStatus: TxnExecutionStatusSUCCEEDED,
		FinalityStatus:  TxnFinalityStatusAcceptedOnL2,
		MessagesSent:    []MsgToL1{},
		Events: []Event{
			{
				FromAddress: utils.TestHexToFelt(t, "0x4718f5a0fc34cc1af16a1cdee98ffb20c31f5cd61d6ab07201858f4287c938d"),
				Data: []*felt.Felt{
					utils.TestHexToFelt(t, "0x3f6f3bc663aedc5285d6013cc3ffcbc4341d86ab488b8b68d297f8258793c41"),
					utils.TestHexToFelt(t, "0x1176a1bd84444c89232ec27754698e5d2e7e1a7f1539f12027f28b23ec9f3d8"),
					utils.TestHexToFelt(t, "0x16d8b4ad4000"),
					utils.TestHexToFelt(t, "0x0"),
				},
				Keys: []*felt.Felt{utils.TestHexToFelt(t, "0x99cd8bde557814842a3121e8ddfd433a539b8c9f14bf31ebf108d12e6196e9")},
			},
			{
				FromAddress: utils.TestHexToFelt(t, "0x4718f5a0fc34cc1af16a1cdee98ffb20c31f5cd61d6ab07201858f4287c938d"),
				Data: []*felt.Felt{
					utils.TestHexToFelt(t, "0x1176a1bd84444c89232ec27754698e5d2e7e1a7f1539f12027f28b23ec9f3d8"),
					utils.TestHexToFelt(t, "0x18ad8494375bc00"),
					utils.TestHexToFelt(t, "0x0"),
					utils.TestHexToFelt(t, "0x18aef21f822fc00"),
					utils.TestHexToFelt(t, "0x0"),
				},
				Keys: []*felt.Felt{utils.TestHexToFelt(t, "0xa9fa878c35cd3d0191318f89033ca3e5501a3d90e21e3cc9256bdd5cd17fdd")},
			},
		},
		ExecutionResources: ExecutionResources{
			ComputationResources: ComputationResources{
				Steps:          615,
				MemoryHoles:    4,
				RangeCheckApps: 19,
			},
		},
	})

	testSet := map[string][]testSetType{
		"mock": {},
		"testnet": {
			{
				TxnHash: utils.TestHexToFelt(t, "0xf2f3d50192637e8d5e817363460c39d3a668fe12f117ecedb9749466d8352b"),
				ExpectedResp: TransactionReceiptWithBlockInfo{
					TransactionReceipt: receiptTxn52767_16,
					BlockNumber:        52767,
					BlockHash:          utils.TestHexToFelt(t, "0x4ae5d52c75e4dea5694f456069f830cfbc7bec70427eee170c3385f751b8564"),
				},
			},
		},
		"mainnet": {},
		"integration": {
			{
				TxnHash: utils.TestHexToFelt(t, "0x49728601e0bb2f48ce506b0cbd9c0e2a9e50d95858aa41463f46386dca489fd"),
				ExpectedResp: TransactionReceiptWithBlockInfo{
					TransactionReceipt: receiptTxnIntegration,
					BlockNumber:        319132,
					BlockHash:          utils.TestHexToFelt(t, "0x50e864db6b81ce69fbeb70e6a7284ee2febbb9a2e707415de7adab83525e9cd"),
				},
			},
		}}[testEnv]

	for _, test := range testSet {
		spy := NewSpy(testConfig.provider.c)
		testConfig.provider.c = spy
		txReceiptWithBlockInfo, err := testConfig.provider.TransactionReceipt(context.Background(), test.TxnHash)
		require.Nil(t, err)
		require.Equal(t, txReceiptWithBlockInfo.BlockNumber, test.ExpectedResp.BlockNumber)
		require.Equal(t, txReceiptWithBlockInfo.BlockHash, test.ExpectedResp.BlockHash)

	}
}

// TestGetTransactionStatus tests starknet_getTransactionStatus
func TestGetTransactionStatus(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		TxnHash      *felt.Felt
		ExpectedResp TxnStatusResp
	}

	testSet := map[string][]testSetType{
		"mock": {},
		"testnet": {
			{
				TxnHash:      utils.TestHexToFelt(t, "0xd109474cd037bad60a87ba0ccf3023d5f2d1cd45220c62091d41a614d38eda"),
				ExpectedResp: TxnStatusResp{FinalityStatus: TxnStatus_Accepted_On_L1, ExecutionStatus: TxnExecutionStatusSUCCEEDED},
			},
		},
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		resp, err := testConfig.provider.GetTransactionStatus(context.Background(), test.TxnHash)
		require.Nil(t, err)
		require.Equal(t, *resp, test.ExpectedResp)
	}
}
