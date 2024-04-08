package rpc

import (
	"context"
	"fmt"
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
func TestBlockWithTxsAndInvokeTXNV0(t *testing.T) {
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
			BlockHash:        utils.TestHexToFelt(t, "0x2488a7149327b4dcd200f05a131911bd44f946021539648642eaa7d6e82f289"),
			ParentHash:       utils.TestHexToFelt(t, "0x2adc07a26d70e72a16775e26c45074f0216bc2e86e35bfe53743968480e4c1b"),
			SequencerAddress: utils.TestHexToFelt(t, "0x1176a1bd84444c89232ec27754698e5d2e7e1a7f1539f12027f28b23ec9f3d8"),
			BlockNumber:      30436,
			NewRoot:          utils.TestHexToFelt(t, "0x4bbb6bec9488d70d9c9e96862cf50e22331a5e8d7b33a56712f56cd04c16e06"),
			Timestamp:        1661450764,
		},
		Status: "ACCEPTED_ON_L1",
		Transactions: []BlockTransaction{

			BlockInvokeTxnV0{
				TransactionHash: utils.TestHexToFelt(t, "0x10d2059db6d261fee740b515ed8b9c50955f03dd43c4729b24dc63278641926"),
				InvokeTxnV0: InvokeTxnV0{
					Type:    "INVOKE",
					MaxFee:  utils.TestHexToFelt(t, "0x470de4df820000"),
					Version: TransactionV0,
					Signature: []*felt.Felt{
						utils.TestHexToFelt(t, "0x7fb440b1dee35c5259bf10d55782bc973434d195bb5c8b95ac7d3d8e2a8a0e4"),
						utils.TestHexToFelt(t, "0x13089e5f38de4ea98e9275be7fadc915946be15c14a8fed7c55202818527bea"),
					},
					FunctionCall: FunctionCall{
						ContractAddress:    nil,
						EntryPointSelector: nil,
						Calldata: []*felt.Felt{
							utils.TestHexToFelt(t, "0x1"),
							utils.TestHexToFelt(t, "0x3fe8e4571772bbe0065e271686bd655efd1365a5d6858981e582f82f2c10313"),
							utils.TestHexToFelt(t, "0x2468d193cd15b621b24c2a602b8dbcfa5eaa14f88416c40c09d7fd12592cb4b"),
							utils.TestHexToFelt(t, "0x0"),
						},
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
				BlockID:       WithBlockHash(utils.TestHexToFelt(t, "0x2488a7149327b4dcd200f05a131911bd44f946021539648642eaa7d6e82f289")),
				ExpectedError: nil,
				want:          &fullBlockSepolia30436,
			},
			{
				BlockID:       WithBlockNumber(30436),
				ExpectedError: nil,
				want:          &fullBlockSepolia30436,
			},
		},
		"mainnet": {},
	}[testEnv]

	for i, test := range testSet {
		fmt.Println("Test case", i)
		fmt.Println("BlockID is", test.BlockID)
		fmt.Println("LookuptxnPositionInOriginal is", test.LookupTxnPositionInOriginal)
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

			invokeV0Want, ok := (*test.want).Transactions[test.LookupTxnPositionInExpected].(BlockInvokeTxnV0)
			if !ok {
				t.Fatal("expected invoke v0 transaction")
			}
			invokeV0Block, ok := blockWithTxs.Transactions[test.LookupTxnPositionInOriginal].(BlockInvokeTxnV0)
			if !ok {
				t.Fatal("expected invoke v0 transaction")
			}
			require.Equal(t, invokeV0Want.TransactionHash, invokeV0Block.TransactionHash, "expected equal TransactionHash")
			require.Equal(t, invokeV0Want.InvokeTxnV0.MaxFee, invokeV0Block.InvokeTxnV0.MaxFee, "expected equal maxfee")
			require.Equal(t, invokeV0Want.InvokeTxnV0.EntryPointSelector, invokeV0Block.InvokeTxnV0.EntryPointSelector, "expected equal eps")

		}

	}
}

// TestBlockWithTxsAndDeployOrDeclare tests BlockWithTxs with Deploy or Declare TXN

