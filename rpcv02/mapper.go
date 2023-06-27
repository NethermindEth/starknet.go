package rpcv02

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/NethermindEth/caigo/types"
)

type functionInvoke types.FunctionInvoke

func (f functionInvoke) MarshalJSON() ([]byte, error) {
	output := map[string]interface{}{}
	sigs := []string{}
	for _, sig := range f.Signature {
		sigs = append(sigs, fmt.Sprintf("0x%x", sig))
	}
	output["signature"] = sigs
	v, err := json.Marshal(f)
	if err != nil {
		return nil, err
	}
	functionCall := map[string]json.RawMessage{}
	err = json.Unmarshal(v, &functionCall)
	if err != nil {
		return nil, err
	}
	output["contract_address"] = functionCall["contract_address"]
	if selector, ok := functionCall["entry_point_selector"]; ok {
		output["entry_point_selector"] = selector
	}
	output["calldata"] = functionCall["calldata"]
	if f.Nonce != nil {
		output["nonce"] = json.RawMessage(
			strconv.Quote(fmt.Sprintf("0x%x", f.Nonce)),
		)
	}
	if f.MaxFee != nil {
		output["max_fee"] = json.RawMessage(
			strconv.Quote(fmt.Sprintf("0x%x", f.MaxFee)),
		)
	}
	output["version"] = json.RawMessage(strconv.Quote(fmt.Sprintf("0x%x", f.Version)))
	return json.Marshal(output)
}
