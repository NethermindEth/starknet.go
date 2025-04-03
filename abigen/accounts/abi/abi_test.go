package abi

import (
	"strings"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/utils"
)

func TestCairoABIParsing(t *testing.T) {
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

	reader := strings.NewReader(abiJSON)
	parsedABI, err := JSON(reader)
	if err != nil {
		t.Fatalf("ABI parsing failed: %v", err)
	}

	if len(parsedABI.Methods) != 2 {
		t.Fatalf("expected 2 methods, got %d", len(parsedABI.Methods))
	}

	increaseBalance, exists := parsedABI.Methods["increase_balance"]
	if !exists {
		t.Fatalf("method 'increase_balance' not found in parsed ABI")
	}
	if increaseBalance.StateMutability != "external" {
		t.Errorf("expected state mutability 'external', got '%s'", increaseBalance.StateMutability)
	}
	if len(increaseBalance.Inputs) != 1 {
		t.Errorf("expected 1 input, got %d", len(increaseBalance.Inputs))
	}
	if increaseBalance.Inputs[0].Name != "amount" {
		t.Errorf("expected input name 'amount', got '%s'", increaseBalance.Inputs[0].Name)
	}
	if increaseBalance.Inputs[0].Type != "core::felt252" {
		t.Errorf("expected input type 'core::felt252', got '%s'", increaseBalance.Inputs[0].Type)
	}

	getBalance, exists := parsedABI.Methods["get_balance"]
	if !exists {
		t.Fatalf("method 'get_balance' not found in parsed ABI")
	}
	if getBalance.StateMutability != "view" {
		t.Errorf("expected state mutability 'view', got '%s'", getBalance.StateMutability)
	}
	if len(getBalance.Inputs) != 0 {
		t.Errorf("expected 0 inputs, got %d", len(getBalance.Inputs))
	}
	if len(getBalance.Outputs) != 1 {
		t.Errorf("expected 1 output, got %d", len(getBalance.Outputs))
	}
	if getBalance.Outputs[0].Type != "core::felt252" {
		t.Errorf("expected output type 'core::felt252', got '%s'", getBalance.Outputs[0].Type)
	}

	if len(parsedABI.Events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(parsedABI.Events))
	}
	eventName := "contracts_v2::hello_starknet::HelloStarknet::Event"
	_, exists = parsedABI.Events[eventName]
	if !exists {
		t.Fatalf("event '%s' not found in parsed ABI", eventName)
	}
}

