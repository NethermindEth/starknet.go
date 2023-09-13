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

// TestBlockNumber tests BlockNumber and check the returned value is strictly positive
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

// TestBlockHashAndNumber tests BlockHashAndNumber and check the returned value is strictly positive
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

// TestBlockWithTxHashes tests BlockWithTxHashes
func TestBlockWithTxHashes(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		BlockID                          BlockID
		ExpectedError                    error
		ExpectedBlockWithTxHashes        *Block
		ExpectedPendingBlockWithTxHashes *PendingBlock
	}

	var blockGoerli310370 = Block{
		BlockHeader: BlockHeader{
			BlockHash:        utils.TestHexToFelt(t, "0x6c2fe3db009a2e008c2d65fca14204f3405cb74742fcf685f02473acaf70c72"),
			ParentHash:       utils.TestHexToFelt(t, "0x1ce6fa8ef59dfa1ad8f7ce7c3a4e6752e2d8ae6274f8257345f680e6ae0b5b5"),
			SequencerAddress: utils.TestHexToFelt(t, "0x46a89ae102987331d369645031b49c27738ed096f2789c24449966da4c6de6b"),
			BlockNumber:      310370,
			NewRoot:          utils.TestHexToFelt(t, "0x5cd7a08312635206c0210b8c90e61ceac27cb09629064e12266fd79e4c05a3d"),
			Timestamp:        1661450764,
		},
		Status: "ACCEPTED_ON_L1",
		Transactions: []Transaction{
			TransactionHash{TransactionHash: utils.TestHexToFelt(t, "0x40c82f79dd2bc1953fc9b347a3e7ab40fe218ed5740bf4e120f74e8a3c9ac99")},
			TransactionHash{TransactionHash: utils.TestHexToFelt(t, "0x28981b14353a28bc46758dff412ac544d16f2ffc8dde31867855592ea054ab1")},
			TransactionHash{TransactionHash: utils.TestHexToFelt(t, "0x41176c650076712f1618a141fc1cf9a8c39f0d9548a3458f29cf363310a1e72")},
			TransactionHash{TransactionHash: utils.TestHexToFelt(t, "0x43cd66f3ddbfbf681ab99bb57bf9d94c83d6e9b586bdbde78ab2deb0328ebd5")},
			TransactionHash{TransactionHash: utils.TestHexToFelt(t, "0x7602cfebe4f3cb3ef4c8b8c6d7dda2efaf4a500723020066f5db50acd5095cd")},
			TransactionHash{TransactionHash: utils.TestHexToFelt(t, "0x2612f3f870ee7e7617d4f9efdc41fa8fd571f9720b059b1aa14c1bf15d3a92a")},
			TransactionHash{TransactionHash: utils.TestHexToFelt(t, "0x1a7810a6c68adf0621ed384d915409c936efa0c9d436683ea0cf7ea171719b")},
			TransactionHash{TransactionHash: utils.TestHexToFelt(t, "0x26683aeef3e9d9bcc1f0d45a5f0b67d0aa1919726524b2a8dc59504dacfd1f4")},
			TransactionHash{TransactionHash: utils.TestHexToFelt(t, "0x1d374aa073435cdde1ec1caf972f7c175fd23438bb220848e71720e00fd7474")},
			TransactionHash{TransactionHash: utils.TestHexToFelt(t, "0xfc13eabaa2f38981e68bb010370cad7a7d0b65a59101ec816042adca0d6841")},
			TransactionHash{TransactionHash: utils.TestHexToFelt(t, "0x672d007224128b99bcc145cd3dbd8930a944b6a5fff5c27e3b158a6ff701509")},
			TransactionHash{TransactionHash: utils.TestHexToFelt(t, "0x24795cbca6d2eba941082cea3f686bc86ef27dd46fdf84b32f9ba25bbeddb28")},
			TransactionHash{TransactionHash: utils.TestHexToFelt(t, "0x69281a4dd58c260a06b3266554c0cf1a4f19b79d8488efef2a1f003d67506ed")},
			TransactionHash{TransactionHash: utils.TestHexToFelt(t, "0x62211cc3c94d612b580eb729410e52277f838f962d91af91fb2b0526704c04d")},
			TransactionHash{TransactionHash: utils.TestHexToFelt(t, "0x5e4128b7680db32de4dff7bc57cb11c9f222752b1f875e84b29785b4c284e2a")},
			TransactionHash{TransactionHash: utils.TestHexToFelt(t, "0xdb8ad2b7d008fd2ad7fba4315b193032dee85e17346c80276a2e08c7f09f80")},
			TransactionHash{TransactionHash: utils.TestHexToFelt(t, "0x67b9541ca879abc29fa24a0fa070285d1899fc044159521c827f6b6aa09bbd6")},
			TransactionHash{TransactionHash: utils.TestHexToFelt(t, "0x5d9c0ab1d4ed6e9376c8ab45ee02b25dd0adced12941aafe8ce37369d19d9c2")},
			TransactionHash{TransactionHash: utils.TestHexToFelt(t, "0x4e52da53e23d92d9818908aeb104b007ea24d3cd4a5aa43144d2db1011e314f")},
			TransactionHash{TransactionHash: utils.TestHexToFelt(t, "0x6cc05f5ab469a3675acb5885c274d5143dca75dd9835c582f59e85ab0642d39")},
			TransactionHash{TransactionHash: utils.TestHexToFelt(t, "0x561ed983d1d9c37c964a96f80ccaf3de772e2b73106d6f49dd7c3f7ed8483d9")},
		},
	}

	testSet := map[string][]testSetType{
		"mock": {{
			BlockID: BlockID{Tag: "latest"},
			ExpectedPendingBlockWithTxHashes: &PendingBlock{
				ParentHash:       &felt.Zero,
				Timestamp:        123,
				SequencerAddress: &felt.Zero,
			},
		},
			{
				BlockID: BlockID{Hash: &felt.Zero},
				ExpectedBlockWithTxHashes: &Block{
					BlockHeader: BlockHeader{
						ParentHash:       &felt.Zero,
						Timestamp:        123,
						SequencerAddress: &felt.Zero},
					Status: BlockStatus_AcceptedOnL1,
				},
			}},
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
		switch resultBlock := result.(type) {
		case Block:
			block, ok := result.(*Block)
			if !ok {
				t.Fatalf("should return *Block, instead: %T\n", result)
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

				if !cmp.Equal(*test.ExpectedBlockWithTxHashes, *block) {
					t.Fatalf("the expected transaction blocks to match, instead: %s", cmp.Diff(test.ExpectedBlockWithTxHashes, block))
				}
			}
		case PendingBlock:
			require.Equal(t, resultBlock.ParentHash, test.ExpectedPendingBlockWithTxHashes.ParentHash, "Error in PendingBlock ParentHash")
			require.Equal(t, resultBlock.SequencerAddress, test.ExpectedPendingBlockWithTxHashes.SequencerAddress, "Error in PendingBlock SequencerAddress")
			require.Equal(t, resultBlock.Timestamp, test.ExpectedPendingBlockWithTxHashes.Timestamp, "Error in PendingBlock Timestamp")
		}

	}
}

