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

// DeployAndWaitNoWallet run the DEPLOY transaction to deploy a contract and
// wait for it to be final with the blockchain.
//
// Deprecated: this command should be replaced by an Invoke on a class or a
// DEPLOY_ACCOUNT for an account.
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
