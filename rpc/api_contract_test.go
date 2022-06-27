package rpc

import (
	"context"
	"testing"
)

// TestCodeAt tests code for a contract instance. This will be deprecated.
func TestCodeAt(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		ContractHash      string
		ExpectedBytecode0 string
	}
	testSet := map[string][]testSetType{
		"mock": {
			{
				ContractHash:      "0xdeadbeef",
				ExpectedBytecode0: "0xdeadbeef",
			},
		},
		"testnet": {
			{
				ContractHash:      "0x6fbd460228d843b7fbef670ff15607bf72e19fa94de21e29811ada167b4ca39",
				ExpectedBytecode0: "0x480680017fff8000",
			},
		},
		"mainnet": {
			{
				ContractHash:      "0x3b4be7def2fc08589348966255e101824928659ebb724855223ff3a8c831efa",
				ExpectedBytecode0: "0x40780017fff7fff",
			},
		},
	}[testEnv]

	for _, test := range testSet {
		code, err := testConfig.client.CodeAt(context.Background(), test.ContractHash)
		if err != nil {
			t.Fatal(err)
		}
		if code == nil || len(code.Bytecode) == 0 {
			t.Fatal("code should exist")
		}
		if code.Bytecode[0] != test.ExpectedBytecode0 {
			t.Fatalf("expecting bytecode[0] %s, got %s", test.ExpectedBytecode0, code.Bytecode[0])
		}
	}
}

// TestClassAt tests code for a class.
func TestClassAt(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		ContractHash      string
		ExpectedOperation string
	}
	testSet := map[string][]testSetType{
		"mock": {
			{
				ContractHash:      "0xdeadbeef",
				ExpectedOperation: "0xdeadbeef",
			},
		},
		"testnet": {
			{
				ContractHash:      "0x6fbd460228d843b7fbef670ff15607bf72e19fa94de21e29811ada167b4ca39",
				ExpectedOperation: "0x480680017fff8000",
			},
		},
		"mainnet": {
			{
				ContractHash:      "0x028105caf03e1c4eb96b1c18d39d9f03bd53e5d2affd0874792e5bf05f3e529f",
				ExpectedOperation: "0x20780017fff7ffd",
			},
		},
	}[testEnv]

	for _, test := range testSet {
		class, err := testConfig.client.ClassAt(context.Background(), test.ContractHash)
		if err != nil {
			t.Fatal(err)
		}
		if class == nil || class.Program == nil {
			t.Fatal("code should exist")
		}
		ops, ok := class.Program.([]interface{})
		if !ok {
			t.Fatalf("program should return []interface{}, instead %T", class.Program)
		}
		if len(ops) == 0 {
			t.Fatal("program have several operations")
		}
		op, _ := ops[0].(string)
		if op != test.ExpectedOperation {
			t.Fatalf("op expected %s, got %s", test.ExpectedOperation, op)
		}
	}
}

// TestClassHashAt tests code for a getClassHashAt.
func TestClassHashAt(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		ContractHash      string
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
		classhash, err := testConfig.client.ClassHashAt(context.Background(), test.ContractHash)
		if err != nil {
			t.Fatal(err)
		}
		if classhash != test.ExpectedClassHash {
			t.Fatalf("class expect %s, got %s", test.ExpectedClassHash, classhash)
		}
	}
}

// TestClass tests code for a class.
func TestClass(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		ClassHash         string
		ExpectedOperation string
	}
	testSet := map[string][]testSetType{
		"mock": {
			{
				ClassHash:         "0xdeadbeef",
				ExpectedOperation: "0xdeadbeef",
			},
		},
		"testnet": {
			{
				ClassHash:         "0x493af3546940eb96471cf95ae3a5aa1286217b07edd1e12d00143010ca904b1",
				ExpectedOperation: "0x40780017fff7fff",
			},
		},
		"mainnet": {
			{
				ClassHash:         "0x4c53698c9a42341e4123632e87b752d6ae470ddedeb8b0063eaa2deea387eeb",
				ExpectedOperation: "0x40780017fff7fff",
			},
		},
	}[testEnv]

	for _, test := range testSet {
		class, err := testConfig.client.Class(context.Background(), test.ClassHash)
		if err != nil {
			t.Fatal(err)
		}
		if class == nil || class.Program == nil {
			t.Fatal("code should exist")
		}
		ops, ok := class.Program.([]interface{})
		if !ok {
			t.Fatalf("program should return []interface{}, instead %T", class.Program)
		}
		if len(ops) == 0 {
			t.Fatal("program have several operations")
		}
		op, _ := ops[0].(string)
		if op != test.ExpectedOperation {
			t.Fatalf("op expected %s, got %s", test.ExpectedOperation, op)
		}
	}
}