// TestBlockWithTxsAndInvokeTXNV0 tests block with Invoke TXN V0
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
		Transactions: []Transaction{
			InvokeTxnV0{
				CommonTransaction: CommonTransaction{
					TransactionHash: utils.TestHexToFelt(t, "0x40c82f79dd2bc1953fc9b347a3e7ab40fe218ed5740bf4e120f74e8a3c9ac99"),
					BroadcastedTxnCommonProperties: BroadcastedTxnCommonProperties{
						Type:    "INVOKE",
						MaxFee:  utils.TestHexToFelt(t, "0xde0b6b3a7640000"),
						Version: TransactionV0,
						Signature: []*felt.Felt{
							utils.TestHexToFelt(t, "0x7bc0a22005a54ec6a005c1e89ab0201cbd0819621edd9fe4d5ef177a4ff33dd"),
							utils.TestHexToFelt(t, "0x13089e5f38de4ea98e9275be7fadc915946be15c14a8fed7c55202818527bea"),
						},
						Nonce: nil,
					},
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

			invokeV0Want, ok := (*test.want).Transactions[test.LookupTxnPositionInExpected].(InvokeTxnV0)
			if !ok {
				t.Fatal("expected invoke v0 transaction")
			}
			invokeV0Block, ok := blockWithTxs.Transactions[test.LookupTxnPositionInOriginal].(InvokeTxnV0)
			if !ok {
				t.Fatal("expected invoke v0 transaction")
			}

			require.Equal(t, invokeV0Want.Hash(), invokeV0Block.Hash(), "expected equal hash")
			require.Equal(t, invokeV0Want.Nonce, invokeV0Block.Nonce, "expected equal nonce")
			require.Equal(t, invokeV0Want.MaxFee, invokeV0Block.MaxFee, "expected equal maxfee")
			require.Equal(t, invokeV0Want.EntryPointSelector, invokeV0Block.EntryPointSelector, "expected equal eps")

		}

	}
}

