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
		BlockID                   BlockID
		ExpectedError             error
		ExpectedBlockWithTxHashes *Block
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
		"mock": {},
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

	var fullBlockGoerli300114 = Block{
		BlockHeader: BlockHeader{
			BlockHash:        utils.TestHexToFelt(t, "0x184268bfbce24766fa53b65c9c8b30b295e145e8281d543a015b46308e27fdf"),
			ParentHash:       utils.TestHexToFelt(t, "0x7307cb0d7fa65c111e71cdfb6209bdc90d2454d4c0f34d8bf5a3fe477826c3c"),
			SequencerAddress: utils.TestHexToFelt(t, "0x46a89ae102987331d369645031b49c27738ed096f2789c24449966da4c6de6b"),
			BlockNumber:      300114,
			NewRoot:          utils.TestHexToFelt(t, "0x239a44410e78665f41f7a65ef3b5ed244ce411965498a83f80f904e22df1045"),
			Timestamp:        1660701246,
		},
		Status: "ACCEPTED_ON_L1",
		Transactions: []Transaction{
			DeclareTxn{
				CommonTransaction: CommonTransaction{
					TransactionHash: utils.TestHexToFelt(t, "0x46a9f52a96b2d226407929e04cb02507e531f7c78b9196fc8c910351d8c33f3"),
					BroadcastedTxnCommonProperties: BroadcastedTxnCommonProperties{
						Type:      TransactionType_Declare,
						MaxFee:    &felt.Zero,
						Version:   TransactionV0,
						Signature: []*felt.Felt{},
						Nonce:     &felt.Zero,
					},
				},
				ClassHash:     utils.TestHexToFelt(t, "0x6feb117d1c3032b0ae7bd3b50cd8ec4a78c621dca0d63ddc17890b78a6c3b49"),
				SenderAddress: utils.TestHexToFelt(t, "0x1"),
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
				BlockID:                     WithBlockNumber(300114),
				ExpectedError:               nil,
				LookupTxnPositionInOriginal: 3,
				LookupTxnPositionInExpected: 0,
				ExpectedBlockWithTxs:        &fullBlockGoerli300114,
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
				StartBlock: 375919,
				EndBlock:   376000,
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
				_, okdec := v.(DeclareTxn)
				_, okdep := v.(DeployTxn)
				_, okdepac := v.(DeployAccountTxn)
				if !okv0 && !okv1 && !okl1 && !okdec && !okdep && !okdepac {
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
