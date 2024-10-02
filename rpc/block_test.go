package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/NethermindEth/starknet.go/utils"
	"github.com/stretchr/testify/require"
)

// TestBlockNumber is a test function to check the behavior of the BlockNumber function and check if there is no errors.
//
// Parameters:
// - t: the testing object for running the test cases
// Returns:
//
//	none
func TestBlockNumber(t *testing.T) {
	testConfig := beforeEach(t)

	blockNumber, err := testConfig.provider.BlockNumber(context.Background())
	require.NoError(t, err, "BlockNumber should not return an error")
	if testEnv == "mock" {
		require.Equal(t, uint64(1234), blockNumber)
	}
}

// TestBlockHashAndNumber is a test function that tests the BlockHashAndNumber function and check if there is no errors.
// Parameters:
// - t: the testing object for running the test cases
// Returns:
//
//	none
func TestBlockHashAndNumber(t *testing.T) {
	testConfig := beforeEach(t)

	blockHashAndNumber, err := testConfig.provider.BlockHashAndNumber(context.Background())
	require.NoError(t, err, "BlockHashAndNumber should not return an error")
	require.True(t, strings.HasPrefix(blockHashAndNumber.BlockHash.String(), "0x"), "current block hash should return a string starting with 0x")

	if testEnv == "mock" {
		require.Equal(t, &BlockHashAndNumberOutput{BlockNumber: 1234, BlockHash: utils.RANDOM_FELT}, blockHashAndNumber)
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

	var blockSepolia64159 BlockTxHashes
	block, err := os.ReadFile("tests/blockWithHashes/sepoliaBlockWithHashes64159.json")
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(block, &blockSepolia64159))

	txHashes := utils.TestHexArrToFelt(t, []string{
		"0x5754961d70d6f39d0e2c71a1a4ff5df0a26b1ceda4881ca82898994379e1e73",
		"0x692381bba0e8505a8e0b92d0f046c8272de9e65f050850df678a0c10d8781d",
	})
	fakeFelt := utils.TestHexToFelt(t, "0xbeef")

	testSet := map[string][]testSetType{
		"mock": {
			{
				BlockID:     BlockID{Tag: "latest"},
				ExpectedErr: nil,
				ExpectedPendingBlockWithTxHashes: &PendingBlockTxHashes{
					PendingBlockHeader{
						ParentHash:       fakeFelt,
						Timestamp:        123,
						SequencerAddress: fakeFelt},
					txHashes,
				},
			},
			{
				BlockID: BlockID{Hash: fakeFelt},
				ExpectedBlockWithTxHashes: &BlockTxHashes{
					BlockHeader: BlockHeader{
						BlockHash:        fakeFelt,
						ParentHash:       fakeFelt,
						Timestamp:        124,
						SequencerAddress: fakeFelt},
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
				BlockID:     WithBlockTag("pending"),
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
		result, err := testConfig.provider.BlockWithTxHashes(context.Background(), test.BlockID)
		require.Equal(t, test.ExpectedErr, err, "Error in BlockWithTxHashes")
		switch resultType := result.(type) {
		case *BlockTxHashes:
			block, ok := result.(*BlockTxHashes)
			require.Truef(t, ok, "should return *BlockTxHashes, instead: %T\n", result)

			if test.ExpectedErr != nil {
				continue
			}

			require.Truef(t, strings.HasPrefix(block.BlockHash.String(), "0x"), "Block Hash should start with \"0x\", instead: %s", block.BlockHash)
			require.NotEmpty(t, block.Transactions, "the number of transactions should not be 0")

			if test.ExpectedBlockWithTxHashes != nil {
				require.Exactly(t, test.ExpectedBlockWithTxHashes, block)
			}
		case *PendingBlockTxHashes:
			pBlock, ok := result.(*PendingBlockTxHashes)
			require.Truef(t, ok, "should return *PendingBlockTxHashes, instead: %T\n", result)

			if test.ExpectedPendingBlockWithTxHashes == nil {
				validatePendingBlockHeader(t, &pBlock.PendingBlockHeader)
			} else {
				require.Exactly(t, test.ExpectedPendingBlockWithTxHashes, pBlock)
			}
		default:
			t.Fatalf("unexpected block type, found: %T\n", resultType)
		}
	}
}

// TestBlockWithTxs tests the BlockWithTxs function.
//
// The function tests the BlockWithTxs function by setting up a test configuration and a test set type.
// It then initializes Block type variables and invokes the BlockWithTxs function with different test scenarios.
// The function checks if the BlockWithTxs function returns the correct block data.
// It also verifies the block hash, the number of transactions in the block, and the details of a specific transaction.
//
// Parameters:
// - t: The t testing object
// Returns:
//
//	none
func TestBlockWithTxs(t *testing.T) {
	testConfig := beforeEach(t)
	require := require.New(t)

	type testSetType struct {
		BlockID              BlockID
		ExpectedBlock        *Block
		ExpectedPendingBlock *PendingBlock
		InvokeV0Index        int //TODO: implement mainnet testcases as Sepolia doesn't contains V0 transactions
		InvokeV1Index        int
		InvokeV3Index        int
		DeclareV0Index       int //TODO: implement mainnet testcases as Sepolia doesn't contains V0 transactions
		DeclareV1Index       int
		DeclareV2Index       int
		DeclareV3Index       int //TODO: implement testcase
		DeployAccountV1Index int
		DeployAccountV3Index int //TODO: implement testcase
		L1HandlerV0Index     int
	}

	var fullBlockSepolia65083 Block
	read, err := os.ReadFile("tests/block/sepoliaBlockTxs65083.json")
	require.NoError(err)
	require.NoError(json.Unmarshal(read, &fullBlockSepolia65083))

	var fullBlockSepolia122476 Block
	read, err = os.ReadFile("tests/block/sepoliaBlockTxs122476.json")
	require.NoError(err)
	require.NoError(json.Unmarshal(read, &fullBlockSepolia122476))

	testSet := map[string][]testSetType{
		"mock": {
			{
				BlockID: WithBlockTag("latest"),
			},
			{
				BlockID: WithBlockTag("pending"),
				ExpectedPendingBlock: &PendingBlock{
					PendingBlockHeader{
						ParentHash:       utils.RANDOM_FELT,
						Timestamp:        123,
						SequencerAddress: utils.RANDOM_FELT,
					},
					nil,
				},
			},
			{
				BlockID:       WithBlockNumber(65083),
				ExpectedBlock: &fullBlockSepolia65083,
				InvokeV1Index: 1,
			},
			{
				BlockID:       WithBlockHash(utils.TestHexToFelt(t, "0x549770b5b74df90276277ff7a11af881c998dffa452f4156f14446db6005174")),
				ExpectedBlock: &fullBlockSepolia65083,
				InvokeV1Index: 1,
			},
		},
		"testnet": {
			{
				BlockID: WithBlockTag("latest"),
			},
			{
				BlockID: WithBlockTag("pending"),
			},
			{
				BlockID:              WithBlockNumber(65083),
				ExpectedBlock:        &fullBlockSepolia65083,
				InvokeV1Index:        1,
				InvokeV3Index:        3,
				DeclareV1Index:       26,
				DeclareV2Index:       25,
				DeployAccountV1Index: 42,
			},
			{
				BlockID:          WithBlockHash(utils.TestHexToFelt(t, "0x56a71e0443d2fbfaa91b1000e830b516ca0d4a424abb9c970d23801957dbfa3")),
				ExpectedBlock:    &fullBlockSepolia122476,
				L1HandlerV0Index: 4,
			},
		},
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		blockWithTxsInterface, err := testConfig.provider.BlockWithTxs(context.Background(), test.BlockID)
		require.NoError(err, "Unable to fetch the given block.")

		switch block := blockWithTxsInterface.(type) {
		case *PendingBlock:
			if test.ExpectedPendingBlock == nil {
				validatePendingBlockHeader(t, &block.PendingBlockHeader)
			} else {
				require.Exactly(test.ExpectedPendingBlock, block)
			}
		case *Block:
			if test.ExpectedBlock == nil {
				require.Equal(block.BlockHash.String()[:2], "0x", "Block Hash should start with \"0x\".")
				require.NotEmpty(block.Transactions, "The number of transaction should not be 0.")
			} else {
				require.Exactly(test.ExpectedBlock, block)

				//validates an BlockInvokeV1 transaction
				if test.InvokeV1Index > 0 {
					invokeV1Expected, ok := (*test.ExpectedBlock).Transactions[test.InvokeV1Index].(BlockInvokeTxnV1)
					require.True(ok, "Expected invoke v1 transaction.")
					invokeV1Block, ok := block.Transactions[test.InvokeV1Index].(BlockInvokeTxnV1)
					require.True(ok, "Expected invoke v1 transaction.")

					require.Exactly(invokeV1Expected, invokeV1Block)
				}

				//validates an BlockInvokeV3 transaction
				if test.InvokeV3Index > 0 {
					invokeV3Expected, ok := (*test.ExpectedBlock).Transactions[test.InvokeV3Index].(BlockInvokeTxnV3)
					require.True(ok, "Expected invoke v3 transaction.")
					invokeV3Block, ok := block.Transactions[test.InvokeV3Index].(BlockInvokeTxnV3)
					require.True(ok, "Expected invoke v3 transaction.")

					require.Exactly(invokeV3Expected, invokeV3Block)
				}

				//validates an BlockDeclareV1 transaction
				if test.DeclareV1Index > 0 {
					declareV1Expected, ok := (*test.ExpectedBlock).Transactions[test.DeclareV1Index].(BlockDeclareTxnV1)
					require.True(ok, "Expected declare v1 transaction.")
					declareV1Block, ok := block.Transactions[test.DeclareV1Index].(BlockDeclareTxnV1)
					require.True(ok, "Expected declare v1 transaction.")

					require.Exactly(declareV1Expected, declareV1Block)
				}

				//validates an BlockDeclareV2 transaction
				if test.DeclareV2Index > 0 {
					declareV2Expected, ok := (*test.ExpectedBlock).Transactions[test.DeclareV2Index].(BlockDeclareTxnV2)
					require.True(ok, "Expected declare v2 transaction.")
					declareV2Block, ok := block.Transactions[test.DeclareV2Index].(BlockDeclareTxnV2)
					require.True(ok, "Expected declare v2 transaction.")

					require.Exactly(declareV2Expected, declareV2Block)
				}

				//validates an BlockDeployAccountV1 transaction
				if test.DeployAccountV1Index > 0 {
					deployAccountV1Expected, ok := (*test.ExpectedBlock).Transactions[test.DeployAccountV1Index].(BlockDeployAccountTxn)
					require.True(ok, "Expected declare v2 transaction.")
					deployAccountV1Block, ok := block.Transactions[test.DeployAccountV1Index].(BlockDeployAccountTxn)
					require.True(ok, "Expected declare v2 transaction.")

					require.Exactly(deployAccountV1Expected, deployAccountV1Block)
				}

				//validates an BlockL1HandlerV0 transaction
				if test.L1HandlerV0Index > 0 {
					l1HandlerV0Expected, ok := (*test.ExpectedBlock).Transactions[test.L1HandlerV0Index].(BlockL1HandlerTxn)
					require.True(ok, "Expected L1 handler transaction.")
					l1HandlerV0Block, ok := block.Transactions[test.L1HandlerV0Index].(BlockL1HandlerTxn)
					require.True(ok, "Expected L1 handler transaction.")

					require.Exactly(l1HandlerV0Expected, l1HandlerV0Block)
				}
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
		ExpectedError error
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
			{
				BlockID:       WithBlockNumber(7338746823462834783),
				ExpectedError: ErrBlockNotFound,
			},
		},
		"mainnet": {},
	}[testEnv]
	for _, test := range testSet {
		count, err := testConfig.provider.BlockTransactionCount(context.Background(), test.BlockID)
		if err != nil {
			require.EqualError(t, test.ExpectedError, err.Error())
		} else {
			require.Equalf(t, test.ExpectedCount, count, "structure expecting %d, instead: %d", test.ExpectedCount, count)
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

func validatePendingBlockHeader(t *testing.T, pBlock *PendingBlockHeader) {
	require.NotZero(t, pBlock.ParentHash)
	require.NotZero(t, pBlock.Timestamp)
	require.NotZero(t, pBlock.SequencerAddress)
	require.NotZero(t, pBlock.L1GasPrice)
	require.NotZero(t, pBlock.StarknetVersion)
	require.NotZero(t, pBlock.L1DataGasPrice)
	require.NotNil(t, pBlock.L1DAMode)
}
