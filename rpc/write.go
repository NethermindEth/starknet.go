package rpc

import (
	"context"
	"fmt"

	"github.com/dontpanicdao/caigo"
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
func (sc *Client) AddInvokeTransaction(ctx context.Context, call types.FunctionCall, signature []string, maxFee string, version string) (*AddInvokeTransactionOutput, error) {
	call.EntryPointSelector = fmt.Sprintf("0x%s", caigo.GetSelectorFromName(call.EntryPointSelector).Text(16))
	var output AddInvokeTransactionOutput
	if err := sc.do(ctx, "starknet_addInvokeTransaction", &output, call, signature, maxFee, version); err != nil {
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
