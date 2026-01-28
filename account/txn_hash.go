package account

import (
	"fmt"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/hash"
	"github.com/NethermindEth/starknet.go/rpc/rpcv10"
	"github.com/NethermindEth/starknet.go/rpc/rpcv9"
)

// @changed now all these methods accept 'any' as the parameter
// TransactionHashDeployAccount calculates the transaction hash for a deploy
// account transaction.
//
// Parameters:
//   - tx: A pointer to a deploy account transaction to calculate the hash.
//   - contractAddress: The contract address as parameters as a *felt.Felt
//
// Returns:
//   - *felt.Felt: the calculated transaction hash
//   - error: an error if any
//
// Deprecated: This method will be removed soon. Use the functions available in
// the `hash` package instead.
func (account *Account) TransactionHashDeployAccount(
	tx any,
	contractAddress *felt.Felt,
) (*felt.Felt, error) {
	switch txn := tx.(type) {
	// deployAccTxn v1
	case *rpcv9.DeployAccountTxnV1:
		return hash.TransactionHashDeployAccountV1(txn, contractAddress, account.ChainID)
	case *rpcv10.DeployAccountTxnV1:
		return hash.TransactionHashDeployAccountV1(txn, contractAddress, account.ChainID)
	// deployAccTxn v3
	case *rpcv9.DeployAccountTxnV3:
		return hash.TransactionHashDeployAccountV3(txn, contractAddress, account.ChainID)
	case *rpcv10.DeployAccountTxnV3:
		return hash.TransactionHashDeployAccountV3(txn, contractAddress, account.ChainID)
	default:
		return nil, fmt.Errorf(
			"%w: got '%T' instead of a deploy account txn pointer",
			ErrTxnTypeUnSupported,
			txn,
		)
	}
}

// @todo move the tests to the hash package
// TransactionHashInvoke calculates the transaction hash for the given invoke
// transaction.
//
// Parameters:
//   - tx: A pointer to an invoke transaction to calculate the hash.
//
// Returns:
//   - *felt.Felt: The calculated transaction hash as a *felt.Felt
//   - error: an error, if any
//
// Deprecated: This method will be removed soon. Use the functions available in
// the `hash` package instead.
func (account *Account) TransactionHashInvoke(tx any) (*felt.Felt, error) {
	switch txn := tx.(type) {
	// invoke v0
	case *rpcv9.InvokeTxnV0:
		return hash.TransactionHashInvokeV0(txn, account.ChainID)
	case *rpcv10.InvokeTxnV0:
		return hash.TransactionHashInvokeV0(txn, account.ChainID)
	// invoke v1
	case *rpcv9.InvokeTxnV1:
		return hash.TransactionHashInvokeV1(txn, account.ChainID)
	case *rpcv10.InvokeTxnV1:
		return hash.TransactionHashInvokeV1(txn, account.ChainID)
	// invoke v3
	case *rpcv9.InvokeTxnV3:
		return hash.TransactionHashInvokeV3(txn, account.ChainID)
	case *rpcv10.InvokeTxnV3:
		return hash.TransactionHashInvokeV3(txn, account.ChainID)
	default:
		return nil, fmt.Errorf(
			"%w: got '%T' instead of an invoke txn pointer",
			ErrTxnTypeUnSupported,
			txn,
		)
	}
}

// TransactionHashDeclare calculates the transaction hash for declaring a
// transaction type.
//
// Parameters:
//   - tx: A pointer to a declare transaction to calculate the hash.
//
// Returns:
//   - *felt.Felt: the calculated transaction hash as `*felt.Felt` value
//   - error: an error, if any
//
// Deprecated: This method will be removed soon. Use the functions available in
// the `hash` package instead.
func (account *Account) TransactionHashDeclare(tx any) (*felt.Felt, error) {
	switch txn := tx.(type) {
	// declare v0
	case *rpcv9.DeclareTxnV0:
		return hash.TransactionHashDeclareV0(txn, account.ChainID)
	case *rpcv10.DeclareTxnV0:
		return hash.TransactionHashDeclareV0(txn, account.ChainID)
	// declare v1
	case *rpcv9.DeclareTxnV1:
		return hash.TransactionHashDeclareV1(txn, account.ChainID)
	case *rpcv10.DeclareTxnV1:
		return hash.TransactionHashDeclareV1(txn, account.ChainID)
	// declare v2
	case *rpcv9.DeclareTxnV2:
		return hash.TransactionHashDeclareV2(txn, account.ChainID)
	case *rpcv10.DeclareTxnV2:
		return hash.TransactionHashDeclareV2(txn, account.ChainID)
	// declare v3
	case *rpcv9.DeclareTxnV3:
		return hash.TransactionHashDeclareV3(txn, account.ChainID, txn.ClassHash)
	case *rpcv10.DeclareTxnV3:
		return hash.TransactionHashDeclareV3(txn, account.ChainID, txn.ClassHash)
	// broadcast declare v3
	case *rpcv9.BroadcastDeclareTxnV3:
		return hash.TransactionHashDeclareV3(
			txn, account.ChainID, hash.ClassHash(txn.ContractClass))
	case *rpcv10.BroadcastDeclareTxnV3:
		return hash.TransactionHashDeclareV3(
			txn, account.ChainID, hash.ClassHash(txn.ContractClass))
	default:
		return nil, fmt.Errorf(
			"%w: got '%T' instead of a declare txn pointer",
			ErrTxnTypeUnSupported,
			txn,
		)
	}
}
