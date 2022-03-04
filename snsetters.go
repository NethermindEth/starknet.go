package caigo

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"strings"
)

type StarkResp struct {
	Result []string `json:"result"`
}

type AddTxResponse struct {
	Code            string `json:"code"`
	TransactionHash string `json:"transaction_hash"`
}

type RawContractDefinition struct {
	ABI               []ABI                  `json:"abi"`
	EntryPointsByType EntryPointsByType      `json:"entry_points_by_type"`
	Program           map[string]interface{} `json:"program"`
}

type DeployRequest struct {
	Type                string   `json:"type"`
	ContractAddressSalt string   `json:"contract_address_salt"`
	ConstructorCalldata []string `json:"constructor_calldata"`
	ContractDefinition  struct {
		ABI               []ABI             `json:"abi"`
		EntryPointsByType EntryPointsByType `json:"entry_points_by_type"`
		Program           string            `json:"program"`
	} `json:"contract_definition"`
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

type EntryPointsByType struct {
	Constructor []struct {
		Offset   string `json:"offset"`
		Selector string `json:"selector"`
	} `json:"CONSTRUCTOR"`
	External []struct {
		Offset   string `json:"offset"`
		Selector string `json:"selector"`
	} `json:"EXTERNAL"`
	L1Handler []interface{} `json:"L1_HANDLER"`
}

func (sg StarknetGateway) Call(sn StarknetRequest, blockId ...string) (resp []string, err error) {
	bid := ""
	if len(blockId) == 1 {
		bid = fmtBlockId(blockId[0])
	}

	url := fmt.Sprintf("%s/call_contract%s", sg.Feeder, strings.Replace(bid, "&", "?", 1))

	pay, err := json.Marshal(sn)
	if err != nil {
		return resp, err
	}

	rawResp, err := postHelper(pay, url)
	if err != nil {
		return resp, err
	}

	var snResp StarkResp
	err = json.Unmarshal(rawResp, &snResp)
	return snResp.Result, err
}

func (sg StarknetGateway) Invoke(sn StarknetRequest) (addResp AddTxResponse, err error) {
	url := fmt.Sprintf("%s/add_transaction", sg.Gateway)

	pay, err := json.Marshal(sn)
	if err != nil {
		return addResp, err
	}

	rawResp, err := postHelper(pay, url)
	if err != nil {
		return addResp, err
	}

	err = json.Unmarshal(rawResp, &addResp)
	return addResp, err
}

// {
// 	"type": "INVOKE_FUNCTION",
// 	"contract_address": "0x0217d176acd37d6d456c433dd5246af96afc03d9f4d9241e815917ad81d639a1",
// 	"entry_point_selector": "0x240060cdb34fcc260f41eac7474ee1d7c80b7e3607daff9ac67c7ea2ebb1c44",
// 	"calldata": ["432476587373182848195563808259184426232857568687880284101194381213081650114", "216030643445273762074482936742625134427639679021380938148798651889117677069", "0", "26"],
//  "signature": ["1578457523021167749085824996732041607588313391654072972372699055154472126961", "1347141126837834920895919552157519075364938231173830585544694984729726519574"]
//   }
// func (sg StarknetGateway) Execute(address *big.Int, calldata []*big.Int) (addResp AddTxResponse, err error) {
// 	url := fmt.Sprintf("%s/add_transaction", sg.Gateway)

// 	rawResp, err := sn.postHelper(url)
// 	if err != nil {
// 		return addResp, err
// 	}

// 	err = json.Unmarshal(rawResp, &addResp)
// 	return addResp, err
// }

func (sg StarknetGateway) Deploy(filePath string, dr DeployRequest) (addResp AddTxResponse, err error) {
	url := fmt.Sprintf("%s/add_transaction", sg.Gateway)

	dat, err := os.ReadFile(filePath)
	if err != nil {
		return addResp, err
	}
	var rawDef RawContractDefinition
	err = json.Unmarshal(dat, &rawDef)
	if err != nil {
		return addResp, err
	}
	fmt.Println("RAW DEF: ", rawDef.ABI, rawDef.EntryPointsByType)

	dr.ContractDefinition.ABI = rawDef.ABI
	dr.ContractDefinition.EntryPointsByType = rawDef.EntryPointsByType
	dr.ContractDefinition.Program, err = CompressCompiledContract(rawDef.Program)
	if err != nil {
		return addResp, err
	}
	fmt.Println("PROG: ", dr.ContractDefinition.Program)

	pay, err := json.Marshal(dr)
	if err != nil {
		return addResp, err
	}

	rawResp, err := postHelper(pay, url)
	if err != nil {
		return addResp, err
	}

	err = json.Unmarshal(rawResp, &addResp)
	return addResp, err
}

func postHelper(pay []byte, url string) (resp []byte, err error) {
	method := "POST"

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

func CompressCompiledContract(program map[string]interface{}) (cc string, err error) {
	pay, err := json.Marshal(program)
	if err != nil {
		return cc, err
	}

	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	_, err = zw.Write(pay)
	if err != nil {
		return cc, err
	}
	if err := zw.Close(); err != nil {
		return cc, err
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
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
