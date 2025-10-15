package account

import (
	"fmt"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/hash"
	"github.com/NethermindEth/starknet.go/rpc"
)

// TransactionHashDeployAccount calculates the transaction hash for a deploy
// account transaction.
//
// Parameters:
//   - tx: The deploy account transaction to calculate the hash for. Can be of
//     type DeployAccountTxn or DeployAccountTxnV3.
//   - contractAddress: The contract address as parameters as a *felt.Felt
//
// Returns:
//   - *felt.Felt: the calculated transaction hash
//   - error: an error if any
func (account *Account) TransactionHashDeployAccount(
	tx rpc.DeployAccountType,
	contractAddress *felt.Felt,
) (*felt.Felt, error) {
	//nolint:lll // The link would be unclickable if we break the line.
	// https://docs.starknet.io/architecture-and-concepts/network-architecture/transactions/#deploy_account_transaction
	switch txn := tx.(type) {
	// deployAccTxn v1, pointer and struct
	case *rpc.DeployAccountTxnV1:
		return hash.TransactionHashDeployAccountV1(txn, contractAddress, account.ChainID)
	case rpc.DeployAccountTxnV1:
		return hash.TransactionHashDeployAccountV1(&txn, contractAddress, account.ChainID)
	// deployAccTxn v3, pointer and struct
	case *rpc.DeployAccountTxnV3:
		return hash.TransactionHashDeployAccountV3(txn, contractAddress, account.ChainID)
	case rpc.DeployAccountTxnV3:
		return hash.TransactionHashDeployAccountV3(&txn, contractAddress, account.ChainID)
	default:
		return nil, fmt.Errorf(
			"%w: got '%T' instead of a valid invoke txn type",
			ErrTxnTypeUnSupported,
			txn,
		)
	}
}

// TransactionHashInvoke calculates the transaction hash for the given invoke
// transaction.
//
// Parameters:
//   - tx: The invoke transaction to calculate the hash for. Can be of type
//     InvokeTxnV0, InvokeTxnV1, or InvokeTxnV3.
//
// Returns:
//   - *felt.Felt: The calculated transaction hash as a *felt.Felt
//   - error: an error, if any
//
// If the transaction type is unsupported, the function returns an error.
func (account *Account) TransactionHashInvoke(tx rpc.InvokeTxnType) (*felt.Felt, error) {
	switch txn := tx.(type) {
	// invoke v0, pointer and struct
	case *rpc.InvokeTxnV0:
		return hash.TransactionHashInvokeV0(txn, account.ChainID)
	case rpc.InvokeTxnV0:
		return hash.TransactionHashInvokeV0(&txn, account.ChainID)
	// invoke v1, pointer and struct
	case *rpc.InvokeTxnV1:
		return hash.TransactionHashInvokeV1(txn, account.ChainID)
	case rpc.InvokeTxnV1:
		return hash.TransactionHashInvokeV1(&txn, account.ChainID)
	// invoke v3, pointer and struct
	case *rpc.InvokeTxnV3:
		return hash.TransactionHashInvokeV3(txn, account.ChainID)
	case rpc.InvokeTxnV3:
		return hash.TransactionHashInvokeV3(&txn, account.ChainID)
	default:
		return nil, fmt.Errorf(
			"%w: got '%T' instead of a valid invoke txn type",
			ErrTxnTypeUnSupported,
			txn,
		)
	}
}

// TransactionHashDeclare calculates the transaction hash for declaring a
// transaction type.
//
// Parameters:
//   - tx: The `tx` parameter of type `rpc.DeclareTxnType`. Can be one of the
//     types DeclareTxnV1/V2/V3, and BroadcastDeclareTxnV3
//
// Returns:
//   - *felt.Felt: the calculated transaction hash as `*felt.Felt` value
//   - error: an error, if any
//
// If the `tx` parameter is not one of the supported types, the function returns
// an error `ErrTxnTypeUnSupported`.
func (account *Account) TransactionHashDeclare(tx rpc.DeclareTxnType) (*felt.Felt, error) {
	switch txn := tx.(type) {
	// Due to inconsistencies in version 0 hash calculation we don't calculate the hash
	case *rpc.DeclareTxnV0, rpc.DeclareTxnV0:
		return nil, ErrTxnVersionUnSupported
	// declare v1, pointer and struct
	case *rpc.DeclareTxnV1:
		return hash.TransactionHashDeclareV1(txn, account.ChainID)
	case rpc.DeclareTxnV1:
		return hash.TransactionHashDeclareV1(&txn, account.ChainID)
	// declare v2, pointer and struct
	case *rpc.DeclareTxnV2:
		return hash.TransactionHashDeclareV2(txn, account.ChainID)
	case rpc.DeclareTxnV2:
		return hash.TransactionHashDeclareV2(&txn, account.ChainID)
	// declare v3, pointer and struct
	case *rpc.DeclareTxnV3:
		return hash.TransactionHashDeclareV3(txn, account.ChainID)
	case rpc.DeclareTxnV3:
		return hash.TransactionHashDeclareV3(&txn, account.ChainID)
	// broadcast declare v3, pointer and struct
	case *rpc.BroadcastDeclareTxnV3:
		return hash.TransactionHashBroadcastDeclareV3(txn, account.ChainID)
	case rpc.BroadcastDeclareTxnV3:
		return hash.TransactionHashBroadcastDeclareV3(&txn, account.ChainID)
	default:
		return nil, fmt.Errorf(
			"%w: got '%T' instead of a valid declare txn type",
			ErrTxnTypeUnSupported,
			txn,
		)
	}
}
