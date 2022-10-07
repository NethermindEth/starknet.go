package gateway

import (
	"context"
	"math/big"
	"net/http"
	"net/url"

	"github.com/dontpanicdao/caigo/types"
)

// Gets a contracts code.
//
// [Reference](https://github.com/starkware-libs/cairo-lang/blob/fc97bdd8322a7df043c87c371634b26c15ed6cee/src/starkware/starknet/services/api/feeder_gateway/feeder_gateway_client.py#L55)
func (sg *Gateway) CodeAt(ctx context.Context, contract string, blockNumber *big.Int) (*types.Code, error) {
	req, err := sg.newRequest(ctx, http.MethodGet, "/get_code", nil)
	if err != nil {
		return nil, err
	}

	appendQueryValues(req, url.Values{"contractAddress": []string{contract}})

	if blockNumber != nil {
		appendQueryValues(req, url.Values{"blockNumber": []string{blockNumber.String()}})
	}

	var resp types.Code
	return &resp, sg.do(req, &resp)
}

func (sg *Gateway) FullContract(ctx context.Context, contract string) (*types.ContractClass, error) {
	req, err := sg.newRequest(ctx, http.MethodGet, "/get_full_contract", nil)
	if err != nil {
		return nil, err
	}

	appendQueryValues(req, url.Values{"contractAddress": []string{contract}})

	var resp types.ContractClass
	return &resp, sg.do(req, &resp)
}