func TestCairoABIWithComplexTypes(t *testing.T) {
	abiJSON := `[
		{
			"type": "function",
			"name": "test_complex_types",
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
		},
		{
			"type": "struct",
			"name": "MyStruct",
			"members": [
				{
					"name": "field1",
					"type": "core::felt252"
				},
				{
					"name": "field2",
					"type": "core::integer::u256"
				}
			]
		}
	]`

	reader := strings.NewReader(abiJSON)
	parsedABI, err := JSON(reader)
	if err != nil {
		t.Fatalf("ABI parsing failed: %v", err)
	}

	method, exists := parsedABI.Methods["test_complex_types"]
	if !exists {
		t.Fatalf("method 'test_complex_types' not found in parsed ABI")
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

	if len(parsedABI.Structs) != 1 {
		t.Fatalf("expected 1 struct, got %d", len(parsedABI.Structs))
	}

	structType, exists := parsedABI.Structs["MyStruct"]
	if !exists {
		t.Fatalf("struct 'MyStruct' not found in parsed ABI")
	}

	if len(structType.Members) != 2 {
		t.Fatalf("expected 2 struct members, got %d", len(structType.Members))
	}

	if structType.Members[0].Name != "field1" {
		t.Errorf("expected struct member name 'field1', got '%s'", structType.Members[0].Name)
	}
	if structType.Members[0].Type != "core::felt252" {
		t.Errorf("expected struct member type 'core::felt252', got '%s'", structType.Members[0].Type)
	}

	if structType.Members[1].Name != "field2" {
		t.Errorf("expected struct member name 'field2', got '%s'", structType.Members[1].Name)
	}
	if structType.Members[1].Type != "core::integer::u256" {
		t.Errorf("expected struct member type 'core::integer::u256', got '%s'", structType.Members[1].Type)
	}
}

func TestPackArguments(t *testing.T) {
	tests := []struct {
		name     string
		args     []Argument
		values   []interface{}
		expected []*felt.Felt
		wantErr  bool
	}{
		{
			name: "felt252",
			args: []Argument{
				{Name: "param", Type: "core::felt252"},
			},
			values: []interface{}{
				utils.Uint64ToFelt(123),
			},
			expected: []*felt.Felt{
				utils.Uint64ToFelt(123),
			},
		},
		{
			name: "bool",
			args: []Argument{
				{Name: "param", Type: "core::bool"},
			},
			values: []interface{}{
				true,
			},
			expected: []*felt.Felt{
				utils.Uint64ToFelt(1),
			},
		},
		{
			name: "u32",
			args: []Argument{
				{Name: "param", Type: "core::integer::u32"},
			},
			values: []interface{}{
				uint32(42),
			},
			expected: []*felt.Felt{
				utils.Uint64ToFelt(42),
			},
		},
		{
			name: "multiple arguments",
			args: []Argument{
				{Name: "param1", Type: "core::felt252"},
				{Name: "param2", Type: "core::bool"},
				{Name: "param3", Type: "core::integer::u32"},
			},
			values: []interface{}{
				utils.Uint64ToFelt(123),
				true,
				uint32(42),
			},
			expected: []*felt.Felt{
				utils.Uint64ToFelt(123),
				utils.Uint64ToFelt(1),
				utils.Uint64ToFelt(42),
			},
		},
		{
			name: "argument count mismatch",
			args: []Argument{
				{Name: "param1", Type: "core::felt252"},
				{Name: "param2", Type: "core::bool"},
			},
			values: []interface{}{
				utils.Uint64ToFelt(123),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := PackArguments(tt.args, tt.values)
			if (err != nil) != tt.wantErr {
				t.Errorf("PackArguments() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if len(result) != len(tt.expected) {
				t.Errorf("PackArguments() result length = %d, want %d", len(result), len(tt.expected))
				return
			}
			for i, v := range result {
				if v.String() != tt.expected[i].String() {
					t.Errorf("PackArguments() result[%d] = %s, want %s", i, v.String(), tt.expected[i].String())
				}
			}
		})
	}
}

func TestUnpackValues(t *testing.T) {
	tests := []struct {
		name     string
		args     []Argument
		data     []*felt.Felt
		expected []interface{}
		wantErr  bool
	}{
		{
			name: "felt252",
			args: []Argument{
				{Type: "core::felt252"},
			},
			data: []*felt.Felt{
				utils.Uint64ToFelt(123),
			},
			expected: []interface{}{
				utils.Uint64ToFelt(123),
			},
		},
		{
			name: "bool",
			args: []Argument{
				{Type: "core::bool"},
			},
			data: []*felt.Felt{
				utils.Uint64ToFelt(1),
			},
			expected: []interface{}{
				true,
			},
		},
		{
			name: "u32",
			args: []Argument{
				{Type: "core::integer::u32"},
			},
			data: []*felt.Felt{
				utils.Uint64ToFelt(42),
			},
			expected: []interface{}{
				uint32(42),
			},
		},
		{
			name: "multiple values",
			args: []Argument{
				{Type: "core::felt252"},
				{Type: "core::bool"},
				{Type: "core::integer::u32"},
			},
			data: []*felt.Felt{
				utils.Uint64ToFelt(123),
				utils.Uint64ToFelt(1),
				utils.Uint64ToFelt(42),
			},
			expected: []interface{}{
				utils.Uint64ToFelt(123),
				true,
				uint32(42),
			},
		},
		{
			name: "insufficient data",
			args: []Argument{
				{Type: "core::felt252"},
				{Type: "core::bool"},
			},
			data: []*felt.Felt{
				utils.Uint64ToFelt(123),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := UnpackValues(tt.args, tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnpackValues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if len(result) != len(tt.expected) {
				t.Errorf("UnpackValues() result length = %d, want %d", len(result), len(tt.expected))
				return
			}

			for i, v := range result {
				switch expected := tt.expected[i].(type) {
				case *felt.Felt:
					got, ok := v.(*felt.Felt)
					if !ok {
						t.Errorf("UnpackValues() result[%d] type = %T, want *felt.Felt", i, v)
						continue
					}
					if got.String() != expected.String() {
						t.Errorf("UnpackValues() result[%d] = %s, want %s", i, got.String(), expected.String())
					}
				case bool:
					got, ok := v.(bool)
					if !ok {
						t.Errorf("UnpackValues() result[%d] type = %T, want bool", i, v)
						continue
					}
					if got != expected {
						t.Errorf("UnpackValues() result[%d] = %v, want %v", i, got, expected)
					}
				case uint32:
					got, ok := v.(uint32)
					if !ok {
						t.Errorf("UnpackValues() result[%d] type = %T, want uint32", i, v)
						continue
					}
					if got != expected {
						t.Errorf("UnpackValues() result[%d] = %v, want %v", i, got, expected)
					}
				default:
					t.Errorf("Unexpected type in test case: %T", expected)
				}
			}
		})
	}
}

func TestGetSelector(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{
			name:     "increase_balance",
			expected: "0x362398bec32bc0ebb411203221a35a0301193a96f317ebe5e40be9f60d15320",
		},
		{
			name:     "get_balance",
			expected: "0x39e11d48192e4333233c7eb19d10ad67c362bb28580c604d67884c85da39695",
		},
		{
			name:     "transfer",
			expected: "0x83afd3f4caedc6eebf44246fe54e38c95e3179a5ec9ea81740eca5b482d12e",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			selector := GetSelector(tt.name)
			expectedFelt, _ := utils.HexToFelt(tt.expected)
			if selector.String() != expectedFelt.String() {
				t.Errorf("GetSelector(%s) = %s, want %s", tt.name, selector.String(), expectedFelt.String())
			}
		})
	}
}
