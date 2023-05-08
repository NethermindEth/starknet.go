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

	"github.com/smartcontractkit/caigo/rpcv02"
	"github.com/smartcontractkit/caigo/types"
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
			Type:    "DECLARE",
			MaxFee:  big.NewInt(10000),
			Version: "0x01",
			Nonce:   big.NewInt(1), // TODO: nonce handling
		},
		ContractClass: class,
		SenderAddress: types.HexToHash("0x01"), // TODO: contant devnet address
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

func (p *RPCv02Provider) deployAccountAndWaitNoWallet(ctx context.Context, classHash types.Hash, compiledClass []byte, salt string, inputs []string) (*DeployOutput, error) {
	provider := rpcv02.Provider(*p)
	class := types.ContractClass{}
	if err := json.Unmarshal(compiledClass, &class); err != nil {
		return nil, err
	}
	fmt.Printf("classHash %v\n", classHash.String())
	tx, err := provider.AddDeployAccountTransaction(ctx, rpcv02.BroadcastedDeployAccountTransaction{
		BroadcastedTxnCommonProperties: rpcv02.BroadcastedTxnCommonProperties{
			MaxFee:  big.NewInt(1),
			Version: "0x01",
			Nonce:   big.NewInt(2), // TODO: nonce handling
		},
		ContractAddressSalt: salt,
		ConstructorCalldata: inputs,
		ClassHash:           classHash,
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
	fmt.Println("a")
	// TODO: use UDC via account
	tx, err := provider.AddDeployTransaction(ctx, rpcv02.BroadcastedDeployTransaction{
		Version:             big.NewInt(1),
		ContractAddressSalt: salt,
		ConstructorCalldata: inputs,
		ContractClass:       class,
	})
	fmt.Println("b")
	if err != nil {
		return nil, err
	}
	status, err := provider.WaitForTransaction(ctx, types.HexToHash(tx.TransactionHash), 8*time.Second)
	fmt.Println("c")
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