// TestBlockWithTxsAndDeployOrDeclare tests BlockWithTxs with Deploy or Declare TXN
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
		Transactions: []Transaction{
			DeployTxn{
				TransactionHash: utils.TestHexToFelt(t, "0x35bd2978d2061b3463498f83c09322ed6a82e4b2a188506525e272a7adcdf6a"),
				ClassHash:       utils.TestHexToFelt(t, "0x1ca349f9721a2bf05012bb475b404313c497ca7d6d5f80c03e98ff31e9867f5"),
				DeployTransactionProperties: DeployTransactionProperties{
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
		Transactions: []Transaction{
			DeclareTxnV1{
				CommonTransaction: CommonTransaction{
					TransactionHash: utils.TestHexToFelt(t, "0x5ad2f85499ea92d33d4a44c8cd4640d1ee4e25c3ee6df0bdf0a76c12c052f0a"),
					BroadcastedTxnCommonProperties: BroadcastedTxnCommonProperties{
						Type:      TransactionType_Declare,
						MaxFee:    utils.TestHexToFelt(t, "0x27a64c6e425"),
						Version:   TransactionV1,
						Signature: []*felt.Felt{utils.TestHexToFelt(t, "0x1454ab28f0bf18f0fd8002bc92169e6443feba6c605728c86850c0dcc9f6f9a"), utils.TestHexToFelt(t, "0xf545949c899ff1d16c61629996e898db2697a2e3e7fa9071b016500ca5c1d1")},
						Nonce:     utils.TestHexToFelt(t, "0x333"),
					},
				},
				ClassHash:     utils.TestHexToFelt(t, "0x681076f783aa2b3faec6ce80bb5485a260ed1672007925e1d502b003aff2232"),
				SenderAddress: utils.TestHexToFelt(t, "0x45dba6ce6a4dc3d2f31aa6da5f51007f1e43e84a1e62c4481bac5454dea4e6d"),
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
		Transactions: []Transaction{
			DeclareTxnV2{
				CommonTransaction: CommonTransaction{
					TransactionHash: utils.TestHexToFelt(t, "0x45d04652ba51685b7b82fc17b3d5741a7c43992369c0b0aebd60916fa23b9b2"),
					BroadcastedTxnCommonProperties: BroadcastedTxnCommonProperties{
						Type:      TransactionType_Declare,
						MaxFee:    utils.TestHexToFelt(t, "0x50c8f30287c"),
						Version:   TransactionV2,
						Signature: []*felt.Felt{utils.TestHexToFelt(t, "0x6be01a56087382337a29fd4577dd20fd82cc9f38f69b8d19e07fc101c3c5ad9"), utils.TestHexToFelt(t, "0x4c633a5582d3932fbfcea8abd45c7453e88a562f1a38877b9575d6a6b926ea2")},
						Nonce:     utils.TestHexToFelt(t, "0xd"),
					},
				},
				ClassHash:         utils.TestHexToFelt(t, "0x6fda8f6630f44571cd6b398795351b37daf27adacbf6fe9357bd23ad19b22f3"),
				CompiledClassHash: utils.TestHexToFelt(t, "0x4380d7c6511f81668530570a8b07bd2148808f90e681bb769549ec4faafef65"),
				SenderAddress:     utils.TestHexToFelt(t, "0x6ef69146f56205e27624a9933f31d6009198c1ea480070a790f16a5d928be92"),
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

// TestBlockTransactionCount tests BlockTransactionCount
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
				_, okv1 := v.(InvokeTxnV1)
				_, okv0 := v.(InvokeTxnV0)
				_, okl1 := v.(L1HandlerTxn)
				_, okdec0 := v.(DeclareTxnV0)
				_, okdec1 := v.(DeclareTxnV1)
				_, okdec2 := v.(DeclareTxnV2)
				_, okdep := v.(DeployTxn)
				_, okdepac := v.(DeployAccountTxn)
				if !okv0 && !okv1 && !okl1 && !okdec0 && !okdec1 && !okdec2 && !okdep && !okdepac {
					t.Fatalf("New Type Detected %T at Block(%d)/Txn(%d)", v, i, k)
				}
			}
		}
	}
}

// TODO: Find a block with such a Txn
// TestBlockWithTxsAndInvokeTXNV1 tests BlockWithTxs with Invoke V1
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

// TestStateUpdate tests StateUpdateByHash
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
