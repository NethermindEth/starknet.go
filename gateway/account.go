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
	resp, err := sg.Call(ctx, types.Transaction{
		ContractAddress:    address,
		EntryPointSelector: "get_nonce",
	}, nil)
	if err != nil {
		return nil, err
	}
	if len(resp) == 0 {
		return nil, fmt.Errorf("no resp in contract call 'get_nonce' %v", address)
	}

	return caigo.HexToBN(resp[0]), nil
}

func (sg *Gateway) EstimateFee(ctx context.Context, tx types.Transaction) (*types.FeeEstimate, error) {
	tx.EntryPointSelector = caigo.BigToHex(caigo.GetSelectorFromName(tx.EntryPointSelector))
	req, err := sg.newRequest(ctx, http.MethodPost, "/estimate_fee", tx)
	if err != nil {
		return nil, err
	}

	var resp types.FeeEstimate
	return &resp, sg.do(req, &resp)
}
