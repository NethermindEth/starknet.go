package gateway

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/dontpanicdao/caigo"
	"github.com/google/go-querystring/query"
)

type TransactionOptions struct {
	TransactionId   uint64 `url:"transactionId,omitempty"`
	TransactionHash string `url:"transactionHash,omitempty"`
}

// Gets the transaction information from a tx id.
//
// [Reference](https://github.com/starkware-libs/cairo-lang/blob/f464ec4797361b6be8989e36e02ec690e74ef285/src/starkware/starknet/services/api/feeder_gateway/feeder_gateway_client.py#L54-L58)
func (gw *StarknetGateway) Transaction(ctx context.Context, opts TransactionOptions) (*caigo.StarknetTransaction, error) {
	req, err := gw.newRequest(ctx, http.MethodGet, "/get_transaction", nil)
	if err != nil {
		return nil, err
	}
	vs, err := query.Values(opts)
	if err != nil {
		return nil, err
	}
	appendQueryValues(req, vs)

	var resp caigo.StarknetTransaction
	return &resp, gw.do(req, &resp)
}

type TransactionStatusOptions struct {
	TransactionId   uint64 `url:"transactionId,omitempty"`
	TransactionHash string `url:"transactionHash,omitempty"`
}

// Gets the transaction status from a txn.
//
// [Reference](https://github.com/starkware-libs/cairo-lang/blob/fc97bdd8322a7df043c87c371634b26c15ed6cee/src/starkware/starknet/services/api/feeder_gateway/feeder_gateway_client.py#L87)
func (gw *StarknetGateway) TransactionStatus(ctx context.Context, opts TransactionStatusOptions) (*caigo.TransactionStatus, error) {
	req, err := gw.newRequest(ctx, http.MethodGet, "/get_transaction_status", nil)
	if err != nil {
		return nil, err
	}
	vs, err := query.Values(opts)
	if err != nil {
		return nil, err
	}
	appendQueryValues(req, vs)

	var resp caigo.TransactionStatus
	return &resp, gw.do(req, &resp)
}

// Gets the transaction id from its hash.
//
// [Reference](https://github.com/starkware-libs/cairo-lang/blob/fc97bdd8322a7df043c87c371634b26c15ed6cee/src/starkware/starknet/services/api/feeder_gateway/feeder_gateway_client.py#L137)
func (gw *StarknetGateway) TransactionID(ctx context.Context, hash string) (*big.Int, error) {
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
func (gw *StarknetGateway) TransactionHash(ctx context.Context, id *big.Int) (string, error) {
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
func (gw *StarknetGateway) TransactionReceipt(ctx context.Context, txHash string) (*caigo.TransactionReceipt, error) {
	req, err := gw.newRequest(ctx, http.MethodGet, "/get_transaction_receipt", nil)
	if err != nil {
		return nil, err
	}

	appendQueryValues(req, url.Values{
		"transactionHash": []string{txHash},
	})

	var resp caigo.TransactionReceipt
	return &resp, gw.do(req, &resp)
}

// Long poll a transaction for specificed interval and max polls until the desired TxStatus has been achieved
// or the transaction reverts
func (gw *StarknetGateway) PollTx(ctx context.Context, txHash string, threshold caigo.TxStatus, interval, maxPoll int) (n int, status string, err error) {
	err = fmt.Errorf("could not find tx status for tx:  %s", txHash)

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	cow := 0
	for range ticker.C {
		if cow >= maxPoll {
			return cow, status, err
		}
		cow++

		stat, err := gw.TransactionStatus(ctx, TransactionStatusOptions{
			TransactionHash: txHash,
		})
		if err != nil {
			return cow, status, err
		}
		sInt := FindTxStatus(stat.TxStatus)
		if sInt == 1 {
			return cow, status, fmt.Errorf(stat.TxFailureReason.ErrorMessage)
		} else if sInt >= int(threshold) {
			return cow, stat.TxStatus, nil
		}
	}
	return cow, status, err
}

func FindTxStatus(stat string) int {
	for i, val := range caigo.TxStatuses {
		if val == strings.ToUpper(stat) {
			return i
		}
	}
	return 0
}
