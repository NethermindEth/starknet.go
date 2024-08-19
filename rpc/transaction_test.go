package rpc

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/utils"
	"github.com/stretchr/testify/require"
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
		ExpectedTxn BlockTransaction
	}

	var BlockDeclareTxnV2Example = BlockTransaction{
		BlockDeclareTxnV2{
			utils.TestHexToFelt(t, "0xd109474cd037bad60a87ba0ccf3023d5f2d1cd45220c62091d41a614d38eda"),
			DeclareTxnV2{
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
			},
		},
	}

	testSet := map[string][]testSetType{
		"mock": {
			{
				TxHash:      utils.TestHexToFelt(t, "0xd109474cd037bad60a87ba0ccf3023d5f2d1cd45220c62091d41a614d38eda"),
				ExpectedTxn: BlockDeclareTxnV2Example,
			},
		},
		"testnet": {
			{
				TxHash:      utils.TestHexToFelt(t, "0xd109474cd037bad60a87ba0ccf3023d5f2d1cd45220c62091d41a614d38eda"),
				ExpectedTxn: BlockDeclareTxnV2Example,
			},
		},
		"mainnet": {},
	}[testEnv]
	for _, test := range testSet {
		tx, err := testConfig.provider.TransactionByHash(context.Background(), test.TxHash)
		require.NoError(t, err)
		require.NotNil(t, tx)

		txCasted, ok := (tx.IBlockTransaction).(BlockDeclareTxnV2)
		require.True(t, ok)
		require.Equal(t, test.ExpectedTxn.IBlockTransaction, txCasted)
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
		ExpectedTxn BlockTransaction
	}

	var InvokeTxnV3example BlockTransaction
	read, err := os.ReadFile("tests/transactions/sepoliaTx_0x6a4a9c4f1a530f7d6dd7bba9b71f090a70d1e3bbde80998fde11a08aab8b282.json")
	require.NoError(t, err)
	err = json.Unmarshal(read, &InvokeTxnV3example)
	require.NoError(t, err)

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

		tx, err := testConfig.provider.TransactionByBlockIdAndIndex(context.Background(), test.BlockID, test.Index)
		require.NoError(t, err)
		require.NotNil(t, tx)
		txCasted, ok := (tx.IBlockTransaction).(BlockInvokeTxnV3)
		require.True(t, ok)
		require.Equal(t, test.ExpectedTxn.IBlockTransaction, txCasted)
	}
}

func TestTransactionReceipt(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		TxnHash      *felt.Felt
		ExpectedResp TransactionReceiptWithBlockInfo
	}
	var receiptTxn52767_16 TransactionReceiptWithBlockInfo
	read, err := os.ReadFile("tests/receipt/sepoliaRec_0xf2f3d50192637e8d5e817363460c39d3a668fe12f117ecedb9749466d8352b.json")
	require.NoError(t, err)
	err = json.Unmarshal(read, &receiptTxn52767_16)
	require.NoError(t, err)

	// // https://voyager.online/tx/0x74011377f326265f5a54e27a27968355e7033ad1de11b77b225374875aff519
	var receiptL1Handler TransactionReceiptWithBlockInfo
	read, err = os.ReadFile("tests/receipt/mainnetRc_0x74011377f326265f5a54e27a27968355e7033ad1de11b77b225374875aff519.json")
	require.NoError(t, err)
	err = json.Unmarshal(read, &receiptL1Handler)
	require.NoError(t, err)

	testSet := map[string][]testSetType{
		"mock": {
			{
				TxnHash:      utils.TestHexToFelt(t, "0xf2f3d50192637e8d5e817363460c39d3a668fe12f117ecedb9749466d8352b"),
				ExpectedResp: receiptTxn52767_16,
			},
			{
				TxnHash:      utils.TestHexToFelt(t, "0x74011377f326265f5a54e27a27968355e7033ad1de11b77b225374875aff519"),
				ExpectedResp: receiptL1Handler,
			},
		},
		"testnet": {
			{
				TxnHash:      utils.TestHexToFelt(t, "0xf2f3d50192637e8d5e817363460c39d3a668fe12f117ecedb9749466d8352b"),
				ExpectedResp: receiptTxn52767_16,
			},
		},
		"mainnet":     {},
		"integration": {},
	}[testEnv]

	for _, test := range testSet {
		spy := NewSpy(testConfig.provider.c)
		testConfig.provider.c = spy
		txReceiptWithBlockInfo, err := testConfig.provider.TransactionReceipt(context.Background(), test.TxnHash)
		require.Nil(t, err)
		require.Equal(t, test.ExpectedResp, *txReceiptWithBlockInfo)
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
