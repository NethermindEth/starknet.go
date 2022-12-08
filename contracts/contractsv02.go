package contracts

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/dontpanicdao/caigo/rpcv02"
	"github.com/dontpanicdao/caigo/types"
)

type RPCv02Provider rpcv02.Provider

func (p *RPCv02Provider) declareAndWaitWithWallet(ctx context.Context, compiledClass []byte) (*DeclareOutput, error) {
	provider := rpcv02.Provider(*p)
	class := types.ContractClass{}
	if err := json.Unmarshal(compiledClass, &class); err != nil {
		return nil, err
	}
	tx, err := provider.AddDeclareTransaction(ctx, rpcv02.BroadcastedDeclareTransaction{
		BroadcastedTxnCommonProperties: rpcv02.BroadcastedTxnCommonProperties{
			Type: "DECLARE",
		},
		ContractClass: class,
		SenderAddress: types.HexToHash("0x0"),
	})
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

func (p *RPCv02Provider) deployAccountAndWaitWithoutWallet(ctx context.Context, compiledClass []byte, salt string, inputs []string) (*DeployOutput, error) {
	provider := rpcv02.Provider(*p)
	class := types.ContractClass{}
	if err := json.Unmarshal(compiledClass, &class); err != nil {
		return nil, err
	}
	tx, err := provider.AddDeployAccountTransaction(ctx, rpcv02.BroadcastedDeployAccountTransaction{
		BroadcastedTxnCommonProperties: rpcv02.BroadcastedTxnCommonProperties{},
		ContractAddressSalt: salt,
		ConstructorCalldata: inputs,
		
	})
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

func (p *RPCv02Provider) deployAndWaitWithWallet(ctx context.Context, compiledClass []byte, salt string, inputs []string) (*DeployOutput, error) {
	provider := rpcv02.Provider(*p)
	class := types.ContractClass{}
	if err := json.Unmarshal(compiledClass, &class); err != nil {
		return nil, err
	}
	tx, err := provider.AddDeployTransaction(ctx, rpcv02.BroadcastedDeployTransaction{
		ContractAddressSalt: salt,
		ConstructorCalldata: inputs,
		ContractClass:       class,
	})
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
