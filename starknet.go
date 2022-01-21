package caigo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type StarkNetRequest struct {
	ContractAddress    string   `json:"contract_address"`
	EntryPointSelector string   `json:"entry_point_selector"`
	Calldata           []string `json:"calldata"`
	Signature          []string `json:"signature"`
}

type StarkResp struct {
	Result []string `json:"result"`
}

type TransactionStatus struct {
	TxStatus  string `json:"tx_status"`
	BlockHash string `json:"block_hash"`
}

type Transaction struct {
	TransactionIndex int `json:"transaction_index"`
	BlockNumber      int `json:"block_number"`
	Transaction      struct {
		Signature          []string `json:"signature"`
		EntryPointType     string   `json:"entry_point_type"`
		TransactionHash    string   `json:"transaction_hash"`
		Calldata           []string `json:"calldata"`
		EntryPointSelector string   `json:"entry_point_selector"`
		ContractAddress    string   `json:"contract_address"`
		Type               string   `json:"type"`
	} `json:"transaction"`
	BlockHash string `json:"block_hash"`
	Status    string `json:"status"`
}

func (sn StarkNetRequest) Call(providerBaseUri string) (resp []string, err error) {
	url := fmt.Sprintf("%s/call_contract", providerBaseUri)

	method := "POST"

	pay, err := json.Marshal(sn)
	if err != nil {
		return resp, err
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(pay))
	if err != nil {
		return resp, err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return resp, err
	}
	defer res.Body.Close()

	var sRes StarkResp
	json.NewDecoder(res.Body).Decode(&sRes)
	return sRes.Result, nil
}

func GetTransactionStatus(providerBaseUri, txHash string) (status TransactionStatus, err error) {
	url := fmt.Sprintf("%s/get_transaction_status?transactionHash=%s", providerBaseUri, txHash)

	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return status, err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return status, err
	}
	defer res.Body.Close()

	json.NewDecoder(res.Body).Decode(&status)
	return status, nil
}

func GetTransaction(providerBaseUri, txHash string) (tx Transaction, err error) {
	url := fmt.Sprintf("%s/get_transaction?transactionHash=%s", providerBaseUri, txHash)
	
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return tx, err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return tx, err
	}
	defer res.Body.Close()

	json.NewDecoder(res.Body).Decode(&tx)
	return tx, nil
}
