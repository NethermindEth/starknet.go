package rpc

import (
	"context"
	"testing"
)

// TestEvents tests Events
func TestEvents(t *testing.T) {
	testConfig := beforeEach(t)
	defer testConfig.client.Close()

	type testSetType struct {
		FromBlockNumber    uint64
		ExpectedEventCount int
	}
	testSet := map[string][]testSetType{
		"mock": {
			{
				FromBlockNumber:    1000,
				ExpectedEventCount: 1,
			},
		},
		"testnet": {
			{
				FromBlockNumber:    250000,
				ExpectedEventCount: 142,
			},
		},
		"mainnet": {
			{
				FromBlockNumber:    1000,
				ExpectedEventCount: 1,
			},
		},
	}[testEnv]

	for _, test := range testSet {
		p := EventParams{
			FromBlock:  test.FromBlockNumber,
			ToBlock:    test.FromBlockNumber,
			PageSize:   1000,
			PageNumber: 0,
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
