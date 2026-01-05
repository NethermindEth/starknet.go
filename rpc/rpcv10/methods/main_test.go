package methods

// @todo remove this file later

import (
	"os"
	"testing"

	"github.com/NethermindEth/starknet.go/internal/tests"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/rpc/rpcv10"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	tests.LoadEnv()

	os.Exit(m.Run())
}

// GetCommonBlockIDs returns a list of common block IDs to use in some RPC tests.
// It includes all block tags, a range of block numbers and the latest block hash.
func GetCommonBlockIDs(t *testing.T, caller rpc.Caller) []rpcv10.BlockID {
	t.Helper()

	// *** all valid block tags ***
	commonBlockIDs := []rpcv10.BlockID{
		rpcv10.WithBlockTag(rpcv10.BlockTagLatest),
		rpcv10.WithBlockTag(rpcv10.BlockTagPreConfirmed),
		rpcv10.WithBlockTag(rpcv10.BlockTagL1Accepted),
	}

	// *** getting the common block number range ***

	// 5 blocks from the first 1M blocks of the network
	// (a lot of changes in the first blocks)
	commonBlockIDs = append(commonBlockIDs, []rpcv10.BlockID{
		rpcv10.WithBlockNumber(0),
		rpcv10.WithBlockNumber(200_000),
		rpcv10.WithBlockNumber(400_000),
		rpcv10.WithBlockNumber(600_000),
		rpcv10.WithBlockNumber(800_000),
		rpcv10.WithBlockNumber(1_000_000),
	}...)

	// get the latest block number of the network
	blockHashAndNumber, err := BlockHashAndNumber(t.Context(), caller)
	require.NoError(t, err, "failed to get the block number")

	// after the block 1_000_000, we add one block every 500_000 blocks
	// until the latest block
	for i := uint64(1_500_000); i < blockHashAndNumber.Number; i += 500_000 {
		commonBlockIDs = append(commonBlockIDs, rpcv10.WithBlockNumber(i))
	}

	// add the latest block hash
	commonBlockIDs = append(commonBlockIDs, rpcv10.WithBlockHash(blockHashAndNumber.Hash))

	return commonBlockIDs
}
