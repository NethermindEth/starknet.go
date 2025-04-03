package abi

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/utils"
)

type ABI struct {
	Constructor Method
	Methods     map[string]Method
	Events      map[string]Event
	Structs     map[string]Struct
}

func JSON(reader io.Reader) (ABI, error) {
	dec := json.NewDecoder(reader)

	var abiEntries []json.RawMessage
	if err := dec.Decode(&abiEntries); err != nil {
		return ABI{}, err
	}

	abi := ABI{
		Methods: make(map[string]Method),
		Events:  make(map[string]Event),
		Structs: make(map[string]Struct),
	}

	for _, entry := range abiEntries {
		var entryMap map[string]interface{}
		if err := json.Unmarshal(entry, &entryMap); err != nil {
			return ABI{}, err
		}

		entryType, ok := entryMap["type"].(string)
		if !ok {
			return ABI{}, fmt.Errorf("ABI entry missing type field")
		}

		switch entryType {
		case "constructor":
			var method Method
			if err := json.Unmarshal(entry, &method); err != nil {
				return ABI{}, err
			}
			abi.Constructor = method
		case "function", "l1_handler":
			var method Method
			if err := json.Unmarshal(entry, &method); err != nil {
				return ABI{}, err
			}
			abi.Methods[method.Name] = method
		case "event":
			var event Event
			if err := json.Unmarshal(entry, &event); err != nil {
				return ABI{}, err
			}
			abi.Events[event.Name] = event
		case "struct":
			var structType Struct
			if err := json.Unmarshal(entry, &structType); err != nil {
				return ABI{}, err
			}
			abi.Structs[structType.Name] = structType
		default:
			return ABI{}, fmt.Errorf("unknown ABI entry type: %s", entryType)
		}
	}

	return abi, nil
}

type Method struct {
	Type           string         `json:"type"`
	Name           string         `json:"name"`
	StateMutability string         `json:"state_mutability,omitempty"`
	Inputs         []Argument     `json:"inputs"`
	Outputs        []Argument     `json:"outputs"`
	Selector       *felt.Felt     // Computed selector
}

type Event struct {
	Type     string     `json:"type"`
	Name     string     `json:"name"`
	Kind     string     `json:"kind,omitempty"`
	Keys     []Argument `json:"keys"`
	Data     []Argument `json:"data"`
	Variants []Variant  `json:"variants,omitempty"`
}

type Variant struct {
	Name string     `json:"name"`
	Type string     `json:"type,omitempty"`
	Keys []Argument `json:"keys,omitempty"`
	Data []Argument `json:"data,omitempty"`
}

type Struct struct {
	Type    string     `json:"type"`
	Name    string     `json:"name"`
	Size    uint64     `json:"size,omitempty"`
	Members []Argument `json:"members"`
}

type Argument struct {
	Name   string `json:"name,omitempty"`
	Type   string `json:"type"`
	Offset int64  `json:"offset,omitempty"` // Only for struct members
}

func GetSelector(name string) *felt.Felt {
	return utils.GetSelectorFromNameFelt(name)
}

func (abi *ABI) PackArguments(args []Argument, values []interface{}) ([]*felt.Felt, error) {
	return PackArguments(args, values)
}
