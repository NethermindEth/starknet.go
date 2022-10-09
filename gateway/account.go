package gateway

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/dontpanicdao/caigo/types"
)

func (sg *Gateway) AccountNonce(ctx context.Context, address types.Hash) (*big.Int, error) {
	resp, err := sg.Call(ctx, types.FunctionCall{
		ContractAddress:    address,
		EntryPointSelector: "get_nonce",
	}, "")
	if err != nil {
		return nil, err
	}
	if len(resp) == 0 {
		return nil, fmt.Errorf("no resp in contract call 'get_nonce' %v", address)
	}

	return types.HexToBN(resp[0]), nil
}

func (sg *Gateway) Nonce(ctx context.Context, contractAddress, blockHashOrTag string) (*big.Int, error) {

	req, err := sg.newRequest(ctx, http.MethodGet, "/get_nonce", nil)
	if err != nil {
		return nil, err
	}

	appendQueryValues(req, url.Values{
		"contractAddress": []string{contractAddress},
	})
	switch {
	case strings.HasPrefix(blockHashOrTag, "0x"):
		appendQueryValues(req, url.Values{
			"blockHash": []string{blockHashOrTag},
		})
	case blockHashOrTag == "":
		appendQueryValues(req, url.Values{
			"blockNumber": []string{"pending"},
		})
	default:
		appendQueryValues(req, url.Values{
			"blockNumber": []string{blockHashOrTag},
		})
	}

	var resp string
	err = sg.do(req, &resp)
	if err != nil {
		return nil, err
	}
	nonce, ok := big.NewInt(0).SetString(resp, 0)
	if !ok {
		return nil, errors.New("nonce not found")
	}
	return nonce, nil
}

type functionInvoke types.FunctionInvoke

func (f functionInvoke) MarshalJSON() ([]byte, error) {
	output := map[string]interface{}{}
	sigs := []string{}
	for _, sig := range f.Signature {
		sigs = append(sigs, sig.Text(10))
	}
	output["signature"] = sigs
	v, err := json.Marshal(f.FunctionCall)
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
	calldataSlice := []string{}
	err = json.Unmarshal(functionCall["calldata"], &calldataSlice)
	if err != nil {
		return nil, err
	}

	calldata := []string{}
	for _, v := range calldataSlice {
		data, _ := big.NewInt(0).SetString(v, 0)
		calldata = append(calldata, data.Text(10))
	}
	output["calldata"] = calldata
	if f.Nonce != nil {
		output["nonce"] = json.RawMessage(
			strconv.Quote(fmt.Sprintf("0x%s", f.Nonce.Text(16))),
		)
	}
	if f.MaxFee != nil {
		output["max_fee"] = json.RawMessage(
			strconv.Quote(fmt.Sprintf("0x%s", f.MaxFee.Text(16))),
		)
	}
	output["version"] = json.RawMessage(strconv.Quote(fmt.Sprintf("0x%d", f.Version)))
	return json.Marshal(output)
}

func (sg *Gateway) EstimateFee(ctx context.Context, call types.FunctionInvoke, hash string) (*types.FeeEstimate, error) {
	call.EntryPointSelector = types.BigToHex(types.GetSelectorFromName(call.EntryPointSelector))
	c := functionInvoke(call)
	req, err := sg.newRequest(ctx, http.MethodPost, "/estimate_fee", c)
	if err != nil {
		return nil, err
	}

	if hash != "" {
		appendQueryValues(req, url.Values{
			"blockHash": []string{hash},
		})
	}
	output := map[string]interface{}{}
	err = sg.do(req, &output)
	if err != nil {
		return nil, err
	}
	gasPrice, _ := output["gas_price"].(int)
	gasConsumed, _ := output["gas_usage"].(int)
	overallFee, _ := output["overall_fee"].(int)
	resp := types.FeeEstimate{
		GasConsumed: types.NumAsHex("0x" + big.NewInt(int64(gasConsumed)).Text(16)),
		GasPrice:    types.NumAsHex("0x" + big.NewInt(int64(gasPrice)).Text(16)),
		OverallFee:  types.NumAsHex("0x" + big.NewInt(int64(overallFee)).Text(16)),
	}
	return &resp, nil
}
