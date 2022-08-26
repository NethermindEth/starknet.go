package rpc

import (
	"context"
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/dontpanicdao/caigo/types"
)

//go:embed tests/counter.json
var counterFile []byte

// TestAddDeployTransaction tests AddDeployTransaction
func TestAddDeployTransaction(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		BroadcastedDeployTxn    BroadcastedDeployTxn
		ExpectedTransactionHash string
		ExpectedContractAddress string
	}
	var contract types.ContractClass

	if err := json.Unmarshal(counterFile, &contract); err != nil {
		t.Fatal("error loading contract:", err)
	}
	testSet := map[string][]testSetType{
		"mock": {
			{
				BroadcastedDeployTxn: BroadcastedDeployTxn{
					ContractClass: ContractClass{
						Program: "",
						EntryPointsByType: struct {
							CONSTRUCTOR ContractEntryPointList "json:\"CONSTRUCTOR\""
							EXTERNAL    ContractEntryPointList "json:\"EXTERNAL\""
							L1_HANDLER  ContractEntryPointList "json:\"L1_HANDLER\""
						}{},
						Abi: &ContractABI{},
					},
				},
				ExpectedTransactionHash: "0xdeadbeef",
				ExpectedContractAddress: "0xdeadbeef",
			},
		},
		"testnet": {
			{
				BroadcastedDeployTxn: BroadcastedDeployTxn{
					ContractClass: ContractClass{
						Program: "",
						EntryPointsByType: struct {
							CONSTRUCTOR ContractEntryPointList "json:\"CONSTRUCTOR\""
							EXTERNAL    ContractEntryPointList "json:\"EXTERNAL\""
							L1_HANDLER  ContractEntryPointList "json:\"L1_HANDLER\""
						}{},
						Abi: &ContractABI{},
					},
				},
				ExpectedTransactionHash: "0x2149bf99d96ed687a488091ea0d2b1e0b24f73fd7ab96809c2640ae2fc0c791",
				ExpectedContractAddress: "0x30b0fc513edb49b5602f985f5515540a63c8884f3c23a9a9b70f3c14eab7255",
			},
		},
		// TODO: add tests for mainnet when possible or when figure out how to
		// create a white-listed contract. For now, the output is:
		//   - code: NotPermittedContract
		//   - message: The contract class attempted to be deployed is not permitted.
		// This behavior isis on purpose due to the fact mainnet is under limited
		// access. For more details, check this discord
		// [conversation](https://discord.com/channels/793094838509764618/793094838987128844/990692360608444508)
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		output, err := testConfig.client.AddDeployTransaction(context.Background(), test.BroadcastedDeployTxn)
		if err != nil {
			t.Fatal(err)
		}
		if output.TransactionHash != test.ExpectedTransactionHash {
			t.Fatalf("tx expected %s, got %s", test.ExpectedTransactionHash, output.TransactionHash)
		}
		if output.ContractAddress != test.ExpectedContractAddress {
			t.Fatalf("contract expected %s, got %s", test.ExpectedContractAddress, output.ContractAddress)
		}
	}
}

// TestAddDeclareTransaction tests AddDeclareTransaction
func TestAddDeclareTransaction(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		BroadcastedDeclareTxn   BroadcastedDeclareTxn
		ExpectedTransactionHash string
		ExpectedClassHash       string
	}
	var contract types.ContractClass

	if err := json.Unmarshal(counterFile, &contract); err != nil {
		t.Fatal("error loading contract:", err)
	}
	testSet := map[string][]testSetType{
		"mock": {
			{
				BroadcastedDeclareTxn: BroadcastedDeclareTxn{
					BroadcastedCommonTxnProperties: BroadcastedCommonTxnProperties{
						Type:      TxnType(""),
						MaxFee:    "",
						Version:   NumAsHex(""),
						Signature: Signature{},
						Nonce:     "",
					},
					ContractClass: ContractClass{
						Program: "",
						EntryPointsByType: struct {
							CONSTRUCTOR ContractEntryPointList "json:\"CONSTRUCTOR\""
							EXTERNAL    ContractEntryPointList "json:\"EXTERNAL\""
							L1_HANDLER  ContractEntryPointList "json:\"L1_HANDLER\""
						}{},
						Abi: &ContractABI{},
					},
				},
				ExpectedTransactionHash: "0xdeadbeef",
				ExpectedClassHash:       "0xdeadbeef",
			},
		},
		"testnet": {
			{
				BroadcastedDeclareTxn: BroadcastedDeclareTxn{
					BroadcastedCommonTxnProperties: BroadcastedCommonTxnProperties{
						Type:      TxnType(""),
						MaxFee:    "",
						Version:   NumAsHex(""),
						Signature: Signature{},
						Nonce:     "",
					},
					ContractClass: ContractClass{
						Program: "",
						EntryPointsByType: struct {
							CONSTRUCTOR ContractEntryPointList "json:\"CONSTRUCTOR\""
							EXTERNAL    ContractEntryPointList "json:\"EXTERNAL\""
							L1_HANDLER  ContractEntryPointList "json:\"L1_HANDLER\""
						}{},
						Abi: &ContractABI{},
					},
				},
				ExpectedTransactionHash: "0x3d570dbde5ed56ddcb5f69578fb5f83b362c4af8b2a60e2be33ed229148e10a",
				ExpectedClassHash:       "0x646552d8029a8fe940dbbe2847bce558d3d1b3e78a5519e970395df6a2b2cc9",
			},
		},
		// TODO: add tests for mainnet when possible or when figure out how to
		// create a white-listed contract.
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		output, err := testConfig.client.AddDeclareTransaction(context.Background(), test.BroadcastedDeclareTxn)
		if err != nil {
			t.Fatal(err)
		}
		if output.TransactionHash != test.ExpectedTransactionHash {
			t.Fatalf("tx expected %s, got %s", test.ExpectedTransactionHash, output.TransactionHash)
		}
		if output.ClassHash != test.ExpectedClassHash {
			t.Fatalf("class expected %s, got %s", test.ExpectedClassHash, output.ClassHash)
		}
	}
}