// TestBlockWithTxsAndDeployOrDeclare is a test function that tests the functionality of the BlockWithTxs function.
// It creates a test configuration, defines a testSetType struct, and initializes three Blocks (fullBlockSepolia310843, fullBlockSepolia848622 and fullBlockSepolia849399).
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

	// TODO : re-add test for deploy account transaction
	var fullBlockSepolia310843 = Block{
		BlockHeader: BlockHeader{
			BlockHash:        utils.TestHexToFelt(t, "0x424fba26a7760b63895abe0c366c2d254cb47090c6f9e91ba2b3fa0824d4fc9"),
			ParentHash:       utils.TestHexToFelt(t, "0x30e34dedf00bb35a9076b2b0f50a5a74fd2501f62094b6e687277be6ef3d444"),
			SequencerAddress: utils.TestHexToFelt(t, "0x46a89ae102987331d369645031b49c27738ed096f2789c24449966da4c6de6b"),
			BlockNumber:      310843,
			NewRoot:          utils.TestHexToFelt(t, "0x32bd4ff21288c898d4d3b6a7aea4ebdb3f1c7089cd52bde98316b4ecb8a50be"),
			Timestamp:        1661486036,
		},
		Status: "ACCEPTED_ON_L1",
		Transactions: []BlockTransaction{
			BlockDeployTxn{
				TransactionHash: utils.TestHexToFelt(t, "0x35bd2978d2061b3463498f83c09322ed6a82e4b2a188506525e272a7adcdf6a"),
				DeployTxn: DeployTxn{
					ClassHash: utils.TestHexToFelt(t, "0x1ca349f9721a2bf05012bb475b404313c497ca7d6d5f80c03e98ff31e9867f5"),
					ConstructorCalldata: []*felt.Felt{
						utils.TestHexToFelt(t, "0x31ad196615d50956d98be085eb1774624106a6936c7c38696e730d2a6df735a"),
						utils.TestHexToFelt(t, "0x736affc32af71f8d361c855b38ffef58ec151bd8361a3b160017b90ada1068e"),
					},
					ContractAddressSalt: utils.TestHexToFelt(t, "0x4241e90ee6a33a1e2e70b088f7e4acfb3d6630964c1a85e96fa96acd56dcfcf"),

					Type:    "DEPLOY",
					Version: TransactionV0,
				},
			},
		},
	}

	var fullBlockSepolia848622 = Block{
		BlockHeader: BlockHeader{
			BlockHash:        utils.TestHexToFelt(t, "0x32964e2e407bb9e71b2de8d9d9829b0537df7c4624e1816e6cece80781ab9cc"),
			ParentHash:       utils.TestHexToFelt(t, "0xecbed6cfe85c77f2f8acefe2effbda817f71ca7457f7ece8262d65cc87a9f7"),
			SequencerAddress: utils.TestHexToFelt(t, "0x1176a1bd84444c89232ec27754698e5d2e7e1a7f1539f12027f28b23ec9f3d8"),
			BlockNumber:      848622,
			NewRoot:          utils.TestHexToFelt(t, "07c4302f09f6a72129679378e9b8d6c67774c5c4e80b1fc186da114f71637b2e"),
			Timestamp:        1692416283,
		},
		Status: "ACCEPTED_ON_L1",
		Transactions: []BlockTransaction{
			BlockDeclareTxnV1{
				TransactionHash: utils.TestHexToFelt(t, "0x5ad2f85499ea92d33d4a44c8cd4640d1ee4e25c3ee6df0bdf0a76c12c052f0a"),
				DeclareTxnV1: DeclareTxnV1{
					Type:          TransactionType_Declare,
					MaxFee:        utils.TestHexToFelt(t, "0x27a64c6e425"),
					Version:       TransactionV1,
					Signature:     []*felt.Felt{utils.TestHexToFelt(t, "0x1454ab28f0bf18f0fd8002bc92169e6443feba6c605728c86850c0dcc9f6f9a"), utils.TestHexToFelt(t, "0xf545949c899ff1d16c61629996e898db2697a2e3e7fa9071b016500ca5c1d1")},
					Nonce:         utils.TestHexToFelt(t, "0x333"),
					ClassHash:     utils.TestHexToFelt(t, "0x681076f783aa2b3faec6ce80bb5485a260ed1672007925e1d502b003aff2232"),
					SenderAddress: utils.TestHexToFelt(t, "0x45dba6ce6a4dc3d2f31aa6da5f51007f1e43e84a1e62c4481bac5454dea4e6d"),
				},
			},
		},
	}

	var fullBlockSepolia849399 = Block{
		BlockHeader: BlockHeader{
			BlockHash:        utils.TestHexToFelt(t, "0x6e5b26127400bac0cd1f3c2ab6e76850ec457c71b1f2fc7cda755bee8a1102a"),
			ParentHash:       utils.TestHexToFelt(t, "0x7cd085d4ab95b3307303cb836ab49c0fbc8d1f9befdcfdc65292d99c9466d05"),
			SequencerAddress: utils.TestHexToFelt(t, "0x1176a1bd84444c89232ec27754698e5d2e7e1a7f1539f12027f28b23ec9f3d8"),
			BlockNumber:      849399,
			NewRoot:          utils.TestHexToFelt(t, "0x239a44410e78665f41f7a65ef3b5ed244ce411965498a83f80f904e22df1045"),
			Timestamp:        1692560305,
		},
		Status: "ACCEPTED_ON_L1",
		Transactions: []BlockTransaction{
			BlockDeclareTxnV2{
				TransactionHash: utils.TestHexToFelt(t, "0x45d04652ba51685b7b82fc17b3d5741a7c43992369c0b0aebd60916fa23b9b2"),
				DeclareTxnV2: DeclareTxnV2{
					Type:              TransactionType_Declare,
					MaxFee:            utils.TestHexToFelt(t, "0x50c8f30287c"),
					Version:           TransactionV2,
					Signature:         []*felt.Felt{utils.TestHexToFelt(t, "0x6be01a56087382337a29fd4577dd20fd82cc9f38f69b8d19e07fc101c3c5ad9"), utils.TestHexToFelt(t, "0x4c633a5582d3932fbfcea8abd45c7453e88a562f1a38877b9575d6a6b926ea2")},
					Nonce:             utils.TestHexToFelt(t, "0xd"),
					ClassHash:         utils.TestHexToFelt(t, "0x6fda8f6630f44571cd6b398795351b37daf27adacbf6fe9357bd23ad19b22f3"),
					CompiledClassHash: utils.TestHexToFelt(t, "0x4380d7c6511f81668530570a8b07bd2148808f90e681bb769549ec4faafef65"),
					SenderAddress:     utils.TestHexToFelt(t, "0x6ef69146f56205e27624a9933f31d6009198c1ea480070a790f16a5d928be92"),
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
				BlockID:                     WithBlockHash(utils.TestHexToFelt(t, "0x424fba26a7760b63895abe0c366c2d254cb47090c6f9e91ba2b3fa0824d4fc9")),
				ExpectedError:               nil,
				LookupTxnPositionInOriginal: 14,
				LookupTxnPositionInExpected: 0,
				ExpectedBlockWithTxs:        &fullBlockSepolia310843,
			},
			{
				BlockID:                     WithBlockNumber(310843),
				ExpectedError:               nil,
				LookupTxnPositionInOriginal: 14,
				LookupTxnPositionInExpected: 0,
				ExpectedBlockWithTxs:        &fullBlockSepolia310843,
			},
			{
				BlockID:                     WithBlockNumber(849399),
				ExpectedError:               nil,
				LookupTxnPositionInOriginal: 71,
				LookupTxnPositionInExpected: 0,
				ExpectedBlockWithTxs:        &fullBlockSepolia849399,
			},
			{
				BlockID:                     WithBlockNumber(848622),
				ExpectedError:               nil,
				LookupTxnPositionInOriginal: 6,
				LookupTxnPositionInExpected: 0,
				ExpectedBlockWithTxs:        &fullBlockSepolia848622,
			},
			{
				BlockID:                     WithBlockNumber(849399),
				ExpectedError:               nil,
				LookupTxnPositionInOriginal: 71,
				LookupTxnPositionInExpected: 0,
				ExpectedBlockWithTxs:        &fullBlockSepolia849399,
			},
			{
				BlockID:                     WithBlockNumber(848622),
				ExpectedError:               nil,
				LookupTxnPositionInOriginal: 6,
				LookupTxnPositionInExpected: 0,
				ExpectedBlockWithTxs:        &fullBlockSepolia848622,
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
				BlockID:       WithBlockNumber(300000),
				ExpectedCount: 23,
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

// TestBlockWithTxsAndInvokeTXNV1 is a test function that tests the behavior of the BlockWithTxsAndInvokeTXNV1 function.
//
// Parameters:
// - t: the testing object for running the test cases
// Returns:
//
//	none
//
// TODO: Find a block with such a Txn
func TestBlockWithTxsAndInvokeTXNV1(t *testing.T) {
	_ = beforeEach(t)

	type testSetType struct {
		check bool
	}
	testSet := map[string][]testSetType{
		"mock": {},
		"testnet": {
			{
				check: false,
			},
		},
		"mainnet": {},
	}[testEnv]
	for _, test := range testSet {
		if test.check {
			t.Fatalf("error running test: %v", ErrNotImplemented)
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
