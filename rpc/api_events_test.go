package rpc

import (
	"context"
	"testing"
)

// TestEvents with blockId as number
func TestEventsBlockNumber(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		FromBlockNumber    BlockNumber
		ToBlockNumber      BlockNumber
		ExpectedEventCount int
	}
	testSet := map[string][]testSetType{
		"mock": {
			{
				FromBlockNumber:    BlockNumber{BlockNumber: 1000},
				ToBlockNumber:      BlockNumber{BlockNumber: 1001},
				ExpectedEventCount: 1,
			},
		},
		"testnet": {
			{
				FromBlockNumber:    BlockNumber{BlockNumber: 250000},
				ToBlockNumber:      BlockNumber{BlockNumber: 250001},
				ExpectedEventCount: 142,
			},
		},
		"mainnet": {
			{
				FromBlockNumber:    BlockNumber{BlockNumber: 1000},
				ToBlockNumber:      BlockNumber{BlockNumber: 1002},
				ExpectedEventCount: 2,
			},
		},
	}[testEnv]

	for _, test := range testSet {
		p := EventParams{
			FromBlock:  test.FromBlockNumber,
			ToBlock:    test.ToBlockNumber,
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

// TestEvents with blockId as number
func TestEventsBlockHash(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		FromBlockHash      BlockHash
		ToBlockHash        BlockHash
		ExpectedEventCount int
	}
	testSet := map[string][]testSetType{
		"mock": {
			{
				FromBlockHash:      BlockHash{BlockHash: "0x318b1c6c894f7f4b0c930534ab1bca261cb519c3b178a140302bd576d96395"},
				ToBlockHash:        BlockHash{BlockHash: "0x4865b60b759b0b838433610bd79f429be19ab70aa3139716424c50efb8c570a"},
				ExpectedEventCount: 1,
			},
		},
		"testnet": {
			{
				FromBlockHash:      BlockHash{BlockHash: "0x318b1c6c894f7f4b0c930534ab1bca261cb519c3b178a140302bd576d96395"},
				ToBlockHash:        BlockHash{BlockHash: "0x4865b60b759b0b838433610bd79f429be19ab70aa3139716424c50efb8c570a"},
				ExpectedEventCount: 142,
			},
		},
		"mainnet": {
			{
				FromBlockHash:      BlockHash{BlockHash: "0x37cb14332210a0eb0088c914d6516bae855c0012f499cef87f2109566180a8e"},
				ToBlockHash:        BlockHash{BlockHash: "0x1b89a82acd63995178f0375bb7003e0ee6423fe289745e5289f87ecb80bda45"},
				ExpectedEventCount: 2,
			},
		},
	}[testEnv]

	for _, test := range testSet {
		p := EventParams{
			FromBlock:  test.FromBlockHash,
			ToBlock:    test.ToBlockHash,
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

// TestEvents with blockId as tag
func TestEventsBlockTag(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		FromBlockTag string
		ToBlockTag   string
	}
	testSet := map[string][]testSetType{
		"mock": {
			{
				FromBlockTag: "latest",
				ToBlockTag:   "pending",
			},
		},
		"testnet": {
			{
				FromBlockTag: "latest",
				ToBlockTag:   "pending",
			},
		},
		"mainnet": {
			{
				FromBlockTag: "latest",
				ToBlockTag:   "pending",
			},
		},
	}[testEnv]

	for _, test := range testSet {
		p := EventParams{
			FromBlock:  test.FromBlockTag,
			ToBlock:    test.ToBlockTag,
			PageSize:   1000,
			PageNumber: 0,
		}
		events, err := testConfig.client.Events(context.Background(), p)
		if err != nil && events == nil {
			t.Fatal(err)
		}
	}
}
