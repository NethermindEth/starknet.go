package caigo

import (
	"math/big"
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

/*
	GETTER Models
*/
type ContractCode struct {
	Bytecode []string `json:"bytecode"`
	Abi      []ABI    `json:"abi"`
}

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

/*
	SETTER Models
*/

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
	Gateway *StarknetGateway
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
