package rpc

import (
	"encoding/json"
	"os"
	"testing"
)

const (
	validDeprecatedContractCompiledPath = "./tests/contract/0x1efa8f84fd4dff9e2902ec88717cf0dafc8c188f80c3450615944a469428f7f.json"
	validContractCompiledPath           = "./tests/contract/0x03a8Bad0A71696fC3eB663D0513Dc165Bb42cD4b662e633e3F87a49627CF3AEF.json"
	invalidContractCompiledPath         = "./tests/0xFakeContract.json"
)

// TestDeprecatedContractClass_UnmarshalValidJSON_Successful is a test function that
// tests the successful unmarshalling of valid JSON into a DeprecatedContractClass
// object.
//
// It reads the content of a file, then unmarshals the content into a 
// DeprecatedContractClass object using the json.Unmarshal function. If any error
// occurs during the process, the test fails.
//
// Parameters:
// - t: The testing.T object used for reporting test failures and logging.
// Returns:
//  none
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

// TestContractClass_UnmarshalValidJSON_Successful is a test function that validates the successful unmarshalling of valid JSON data into a ContractClass object.
//
// The function does the following:
// - Reads the content of a file specified by the validContractCompiledPath variable.
// - Unmarshals the content into a ContractClass object using the json.Unmarshal function.
//
// Parameters:
// - t: The testing.T object used for reporting test failures and logging.
// Returns:
//  none
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
