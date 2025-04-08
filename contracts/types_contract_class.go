package contracts

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/NethermindEth/juno/core/felt"
)

type DeprecatedContractClass struct {
	// Program A base64 representation of the compressed program code
	Program string `json:"program"`

	DeprecatedEntryPointsByType DeprecatedEntryPointsByType `json:"entry_points_by_type"`

	ABI *ABI `json:"abi,omitempty"`
}

// UnmarshalJSON unmarshals JSON content into the DeprecatedContractClass struct.
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

	depEntryPointsByType := DeprecatedEntryPointsByType{}
	if err := json.Unmarshal(data, &depEntryPointsByType); err != nil {
		return err
	}
	c.DeprecatedEntryPointsByType = depEntryPointsByType

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

// encodeProgram encodes the content byte array using gzip compression and base64 encoding.
//
// It takes a content byte array as a parameter and returns the encoded program string and an error.
//
// Parameters:
// - content: byte array to be encoded
// Returns:
// - string: the encoded program
// - error: the error if any
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

type DeprecatedEntryPointsByType struct {
	Constructor []DeprecatedCairoEntryPoint `json:"CONSTRUCTOR"`
	External    []DeprecatedCairoEntryPoint `json:"EXTERNAL"`
	L1Handler   []DeprecatedCairoEntryPoint `json:"L1_HANDLER"`
}

type DeprecatedCairoEntryPoint struct {
	// The offset of the entry point in the program
	Offset NumAsHex `json:"offset"`
	// A unique  identifier of the entry point (function) in the program
	Selector *felt.Felt `json:"selector"`
}

type ABI []ABIEntry

type ContractClass struct {
	// The list of Sierra instructions of which the program consists
	SierraProgram []*felt.Felt `json:"sierra_program"`

	// The version of the contract class object. Currently, the Starknet OS supports version 0.1.0
	ContractClassVersion string `json:"contract_class_version"`

	EntryPointsByType SierraEntryPointsByType `json:"entry_points_by_type"`

	ABI NestedString `json:"abi,omitempty"`
}

type NestedString string

func (ns *NestedString) UnmarshalJSON(data []byte) error {
	var temp interface{}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	if value, ok := temp.(string); ok {
		// For cairo compiler prior to 2.7.0, the ABI is a string
		*ns = NestedString(value)
	} else {

		var out bytes.Buffer
		err := json.Indent(&out, data, "", "")

		if err != nil {
			return err
		}

		// Replace '\n' to ''
		out_str := bytes.ReplaceAll(out.Bytes(), []byte{10}, []byte{})
		// Replace ',"' to ', "'
		out_str = bytes.ReplaceAll(out_str, []byte{44, 34}, []byte{44, 32, 34})
		// Replace ',{' to ', {'
		out_str = bytes.ReplaceAll(out_str, []byte{44, 123}, []byte{44, 32, 123})

		*ns = NestedString(out_str)
	}

	return nil
}

type SierraEntryPoint struct {
	// The index of the function in the program
	FunctionIdx int `json:"function_idx"`
	// A unique  identifier of the entry point (function) in the program
	Selector *felt.Felt `json:"selector"`
}

type SierraEntryPointsByType struct {
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

type FunctionStateMutability string

const (
	FuncStateMutVIEW FunctionStateMutability = "view"
)

type FunctionABIEntry struct {
	// The function type
	Type ABIType `json:"type"`

	// The function name
	Name string `json:"name"`

	StateMutability FunctionStateMutability `json:"stateMutability,omitempty"`

	Inputs []TypedParameter `json:"inputs"`

	Outputs []TypedParameter `json:"outputs"`
}

// IsType returns the ABIType of the StructABIEntry.
func (s *StructABIEntry) IsType() ABIType {
	return s.Type
}

// IsType returns the ABIType of the EventABIEntry.
func (e *EventABIEntry) IsType() ABIType {
	return e.Type
}

// IsType returns the ABIType of the FunctionABIEntry.
func (f *FunctionABIEntry) IsType() ABIType {
	return f.Type
}

type TypedParameter struct {
	// The parameter's name
	Name string `json:"name"`

	// The parameter's type
	Type string `json:"type"`
}
