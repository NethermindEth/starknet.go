package caigo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"strings"
)

var (
	INVOKE *big.Int = big.NewInt(115923154332517) //"invoke"
	DEPLOY *big.Int = big.NewInt(110386840629113) //"deploy"
)

type StarkResp struct {
	Result []string `json:"result"`
}

type AddTxResponse struct {
	Code            string `json:"code"`
	TransactionHash string `json:"transaction_hash"`
}

type StarknetRequest struct {
	ContractAddress    string   `json:"contract_address"`
	EntryPointSelector string   `json:"entry_point_selector"`
	Calldata           []string `json:"calldata"`
	Signature          []string `json:"signature"`
	Type               string   `json:"type,omitempty"`
	Nonce              string   `json:"nonce,omitempty"`
}

// struct to catch starknet.js transaction payloads
type JSTransaction struct {
	Calldata           []string `json:"calldata"`
	ContractAddress    string   `json:"contract_address"`
	EntryPointSelector string   `json:"entry_point_selector"`
	EntryPointType     string   `json:"entry_point_type"`
	JSSignature        []string `json:"signature"`
	TransactionHash    string   `json:"transaction_hash"`
	Type               string   `json:"type"`
	Nonce              string   `json:"nonce"`
}

func (sg StarknetGateway) Call(sn StarknetRequest, blockId ...string) (resp []string, err error) {
	bid := ""
	if len(blockId) == 1 {
		bid = fmtBlockId(blockId[0])
	}

	url := fmt.Sprintf("%s/call_contract%s", sg.Feeder, strings.Replace(bid, "&", "?", 1))

	rawResp, err := sn.postHelper(url)
	if err != nil {
		return resp, err
	}

	var snResp StarkResp
	err = json.Unmarshal(rawResp, &snResp)
	return snResp.Result, err
}

func (sg StarknetGateway) Invoke(sn StarknetRequest) (addResp AddTxResponse, err error) {
	url := fmt.Sprintf("%s/add_transaction", sg.Gateway)

	rawResp, err := sn.postHelper(url)
	if err != nil {
		return addResp, err
	}

	err = json.Unmarshal(rawResp, &addResp)
	return addResp, err
}

func (sn StarknetRequest) postHelper(url string) (resp []byte, err error) {
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

	resp, err = ioutil.ReadAll(res.Body)
	if res.StatusCode >= 400 {
		return resp, fmt.Errorf("%s", string(resp))
	}

	return resp, err
}

func (jtx JSTransaction) ConvertTx() (tx Transaction) {
	tx = Transaction{
		ContractAddress:    jsToBN(jtx.ContractAddress),
		EntryPointSelector: jsToBN(jtx.EntryPointSelector),
		EntryPointType:     jtx.EntryPointType,
		TransactionHash:    jsToBN(jtx.TransactionHash),
		Type:               jtx.Type,
		Nonce:              jsToBN(jtx.Nonce),
	}
	for _, cd := range jtx.Calldata {
		tx.Calldata = append(tx.Calldata, jsToBN(cd))
	}
	for _, sigElem := range jtx.JSSignature {
		tx.Signature = append(tx.Signature, jsToBN(sigElem))
	}
	return tx
}

func jsToBN(str string) *big.Int {
	if strings.Contains(str, "0x") {
		return HexToBN(str)
	} else {
		return StrToBig(str)
	}
}

// Adheres to 'starknet.js' hash non typedData
func (sc StarkCurve) HashMsg(pubkey *big.Int, tx Transaction) (hash *big.Int, err error) {
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

func (sc StarkCurve) HashTx(prefix, chainId *big.Int, tx Transaction) (hash *big.Int, err error) {
	tx.Calldata = append(tx.Calldata, big.NewInt(int64(len(tx.Calldata))))
	cdHash, err := sc.HashElements(tx.Calldata)
	if err != nil {
		return hash, err
	}

	txHashData := []*big.Int{
		prefix,
		tx.ContractAddress,
		tx.EntryPointSelector,
		cdHash,
		// chainId,
		tx.Nonce,
	}

	txHashData = append(txHashData, big.NewInt(int64(len(txHashData))))
	hash, err = sc.HashElements(txHashData)
	return hash, err
}
