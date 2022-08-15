package gateway

import (
	"context"
	"fmt"
	"math/big"
	"net/http"

	"github.com/dontpanicdao/caigo/types"
	"github.com/google/go-querystring/query"
)

type Block struct {
	BlockHash           types.Felt           `json:"block_hash"`
	ParentBlockHash     *types.Felt          `json:"parent_block_hash"`
	BlockNumber         uint64               `json:"block_number"`
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

type tempBlockOptions struct {
	BlockNumber uint64 `url:"blockId,omitempty"`
	BlockHash   string `url:"blockHash,omitempty"`
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

func (sg *Gateway) BlockHashByID(ctx context.Context, id uint64) (block *types.Felt, err error) {
	req, err := sg.newRequest(ctx, http.MethodGet, "/get_block_hash_by_id", nil)
	if err != nil {
		return nil, err
	}

	out := tempBlockOptions{}

	out.BlockNumber = id

	vs, err := query.Values(out)
	if err != nil {
		return nil, err
	}
	appendQueryValues(req, vs)

	var resp *types.Felt
	return resp, sg.do(req, &resp)
}

func (sg *Gateway) BlockIDByHash(ctx context.Context, hash *types.Felt) (block uint64, err error) {
	req, err := sg.newRequest(ctx, http.MethodGet, "/get_block_id_by_hash", nil)
	if err != nil {
		return 0, err
	}

	out := tempBlockOptions{}
	if hash != nil && hash.Int != nil {
		out.BlockHash = fmt.Sprintf("0x%s", hash.Int.Text(16))
	}

	vs, err := query.Values(out)
	if err != nil {
		return 0, err
	}
	appendQueryValues(req, vs)

	var resp uint64
	return resp, sg.do(req, &resp)
}

func (sg *Gateway) BlockByHash(context.Context, *types.Felt, string) (*types.Block, error) {
	panic("not implemented")
}

func (sg *Gateway) BlockByNumber(context.Context, *big.Int, string) (*types.Block, error) {
	panic("not implemented")
}
