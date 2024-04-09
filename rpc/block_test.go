package rpc

import (
	"context"
	"strings"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/utils"
	"github.com/google/go-cmp/cmp"
	"github.com/test-go/testify/require"
)

// TestBlockNumber is a test function to check the behavior of the BlockNumber function and check the returned value is strictly positive.
//
// The function performs the following steps:
// 1. Sets up the test configuration.
// 2. Defines a test set.
// 3. Loops over the test set.
// 4. Creates a new spy.
// 5. Calls the BlockNumber function on the test provider.
// 6. Validates the returned block number.
//
// Parameters:
// - t: the testing object for running the test cases
// Returns:
//
//	none
func TestBlockNumber(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct{}

	testSet := map[string][]testSetType{
		"mock":    {},
		"testnet": {{}},
		"mainnet": {{}},
		"devnet":  {},
	}[testEnv]

	for range testSet {
		spy := NewSpy(testConfig.provider.c)
		testConfig.provider.c = spy
		blockNumber, err := testConfig.provider.BlockNumber(context.Background())
		if err != nil {
			t.Fatal("BlockWithTxHashes match the expected error:", err)
		}
		if diff, err := spy.Compare(blockNumber, false); err != nil || diff != "FullMatch" {
			t.Fatal("expecting to match", err)
		}
		if blockNumber <= 3000 {
			t.Fatal("Block number should be > 3000, instead: ", blockNumber)
		}
	}
}

// TestBlockHashAndNumber is a test function that tests the BlockHashAndNumber function and check the returned value is strictly positive.
//
// It sets up a test configuration and creates a test set based on the test environment.
// Then it iterates through the test set and performs the following steps:
//   - Creates a new spy using the testConfig provider.
//   - Sets the testConfig provider to the spy.
//   - Calls the BlockHashAndNumber function of the testConfig provider with a context.
//   - Checks if there is an error and if it matches the expected error.
//   - Compares the result with the spy and checks if it matches the expected result.
//   - Checks if the block number is greater than 3000.
//   - Checks if the block hash starts with "0x".
//
// Parameters:
// - t: the testing object for running the test cases
// Returns:
//
//	none
func TestBlockHashAndNumber(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct{}

	testSet := map[string][]testSetType{
		"mock":    {},
		"testnet": {{}},
		"mainnet": {{}},
		"devnet":  {},
	}[testEnv]

	for range testSet {
		spy := NewSpy(testConfig.provider.c)
		testConfig.provider.c = spy
		blockHashAndNumber, err := testConfig.provider.BlockHashAndNumber(context.Background())
		if err != nil {
			t.Fatal("BlockHashAndNumber match the expected error:", err)
		}
		if diff, err := spy.Compare(blockHashAndNumber, false); err != nil || diff != "FullMatch" {
			t.Fatal("expecting to match", err)
		}
		if blockHashAndNumber.BlockNumber < 3000 {
			t.Fatal("Block number should be > 3000, instead: ", blockHashAndNumber.BlockNumber)
		}
		if !strings.HasPrefix(blockHashAndNumber.BlockHash.String(), "0x") {
			t.Fatal("current block hash should return a string starting with 0x")
		}
	}
}

