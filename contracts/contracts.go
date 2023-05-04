package contracts

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/dontpanicdao/caigo/gateway"
	"github.com/dontpanicdao/caigo/rpcv01"
	"github.com/dontpanicdao/caigo/types"
)

type RPCv01Provider rpcv01.Provider

type DeclareOutput struct {
	classHash       string
	transactionHash string
}

func (p *RPCv01Provider) declareAndWaitNoWallet(ctx context.Context, compiledClass []byte) (*DeclareOutput, error) {
	provider := rpcv01.Provider(*p)
	class := types.ContractClass{}
	if err := json.Unmarshal(compiledClass, &class); err != nil {
		return nil, err
	}
	tx, err := provider.AddDeclareTransaction(ctx, class, "0x0")
	if err != nil {
		return nil, err
	}
	status, err := provider.WaitForTransaction(ctx, types.HexToHash(tx.TransactionHash), 8*time.Second)
	if err != nil {
		log.Printf("class Hash: %s\n", tx.ClassHash)
		log.Printf("transaction Hash: %s\n", tx.TransactionHash)
		return nil, err
	}
	if status == types.TransactionRejected {
		log.Printf("class Hash: %s\n", tx.ClassHash)
		log.Printf("transaction Hash: %s\n", tx.TransactionHash)
		return nil, errors.New("declare rejected")
	}
	return &DeclareOutput{
		classHash:       tx.ClassHash,
		transactionHash: tx.TransactionHash,
	}, nil
}

type DeployOutput struct {
	ContractAddress string
	ClassHash       string
	TransactionHash string
}

// DeployAndWaitNoWallet run the DEPLOY transaction to deploy a contract and
// wait for it to be final with the blockchain.
//
// Deprecated: this command should be replaced by an Invoke on a class or a
// DEPLOY_ACCOUNT for an account.
func (p *RPCv01Provider) deployAndWaitNoWallet(ctx context.Context, compiledClass []byte, salt string, inputs []string) (*DeployOutput, error) {
	provider := rpcv01.Provider(*p)
	class := types.ContractClass{}
	if err := json.Unmarshal(compiledClass, &class); err != nil {
		return nil, err
	}
	tx, err := provider.AddDeployTransaction(ctx, salt, inputs, class)
	if err != nil {
		return nil, err
	}
	status, err := provider.WaitForTransaction(ctx, types.HexToHash(tx.TransactionHash), 8*time.Second)
	if err != nil {
		log.Printf("contract Address: %s\n", tx.ContractAddress)
		log.Printf("transaction Hash: %s\n", tx.TransactionHash)
		return nil, err
	}
	if status == types.TransactionRejected {
		log.Printf("contract Address: %s\n", tx.ContractAddress)
		log.Printf("transaction Hash: %s\n", tx.TransactionHash)
		return nil, errors.New("deploy rejected")
	}
	return &DeployOutput{
		ContractAddress: tx.ContractAddress,
		TransactionHash: tx.TransactionHash,
	}, nil
}

type GatewayProvider gateway.GatewayProvider

func (p *GatewayProvider) declareAndWaitWithWallet(ctx context.Context, compiledClass []byte) (*DeclareOutput, error) {
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

// TODO: remove compiledClass from the interface
func (p *GatewayProvider) deployAccountAndWaitNoWallet(ctx context.Context, classHash types.Hash, compiledClass []byte, salt string, inputs []string) (*DeployOutput, error) {
	provider := gateway.GatewayProvider(*p)
	class := types.ContractClass{}

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
