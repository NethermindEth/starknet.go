package gateway

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"net/url"

	"github.com/dontpanicdao/caigo/types"
	"github.com/google/go-querystring/query"
)

type Block struct {
	BlockHash           *types.Felt          `json:"block_hash"`
	ParentBlockHash     string               `json:"parent_block_hash"`
	BlockNumber         int                  `json:"block_number"`
	StateRoot           string               `json:"state_root"`
	Status              string               `json:"status"`
	Transactions        []Transaction        `json:"transactions"`
	Timestamp           int                  `json:"timestamp"`
	TransactionReceipts []TransactionReceipt `json:"transaction_receipts"`
}

func (b Block) Normalize() *types.Block {
	normalized := types.Block{
		BlockHash:       b.BlockHash,
		ParentBlockHash: b.ParentBlockHash,
		BlockNumber:     b.BlockNumber,
		NewRoot:         b.StateRoot,
		Status:          b.Status,
		AcceptedTime:    uint64(b.Timestamp),
	}

	for _, txn := range b.Transactions {
		normalized.Transactions = append(normalized.Transactions, txn.Normalize())
	}

	return &normalized
}

type BlockOptions struct {
	BlockNumber uint64      `url:"blockNumber,omitempty"`
	BlockHash   *types.Felt `url:"blockHash,omitempty"`
}

// Gets the block information from a block ID.
//
// [Reference](https://github.com/starkware-libs/cairo-lang/blob/f464ec4797361b6be8989e36e02ec690e74ef285/src/starkware/starknet/services/api/feeder_gateway/feeder_gateway_client.py#L27-L31)
func (sg *Gateway) Block(ctx context.Context, opts *BlockOptions) (*Block, error) {
	req, err := sg.newRequest(ctx, http.MethodGet, "/get_block", nil)
	if err != nil {
		return nil, err
	}

	// Example of an implementation change to map the real API
	type tempBlockOptions struct {
		BlockNumber uint64 `url:"blockNumber,omitempty"`
		BlockHash   string `url:"blockHash,omitempty"`
	}

	out := tempBlockOptions{}
	if opts != nil {
		out.BlockNumber = opts.BlockNumber
		if opts.BlockHash != nil && opts.BlockHash.Int != nil {
			out.BlockHash = fmt.Sprintf("0x%s", opts.BlockHash.Int.Text(16))
		}

		vs, err := query.Values(out)
		if err != nil {
			return nil, err
		}
		appendQueryValues(req, vs)
	}

	var resp Block
	return &resp, sg.do(req, &resp)
}

func (sg *Gateway) BlockHashByID(ctx context.Context, id uint64) (block string, err error) {
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

func (sg *Gateway) BlockIDByHash(ctx context.Context, hash string) (block uint64, err error) {
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

func (sg *Gateway) BlockByHash(context.Context, *types.Felt, string) (*types.Block, error) {
	panic("not implemented")
}

func (sg *Gateway) BlockByNumber(context.Context, *big.Int, string) (*types.Block, error) {
	panic("not implemented")
}
