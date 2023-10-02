package rpc

import (
	"encoding/json"
	"os"
	"testing"
)

const (
	validDeprecatedContractCompiledPath = "./tests/0x1efa8f84fd4dff9e2902ec88717cf0dafc8c188f80c3450615944a469428f7f.json"
	validContractCompiledPath           = "./tests/0x03a8Bad0A71696fC3eB663D0513Dc165Bb42cD4b662e633e3F87a49627CF3AEF.json"
	invalidContractCompiledPath         = "./tests/0xFakeContract.json"
)

// TestDeprecatedContractClass_UnmarshalValidJSON_Successful is a test function that checks the successful unmarshalling of valid JSON into a DeprecatedContractClass instance.
//
// The function does the following:
// - Reads the content of a file specified by validDeprecatedContractCompiledPath.
// - Unmarshals the content into a DeprecatedContractClass instance.
//
// Parameters:
// - t: A pointer to a testing.T instance.
//
// Return type: None.
func TestDeprecatedContractClass_UnmarshalValidJSON_Successful(t *testing.T) {
	content, err := os.ReadFile(validDeprecatedContractCompiledPath)
	if err != nil {
		t.Fatal("should be able to read file", err)
	}

	contractClass := DeprecatedContractClass{}
	if err := json.Unmarshal(content, &contractClass); err != nil {
		t.Fatal("should be able unmarshall Class", err)
	}
}

// TestContractClass_UnmarshalValidJSON_Successful is a test function that unmarshalls a valid JSON into a ContractClass struct and checks for successful execution.
//
// The function does the following:
// - Reads the content of a file using os.ReadFile.
// - Unmarshalls the content into a ContractClass struct using json.Unmarshal.
//
// Parameters:
// - t: a testing.T object used for testing.
//
// Return type: None.
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
