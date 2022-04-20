package caigo

import (
	"context"
	"math/big"
	"net/http"
	"net/url"
)

type Code struct {
	Bytecode []string `json:"bytecode"`
	Abi      []struct {
		Inputs []struct {
			Name string `json:"name"`
			Type string `json:"type"`
		} `json:"inputs"`
		Name            string        `json:"name"`
		Outputs         []interface{} `json:"outputs"`
		Type            string        `json:"type"`
		StateMutability string        `json:"stateMutability,omitempty"`
	} `json:"abi"`
}

// Gets a contracts code.
//
// [Reference](https://github.com/starkware-libs/cairo-lang/blob/fc97bdd8322a7df043c87c371634b26c15ed6cee/src/starkware/starknet/services/api/feeder_gateway/feeder_gateway_client.py#L55)
func (sg *StarknetGateway) CodeAt(ctx context.Context, contract string, blockNumber *big.Int) (*Code, error) {
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
