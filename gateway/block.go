package gateway

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/dontpanicdao/caigo"
	"github.com/google/go-querystring/query"
)

type BlockOptions struct {
	BlockNumber uint64 `url:"blockNumber,omitempty"`
	BlockHash   string `url:"blockHash,omitempty"`
}

// Gets the block information from a block ID.
//
// [Reference](https://github.com/starkware-libs/cairo-lang/blob/f464ec4797361b6be8989e36e02ec690e74ef285/src/starkware/starknet/services/api/feeder_gateway/feeder_gateway_client.py#L27-L31)
func (sg *StarknetGateway) Block(ctx context.Context, opts *BlockOptions) (*caigo.Block, error) {
	req, err := sg.newRequest(ctx, http.MethodGet, "/get_block", nil)
	if err != nil {
		return nil, err
	}
	if opts != nil {
		vs, err := query.Values(opts)
		if err != nil {
			return nil, err
		}
		appendQueryValues(req, vs)
	}

	var resp caigo.Block
	return &resp, sg.do(req, &resp)
}

func (sg *StarknetGateway) BlockHashByID(ctx context.Context, id uint64) (block string, err error) {
	req, err := sg.newRequest(ctx, http.MethodGet, "/get_block_hash_by_id", nil)
	if err != nil {
		return "", err
	}

	appendQueryValues(req, url.Values{
		"blockId": []string{fmt.Sprint(id)},
	})

	var resp string
	return resp, sg.do(req, &resp)
}

func (sg *StarknetGateway) BlockIDByHash(ctx context.Context, hash string) (block uint64, err error) {
	req, err := sg.newRequest(ctx, http.MethodGet, "/get_block_id_by_hash", nil)
	if err != nil {
		return 0, err
	}

	appendQueryValues(req, url.Values{
		"blockHash": []string{hash},
	})

	var resp uint64
	return resp, sg.do(req, &resp)
}

func (sg *StarknetGateway) BlockByHash(ctx context.Context, hash string) (*caigo.Block, error) {
	return sg.Block(ctx, &BlockOptions{BlockHash: hash})
}

func (sg *StarknetGateway) BlockByNumber(ctx context.Context, number uint64) (*caigo.Block, error) {
	return sg.Block(ctx, &BlockOptions{BlockNumber: number})
}
