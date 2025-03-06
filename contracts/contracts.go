package contracts

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/curve"
)

var PREFIX_CONTRACT_ADDRESS = new(felt.Felt).SetBytes([]byte("STARKNET_CONTRACT_ADDRESS"))

type NestedUints struct {
	IsArray bool
	Value   *uint64
	Values  []NestedUints
}

func toNestedInts(values []any) ([]NestedUints, error) {

	var res []NestedUints = make([]NestedUints, 0)

	for _, value := range values {
		// check if the value is a single number
		if numeric, ok := value.(float64); ok {
			intVal := uint64(numeric)
			res = append(res, NestedUints{
				IsArray: false,
				Value:   &intVal,
				Values:  nil,
			})
			continue
		}

		// check if the value is an array
		if arrVal, ok := value.([]any); ok {
			nested, err := toNestedInts(arrVal)
			if err != nil {
				return nil, err
			}

			res = append(res, NestedUints{
				IsArray: true,
				Value:   nil,
				Values:  nested,
			})
			continue
		}

		return nil, errors.New("invalid NestedUint type")
	}

	return res, nil
}

func (ns *NestedUints) UnmarshalJSON(data []byte) error {
	var temp []any
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	nested, err := toNestedInts(temp)

	if err != nil {
		return err
	}

	*ns = NestedUints{
		IsArray: true,
		Value:   nil,
		Values:  nested,
	}

	return nil
}

// MarshalJSON implements the json.Marshaler interface for NestedUints.
// It converts the NestedUints structure back into a JSON array format.
func (ns NestedUints) MarshalJSON() ([]byte, error) {
	if !ns.IsArray {
		if ns.Value == nil {
			return nil, errors.New("invalid NestedUints: non-array type must have a value")
		}
		return json.Marshal(*ns.Value)
	}

	result := make([]any, len(ns.Values))
	for i, v := range ns.Values {
		if !v.IsArray {
			if v.Value == nil {
				return nil, errors.New("invalid NestedUints: non-array type must have a value")
			}
			result[i] = *v.Value
		} else {
			nestedJSON, err := v.MarshalJSON()
			if err != nil {
				return nil, err
			}
			var nestedValue any
			if err := json.Unmarshal(nestedJSON, &nestedValue); err != nil {
				return nil, err
			}
			result[i] = nestedValue
		}
	}
	return json.Marshal(result)
}

type CasmClass struct {
	Prime                  string                     `json:"prime"`
	Version                string                     `json:"compiler_version"`
	ByteCode               []*felt.Felt               `json:"bytecode"`
	EntryPointByType       CasmClassEntryPointsByType `json:"entry_points_by_type"`
	BytecodeSegmentLengths *NestedUints               `json:"bytecode_segment_lengths,omitempty"`
}

type CasmClassEntryPointsByType struct {
	Constructor []CasmClassEntryPoint `json:"CONSTRUCTOR"`
	External    []CasmClassEntryPoint `json:"EXTERNAL"`
	L1Handler   []CasmClassEntryPoint `json:"L1_HANDLER"`
}

type CasmClassEntryPoint struct {
	Selector *felt.Felt `json:"selector"`
	Offset   int        `json:"offset"`
	Builtins []string   `json:"builtins"`
}

// UnmarshalCasmClass is a function that unmarshals a CasmClass object from a file.
// CASM = Cairo instructions
//
// It takes a file path as a parameter and returns a pointer to the unmarshaled CasmClass object and an error.
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
func PrecomputeAddress(deployerAddress *felt.Felt, salt *felt.Felt, classHash *felt.Felt, constructorCalldata []*felt.Felt) *felt.Felt {
	return curve.PedersenArray(
		PREFIX_CONTRACT_ADDRESS,
		deployerAddress,
		salt,
		classHash,
		curve.PedersenArray(constructorCalldata...),
	)
}