// TestBlockWithTxHashes tests the functionality of the BlockWithTxHashes function.
//
// The function takes a testing.T object as a parameter and initializes a testConfig object.
// It defines a testSetType struct that contains several fields including BlockID, ExpectedError, ExpectedBlockWithTxHashes, and ExpectedPendingBlockWithTxHashes.
// The function then initializes a blockSepolia30436 variable of type BlockTxHashes with a predefined set of values.
// It also initializes a txHashes variable of type []felt.Felt and a blockHash variable of type felt.Felt.
//
// The function defines a testSet map that has three keys: "mock", "testnet", and "mainnet".
// Each key corresponds to a slice of testSetType objects.
// The "mock" key has two testSetType objects with different field values.
// The "testnet" key has three testSetType objects with different field values.
// The "mainnet" key does not have any testSetType objects.
//
// The function then iterates over the testSet map and performs the following steps for each testSetType object:
// - It creates a new Spy object and assigns it to the testConfig.provider.c field.
// - It calls the BlockWithTxHashes function with the provided BlockID and stores the result in the result variable.
// - It checks if the returned error matches the expected error. If not, it calls the Fatal function of the testing.T object with an error message.
// - It checks the type of the result variable and performs specific assertions based on the type.
//   - If the result is of type *BlockTxHashes, it checks various fields of the BlockTxHashes object against the expected values.
//   - If the result is of type *PendingBlockTxHashes, it checks various fields of the PendingBlockTxHashes object against the expected values.
//   - If the result is of any other type, it calls the Fatal function of the testing.T object with an error message.
//
// Parameters:
// - t: the testing object for running the test cases
// Returns:
//
//	none
func TestBlockWithTxHashes(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		BlockID                          BlockID
		ExpectedErr                      error
		ExpectedBlockWithTxHashes        *BlockTxHashes
		ExpectedPendingBlockWithTxHashes *PendingBlockTxHashes
	}

	var blockSepolia30436 = BlockTxHashes{
		BlockHeader: BlockHeader{
			BlockHash:        utils.TestHexToFelt(t, "0x2488a7149327b4dcd200f05a131911bd44f946021539648642eaa7d6e82f289"),
			ParentHash:       utils.TestHexToFelt(t, "0x2adc07a26d70e72a16775e26c45074f0216bc2e86e35bfe53743968480e4c1b"),
			SequencerAddress: utils.TestHexToFelt(t, "0x1176a1bd84444c89232ec27754698e5d2e7e1a7f1539f12027f28b23ec9f3d8"),
			BlockNumber:      30436,
			NewRoot:          utils.TestHexToFelt(t, "0x4bbb6bec9488d70d9c9e96862cf50e22331a5e8d7b33a56712f56cd04c16e06"),
			Timestamp:        1707158969,
		},
		Status: "ACCEPTED_ON_L1",
		Transactions: utils.TestHexArrToFelt(t, []string{
			"0x10d2059db6d261fee740b515ed8b9c50955f03dd43c4729b24dc63278641926",
			"0x50d39b3fd4931343aeb6bf325493f7b8c05d8bf2819e4dad465e42751e2412",
			"0x3761dfe1ce22d16eab6339d8ba2ab4c56008182d979979f8e68217920193996",
			"0x7e3f4756d90b1f6f8249185857b4beab0c0dd3d3b207ad73fd249c6267ecea5",
		}),
	}

	txHashes := utils.TestHexArrToFelt(t, []string{
		"0x10d2059db6d261fee740b515ed8b9c50955f03dd43c4729b24dc63278641926",
		"0x50d39b3fd4931343aeb6bf325493f7b8c05d8bf2819e4dad465e42751e2412",
		"0x3761dfe1ce22d16eab6339d8ba2ab4c56008182d979979f8e68217920193996",
	})
	blockHash := utils.TestHexToFelt(t, "0xbeef")

	testSet := map[string][]testSetType{
		"mock": {
			{
				BlockID:     BlockID{Tag: "latest"},
				ExpectedErr: nil,
				ExpectedPendingBlockWithTxHashes: &PendingBlockTxHashes{
					PendingBlockHeader{
						ParentHash:       &felt.Zero,
						Timestamp:        123,
						SequencerAddress: &felt.Zero},
					txHashes,
				},
			},
			{
				BlockID: BlockID{Hash: blockHash},
				ExpectedBlockWithTxHashes: &BlockTxHashes{
					BlockHeader: BlockHeader{
						BlockHash:        blockHash,
						ParentHash:       &felt.Zero,
						Timestamp:        124,
						SequencerAddress: &felt.Zero},
					Status:       BlockStatus_AcceptedOnL1,
					Transactions: txHashes,
				},
			},
		},
		"testnet": {
			{
				BlockID:     WithBlockTag("latest"),
				ExpectedErr: nil,
			},
			{
				BlockID:                   WithBlockHash(utils.TestHexToFelt(t, "0x2488a7149327b4dcd200f05a131911bd44f946021539648642eaa7d6e82f289")),
				ExpectedErr:               nil,
				ExpectedBlockWithTxHashes: &blockSepolia30436,
			},
			{
				BlockID:                   WithBlockNumber(30436),
				ExpectedErr:               nil,
				ExpectedBlockWithTxHashes: &blockSepolia30436,
			},
		},
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		spy := NewSpy(testConfig.provider.c)
		testConfig.provider.c = spy
		result, err := testConfig.provider.BlockWithTxHashes(context.Background(), test.BlockID)
		require.Equal(t, test.ExpectedErr, err, "Error in BlockWithTxHashes")
		switch resultType := result.(type) {
		case *BlockTxHashes:
			block, ok := result.(*BlockTxHashes)
			if !ok {
				t.Fatalf("should return *BlockTxHashes, instead: %T\n", result)
			}
			if test.ExpectedErr != nil {
				continue
			}
			if !strings.HasPrefix(block.BlockHash.String(), "0x") {
				t.Fatal("Block Hash should start with \"0x\", instead", block.BlockHash)
			}

			if len(block.Transactions) == 0 {
				t.Fatal("the number of transaction should not be 0")
			}
			if test.ExpectedBlockWithTxHashes != nil {
				if (*test.ExpectedBlockWithTxHashes).BlockHash == &felt.Zero {
					continue
				}

				require.Equal(t, block.BlockHeader.BlockHash, test.ExpectedBlockWithTxHashes.BlockHeader.BlockHash, "Error in BlockTxHash BlockHash")
				require.Equal(t, block.BlockHeader.ParentHash, test.ExpectedBlockWithTxHashes.BlockHeader.ParentHash, "Error in BlockTxHash ParentHash")
				require.Equal(t, block.BlockHeader.Timestamp, test.ExpectedBlockWithTxHashes.BlockHeader.Timestamp, "Error in BlockTxHash Timestamp")
				require.Equal(t, block.BlockHeader.SequencerAddress, test.ExpectedBlockWithTxHashes.BlockHeader.SequencerAddress, "Error in BlockTxHash SequencerAddress")
				require.Equal(t, block.Status, test.ExpectedBlockWithTxHashes.Status, "Error in BlockTxHash Status")
				require.Equal(t, block.Transactions, test.ExpectedBlockWithTxHashes.Transactions, "Error in BlockTxHash Transactions")
			}
		case *PendingBlockTxHashes:
			pBlock, ok := result.(*PendingBlockTxHashes)
			if !ok {
				t.Fatalf("should return *PendingBlockTxHashes, instead: %T\n", result)
			}

			require.Equal(t, pBlock.ParentHash, test.ExpectedPendingBlockWithTxHashes.ParentHash, "Error in PendingBlockTxHashes ParentHash")
			require.Equal(t, pBlock.SequencerAddress, test.ExpectedPendingBlockWithTxHashes.SequencerAddress, "Error in PendingBlockTxHashes SequencerAddress")
			require.Equal(t, pBlock.Timestamp, test.ExpectedPendingBlockWithTxHashes.Timestamp, "Error in PendingBlockTxHashes Timestamp")
			require.Equal(t, pBlock.Transactions, test.ExpectedPendingBlockWithTxHashes.Transactions, "Error in PendingBlockTxHashes Transactions")
		default:
			t.Fatalf("unexpected block type, found: %T\n", resultType)
		}
	}
}

// TestBlockWithTxsAndInvokeTXNV0 tests the BlockWithTxsAndInvokeTXNV0 function.
//
// The function tests the BlockWithTxsAndInvokeTXNV0 function by setting up a test configuration and a test set type.
// It then initializes a fullBlockSepolia30436 variable with a Block struct and invokes the BlockWithTxs function with different test scenarios.
// The function compares the expected error with the actual error and checks if the BlockWithTxs function returns the correct block data.
// It also verifies the block hash, the number of transactions in the block, and the details of a specific transaction.
//
// Parameters:
// - t: The t testing object
// Returns:
//
//	none
//
// Todo: Find a V0 transaction on SN_SEPOLIA
// func TestBlockWithTxsAndInvokeTXNV0(t *testing.T) {
// 	testConfig := beforeEach(t)

// 	type testSetType struct {
// 		BlockID                     BlockID
// 		ExpectedError               error
// 		LookupTxnPositionInOriginal int
// 		LookupTxnPositionInExpected int
// 		want                        *Block
// 	}

