package rpc

import (
	"context"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestTransactionByHash tests transaction by hash
//
// Parameters:
//   - t: the testing object for running the test cases
//
// Returns:
// none
func TestTransactionByHash(t *testing.T) {
	testConfig := beforeEach(t, false)

	type testSetType struct {
		TxHash      *felt.Felt
		ExpectedTxn BlockTransaction
	}

	BlockDeclareTxnV2Example := BlockTransaction{
		Hash: internalUtils.TestHexToFelt(t, "0xd109474cd037bad60a87ba0ccf3023d5f2d1cd45220c62091d41a614d38eda"),
		Transaction: DeclareTxnV2{
			Type:              TransactionType_Declare,
			Version:           TransactionV2,
			MaxFee:            internalUtils.TestHexToFelt(t, "0x4a0fbb2d7a43"),
			ClassHash:         internalUtils.TestHexToFelt(t, "0x79b7ec8fdf40a4ff6ed47123049dfe36b5c02db93aa77832682344775ef70c6"),
			CompiledClassHash: internalUtils.TestHexToFelt(t, "0x7130f75fc2f1400813d1e96ea7ebee334b568a87b645a62aade0eb2fa2cf252"),
			Nonce:             internalUtils.TestHexToFelt(t, "0x16e"),
			Signature: []*felt.Felt{
				internalUtils.TestHexToFelt(t, "0x5569787df42fece1184537b0d480900a403386355b9d6a59e7c7a7e758287f0"),
				internalUtils.TestHexToFelt(t, "0x2acaeea2e0817da33ed5dbeec295b0177819b5a5a50b0a669e6eecd88e42e92"),
			},
			SenderAddress: internalUtils.TestHexToFelt(t, "0x5fd4befee268bf6880f955875cbed3ade8346b1f1e149cc87b317e62b6db569"),
		},
	}

	testSet := map[string][]testSetType{
		"mock": {
			{
				TxHash:      internalUtils.TestHexToFelt(t, "0xd109474cd037bad60a87ba0ccf3023d5f2d1cd45220c62091d41a614d38eda"),
				ExpectedTxn: BlockDeclareTxnV2Example,
			},
		},
		"testnet": {
			{
				TxHash:      internalUtils.TestHexToFelt(t, "0xd109474cd037bad60a87ba0ccf3023d5f2d1cd45220c62091d41a614d38eda"),
				ExpectedTxn: BlockDeclareTxnV2Example,
			},
		},
		"mainnet": {},
	}[testEnv]
	for _, test := range testSet {
		tx, err := testConfig.provider.TransactionByHash(context.Background(), test.TxHash)
		require.NoError(t, err)
		require.NotNil(t, tx)

		assert.Equal(t, test.ExpectedTxn, *tx)
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
//   - t: the testing object for running the test cases
//
// Returns:
//
//	none
func TestTransactionByBlockIdAndIndex(t *testing.T) {
	testConfig := beforeEach(t, false)

	type testSetType struct {
		BlockID     BlockID
		Index       uint64
		ExpectedTxn BlockTransaction
	}

	InvokeTxnV3example := *internalUtils.TestUnmarshalJSONFileToType[BlockTransaction](t, "./tests/transactions/sepoliaBlockInvokeTxV3_0x265f6a59e7840a4d52cec7db37be5abd724fdfd72db9bf684f416927a88bc89.json", "")

	testSet := map[string][]testSetType{
		"mock": {
			{
				BlockID:     WithBlockHash(internalUtils.TestHexToFelt(t, "0x873a3d4e1159ccecec5488e07a31c9a4ba8c6d2365b6aa48d39f5fd54e6bd0")),
				Index:       0,
				ExpectedTxn: InvokeTxnV3example,
			},
		},
		"testnet": {
			{
				BlockID:     WithBlockHash(internalUtils.TestHexToFelt(t, "0x873a3d4e1159ccecec5488e07a31c9a4ba8c6d2365b6aa48d39f5fd54e6bd0")),
				Index:       3,
				ExpectedTxn: InvokeTxnV3example,
			},
		},
		"mainnet": {},
	}[testEnv]
	for _, test := range testSet {
		tx, err := testConfig.provider.TransactionByBlockIdAndIndex(context.Background(), test.BlockID, test.Index)
		require.NoError(t, err)
		require.NotNil(t, tx)
		assert.Equal(t, test.ExpectedTxn, *tx)
	}
}

func TestTransactionReceipt(t *testing.T) {
	testConfig := beforeEach(t, false)

	type testSetType struct {
		TxnHash      *felt.Felt
		ExpectedResp TransactionReceiptWithBlockInfo
	}

	receiptTxn52767_16 := *internalUtils.TestUnmarshalJSONFileToType[TransactionReceiptWithBlockInfo](t, "./tests/receipt/sepoliaRec_0xf2f3d50192637e8d5e817363460c39d3a668fe12f117ecedb9749466d8352b.json", "")

	// https://voyager.online/tx/0x74011377f326265f5a54e27a27968355e7033ad1de11b77b225374875aff519
	receiptL1Handler := *internalUtils.TestUnmarshalJSONFileToType[TransactionReceiptWithBlockInfo](t, "./tests/receipt/mainnetRc_0x74011377f326265f5a54e27a27968355e7033ad1de11b77b225374875aff519.json", "")

	testSet := map[string][]testSetType{
		"mock": {
			{
				TxnHash:      internalUtils.TestHexToFelt(t, "0xf2f3d50192637e8d5e817363460c39d3a668fe12f117ecedb9749466d8352b"),
				ExpectedResp: receiptTxn52767_16,
			},
			{
				TxnHash:      internalUtils.TestHexToFelt(t, "0x74011377f326265f5a54e27a27968355e7033ad1de11b77b225374875aff519"),
				ExpectedResp: receiptL1Handler,
			},
		},
		"testnet": {
			{
				TxnHash:      internalUtils.TestHexToFelt(t, "0xf2f3d50192637e8d5e817363460c39d3a668fe12f117ecedb9749466d8352b"),
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

// TestGetTransactionStatus tests starknet_getTransactionStatus in the GetTransactionStatus function
func TestGetTransactionStatus(t *testing.T) {
	testConfig := beforeEach(t, false)

	type testSetType struct {
		TxnHash      *felt.Felt
		ExpectedResp TxnStatusResult
	}

	testSet := map[string][]testSetType{
		"mock": {},
		"testnet": {
			{
				TxnHash:      internalUtils.TestHexToFelt(t, "0xd109474cd037bad60a87ba0ccf3023d5f2d1cd45220c62091d41a614d38eda"),
				ExpectedResp: TxnStatusResult{FinalityStatus: TxnStatus_Accepted_On_L1, ExecutionStatus: TxnExecutionStatusSUCCEEDED},
			},
			{
				TxnHash: internalUtils.TestHexToFelt(t, "0x5adf825a4b7fc4d2d99e65be934bd85c83ca2b9383f2ff28fc2a4bc2e6382fc"),
				ExpectedResp: TxnStatusResult{
					FinalityStatus:  TxnStatus_Accepted_On_L1,
					ExecutionStatus: TxnExecutionStatusREVERTED,
					FailureReason:   "Transaction execution has failed:\n0: Error in the called contract (contract address: 0x036d67ab362562a97f9fba8a1051cf8e37ff1a1449530fb9f1f0e32ac2da7d06, class hash: 0x061dac032f228abef9c6626f995015233097ae253a7f72d68552db02f2971b8f, selector: 0x015d40a3d6ca2ac30f4031e42be28da9b056fef9bb7357ac5e85627ee876e5ad):\nError at pc=0:4835:\nCairo traceback (most recent call last):\nUnknown location (pc=0:67)\nUnknown location (pc=0:1835)\nUnknown location (pc=0:2554)\nUnknown location (pc=0:3436)\nUnknown location (pc=0:4040)\n\n1: Error in the called contract (contract address: 0x00000000000000000000000000000000000000000000000000000000ffffffff, class hash: 0x0000000000000000000000000000000000000000000000000000000000000000, selector: 0x02f0b3c5710379609eb5495f1ecd348cb28167711b73609fe565a72734550354):\nRequested contract address 0x00000000000000000000000000000000000000000000000000000000ffffffff is not deployed.\n",
				},
			},
		},
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		resp, err := testConfig.provider.GetTransactionStatus(context.Background(), test.TxnHash)
		require.Nil(t, err)
		require.Equal(t, resp.FinalityStatus, test.ExpectedResp.FinalityStatus)
		require.Equal(t, resp.ExecutionStatus, test.ExpectedResp.ExecutionStatus)
		require.Equal(t, resp.FailureReason, test.ExpectedResp.FailureReason)
	}
}

// TestGetMessagesStatus tests starknet_getMessagesStatus in the GetMessagesStatus function
func TestGetMessagesStatus(t *testing.T) {
	testConfig := beforeEach(t, false)

	type testSetType struct {
		TxHash       NumAsHex
		ExpectedResp []MessageStatus
		ExpectedErr  error
	}

	testSet := map[string][]testSetType{
		"mock": {
			{
				TxHash: "0x123",
				ExpectedResp: []MessageStatus{
					{
						Hash:           internalUtils.RANDOM_FELT,
						FinalityStatus: TxnStatus_Accepted_On_L2,
					},
					{
						Hash:           internalUtils.RANDOM_FELT,
						FinalityStatus: TxnStatus_Accepted_On_L2,
					},
				},
			},
			{
				TxHash:      "0xdededededededededededededededededededededededededededededededede",
				ExpectedErr: ErrHashNotFound,
			},
		},
		"testnet": {
			{
				TxHash: "0x06c5ca541e3d6ce35134e1de3ed01dbf106eaa770d92744432b497f59fddbc00",
				ExpectedResp: []MessageStatus{
					{
						Hash:           internalUtils.TestHexToFelt(t, "0x71660e0442b35d307fc07fa6007cf2ae4418d29fd73833303e7d3cfe1157157"),
						FinalityStatus: TxnStatus_Accepted_On_L1,
					},
					{
						Hash:           internalUtils.TestHexToFelt(t, "0x28a3d1f30922ab86bb240f7ce0f5e8cbbf936e5d2fcfe52b8ffbe71e341640"),
						FinalityStatus: TxnStatus_Accepted_On_L1,
					},
				},
			},
			{
				TxHash:      "0xdededededededededededededededededededededededededededededededede",
				ExpectedErr: ErrHashNotFound,
			},
		},
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		resp, err := testConfig.provider.GetMessagesStatus(context.Background(), test.TxHash)
		if test.ExpectedErr != nil {
			require.EqualError(t, err, test.ExpectedErr.Error())
		} else {
			require.Nil(t, err)
			require.Equal(t, test.ExpectedResp, resp)
		}
	}
}
