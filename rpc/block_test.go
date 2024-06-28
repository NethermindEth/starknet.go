package rpc

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/utils"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		require.NoError(t, err, "BlockNumber should not return an error")

		diff, err := spy.Compare(blockNumber, false)
		require.NoError(t, err, "expecting to match")
		require.Equal(t, "FullMatch", diff, "expecting to match, instead %s", diff)

		require.False(t, blockNumber <= 3000, fmt.Sprintf("Block number should be > 3000, instead: %d", blockNumber))
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
		require.NoError(t, err, "BlockHashAndNumber should not return an error")

		diff, err := spy.Compare(blockHashAndNumber, false)
		require.NoError(t, err, "expecting to match")
		require.Equal(t, "FullMatch", diff, "expecting to match, instead %s", diff)

		require.False(t, blockHashAndNumber.BlockNumber <= 3000, "Block number should be > 3000, instead: %d", blockHashAndNumber.BlockNumber)

		require.True(t, strings.HasPrefix(blockHashAndNumber.BlockHash.String(), "0x"), "current block hash should return a string starting with 0x")
	}
}

// TestBlockWithTxHashes tests the functionality of the BlockWithTxHashes function.
//
// The function takes a testing.T object as a parameter and initializes a testConfig object.
// It defines a testSetType struct that contains several fields including BlockID, ExpectedError, ExpectedBlockWithTxHashes, and ExpectedPendingBlockWithTxHashes.
// The function then initializes a blockSepolia64159 variable of type BlockTxHashes with a predefined set of values.
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

	var blockSepolia64159 = BlockTxHashes{
		BlockHeader: BlockHeader{
			BlockHash:        utils.TestHexToFelt(t, "0x6df565874b2ea6a02d346a23f9efb0b26abbf5708b51bb12587f88a49052964"),
			ParentHash:       utils.TestHexToFelt(t, "0x1406ec9385293905d6c20e9c5aa0bbf9f63f87d39cf12fcdfef3ed0d056c0f5"),
			SequencerAddress: utils.TestHexToFelt(t, "0x1176a1bd84444c89232ec27754698e5d2e7e1a7f1539f12027f28b23ec9f3d8"),
			BlockNumber:      64159,
			NewRoot:          utils.TestHexToFelt(t, "0x310be818a18de0d6f6c1391f467d0dbd1a2753e6dde876449448465f8e617f0"),
			Timestamp:        1714901729,
		},
		Status: "ACCEPTED_ON_L1",
		Transactions: utils.TestHexArrToFelt(t, []string{
			"0x5e3fcf2f7dc0f3786a7406dc271cc54a00ba4658f9d0567b25b8e2a90a6250f",
			"0x603fec12cfbc73bcf411fdf6f7780d4698c71b989002192a1277025235b23b9",
			"0x23a24d95872d0eb15bf54cfb432830a3b85ad5c621b5edf849f131a2a45988d",
			"0x4a35d133717f1f0288346432037d7964a16e503100fdab4cc3914a44790f3b4",
			"0x5f14364b746abcfdfc0280877ff6d18c311d363e62264d7f218c5da2d396acc",
			"0x7b49053e9a0bcd28c40d946702c28cf4ba068ff2e1755eff3fd99d62aada1a8",
			"0x173c7a20046ab576667a2581cdea9565160a76e028864102a4d0828ca35a0d3",
			"0x5754961d70d6f39d0e2c71a1a4ff5df0a26b1ceda4881ca82898994379e1e73",
			"0x402f81ce59e3e79ca12c9cffea888c3f02542f2f4926731cb2145a3b8e810d5",
			"0x48fa3e27b725e595e910548f3c0bb1ddfb30d32b31ebbf23fa1e63a66b0e59d",
			"0x3d43ca0ea28f8e412b6abb37b76e75ac33e7df177cc8e4221e361ed0621bcdd",
			"0x692381bba0e8505a8e0b92d0f046c8272de9e65f050850df678a0c10d8781d",
			"0x7a56162702394b43b8fc05c511e1ddbe490749f2fd6786365d5d59797cd2012",
			"0xcaead659d470ac0b572f0a4d8831275d51944d08961aab1a85bd68f9e98409",
			"0x73376f10049aa689a5c6bf78b39b5a8c76ce5fb6611290b3080aa0d4f492d56",
			"0x45061dccdb8cb32e428ec7b25136ae3a691f02bf5d01bd2d30ae9d2d4a29d4e",
			"0x707071e6d5354a254935bf605b4eba4bb289261dc9ce75c1e9d8ad1367a5154",
			"0x5358c68fa172aafabae1a007f5ec71eb1bddd64d4574366880f253287e8a0be",
			"0x27a2e57d2ead1f9b7ce99c2e7b367820ecedf66cae3353572e6e7a7de89d1ce",
			"0x59866d2230404c1c5787b523d3233f7f9d4ea9faf804ff747a0a2d3da77720e",
			"0x5d41f4dec3678156d3888d6b890648c3baa02d866820689d5f8b3e20735521b",
		}),
	}

	txHashes := utils.TestHexArrToFelt(t, []string{
		"0x5754961d70d6f39d0e2c71a1a4ff5df0a26b1ceda4881ca82898994379e1e73",
		"0x692381bba0e8505a8e0b92d0f046c8272de9e65f050850df678a0c10d8781d",
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
				BlockID:                   WithBlockHash(utils.TestHexToFelt(t, "0x6df565874b2ea6a02d346a23f9efb0b26abbf5708b51bb12587f88a49052964")),
				ExpectedErr:               nil,
				ExpectedBlockWithTxHashes: &blockSepolia64159,
			},
			{
				BlockID:                   WithBlockNumber(64159),
				ExpectedErr:               nil,
				ExpectedBlockWithTxHashes: &blockSepolia64159,
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
			require.True(t, ok, fmt.Sprintf("should return *BlockTxHashes, instead: %T\n", result))

			if test.ExpectedErr != nil {
				continue
			}

			require.True(t, strings.HasPrefix(block.BlockHash.String(), "0x"), "Block Hash should start with \"0x\", instead: %s", block.BlockHash)
			require.NotEmpty(t, block.Transactions, "the number of transactions should not be 0")

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
			require.True(t, ok, fmt.Sprintf("should return *PendingBlockTxHashes, instead: %T\n", result))

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
func TestBlockWithTxsAndInvokeTXNV0(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		BlockID                     BlockID
		ExpectedError               error
		LookupTxnPositionInOriginal int
		LookupTxnPositionInExpected int
		want                        *Block
	}

	testSet := map[string][]testSetType{
		"mock":    {},
		"testnet": {},
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		spy := NewSpy(testConfig.provider.c)
		testConfig.provider.c = spy
		blockWithTxsInterface, err := testConfig.provider.BlockWithTxs(context.Background(), test.BlockID)
		require.Equal(t, test.ExpectedError, err, "Error in BlockWithTxHashes doesn't match the expected error")

		if test.ExpectedError != nil && blockWithTxsInterface == nil {
			continue
		}
		blockWithTxs, ok := blockWithTxsInterface.(*Block)
		require.True(t, ok, fmt.Sprintf("expecting *rpv02.Block, instead %T", blockWithTxsInterface))

		_, err = spy.Compare(blockWithTxs, false)
		require.NoError(t, err, "expecting to match")

		require.True(t, strings.HasPrefix(blockWithTxs.BlockHash.String(), "0x"), "Block Hash should start with \"0x\", instead: %s", blockWithTxs.BlockHash)
		require.NotEmpty(t, blockWithTxs.Transactions, "the number of transactions should not be 0")

		if test.want != nil {
			if (*test.want).BlockHash == &felt.Zero {
				continue
			}

			invokeV0Want, ok := (*test.want).Transactions[test.LookupTxnPositionInExpected].(BlockInvokeTxnV0)
			require.True(t, ok, "expected invoke v0 transaction")

			invokeV0Block, ok := blockWithTxs.Transactions[test.LookupTxnPositionInOriginal].(BlockInvokeTxnV0)
			require.True(t, ok, "expected invoke v0 transaction")

			require.Equal(t, invokeV0Want.TransactionHash, invokeV0Block.TransactionHash, "expected equal TransactionHash")
			require.Equal(t, invokeV0Want.InvokeTxnV0.MaxFee, invokeV0Block.InvokeTxnV0.MaxFee, "expected equal maxfee")
			require.Equal(t, invokeV0Want.InvokeTxnV0.EntryPointSelector, invokeV0Block.InvokeTxnV0.EntryPointSelector, "expected equal eps")

		}

	}
}

// TestBlockWithTxsAndInvokeTXNV1 tests the BlockWithTxsAndInvokeTXNV1 function.
//
// The function tests the BlockWithTxsAndInvokeTXNV1 function by setting up a test configuration and a test set type.
// It then initializes a fullBlockSepolia64159 variable with a Block struct and invokes the BlockWithTxs function with different test scenarios.
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

	var fullBlockSepolia64159 = Block{
		BlockHeader: BlockHeader{
			BlockHash:        utils.TestHexToFelt(t, "0x6df565874b2ea6a02d346a23f9efb0b26abbf5708b51bb12587f88a49052964"),
			ParentHash:       utils.TestHexToFelt(t, "0x1406ec9385293905d6c20e9c5aa0bbf9f63f87d39cf12fcdfef3ed0d056c0f5"),
			SequencerAddress: utils.TestHexToFelt(t, "0x1176a1bd84444c89232ec27754698e5d2e7e1a7f1539f12027f28b23ec9f3d8"),
			BlockNumber:      64159,
			NewRoot:          utils.TestHexToFelt(t, "0x310be818a18de0d6f6c1391f467d0dbd1a2753e6dde876449448465f8e617f0"),
			Timestamp:        1714901729,
		},
		Status: "ACCEPTED_ON_L1",
		Transactions: []BlockTransaction{

			BlockInvokeTxnV1{
				TransactionHash: utils.TestHexToFelt(t, "0x5f14364b746abcfdfc0280877ff6d18c311d363e62264d7f218c5da2d396acc"),
				InvokeTxnV1: InvokeTxnV1{
					Type:          "INVOKE",
					Version:       TransactionV1,
					Nonce:         utils.TestHexToFelt(t, "0x33"),
					MaxFee:        utils.TestHexToFelt(t, "0x1bad55a98e1c1"),
					SenderAddress: utils.TestHexToFelt(t, "0x3543d2f0290e39a08cfdf2245f14aec7dca60672b7c7458375f3cb3834e1067"),
					Signature: []*felt.Felt{
						utils.TestHexToFelt(t, "0x1"),
						utils.TestHexToFelt(t, "0x7bc0a22005a54ec6a005c1e89ab0201cbd0819621edd9fe4d5ef177a4ff33dd"),
						utils.TestHexToFelt(t, "0x13089e5f38de4ea98e9275be7fadc915946be15c14a8fed7c55202818527bea"),
					},
					Calldata: []*felt.Felt{
						utils.TestHexToFelt(t, "0x1"),
						utils.TestHexToFelt(t, "0x517567ac7026ce129c950e6e113e437aa3c83716cd61481c6bb8c5057e6923e"),
						utils.TestHexToFelt(t, "0xcaffbd1bd76bd7f24a3fa1d69d1b2588a86d1f9d2359b13f6a84b7e1cbd126"),
						utils.TestHexToFelt(t, "0x6"),
						utils.TestHexToFelt(t, "0x5265706f73736573734275696c64696e67"),
						utils.TestHexToFelt(t, "0x4"),
						utils.TestHexToFelt(t, "0x5"),
						utils.TestHexToFelt(t, "0x1b48"),
						utils.TestHexToFelt(t, "0x1"),
						utils.TestHexToFelt(t, "0xe52"),
					},
				},
			},
		},
	}

	testSet := map[string][]testSetType{
		"mock": {},
		"testnet": {
			{
				BlockID:                     WithBlockNumber(64159),
				ExpectedError:               nil,
				want:                        &fullBlockSepolia64159,
				LookupTxnPositionInExpected: 0,
				LookupTxnPositionInOriginal: 4,
			},
			{
				BlockID:                     WithBlockHash(utils.TestHexToFelt(t, "0x6df565874b2ea6a02d346a23f9efb0b26abbf5708b51bb12587f88a49052964")),
				ExpectedError:               nil,
				want:                        &fullBlockSepolia64159,
				LookupTxnPositionInExpected: 0,
				LookupTxnPositionInOriginal: 4,
			},
		},
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		spy := NewSpy(testConfig.provider.c)
		testConfig.provider.c = spy
		blockWithTxsInterface, err := testConfig.provider.BlockWithTxs(context.Background(), test.BlockID)
		require.NoError(t, err, "Unable to fetch the given block.")

		blockWithTxs, ok := blockWithTxsInterface.(*Block)
		require.True(t, ok, "Failed to assert the Interface as *Block.")
		require.Equal(t, blockWithTxs.BlockHash.String()[:2], "0x", "Block Hash should start with \"0x\".")
		require.NotEqual(t, len(blockWithTxs.Transactions), 0, "The number of transaction should not be 0.")

		invokeV1Want, ok := (*test.want).Transactions[test.LookupTxnPositionInExpected].(BlockInvokeTxnV1)
		require.True(t, ok, "Expected invoke v1 transaction.")

		invokeV1Block, ok := blockWithTxs.Transactions[test.LookupTxnPositionInOriginal].(BlockInvokeTxnV1)
		require.True(t, ok, "Expected invoke v1 transaction.")

		require.Equal(t, invokeV1Want.TransactionHash.String(), invokeV1Block.TransactionHash.String(), "Expected equal TransactionHash.")
		require.Equal(t, invokeV1Want.InvokeTxnV1.MaxFee.String(), invokeV1Block.InvokeTxnV1.MaxFee.String(), "Expected equal maxfee.")
		require.Equal(t, invokeV1Want.InvokeTxnV1.Calldata[1].String(), invokeV1Block.InvokeTxnV1.Calldata[1].String(), "Expected equal calldatas.")
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
					Nonce:         utils.TestHexToFelt(t, "0x49f6"),
					Version:       TransactionV3,
					Signature: []*felt.Felt{
						utils.TestHexToFelt(t, "0x7f14bb280b602d0c0e22e91ea5d987371554664f68c57c0acd16bf9f8be36b4"),
						utils.TestHexToFelt(t, "0x22c57ce8eb211c7fe0f04e7da338f579fcbc9e8997ec432fac7738c80fd56ad"),
					},
					Calldata: []*felt.Felt{
						utils.TestHexToFelt(t, "0x1"),
						utils.TestHexToFelt(t, "0x3b4a36565470dadabb8db607d894c2c5c74c97c1da2a02e68f7b2c1e816c546"),
						utils.TestHexToFelt(t, "0xc844fd57777b0cd7e75c8ea68deec0adf964a6308da7a58de32364b7131cc8"),
						utils.TestHexToFelt(t, "0x13"),
						utils.TestHexToFelt(t, "0x4e8ccea93954c392e93c650e4203c4bbb14bbd49172cf11ab344f3307ad8c"),
						utils.TestHexToFelt(t, "0x54c04"),
						utils.TestHexToFelt(t, "0x65c358022a33f92065dbc77f4ead316ca6140b77d3b9c85aeb533f3c6f0016"),
						utils.TestHexToFelt(t, "0x6600d885"),
						utils.TestHexToFelt(t, "0x304010200000000000000000000000000000000000000000000000000000000"),
						utils.TestHexToFelt(t, "0x4"),
						utils.TestHexToFelt(t, "0x4fdfb0a7f6"),
						utils.TestHexToFelt(t, "0x4fdfb0a7f6"),
						utils.TestHexToFelt(t, "0x4fe0bd9fe0"),
						utils.TestHexToFelt(t, "0x4fe0bd9fe0"),
						utils.TestHexToFelt(t, "0xa0663181cca796872"),
						utils.TestHexToFelt(t, "0x1"),
						utils.TestHexToFelt(t, "0x2"),
						utils.TestHexToFelt(t, "0x6b34b2fe5d97c7d5dde2bd56b8344cb40b984d88fa7b0d9cabec6bcf882e072"),
						utils.TestHexToFelt(t, "0x379c65ebd06927d01ab023e275bb4e45fb9ff9da8d27bf70bea6c2e5d51839b"),
						utils.TestHexToFelt(t, "0x2cb74dff29a13dd5d855159349ec92f943bacf0547ff3734e7d84a15d08cbc5"),
						utils.TestHexToFelt(t, "0xf9c67a4272ebcfcd527c103c7ee33946645e4df836a1c8c6f95c24d1e0ff94"),
						utils.TestHexToFelt(t, "0x3deb1e9016f2a77b97802dfc4f39824cbe08cafaca8adc29d0467fc119a6673"),
						utils.TestHexToFelt(t, "0x2e7dc996ebf724c1cf18d668fc3455df4245749ebc0724101cbc6c9cb13c962"),
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
		require.NoError(t, err, "Unable to fetch the given block.")

		blockWithTxs, ok := blockWithTxsInterface.(*Block)
		require.True(t, ok, "Failed to assert the Interface as *Block.")
		require.Equal(t, blockWithTxs.BlockHash.String()[:2], "0x", "Block Hash should start with \"0x\".")
		require.NotEqual(t, len(blockWithTxs.Transactions), 0, "The number of transaction should not be 0.")

		invokeV3Want, ok := (*test.want).Transactions[test.LookupTxnPositionInExpected].(BlockInvokeTxnV3)
		require.True(t, ok, "Expected invoke v3 transaction.")

		invokeV3Block, ok := blockWithTxs.Transactions[test.LookupTxnPositionInOriginal].(BlockInvokeTxnV3)
		require.True(t, ok, "Expected invoke v3 transaction.")

		require.Equal(t, invokeV3Want.TransactionHash.String(), invokeV3Block.TransactionHash.String(), "Expected equal TransactionHash.")
		require.Equal(t, invokeV3Want.InvokeTxnV3.Nonce, invokeV3Block.InvokeTxnV3.Nonce, "Expected equal nonce.")
		require.Equal(t, invokeV3Want.InvokeTxnV3.Calldata[1].String(), invokeV3Block.InvokeTxnV3.Calldata[1].String(), "Expected equal calldatas.")
	}
}

// TestBlockWithTxsAndDeployOrDeclare tests BlockWithTxs with Deploy or Declare TXN

// TestBlockWithTxsAndDeployOrDeclare is a test function that tests the functionality of the BlockWithTxs function.
// It creates a test configuration, defines a testSetType struct, and initializes three Blocks (fullBlockSepolia65204, fullBlockSepolia65083 and fullBlockSepoliai65212).
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
func TestBlockWithTxsAndDeployOrDeclare(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		BlockID                     BlockID
		ExpectedError               error
		LookupTxnPositionInOriginal int
		LookupTxnPositionInExpected int
		ExpectedBlockWithTxs        *Block
	}

	var fullBlockSepolia65204 = Block{
		BlockHeader: BlockHeader{
			BlockHash:        utils.TestHexToFelt(t, "0x2b0d32dbe49e0ffdc6d5cb36896198c242276ec63a60106a17963427f7cf10e"),
			ParentHash:       utils.TestHexToFelt(t, "0x7c703a06d7476055adaa518cb9682621755e7a95f144a54853ace320800e82e"),
			SequencerAddress: utils.TestHexToFelt(t, "0x1176a1bd84444c89232ec27754698e5d2e7e1a7f1539f12027f28b23ec9f3d8"),
			BlockNumber:      65204,
			NewRoot:          utils.TestHexToFelt(t, "0x1a24db193c0ebab47cb29c59ef277033d2e37dd909297073cb5c9594d892344"),
			Timestamp:        1715289398,
		},
		Status: "ACCEPTED_ON_L1",
		Transactions: []BlockTransaction{
			BlockDeployAccountTxn{
				TransactionHash: utils.TestHexToFelt(t, "0x26e30d2ed579c1ff575710d8ce29d9056e67ac08ab261a7221d384734d6ad5a"),
				DeployAccountTxn: DeployAccountTxn{
					Type:                TransactionType_DeployAccount,
					Version:             TransactionV1,
					Nonce:               utils.TestHexToFelt(t, "0x0"),
					MaxFee:              utils.TestHexToFelt(t, "0x4865413596"),
					ContractAddressSalt: utils.TestHexToFelt(t, "0x5ffef5f00daec09457836121dcb7a8bcae97080716a292d3d8ee899d0b2597e"),
					ClassHash:           utils.TestHexToFelt(t, "0x13bfe114fb1cf405bfc3a7f8dbe2d91db146c17521d40dcf57e16d6b59fa8e6"),
					ConstructorCalldata: []*felt.Felt{
						utils.TestHexToFelt(t, "0x5ffef5f00daec09457836121dcb7a8bcae97080716a292d3d8ee899d0b2597e"),
					},
					Signature: []*felt.Felt{
						utils.TestHexToFelt(t, "0x561486f318c01d7d47434875d685a505c282804616ed3c69e91fb59731739fc"),
						utils.TestHexToFelt(t, "0x1beb073a4cac8b344d22e2d2833ea30e140b3701936b3a4af69a12a2f47394b"),
						utils.TestHexToFelt(t, "0x816dd0297efc55dc1e7559020a3a825e81ef734b558f03c83325d4da7e6253"),
						utils.TestHexToFelt(t, "0x0"),
						utils.TestHexToFelt(t, "0x0"),
						utils.TestHexToFelt(t, "0x0"),
						utils.TestHexToFelt(t, "0x0"),
						utils.TestHexToFelt(t, "0x0"),
						utils.TestHexToFelt(t, "0x0"),
						utils.TestHexToFelt(t, "0x0"),
						utils.TestHexToFelt(t, "0x0"),
						utils.TestHexToFelt(t, "0x0"),
						utils.TestHexToFelt(t, "0x534e5f5345504f4c4941"),
						utils.TestHexToFelt(t, "0x1f58325efbae915a4ea70270c536fa1c84c48eedf8ae3b970e0c0270694ddad"),
						utils.TestHexToFelt(t, "0x2c7adff34d5e2b1c4205c80638df0e3a5b6831ff1d9643f6e7b8bcbbe1f71f3"),
					},
				},
			},
		},
	}

	var fullBlockSepolia65083 = Block{
		BlockHeader: BlockHeader{
			BlockHash:        utils.TestHexToFelt(t, "0x549770b5b74df90276277ff7a11af881c998dffa452f4156f14446db6005174"),
			ParentHash:       utils.TestHexToFelt(t, "0x4de1acdff24acba2a537ef651ec8f790e5c0321f92f1115b272a6f2f2d637e8"),
			SequencerAddress: utils.TestHexToFelt(t, "0x1176a1bd84444c89232ec27754698e5d2e7e1a7f1539f12027f28b23ec9f3d8"),
			BlockNumber:      65083,
			NewRoot:          utils.TestHexToFelt(t, "0x7df132d80333fd54a3a059e0cc6e851bda52cc72d0437a8f13a1b0809a17499"),
			Timestamp:        1715244926,
		},
		Status: "ACCEPTED_ON_L1",
		Transactions: []BlockTransaction{
			BlockDeclareTxnV1{
				TransactionHash: utils.TestHexToFelt(t, "0x3c7817502dac2dc90198c6b64b85f3700507d74c75e08af85164e1b35e3a8b5"),
				DeclareTxnV1: DeclareTxnV1{
					Type:    TransactionType_Declare,
					MaxFee:  utils.TestHexToFelt(t, "0xde0b6b3a7640000"),
					Version: TransactionV1,
					Signature: []*felt.Felt{
						utils.TestHexToFelt(t, "0x998b12c82d208af7c4b820626f2f7e015b8ee33ef5ae44e8a04f5254977865"),
						utils.TestHexToFelt(t, "0x55c341329a881afb29462ab32dcebb16d35c56021c3595bbba01b5c563f66fe"),
					},
					Nonce:         utils.TestHexToFelt(t, "0x713"),
					ClassHash:     utils.TestHexToFelt(t, "0x6c5a3f54e8bbfe07167c87f8ace70629573e05c385fe4bea69bc1d323acb8d3"),
					SenderAddress: utils.TestHexToFelt(t, "0x2cc631ca0c544639f6e4403b8f3611696a3d831e8157ea1c946e35429c7ac31"),
				},
			},
		},
	}

	var fullBlockSepoliai65212 = Block{
		BlockHeader: BlockHeader{
			BlockHash:        utils.TestHexToFelt(t, "0x31b785f0f8b258f7b164d13ecc02dc4e06a1c67174f1e39f9368d1b5af43ae"),
			ParentHash:       utils.TestHexToFelt(t, "0x6133d377632092e31b0adad5a6496c8468cb3cb53de10f3ecdfc748f57cf9e3"),
			SequencerAddress: utils.TestHexToFelt(t, "0x1176a1bd84444c89232ec27754698e5d2e7e1a7f1539f12027f28b23ec9f3d8"),
			BlockNumber:      65212,
			NewRoot:          utils.TestHexToFelt(t, "0x2d867e8a51807513d3372cf2658f03b0c212caba0fcad4f036e3023215b27a"),
			Timestamp:        1715292305,
		},
		Status: "ACCEPTED_ON_L1",
		Transactions: []BlockTransaction{
			BlockDeclareTxnV2{
				TransactionHash: utils.TestHexToFelt(t, "0x3e68091ecacab5a880ae8d9847d7b87408bbf05270ded34e082577acfdc3770"),
				DeclareTxnV2: DeclareTxnV2{
					Type:    TransactionType_Declare,
					MaxFee:  utils.TestHexToFelt(t, "0x1108942d5866"),
					Version: TransactionV2,
					Signature: []*felt.Felt{
						utils.TestHexToFelt(t, "0x40cedf50ffc6866050b63e9576333da68ed51e258c6ddfa7e22d5e557f71961"),
						utils.TestHexToFelt(t, "0x4b33458a2005664fbd204331c078df3ac5a16288285ba0fcf1b6fd9c9d9257a"),
					},
					Nonce:             utils.TestHexToFelt(t, "0x2f"),
					ClassHash:         utils.TestHexToFelt(t, "0x190c14ee771ed5c77a006545659baa33d22c083473c8ba136cf1614f3e545b0"),
					CompiledClassHash: utils.TestHexToFelt(t, "0x3bd32e501b720756105de78796bebd5f4412da986811dff38a114be296e6120"),
					SenderAddress:     utils.TestHexToFelt(t, "0x3b2d6d0edcbdbdf6548d2b79531263628887454a0a608762c71056172d36240"),
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
				BlockID:                     WithBlockHash(utils.TestHexToFelt(t, "0x2b0d32dbe49e0ffdc6d5cb36896198c242276ec63a60106a17963427f7cf10e")),
				ExpectedError:               nil,
				LookupTxnPositionInOriginal: 5,
				LookupTxnPositionInExpected: 0,
				ExpectedBlockWithTxs:        &fullBlockSepolia65204,
			},
			{
				BlockID:                     WithBlockNumber(65204),
				ExpectedError:               nil,
				LookupTxnPositionInOriginal: 5,
				LookupTxnPositionInExpected: 0,
				ExpectedBlockWithTxs:        &fullBlockSepolia65204,
			},
			{
				BlockID:                     WithBlockNumber(65212),
				ExpectedError:               nil,
				LookupTxnPositionInOriginal: 24,
				LookupTxnPositionInExpected: 0,
				ExpectedBlockWithTxs:        &fullBlockSepoliai65212,
			},
			{
				BlockID:                     WithBlockNumber(65083),
				ExpectedError:               nil,
				LookupTxnPositionInOriginal: 26,
				LookupTxnPositionInExpected: 0,
				ExpectedBlockWithTxs:        &fullBlockSepolia65083,
			},
			{
				BlockID:                     WithBlockNumber(65212),
				ExpectedError:               nil,
				LookupTxnPositionInOriginal: 24,
				LookupTxnPositionInExpected: 0,
				ExpectedBlockWithTxs:        &fullBlockSepoliai65212,
			},
			{
				BlockID:                     WithBlockNumber(65083),
				ExpectedError:               nil,
				LookupTxnPositionInOriginal: 26,
				LookupTxnPositionInExpected: 0,
				ExpectedBlockWithTxs:        &fullBlockSepolia65083,
			},
		},
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		spy := NewSpy(testConfig.provider.c)
		testConfig.provider.c = spy
		blockWithTxsInterface, err := testConfig.provider.BlockWithTxs(context.Background(), test.BlockID)
		assert.Equal(t, test.ExpectedError, err, "BlockWithTxHashes doesn't match the expected error.")

		if test.ExpectedError != nil && blockWithTxsInterface == nil {
			continue
		}

		blockWithTxs, ok := blockWithTxsInterface.(*Block)
		require.True(t, ok, fmt.Sprintf("Expecting *rpc.Block, instead %T", blockWithTxsInterface))

		diff, err := spy.Compare(blockWithTxs, false)
		require.NoError(t, err, "Expected to compare the BlockWithTxs.")

		if diff != "FullMatch" {
			_, err = spy.Compare(blockWithTxs, false)
			require.NoError(t, err, "Unable to compare the count.")
		}

		require.True(t, strings.HasPrefix(blockWithTxs.BlockHash.String(), "0x"), "Block Hash should start with \"0x\", instead: %s", blockWithTxs.BlockHash)
		require.NotEmpty(t, blockWithTxs.Transactions, "the number of transactions should not be 0")

		if test.ExpectedBlockWithTxs != nil {
			if test.ExpectedBlockWithTxs.BlockHash == &felt.Zero {
				continue
			}
			require.True(t, cmp.Equal(test.ExpectedBlockWithTxs.Transactions[test.LookupTxnPositionInExpected], blockWithTxs.Transactions[test.LookupTxnPositionInOriginal]),
				"the expected transaction blocks to match, instead: %s", cmp.Diff(test.ExpectedBlockWithTxs.Transactions[test.LookupTxnPositionInExpected], blockWithTxs.Transactions[test.LookupTxnPositionInOriginal]))
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
				BlockID:       WithBlockNumber(30000),
				ExpectedCount: 4,
			},
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
		require.NoError(t, err, "Unable to fetch the given block.")

		diff, err := spy.Compare(count, false)
		require.NoError(t, err, "Unable to compare the count.")

		_, err = spy.Compare(count, true)
		require.NoError(t, err, "Unable to compare the count.")

		require.Equal(t, "FullMatch", diff, "structure expecting to be FullMatch, instead %s", diff)

		require.Equal(t, test.ExpectedCount, count, fmt.Sprintf("structure expecting %d, instead: %d", test.ExpectedCount, count))
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
				StartBlock: 52959,
				EndBlock:   52960,
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
				_, okv0 := v.(BlockInvokeTxnV0)
				_, okv1 := v.(BlockInvokeTxnV1)
				_, okv3 := v.(BlockInvokeTxnV3)
				_, okl1 := v.(BlockL1HandlerTxn)
				_, okdec0 := v.(BlockDeclareTxnV0)
				_, okdec1 := v.(BlockDeclareTxnV1)
				_, okdec2 := v.(BlockDeclareTxnV2)
				_, okdec3 := v.(BlockDeclareTxnV3)
				_, okdep := v.(BlockDeployTxn)
				_, okdepac := v.(BlockDeployAccountTxn)
				if !okv0 && !okv1 && !okv3 && !okl1 && !okdec0 && !okdec1 && !okdec2 && !okdec3 && !okdep && !okdepac {
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
				BlockID: WithBlockNumber(30000),
				ExpectedStateUpdateOutput: StateUpdateOutput{
					BlockHash: utils.TestHexToFelt(t, "0x62ab7b3ade3e7c26d0f50cb539c621b679e07440685d639904663213f906938"),
					NewRoot:   utils.TestHexToFelt(t, "0x491250c959067f21177f50cfdfede2bd9c8f2597f4ed071dbdba4a7ee3dabec"),
					PendingStateUpdate: PendingStateUpdate{
						OldRoot: utils.TestHexToFelt(t, "0x19aa982a75263d4c4de4cc4c5d75c3dec32e00b95bef7bbb4d17762a0b138af"),
						StateDiff: StateDiff{
							StorageDiffs: []ContractStorageDiffItem{
								{
									Address: utils.TestHexToFelt(t, "0xe5cc6f2b6d34979184b88334eb64173fe4300cab46ecd3229633fcc45c83d4"),
									StorageEntries: []StorageEntry{
										{
											Key:   utils.TestHexToFelt(t, "0x1813aac5f5e7799684c6dc33e51f44d3627fd748c800724a184ed5be09b713e"),
											Value: utils.TestHexToFelt(t, "0x630b4197"),
										},
									},
								},
							},
						},
					},
				},
			},
		},
		"testnet": {
			{
				BlockID: WithBlockNumber(30000),
				ExpectedStateUpdateOutput: StateUpdateOutput{
					BlockHash: utils.TestHexToFelt(t, "0x62ab7b3ade3e7c26d0f50cb539c621b679e07440685d639904663213f906938"),
					NewRoot:   utils.TestHexToFelt(t, "0x491250c959067f21177f50cfdfede2bd9c8f2597f4ed071dbdba4a7ee3dabec"),
					PendingStateUpdate: PendingStateUpdate{
						OldRoot: utils.TestHexToFelt(t, "0x1d2922de7bb14766d0c3aa323876d9f5a4b1733f6dc199bbe596d06dd8f70e4"),
						StateDiff: StateDiff{
							StorageDiffs: []ContractStorageDiffItem{
								{
									Address: utils.TestHexToFelt(t, "0xe5cc6f2b6d34979184b88334eb64173fe4300cab46ecd3229633fcc45c83d4"),
									StorageEntries: []StorageEntry{
										{
											Key:   utils.TestHexToFelt(t, "0x1813aac5f5e7799684c6dc33e51f44d3627fd748c800724a184ed5be09b713e"),
											Value: utils.TestHexToFelt(t, "0x630b4197"),
										},
									},
								},
								{
									Address: utils.TestHexToFelt(t, "0x49d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"),
									StorageEntries: []StorageEntry{
										{
											Key:   utils.TestHexToFelt(t, "0x7b3303ee433d39925f7c289cd2048052a2d8e2d653bdd7cdfa6a6ab8365445d"),
											Value: utils.TestHexToFelt(t, "0x462893a80b9b5834"),
										},
										{
											Key:   utils.TestHexToFelt(t, "0x5496768776e3db30053404f18067d81a6e06f5a2b0de326e21298fd9d569a9a"),
											Value: utils.TestHexToFelt(t, "0x1bc48439cb7402fb6"),
										},
									},
								},
								{
									Address: utils.TestHexToFelt(t, "0x36031daa264c24520b11d93af622c848b2499b66b41d611bac95e13cfca131a"),
									StorageEntries: []StorageEntry{
										{
											Key:   utils.TestHexToFelt(t, "0x6f64efd140d53af83432093bb6c3d5e8db645bd89feead6dda806955f68ef2a"),
											Value: utils.TestHexToFelt(t, "0x3df78515979000000000000000000000000065bfec4b"),
										},
										{
											Key:   utils.TestHexToFelt(t, "0x143dae8bc0e9898f65cb1eb84f16bfb9cb09431972541141677721dd541f055"),
											Value: utils.TestHexToFelt(t, "0x5f35296000000000000000000000000065bfec4c"),
										},
										{
											Key:   utils.TestHexToFelt(t, "0x2d04b0419a0e89f6b4dabc3dc19b087e71f0dec9f1785606f00517d3468636b"),
											Value: utils.TestHexToFelt(t, "0x5f5f2e4000000000000000000000000065bfec4c"),
										},
										{
											Key:   utils.TestHexToFelt(t, "0x55c3ad197a2fa1dce3a999ae803099406fab085f187b926e7e1f0e38592043d"),
											Value: utils.TestHexToFelt(t, "0x3985cb98c08000000000000000000000000065bfec4b"),
										},
										{
											Key:   utils.TestHexToFelt(t, "0x8653303a2624a587179380e17d7876d346aea7f02dbd57782950500ea7276e"),
											Value: utils.TestHexToFelt(t, "0x3e076b4dfa2000000000000000000000000065bfec4b"),
										},
										{
											Key:   utils.TestHexToFelt(t, "0x56041f8991ff7eff841647cfda1f1cfb9e7321c5a96c53d4a5072497de6b50f"),
											Value: utils.TestHexToFelt(t, "0x23b8c472000000000000000000000000065bfec4c"),
										},
										{
											Key:   utils.TestHexToFelt(t, "0x6a6414ca66551a2324e436ed37d069f1660ef01bc3fe90497fc729ee60781b8"),
											Value: utils.TestHexToFelt(t, "0x3511a8db1d000000000000000000000000065bfec4b"),
										},
										{
											Key:   utils.TestHexToFelt(t, "0x437f038e1991939def57775a3405a3b6f0c0830f09d0e6cfc309393950fa773"),
											Value: utils.TestHexToFelt(t, "0x3d5e14753a000000000000000000000000065bfec4c"),
										},
										{
											Key:   utils.TestHexToFelt(t, "0x7b4de97b546ed17a0d490dab334867e9383e029411c268a8902768b6da6a2eb"),
											Value: utils.TestHexToFelt(t, "0x5f38d0b000000000000000000000000065bfec4b"),
										},
										{
											Key:   utils.TestHexToFelt(t, "0x28e86558bd7c5a9c26fceeafb9570eb7b3011db4a9ff813b318f91129935c37"),
											Value: utils.TestHexToFelt(t, "0xf423c000000000000000000000000065bfec4c"),
										},
										{
											Key:   utils.TestHexToFelt(t, "0x1b3f3d264a9c63c581333d4b97c556b6f20f9a1abf64c7f71e04b35df62cc70"),
											Value: utils.TestHexToFelt(t, "0xf407f000000000000000000000000065bfec4c"),
										},
									},
								},
								{
									Address: utils.TestHexToFelt(t, "0x4718f5a0fc34cc1af16a1cdee98ffb20c31f5cd61d6ab07201858f4287c938d"),
									StorageEntries: []StorageEntry{
										{
											Key:   utils.TestHexToFelt(t, "0x5496768776e3db30053404f18067d81a6e06f5a2b0de326e21298fd9d569a9a"),
											Value: utils.TestHexToFelt(t, "0x5d8da32bae8513cfa"),
										},
										{
											Key:   utils.TestHexToFelt(t, "0x295c615dc08b568dce79348e5dd16f45bc6458ddb026f09e16ce03f3c68e12e"),
											Value: utils.TestHexToFelt(t, "0x218cbd49b5dafd0cec6"),
										},
									},
								},
								{
									Address: utils.TestHexToFelt(t, "0x47ad6a25df680763e5663bd0eba3d2bfd18b24b1e8f6bd36b71c37433c63ed0"),
									StorageEntries: []StorageEntry{
										{
											Key:   utils.TestHexToFelt(t, "0x38b0933d0e83013f5bd5aee82962149fed820534bfc3978a5180646208e7937"),
											Value: utils.TestHexToFelt(t, "0x7df9d41833d7cf135b059ccc165ed4332cc32ac3eddf3f6239594731b0d8c8"),
										},
										{
											Key:   utils.TestHexToFelt(t, "0x38b0933d0e83013f5bd5aee82962149fed820534bfc3978a5180646208e7936"),
											Value: utils.TestHexToFelt(t, "0x3b2f128039c288928ff492627eba9969d760b7fd0b16f3d39aa18f1f8744765"),
										},
									},
								},
								{
									Address: utils.TestHexToFelt(t, "0x1"),
									StorageEntries: []StorageEntry{
										{
											Key:   utils.TestHexToFelt(t, "0x7526"),
											Value: utils.TestHexToFelt(t, "0x32f159b038c06f9829d8ee63db1556a3390265b0b49b89c48235b6f77326339"),
										},
									},
								},
							},
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
		require.NoError(t, err, "Unable to fetch the given block.")

		require.Equal(t,
			test.ExpectedStateUpdateOutput.BlockHash.String(),
			stateUpdate.BlockHash.String(),
			fmt.Sprintf("structure expecting %s, instead: %s", test.ExpectedStateUpdateOutput.BlockHash.String(), stateUpdate.BlockHash.String()),
		)
	}
}
