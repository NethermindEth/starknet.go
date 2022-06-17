package gateway

import (
	"context"
	"net/http"
	"net/url"

	"github.com/google/go-querystring/query"
)

func (sg *Gateway) ClassByHash(ctx context.Context, hash string) (*RawContractDefinition, error) {
	req, err := sg.newRequest(ctx, http.MethodGet, "/get_class_by_hash", nil)
	if err != nil {
		return nil, err
	}

	appendQueryValues(req, url.Values{
		"classHash": []string{hash},
	})

	var resp RawContractDefinition
	return &resp, sg.do(req, &resp)
}

func (sg *Gateway) ClassHashAt(ctx context.Context, address string, opts *BlockOptions) (string, error) {
	req, err := sg.newRequest(ctx, http.MethodGet, "/get_class_hash_at", nil)
	if err != nil {
		return "", err
	}

	appendQueryValues(req, url.Values{
		"contractAddress": []string{address},
	})

	if opts != nil {
		vs, err := query.Values(opts)
		if err != nil {
			return "", err
		}
		appendQueryValues(req, vs)
	}

	var resp string
	return resp, sg.do(req, &resp)
}