// 	var fullBlockSepolia30436 = Block{
// 		BlockHeader: BlockHeader{
// 			BlockHash:        utils.TestHexToFelt(t, "0x2488a7149327b4dcd200f05a131911bd44f946021539648642eaa7d6e82f289"),
// 			ParentHash:       utils.TestHexToFelt(t, "0x2adc07a26d70e72a16775e26c45074f0216bc2e86e35bfe53743968480e4c1b"),
// 			SequencerAddress: utils.TestHexToFelt(t, "0x1176a1bd84444c89232ec27754698e5d2e7e1a7f1539f12027f28b23ec9f3d8"),
// 			BlockNumber:      30436,
// 			NewRoot:          utils.TestHexToFelt(t, "0x4bbb6bec9488d70d9c9e96862cf50e22331a5e8d7b33a56712f56cd04c16e06"),
// 			Timestamp:        1661450764,
// 		},
// 		Status: "ACCEPTED_ON_L1",
// 		Transactions: []BlockTransaction{

// 			BlockInvokeTxnV0{
// 				TransactionHash: utils.TestHexToFelt(t, "0x10d2059db6d261fee740b515ed8b9c50955f03dd43c4729b24dc63278641926"),
// 				InvokeTxnV0: InvokeTxnV0{
// 					Type:    "INVOKE",
// 					MaxFee:  utils.TestHexToFelt(t, "0x470de4df820000"),
// 					Version: TransactionV0,
// 					Signature: []*felt.Felt{
// 						utils.TestHexToFelt(t, "0x7fb440b1dee35c5259bf10d55782bc973434d195bb5c8b95ac7d3d8e2a8a0e4"),
// 						utils.TestHexToFelt(t, "0x13089e5f38de4ea98e9275be7fadc915946be15c14a8fed7c55202818527bea"),
// 					},
// 					FunctionCall: FunctionCall{
// 						ContractAddress:    nil,
// 						EntryPointSelector: nil,
// 						Calldata: []*felt.Felt{
// 							utils.TestHexToFelt(t, "0x1"),
// 							utils.TestHexToFelt(t, "0x3fe8e4571772bbe0065e271686bd655efd1365a5d6858981e582f82f2c10313"),
// 							utils.TestHexToFelt(t, "0x2468d193cd15b621b24c2a602b8dbcfa5eaa14f88416c40c09d7fd12592cb4b"),
// 							utils.TestHexToFelt(t, "0x0"),
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}

// 	testSet := map[string][]testSetType{
// 		"mock": {},
// 		"testnet": {
// 			{
// 				BlockID:       WithBlockTag("latest"),
// 				ExpectedError: nil,
// 			},
// 			{
// 				BlockID:       WithBlockHash(utils.TestHexToFelt(t, "0x2488a7149327b4dcd200f05a131911bd44f946021539648642eaa7d6e82f289")),
// 				ExpectedError: nil,
// 				want:          &fullBlockSepolia30436,
// 			},
// 			{
// 				BlockID:       WithBlockNumber(30436),
// 				ExpectedError: nil,
// 				want:          &fullBlockSepolia30436,
// 			},
// 		},
// 		"mainnet": {},
// 	}[testEnv]

// 	for i, test := range testSet {
// 		fmt.Println("Test case", i)
// 		fmt.Println("BlockID is", test.BlockID)
// 		fmt.Println("LookuptxnPositionInOriginal is", test.LookupTxnPositionInOriginal)
// 		spy := NewSpy(testConfig.provider.c)
// 		testConfig.provider.c = spy
// 		blockWithTxsInterface, err := testConfig.provider.BlockWithTxs(context.Background(), test.BlockID)
// 		if err != test.ExpectedError {
// 			t.Fatal("BlockWithTxHashes match the expected error:", err)
// 		}
// 		if test.ExpectedError != nil && blockWithTxsInterface == nil {
// 			continue
// 		}
// 		blockWithTxs, ok := blockWithTxsInterface.(*Block)
// 		if !ok {
// 			t.Fatalf("expecting *rpv02.Block, instead %T", blockWithTxsInterface)
// 		}
// 		_, err = spy.Compare(blockWithTxs, false)
// 		if err != nil {
// 			t.Fatal("expecting to match", err)
// 		}
// 		if !strings.HasPrefix(blockWithTxs.BlockHash.String(), "0x") {
// 			t.Fatal("Block Hash should start with \"0x\", instead", blockWithTxs.BlockHash)
// 		}

// 		if len(blockWithTxs.Transactions) == 0 {
// 			t.Fatal("the number of transaction should not be 0")
// 		}

// 		if test.want != nil {
// 			if (*test.want).BlockHash == &felt.Zero {
// 				continue
// 			}

// 			invokeV0Want, ok := (*test.want).Transactions[test.LookupTxnPositionInExpected].(BlockInvokeTxnV0)
// 			if !ok {
// 				t.Fatal("expected invoke v0 transaction")
// 			}
// 			invokeV0Block, ok := blockWithTxs.Transactions[test.LookupTxnPositionInOriginal].(BlockInvokeTxnV0)
// 			if !ok {
// 				t.Fatal("expected invoke v0 transaction")
// 			}
// 			require.Equal(t, invokeV0Want.TransactionHash, invokeV0Block.TransactionHash, "expected equal TransactionHash")
// 			require.Equal(t, invokeV0Want.InvokeTxnV0.MaxFee, invokeV0Block.InvokeTxnV0.MaxFee, "expected equal maxfee")
// 			require.Equal(t, invokeV0Want.InvokeTxnV0.EntryPointSelector, invokeV0Block.InvokeTxnV0.EntryPointSelector, "expected equal eps")

// 		}

// 	}
// }

