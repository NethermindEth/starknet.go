package gateway

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strings"
)

type Block struct {
	BlockHash           string               `json:"block_hash"`
	ParentBlockHash     string               `json:"parent_block_hash"`
	BlockNumber         int                  `json:"block_number"`
	StateRoot           string               `json:"state_root"`
	Status              string               `json:"status"`
	Transactions        []Transaction        `json:"transactions"`
	Timestamp           int                  `json:"timestamp"`
	TransactionReceipts []TransactionReceipt `json:"transaction_receipts"`
}

type BlockOptions struct {
	BlockNumber *uint64
	BlockHash   string
	Tag         string
}

var ErrInvalidBlock = errors.New("invalid block")

func (b *BlockOptions) appendQueryValues(req *http.Request) error {
	q := req.URL.Query()
	if b.Tag == "pending" || b.Tag == "latest" {
		q.Add("blockNumber", b.Tag)
		req.URL.RawQuery = q.Encode()
		return nil
	}

	if b.Tag != "" {
		return ErrInvalidBlock
	}

	if b.BlockNumber != nil {
		q.Add("blockNumber", fmt.Sprintf("%d", *b.BlockNumber))
		req.URL.RawQuery = q.Encode()
		return nil
	}

	if b.BlockHash != "" {
		if !strings.HasPrefix(b.BlockHash, "0x") {
			return ErrInvalidBlock
		}
		_, ok := big.NewInt(0).SetString(b.BlockHash, 0)
		if !ok {
			return ErrInvalidBlock
		}
		q.Add("blockHash", b.BlockHash)
		req.URL.RawQuery = q.Encode()
		return nil
	}
	q.Add("blockNumber", "pending")
	req.URL.RawQuery = q.Encode()
	return nil
}

// Gets the block information from a block ID.
//
// [Reference](https://github.com/starkware-libs/cairo-lang/blob/f464ec4797361b6be8989e36e02ec690e74ef285/src/starkware/starknet/services/api/feeder_gateway/feeder_gateway_client.py#L27-L31)
func (sg *Gateway) Block(ctx context.Context, opts *BlockOptions) (*Block, error) {
	req, err := sg.newRequest(ctx, http.MethodGet, "/get_block", nil)
	if err != nil {
		return nil, err
	}

	if err := opts.appendQueryValues(req); err != nil {
		return nil, err
	}

	var resp Block
	return &resp, sg.do(req, &resp)
}

func (sg *Gateway) BlockByHash(context.Context, string, string) (*Block, error) {
	panic("not implemented")
}

func (sg *Gateway) BlockByNumber(context.Context, *big.Int, string) (*Block, error) {
	panic("not implemented")
}
