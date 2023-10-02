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

// TestBlockNumber tests BlockNumber and check the returned value is strictly positive.
//
// It does the following:
// - Initializes the test configuration.
// - Defines a test set with different environments.
// - Iterates through the test set.
// - Creates a spy object with the test configuration provider.
// - Retrieves the block number from the test configuration provider.
// - Checks if the retrieved block number matches the expected value.
// - Checks if the block number is greater than 3000.
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
// The function does the following:
// - Sets up a test configuration.
// - Defines a test set.
// - Iterates over the test set.
// - Initializes a spy.
// - Calls the BlockHashAndNumber function.
// - Checks for the expected error.
// - Compares the result with a spy.
// - Checks the block number.
// - Checks the block hash.
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

// TestBlockWithTxHashes is a unit test function for testing the BlockWithTxHashes function.
//
// It tests the behavior of the BlockWithTxHashes function in different scenarios.
// The function takes a test object of type testSetType and performs the following steps:
// - Sets up the test configuration.
// - Calls the BlockWithTxHashes function with the provided test.BlockID.
// - Compares the result with the expected error and block values.
//
// Parameters:
// - t: A testing.T object used for reporting test failures and logging.
//
// Return type: None.
func TestBlockWithTxHashes(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		BlockID                          BlockID
		ExpectedError                    error
		ExpectedBlockWithTxHashes        *BlockTxHashes
		ExpectedPendingBlockWithTxHashes *PendingBlockTxHashes
	}

	var blockGoerli310370 = BlockTxHashes{
		BlockHeader: BlockHeader{
			BlockHash:        utils.TestHexToFelt(t, "0x6c2fe3db009a2e008c2d65fca14204f3405cb74742fcf685f02473acaf70c72"),
			ParentHash:       utils.TestHexToFelt(t, "0x1ce6fa8ef59dfa1ad8f7ce7c3a4e6752e2d8ae6274f8257345f680e6ae0b5b5"),
			SequencerAddress: utils.TestHexToFelt(t, "0x46a89ae102987331d369645031b49c27738ed096f2789c24449966da4c6de6b"),
			BlockNumber:      310370,
			NewRoot:          utils.TestHexToFelt(t, "0x5cd7a08312635206c0210b8c90e61ceac27cb09629064e12266fd79e4c05a3d"),
			Timestamp:        1661450764,
		},
		Status: "ACCEPTED_ON_L1",
		Transactions: utils.TestHexArrToFelt(t, []string{
			"0x40c82f79dd2bc1953fc9b347a3e7ab40fe218ed5740bf4e120f74e8a3c9ac99",
			"0x28981b14353a28bc46758dff412ac544d16f2ffc8dde31867855592ea054ab1",
			"0x41176c650076712f1618a141fc1cf9a8c39f0d9548a3458f29cf363310a1e72",
			"0x43cd66f3ddbfbf681ab99bb57bf9d94c83d6e9b586bdbde78ab2deb0328ebd5",
			"0x7602cfebe4f3cb3ef4c8b8c6d7dda2efaf4a500723020066f5db50acd5095cd",
			"0x2612f3f870ee7e7617d4f9efdc41fa8fd571f9720b059b1aa14c1bf15d3a92a",
			"0x1a7810a6c68adf0621ed384d915409c936efa0c9d436683ea0cf7ea171719b",
			"0x26683aeef3e9d9bcc1f0d45a5f0b67d0aa1919726524b2a8dc59504dacfd1f4",
			"0x1d374aa073435cdde1ec1caf972f7c175fd23438bb220848e71720e00fd7474",
			"0xfc13eabaa2f38981e68bb010370cad7a7d0b65a59101ec816042adca0d6841",
			"0x672d007224128b99bcc145cd3dbd8930a944b6a5fff5c27e3b158a6ff701509",
			"0x24795cbca6d2eba941082cea3f686bc86ef27dd46fdf84b32f9ba25bbeddb28",
			"0x69281a4dd58c260a06b3266554c0cf1a4f19b79d8488efef2a1f003d67506ed",
			"0x62211cc3c94d612b580eb729410e52277f838f962d91af91fb2b0526704c04d",
			"0x5e4128b7680db32de4dff7bc57cb11c9f222752b1f875e84b29785b4c284e2a",
			"0xdb8ad2b7d008fd2ad7fba4315b193032dee85e17346c80276a2e08c7f09f80",
			"0x67b9541ca879abc29fa24a0fa070285d1899fc044159521c827f6b6aa09bbd6",
			"0x5d9c0ab1d4ed6e9376c8ab45ee02b25dd0adced12941aafe8ce37369d19d9c2",
			"0x4e52da53e23d92d9818908aeb104b007ea24d3cd4a5aa43144d2db1011e314f",
			"0x6cc05f5ab469a3675acb5885c274d5143dca75dd9835c582f59e85ab0642d39",
			"0x561ed983d1d9c37c964a96f80ccaf3de772e2b73106d6f49dd7c3f7ed8483d9",
		}),
	}

	txHashes := utils.TestHexArrToFelt(t, []string{
		"0x40c82f79dd2bc1953fc9b347a3e7ab40fe218ed5740bf4e120f74e8a3c9ac99",
		"0x28981b14353a28bc46758dff412ac544d16f2ffc8dde31867855592ea054ab1",
	})
	blockHash := utils.TestHexToFelt(t, "0xbeef")

	testSet := map[string][]testSetType{
		"mock": {
			{
				BlockID: BlockID{Tag: "latest"},
				ExpectedPendingBlockWithTxHashes: &PendingBlockTxHashes{
					ParentHash:       &felt.Zero,
					Timestamp:        123,
					SequencerAddress: &felt.Zero,
					Transactions:     txHashes,    
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
				BlockID:       WithBlockTag("latest"),
				ExpectedError: nil,
			},
			{
				BlockID:                   WithBlockHash(utils.TestHexToFelt(t, "0x6c2fe3db009a2e008c2d65fca14204f3405cb74742fcf685f02473acaf70c72")),
				ExpectedError:             nil,
				ExpectedBlockWithTxHashes: &blockGoerli310370,
			},
			{
				BlockID:                   WithBlockNumber(310370),
				ExpectedError:             nil,
				ExpectedBlockWithTxHashes: &blockGoerli310370,
			},
		},
		"mainnet": {},
	}[testEnv]
	
	for _, test := range testSet {
		spy := NewSpy(testConfig.provider.c)
		testConfig.provider.c = spy
		result, err := testConfig.provider.BlockWithTxHashes(context.Background(), test.BlockID)
		if err != test.ExpectedError {
			t.Fatal("BlockWithTxHashes match the expected error:", err)
		}
		switch resultType := result.(type) {
		case *BlockTxHashes:
			block, ok := result.(*BlockTxHashes)
			if !ok {
				t.Fatalf("should return *BlockTxHashes, instead: %T\n", result)
			}
			if test.ExpectedError != nil {
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

// TestBlockWithTxsAndInvokeTXNV0 is a test function that tests the functionality of BlockWithTxs and InvokeTXNV0.
//
// The function takes in a testing.T object as a parameter and performs the following steps:
// 1. Sets up the test configuration using the beforeEach function.
// 2. Defines a testSetType struct that holds the input parameters, expected error, and expected output for each test case.
// 3. Initializes a fullBlockGoerli310370 variable of type Block with test data.
// 4. Defines a testSet map that maps a test environment (mock, testnet, mainnet) to an array of testSetType.
// 5. Loops through each test case in the testSet and performs the following steps:
//     a. Initializes a spy object of type Spy using the testConfig provider.
//     b. Calls the BlockWithTxs function of the testConfig provider with the BlockID from the test case.
//     c. Checks if the returned error matches the expected error.
//     d. Checks if the returned blockWithTxs is of type *Block.
//     e. Compares the returned blockWithTxs with the expected output using the spy's Compare function.
//     f. Checks if the BlockHash of the returned blockWithTxs starts with "0x".
//     g. Checks if the number of transactions in the returned blockWithTxs is greater than 0.
//     h. If the expected output is not nil, checks if the TransactionHash, MaxFee, and EntryPointSelector of the BlockInvokeTxnV0 transaction at the specified position in the returned blockWithTxs match the expected values.
//
// The function does not return any values.
func TestBlockWithTxsAndInvokeTXNV0(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		BlockID                     BlockID
		ExpectedError               error
		LookupTxnPositionInOriginal int
		LookupTxnPositionInExpected int
		want                        *Block
	}

	var fullBlockGoerli310370 = Block{
		BlockHeader: BlockHeader{
			BlockHash:        utils.TestHexToFelt(t, "0x6c2fe3db009a2e008c2d65fca14204f3405cb74742fcf685f02473acaf70c72"),
			ParentHash:       utils.TestHexToFelt(t, "0x1ce6fa8ef59dfa1ad8f7ce7c3a4e6752e2d8ae6274f8257345f680e6ae0b5b5"),
			SequencerAddress: utils.TestHexToFelt(t, "0x46a89ae102987331d369645031b49c27738ed096f2789c24449966da4c6de6b"),
			BlockNumber:      310370,
			NewRoot:          utils.TestHexToFelt(t, "0x5cd7a08312635206c0210b8c90e61ceac27cb09629064e12266fd79e4c05a3d"),
			Timestamp:        1661450764,
		},
		Status: "ACCEPTED_ON_L1",
		Transactions: []BlockTransaction{

			BlockInvokeTxnV0{
				TransactionHash: utils.TestHexToFelt(t, "0x40c82f79dd2bc1953fc9b347a3e7ab40fe218ed5740bf4e120f74e8a3c9ac99"),
				InvokeTxnV0: InvokeTxnV0{
					Type:    "INVOKE",
					MaxFee:  utils.TestHexToFelt(t, "0xde0b6b3a7640000"),
					Version: TransactionV0,
					Signature: []*felt.Felt{
						utils.TestHexToFelt(t, "0x7bc0a22005a54ec6a005c1e89ab0201cbd0819621edd9fe4d5ef177a4ff33dd"),
						utils.TestHexToFelt(t, "0x13089e5f38de4ea98e9275be7fadc915946be15c14a8fed7c55202818527bea"),
					},
					FunctionCall: FunctionCall{
						ContractAddress:    utils.TestHexToFelt(t, "0x2e28403d7ee5e337b7d456327433f003aa875c29631906908900058c83d8cb6"),
						EntryPointSelector: utils.TestHexToFelt(t, "0x15d40a3d6ca2ac30f4031e42be28da9b056fef9bb7357ac5e85627ee876e5ad"),
						Calldata: []*felt.Felt{
							utils.TestHexToFelt(t, "0x1"),
							utils.TestHexToFelt(t, "0x33830ce413e4c096eef81b5e6ffa9b9f5d963f57b8cd63c9ae4c839c383c1a6"),
							utils.TestHexToFelt(t, "0x2db698626ed7f60212e1ce6e99afb796b6b423d239c3f0ecef23e840685e866"),
							utils.TestHexToFelt(t, "0x0"),
							utils.TestHexToFelt(t, "0x2"),
							utils.TestHexToFelt(t, "0x2"),
							utils.TestHexToFelt(t, "0x61c6e7484657e5dc8b21677ffa33e4406c0600bba06d12cf1048fdaa55bdbc3"),
							utils.TestHexToFelt(t, "0x6307b990"),
							utils.TestHexToFelt(t, "0x2b81"),
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
				BlockID:       WithBlockHash(utils.TestHexToFelt(t, "0x6c2fe3db009a2e008c2d65fca14204f3405cb74742fcf685f02473acaf70c72")),
				ExpectedError: nil,
				want:          &fullBlockGoerli310370,
			},
			{
				BlockID:       WithBlockNumber(310370),
				ExpectedError: nil,
				want:          &fullBlockGoerli310370,
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

// TestBlockWithTxsAndDeployOrDeclare is a test function that tests the behavior of the BlockWithTxs method with Deploy or Declare TXN.
//
// It sets up a test configuration, defines a testSetType struct, and initializes a fullBlockGoerli310843 variable of type Block.
// Then, it defines a testSet map with testSetType values.
// The function iterates over the testSet and performs the following actions:
// - Creates a new spy instance.
// - Sets the provider.c field of the testConfig to the spy instance.
// - Calls the BlockWithTxs method of the provider with the given test.BlockID argument and assigns the result to blockWithTxsInterface.
// - Compares the error returned by the BlockWithTxs method to the expected error. If they don't match, the function terminates with a fatal error.
// - Checks if the expected error is not nil and the blockWithTxsInterface is nil. If true, the function continues with the next iteration.
// - Asserts that the blockWithTxsInterface is of type *Block. If not, the function terminates with a fatal error.
// - Compares the blockWithTxs with the spy instance. If they don't match, the function terminates with a fatal error.
// - Asserts that the blockWithTxs.BlockHash starts with "0x". If not, the function terminates with a fatal error.
// - Asserts that the number of transactions in blockWithTxs is not zero. If it is zero, the function terminates with a fatal error.
// - Checks if test.ExpectedBlockWithTxs is not nil. If true, it compares the expected transaction block at test.LookupTxnPositionInExpected with the actual transaction block at test.LookupTxnPositionInOriginal. If they don't match, the function terminates with a fatal error.
//
// The function does not have any parameters and does not return any values.
func TestBlockWithTxsAndDeployOrDeclare(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		BlockID                     BlockID
		ExpectedError               error
		LookupTxnPositionInOriginal int
		LookupTxnPositionInExpected int
		ExpectedBlockWithTxs        *Block
	}

	// To do : re-add test for deploy account transaction
	var fullBlockGoerli310843 = Block{
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

	var fullBlockGoerli848622 = Block{
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

	var fullBlockGoerli849399 = Block{
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
				ExpectedBlockWithTxs:        &fullBlockGoerli310843,
			},
			{
				BlockID:                     WithBlockNumber(310843),
				ExpectedError:               nil,
				LookupTxnPositionInOriginal: 14,
				LookupTxnPositionInExpected: 0,
				ExpectedBlockWithTxs:        &fullBlockGoerli310843,
			},
			{
				BlockID:                     WithBlockNumber(849399),
				ExpectedError:               nil,
				LookupTxnPositionInOriginal: 71,
				LookupTxnPositionInExpected: 0,
				ExpectedBlockWithTxs:        &fullBlockGoerli849399,
			},
			{
				BlockID:                     WithBlockNumber(848622),
				ExpectedError:               nil,
				LookupTxnPositionInOriginal: 6,
				LookupTxnPositionInExpected: 0,
				ExpectedBlockWithTxs:        &fullBlockGoerli848622,
			},
			{
				BlockID:                     WithBlockNumber(849399),
				ExpectedError:               nil,
				LookupTxnPositionInOriginal: 71,
				LookupTxnPositionInExpected: 0,
				ExpectedBlockWithTxs:        &fullBlockGoerli849399,
			},
			{
				BlockID:                     WithBlockNumber(848622),
				ExpectedError:               nil,
				LookupTxnPositionInOriginal: 6,
				LookupTxnPositionInExpected: 0,
				ExpectedBlockWithTxs:        &fullBlockGoerli848622,
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
			spy.Compare(blockWithTxs, false)
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

// TestBlockTransactionCount tests the BlockTransactionCount function.
//
// It sets up a test environment and creates a testSet containing mock,
// testnet, and mainnet data. For each test in the testSet, it initializes a
// spy provider and sets the testConfig provider to the spy. Then it calls
// BlockTransactionCount with the test's BlockID and checks if the returned
// count matches the test's ExpectedCount. It also compares the returned count
// with the spy's count to ensure a full match. If any of the checks fail, it
// logs an error and fails the test.
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
			spy.Compare(count, true)
			t.Fatal("structure expecting to be FullMatch, instead", diff)
		}
		if count != test.ExpectedCount {
			t.Fatalf("structure expecting %d, instead: %d", test.ExpectedCount, count)
		}
	}
}

// TestCaptureUnsupportedBlockTxn tests the capture of unsupported block transactions.
//
// This function takes no parameters.
// It does not return anything.
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

// TestBlockWithTxsAndInvokeTXNV1 is a test function for testing the block with transactions and invoking TXNV1.
//
// The function does not take any parameters.
// It does not return anything.
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

// TestStateUpdate is a unit test function that tests the StateUpdate function.
//
// The function takes no parameters.
// It does not return any value.
// TODO: this is not implemented yet with pathfinder as you can see from the
// [code](https://github.com/eqlabs/pathfinder/blob/927183552dad6dcdfebac16c8c1d2baf019127b1/crates/pathfinder/rpc_examples.sh#L37)
// check when it is and test when it is the case.
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