// TestBlockWithTxsAndInvokeTXNV1 tests the BlockWithTxsAndInvokeTXNV1 function.
//
// The function tests the BlockWithTxsAndInvokeTXNV1 function by setting up a test configuration and a test set type.
// It then initializes a fullBlockSepolia30436 variable with a Block struct and invokes the BlockWithTxs function with different test scenarios.
// The function compares the expected error with the actual error and checks if the BlockWithTxs function returns the correct block data.
// It also verifies the block hash, the number of transactions in the block, and the details of a specific transaction.
//
// Parameters:
// - t: The t testing object
// Returns:
//
//	none
func TestBlockWithTxsAndInvokeTXNV1(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		BlockID                     BlockID
		ExpectedError               error
		LookupTxnPositionInOriginal int
		LookupTxnPositionInExpected int
		want                        *Block
	}

	var fullBlockSepolia30436 = Block{
		BlockHeader: BlockHeader{
			BlockHash:        utils.TestHexToFelt(t, "0x10d2059db6d261fee740b515ed8b9c50955f03dd43c4729b24dc63278641926"),
			ParentHash:       utils.TestHexToFelt(t, "0x2adc07a26d70e72a16775e26c45074f0216bc2e86e35bfe53743968480e4c1b"),
			SequencerAddress: utils.TestHexToFelt(t, "0x1176a1bd84444c89232ec27754698e5d2e7e1a7f1539f12027f28b23ec9f3d8"),
			BlockNumber:      30436,
			NewRoot:          utils.TestHexToFelt(t, "0x4bbb6bec9488d70d9c9e96862cf50e22331a5e8d7b33a56712f56cd04c16e06"),
			Timestamp:        1661450764,
		},
		Status: "ACCEPTED_ON_L1",
		Transactions: []BlockTransaction{

			BlockInvokeTxnV1{
				TransactionHash: utils.TestHexToFelt(t, "0x10d2059db6d261fee740b515ed8b9c50955f03dd43c4729b24dc63278641926"),
				InvokeTxnV1: InvokeTxnV1{
					Type:    "INVOKE",
					Nonce:   utils.TestHexToFelt(t, "0x12562"),
					MaxFee:  utils.TestHexToFelt(t, "0x470de4df820000"),
					Version: TransactionV1,
					Signature: []*felt.Felt{
						utils.TestHexToFelt(t, "0x7fb440b1dee35c5259bf10d55782bc973434d195bb5c8b95ac7d3d8e2a8a0e4"),
						utils.TestHexToFelt(t, "0x3e16f111f8a22cb484b09d2554fca1e669cd540c5ad7cf2b9878071a9b95693"),
					},
					SenderAddress: utils.TestHexToFelt(t, "0x35acd6dd6c5045d18ca6d0192af46b335a5402c02d41f46e4e77ea2c951d9a3"),
					Calldata: []*felt.Felt{
						utils.TestHexToFelt(t, "0x1"),
						utils.TestHexToFelt(t, "0x3fe8e4571772bbe0065e271686bd655efd1365a5d6858981e582f82f2c10313"),
						utils.TestHexToFelt(t, "0x2468d193cd15b621b24c2a602b8dbcfa5eaa14f88416c40c09d7fd12592cb4b"),
						utils.TestHexToFelt(t, "0x0"),
					},
				},
			},
		},
	}

	testSet := map[string][]testSetType{
		"mock": {},
		"testnet": {
			{
				BlockID:       WithBlockTag("latest"),
				ExpectedError: nil,
			},
			{
				BlockID:                     WithBlockHash(utils.TestHexToFelt(t, "0x2488a7149327b4dcd200f05a131911bd44f946021539648642eaa7d6e82f289")),
				LookupTxnPositionInExpected: 0,
				LookupTxnPositionInOriginal: 0,
				ExpectedError:               nil,
				want:                        &fullBlockSepolia30436,
			},
			{
				BlockID:                     WithBlockNumber(30436),
				LookupTxnPositionInExpected: 0,
				LookupTxnPositionInOriginal: 0,
				ExpectedError:               nil,
				want:                        &fullBlockSepolia30436,
			},
		},
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		spy := NewSpy(testConfig.provider.c)
		testConfig.provider.c = spy
		blockWithTxsInterface, err := testConfig.provider.BlockWithTxs(context.Background(), test.BlockID)
		if err != test.ExpectedError {
			t.Fatal("BlockWithTxHashes match the expected error:", err)
		}
		if test.ExpectedError != nil && blockWithTxsInterface == nil {
			continue
		}
		blockWithTxs, ok := blockWithTxsInterface.(*Block)
		if !ok {
			t.Fatalf("expecting *rpv02.Block, instead %T", blockWithTxsInterface)
		}
		_, err = spy.Compare(blockWithTxs, false)
		if err != nil {
			t.Fatal("expecting to match", err)
		}
		if !strings.HasPrefix(blockWithTxs.BlockHash.String(), "0x") {
			t.Fatal("Block Hash should start with \"0x\", instead", blockWithTxs.BlockHash)
		}

		if len(blockWithTxs.Transactions) == 0 {
			t.Fatal("the number of transaction should not be 0")
		}

		if test.want != nil {
			if (*test.want).BlockHash == &felt.Zero {
				continue
			}

			invokeV1Want, ok := (*test.want).Transactions[test.LookupTxnPositionInExpected].(BlockInvokeTxnV1)
			if !ok {
				t.Fatal("expected invoke v1 transaction")
			}
			invokeV1Block, ok := blockWithTxs.Transactions[test.LookupTxnPositionInOriginal].(BlockInvokeTxnV1)
			if !ok {
				t.Fatal("expected invoke v1 transaction")
			}
			require.Equal(t, invokeV1Want.TransactionHash, invokeV1Block.TransactionHash, "expected equal TransactionHash")
			require.Equal(t, invokeV1Want.InvokeTxnV1.MaxFee, invokeV1Block.InvokeTxnV1.MaxFee, "expected equal maxfee")
			require.Equal(t, invokeV1Want.InvokeTxnV1.SenderAddress, invokeV1Block.InvokeTxnV1.SenderAddress, "expected equal senders addresses")
		}
	}
}

