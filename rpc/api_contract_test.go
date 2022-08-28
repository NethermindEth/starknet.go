package rpc

import (
	"context"
	"testing"
)

// TestClassAt tests code for a class.
func TestClassAt(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		ContractAddress   Address
		ExpectedOperation string
	}
	testSet := map[string][]testSetType{
		"mock": {
			{
				ContractAddress:   Address("0xdeadbeef"),
				ExpectedOperation: "0xdeadbeef",
			},
		},
		"testnet": {
			{
				ContractAddress:   Address("0x6fbd460228d843b7fbef670ff15607bf72e19fa94de21e29811ada167b4ca39"),
				ExpectedOperation: "0x480680017fff8000",
			},
		},
		"mainnet": {
			{
				ContractAddress:   Address("0x028105caf03e1c4eb96b1c18d39d9f03bd53e5d2affd0874792e5bf05f3e529f"),
				ExpectedOperation: "0x20780017fff7ffd",
			},
		},
	}[testEnv]

	for _, test := range testSet {
		spy := NewSpy(testConfig.client.c)
		testConfig.client.c = spy
		class, err := testConfig.client.ClassAt(context.Background(), WithBlockTag("latest"), test.ContractAddress)
		if err != nil {
			t.Fatal(err)
		}
		diff, err := spy.Compare(class, false)
		if err != nil {
			t.Fatal("expecting to match", err)
		}
		if diff != "FullMatch" {
			spy.Compare(class, true)
			t.Fatal("structure expecting to be FullMatch, instead", diff)
		}
		if class == nil || class.Program == "" {
			t.Fatal("code should exist")
		}
	}
}

// TestClassHashAt tests code for a ClassHashAt.
func TestClassHashAt(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		ContractHash      Address
		ExpectedClassHash string
	}
	testSet := map[string][]testSetType{
		"mock": {
			{
				ContractHash:      "0xdeadbeef",
				ExpectedClassHash: "0xdeadbeef",
			},
		},
		"testnet": {
			{
				ContractHash:      "0x315e364b162653e5c7b23efd34f8da27ba9c069b68e3042b7d76ce1df890313",
				ExpectedClassHash: "0x493af3546940eb96471cf95ae3a5aa1286217b07edd1e12d00143010ca904b1",
			},
		},
		"mainnet": {
			{
				ContractHash:      "0x3b4be7def2fc08589348966255e101824928659ebb724855223ff3a8c831efa",
				ExpectedClassHash: "0x4c53698c9a42341e4123632e87b752d6ae470ddedeb8b0063eaa2deea387eeb",
			},
		},
	}[testEnv]

	for _, test := range testSet {

		spy := NewSpy(testConfig.client.c)
		testConfig.client.c = spy
		classhash, err := testConfig.client.ClassHashAt(context.Background(), WithBlockTag("latest"), test.ContractHash)
		if err != nil {
			t.Fatal(err)
		}
		diff, err := spy.Compare(classhash, false)
		if err != nil {
			t.Fatal("expecting to match", err)
		}
		if diff != "FullMatch" {
			spy.Compare(classhash, true)
			t.Fatal("structure expecting to be FullMatch, instead", diff)
		}
		if classhash == nil {
			t.Fatalf("should return a class, instead %v", classhash)
		}
		if *classhash != test.ExpectedClassHash {
			t.Fatalf("class expect %s, got %s", test.ExpectedClassHash, *classhash)
		}
	}
}

// TestClass tests code for a class.
func TestClass(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		BlockID           BlockID
		ClassHash         string
		ExpectedOperation string
	}
	testSet := map[string][]testSetType{
		"mock": {
			{
				BlockID:           WithBlockTag("pending"),
				ClassHash:         "0xdeadbeef",
				ExpectedOperation: "0xdeadbeef",
			},
		},
		"testnet": {},
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		class, err := testConfig.client.Class(context.Background(), test.BlockID, test.ClassHash)
		if err != nil {
			t.Fatal(err)
		}
		if class == nil || class.Program == "" {
			t.Fatal("code should exist")
		}
	}
}

// TestStorageAt tests StorageAt
func TestStorageAt(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		ContractHash   Address
		StorageKey     string
		BlockHashOrTag BlockID
		ExpectedValue  string
	}
	testSet := map[string][]testSetType{
		"mock": {
			{
				ContractHash:   Address("0xdeadbeef"),
				StorageKey:     "_signer",
				BlockHashOrTag: WithBlockTag("latest"),
				ExpectedValue:  "0xdeadbeef",
			},
		},
		"testnet": {
			{
				ContractHash:   Address("0x6fbd460228d843b7fbef670ff15607bf72e19fa94de21e29811ada167b4ca39"),
				StorageKey:     "balance",
				BlockHashOrTag: WithBlockTag("latest"),
				ExpectedValue:  "0x1e240",
			},
		},
		"mainnet": {
			{
				ContractHash:   Address("0x8d17e6a3B92a2b5Fa21B8e7B5a3A794B05e06C5FD6C6451C6F2695Ba77101"),
				StorageKey:     "_signer",
				BlockHashOrTag: WithBlockTag("latest"),
				ExpectedValue:  "0x7f72660ca40b8ca85f9c0dd38db773f17da7a52f5fc0521cb8b8d8d44e224b8",
			},
		},
	}[testEnv]

	for _, test := range testSet {
		value, err := testConfig.client.StorageAt(context.Background(), test.ContractHash, test.StorageKey, test.BlockHashOrTag)
		if err != nil {
			t.Fatal(err)
		}
		if value != test.ExpectedValue {
			t.Fatalf("expecting value %s, got %s", test.ExpectedValue, value)
		}
	}
}

// TestNonce tests Nonce
func TestNonce(t *testing.T) {
	// Pathfinder 0.3.1 is failing.
	// TODO: get feedback to add the test
	// curl -s -X POST -H 'Content-Type: application/json' \
	// -d '{"jsonrpc":"2.0","id":"0","method":"starknet_getNonce","params":["latest", "0x0207aCC15dc241e7d167E67e30E769719A727d3E0fa47f9E187707289885Dfde"]}' \
	// https://localhost:9545/

	testEnv = "testnet"
	testConfig := beforeEach(t)

	type testSetType struct {
		ContractHash  Address
		Block         BlockID
		ExpectedNonce string
	}
	testSet := map[string][]testSetType{
		"mock":    {},
		"testnet": {},
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		spy := NewSpy(testConfig.client.c)
		testConfig.client.c = spy
		value, err := testConfig.client.Nonce(context.Background(), test.Block, test.ContractHash)
		if err != nil {
			t.Fatal(err)
		}
		diff, err := spy.Compare(value, false)
		if err != nil {
			t.Fatal("expecting to match", err)
		}
		if diff != "FullMatch" {
			spy.Compare(value, true)
			t.Fatal("structure expecting to be FullMatch, instead", diff)
		}
		if *value != test.ExpectedNonce {
			t.Fatalf("expecting value %s, got %s", test.ExpectedNonce, *value)
		}
	}
}
