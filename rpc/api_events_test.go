package rpc

import (
	"context"
	"testing"

	"github.com/dontpanicdao/caigo/rpc/types"
)

// TestEvents tests Events
func TestEvents(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		eventFilter        types.EventFilter
		ExpectedEventCount int
	}
	blockHash := types.HexToHash("0x39a84ddd8f9c15667661fb99b9e6d841d47651560baed288ca7dbeb501c687b")
	testSet := map[string][]testSetType{
		"mock": {},
		"testnet": {{
			eventFilter: types.EventFilter{
				FromBlock: types.BlockID{
					Hash: &blockHash,
				},
				ToBlock: types.BlockID{
					Hash: &blockHash,
				},
				Address:    types.HexToHash("0x49d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"),
				Keys:       []string{"0x99cd8bde557814842a3121e8ddfd433a539b8c9f14bf31ebf108d12e6196e9"},
				PageSize:   1000,
				PageNumber: 0,
			},
			ExpectedEventCount: 28,
		},
		},
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		spy := NewSpy(testConfig.client.c)
		testConfig.client.c = spy
		events, err := testConfig.client.Events(context.Background(), test.eventFilter)
		if err != nil {
			t.Fatal(err)
		}
		if diff, err := spy.Compare(events, false); err != nil || diff != "FullMatch" {
			spy.Compare(events, true)
			t.Fatal("expecting to match", err)
		}
		if events == nil || len(events.Events) == 0 {
			t.Fatal("events should exist")
		}
		if len(events.Events) != test.ExpectedEventCount {
			t.Fatalf("# events expected %d, got %d", test.ExpectedEventCount, len(events.Events))
		}
	}
}
