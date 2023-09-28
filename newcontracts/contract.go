package newcontract

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/NethermindEth/juno/core/felt"
)

type CasmClass struct {
	Prime            string                     `json:"prime"`
	Version          string                     `json:"compiler_version"`
	ByteCode         []*felt.Felt               `json:"bytecode"`
	EntryPointByType CasmClassEntryPointsByType `json:"entry_points_by_type"`
	// Hints            any                        `json:"hints"`
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

func UnmarshalCasmClass(filePath string) (*CasmClass, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
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
