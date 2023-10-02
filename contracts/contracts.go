package contracts

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"math/big"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/gateway"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/types"
)

type DeclareOutput struct {
	classHash       string
	transactionHash string
}

type DeployOutput struct {
	ContractAddress string
	ClassHash       string
	TransactionHash string
}

type GatewayProvider gateway.GatewayProvider

// declareAndWaitWithWallet declares a contract class and waits for the transaction to be finalized.
//
// ctx: The context.Context object for controlling the lifespan of this operation.
// compiledClass: The compiled contract class in the form of a byte array.
//
// Returns a pointer to a DeclareOutput struct and an error if any.
func (p *GatewayProvider) declareAndWaitWithWallet(ctx context.Context, compiledClass []byte) (*DeclareOutput, error) {
	provider := gateway.GatewayProvider(*p)
	class := rpc.DeprecatedContractClass{}
	if err := json.Unmarshal(compiledClass, &class); err != nil {
		return nil, err
	}
	tx, err := provider.Declare(ctx, class, gateway.DeclareRequest{})
	if err != nil {
		return nil, err
	}
	_, receipt, err := (&provider).WaitForTransaction(ctx, tx.TransactionHash, 3, 10)
	if err != nil {
		return nil, err
	}
	if !receipt.Status.IsTransactionFinal() ||
		receipt.Status == types.TransactionState(rpc.TxnExecutionStatusREVERTED) {
		return nil, fmt.Errorf("wrong status: %s", receipt.Status)
	}
	return &DeclareOutput{
		classHash:       tx.ClassHash,
		transactionHash: tx.TransactionHash,
	}, nil
}


// deployAccountAndWaitNoWallet deploys an account and waits for it to be mined without using a wallet.
//
// ctx: The context.Context to use for the request.
// classHash: The hash of the class.
// compiledClass: The compiled class.
// salt: The salt to use for the contract address.
// inputs: The constructor calldata.
//
// *DeployOutput: The deployed contract address and transaction hash.
// error: An error if the deployment fails.
// TODO: remove compiledClass from the interface
func (p *GatewayProvider) deployAccountAndWaitNoWallet(ctx context.Context, classHash *felt.Felt, compiledClass []byte, salt string, inputs []string) (*DeployOutput, error) {
	provider := gateway.GatewayProvider(*p)
	class := rpc.DeprecatedContractClass{}

	if err := json.Unmarshal(compiledClass, &class); err != nil {
		return nil, err
	}
	fmt.Printf("classHash %v\n", classHash.String())
	tx, err := provider.DeployAccount(ctx, types.DeployAccountRequest{
		// MaxFee
		Version:             big.NewInt(1),
		ContractAddressSalt: salt,
		ConstructorCalldata: inputs,
		ClassHash:           classHash.String(),
	})

	if err != nil {
		return nil, err
	}

	_, receipt, err := (&provider).WaitForTransaction(ctx, tx.TransactionHash, 8, 60)

	if err != nil {
		log.Printf("contract Address: %s\n", tx.ContractAddress)
		log.Printf("transaction Hash: %s\n", tx.TransactionHash)
		return nil, err
	}

	if !receipt.Status.IsTransactionFinal() ||

		receipt.Status == types.TransactionRejected {
		return nil, fmt.Errorf("wrong status: %s", receipt.Status)
	}

	return &DeployOutput{
		ContractAddress: tx.ContractAddress,
		TransactionHash: tx.TransactionHash,
	}, nil
}


// deployAndWaitNoWallet deploys a contract using the given compiled class, salt, and inputs.
// It returns the deployed contract's address (wait for it to be final) and transaction hash, or an error if deployment fails.
//
// ctx: The context.Context used for the deployment.
// compiledClass: The compiled class of the contract to be deployed.
// salt: The salt used for generating the contract's address.
// inputs: The constructor calldata for initializing the contract.
//
// Returns a pointer to DeployOutput and an error.
// Deprecated: this command should be replaced by an Invoke on a class or a
// DEPLOY_ACCOUNT for an account.
func (p *GatewayProvider) deployAndWaitNoWallet(ctx context.Context, compiledClass []byte, salt string, inputs []string) (*DeployOutput, error) {
	provider := gateway.GatewayProvider(*p)
	class := rpc.DeprecatedContractClass{}
	if err := json.Unmarshal(compiledClass, &class); err != nil {
		return nil, err
	}
	tx, err := provider.Deploy(ctx, class, rpc.DeployAccountTxn{})

	// tx, err := provider.Deploy(ctx, class, rpc.DeployRequest{
	// 	ContractAddressSalt: salt,
	// 	ConstructorCalldata: inputs,
	// })
	if err != nil {
		return nil, err
	}
	_, receipt, err := (&provider).WaitForTransaction(ctx, tx.TransactionHash, 8, 60)
	if err != nil {
		log.Printf("contract Address: %s\n", tx.ContractAddress)
		log.Printf("transaction Hash: %s\n", tx.TransactionHash)
		return nil, err
	}
	if !receipt.Status.IsTransactionFinal() ||
		receipt.Status == types.TransactionRejected {
		return nil, fmt.Errorf("wrong status: %s", receipt.Status)
	}
	return &DeployOutput{
		ContractAddress: tx.ContractAddress,
		TransactionHash: tx.TransactionHash,
	}, nil
}
