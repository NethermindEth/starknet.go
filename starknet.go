package caigo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"math/big"
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

type StarknetTransaction struct {
	TransactionIndex int `json:"transaction_index"`
	BlockNumber      int `json:"block_number"`
	Transaction      Transaction `json:"transaction"`
	BlockHash string `json:"block_hash"`
	Status    string `json:"status"`
}

// Starknet transaction composition
type Transaction struct {
	Calldata           []*big.Int `json:"calldata"`
	ContractAddress    *big.Int   `json:"contract_address"`
	EntryPointSelector *big.Int   `json:"entry_point_selector"`
	EntryPointType     string     `json:"entry_point_type"`
	Signature          []*big.Int `json:"signature"`
	TransactionHash    *big.Int   `json:"transaction_hash"`
	Type               string     `json:"type"`
	Nonce              *big.Int   `json:"nonce,omitempty"`
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

func GetTransaction(providerBaseUri, txHash string) (tx StarknetTransaction, err error) {
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


// Adheres to 'starknet.js' hash non typedData
func (sc StarkCurve) HashTx(pubkey *big.Int, tx Transaction) (hash *big.Int, err error) {
	tx.Calldata = append(tx.Calldata, big.NewInt(int64(len(tx.Calldata))))
	cdHash, err := sc.HashElements(tx.Calldata)
	if err != nil {
		return hash, err
	}

	txHashData := []*big.Int{
		pubkey,
		tx.ContractAddress,
		tx.EntryPointSelector,
		cdHash,
		tx.Nonce,
	}
	txHashData = append(txHashData, big.NewInt(int64(len(txHashData))))
	hash, err = sc.HashElements(txHashData)
	return hash, err
}