// TestAddInvokeTransaction tests AddInvokeTransaction
func TestAddInvokeTransaction(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		BroadcastedInvokeTxn    BroadcastedInvokeTxn
		ExpectedTransactionHash string
	}
	testSet := map[string][]testSetType{
		"mock": {
			{
				BroadcastedInvokeTxn: BroadcastedInvokeTxnV0{
					BroadcastedCommonTxnProperties: BroadcastedCommonTxnProperties{
						Type: TxnType("0xdeadbeef"),
						Signature: []string{
							"3557065757165699682249469970267166698995647077461960906176449260016084767701",
							"3202126414680946801789588986259466145787792017299869598314522555275920413944",
						},
						MaxFee:  "0x4f388496839",
						Version: "0x0",
					},
					InvokeV0: InvokeV0(
						FunctionCall{
							ContractAddress:    Address("0x23371b227eaecd8e8920cd429d2cd0f3fee6abaacca08d3ab82a7cdd"),
							EntryPointSelector: "0x15d40a3d6ca2ac30f4031e42be28da9b056fef9bb7357ac5e85627ee876e5ad",
							CallData: []string{
								"0x1",
								"0x677bb1cdc050e8d63855e8743ab6e09179138def390676cc03c484daf112ba1",
								"0x362398bec32bc0ebb411203221a35a0301193a96f317ebe5e40be9f60d15320",
								"0x0",
								"0x1",
								"0x1",
								"0x2b",
								"0x0",
							},
						},
					),
				},
				ExpectedTransactionHash: "0xdeadbeef",
			},
		},
		"testnet": {},
		"mainnet": {
			{
				BroadcastedInvokeTxn: BroadcastedInvokeTxnV0{
					BroadcastedCommonTxnProperties: BroadcastedCommonTxnProperties{
						Type: TxnType("0xdeadbeef"),
						Signature: []string{
							"3557065757165699682249469970267166698995647077461960906176449260016084767701",
							"3202126414680946801789588986259466145787792017299869598314522555275920413944",
						},
						MaxFee:  "0x4f388496839",
						Version: "0x0",
					},
					InvokeV0: InvokeV0(
						FunctionCall{
							ContractAddress:    Address("0x23371b227eaecd8e8920cd429d2cd0f3fee6abaacca08d3ab82a7cdd"),
							EntryPointSelector: "0x15d40a3d6ca2ac30f4031e42be28da9b056fef9bb7357ac5e85627ee876e5ad",
							CallData: []string{
								"0x1",
								"0x677bb1cdc050e8d63855e8743ab6e09179138def390676cc03c484daf112ba1",
								"0x362398bec32bc0ebb411203221a35a0301193a96f317ebe5e40be9f60d15320",
								"0x0",
								"0x1",
								"0x1",
								"0x2b",
								"0x0",
							},
						},
					),
				},
				ExpectedTransactionHash: "0xdeadbeef",
			},
		},
	}[testEnv]

	for _, test := range testSet {
		output, err := testConfig.client.AddInvokeTransaction(context.Background(), test.BroadcastedInvokeTxn)
		if err != nil || output == nil {
			t.Fatalf("output is nil, go err %v", err)
		}
		if output.TransactionHash != test.ExpectedTransactionHash {
			t.Fatalf("tx expected %s, got %s", test.ExpectedTransactionHash, output.TransactionHash)
		}
	}
}
