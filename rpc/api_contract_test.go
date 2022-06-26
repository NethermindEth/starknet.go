package rpc

import (
	"context"
	"testing"
)

// func TestContract(t *testing.T) {

// 	// tested against pathfinder @ 0313b14ea1fad8f73635a3002d106908813e57f1
// 	classHash, err := client.ClassHashAt(context.Background(), accountAddr)
// 	if err != nil {
// 		t.Errorf("Could not retrieve class hash: %v\n", err)
// 	}

// 	_, err = client.ClassAt(context.Background(), accountAddr)
// 	if err != nil {
// 		t.Errorf("Could not retrieve class: %v\n", err)
// 	}

// 	_, err = client.Class(context.Background(), classHash)
// 	if err != nil {
// 		t.Errorf("Could not retrieve class: %v\n", err)
// 	}
// }

// TestCodeAt tests code for a contract instance. This will be deprecated.
func TestCodeAt(t *testing.T) {
	testConfig := beforeEach(t)
	defer testConfig.client.Close()

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
