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
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/types"
	"github.com/NethermindEth/starknet.go/utils"
)

type RPCProvider rpc.Provider

func (p *RPCProvider) declareAndWaitWithWallet(ctx context.Context, compiledClass []byte) (*DeclareOutput, error) {
	provider := rpc.Provider(*p)
	class := rpc.DeprecatedContractClass{}
	if err := json.Unmarshal(compiledClass, &class); err != nil {
		return nil, err
	}
	SenderAddress, err := utils.HexToFelt("0x01")
	if err != nil {
		return nil, err
	}
	tx, err := provider.AddDeclareTransaction(ctx, rpc.BroadcastedDeclareTransaction{
		BroadcastedTxnCommonProperties: rpc.BroadcastedTxnCommonProperties{
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

func (p *RPCProvider) deployAccountAndWaitNoWallet(ctx context.Context, classHash *felt.Felt, compiledClass []byte, salt string, inputs []string) (*DeployOutput, error) {
	provider := rpc.Provider(*p)
	class := rpc.DeprecatedContractClass{}
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

	tx, err := provider.AddDeployAccountTransaction(ctx, rpc.BroadcastedDeployAccountTransaction{
		BroadcastedTxnCommonProperties: rpc.BroadcastedTxnCommonProperties{
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

func (p *RPCProvider) deployAndWaitWithWallet(ctx context.Context, compiledClass []byte, salt string, inputs []string) (*DeployOutput, error) {
	provider := rpc.Provider(*p)
	class := rpc.DeprecatedContractClass{}
	if err := json.Unmarshal(compiledClass, &class); err != nil {
		return nil, err
	}
	saltFelt, err := utils.HexToFelt(salt)
	if err != nil {
		return nil, err
	}
	inputsFelt, err := utils.HexArrToFelt(inputs)
	if err != nil {
		return nil, err
	}

	// TODO: use UDC via account
	tx, err := provider.AddDeployAccountTransaction(ctx, rpc.BroadcastedDeployAccountTransaction{
		BroadcastedTxnCommonProperties: rpc.BroadcastedTxnCommonProperties{
			Version: rpc.TransactionV1,
		},
		ContractAddressSalt: saltFelt,
		ConstructorCalldata: inputsFelt,
	})

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
