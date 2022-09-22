package types

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

const validContractCompiledPath = "./testdata/cairo/minimum_contract_compiled.json"
const invalidContractCompiledPath = "./testdata/cairo/invalid_minimum_contract_compiled.json"

func TestContractClass_UnmarshalValidJSON_Successful(t *testing.T) {
	content, err := os.ReadFile(validContractCompiledPath)
	assert.NoError(t, err)

	contractClass := ContractClass{}
	err = contractClass.UnmarshalJSON(content)
	assert.NoError(t, err)
}

func TestContractClass_UnmarshalInvalidJSON_Fails(t *testing.T) {
	content, err := os.ReadFile(invalidContractCompiledPath)
	assert.NoError(t, err)
	contractClass := ContractClass{}
	err = contractClass.UnmarshalJSON(content)
	assert.ErrorContains(t, err, "abi is not iterable")
}
