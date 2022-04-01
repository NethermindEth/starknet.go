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

type FeeEstimate struct {
	Amount *big.Int `json:"amount"`
	Unit   string   `json:"unit"`
}

type RawContractDefinition struct {
	ABI               []ABI                  `json:"abi"`
	EntryPointsByType EntryPointsByType      `json:"entry_points_by_type"`
	Program           map[string]interface{} `json:"program"`
}

type Signer struct {
	private *big.Int
	Curve   StarkCurve
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
		Curve:   sc,
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

	sn.Type = INVOKE
	if len(sn.Calldata) == 0 {
		sn.Calldata = []string{}
	}
	if len(sn.Signature) == 0 {
		sn.Signature = []string{}
	}

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

func (signer Signer) Execute(address *big.Int, txs []Transaction) (addResp AddTxResponse, err error) {
	nonce, err := signer.Gateway.GetAccountNonce(address)
	if err != nil {
		return addResp, err
	}

	maxFee := big.NewInt(0)
	// for _, tx := range txs {
	// 	var cdStrings []string
	// 	for _, data := range tx.Calldata {
	// 		cdStrings = append(cdStrings, data.String())
	// 	}

	// 	innerFee, err := signer.Gateway.EstimateFee(StarknetRequest{
	// 		ContractAddress: BigToHex(tx.ContractAddress),
	// 		EntryPointSelector: BigToHex(tx.EntryPointSelector),
	// 		Calldata: cdStrings,
	// 		Signature: []string{},
	// 	})
	// 	if err == nil {
	// 		maxFee = maxFee.Add(maxFee, innerFee.Amount)
	// 	}
	// }

	hash, err := signer.Curve.HashMulticall(address, nonce, maxFee, UTF8StrToBig(signer.Gateway.ChainId), txs)
	if err != nil {
		return addResp, err
	}

	r, s, err := signer.Curve.Sign(hash, signer.private)
	if err != nil {
		return addResp, err
	}

	req := StarknetRequest{
		ContractAddress:    BigToHex(address),
		EntryPointSelector: BigToHex(GetSelectorFromName(EXECUTE_SELECTOR)),
		Calldata:           FmtExecuteCalldataStrings(nonce, txs),
		Signature:          []string{r.String(), s.String()},
	}

	return signer.Gateway.Invoke(req)
}

func (sg StarknetGateway) Deploy(filePath string, deployRequest DeployRequest) (addResp AddTxResponse, err error) {
	url := fmt.Sprintf("%s/add_transaction", sg.Gateway)

	dat, err := os.ReadFile(filePath)
	if err != nil {
		return addResp, err
	}

	deployRequest.Type = DEPLOY
	if len(deployRequest.ConstructorCalldata) == 0 {
		deployRequest.ConstructorCalldata = []string{}
	}

	var rawDef RawContractDefinition
	err = json.Unmarshal(dat, &rawDef)
	if err != nil {
		return addResp, err
	}

	deployRequest.ContractDefinition.ABI = rawDef.ABI
	deployRequest.ContractDefinition.EntryPointsByType = rawDef.EntryPointsByType
	deployRequest.ContractDefinition.Program, err = CompressCompiledContract(rawDef.Program)
	if err != nil {
		return addResp, err
	}

	pay, err := json.Marshal(deployRequest)
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

func (sg StarknetGateway) EstimateFee(sn StarknetRequest) (fee FeeEstimate, err error) {
	url := fmt.Sprintf("%s/estimate_fee", sg.Feeder)

	pay, err := json.Marshal(sn)
	if err != nil {
		return fee, err
	}

	rawResp, err := postHelper(pay, url)
	if err != nil {
		return fee, err
	}

	err = json.Unmarshal(rawResp, &fee)
	return fee, err
}

func (sg StarknetGateway) GetAccountNonce(address *big.Int) (nonce *big.Int, err error) {
	resp, err := sg.Call(StarknetRequest{
		ContractAddress:    BigToHex(address),
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

func FmtExecuteCalldataStrings(nonce *big.Int, txs []Transaction) (calldataStrings []string) {
	callArray := FmtExecuteCalldata(nonce, txs)
	for _, data := range callArray {
		calldataStrings = append(calldataStrings, data.String())
	}
	return calldataStrings
}

func FmtExecuteCalldata(nonce *big.Int, txs []Transaction) (calldataArray []*big.Int) {
	callArray := []*big.Int{big.NewInt(int64(len(txs)))}

	for _, tx := range txs {
		callArray = append(callArray, tx.ContractAddress, tx.EntryPointSelector)
		if len(tx.Calldata) == 0 {
			callArray = append(callArray, big.NewInt(0), big.NewInt(0))
		} else {
			callArray = append(callArray, big.NewInt(int64(len(calldataArray))), big.NewInt(int64(len(tx.Calldata))))
			calldataArray = append(calldataArray, tx.Calldata...)
		}
	}

	callArray = append(callArray, big.NewInt(int64(len(calldataArray))))
	callArray = append(callArray, calldataArray...)
	callArray = append(callArray, nonce)
	return callArray
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

func (sc StarkCurve) HashMulticall(addr, nonce, maxFee, chainId *big.Int, txs []Transaction) (hash *big.Int, err error) {
	callArray := FmtExecuteCalldata(nonce, txs)
	callArray = append(callArray, big.NewInt(int64(len(callArray))))
	cdHash, err := sc.HashElements(callArray)
	if err != nil {
		return hash, err
	}

	multiHashData := []*big.Int{
		UTF8StrToBig(TRANSACTION_PREFIX),
		big.NewInt(TRANSACTION_VERSION),
		addr,
		GetSelectorFromName(EXECUTE_SELECTOR),
		cdHash,
		maxFee,
		chainId,
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
