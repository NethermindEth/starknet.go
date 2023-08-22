package gateway

import (
	"context"
	"net/http"
	"net/url"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/rpc"
)

// TODO: returns DeprecatedContractClass | SierraContractClass
func (sg *Gateway) ClassByHash(ctx context.Context, hash string) (*rpc.DepcreatedContractClass, error) {
	req, err := sg.newRequest(ctx, http.MethodGet, "/get_class_by_hash", nil)
	if err != nil {
		return nil, err
	}

	appendQueryValues(req, url.Values{
		"classHash": []string{hash},
	})

	var resp rpc.DepcreatedContractClass
	return &resp, sg.do(req, &resp)
}

func (sg *Gateway) ClassHashAt(ctx context.Context, address string) (*felt.Felt, error) {
	req, err := sg.newRequest(ctx, http.MethodGet, "/get_class_hash_at", nil)
	if err != nil {
		return nil, err
	}

	appendQueryValues(req, url.Values{
		"contractAddress": []string{address},
	})

	var resp *felt.Felt
	if err = sg.do(req, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (sg *Gateway) Class(context.Context, string) (*rpc.DepcreatedContractClass, error) {
	panic("not implemented")
}

func (sg *Gateway) ClassAt(context.Context, string) (*rpc.DepcreatedContractClass, error) {
	panic("not implemented")
}
