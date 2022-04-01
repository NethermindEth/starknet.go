package caigo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

/*
	Instantiate a new StarkNet Gateway client
	- defaults to the GOERLI endpoints
*/
func NewGateway(chainId ...string) (sg StarknetGateway) {
	sg = StarknetGateway{
		Base:    GOERLI_BASE,
		Feeder:  GOERLI_BASE + "/feeder_gateway",
		Gateway: GOERLI_BASE + "/gateway",
		ChainId: GOERLI_ID,
	}
	if len(chainId) == 1 {
		if strings.Contains("main", strings.ToLower(chainId[0])) {
			sg = StarknetGateway{
				Base:    MAINNET_BASE,
				Feeder:  MAINNET_BASE + "/feeder_gateway",
				Gateway: MAINNET_BASE + "/gateway",
				ChainId: MAINNET_ID,
			}
		} else if strings.Contains("local", strings.ToLower(chainId[0])) || strings.Contains("dev", strings.ToLower(chainId[0])) {
			sg = StarknetGateway{
				Base:    LOCAL_BASE,
				Feeder:  LOCAL_BASE + "/feeder_gateway",
				Gateway: LOCAL_BASE + "/gateway",
				ChainId: GOERLI_ID,
			}
		} else {
			sg = StarknetGateway{
				Base:    chainId[0],
				Feeder:  chainId[0] + "/feeder_gateway",
				Gateway: chainId[0] + "/gateway",
				ChainId: GOERLI_ID,
			}
		}
	}
	return sg
}

func (sg StarknetGateway) GetBlockHashById(blockId string) (block string, err error) {
	url := fmt.Sprintf("%s/get_block_hash_by_id?blockId=%s", sg.Feeder, blockId)

	resp, err := getHelper(url)
	if err != nil {
		return block, err
	}

	return strings.Replace(string(resp), "\"", "", -1), nil
}

func (sg StarknetGateway) GetBlockIdByHash(blockHash string) (block string, err error) {
	url := fmt.Sprintf("%s/get_block_id_by_hash?blockHash=%s", sg.Feeder, blockHash)

	resp, err := getHelper(url)
	if err != nil {
		return block, err
	}

	return strings.Replace(string(resp), "\"", "", -1), nil
}

func (sg StarknetGateway) GetTransactionHashById(txId string) (tx string, err error) {
	url := fmt.Sprintf("%s/get_transaction_hash_by_id?transactionId=%s", sg.Feeder, txId)

	resp, err := getHelper(url)
	if err != nil {
		return tx, err
	}

	return strings.Replace(string(resp), "\"", "", -1), nil
}

func (sg StarknetGateway) GetTransactionIdByHash(txHash string) (tx string, err error) {
	url := fmt.Sprintf("%s/get_transaction_id_by_hash?transactionHash=%s", sg.Feeder, txHash)

	resp, err := getHelper(url)
	if err != nil {
		return tx, err
	}

	return strings.Replace(string(resp), "\"", "", -1), nil
}

func (sg StarknetGateway) GetStorageAt(contractAddress, key, blockId string) (storage string, err error) {
	url := fmt.Sprintf("%s/get_storage_at?contractAddress=%s&key=%s%s", sg.Feeder, contractAddress, key, fmtBlockId(blockId))

	resp, err := getHelper(url)
	if err != nil {
		return storage, err
	}

	return strings.Replace(string(resp), "\"", "", -1), nil
}

func (sg StarknetGateway) GetCode(contractAddress, blockId string) (code ContractCode, err error) {
	url := fmt.Sprintf("%s/get_code?contractAddress=%s%s", sg.Feeder, contractAddress, fmtBlockId(blockId))

	resp, err := getHelper(url)
	if err != nil {
		return code, err
	}

	err = json.Unmarshal(resp, &code)
	return code, err
}

func (sg StarknetGateway) GetBlock(blockId string) (block Block, err error) {
	bid := fmtBlockId(blockId)

	url := fmt.Sprintf("%s/get_block%s", sg.Feeder, strings.Replace(bid, "&", "?", 1))

	resp, err := getHelper(url)
	if err != nil {
		return block, err
	}

	err = json.Unmarshal(resp, &block)
	return block, err
}

func (sg StarknetGateway) GetTransactionStatus(txHash string) (status TransactionStatus, err error) {
	url := fmt.Sprintf("%s/get_transaction_status?transactionHash=%s", sg.Feeder, txHash)

	resp, err := getHelper(url)
	if err != nil {
		return status, err
	}

	err = json.Unmarshal(resp, &status)
	return status, err
}

func (sg StarknetGateway) GetTransaction(txHash string) (tx StarknetTransaction, err error) {
	url := fmt.Sprintf("%s/get_transaction?transactionHash=%s", sg.Feeder, txHash)

	resp, err := getHelper(url)
	if err != nil {
		return tx, err
	}

	err = json.Unmarshal(resp, &tx)
	return tx, err
}

func (sg StarknetGateway) GetTransactionReceipt(txHash string) (receipt TransactionReceipt, err error) {
	url := fmt.Sprintf("%s/get_transaction_receipt?transactionHash=%s", sg.Feeder, txHash)

	resp, err := getHelper(url)
	if err != nil {
		return receipt, err
	}

	err = json.Unmarshal(resp, &receipt)
	return receipt, err
}

func (sg StarknetGateway) PollTx(txHash string, threshold TxStatus, interval, maxPoll int) (n int, status string, err error) {
	err = fmt.Errorf("could find tx status for tx:  %s\n", txHash)

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	cow := 0
	for range ticker.C {
		if cow >= maxPoll {
			return cow, status, err
		}
		cow++

		stat, err := sg.GetTransactionStatus(txHash)
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

func fmtBlockId(blockId string) string {
	if len(blockId) < 2 {
		return ""
	}

	if blockId[:2] == "0x" {
		return fmt.Sprintf("&blockHash=%s", blockId)
	}
	return fmt.Sprintf("&blockNumber=%s", blockId)
}

func getHelper(url string) (resp []byte, err error) {
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return resp, err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return resp, err
	}
	defer res.Body.Close()

	resp, err = ioutil.ReadAll(res.Body)
	return resp, err
}
