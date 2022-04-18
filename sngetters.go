package caigo

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"strings"
)

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

func (sg *StarknetGateway) AccountNonce(ctx context.Context, address string) (nonce *big.Int, err error) {
	resp, err := sg.Call(ctx, Transaction{
		ContractAddress:    address,
		EntryPointSelector: "get_nonce",
	})
	if err != nil {
		return nonce, err
	}
	if len(resp) == 0 {
		return nonce, fmt.Errorf("no resp in contract call 'get_nonce' %v", address)
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
