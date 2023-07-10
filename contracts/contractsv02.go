package contracts

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/rpcv02"
	"github.com/NethermindEth/starknet.go/types"
	"github.com/NethermindEth/starknet.go/utils"
)

type RPCv02Provider rpcv02.Provider

func (p *RPCv02Provider) declareAndWaitWithWallet(ctx context.Context, compiledClass []byte) (*DeclareOutput, error) {
	provider := rpcv02.Provider(*p)
	class := rpcv02.ContractClass{}
	if err := json.Unmarshal(compiledClass, &class); err != nil {
		return nil, err
	}
	SenderAddress, err := utils.HexToFelt("0x01")
	if err != nil {
		return nil, err
	}
	tx, err := provider.AddDeclareTransaction(ctx, rpcv02.BroadcastedDeclareTransaction{
		BroadcastedTxnCommonProperties: rpcv02.BroadcastedTxnCommonProperties{
			Type:    "DECLARE",
			MaxFee:  new(felt.Felt).SetUint64(10000),
			Version: "0x01",
			Nonce:   new(felt.Felt).SetUint64(1), // TODO: nonce handling
		},
		ContractClass: class,
		SenderAddress: SenderAddress, // TODO: constant devnet address
	})
	if err != nil {
		return nil, err
	}

	status, err := provider.WaitForTransaction(ctx, tx.TransactionHash, 8*time.Second)
	if err != nil {
		log.Printf("class Hash: %s\n", tx.ClassHash)
		log.Printf("transaction Hash: %s\n", tx.TransactionHash)
		return nil, err
	}
	if types.TransactionState(status.String()) == types.TransactionRejected {
		log.Printf("class Hash: %s\n", tx.ClassHash)
		log.Printf("transaction Hash: %s\n", tx.TransactionHash)
		return nil, errors.New("declare rejected")
	}
	return &DeclareOutput{
		classHash:       tx.ClassHash.String(),
		transactionHash: tx.TransactionHash.String(),
	}, nil
}

func (p *RPCv02Provider) deployAccountAndWaitNoWallet(ctx context.Context, classHash *felt.Felt, compiledClass []byte, salt string, inputs []string) (*DeployOutput, error) {
	provider := rpcv02.Provider(*p)
	class := rpcv02.ContractClass{}
	if err := json.Unmarshal(compiledClass, &class); err != nil {
		return nil, err
	}

	fmt.Printf("classHash %v\n", classHash.String())

	saltFelt, err := utils.HexToFelt(salt)
	if err != nil {
		return nil, err
	}
	inputsFelt, err := utils.HexArrToFelt(inputs)
	if err != nil {
		return nil, err
	}

	tx, err := provider.AddDeployAccountTransaction(ctx, rpcv02.BroadcastedDeployAccountTransaction{
		BroadcastedTxnCommonProperties: rpcv02.BroadcastedTxnCommonProperties{
			MaxFee:  new(felt.Felt).SetUint64(1),
			Version: "0x01",
			Nonce:   new(felt.Felt).SetUint64(2), // TODO: nonce handling
		},
		ContractAddressSalt: saltFelt,
		ConstructorCalldata: inputsFelt,
		ClassHash:           classHash,
	})
	if err != nil {
		return nil, err
	}

	status, err := provider.WaitForTransaction(ctx, tx.TransactionHash, 8*time.Second)
	if err != nil {
		log.Printf("contract Address: %s\n", tx.ContractAddress)
		log.Printf("transaction Hash: %s\n", tx.TransactionHash)
		return nil, err
	}
	if types.TransactionState(status.String()) == types.TransactionRejected {
		log.Printf("contract Address: %s\n", tx.ContractAddress)
		log.Printf("transaction Hash: %s\n", tx.TransactionHash)
		return nil, errors.New("deploy rejected")
	}
	return &DeployOutput{
		ContractAddress: tx.ContractAddress.String(),
		TransactionHash: tx.TransactionHash.String(),
	}, nil
}

func (p *RPCv02Provider) deployAndWaitWithWallet(ctx context.Context, compiledClass []byte, salt string, inputs []string) (*DeployOutput, error) {
	provider := rpcv02.Provider(*p)
	class := rpcv02.ContractClass{}
	if err := json.Unmarshal(compiledClass, &class); err != nil {
		return nil, err
	}
	fmt.Println("a")

	saltFelt, err := utils.HexToFelt(salt)
	if err != nil {
		return nil, err
	}
	inputsFelt, err := utils.HexArrToFelt(inputs)
	if err != nil {
		return nil, err
	}

	// TODO: use UDC via account
	tx, err := provider.AddDeployTransaction(ctx, rpcv02.BroadcastedDeployTxn{
		DeployTransactionProperties: rpcv02.DeployTransactionProperties{
			Version:             rpcv02.TransactionV1,
			ContractAddressSalt: saltFelt,
			ConstructorCalldata: inputsFelt,
		},
		ContractClass: class,
	})
	fmt.Println("b")
	if err != nil {
		return nil, err
	}

	status, err := provider.WaitForTransaction(ctx, tx.TransactionHash, 8*time.Second)
	fmt.Println("c")
	if err != nil {
		log.Printf("contract Address: %s\n", tx.ContractAddress)
		log.Printf("transaction Hash: %s\n", tx.TransactionHash)
		return nil, err
	}
	if types.TransactionState(status.String()) == types.TransactionRejected {
		log.Printf("contract Address: %s\n", tx.ContractAddress)
		log.Printf("transaction Hash: %s\n", tx.TransactionHash)
		return nil, errors.New("deploy rejected")
	}
	return &DeployOutput{
		ContractAddress: tx.ContractAddress.String(),
		TransactionHash: tx.TransactionHash.String(),
	}, nil
}
