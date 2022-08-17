package gateway

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/dontpanicdao/caigo/felt"
	"github.com/dontpanicdao/caigo/types"
)

func (sg *Gateway) AccountNonce(ctx context.Context, address felt.Felt) (*felt.Felt, error) {
	selector := felt.GetSelectorFromName("get_nonce")
	resp, err := sg.Call(ctx, types.FunctionCall{
		ContractAddress:    address,
		EntryPointSelector: &selector,
	}, "")
	if err != nil {
		return nil, err
	}
	if len(resp) == 0 {
		return nil, fmt.Errorf("no resp in contract call 'get_nonce' %v", address)
	}
	output, err := felt.TextToFelt(resp[0])
	return output, err
}

func (sg *Gateway) EstimateFee(ctx context.Context, call types.FunctionInvoke, hash string) (*types.FeeEstimate, error) {
	req, err := sg.newRequest(ctx, http.MethodPost, "/estimate_fee", call)
	if err != nil {
		return nil, err
	}

	if hash != "" {
		appendQueryValues(req, url.Values{
			"blockHash": []string{hash},
		})
	}

	var resp types.FeeEstimate
	return &resp, sg.do(req, &resp)
}
