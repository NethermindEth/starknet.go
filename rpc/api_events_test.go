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
		FromBlock          types.BlockID
		ExpectedEventCount int
	}
	testSet := map[string][]testSetType{
		"mock":    {},
		"testnet": {},
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		p := types.EventFilter{
			FromBlock: test.FromBlock,
			ToBlock:   test.FromBlock,
			ChunkSize: 100,
		}
		events, err := testConfig.client.Events(context.Background(), p)
		if err != nil {
			t.Fatal(err)
		}
		if events == nil || len(events.Events) == 0 {
			t.Fatal("events should exist")
		}
		if len(events.Events) != test.ExpectedEventCount {
			t.Fatalf("# events expected %d, got %d", test.ExpectedEventCount, len(events.Events))
		}
	}
}
