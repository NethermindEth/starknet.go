package gateway

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dontpanicdao/caigo/felt"
	"github.com/dontpanicdao/caigo/types"
	"github.com/google/go-querystring/query"
)

type Block struct {
	BlockHash           felt.Felt            `json:"block_hash"`
	ParentBlockHash     *felt.Felt           `json:"parent_block_hash"`
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
		normalizedTxn := txn.Normalize()
		normalized.Transactions = append(normalized.Transactions, *normalizedTxn)
	}

	return &normalized
}

type BlockOptions struct {
	BlockNumber uint64    `url:"blockNumber,omitempty"`
	BlockHash   felt.Felt `url:"blockHash,omitempty"`
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
	if opts != nil {
		out := tempBlockOptions{
			BlockNumber: opts.BlockNumber,
		}
		if !opts.BlockHash.IsNil() {
			out.BlockHash = opts.BlockHash.Hex()
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

func (sg *Gateway) BlockHashByID(ctx context.Context, id uint64) (block *felt.Felt, err error) {
	req, err := sg.newRequest(ctx, http.MethodGet, "/get_block_hash_by_id", nil)
	if err != nil {
		return nil, err
	}

	out := tempBlockOptions{
		BlockNumber: id,
	}

	out.BlockNumber = id

	vs, err := query.Values(out)
	if err != nil {
		return nil, err
	}
	appendQueryValues(req, vs)

	var resp *felt.Felt
	return resp, sg.do(req, &resp)
}

func (sg *Gateway) BlockIDByHash(ctx context.Context, hash felt.Felt) (block uint64, err error) {
	req, err := sg.newRequest(ctx, http.MethodGet, "/get_block_id_by_hash", nil)
	if err != nil {
		return 0, err
	}

	out := tempBlockOptions{}
	if !hash.IsNil() {
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

func (sg *Gateway) BlockByHash(context.Context, felt.Felt, string) (*types.Block, error) {
	panic("not implemented")
}

func (sg *Gateway) BlockByNumber(context.Context, uint64, string) (*types.Block, error) {
	panic("not implemented")
}