// TestBlockWithTxsAndInvokeTXNV3 tests the BlockWithTxsAndInvokeTXNV3 function.
//
// The function tests the BlockWithTxsAndInvokeTXNV3 function by setting up a test configuration and a test set type.
// It then initializes a fullBlockSepolia52767 variable with a Block struct and invokes the BlockWithTxs function with different test scenarios.
// The function compares the expected error with the actual error and checks if the BlockWithTxs function returns the correct block data.
// It also verifies the block hash, the number of transactions in the block, and the details of a specific transaction.
//
// Parameters:
// - t: The t testing object
// Returns:
//
//	none
func TestBlockWithTxsAndInvokeTXNV3(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		BlockID                     BlockID
		ExpectedError               error
		LookupTxnPositionInOriginal int
		LookupTxnPositionInExpected int
		want                        *Block
	}

	var fullBlockSepolia52767 = Block{
		BlockHeader: BlockHeader{
			BlockHash:        utils.TestHexToFelt(t, "0x4ae5d52c75e4dea5694f456069f830cfbc7bec70427eee170c3385f751b8564"),
			ParentHash:       utils.TestHexToFelt(t, "0x7d3a1bc98e49c197b38538fbc351dae6ed4f0ff4e718db119ddacef3088b928"),
			SequencerAddress: utils.TestHexToFelt(t, "0x1176a1bd84444c89232ec27754698e5d2e7e1a7f1539f12027f28b23ec9f3d8"),
			BlockNumber:      52767,
			NewRoot:          utils.TestHexToFelt(t, "0x3d86be8765b9b6ab724fb8c10a64c4e1705bcc6d39032fe9973037abedc113a"),
			Timestamp:        1661450764,
		},
		Status: "ACCEPTED_ON_L1",
		Transactions: []BlockTransaction{
			BlockInvokeTxnV3{
				TransactionHash: utils.TestHexToFelt(t, "0xb91eada292de46f4ec663bac57699c7b8f8fa454a8efad91fde7f35d941199"),
				InvokeTxnV3: InvokeTxnV3{
					Type:          "INVOKE",
					SenderAddress: utils.TestHexToFelt(t, "0x573ea9a8602e03417a4a31d55d115748f37a08bbb23adf6347cb699743a998d"),
					Nonce:         utils.TestHexToFelt(t, "0x470de4df820000"),
					Version:       TransactionV3,
					Signature: []*felt.Felt{
						utils.TestHexToFelt(t, "0x7f14bb280b602d0c0e22e91ea5d987371554664f68c57c0acd16bf9f8be36b4"),
						utils.TestHexToFelt(t, "0x22c57ce8eb211c7fe0f04e7da338f579fcbc9e8997ec432fac7738c80fd56ad"),
					},
					Calldata: []*felt.Felt{
						utils.TestHexToFelt(t, "0x1"),
						utils.TestHexToFelt(t, "0x3fe8e4571772bbe0065e271686bd655efd1365a5d6858981e582f82f2c10313"),
						utils.TestHexToFelt(t, "0x2468d193cd15b621b24c2a602b8dbcfa5eaa14f88416c40c09d7fd12592cb4b"),
						utils.TestHexToFelt(t, "0x0"),
					},
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
					Tip:                   "0x0",
					PayMasterData:         []*felt.Felt{},
					AccountDeploymentData: []*felt.Felt{},
					NonceDataMode:         DAModeL1,
					FeeMode:               DAModeL1,
				},
			},
		},
	}

	testSet := map[string][]testSetType{
		"mock": {},
		"testnet": {
			{
				BlockID:       WithBlockTag("latest"),
				ExpectedError: nil,
			},
			{
				BlockID:                     WithBlockHash(utils.TestHexToFelt(t, "0x4ae5d52c75e4dea5694f456069f830cfbc7bec70427eee170c3385f751b8564")),
				ExpectedError:               nil,
				LookupTxnPositionInExpected: 0,
				LookupTxnPositionInOriginal: 25,
				want:                        &fullBlockSepolia52767,
			},
			{
				BlockID:                     WithBlockNumber(52767),
				ExpectedError:               nil,
				LookupTxnPositionInExpected: 0,
				LookupTxnPositionInOriginal: 25,
				want:                        &fullBlockSepolia52767,
			},
		},
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		spy := NewSpy(testConfig.provider.c)
		testConfig.provider.c = spy
		blockWithTxsInterface, err := testConfig.provider.BlockWithTxs(context.Background(), test.BlockID)
		if err != test.ExpectedError {
			t.Fatal("BlockWithTxHashes match the expected error:", err)
		}
		if test.ExpectedError != nil && blockWithTxsInterface == nil {
			continue
		}
		blockWithTxs, ok := blockWithTxsInterface.(*Block)
		if !ok {
			t.Fatalf("expecting *rpv02.Block, instead %T", blockWithTxsInterface)
		}
		_, err = spy.Compare(blockWithTxs, false)
		if err != nil {
			t.Fatal("expecting to match", err)
		}
		if !strings.HasPrefix(blockWithTxs.BlockHash.String(), "0x") {
			t.Fatal("Block Hash should start with \"0x\", instead", blockWithTxs.BlockHash)
		}

		if len(blockWithTxs.Transactions) == 0 {
			t.Fatal("the number of transaction should not be 0")
		}

		if test.want != nil {
			if (*test.want).BlockHash == &felt.Zero {
				continue
			}

			invokeV3Want, ok := (*test.want).Transactions[test.LookupTxnPositionInExpected].(BlockInvokeTxnV3)
			if !ok {
				t.Fatal("expected invoke v3 transaction")
			}

			invokeV3Block, ok := blockWithTxs.Transactions[test.LookupTxnPositionInOriginal].(BlockInvokeTxnV3)
			if !ok {
				t.Fatal("expected invoke v3 transaction")
			}

			require.Equal(t, invokeV3Want.TransactionHash, invokeV3Block.TransactionHash, "expected equal TransactionHash")
			require.Equal(t, invokeV3Want.InvokeTxnV3.NonceDataMode, invokeV3Block.InvokeTxnV3.NonceDataMode, "expected equal nonceDataMode")
			require.Equal(t, invokeV3Want.InvokeTxnV3.FeeMode, invokeV3Block.InvokeTxnV3.FeeMode, "expected equal feeMode")
		}
	}
}

// TestBlockWithTxsAndDeployOrDeclare tests BlockWithTxs with Deploy or Declare TXN

