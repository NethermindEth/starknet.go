package gateway

import (
	"context"
	"fmt"
	"math/big"
	"net/http"

	"github.com/dontpanicdao/caigo"
	"github.com/dontpanicdao/caigo/types"
)

func (sg *Gateway) AccountNonce(ctx context.Context, address string) (*big.Int, error) {
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
	b, _ := big.NewInt(0).SetString(resp[0], 0)
	return b, nil
}

func (sg *Gateway) EstimateFee(ctx context.Context, tx types.Transaction) (fee caigo.FeeEstimate, err error) {
	req, err := sg.newRequest(ctx, http.MethodPost, "/estimate_fee", tx)
	if err != nil {
		return fee, err
	}

	return fee, sg.do(req, &fee)
}
