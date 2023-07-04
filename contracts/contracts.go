package contracts

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"math/big"

	"github.com/NethermindEth/caigo/gateway"
	"github.com/NethermindEth/caigo/rpcv02"
	"github.com/NethermindEth/caigo/types"
	"github.com/NethermindEth/juno/core/felt"
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

func (p *GatewayProvider) declareAndWaitWithWallet(ctx context.Context, compiledClass []byte) (*DeclareOutput, error) {
	provider := gateway.GatewayProvider(*p)
	class := rpcv02.ContractClass{}
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
		rpcv02.TransactionState(receipt.Status.String()) == rpcv02.TransactionRejected {
		return nil, fmt.Errorf("wrong status: %s", receipt.Status)
	}
	return &DeclareOutput{
		classHash:       tx.ClassHash,
		transactionHash: tx.TransactionHash,
	}, nil
}

// TODO: remove compiledClass from the interface
func (p *GatewayProvider) deployAccountAndWaitNoWallet(ctx context.Context, classHash *felt.Felt, compiledClass []byte, salt string, inputs []string) (*DeployOutput, error) {
	provider := gateway.GatewayProvider(*p)
	class := rpcv02.ContractClass{}

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

// DeployAndWaitNoWallet run the DEPLOY transaction to deploy a contract and
// wait for it to be final with the blockchain.
//
// Deprecated: this command should be replaced by an Invoke on a class or a
// DEPLOY_ACCOUNT for an account.
func (p *GatewayProvider) deployAndWaitNoWallet(ctx context.Context, compiledClass []byte, salt string, inputs []string) (*DeployOutput, error) {
	provider := gateway.GatewayProvider(*p)
	class := rpcv02.ContractClass{}
	if err := json.Unmarshal(compiledClass, &class); err != nil {
		return nil, err
	}
	tx, err := provider.Deploy(ctx, class, rpcv02.DeployAccountTxn{})

	// tx, err := provider.Deploy(ctx, class, rpcv02.DeployRequest{
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