// TestBlockWithTxsAndDeployOrDeclare is a test function that tests the functionality of the BlockWithTxs function.
// It creates a test configuration, defines a testSetType struct, and initializes three Blocks (fullBlockSepolia52959, fullBlockSepolia848622 and fullBlockSepolia849399).
// It then defines a testSet map with different test scenarios for the BlockWithTxs function.
// The function iterates over the testSet and performs the BlockWithTxs operation on each test case.
// It compares the returned blockWithTxs with the expected result and verifies the correctness of the operation.
// The function also checks the block hash, the number of transactions, and other properties of the returned blockWithTxs.
// The function returns an error if the actual result does not match the expected result.
// It uses the Spy object to compare the blockWithTxs with the expected result and returns an error if they don't match.
// The function also checks the block hash to ensure it starts with "0x" and verifies that the number of transactions is not zero.
// Finally, the function compares the transactions of the returned blockWithTxs with the expected transactions and returns an error if they don't match.
//
// Parameters:
// - t: *testing.T - the testing object for running the test cases
// Returns:
//
//	none
func TestBlockWithTxsAndDeployAccountOrDeclare(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		BlockID                     BlockID
		ExpectedError               error
		LookupTxnPositionInOriginal int
		LookupTxnPositionInExpected int
		ExpectedBlockWithTxs        *Block
	}

	var fullBlockSepolia52959 = Block{
		BlockHeader: BlockHeader{
			BlockHash:        utils.TestHexToFelt(t, "0x5e4c92970bb2bc51a3824a8357078ef00e0c089313c4ac1d9004166d9adc6aa"),
			ParentHash:       utils.TestHexToFelt(t, "0x469f2b163bd62e042e77390ae3d1fa278212e279408163e14624d39d6529bd5"),
			SequencerAddress: utils.TestHexToFelt(t, "0x1176a1bd84444c89232ec27754698e5d2e7e1a7f1539f12027f28b23ec9f3d8"),
			BlockNumber:      52959,
			NewRoot:          utils.TestHexToFelt(t, "0x5aebb84ffae11caa645da6ade011df2e5d60a1943d9533fcd8326ff5ca000b2"),
			Timestamp:        1711378335,
		},
		Status: "ACCEPTED_ON_L1",
		Transactions: []BlockTransaction{
			BlockDeployAccountTxn{
				TransactionHash: utils.TestHexToFelt(t, "0x38d17f7aaa320f43a638a0097a43332614d3306d91036cab258a07441a14a10"),
				DeployAccountTxn: DeployAccountTxn{
					ClassHash: utils.TestHexToFelt(t, "0x29927c8af6bccf3f6fda035981e765a7bdbf18a2dc0d630494f8758aa908e2b"),
					Signature: []*felt.Felt{
						utils.TestHexToFelt(t, "0x20b878ad13672908aadeb500348d3c172a9d2534474d5f87788de5268b29c4e"),
						utils.TestHexToFelt(t, "0xd4b95c47f5234c220e2c9bc96b15677aad0cc47d74ed647c53ab121f632e95"),
					},
					ConstructorCalldata: []*felt.Felt{
						utils.TestHexToFelt(t, "0x7e02fa1096a3292b92f28167a180dc6b0944f1d0cc1b4919c88cc2ca961848"),
						utils.TestHexToFelt(t, "0x0"),
					},
					ContractAddressSalt: utils.TestHexToFelt(t, "0x7e02fa1096a3292b92f28167a180dc6b0944f1d0cc1b4919c88cc2ca961848"),
					MaxFee:              utils.TestHexToFelt(t, "0x12aff9cd5fac"),
					Type:                "DEPLOY_ACCOUNT",
					Nonce:               utils.TestHexToFelt(t, "0x0"),
					Version:             TransactionV1,
				},
			},
		},
	}

	var fullBlockSepolia53617 = Block{
		BlockHeader: BlockHeader{
			BlockHash:        utils.TestHexToFelt(t, "0x526355d90ef6118d4f871ccdb3a3d0ea27d10b0a02b3005e7697dd321b52ada"),
			ParentHash:       utils.TestHexToFelt(t, "0x645b998412c408e979577b357f4afc0902652fed85bdd7f9d81b7fa2ffc9506"),
			SequencerAddress: utils.TestHexToFelt(t, "0x1176a1bd84444c89232ec27754698e5d2e7e1a7f1539f12027f28b23ec9f3d8"),
			BlockNumber:      53617,
			NewRoot:          utils.TestHexToFelt(t, "0x703f5a9be28fc4e57f3dd84abd05acb5bbdc957b77ce9e134b0d5d08d04d9ec"),
			Timestamp:        1711538993,
		},
		Status: "ACCEPTED_ON_L1",
		Transactions: []BlockTransaction{
			BlockDeclareTxnV1{
				TransactionHash: utils.TestHexToFelt(t, "0x7ff4942928fb6383514fd18d8b4ddc8a154b496d2204f60ee47e6272305172c"),
				DeclareTxnV1: DeclareTxnV1{
					Type:    TransactionType_Declare,
					MaxFee:  utils.TestHexToFelt(t, "0x3365d77ce5000"),
					Version: TransactionV1,
					Signature: []*felt.Felt{
						utils.TestHexToFelt(t, "0x106b85d29dfe14af94d46e8758491c468f0812227c8969de7f7c382355bd72c"),
						utils.TestHexToFelt(t, "0x28a8bdedf588daa657057240e952fb93070290b6d906bc80c9e9e3f3e3be6da"),
					},
					Nonce:         utils.TestHexToFelt(t, "0x128"),
					ClassHash:     utils.TestHexToFelt(t, "0x6b433af0dff031b3578de26217fa7cd2a8ac0d70c25f8fbf332fc603a5dcf2d"),
					SenderAddress: utils.TestHexToFelt(t, "0x2cc631ca0c544639f6e4403b8f3611696a3d831e8157ea1c946e35429c7ac31"),
				},
			},
		},
	}

	var fullBlockSepolia53794 = Block{
		BlockHeader: BlockHeader{
			BlockHash:        utils.TestHexToFelt(t, "0x5d216dd9e394a088330b5d77a4a282618b7f0796a2ef8f081c5aa01e3ace6f0"),
			ParentHash:       utils.TestHexToFelt(t, "0x6e31bf89281805f949a0ab5d7bf6fa44becf7c4cff4d426626a1647376e461f"),
			SequencerAddress: utils.TestHexToFelt(t, "0x1176a1bd84444c89232ec27754698e5d2e7e1a7f1539f12027f28b23ec9f3d8"),
			BlockNumber:      53794,
			NewRoot:          utils.TestHexToFelt(t, "0x1c6f8a85528d03dcf88a4087bb191bd608f301a7eb3e05d236b0c441fc30174"),
			Timestamp:        1711582037,
		},
		Status: "ACCEPTED_ON_L1",
		Transactions: []BlockTransaction{
			BlockDeclareTxnV2{
				TransactionHash: utils.TestHexToFelt(t, "0x7c0108477e9ece3dbd421f74e73bd4e15c7fcb496784e99d4b6f6710463b6f3"),
				DeclareTxnV2: DeclareTxnV2{
					Type:    TransactionType_Declare,
					MaxFee:  utils.TestHexToFelt(t, "0xde0b6b3a7640000"),
					Version: TransactionV2,
					Signature: []*felt.Felt{
						utils.TestHexToFelt(t, "0x2dd7cf8724045aa56edce10f96d1ec62a89492dc69e43d4b21c62a3708f3eb6"),
						utils.TestHexToFelt(t, "0x4918acc659aa5e4e4b4b8a7ffdc3686fefd219c556ef7582904675cae6b1028"),
					},
					Nonce:             utils.TestHexToFelt(t, "0x353"),
					ClassHash:         utils.TestHexToFelt(t, "0x36f7463a30a3ff34ecdcfeea4aa880b97c366523925b1ba069faadddbef0e02"),
					CompiledClassHash: utils.TestHexToFelt(t, "0x5047109bf7eb550c5e6b0c37714f6e0db4bb8b5b227869e0797ecfc39240aa7"),
					SenderAddress:     utils.TestHexToFelt(t, "0xaf46a3d75c1abc204cbe7e08f220680958bd8aca2c3cfc2ef34c686148ecf7"),
				},
			},
		},
	}

	testSet := map[string][]testSetType{
		"mock": {},
		"testnet": {
			{
				BlockID:       WithBlockTag("latest"),
				ExpectedError: nil,
			},
			{
				BlockID:                     WithBlockHash(utils.TestHexToFelt(t, "0x5e4c92970bb2bc51a3824a8357078ef00e0c089313c4ac1d9004166d9adc6aa")),
				ExpectedError:               nil,
				LookupTxnPositionInOriginal: 4,
				LookupTxnPositionInExpected: 0,
				ExpectedBlockWithTxs:        &fullBlockSepolia52959,
			},
			{
				BlockID:                     WithBlockNumber(52959),
				ExpectedError:               nil,
				LookupTxnPositionInOriginal: 4,
				LookupTxnPositionInExpected: 0,
				ExpectedBlockWithTxs:        &fullBlockSepolia52959,
			},
			{
				BlockID:                     WithBlockHash(utils.TestHexToFelt(t, "0x5d216dd9e394a088330b5d77a4a282618b7f0796a2ef8f081c5aa01e3ace6f0")),
				ExpectedError:               nil,
				LookupTxnPositionInOriginal: 72,
				LookupTxnPositionInExpected: 0,
				ExpectedBlockWithTxs:        &fullBlockSepolia53794,
			},
			{
				BlockID:                     WithBlockNumber(53794),
				ExpectedError:               nil,
				LookupTxnPositionInOriginal: 72,
				LookupTxnPositionInExpected: 0,
				ExpectedBlockWithTxs:        &fullBlockSepolia53794,
			},
			{
				BlockID:                     WithBlockHash(utils.TestHexToFelt(t, "0x526355d90ef6118d4f871ccdb3a3d0ea27d10b0a02b3005e7697dd321b52ada")),
				ExpectedError:               nil,
				LookupTxnPositionInOriginal: 72,
				LookupTxnPositionInExpected: 0,
				ExpectedBlockWithTxs:        &fullBlockSepolia53617,
			},
			{
				BlockID:                     WithBlockNumber(53617),
				ExpectedError:               nil,
				LookupTxnPositionInOriginal: 72,
				LookupTxnPositionInExpected: 0,
				ExpectedBlockWithTxs:        &fullBlockSepolia53617,
			},
		},
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		spy := NewSpy(testConfig.provider.c)
		testConfig.provider.c = spy
		blockWithTxsInterface, err := testConfig.provider.BlockWithTxs(context.Background(), test.BlockID)
		if err != test.ExpectedError {
			t.Fatal("BlockWithTxHashes match the expected error:", err)
		}
		if test.ExpectedError != nil && blockWithTxsInterface == nil {
			continue
		}
		blockWithTxs, ok := blockWithTxsInterface.(*Block)
		if !ok {
			t.Fatalf("expecting *rpc.Block, instead %T", blockWithTxsInterface)
		}
		diff, err := spy.Compare(blockWithTxs, false)
		if err != nil {
			t.Fatal("expecting to match", err)
		}
		if diff != "FullMatch" {
			if _, err := spy.Compare(blockWithTxs, false); err != nil {
				t.Fatal(err)
			}
		}
		if !strings.HasPrefix(blockWithTxs.BlockHash.String(), "0x") {
			t.Fatal("Block Hash should start with \"0x\", instead", blockWithTxs.BlockHash)
		}

		if len(blockWithTxs.Transactions) == 0 {
			t.Fatal("the number of transaction should not be 0")
		}
		if test.ExpectedBlockWithTxs != nil {
			if test.ExpectedBlockWithTxs.BlockHash == &felt.Zero {
				continue
			}
			if !cmp.Equal(test.ExpectedBlockWithTxs.Transactions[test.LookupTxnPositionInExpected], blockWithTxs.Transactions[test.LookupTxnPositionInOriginal]) {
				t.Fatalf("the expected transaction blocks to match, instead: %s", cmp.Diff(test.ExpectedBlockWithTxs.Transactions[test.LookupTxnPositionInExpected], blockWithTxs.Transactions[test.LookupTxnPositionInOriginal]))
			}
		}

	}
}

