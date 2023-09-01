package rpc

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/NethermindEth/juno/core/felt"
)

type NumAsHex string

type DeprecatedCairoEntryPoint struct {
	// The offset of the entry point in the program
	Offset NumAsHex `json:"offset"`
	// A unique  identifier of the entry point (function) in the program
	Selector *felt.Felt `json:"selector"`
}

type ClassOutput interface{}

var _ ClassOutput = &DeprecatedContractClass{}
var _ ClassOutput = &ContractClass{}

type ABI []ABIEntry

type DeprecatedEntryPointsByType struct {
	Constructor []DeprecatedCairoEntryPoint `json:"CONSTRUCTOR"`
	External    []DeprecatedCairoEntryPoint `json:"EXTERNAL"`
	L1Handler   []DeprecatedCairoEntryPoint `json:"L1_HANDLER"`
}
type DeprecatedContractClass struct {
	// Program A base64 representation of the compressed program code
	Program string `json:"program"`

	EntryPointsByType DeprecatedEntryPointsByType `json:"entry_points_by_type"`

	ABI *ABI `json:"abi,omitempty"`
}

func (c *DeprecatedContractClass) UnmarshalJSON(content []byte) error {
	v := map[string]json.RawMessage{}
	if err := json.Unmarshal(content, &v); err != nil {
		return err
	}

	// process 'program'. If it is a string, keep it, otherwise encode it.
	data, ok := v["program"]
	if !ok {
		return fmt.Errorf("missing program in json object")
	}
	program := ""
	if err := json.Unmarshal(data, &program); err != nil {
		if program, err = encodeProgram(data); err != nil {
			return err
		}
	}
	c.Program = program

	// process 'entry_points_by_type'
	data, ok = v["entry_points_by_type"]
	if !ok {
		return fmt.Errorf("missing entry_points_by_type in json object")
	}

	entryPointsByType := DeprecatedEntryPointsByType{}
	if err := json.Unmarshal(data, &entryPointsByType); err != nil {
		return err
	}
	c.EntryPointsByType = entryPointsByType

	// process 'abi'
	data, ok = v["abi"]
	if !ok {
		// contractClass can have an empty ABI for instance with ClassAt
		return nil
	}

	abis := []interface{}{}
	if err := json.Unmarshal(data, &abis); err != nil {
		return err
	}

	abiPointer := ABI{}
	for _, abi := range abis {
		if checkABI, ok := abi.(map[string]interface{}); ok {
			var ab ABIEntry
			abiType, ok := checkABI["type"].(string)
			if !ok {
				return fmt.Errorf("unknown abi type %v", checkABI["type"])
			}
			switch abiType {
			case string(ABITypeConstructor), string(ABITypeFunction), string(ABITypeL1Handler):
				ab = &FunctionABIEntry{}
			case string(ABITypeStruct):
				ab = &StructABIEntry{}
			case string(ABITypeEvent):
				ab = &EventABIEntry{}
			default:
				return fmt.Errorf("unknown ABI type %v", checkABI["type"])
			}
			data, err := json.Marshal(checkABI)
			if err != nil {
				return err
			}
			err = json.Unmarshal(data, ab)
			if err != nil {
				return err
			}
			abiPointer = append(abiPointer, ab)
		}
	}

	c.ABI = &abiPointer
	return nil
}

// https://github.com/starkware-libs/starknet-specs/blob/v0.3.0/api/starknet_api_openrpc.json#L2372
type ContractClass struct {
	// The list of Sierra instructions of which the program consists
	SierraProgram []*felt.Felt `json:"sierra_program"`

	// The version of the contract class object. Currently, the Starknet OS supports version 0.1.0
	Version string `json:"contract_class_version"`

	EntryPointsByType EntryPointsByType `json:"entry_points_by_type"`

	ABI string `json:"abi,omitempty"`
}

type SierraEntryPoint struct {
	// The index of the function in the program
	FunctionIdx int `json:"function_idx"`
	// A unique  identifier of the entry point (function) in the program
	Selector *felt.Felt `json:"selector"`
}

type EntryPointsByType struct {
	Constructor []SierraEntryPoint `json:"CONSTRUCTOR"`
	External    []SierraEntryPoint `json:"EXTERNAL"`
	L1Handler   []SierraEntryPoint `json:"L1_HANDLER"`
}

type ABIEntry interface {
	IsType() ABIType
}

type ABIType string

const (
	ABITypeConstructor ABIType = "constructor"
	ABITypeFunction    ABIType = "function"
	ABITypeL1Handler   ABIType = "l1_handler"
	ABITypeEvent       ABIType = "event"
	ABITypeStruct      ABIType = "struct"
)

type StructABIEntry struct {
	// The event type
	Type ABIType `json:"type"`

	// The event name
	Name string `json:"name"`

	// todo(minumum size should be 1)
	Size uint64 `json:"size"`

	Members []Member `json:"members"`
}

type Member struct {
	TypedParameter
	Offset int64 `json:"offset"`
}

type EventABIEntry struct {
	// The event type
	Type ABIType `json:"type"`

	// The event name
	Name string `json:"name"`

	Keys []TypedParameter `json:"keys"`

	Data []TypedParameter `json:"data"`
}

type FunctionABIEntry struct {
	// The function type
	Type ABIType `json:"type"`

	// The function name
	Name string `json:"name"`

	StateMutability *string `json:"stateMutability,omitempty"`

	Inputs []TypedParameter `json:"inputs"`

	Outputs []TypedParameter `json:"outputs"`
}

func (s *StructABIEntry) IsType() ABIType {
	return s.Type
}

func (e *EventABIEntry) IsType() ABIType {
	return e.Type
}

func (f *FunctionABIEntry) IsType() ABIType {
	return f.Type
}

type TypedParameter struct {
	// The parameter's name
	Name string `json:"name"`

	// The parameter's type
	Type string `json:"type"`
}

// encodeProgram compress a program to send it to the API
func encodeProgram(content []byte) (string, error) {
	buf := bytes.NewBuffer(nil)
	gzipContent := gzip.NewWriter(buf)
	_, err := gzipContent.Write(content)
	if err != nil {
		return "", err
	}
	gzipContent.Close()
	program := base64.StdEncoding.EncodeToString(buf.Bytes())
	return program, nil
}
