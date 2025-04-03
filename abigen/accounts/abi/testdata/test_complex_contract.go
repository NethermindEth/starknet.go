package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/NethermindEth/starknet.go/abigen/accounts/abi"
)

func TestComplexContractABI(t *testing.T) {
	abiFile, err := os.Open(filepath.Join("testdata", "complex_contract.json"))
	if err != nil {
		t.Fatalf("Failed to open ABI file: %v", err)
	}
	defer abiFile.Close()
	
	abiBytes, err := io.ReadAll(abiFile)
	if err != nil {
		t.Fatalf("Failed to read ABI file: %v", err)
	}
	
	reader := strings.NewReader(string(abiBytes))
	parsedABI, err := abi.JSON(reader)
	if err != nil {
		t.Fatalf("ABI parsing failed: %v", err)
	}

	if len(parsedABI.Methods) != 3 {
		t.Fatalf("expected 3 methods, got %d", len(parsedABI.Methods))
	}

	transfer, exists := parsedABI.Methods["transfer"]
	if !exists {
		t.Fatalf("method 'transfer' not found in parsed ABI")
	}
	if transfer.StateMutability != "external" {
		t.Errorf("expected state mutability 'external', got '%s'", transfer.StateMutability)
	}
	if len(transfer.Inputs) != 2 {
		t.Errorf("expected 2 inputs, got %d", len(transfer.Inputs))
	}
	if transfer.Inputs[0].Name != "recipient" {
		t.Errorf("expected input name 'recipient', got '%s'", transfer.Inputs[0].Name)
	}
	if transfer.Inputs[0].Type != "core::starknet::ContractAddress" {
		t.Errorf("expected input type 'core::starknet::ContractAddress', got '%s'", transfer.Inputs[0].Type)
	}
	if transfer.Inputs[1].Name != "amount" {
		t.Errorf("expected input name 'amount', got '%s'", transfer.Inputs[1].Name)
	}
	if transfer.Inputs[1].Type != "core::integer::u256" {
		t.Errorf("expected input type 'core::integer::u256', got '%s'", transfer.Inputs[1].Type)
	}

	getBalances, exists := parsedABI.Methods["get_balances"]
	if !exists {
		t.Fatalf("method 'get_balances' not found in parsed ABI")
	}
	if getBalances.StateMutability != "view" {
		t.Errorf("expected state mutability 'view', got '%s'", getBalances.StateMutability)
	}
	if len(getBalances.Inputs) != 1 {
		t.Errorf("expected 1 input, got %d", len(getBalances.Inputs))
	}
	if getBalances.Inputs[0].Type != "core::array::Array<core::starknet::ContractAddress>" {
		t.Errorf("expected input type 'core::array::Array<core::starknet::ContractAddress>', got '%s'", getBalances.Inputs[0].Type)
	}
	if len(getBalances.Outputs) != 1 {
		t.Errorf("expected 1 output, got %d", len(getBalances.Outputs))
	}
	if getBalances.Outputs[0].Type != "core::array::Array<core::integer::u256>" {
		t.Errorf("expected output type 'core::array::Array<core::integer::u256>', got '%s'", getBalances.Outputs[0].Type)
	}

	if len(parsedABI.Events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(parsedABI.Events))
	}
	
	transferEvent, exists := parsedABI.Events["Transfer"]
	if !exists {
		t.Fatalf("event 'Transfer' not found in parsed ABI")
	}
	if len(transferEvent.Keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(transferEvent.Keys))
	}
	if len(transferEvent.Data) != 1 {
		t.Errorf("expected 1 data field, got %d", len(transferEvent.Data))
	}
}
