package gateway

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"net/url"

	"github.com/dontpanicdao/caigo"
	"github.com/dontpanicdao/caigo/types"
)

func (sg *Gateway) AccountNonce(ctx context.Context, address string) (*big.Int, error) {
	resp, err := sg.Call(ctx, types.FunctionCall{
		ContractAddress:    types.StrToFelt(address),
		EntryPointSelector: "get_nonce",
	}, "")
	if err != nil {
		return nil, err
	}
	if len(resp) == 0 {
		return nil, fmt.Errorf("no resp in contract call 'get_nonce' %v", address)
	}

	return caigo.HexToBN(resp[0]), nil
}

func (sg *Gateway) EstimateFee(ctx context.Context, call types.FunctionInvoke, hash string) (*types.FeeEstimate, error) {
	call.EntryPointSelector = caigo.BigToHex(caigo.GetSelectorFromName(call.EntryPointSelector))

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
