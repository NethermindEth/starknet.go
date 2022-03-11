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
	"strconv"
)

const PREFIX_TRANSACTION = "StarkNet Transaction"

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

type Signer struct {
	private *big.Int
	Curve StarkCurve
	Gateway StarknetGateway
	PublicX *big.Int
	PublicY *big.Int
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

func (sc StarkCurve) NewSigner(private, pubX, pubY *big.Int, chainId ...string) (signer Signer, err error) {
	if len(sc.ConstantPoints) == 0 {
		return signer, fmt.Errorf("must initiate precomputed constant points")
	}
	var gw StarknetGateway
	if len(chainId) == 1 {
		gw = NewGateway(chainId[0])
	} else {
		gw = NewGateway()
	}

	return Signer{
		private: private,
		Curve: sc,
		Gateway: gw,
		PublicX: pubX,
		PublicY: pubY,
	}, nil
}

func (sg StarknetGateway) Call(sn StarknetRequest, blockId ...string) (resp []string, err error) {
	bid := ""
	if len(blockId) == 1 {
		bid = fmtBlockId(blockId[0])
	}

	url := fmt.Sprintf("%s/call_contract%s", sg.Feeder, strings.Replace(bid, "&", "?", 1))

	if len(sn.Calldata) == 0 {
		sn.Calldata = []string{}
	}
	if len(sn.Signature) == 0 {
		sn.Signature = []string{}
	}

	pay, err := json.Marshal(sn)
	if err != nil {
		return resp, err
	}
	fmt.Println("PAY: ", string(pay))

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
	fmt.Println("DIS: ", string(pay))

	rawResp, err := postHelper(pay, url)
	if err != nil {
		return addResp, err
	}

	err = json.Unmarshal(rawResp, &addResp)
	return addResp, err
}

func (signer Signer) Execute(address *big.Int, txs []Transaction) (addResp AddTxResponse, err error) {
	nonce, err := signer.Gateway.GetAccountNonce(address)
	if err != nil {
		return addResp, err
	}

	hash, err := signer.Curve.HashMulticall(address, nonce, big.NewInt(0), big.NewInt(0), txs)
	if err != nil {
		return addResp, err
	}
	exp, _ := new(big.Int).SetString("303039478180100935883132162839533500076307020484752911895848045150679512896", 10)
	fmt.Println("GOT: ", hash)
	fmt.Println("EXP: ", exp)


	r, s, err := signer.Curve.Sign(hash, signer.private)
	if err != nil {
		return addResp, err
	}

	req := StarknetRequest{
		Type:               "INVOKE_FUNCTION",
		ContractAddress:    BigToHex(address),
		EntryPointSelector: BigToHex(GetSelectorFromName("__execute__")),
		Calldata:           FmtExecuteCalldata(nonce, txs),
		Signature:          []string{r.String(), s.String()},
	}

	return signer.Gateway.Invoke(req)
}

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

	dr.ContractDefinition.ABI = rawDef.ABI
	dr.ContractDefinition.EntryPointsByType = rawDef.EntryPointsByType
	dr.ContractDefinition.Program, err = CompressCompiledContract(rawDef.Program)
	if err != nil {
		return addResp, err
	}

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

func (sg StarknetGateway) GetAccountNonce(address *big.Int) (nonce *big.Int, err error) {
	resp, err := sg.Call(StarknetRequest{
		ContractAddress: BigToHex(address),
		EntryPointSelector: BigToHex(GetSelectorFromName("get_nonce")),
	})
	if err != nil {
		return nonce, err
	}
	if len(resp) == 0 {
		return nonce, fmt.Errorf("no resp in contract call 'get_nonce' %v\n", BigToHex(address))
	}

	return HexToBN(resp[0]), nil
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

func FmtExecuteCalldata(nonce *big.Int, txs []Transaction) (calldata []string) {
	var callArray, calldataArray []string
	calldata = append(calldata, strconv.Itoa(len(txs)))

	for i, tx := range txs {
		callArray = append(callArray, tx.ContractAddress.String(), tx.EntryPointSelector.String(), strconv.Itoa(i), strconv.Itoa(len(tx.Calldata)))
		if len(tx.Calldata) == 0 {
			calldataArray = append(calldataArray, "0")
		} else {
			for _, val :=  range tx.Calldata {
				calldataArray = append(calldataArray, val.String())
			}
		}
	}

	calldata = append(calldata, callArray...)
	calldata = append(calldata, strconv.Itoa(len(calldataArray)))
	calldata = append(calldata, calldataArray...)
	calldata = append(calldata, nonce.String())
	return calldata 
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

func (sc StarkCurve) HashMulticall(addr, nonce, max_fee, version *big.Int, txs []Transaction) (hash *big.Int, err error) {
	for _, tx := range txs {
		hash, err = sc.HashTx(addr, tx)
	}

	multiHashData := []*big.Int{
		UTF8StrToBig(PREFIX_TRANSACTION),
		addr,
		new(big.Int).Set(hash),
		nonce,
		max_fee,
		version,
	}

	multiHashData = append(multiHashData, big.NewInt(int64(len(multiHashData))))	
	hash, err = sc.HashElements(multiHashData)
	return hash, err
}

// Adheres to 'starknet.js' hash non typedData
func (sc StarkCurve) HashMsg(addr *big.Int, tx Transaction) (hash *big.Int, err error) {
	tx.Calldata = append(tx.Calldata, big.NewInt(int64(len(tx.Calldata))))
	cdHash, err := sc.HashElements(tx.Calldata)
	if err != nil {
		return hash, err
	}

	txHashData := []*big.Int{
		addr,
		tx.ContractAddress,
		tx.EntryPointSelector,
		cdHash,
		tx.Nonce,
	}

	txHashData = append(txHashData, big.NewInt(int64(len(txHashData))))	
	hash, err = sc.HashElements(txHashData)
	return hash, err
}

// Adheres to 'starknet.js' hash non typedData
func (sc StarkCurve) HashTx(addr *big.Int, tx Transaction) (hash *big.Int, err error) {
	tx.Calldata = append(tx.Calldata, big.NewInt(int64(len(tx.Calldata))))
	cdHash, err := sc.HashElements(tx.Calldata)
	if err != nil {
		return hash, err
	}

	txHashData := []*big.Int{
		tx.ContractAddress,
		tx.EntryPointSelector,
		cdHash,
	}

	txHashData = append(txHashData, big.NewInt(int64(len(txHashData))))	
	hash, err = sc.HashElements(txHashData)
	return hash, err
}
