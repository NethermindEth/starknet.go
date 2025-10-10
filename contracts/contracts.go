package contracts

import (
	"encoding/json"
	"os"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/curve"
)

var prefixContractAddress = new(felt.Felt).SetBytes([]byte("STARKNET_CONTRACT_ADDRESS"))

// UnmarshalCasmClass is a function that unmarshals a CasmClass object from a file.
// CASM = Cairo instructions
//
// It takes a file path as a parameter and returns a pointer to the
// unmarshaled CasmClass object and an error.
func UnmarshalCasmClass(filePath string) (*CasmClass, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var casmClass CasmClass
	err = json.Unmarshal(content, &casmClass)
	if err != nil {
		return nil, err
	}

	return &casmClass, nil
}

// PrecomputeAddress calculates the precomputed address for a contract instance.
// ref: https://github.com/starkware-libs/cairo-lang/blob/master/src/starkware/starknet/core/os/contract_address/contract_address.py
//
// Parameters:
//   - deployerAddress: the deployer address
//   - salt: the salt
//   - classHash: the class hash
//   - constructorCalldata: the constructor calldata
//
// Returns:
//   - *felt.Felt: the precomputed address as a *felt.Felt
//
//nolint:lll // The link would be unclickable if we break the line.
func PrecomputeAddress(
	deployerAddress, salt, classHash *felt.Felt,
	constructorCalldata []*felt.Felt,
) *felt.Felt {
	return curve.PedersenArray(
		prefixContractAddress,
		deployerAddress,
		salt,
		classHash,
		curve.PedersenArray(constructorCalldata...),
	)
}
