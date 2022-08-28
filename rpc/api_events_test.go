package rpc

import (
	"context"
	"testing"
)

// TestEvents tests Events
func TestEvents(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		FromBlock          BlockID
		ExpectedEventCount int
	}
	testSet := map[string][]testSetType{
		"mock":    {},
		"testnet": {},
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		p := EventFilterParams{
			EventFilter{
				FromBlock: test.FromBlock,
				ToBlock:   test.FromBlock,
			},
			ResultPageRequest{
				ChunkSize: 100,
			},
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
