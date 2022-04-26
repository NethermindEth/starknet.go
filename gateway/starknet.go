package gateway

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"

	"github.com/dontpanicdao/caigo"
	"github.com/dontpanicdao/caigo/types"
	"github.com/google/go-querystring/query"
)

type StarkResp struct {
	Result []string `json:"result"`
}

/*
	'call_contract' wrapper and can accept a blockId in the hash or height format
*/
func (sg *Gateway) Call(ctx context.Context, tx Transaction, opts *BlockOptions) ([]string, error) {
	tx.EntryPointSelector = caigo.BigToHex(caigo.GetSelectorFromName(tx.EntryPointSelector))
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
func (sg *Gateway) Invoke(ctx context.Context, tx types.Transaction) (*types.AddTxResponse, error) {
	tx.EntryPointSelector = caigo.BigToHex(caigo.GetSelectorFromName(tx.EntryPointSelector))
	tx.Type = INVOKE

	if len(tx.Calldata) == 0 {
		tx.Calldata = []string{}
	}
	if len(tx.Signature) == 0 {
		tx.Signature = []string{}
	}

	req, err := sg.newRequest(ctx, http.MethodPost, "/add_transaction", tx)
	if err != nil {
		return nil, err
	}

	var resp types.AddTxResponse
	return &resp, sg.do(req, &resp)
}

/*
	'add_transaction' wrapper for compressing and deploying a compiled StarkNet contract
*/
func (sg *Gateway) Deploy(ctx context.Context, filePath string, deployRequest types.DeployRequest) (resp types.AddTxResponse, err error) {
	dat, err := os.ReadFile(filePath)
	if err != nil {
		return resp, err
	}

	deployRequest.Type = DEPLOY
	if len(deployRequest.ConstructorCalldata) == 0 {
		deployRequest.ConstructorCalldata = []string{}
	}

	var rawDef types.RawContractDefinition
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
