package contracts

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"

	"github.com/dontpanicdao/caigo/gateway"
	"github.com/dontpanicdao/caigo/types"
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

func (p *GatewayProvider) declareAndWaitNoWallet(ctx context.Context, compiledClass []byte) (*DeclareOutput, error) {
	provider := gateway.GatewayProvider(*p)
	class := types.ContractClass{}
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
		receipt.Status == types.TransactionRejected {
		return nil, fmt.Errorf("wrong status: %s", receipt.Status)
	}
	return &DeclareOutput{
		classHash:       tx.ClassHash,
		transactionHash: tx.TransactionHash,
	}, nil
}

// DeployAndWaitNoWallet run the DEPLOY transaction to deploy a contract and
// wait for it to be final with the blockchain.
//
// Deprecated: this command should be replaced by an Invoke on a class or a
// DEPLOY_ACCOUNT for an account.
func (p *GatewayProvider) deployAndWaitNoWallet(ctx context.Context, compiledClass []byte, salt string, inputs []string) (*DeployOutput, error) {
	provider := gateway.GatewayProvider(*p)
	class := types.ContractClass{}
	if err := json.Unmarshal(compiledClass, &class); err != nil {
		return nil, err
	}
	tx, err := provider.Deploy(ctx, class, types.DeployRequest{
		ContractAddressSalt: salt,
		ConstructorCalldata: inputs,
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
