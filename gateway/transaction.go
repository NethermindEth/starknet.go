package gateway

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/dontpanicdao/caigo/types"
	"github.com/google/go-querystring/query"
)

type StarknetTransaction struct {
	TransactionIndex int         `json:"transaction_index"`
	BlockNumber      int         `json:"block_number"`
	Transaction      Transaction `json:"transaction"`
	BlockHash        *types.Felt `json:"block_hash"`
	Status           string      `json:"status"`
}

type Transaction struct {
	TransactionHash    *types.Felt   `json:"transaction_hash,omitempty"`
	ClassHash          string        `json:"class_hash,omitempty"`
	ContractAddress    *types.Felt   `json:"contract_address,omitempty"`
	SenderAddress      *types.Felt   `json:"sender_address,omitempty"`
	EntryPointSelector *types.Felt   `json:"entry_point_selector"`
	Calldata           []*types.Felt `json:"calldata"`
	Signature          []*types.Felt `json:"signature"`
	EntryPointType     *types.Felt   `json:"entry_point_type,omitempty"`
	MaxFee             *types.Felt   `json:"max_fee,omitempty"`
	Nonce              *types.Felt   `json:"nonce,omitempty"`
	Version            string        `json:"version,omitempty"`
	Type               string        `json:"type,omitempty"`
}

func (t Transaction) Normalize() *types.Transaction {
	return &types.Transaction{
		TransactionHash:    t.TransactionHash,
		ClassHash:          t.ClassHash,
		ContractAddress:    t.ContractAddress,
		SenderAddress:      t.SenderAddress,
		EntryPointSelector: t.EntryPointSelector,
		Calldata:           t.Calldata,
		Signature:          t.Signature,
		MaxFee:             t.MaxFee,
		Nonce:              t.Nonce,
		Version:            t.Version,
		Type:               t.Type,
	}
}

type TransactionReceipt struct {
	Status                string `json:"status"`
	BlockHash             string `json:"block_hash"`
	BlockNumber           int    `json:"block_number"`
	TransactionIndex      int    `json:"transaction_index"`
	TransactionHash       string `json:"transaction_hash"`
	L1ToL2ConsumedMessage struct {
		FromAddress string   `json:"from_address"`
		ToAddress   string   `json:"to_address"`
		Selector    string   `json:"selector"`
		Payload     []string `json:"payload"`
	} `json:"l1_to_l2_consumed_message"`
	L2ToL1Messages     []interface{}            `json:"l2_to_l1_messages"`
	Events             []interface{}            `json:"events"`
	ExecutionResources types.ExecutionResources `json:"execution_resources"`
}

type TransactionOptions struct {
	TransactionId   uint64 `url:"transactionId,omitempty"`
	TransactionHash string `url:"transactionHash,omitempty"`
}

func (gw *Gateway) TransactionByHash(ctx context.Context, hash string) (*types.Transaction, error) {
	t, err := gw.Transaction(ctx, TransactionOptions{TransactionHash: hash})
	if err != nil {
		return nil, err
	}

	return t.Transaction.Normalize(), nil
}

// Gets the transaction information from a tx id.
//
// [Reference](https://github.com/starkware-libs/cairo-lang/blob/f464ec4797361b6be8989e36e02ec690e74ef285/src/starkware/starknet/services/api/feeder_gateway/feeder_gateway_client.py#L54-L58)
func (gw *Gateway) Transaction(ctx context.Context, opts TransactionOptions) (*StarknetTransaction, error) {
	req, err := gw.newRequest(ctx, http.MethodGet, "/get_transaction", nil)
	if err != nil {
		return nil, err
	}
	vs, err := query.Values(opts)
	if err != nil {
		return nil, err
	}
	appendQueryValues(req, vs)

	var resp StarknetTransaction
	return &resp, gw.do(req, &resp)
}

type TransactionStatusOptions struct {
	TransactionId   uint64 `url:"transactionId,omitempty"`
	TransactionHash string `url:"transactionHash,omitempty"`
}

