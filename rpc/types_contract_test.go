package rpc

import (
	"encoding/json"
	"os"
	"testing"
)

const (
	validContractCompiledPath   = "./tests/0x1efa8f84fd4dff9e2902ec88717cf0dafc8c188f80c3450615944a469428f7f.json"
	invalidContractCompiledPath = "./tests/0xFakeContract.json"
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
