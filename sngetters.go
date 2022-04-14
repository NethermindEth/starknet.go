package caigo

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"strings"
	"time"
)

/*
	Instantiate a new StarkNet Gateway client
	- defaults to the GOERLI endpoints
*/
func NewGateway(opts ...GatewayOption) *StarknetGateway {
	gopts := gatewayOptions{
		chainID: GOERLI_ID,
		client:  http.DefaultClient,
	}

	for _, opt := range opts {
		opt.apply(&gopts)
	}

	var sg *StarknetGateway
	switch id := strings.ToLower(gopts.chainID); {
	case strings.Contains("main", id):
		sg = &StarknetGateway{
			Base:    MAINNET_BASE,
			Feeder:  MAINNET_BASE + "/feeder_gateway",
			Gateway: MAINNET_BASE + "/gateway",
			ChainId: MAINNET_ID,
		}
	case strings.Contains("local", id):
		fallthrough
	case strings.Contains("dev", id):
		sg = &StarknetGateway{
			Base:    LOCAL_BASE,
			Feeder:  LOCAL_BASE + "/feeder_gateway",
			Gateway: LOCAL_BASE + "/gateway",
			ChainId: GOERLI_ID,
		}
	default:
		sg = &StarknetGateway{
			Base:    GOERLI_BASE,
			Feeder:  GOERLI_BASE + "/feeder_gateway",
			Gateway: GOERLI_BASE + "/gateway",
			ChainId: GOERLI_ID,
		}
	}

	sg.client = gopts.client

	return sg
}

func (sg *StarknetGateway) BlockHashById(ctx context.Context, blockId string) (block string, err error) {
	url := fmt.Sprintf("%s/get_block_hash_by_id?blockId=%s", sg.Feeder, blockId)

	resp, err := sg.getHelper(ctx, url)
	if err != nil {
		return block, err
	}

	return strings.Replace(string(resp), "\"", "", -1), nil
}

func (sg *StarknetGateway) BlockIdByHash(ctx context.Context, blockHash string) (block string, err error) {
	url := fmt.Sprintf("%s/get_block_id_by_hash?blockHash=%s", sg.Feeder, blockHash)

	resp, err := sg.getHelper(ctx, url)
	if err != nil {
		return block, err
	}

	return strings.Replace(string(resp), "\"", "", -1), nil
}

func (sg *StarknetGateway) TransactionHashById(ctx context.Context, txId string) (tx string, err error) {
	url := fmt.Sprintf("%s/get_transaction_hash_by_id?transactionId=%s", sg.Feeder, txId)

	resp, err := sg.getHelper(ctx, url)
	if err != nil {
		return tx, err
	}

	return strings.Replace(string(resp), "\"", "", -1), nil
}

func (sg *StarknetGateway) TransactionIdByHash(ctx context.Context, txHash string) (tx string, err error) {
	url := fmt.Sprintf("%s/get_transaction_id_by_hash?transactionHash=%s", sg.Feeder, txHash)

	resp, err := sg.getHelper(ctx, url)
	if err != nil {
		return tx, err
	}

	return strings.Replace(string(resp), "\"", "", -1), nil
}

func (sg *StarknetGateway) StorageAt(ctx context.Context, contractAddress, key, blockId string) (storage string, err error) {
	url := fmt.Sprintf("%s/get_storage_at?contractAddress=%s&key=%s%s", sg.Feeder, contractAddress, key, fmtBlockId(blockId))

	resp, err := sg.getHelper(ctx, url)
	if err != nil {
		return storage, err
	}

	return strings.Replace(string(resp), "\"", "", -1), nil
}

func (sg *StarknetGateway) Code(ctx context.Context, contractAddress, blockId string) (code ContractCode, err error) {
	url := fmt.Sprintf("%s/get_code?contractAddress=%s%s", sg.Feeder, contractAddress, fmtBlockId(blockId))

	resp, err := sg.getHelper(ctx, url)
	if err != nil {
		return code, err
	}

	err = json.Unmarshal(resp, &code)
	return code, err
}

func (sg *StarknetGateway) Block(ctx context.Context, blockId string) (block Block, err error) {
	bid := fmtBlockId(blockId)

	url := fmt.Sprintf("%s/get_block%s", sg.Feeder, strings.Replace(bid, "&", "?", 1))

	resp, err := sg.getHelper(ctx, url)
	if err != nil {
		return block, err
	}

	err = json.Unmarshal(resp, &block)
	return block, err
}

func (sg *StarknetGateway) TransactionStatus(ctx context.Context, txHash string) (status TransactionStatus, err error) {
	url := fmt.Sprintf("%s/get_transaction_status?transactionHash=%s", sg.Feeder, txHash)

	resp, err := sg.getHelper(ctx, url)
	if err != nil {
		return status, err
	}

	err = json.Unmarshal(resp, &status)
	return status, err
}

func (sg *StarknetGateway) Transaction(ctx context.Context, txHash string) (tx StarknetTransaction, err error) {
	url := fmt.Sprintf("%s/get_transaction?transactionHash=%s", sg.Feeder, txHash)

	resp, err := sg.getHelper(ctx, url)
	if err != nil {
		return tx, err
	}

	err = json.Unmarshal(resp, &tx)
	return tx, err
}

func (sg *StarknetGateway) TransactionReceipt(ctx context.Context, txHash string) (receipt TransactionReceipt, err error) {
	url := fmt.Sprintf("%s/get_transaction_receipt?transactionHash=%s", sg.Feeder, txHash)

	resp, err := sg.getHelper(ctx, url)
	if err != nil {
		return receipt, err
	}

	err = json.Unmarshal(resp, &receipt)
	return receipt, err
}

func (sg *StarknetGateway) PollTx(ctx context.Context, txHash string, threshold TxStatus, interval, maxPoll int) (n int, status string, err error) {
	err = fmt.Errorf("could find tx status for tx:  %s", txHash)

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	cow := 0
	for range ticker.C {
		if cow >= maxPoll {
			return cow, status, err
		}
		cow++

		stat, err := sg.TransactionStatus(ctx, txHash)
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

func (sg *StarknetGateway) AccountNonce(ctx context.Context, address *big.Int) (nonce *big.Int, err error) {
	resp, err := sg.Call(ctx, StarknetRequest{
		ContractAddress:    BigToHex(address),
		EntryPointSelector: BigToHex(GetSelectorFromName("get_nonce")),
	})
	if err != nil {
		return nonce, err
	}
	if len(resp) == 0 {
		return nonce, fmt.Errorf("no resp in contract call 'get_nonce' %v", BigToHex(address))
	}

	return HexToBN(resp[0]), nil
}

func fmtBlockId(blockId string) string {
	if len(blockId) < 2 {
		return ""
	}

	if blockId[:2] == "0x" {
		return fmt.Sprintf("&blockHash=%s", blockId)
	}
	return fmt.Sprintf("&blockNumber=%s", blockId)
}

func (sg *StarknetGateway) getHelper(ctx context.Context, url string) (resp []byte, err error) {
	method := "GET"

	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return resp, err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := sg.client.Do(req)
	if err != nil {
		return resp, err
	}
	defer res.Body.Close()

	resp, err = ioutil.ReadAll(res.Body)
	return resp, err
}
