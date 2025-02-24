package hash_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/NethermindEth/starknet.go/contracts"
	"github.com/NethermindEth/starknet.go/hash"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestUnmarshalCasmClassHash is a test function that tests the unmarshaling of a CasmClass hash.
//
// It reads the content of the "./tests/hello_starknet_compiled.casm.json" file and unmarshals it into a contracts.CasmClass variable.
// The function returns an assertion error if there is an error reading the file or unmarshaling the content.
// Parameters:
// - t: A testing.T object used for running the test and reporting any failures.
// Returns:
//
//	none
func TestUnmarshalCasmClassHash(t *testing.T) {
	type testSetType struct {
		Description string
		CasmPath    string
	}
	testSet := []testSetType{
		{
			Description: "Compile: 2.1.0, no 'bytecode_segment_lengths' field",
			CasmPath:    "./tests/hello_starknet_compiled.compiled_contract_class.json",
		},
		{
			Description: "Compile: 2.7.0, 'bytecode_segment_lengths' with uint64 values",
			CasmPath:    "./tests/contracts_v2_TestContract.compiled_contract_class.json",
		},
	}

	for _, test := range testSet {
		t.Run(test.Description, func(t *testing.T) {
			content, err := os.ReadFile(test.CasmPath)
			require.NoError(t, err)

			var casmClass contracts.CasmClass
			err = json.Unmarshal(content, &casmClass)
			require.NoError(t, err)

			// TODO: uncomment and test this when merged with v.0.8.0 (v.0.8.0 branch has the missing fields of this branch, like 'Hints')
			// jsonResult, err := json.Marshal(casmClass)
			// require.NoError(t, err)

			// assert.JSONEq(t, string(content), string(jsonResult))
		})
	}
}

// TestClassHashes is a test function that verifies the correctness of the ClassHash and CompiledClassHash functions.
//
// It reads the contents of a file, unmarshals it into a SierraClass or CasmClass object, computes the hash using the ClassHash or CompiledClassHash function,
// and asserts that the computed hash matches the expected hash.
//
// Parameters:
func TestClassHashes(t *testing.T) {
	type testSetType struct {
		FileNameWithoutExtensions string
		ExpectedClassHash         string
		ExpectedCompiledClassHash string
	}

	// Ref ClassHash: https://github.com/software-mansion/starknet.py/blob/39af414389984efbc6edc48b0fe1f914ea5b9a77/starknet_py/tests/unit/hash/sierra_class_hash_test.py
	// Ref CompiledClassHash: https://github.com/software-mansion/starknet.py/blob/39af414389984efbc6edc48b0fe1f914ea5b9a77/starknet_py/tests/unit/hash/casm_class_hash_test.py
	testSet := []testSetType{
		{ // internal case, with "abi" field as string
			FileNameWithoutExtensions: "hello_starknet_compiled",
			ExpectedClassHash:         "0x4ec2ecf58014bc2ffd7c84843c3525e5ecb0a2cac33c47e9c347f39fc0c0944",
			ExpectedCompiledClassHash: "0x785fa5f2bacf0bfe3bc413be5820a61e1ea63f2ec27ef00331ee9f46ad07603",
		},
		{
			FileNameWithoutExtensions: "contracts_v2_Account",
			ExpectedClassHash:         "0x183078afce57a1d33b948ea6cd9ab0769dd08ca93a6afe4c23637b08aa893c1",
			ExpectedCompiledClassHash: "0x108977ab61715437fc7097b6499b3cf9491361eb6a8ce6df6c8536b7feec508",
		},
		{
			FileNameWithoutExtensions: "contracts_v2_ERC20",
			ExpectedClassHash:         "0x746248ba570006607113ae3f4dbb4130e81233fb818d15329c6a4aaccf94812",
			ExpectedCompiledClassHash: "0x5adc857416202a5902c01168542e188c3aa6380f57c911ae98cf20bc52be367",
		},
		{
			FileNameWithoutExtensions: "contracts_v2_HelloStarknet",
			ExpectedClassHash:         "0x224518978adb773cfd4862a894e9d333192fbd24bc83841dc7d4167c09b89c5",
			ExpectedCompiledClassHash: "0x6ff9f7df06da94198ee535f41b214dce0b8bafbdb45e6c6b09d4b3b693b1f17",
		},
		{
			FileNameWithoutExtensions: "contracts_v2_TestContract",
			ExpectedClassHash:         "0x3adac8a417b176d27e11b420aa1063b07a6b54bbb21091ad77b2a9156af7a3b",
			ExpectedCompiledClassHash: "0x2193add92c182c9236f0c156f11dc4f18d5a78fd9b763a3c0f4a1d3bd8b87d4",
		},
		{
			FileNameWithoutExtensions: "contracts_v2_TokenBridge",
			ExpectedClassHash:         "0x3d138e923f01b7ed1bb82b9b4e7f6df64e0c429faf8b27539addc71c1407237",
			ExpectedCompiledClassHash: "0x41d26534c7ca29e212ae48acfb9f86f69a9624977c979697c15f587fa95204",
		},
	}

	t.Run("Test Sierra ClassHash:", func(t *testing.T) {
		for _, test := range testSet {

			t.Run(test.FileNameWithoutExtensions, func(t *testing.T) {
				sierraClass := *utils.TestUnmarshallJSONToType[rpc.ContractClass](t, "./tests/"+test.FileNameWithoutExtensions+".contract_class.json", "")

				hash := hash.ClassHash(sierraClass)
				assert.Equal(t, test.ExpectedClassHash, hash.String())
			})
		}
	})

	t.Run("Test CompiledClassHash:", func(t *testing.T) {
		for _, test := range testSet {

			t.Run(test.FileNameWithoutExtensions, func(t *testing.T) {
				casmClass := *utils.TestUnmarshallJSONToType[contracts.CasmClass](t, "./tests/"+test.FileNameWithoutExtensions+".compiled_contract_class.json", "")

				hash, err := hash.CompiledClassHash(casmClass)
				require.NoError(t, err)
				assert.Equal(t, test.ExpectedCompiledClassHash, hash.String())
			})
		}
	})
}
