package methods

import (
	"testing"

	"github.com/NethermindEth/starknet.go/rpc/callers"
	. "github.com/NethermindEth/starknet.go/rpc/rpcv10"
	"github.com/stretchr/testify/require"
)

// GetCommonBlockIDs returns a list of common block IDs to use in some RPC tests.
// It includes all block tags, a range of block numbers and the latest block hash.
func GetCommonBlockIDs(t *testing.T, caller callers.Caller) []BlockID {
	t.Helper()

	// *** all valid block tags ***
	commonBlockIDs := []BlockID{
		WithBlockTag(BlockTagLatest),
		WithBlockTag(BlockTagPreConfirmed),
		WithBlockTag(BlockTagL1Accepted),
	}

	// *** getting the common block number range ***

	// 5 blocks from the first 1M blocks of the network
	// (a lot of changes in the first blocks)
	commonBlockIDs = append(commonBlockIDs, []BlockID{
		WithBlockNumber(0),
		WithBlockNumber(200_000),
		WithBlockNumber(400_000),
		WithBlockNumber(600_000),
		WithBlockNumber(800_000),
		WithBlockNumber(1_000_000),
	}...)

	// get the latest block number of the network
	blockHashAndNumber, err := BlockHashAndNumber(t.Context(), caller)
	require.NoError(t, err, "failed to get the block number")

	// after the block 1_000_000, we add one block every 500_000 blocks
	// until the latest block
	for i := uint64(1_500_000); i < blockHashAndNumber.Number; i += 500_000 {
		commonBlockIDs = append(commonBlockIDs, WithBlockNumber(i))
	}

	// add the latest block hash
	commonBlockIDs = append(commonBlockIDs, WithBlockHash(blockHashAndNumber.Hash))

	return commonBlockIDs
}
