package types

import (
	"encoding/json"
	"os"
	"testing"
)

const (
	validContractCompiledPath   = "./testdata/cairo/minimum_contract_compiled.json"
	invalidContractCompiledPath = "./testdata/cairo/invalid_minimum_contract_compiled.json"
)

func TestContractClass_UnmarshalValidJSON_Successful(t *testing.T) {
	content, err := os.ReadFile(validContractCompiledPath)
	if err != nil {
		t.Fatal("should be able to read file", err)
	}

	contractClass := ContractClass{}
	if err := json.Unmarshal(content, &contractClass); err != nil {
		t.Fatal("should be able unmarshall Class", err)
	}
}

func TestContractClass_UnmarshalInvalidJSON_Fails(t *testing.T) {
	content, err := os.ReadFile(invalidContractCompiledPath)
	if err != nil {
		t.Fatal("should be able to read file", err)
	}

	contractClass := ContractClass{}
	if err := json.Unmarshal(content, &contractClass); err != nil {
		t.Fatal("should be able unmarshall Class", err)
	}
}
