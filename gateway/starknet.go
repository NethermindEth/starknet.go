package gateway

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

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

type GatewayFunctionCall struct {
	types.FunctionCall
	Signature []string `json:"signature"`
}

/*
'call_contract' wrapper and can accept a blockId in the hash or height format
*/
func (sg *Gateway) Call(ctx context.Context, call types.FunctionCall, blockHashOrTag string) ([]string, error) {
	gc := GatewayFunctionCall{
		FunctionCall: call,
	}
	gc.EntryPointSelector = types.BigToHex(types.GetSelectorFromName(gc.EntryPointSelector))
	if len(gc.Calldata) == 0 {
		gc.Calldata = []string{}
	}

	if len(gc.Signature) == 0 {
		gc.Signature = []string{"0", "0"} // allows rpc and http clients to implement(has to be a better way)
	}

	req, err := sg.newRequest(ctx, http.MethodPost, "/call_contract", gc)
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
func (sg *Gateway) Invoke(ctx context.Context, invoke types.FunctionInvoke) (*types.AddInvokeTransactionOutput, error) {
	tx := Transaction{
		Type:            INVOKE,
		ContractAddress: invoke.ContractAddress.Hex(),
		Version:         fmt.Sprintf("0x%d", invoke.Version),
		MaxFee:          invoke.MaxFee.String(),
	}
	if invoke.EntryPointSelector != "" {
		tx.EntryPointSelector = types.BigToHex(types.GetSelectorFromName(invoke.EntryPointSelector))
	}
	if invoke.Nonce != nil {
		tx.Nonce = invoke.Nonce.String()
	}

	if len(invoke.Calldata) == 0 {
		tx.Calldata = []string{}
	} else {
		tx.Calldata = invoke.Calldata
	}

	if len(invoke.Signature) == 0 {
		tx.Signature = []string{}
	} else {
		// stop-gap before full types.Felt cutover
		tx.Signature = []string{invoke.Signature[0].String(), invoke.Signature[1].String()}
	}

	req, err := sg.newRequest(ctx, http.MethodPost, "/add_transaction", tx)
	if err != nil {
		return nil, err
	}

	var resp types.AddInvokeTransactionOutput
	return &resp, sg.do(req, &resp)
}

/*
'add_transaction' wrapper for compressing and deploying a compiled StarkNet contract
*/
func (sg *Gateway) Deploy(ctx context.Context, filePath string, deployRequest types.DeployRequest) (resp types.AddDeployResponse, err error) {
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

	var rawDef types.ContractClass
	if err = json.Unmarshal(dat, &rawDef); err != nil {
		return resp, err
	}

	deployRequest.ContractDefinition = rawDef
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
func (sg *Gateway) Declare(ctx context.Context, filePath string, declareRequest DeclareRequest) (resp types.AddDeclareResponse, err error) {
	dat, err := os.ReadFile(filePath)
	if err != nil {
		return resp, err
	}

	declareRequest.Type = DECLARE
	declareRequest.SenderAddress = "0x1"
	declareRequest.MaxFee = "0x0"
	declareRequest.Nonce = "0x0"
	declareRequest.Signature = []string{}

	var rawDef types.ContractClass
	if err = json.Unmarshal(dat, &rawDef); err != nil {
		return resp, err
	}

	declareRequest.ContractClass = rawDef
	if err != nil {
		return resp, err
	}

	req, err := sg.newRequest(ctx, http.MethodPost, "/add_transaction", declareRequest)
	if err != nil {
		return resp, err
	}

	return resp, sg.do(req, &resp)
}

type DeclareRequest struct {
	Type          string              `json:"type"`
	SenderAddress string              `json:"sender_address"`
	MaxFee        string              `json:"max_fee"`
	Nonce         string              `json:"nonce"`
	Signature     []string            `json:"signature"`
	ContractClass types.ContractClass `json:"contract_class"`
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

func (sg *Gateway) ContractAddresses(ctx context.Context) (*ContractAddresses, error) {
	req, err := sg.newRequest(ctx, http.MethodGet, "/get_contract_addresses", nil)
	if err != nil {
		return nil, err
	}

	var resp ContractAddresses
	return &resp, sg.do(req, &resp)
}

type ContractAddresses struct {
	Starknet             string `json:"Starknet"`
	GpsStatementVerifier string `json:"GpsStatementVerifier"`
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
