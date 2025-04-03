package abigen

import (
	"strings"
	"testing"

	"github.com/NethermindEth/starknet.go/abigen/accounts/abi"
)

func TestBindCairo(t *testing.T) {
	t.Skip("Skipping test until template embedding is fixed")
	
	abiJSON := `[
		{
			"type": "function",
			"name": "increase_balance",
			"inputs": [
				{
					"name": "amount",
					"type": "core::felt252"
				}
			],
			"outputs": [],
			"state_mutability": "external"
		},
		{
			"type": "function",
			"name": "get_balance",
			"inputs": [],
			"outputs": [
				{
					"type": "core::felt252"
				}
			],
			"state_mutability": "view"
		},
		{
			"type": "event",
			"name": "contracts_v2::hello_starknet::HelloStarknet::Event",
			"kind": "enum",
			"variants": []
		}
	]`

	types := []string{"HelloStarknet"}
	abis := []string{abiJSON}
	bytecodes := []string{""}
	pkg := "test"

	code, err := BindCairo(types, abis, bytecodes, pkg)
	if err != nil {
		t.Fatalf("binding failed: %v", err)
	}

	if !strings.Contains(code, "package test") {
		t.Errorf("generated code does not contain package declaration")
	}

	if !strings.Contains(code, "type HelloStarknet struct") {
		t.Errorf("generated code does not contain contract struct")
	}

	if !strings.Contains(code, "func (_HelloStarknet *HelloStarknetCaller) GetBalance") {
		t.Errorf("generated code does not contain view method")
	}

	if !strings.Contains(code, "func (_HelloStarknet *HelloStarknetTransactor) IncreaseBalance") {
		t.Errorf("generated code does not contain external method")
	}
}

func TestCairoTypeConversion(t *testing.T) {
	tests := []struct {
		cairoType string
		goType    string
	}{
		{"core::felt252", "*felt.Felt"},
		{"core::integer::u8", "uint32"},
		{"core::integer::u16", "uint32"},
		{"core::integer::u32", "uint32"},
		{"core::integer::u64", "uint64"},
		{"core::integer::u128", "uint64"},
		{"core::integer::u256", "*big.Int"},
		{"core::bool", "bool"},
		{"core::starknet::ContractAddress", "*felt.Felt"},
		{"core::array::Array<core::felt252>", "[]*felt.Felt"},
		{"core::array::Array<core::integer::u64>", "[]uint64"},
	}

	for _, tt := range tests {
		got := bindCairoType(tt.cairoType)
		if got != tt.goType {
			t.Errorf("bindCairoType(%q) = %q, want %q", tt.cairoType, got, tt.goType)
		}
	}
}

func TestCairoABIParsing(t *testing.T) {
	abiJSON := `[
		{
			"type": "function",
			"name": "test_types",
			"inputs": [
				{
					"name": "felt_param",
					"type": "core::felt252"
				},
				{
					"name": "u256_param",
					"type": "core::integer::u256"
				},
				{
					"name": "bool_param",
					"type": "core::bool"
				},
				{
					"name": "array_param",
					"type": "core::array::Array<core::felt252>"
				}
			],
			"outputs": [
				{
					"type": "core::felt252"
				}
			],
			"state_mutability": "view"
		}
	]`

	reader := strings.NewReader(abiJSON)
	parsedABI, err := abi.JSON(reader)
	if err != nil {
		t.Fatalf("ABI parsing failed: %v", err)
	}

	method, exists := parsedABI.Methods["test_types"]
	if !exists {
		t.Fatalf("method 'test_types' not found in parsed ABI")
	}

	expectedInputs := []struct {
		name string
		typ  string
	}{
		{"felt_param", "core::felt252"},
		{"u256_param", "core::integer::u256"},
		{"bool_param", "core::bool"},
		{"array_param", "core::array::Array<core::felt252>"},
	}

	if len(method.Inputs) != len(expectedInputs) {
		t.Fatalf("expected %d inputs, got %d", len(expectedInputs), len(method.Inputs))
	}

	for i, expected := range expectedInputs {
		if method.Inputs[i].Name != expected.name {
			t.Errorf("input %d: expected name %q, got %q", i, expected.name, method.Inputs[i].Name)
		}
		if method.Inputs[i].Type != expected.typ {
			t.Errorf("input %d: expected type %q, got %q", i, expected.typ, method.Inputs[i].Type)
		}
	}

	if len(method.Outputs) != 1 {
		t.Fatalf("expected 1 output, got %d", len(method.Outputs))
	}
	if method.Outputs[0].Type != "core::felt252" {
		t.Errorf("expected output type %q, got %q", "core::felt252", method.Outputs[0].Type)
	}

	if method.StateMutability != "view" {
		t.Errorf("expected state mutability %q, got %q", "view", method.StateMutability)
	}
}

func TestCairoBinderCreation(t *testing.T) {
	abiJSON := `[
		{
			"type": "function",
			"name": "view_method",
			"inputs": [],
			"outputs": [{"type": "core::felt252"}],
			"state_mutability": "view"
		},
		{
			"type": "function",
			"name": "external_method",
			"inputs": [{"name": "param", "type": "core::felt252"}],
			"outputs": [],
			"state_mutability": "external"
		},
		{
			"type": "event",
			"name": "test_event",
			"keys": [{"name": "key", "type": "core::felt252"}],
			"data": [{"name": "data", "type": "core::felt252"}]
		}
	]`

	reader := strings.NewReader(abiJSON)
	parsedABI, err := abi.JSON(reader)
	if err != nil {
		t.Fatalf("ABI parsing failed: %v", err)
	}

	binder := newCairoBinder(parsedABI)

	if len(binder.methods) != 2 {
		t.Fatalf("expected 2 methods, got %d", len(binder.methods))
	}

	viewMethod, exists := binder.methods["view_method"]
	if !exists {
		t.Fatalf("method 'view_method' not found in binder")
	}
	if !viewMethod.Const {
		t.Errorf("expected 'view_method' to be constant")
	}
	if viewMethod.Normalized.Name != "ViewMethod" {
		t.Errorf("expected normalized name 'ViewMethod', got %q", viewMethod.Normalized.Name)
	}

	externalMethod, exists := binder.methods["external_method"]
	if !exists {
		t.Fatalf("method 'external_method' not found in binder")
	}
	if externalMethod.Const {
		t.Errorf("expected 'external_method' to be non-constant")
	}
	if externalMethod.Normalized.Name != "ExternalMethod" {
		t.Errorf("expected normalized name 'ExternalMethod', got %q", externalMethod.Normalized.Name)
	}

	if len(binder.events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(binder.events))
	}

	event, exists := binder.events["test_event"]
	if !exists {
		t.Fatalf("event 'test_event' not found in binder")
	}
	if event.Normalized.Name != "TestEvent" {
		t.Errorf("expected normalized name 'TestEvent', got %q", event.Normalized.Name)
	}
}
