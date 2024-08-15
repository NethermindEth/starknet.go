package rpc

import (
	"encoding/json"
	"os"
	"testing"
)

const (
	validDeprecatedContractCompiledPath = "./tests/contract/0x01b661756bf7d16210fc611626e1af4569baa1781ffc964bd018f4585ae241c1.json"
	validContractCompiledPath           = "./tests/contract/0x03e9b96873987da76121f74a3df71e38c44527d8ce2ad115bcfda3cba0548cc3.json"
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
//
//	none
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
//
//	none
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
