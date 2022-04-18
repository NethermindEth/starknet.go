package caigo

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"strings"
)

/*
	Instantiate a new StarkNet Signer which includes structures for calling the network and signing transactions:
	- private signing key
	- stark curve definition
	- full StarknetGateway definition
	- public key pair for signature verifications
*/
func (sc StarkCurve) NewSigner(private, pubX, pubY *big.Int, chainId ...string) (signer Signer, err error) {
	if len(sc.ConstantPoints) == 0 {
		return signer, fmt.Errorf("must initiate precomputed constant points")
	}
	var gw *StarknetGateway
	if len(chainId) == 1 {
		gw = NewGateway(WithChain(chainId[0]))
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

/*
	'call_contract' wrapper and can accept a blockId in the hash or height format
*/
func (sg *StarknetGateway) Call(ctx context.Context, tx Transaction, blockId ...string) (resp []string, err error) {
	bid := ""
	if len(blockId) == 1 {
		bid = fmtBlockId(blockId[0])
	}

	url := fmt.Sprintf("%s/call_contract%s", sg.Feeder, strings.Replace(bid, "&", "?", 1))
	tx.EntryPointSelector = BigToHex(GetSelectorFromName(tx.EntryPointSelector))

	if len(tx.Calldata) == 0 {
		tx.Calldata = []string{}
	}
	if len(tx.Signature) == 0 {
		tx.Signature = []string{}
	}

	pay, err := json.Marshal(tx)
	if err != nil {
		return resp, err
	}

	rawResp, err := sg.postHelper(ctx, pay, url)
	if err != nil {
		return resp, err
	}

	var snResp StarkResp
	err = json.Unmarshal(rawResp, &snResp)
	return snResp.Result, err
}

/*
	'add_transaction' wrapper for invokation requests
*/
func (sg *StarknetGateway) Invoke(ctx context.Context, tx Transaction) (addResp AddTxResponse, err error) {
	url := fmt.Sprintf("%s/add_transaction", sg.Gateway)
	tx.EntryPointSelector = BigToHex(GetSelectorFromName(tx.EntryPointSelector))
	tx.Type = INVOKE

	if len(tx.Calldata) == 0 {
		tx.Calldata = []string{}
	}
	if len(tx.Signature) == 0 {
		tx.Signature = []string{}
	}

	pay, err := json.Marshal(tx)
	if err != nil {
		return addResp, err
	}

	rawResp, err := sg.postHelper(ctx, pay, url)
	if err != nil {
		return addResp, err
	}

	err = json.Unmarshal(rawResp, &addResp)
	return addResp, err
}

/*
	'add_transaction' wrapper for compressing and deploying a compiled StarkNet contract
*/
func (sg *StarknetGateway) Deploy(ctx context.Context, filePath string, deployRequest DeployRequest) (addResp AddTxResponse, err error) {
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

	rawResp, err := sg.postHelper(ctx, pay, url)
	if err != nil {
		return addResp, err
	}

	err = json.Unmarshal(rawResp, &addResp)
	return addResp, err
}

/*
	invocation wrapper for StarkNet account calls to '__execute__' contact calls through an account abstraction
	- implementation has been tested against OpenZeppelin Account contract as of: https://github.com/OpenZeppelin/cairo-contracts/blob/4116c1ecbed9f821a2aa714c993a35c1682c946e/src/openzeppelin/account/Account.cairo
	- accepts a multicall
*/
func (signer Signer) Execute(ctx context.Context, address string, txs []Transaction) (addResp AddTxResponse, err error) {
	nonce, err := signer.Gateway.AccountNonce(ctx, address)
	if err != nil {
		return addResp, err
	}

	maxFee := big.NewInt(0)
	// for _, tx := range txs {
	// 	var cdStrings []string
	// 	for _, data := range tx.Calldata {
	// 		cdStrings = append(cdStrings, data.String())
	// 	}

	// 	innerFee, err := signer.Gateway.EstimateFee(Transaction{
	// 		ContractAddress: tx.ContractAddress,
	// 		EntryPointSelector: tx.EntryPointSelector,
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

	req := Transaction{
		ContractAddress:    address,
		EntryPointSelector: BigToHex(GetSelectorFromName(EXECUTE_SELECTOR)),
		Calldata:           FmtExecuteCalldataStrings(nonce, txs),
		Signature:          []string{r.String(), s.String()},
	}

	return signer.Gateway.Invoke(ctx, req)
}

func (sg *StarknetGateway) EstimateFee(ctx context.Context, tx Transaction) (fee FeeEstimate, err error) {
	url := fmt.Sprintf("%s/estimate_fee", sg.Feeder)

	pay, err := json.Marshal(tx)
	if err != nil {
		return fee, err
	}

	rawResp, err := sg.postHelper(ctx, pay, url)
	if err != nil {
		return fee, err
	}

	err = json.Unmarshal(rawResp, &fee)
	return fee, err
}

func (sg *StarknetGateway) postHelper(ctx context.Context, pay []byte, url string) (resp []byte, err error) {
	method := "POST"

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(pay))
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

/*
	Formats the multicall transactions in a format which can be signed and verified by the network and OpenZeppelin account contracts
*/
func FmtExecuteCalldata(nonce *big.Int, txs []Transaction) (calldataArray []*big.Int) {
	callArray := []*big.Int{big.NewInt(int64(len(txs)))}

	for _, tx := range txs {
		callArray = append(callArray, SNValToBN(tx.ContractAddress), GetSelectorFromName(tx.EntryPointSelector))
		if len(tx.Calldata) == 0 {
			callArray = append(callArray, big.NewInt(0), big.NewInt(0))
		} else {
			callArray = append(callArray, big.NewInt(int64(len(calldataArray))), big.NewInt(int64(len(tx.Calldata))))
			for _, cd := range tx.Calldata {
				calldataArray = append(calldataArray, SNValToBN(cd))
			}
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

func (sc StarkCurve) HashMulticall(addr string, nonce, maxFee, chainId *big.Int, txs []Transaction) (hash *big.Int, err error) {
	callArray := FmtExecuteCalldata(nonce, txs)
	callArray = append(callArray, big.NewInt(int64(len(callArray))))
	cdHash, err := sc.HashElements(callArray)
	if err != nil {
		return hash, err
	}

	multiHashData := []*big.Int{
		UTF8StrToBig(TRANSACTION_PREFIX),
		big.NewInt(TRANSACTION_VERSION),
		SNValToBN(addr),
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
	calldataArray := []*big.Int{big.NewInt(int64(len(tx.Calldata)))}
	for _, cd := range tx.Calldata {
		calldataArray = append(calldataArray, HexToBN(cd))
	}

	cdHash, err := sc.HashElements(calldataArray)
	if err != nil {
		return hash, err
	}

	txHashData := []*big.Int{
		addr,
		SNValToBN(tx.ContractAddress),
		GetSelectorFromName(tx.EntryPointSelector),
		cdHash,
		SNValToBN(tx.Nonce),
	}

	txHashData = append(txHashData, big.NewInt(int64(len(txHashData))))
	hash, err = sc.HashElements(txHashData)
	return hash, err
}

// Adheres to 'starknet.js' hash non typedData
func (sc StarkCurve) HashTx(addr *big.Int, tx Transaction) (hash *big.Int, err error) {
	calldataArray := []*big.Int{big.NewInt(int64(len(tx.Calldata)))}
	for _, cd := range tx.Calldata {
		calldataArray = append(calldataArray, SNValToBN(cd))
	}

	cdHash, err := sc.HashElements(calldataArray)
	if err != nil {
		return hash, err
	}

	txHashData := []*big.Int{
		SNValToBN(tx.ContractAddress),
		GetSelectorFromName(tx.EntryPointSelector),
		cdHash,
	}

	txHashData = append(txHashData, big.NewInt(int64(len(txHashData))))
	hash, err = sc.HashElements(txHashData)
	return hash, err
}
