package contracts

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestUnmarshalContractClass is a test function to unmarshal a contract class.
//
// It reads the content of a file and unmarshals it into a ContractClass struct.
// Then it asserts the equality of certain values within the struct.
//
// Parameters:
// - t: The testing.T instance for running the test
// Returns:
//
//	none
func TestUnmarshalContractClass(t *testing.T) {
	content, err := os.ReadFile("./tests/hello_starknet_compiled.sierra.json")
	require.NoError(t, err)

	var class rpc.ContractClass
	err = json.Unmarshal(content, &class)
	require.NoError(t, err)
	assert.Equal(t, class.SierraProgram[0].String(), "0x1")
	assert.Equal(t, class.SierraProgram[1].String(), "0x3")
}

// TestUnmarshalCasmClass tests the UnmarshalCasmClass function.
//
// It reads the content of a file and unmarshals it into a CasmClass struct.
// The function uses the require.NoError and assert.Equal functions from the github.com/stretchr/testify/assert package to perform the assertions.
// It is a test function and is meant to be used with the Go testing framework.
//
// Parameters:
// - t: The testing.T instance for running the test
// Returns:
//
//	none
func TestUnmarshalCasmClass(t *testing.T) {
	casmClass, err := UnmarshalCasmClass("./tests/hello_starknet_compiled.casm.json")
	require.NoError(t, err)
	assert.Equal(t, casmClass.Prime, "0x800000000000011000000000000000000000000000000000000000000000001")
	assert.Equal(t, casmClass.Version, "2.1.0")
	assert.Equal(t, casmClass.EntryPointByType.External[0].Selector.String(), "0x362398bec32bc0ebb411203221a35a0301193a96f317ebe5e40be9f60d15320")
	assert.Equal(t, casmClass.EntryPointByType.External[1].Offset, 130)
	assert.Equal(t, casmClass.EntryPointByType.External[1].Builtins[0], "range_check")
}

// TestPrecomputeAddress tests the PrecomputeAddress function.
//
// It calls the PrecomputeAddress with predefined parameter values and compares the result with predefined expected results.
// The function uses the 'require' .NoError and .Equal functions from the github.com/stretchr/testify/assert package to perform the assertions.
// It is a test function and is meant to be used with the Go testing framework.
//
// Parameters:
// - t: The testing.T instance for running the test
// Returns:
//
//	none
func TestPrecomputeAddress(t *testing.T) {
	type testSetType struct {
		DeployerAddress            string
		Salt                       string
		ClassHash                  string
		ConstructorCalldata        []*felt.Felt
		ExpectedPrecomputedAddress string
	}

	testSet := []testSetType{
		{ //https://sepolia.voyager.online/tx/0x3789fe05652c9b18b98750b840e64cd3cc737592012c40d3233170d099db507
			DeployerAddress: "0",
			Salt:            "0x0702e82f1ec15656ad4502268dad530197141f3b59f5529835af9318ef399da5",
			ClassHash:       "0x064728e0c0713811c751930f8d3292d683c23f107c89b0a101425d9e80adb1c0",
			ConstructorCalldata: []*felt.Felt{
				utils.TestHexToFelt(t, "0x022f3e55b61d86c2ac5239fa3b3b8761f26b9a5c0b5f61ddbd5d756ced498b46"),
			},
			ExpectedPrecomputedAddress: "0x31463b5263a6631be4d1fe92d64d13e3a8498c440bf789e69ccb951eb8ad5da",
		},
		{ //https://sepolia.voyager.online/tx/0x7a4458b402a172e730c947b293a499d310a7ae6cfb18b5d9774fc10625927e5
			DeployerAddress: "0",
			Salt:            "0x023a851e8aeba201772098e1a1db3448f6238b20f928527242eb383905d91a87",
			ClassHash:       "0x061dac032f228abef9c6626f995015233097ae253a7f72d68552db02f2971b8f",
			ConstructorCalldata: []*felt.Felt{
				utils.TestHexToFelt(t, "0x023a851e8aeba201772098e1a1db3448f6238b20f928527242eb383905d91a87"),
			},
			ExpectedPrecomputedAddress: "0x28771beb7a2522a07d2ae6fc1fa5af942e8e863f70e6d7d74f9600ea3d5c242",
		},
		{ //https://sepolia.voyager.online/tx/0x2419a80d80045dd08cdb2606850c4eaf0ed8e705ee07bb1837d8daf12263bc0
			DeployerAddress: "0",
			Salt:            "0x0702e82f1ec15656ad4502268dad530197141f3b59f5529835af9318ef399da5",
			ClassHash:       "0xf6f44afb3cacbcc01a371aff62c86ca9a45feba065424c99f7cd8637514d8f",
			ConstructorCalldata: []*felt.Felt{
				utils.TestHexToFelt(t, "0x022f3e55b61d86c2ac5239fa3b3b8761f26b9a5c0b5f61ddbd5d756ced498b46"),
			},
			ExpectedPrecomputedAddress: "0x50cb9257feb7e960c8ab7d1cf48f33cfbe21de138409be476f63203383ece63",
		},
	}

	for _, test := range testSet {
		precomputedAddress := PrecomputeAddress(
			utils.TestHexToFelt(t, test.DeployerAddress),
			utils.TestHexToFelt(t, test.Salt),
			utils.TestHexToFelt(t, test.ClassHash),
			test.ConstructorCalldata,
		)
		require.Equal(t, test.ExpectedPrecomputedAddress, precomputedAddress.String())
	}
}
