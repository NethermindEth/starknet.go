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
// If the transaction is not of type InvokeTxn3, it fails the test. Finally, it asserts
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

	var InvokeTxnV3example = InvokeTxnV3{
		Type:    TransactionType_Invoke,
		Version: TransactionV3,
		Nonce:   utils.TestHexToFelt(t, "0x359d"),
		Signature: []*felt.Felt{
			utils.TestHexToFelt(t, "0x665f0c67ed4d9565f63857b1a55974b98b2411f579c53c9f903fd21a3edb3d1"),
			utils.TestHexToFelt(t, "0x549c4480aba4753c58f757c92b5a1d3d67b2ced4bf06076825af3f52f738d6d"),
		},
		SenderAddress:         utils.TestHexToFelt(t, "0x143fe26927dd6a302522ea1cd6a821ab06b3753194acee38d88a85c93b3cbc6"),
		NonceDataMode:         DAModeL1,
		FeeMode:               DAModeL1,
		PayMasterData:         []*felt.Felt{},
		AccountDeploymentData: []*felt.Felt{},
		ResourceBounds: ResourceBoundsMapping{
			L1Gas: ResourceBounds{
				MaxAmount:       "0x3bb2",
				MaxPricePerUnit: "0x2ba7def30000",
			},
			L2Gas: ResourceBounds{
				MaxAmount:       "0x0",
				MaxPricePerUnit: "0x0",
			},
		},
		Calldata: []*felt.Felt{
			utils.TestHexToFelt(t, "0x1"),
			utils.TestHexToFelt(t, "0x6b74c515944ef1ef630ee1cf08a22e110c39e217fa15554a089182a11f78ed"),
			utils.TestHexToFelt(t, "0xc844fd57777b0cd7e75c8ea68deec0adf964a6308da7a58de32364b7131cc8"),
			utils.TestHexToFelt(t, "0x13"),
			utils.TestHexToFelt(t, "0x41bbf1eff2ac123d9e01004a385329369cbc1c309838562f030b3faa2caa4"),
			utils.TestHexToFelt(t, "0x54103"),
			utils.TestHexToFelt(t, "0x7e430a7a59836b5969859b25379c640a8ccb66fb142606d7acb1a5563c2ad9"),
			utils.TestHexToFelt(t, "0x6600d829"),
			utils.TestHexToFelt(t, "0x103020400000000000000000000000000000000000000000000000000000000"),
			utils.TestHexToFelt(t, "0x4"),
			utils.TestHexToFelt(t, "0x5f5e100"),
			utils.TestHexToFelt(t, "0x5f60fc2"),
			utils.TestHexToFelt(t, "0x5f60fc2"),
			utils.TestHexToFelt(t, "0x5f6570d"),
			utils.TestHexToFelt(t, "0xa07695b6574c60c37"),
			utils.TestHexToFelt(t, "0x1"),
			utils.TestHexToFelt(t, "0x2"),
			utils.TestHexToFelt(t, "0x7afe11c6cdf846e8e33ff55c6e8310293b81aa58da4618af0c2fb29db2515c7"),
			utils.TestHexToFelt(t, "0x1200966b0f9a5cd1bf7217b202c3a4073a1ff583e4779a3a3ffb97a532fe0c"),
			utils.TestHexToFelt(t, "0x2cb74dff29a13dd5d855159349ec92f943bacf0547ff3734e7d84a15d08cbc5"),
			utils.TestHexToFelt(t, "0x460769330eab4b3269a5c07369382fcc09fbfc92458c63f77292425c72272f9"),
			utils.TestHexToFelt(t, "0x10ebdb197fd1017254b927b01073c64a368db45534413b539895768e57b72ba"),
			utils.TestHexToFelt(t, "0x2e7dc996ebf724c1cf18d668fc3455df4245749ebc0724101cbc6c9cb13c962"),
		},
		Tip: "0x0",
	}

	testSet := map[string][]testSetType{
		"mock": {
			{
				BlockID:     WithBlockHash(utils.TestHexToFelt(t, "0x4ae5d52c75e4dea5694f456069f830cfbc7bec70427eee170c3385f751b8564")),
				Index:       0,
				ExpectedTxn: InvokeTxnV3example,
			},
		},
		"testnet": {
			{
				BlockID:     WithBlockHash(utils.TestHexToFelt(t, "0x4ae5d52c75e4dea5694f456069f830cfbc7bec70427eee170c3385f751b8564")),
				Index:       15,
				ExpectedTxn: InvokeTxnV3example,
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
		txCasted, ok := (tx).(InvokeTxnV3)
		if !ok {
			t.Fatalf("transaction should be InvokeTxnV3, instead %T", tx)
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
				FromAddress: utils.TestHexToFelt(t, "0x4718f5a0fc34cc1af16a1cdee98ffb20c31f5cd61d6ab07201858f4287c938d"),
				Data: []*felt.Felt{
					utils.TestHexToFelt(t, "0x6016d919abf2ddefe03dacc2ff5c8f42eb80cf65add1e90dd73c5c5e06ef3e2"),
					utils.TestHexToFelt(t, "0x1176a1bd84444c89232ec27754698e5d2e7e1a7f1539f12027f28b23ec9f3d8"),
					utils.TestHexToFelt(t, "0x16409a78a10b00"),
					utils.TestHexToFelt(t, "0x0"),
				},
				Keys: []*felt.Felt{
					utils.TestHexToFelt(t, "0x99cd8bde557814842a3121e8ddfd433a539b8c9f14bf31ebf108d12e6196e9"),
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
			DataAvailability: DataAvailability{
				L1Gas:     0,
				L1DataGas: 128,
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
		"mock": {
			{
				TxnHash: utils.TestHexToFelt(t, "0xf2f3d50192637e8d5e817363460c39d3a668fe12f117ecedb9749466d8352b"),
				ExpectedResp: TransactionReceiptWithBlockInfo{
					TransactionReceipt: receiptTxn52767_16,
					BlockNumber:        52767,
					BlockHash:          utils.TestHexToFelt(t, "0x4ae5d52c75e4dea5694f456069f830cfbc7bec70427eee170c3385f751b8564"),
				},
			},
		},
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
		require.Equal(t, test.ExpectedResp, txReceiptWithBlockInfo)
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
