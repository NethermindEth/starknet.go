package rpc

import (
	"context"
	"testing"
)

// TestEvents tests Events
func TestEvents(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		eventFilter        EventFilter
		ExpectedEventCount int
	}
	testSet := map[string][]testSetType{
		"mock":    {},
		"testnet": {
			// TODO: add the test back. Right now, response time is very bad with PathFinder
			// so it is working but too slow to sustain
			// 	{
			// 	eventFilter: types.EventFilter{
			// 		FromBlock: types.BlockID{
			// 			Hash: &blockHash,
			// 		},
			// 		ToBlock: types.BlockID{
			// 			Hash: &blockHash,
			// 		},
			// 		Address:    types.StrToFelt("0x49d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"),
			// 		Keys:       []string{"0x99cd8bde557814842a3121e8ddfd433a539b8c9f14bf31ebf108d12e6196e9"},
			// 		PageSize:   1000,
			// 		PageNumber: 0,
			// 	},
			// 	ExpectedEventCount: 28,
			// },
		},
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		spy := NewSpy(testConfig.provider.c)
		testConfig.provider.c = spy
		eventInput := EventsInput{
			EventFilter: test.eventFilter,
		}
		events, err := testConfig.provider.Events(context.Background(), eventInput)
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
