package gateway

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"
	"os"

	"github.com/dontpanicdao/caigo"
	"github.com/dontpanicdao/caigo/types"
	"github.com/google/go-querystring/query"
)

type StarkResp struct {
	Result []string `json:"result"`
}

type StateUpdate struct {
	BlockHash string `json:"block_hash"`
	NewRoot   string `json:"new_root"`
	OldRoot   string `json:"old_root"`
	StateDiff struct {
		StorageDiffs      map[string]interface{} `json:"storage_diffs"`
		DeployedContracts []struct {
			Address   string `json:"address"`
			ClassHash string `json:"class_hash"`
		} `json:"deployed_contracts"`
	} `json:"state_diff"`
}

func (sg *Gateway) ChainID(context.Context) (string, error) {
	return sg.ChainId, nil
}

/*
	'call_contract' wrapper and can accept a blockId in the hash or height format
*/
func (sg *Gateway) Call(ctx context.Context, call types.FunctionCall, blockHashOrTag string) ([]string, error) {
	call.EntryPointSelector = caigo.BigToHex(caigo.GetSelectorFromName(call.EntryPointSelector))
	if len(call.Calldata) == 0 {
		call.Calldata = []string{}
	}

	req, err := sg.newRequest(ctx, http.MethodPost, "/call_contract", call)
	if err != nil {
		return nil, err
	}

	if blockHashOrTag != "" {
		appendQueryValues(req, url.Values{
			"blockHash": []string{blockHashOrTag},
		})
	}

	var resp StarkResp
	return resp.Result, sg.do(req, &resp)
}

/*
	'add_transaction' wrapper for invokation requests
*/
func (sg *Gateway) Invoke(ctx context.Context, invoke types.FunctionInvoke) (*types.AddTxResponse, error) {
	var tx types.Transaction
	tx.EntryPointSelector = caigo.BigToHex(caigo.GetSelectorFromName(tx.EntryPointSelector))
	tx.Type = INVOKE

	if len(invoke.Calldata) == 0 {
		invoke.Calldata = []string{}
	}
	if len(invoke.Signature) == 0 {
		invoke.Signature = []*types.Felt{}
	}

	req, err := sg.newRequest(ctx, http.MethodPost, "/add_transaction", invoke)
	if err != nil {
		return nil, err
	}

	var resp types.AddTxResponse
	return &resp, sg.do(req, &resp)
}

type RawContractDefinition struct {
	ABI               []types.ABI             `json:"abi"`
	EntryPointsByType types.EntryPointsByType `json:"entry_points_by_type"`
	Program           map[string]interface{}  `json:"program"`
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
	if deployRequest.ContractAddressSalt == "" {
		deployRequest.ContractAddressSalt = "0x0"
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

/*
	'add_transaction' wrapper for compressing and declaring a contract class
*/
func (sg *Gateway) Declare(ctx context.Context, filePath string, declareRequest types.DeclareRequest) (resp types.AddTxResponse, err error) {
	dat, err := os.ReadFile(filePath)
	if err != nil {
		return resp, err
	}

	declareRequest.Type = DECLARE
	declareRequest.SenderAddress = "0x1"
	declareRequest.MaxFee = "0x0"
	declareRequest.Nonce = "0x0"
	declareRequest.Signature = []string{}

	var rawDef RawContractDefinition
	if err = json.Unmarshal(dat, &rawDef); err != nil {
		return resp, err
	}

	declareRequest.ContractClass.ABI = rawDef.ABI
	declareRequest.ContractClass.EntryPointsByType = rawDef.EntryPointsByType
	declareRequest.ContractClass.Program, err = CompressCompiledContract(rawDef.Program)
	if err != nil {
		return resp, err
	}

	req, err := sg.newRequest(ctx, http.MethodPost, "/add_transaction", declareRequest)
	if err != nil {
		return resp, err
	}

	return resp, sg.do(req, &resp)
}

func (sg *Gateway) StateUpdate(ctx context.Context, opts *BlockOptions) (*StateUpdate, error) {
	req, err := sg.newRequest(ctx, http.MethodGet, "/get_state_update", nil)
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

	var resp StateUpdate
	return &resp, sg.do(req, &resp)
}

func (sg *Gateway) ContractAddresses(ctx context.Context) (*types.ContractAddresses, error) {
	req, err := sg.newRequest(ctx, http.MethodGet, "/get_contract_addresses", nil)
	if err != nil {
		return nil, err
	}

	var resp types.ContractAddresses
	return &resp, sg.do(req, &resp)
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
