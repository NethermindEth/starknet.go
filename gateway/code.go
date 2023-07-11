package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"net/url"

	"github.com/NethermindEth/starknet.go/rpcv02"
)

type Bytecode []string

type Code struct {
	Bytecode Bytecode    `json:"bytecode"`
	Abi      *rpcv02.ABI `json:"abi"`
}

func (c *Code) UnmarshalJSON(content []byte) error {
	v := map[string]json.RawMessage{}
	if err := json.Unmarshal(content, &v); err != nil {
		return err
	}

	// process 'bytecode'.
	data, ok := v["bytecode"]
	if !ok {
		return fmt.Errorf("missing bytecode in json object")
	}
	bytecode := []string{}
	if err := json.Unmarshal(data, &bytecode); err != nil {
		return err
	}
	c.Bytecode = bytecode

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

	abiPointer := rpcv02.ABI{}
	for _, abi := range abis {
		if checkABI, ok := abi.(map[string]interface{}); ok {
			var ab rpcv02.ABIEntry
			abiType, ok := checkABI["type"].(string)
			if !ok {
				return fmt.Errorf("unknown abi type %v", checkABI["type"])
			}
			switch abiType {
			case string(rpcv02.ABITypeConstructor), string(rpcv02.ABITypeFunction), string(rpcv02.ABITypeL1Handler):
				ab = &rpcv02.FunctionABIEntry{}
			case string(rpcv02.ABITypeStruct):
				ab = &rpcv02.StructABIEntry{}
			case string(rpcv02.ABITypeEvent):
				ab = &rpcv02.EventABIEntry{}
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

// Gets a contracts code.
//
// [Reference](https://github.com/starkware-libs/cairo-lang/blob/fc97bdd8322a7df043c87c371634b26c15ed6cee/src/starkware/starknet/services/api/feeder_gateway/feeder_gateway_client.py#L55)
func (sg *Gateway) CodeAt(ctx context.Context, contract string, blockNumber *big.Int) (*Code, error) {
	req, err := sg.newRequest(ctx, http.MethodGet, "/get_code", nil)
	if err != nil {
		return nil, err
	}

	appendQueryValues(req, url.Values{"contractAddress": []string{contract}})

	if blockNumber != nil {
		appendQueryValues(req, url.Values{"blockNumber": []string{blockNumber.String()}})
	}

	var resp Code
	return &resp, sg.do(req, &resp)
}

func (sg *Gateway) FullContract(ctx context.Context, contract string) (*rpcv02.ContractClass, error) {
	req, err := sg.newRequest(ctx, http.MethodGet, "/get_full_contract", nil)
	if err != nil {
		return nil, err
	}

	appendQueryValues(req, url.Values{"contractAddress": []string{contract}})

	var resp rpcv02.ContractClass
	return &resp, sg.do(req, &resp)
}
