package rpc

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/dontpanicdao/caigo/rpc/types"
)

// AddDeclareTransactionOutput provides the output for AddDeclareTransaction.
type AddDeclareTransactionOutput struct {
	TransactionHash string `json:"transaction_hash"`
	ClassHash       string `json:"class_hash"`
}

// AddDeployTransactionOutput provides the output for AddDeployTransaction.
type AddDeployTransactionOutput struct {
	TransactionHash string `json:"transaction_hash"`
	ContractAddress string `json:"contract_address"`
}

// AddInvokeTransactionOutput provides the output for AddInvokeTransaction.
type AddInvokeTransactionOutput struct {
	TransactionHash string `json:"transaction_hash"`
}

// AddInvokeTransaction estimates the fee for a given StarkNet transaction.
func (sc *Client) AddInvokeTransaction(ctx context.Context, broadcastedInvokeTxn types.BroadcastedInvokeTxn) (*AddInvokeTransactionOutput, error) {
	// TODO: We might have to gzip/base64 the program and provide helpers to call
	// this API
	var output AddInvokeTransactionOutput
	if err := sc.do(ctx, "starknet_addInvokeTransaction", &output, broadcastedInvokeTxn); err != nil {
		return nil, err
	}
	return &output, nil
}

// AddDeclareTransaction submits a new class declaration transaction.
func (sc *Client) AddDeclareTransaction(ctx context.Context, contractClass types.ContractClass, version string) (*AddDeclareTransactionOutput, error) {
	// TODO: We might have to gzip/base64 the program and provide helpers to call
	// this API
	var result AddDeclareTransactionOutput
	if err := sc.do(ctx, "starknet_addDeclareTransaction", &result, contractClass, version); err != nil {
		return nil, err
	}
	return &result, nil
}

// AddDeployTransaction allows to declare a class and instantiate the
// associated contract in one command. This function will be deprecated and
// replaced by AddDeclareTransaction to declare a class, followed by
// AddInvokeTransaction to instantiate the contract. For now, it remains the only
// way to deploy an account without being charged for it.
func (sc *Client) AddDeployTransaction(ctx context.Context, salt string, inputs []string, contractClass types.ContractClass) (*AddDeployTransactionOutput, error) {
	var result AddDeployTransactionOutput
	if err := sc.do(ctx, "starknet_addDeployTransaction", &result, salt, inputs, contractClass); err != nil {
		return nil, err
	}
	return &result, nil
}

// Keep that function to build helper with broadcastedDeployTxn and broadcastedDeclareTxn
func encodeProgram(content []byte) (string, error) {
	buf := bytes.NewBuffer(nil)
	gzipContent := gzip.NewWriter(buf)
	_, err := gzipContent.Write(content)
	if err != nil {
		return "", err
	}
	gzipContent.Close()
	program := base64.StdEncoding.EncodeToString(buf.Bytes())
	return program, nil
}

// TODO: replace this function with an Unmarshal function.
func guessABI(abis []interface{}) (*types.ABI, error) {
	output := types.ABI{}
	for _, abi := range abis {
		if checkABI, ok := abi.(map[string]interface{}); ok {
			var ab types.ABIEntry
			switch checkABI["type"] {
			case "constructor", "function", "l1_handler":
				ab = &types.FunctionABIEntry{}
			case "struct":
				ab = &types.StructABIEntry{}
			case "event":
				ab = &types.EventABIEntry{}
			default:
				return nil, fmt.Errorf("unknown ABI type %v", checkABI["type"])
			}
			data, err := json.Marshal(checkABI)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(data, ab)
			if err != nil {
				return nil, err
			}
			output = append(output, ab)
		}
	}
	return &output, nil
}