// TestBlockTransactionCount tests the function that calculates the number of transactions in a block.
//
// This function tests the BlockTransactionCount function by running it with different test cases.
// It verifies that the function returns the expected count of transactions for each test case.
// The test cases include a mock environment, a testnet environment, and a mainnet environment.
// For each test case, the function sets up a spy provider and calls BlockTransactionCount with a specific block ID.
// It then compares the returned count with the expected count and verifies that they match.
// If the counts do not match, the function reports an error and provides additional information.
// Finally, the function terminates if all test cases pass.
//
// Parameters:
// - t: the testing object for running the test cases
// Returns:
//
//	none
func TestBlockTransactionCount(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		BlockID       BlockID
		ExpectedCount uint64
	}
	testSet := map[string][]testSetType{
		"mock": {
			{
				BlockID:       WithBlockNumber(300000),
				ExpectedCount: 10,
			},
		},
		"testnet": {
			{
				BlockID:       WithBlockNumber(52959),
				ExpectedCount: 58,
			},
		},
		"mainnet": {},
	}[testEnv]
	for _, test := range testSet {
		spy := NewSpy(testConfig.provider.c)
		testConfig.provider.c = spy
		count, err := testConfig.provider.BlockTransactionCount(context.Background(), test.BlockID)
		if err != nil {
			t.Fatal(err)
		}
		diff, err := spy.Compare(count, false)
		if err != nil {
			t.Fatal("expecting to match", err)
		}
		if diff != "FullMatch" {
			if _, err := spy.Compare(count, true); err != nil {
				t.Fatal(err)
			}
			t.Fatal("structure expecting to be FullMatch, instead", diff)
		}
		if count != test.ExpectedCount {
			t.Fatalf("structure expecting %d, instead: %d", test.ExpectedCount, count)
		}
	}
}

