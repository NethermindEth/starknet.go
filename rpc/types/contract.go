package types

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

type EntryPoint struct {
	// The offset of the entry point in the program
	Offset NumAsHex `json:"offset"`
	// A unique identifier of the entry point (function) in the program
	Selector string `json:"selector"`
}

type ABI []ABIEntry

type EntryPointsByType struct {
	Constructor []EntryPoint `json:"CONSTRUCTOR"`
	External    []EntryPoint `json:"EXTERNAL"`
	L1Handler   []EntryPoint `json:"L1_HANDLER"`
}

type ContractClass struct {
	// Program A base64 representation of the compressed program code
	Program string `json:"program"`

	EntryPointsByType EntryPointsByType `json:"entry_points_by_type"`

	Abi *ABI `json:"abi,omitempty"`
}

func (c *ContractClass) UnmarshalJSON(content []byte) error {
	v := map[string]json.RawMessage{}
	if err := json.Unmarshal(content, &v); err != nil {
		return err
	}

	// process 'program'
	data, ok := v["program"]
	if !ok {
		return fmt.Errorf("missing program in json object")
	}

	program, err := encodeProgram(data)
	if err != nil {
		return err
	}
	c.Program = program

	// process 'entry_points_by_type'
	data, ok = v["entry_points_by_type"]
	if !ok {
		return fmt.Errorf("missing entry_points_by_type in json object")
	}

	entryPointsByType := EntryPointsByType{}
	err = json.Unmarshal(data, &entryPointsByType)
	if err != nil {
		return err
	}
	c.EntryPointsByType = entryPointsByType

	// process 'abi'
	data, ok = v["abi"]
	if !ok {
		return fmt.Errorf("missing abi in json object")
	}

	abis := []interface{}{}
	err = json.Unmarshal(data, &entryPointsByType)
	if err != nil {
		return err
	}

	abiPointer := ABI{}
	for _, abi := range abis {
		if checkABI, ok := abi.(map[string]interface{}); ok {
			var ab ABIEntry
			switch checkABI["type"] {
			case "constructor", "function", "l1_handler":
				ab = &FunctionABIEntry{}
			case "struct":
				ab = &StructABIEntry{}
			case "event":
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

	c.Abi = &abiPointer
	return nil
}

type ABIEntry interface {
	IsType() string
}

type StructABIType string

const (
	StructABITypeEvent StructABIType = "struct"
)

type EventABIType string

const (
	EventABITypeEvent EventABIType = "event"
)

type FunctionABIType string

const (
	FunctionABITypeFunction  FunctionABIType = "function"
	FunctionABITypeL1Handler FunctionABIType = "l1_handler"
)

type StructABIEntry struct {
	// The event type
	Type StructABIType `json:"type"`

	// The event name
	Name string `json:"name"`

	Size uint64 `json:"size"`

	Members []StructMember `json:"members"`
}

type StructMember struct {
	TypedParameter
	Offset uint64 `json:"offset"`
}

type EventABIEntry struct {
	// The event type
	Type EventABIType `json:"type"`

	// The event name
	Name string `json:"name"`

	Keys []TypedParameter `json:"keys"`

	Data TypedParameter `json:"data"`
}

type FunctionABIEntry struct {
	// The function type
	Type FunctionABIType `json:"type"`

	// The function name
	Name string `json:"name"`

	Inputs []TypedParameter `json:"inputs"`

	Outputs []TypedParameter `json:"outputs"`
}

type TypedParameter struct {
	// The parameter's name
	Name string `json:"name"`

	// The parameter's type
	Type string `json:"type"`
}

// encodeProgram Keep that function to build helper with broadcastedDeployTxn and broadcastedDeclareTxn
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
