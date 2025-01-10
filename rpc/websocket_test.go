package rpc

import (
	"context"
	"fmt"
	"testing"

	"github.com/NethermindEth/starknet.go/client"
	"github.com/stretchr/testify/require"
)

func TestSubscribeNewHeads(t *testing.T) {
	if testEnv != "testnet" {
		t.Skip("Skipping test as it requires a testnet environment")
	}

	testConfig := beforeEach(t)
	require.NotNil(t, testConfig.wsBase, "wsProvider base is not set")

	type testSetType struct {
		headers         chan *BlockHeader
		blockID         []BlockID
		counter         int
		isErrorExpected bool
	}

	provider := testConfig.provider
	blockNumber, err := provider.BlockNumber(context.Background())
	require.NoError(t, err)

	latestBlockNumbers := []uint64{blockNumber, blockNumber + 1} // for the case the latest block number is updated

	testSet := map[string][]testSetType{
		"testnet": {
			{ // normal
				headers:         make(chan *BlockHeader),
				isErrorExpected: false,
			},
			{ // with tag latest
				headers:         make(chan *BlockHeader),
				blockID:         []BlockID{WithBlockTag("latest")},
				isErrorExpected: false,
			},
			{ // with tag pending
				headers:         make(chan *BlockHeader),
				blockID:         []BlockID{WithBlockTag("pending")},
				isErrorExpected: true,
			},
			{ // with block number within the range of 1024 blocks
				headers:         make(chan *BlockHeader),
				blockID:         []BlockID{WithBlockNumber(blockNumber - 100)},
				counter:         100,
				isErrorExpected: false,
			},
			{ // invalid, with block number out of the range of 1024 blocks
				headers:         make(chan *BlockHeader),
				blockID:         []BlockID{WithBlockNumber(blockNumber - 1025)},
				isErrorExpected: true,
			},
			{ // invalid, more than one blockID parameter
				headers:         make(chan *BlockHeader),
				blockID:         []BlockID{WithBlockTag("latest"), WithBlockTag("latest")},
				isErrorExpected: true,
			},
		},
	}[testEnv]

	for index, test := range testSet {
		t.Run(fmt.Sprintf("test %d", index+1), func(t *testing.T) {

			wsProvider, err := NewWebsocketProvider(testConfig.wsBase)
			require.NoError(t, err)
			defer wsProvider.Close()

			var sub *client.ClientSubscription
			if len(test.blockID) == 0 {
				sub, err = wsProvider.SubscribeNewHeads(context.Background(), test.headers)
			} else {
				sub, err = wsProvider.SubscribeNewHeads(context.Background(), test.headers, test.blockID...)
			}

			if test.isErrorExpected {
				require.Error(t, err)
				return
			} else {
				require.NoError(t, err)
			}

			require.NotNil(t, sub)
			defer sub.Unsubscribe()

			for {
				select {
				case resp := <-test.headers:
					require.IsType(t, &BlockHeader{}, resp)

					if test.counter != 0 {
						if test.counter == 1 {
							require.Contains(t, latestBlockNumbers, resp.BlockNumber+1)
							return
						} else {
							test.counter--
						}
					} else {
						return
					}
				case err := <-sub.Err():
					require.NoError(t, err)
				}
			}
		})
	}
}