// TestCaptureUnsupportedBlockTxn tests the functionality of capturing unsupported block transactions.
//
// Parameters:
// - t: the testing object for running the test cases
// Returns:
//
//	none
func TestCaptureUnsupportedBlockTxn(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		StartBlock uint64
		EndBlock   uint64
	}
	testSet := map[string][]testSetType{
		"mock": {},
		"testnet": {
			{
				StartBlock: 381000,
				EndBlock:   381001,
			},
		},
		"mainnet": {},
	}[testEnv]
	for _, test := range testSet {
		for i := test.StartBlock; i < test.EndBlock; i++ {
			blockWithTxsInterface, err := testConfig.provider.BlockWithTxs(context.Background(), WithBlockNumber(i))
			if err != nil {
				t.Fatal("BlockWithTxHashes match the expected error:", err)
			}
			blockWithTxs, ok := blockWithTxsInterface.(*Block)
			if !ok {
				t.Fatalf("expecting *rpc.Block, instead %T", blockWithTxsInterface)
			}
			for k, v := range blockWithTxs.Transactions {
				_, okv1 := v.(BlockInvokeTxnV1)
				_, okv0 := v.(BlockInvokeTxnV0)
				_, okl1 := v.(BlockL1HandlerTxn)
				_, okdec0 := v.(BlockDeclareTxnV0)
				_, okdec1 := v.(BlockDeclareTxnV1)
				_, okdec2 := v.(BlockDeclareTxnV2)
				_, okdep := v.(BlockDeployTxn)
				_, okdepac := v.(BlockDeployAccountTxn)
				if !okv0 && !okv1 && !okl1 && !okdec0 && !okdec1 && !okdec2 && !okdep && !okdepac {
					t.Fatalf("New Type Detected %T at Block(%d)/Txn(%d)", v, i, k)
				}
			}
		}
	}
}

// TestStateUpdate is a test function for the StateUpdate method.
//
// It tests the StateUpdate method by creating a test set and iterating through each test case.
// For each test case, it creates a spy object and sets it as the provider.
// Then, it calls the StateUpdate method with the given test block ID.
// If there is an error, it fails the test.
// If the returned block hash does not match the expected block hash, it fails the test.
// TODO: this is not implemented yet with pathfinder as you can see from the
// [code](https://github.com/eqlabs/pathfinder/blob/927183552dad6dcdfebac16c8c1d2baf019127b1/crates/pathfinder/rpc_examples.sh#L37)
// check when it is and test when it is the case.
//
// Parameters:
// - t: the testing object for running the test cases
// Returns:
//
//	none
func TestStateUpdate(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		BlockID                   BlockID
		ExpectedStateUpdateOutput StateUpdateOutput
	}

	testSet := map[string][]testSetType{
		"mock": {
			{
				BlockID: WithBlockNumber(300000),
				ExpectedStateUpdateOutput: StateUpdateOutput{
					BlockHash: utils.TestHexToFelt(t, "0x4f1cee281edb6cb31b9ba5a8530694b5527cf05c5ac6502decf3acb1d0cec4"),
					NewRoot:   utils.TestHexToFelt(t, "0x70677cda9269d47da3ff63bc87cf1c87d0ce167b05da295dc7fc68242b250b"),
					PendingStateUpdate: PendingStateUpdate{
						OldRoot: utils.TestHexToFelt(t, "0x19aa982a75263d4c4de4cc4c5d75c3dec32e00b95bef7bbb4d17762a0b138af"),
						StateDiff: StateDiff{
							StorageDiffs: []ContractStorageDiffItem{{
								Address: utils.TestHexToFelt(t, "0xe5cc6f2b6d34979184b88334eb64173fe4300cab46ecd3229633fcc45c83d4"),
								StorageEntries: []StorageEntry{
									{
										Key:   utils.TestHexToFelt(t, "0x1813aac5f5e7799684c6dc33e51f44d3627fd748c800724a184ed5be09b713e"),
										Value: utils.TestHexToFelt(t, "0x630b4197"),
									},
								},
							}},
						},
					},
				},
			},
		},
		"testnet": {
			{
				BlockID: WithBlockNumber(300000),
				ExpectedStateUpdateOutput: StateUpdateOutput{
					BlockHash: utils.TestHexToFelt(t, "0x03b6d94b246815960f38b7dffc53cda192e7d1dcfff61e1bc042fb57e95f8349"),
					NewRoot:   utils.TestHexToFelt(t, "0x70677cda9269d47da3ff63bc87cf1c87d0ce167b05da295dc7fc68242b250b"),
					PendingStateUpdate: PendingStateUpdate{
						OldRoot: utils.TestHexToFelt(t, "0x19aa982a75263d4c4de4cc4c5d75c3dec32e00b95bef7bbb4d17762a0b138af"),
						StateDiff: StateDiff{
							StorageDiffs: []ContractStorageDiffItem{{
								Address: utils.TestHexToFelt(t, "0xe5cc6f2b6d34979184b88334eb64173fe4300cab46ecd3229633fcc45c83d4"),
								StorageEntries: []StorageEntry{
									{
										Key:   utils.TestHexToFelt(t, "0x1813aac5f5e7799684c6dc33e51f44d3627fd748c800724a184ed5be09b713e"),
										Value: utils.TestHexToFelt(t, "0x630b4197"),
									},
								},
							}},
						},
					},
				},
			},
		},
		"mainnet": {},
	}[testEnv]
	for _, test := range testSet {
		spy := NewSpy(testConfig.provider.c)
		testConfig.provider.c = spy
		stateUpdate, err := testConfig.provider.StateUpdate(context.Background(), test.BlockID)
		if err != nil {
			t.Fatal(err)
		}
		if stateUpdate.BlockHash.String() != test.ExpectedStateUpdateOutput.BlockHash.String() {
			t.Fatalf("structure expecting %s, instead: %s", test.ExpectedStateUpdateOutput.BlockHash.String(), stateUpdate.BlockHash.String())
		}
	}
}