// Gets the transaction status from a txn.
//
// [Reference](https://github.com/starkware-libs/cairo-lang/blob/fc97bdd8322a7df043c87c371634b26c15ed6cee/src/starkware/starknet/services/api/feeder_gateway/feeder_gateway_client.py#L87)
func (gw *Gateway) TransactionStatus(ctx context.Context, opts TransactionStatusOptions) (*types.TransactionStatus, error) {
	req, err := gw.newRequest(ctx, http.MethodGet, "/get_transaction_status", nil)
	if err != nil {
		return nil, err
	}
	vs, err := query.Values(opts)
	if err != nil {
		return nil, err
	}
	appendQueryValues(req, vs)

	var resp types.TransactionStatus
	return &resp, gw.do(req, &resp)
}

// Gets the transaction id from its hash.
//
// [Reference](https://github.com/starkware-libs/cairo-lang/blob/fc97bdd8322a7df043c87c371634b26c15ed6cee/src/starkware/starknet/services/api/feeder_gateway/feeder_gateway_client.py#L137)
func (gw *Gateway) TransactionID(ctx context.Context, hash string) (*big.Int, error) {
	req, err := gw.newRequest(ctx, http.MethodGet, "/get_transaction_id_by_hash", nil)
	if err != nil {
		return nil, err
	}

	appendQueryValues(req, url.Values{
		"transactionHash": []string{hash},
	})

	var resp big.Int
	return &resp, gw.do(req, &resp)
}

// Gets the transaction hash from its id.
//
// [Reference](https://github.com/starkware-libs/cairo-lang/blob/fc97bdd8322a7df043c87c371634b26c15ed6cee/src/starkware/starknet/services/api/feeder_gateway/feeder_gateway_client.py#L130)
func (gw *Gateway) TransactionHash(ctx context.Context, id *big.Int) (string, error) {
	req, err := gw.newRequest(ctx, http.MethodGet, "/get_transaction_hash_by_id", nil)
	if err != nil {
		return "", err
	}

	appendQueryValues(req, url.Values{
		"transactionId": []string{id.String()},
	})

	var resp string
	if err := gw.do(req, &resp); err != nil {
		return "", err
	}

	return resp, nil
}

// Get transaction receipt for specific tx
//
// [Reference](https://github.com/starkware-libs/cairo-lang/blob/fc97bdd8322a7df043c87c371634b26c15ed6cee/src/starkware/starknet/services/api/feeder_gateway/feeder_gateway_client.py#L104)
func (gw *Gateway) TransactionReceipt(ctx context.Context, txHash string) (*types.TransactionReceipt, error) {
	req, err := gw.newRequest(ctx, http.MethodGet, "/get_transaction_receipt", nil)
	if err != nil {
		return nil, err
	}

	appendQueryValues(req, url.Values{
		"transactionHash": []string{txHash},
	})

	var resp types.TransactionReceipt
	return &resp, gw.do(req, &resp)
}

func (gw *Gateway) TransactionTrace(ctx context.Context, txHash string) (*types.TransactionTrace, error) {
	req, err := gw.newRequest(ctx, http.MethodGet, "/get_transaction_trace", nil)
	if err != nil {
		return nil, err
	}

	appendQueryValues(req, url.Values{
		"transactionHash": []string{txHash},
	})

	var resp types.TransactionTrace
	return &resp, gw.do(req, &resp)
}

// Long poll a transaction for specificed interval and max polls until the desired TxStatus has been achieved
// or the transaction reverts
func (gw *Gateway) PollTx(ctx context.Context, txHash string, threshold types.TxStatus, interval, maxPoll int) (n int, receipt *types.TransactionReceipt, err error) {
	err = fmt.Errorf("could not find tx status for tx:  %s", txHash)

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	cow := 0
	for range ticker.C {
		if cow >= maxPoll {
			return cow, receipt, err
		}
		cow++

		receipt, err = gw.TransactionReceipt(ctx, txHash)
		if err != nil {
			return cow, receipt, err
		}
		sInt := FindTxStatus(receipt.Status)
		if sInt == 1 {
			return cow, receipt, fmt.Errorf(receipt.StatusData)
		} else if sInt >= int(threshold) {
			return cow, receipt, nil
		}
	}
	return cow, receipt, err
}

func FindTxStatus(stat string) int {
	for i, val := range types.TxStatuses {
		if val == strings.ToUpper(stat) {
			return i
		}
	}
	return 0
}
