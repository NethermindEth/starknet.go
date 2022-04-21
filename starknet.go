package caigo

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"net/http"
	"os"

	"github.com/google/go-querystring/query"
)

const (
	INVOKE              string = "INVOKE_FUNCTION"
	DEPLOY              string = "DEPLOY"
	GOERLI_ID           string = "SN_GOERLI"
	MAINNET_ID          string = "SN_MAIN"
	LOCAL_BASE          string = "http://localhost:5000"
	GOERLI_BASE         string = "https://alpha4.starknet.io"
	MAINNET_BASE        string = "https://alpha-mainnet.starknet.io"
	EXECUTE_SELECTOR    string = "__execute__"
	TRANSACTION_PREFIX  string = "invoke"
	TRANSACTION_VERSION int64  = 0
)

type ABI struct {
	Members []struct {
		Name   string `json:"name"`
		Offset int    `json:"offset"`
		Type   string `json:"type"`
	} `json:"members,omitempty"`
	Name   string `json:"name"`
	Size   int    `json:"size,omitempty"`
	Type   string `json:"type"`
	Inputs []struct {
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"inputs,omitempty"`
	Outputs         []interface{} `json:"outputs,omitempty"`
	StateMutability string        `json:"stateMutability,omitempty"`
}

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

/*
	'call_contract' wrapper and can accept a blockId in the hash or height format
*/
func (sg *StarknetGateway) Call(ctx context.Context, tx Transaction, opts *BlockOptions) ([]string, error) {
	tx.EntryPointSelector = BigToHex(GetSelectorFromName(tx.EntryPointSelector))
	if len(tx.Calldata) == 0 {
		tx.Calldata = []string{}
	}
	if len(tx.Signature) == 0 {
		tx.Signature = []string{}
	}

	req, err := sg.newRequest(ctx, http.MethodPost, "/call_contract", tx)
	if err != nil {
		return nil, err
	}

	if opts != nil {
		vs, err := query.Values(opts)
		if err != nil {
			return nil, err
		}
		appendQueryValues(req, vs)
	}

	var resp StarkResp
	return resp.Result, sg.do(req, &resp)
}

/*
	'add_transaction' wrapper for invokation requests
*/
func (sg *StarknetGateway) Invoke(ctx context.Context, tx Transaction) (resp AddTxResponse, err error) {
	tx.EntryPointSelector = BigToHex(GetSelectorFromName(tx.EntryPointSelector))
	tx.Type = INVOKE

	if len(tx.Calldata) == 0 {
		tx.Calldata = []string{}
	}
	if len(tx.Signature) == 0 {
		tx.Signature = []string{}
	}

	req, err := sg.newRequest(ctx, http.MethodPost, "/add_transaction", tx)
	if err != nil {
		return resp, err
	}

	return resp, sg.do(req, &resp)
}

/*
	'add_transaction' wrapper for compressing and deploying a compiled StarkNet contract
*/
func (sg *StarknetGateway) Deploy(ctx context.Context, filePath string, deployRequest DeployRequest) (resp AddTxResponse, err error) {
	dat, err := os.ReadFile(filePath)
	if err != nil {
		return resp, err
	}

	deployRequest.Type = DEPLOY
	if len(deployRequest.ConstructorCalldata) == 0 {
		deployRequest.ConstructorCalldata = []string{}
	}

	var rawDef RawContractDefinition
	if err = json.Unmarshal(dat, &rawDef); err != nil {
		return resp, err
	}

	deployRequest.ContractDefinition.ABI = rawDef.ABI
	deployRequest.ContractDefinition.EntryPointsByType = rawDef.EntryPointsByType
	deployRequest.ContractDefinition.Program, err = CompressCompiledContract(rawDef.Program)
	if err != nil {
		return resp, err
	}

	req, err := sg.newRequest(ctx, http.MethodPost, "/add_transaction", deployRequest)
	if err != nil {
		return resp, err
	}

	return resp, sg.do(req, &resp)
}

func CompressCompiledContract(program map[string]interface{}) (cc string, err error) {
	pay, err := json.Marshal(program)
	if err != nil {
		return cc, err
	}

	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	if _, err = zw.Write(pay); err != nil {
		return cc, err
	}
	if err := zw.Close(); err != nil {
		return cc, err
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
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
