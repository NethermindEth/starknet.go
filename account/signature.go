package account

import (
	"context"
	"errors"
	"fmt"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/curve"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/NethermindEth/starknet.go/rpc/rpcv10"
	"github.com/NethermindEth/starknet.go/rpc/rpcv9"
)

// Sign signs the given felt message using the account's private key.
//
// Parameters:
//   - ctx: is the context used for the signing operation
//   - msg: is the felt message to be signed
//
// Returns:
//   - []*felt.Felt: an array of signed felt messages
//   - error: an error, if any
func (account *Account) Sign(ctx context.Context, msg *felt.Felt) ([]*felt.Felt, error) {
	msgBig := internalUtils.FeltToBigInt(msg)

	s1, s2, err := account.ks.Sign(ctx, account.publicKey, msgBig)
	if err != nil {
		return nil, err
	}
	s1Felt := internalUtils.BigIntToFelt(s1)
	s2Felt := internalUtils.BigIntToFelt(s2)

	return []*felt.Felt{s1Felt, s2Felt}, nil
}

// SignInvokeTransaction signs and invokes a transaction.
//
// Parameters:
//   - ctx: the context.Context for the function execution.
//   - invokeTx: a pointer to the invoke transaction to be signed.
//
// Returns:
//   - error: an error if any
func (account *Account) SignInvokeTransaction(
	ctx context.Context,
	invokeTx any,
) error {
	txHash, err := account.TransactionHashInvoke(invokeTx)
	if err != nil {
		return err
	}
	signature, err := account.Sign(ctx, txHash)
	if err != nil {
		return err
	}

	switch invoke := invokeTx.(type) {
	// invoke v0
	case *rpcv9.InvokeTxnV0:
		invoke.Signature = signature
	case *rpcv10.InvokeTxnV0:
		invoke.Signature = signature
	// invoke v1
	case *rpcv9.InvokeTxnV1:
		invoke.Signature = signature
	case *rpcv10.InvokeTxnV1:
		invoke.Signature = signature
	// invoke v3
	case *rpcv9.InvokeTxnV3:
		invoke.Signature = signature
	case *rpcv10.InvokeTxnV3:
		invoke.Signature = signature
	default:
		return fmt.Errorf(
			"invalid invoke txn of type %T, did you pass a valid invoke txn pointer?",
			invoke,
		)
	}

	return nil
}

// SignDeployAccountTransaction signs a deploy account transaction.
//
// Parameters:
//   - ctx: the context.Context for the function execution
//   - tx: a pointer to the deploy account transaction to be signed
//   - precomputeAddress: the precomputed address for the transaction
//
// Returns:
//   - error: an error if any
func (account *Account) SignDeployAccountTransaction(
	ctx context.Context,
	tx any,
	precomputeAddress *felt.Felt,
) error {
	txHash, err := account.TransactionHashDeployAccount(tx, precomputeAddress)
	if err != nil {
		return err
	}
	signature, err := account.Sign(ctx, txHash)
	if err != nil {
		return err
	}

	switch deployAcc := tx.(type) {
	// deployAcc v1
	case *rpcv9.DeployAccountTxnV1:
		deployAcc.Signature = signature
	case *rpcv10.DeployAccountTxnV1:
		deployAcc.Signature = signature
	// deployAcc v3
	case *rpcv9.DeployAccountTxnV3:
		deployAcc.Signature = signature
	case *rpcv10.DeployAccountTxnV3:
		deployAcc.Signature = signature
	default:
		return fmt.Errorf(
			"invalid deploy account txn of type %T, did you pass a valid deploy account txn pointer?",
			deployAcc,
		)
	}

	return nil
}

// SignDeclareTransaction signs a declare transaction using the provided Account.
//
// Parameters:
//   - ctx: the context.Context
//   - tx: the pointer to a Declare or BroadcastDeclare txn
//
// Returns:
//   - error: an error if any
func (account *Account) SignDeclareTransaction(ctx context.Context, tx any) error {
	txHash, err := account.TransactionHashDeclare(tx)
	if err != nil {
		return err
	}
	signature, err := account.Sign(ctx, txHash)
	if err != nil {
		return err
	}

	switch declare := tx.(type) {
	// declare v1
	case *rpcv9.DeclareTxnV1:
		declare.Signature = signature
	case *rpcv10.DeclareTxnV1:
		declare.Signature = signature
	// declare v2
	case *rpcv9.DeclareTxnV2:
		declare.Signature = signature
	case *rpcv10.DeclareTxnV2:
		declare.Signature = signature
	// declare v3
	case *rpcv9.DeclareTxnV3:
		declare.Signature = signature
	case *rpcv10.DeclareTxnV3:
		declare.Signature = signature
	// broadcast declare v3
	case *rpcv9.BroadcastDeclareTxnV3:
		declare.Signature = signature
	case *rpcv10.BroadcastDeclareTxnV3:
		declare.Signature = signature
	default:
		return fmt.Errorf(
			"invalid declare txn of type %T, did you pass a valid declare txn pointer?",
			declare,
		)
	}

	return nil
}

// Verifies the validity of the signature for a given message hash using the
// account's public key.
//
// Parameters:
//   - msgHash: The message hash to be verified
//   - signature: A slice of felt.Felt containing the two signature components
//
// Returns:
//   - bool: true if the signature is valid, false otherwise
//   - error: An error if any occurred during the verification process
func (account *Account) Verify(msgHash *felt.Felt, signature []*felt.Felt) (bool, error) {
	publicKeyFelt, err := new(felt.Felt).SetString(account.publicKey)
	if err != nil {
		return false, errors.Join(errors.New("failed to convert public key to felt"), err)
	}

	return curve.VerifyFelts(msgHash, signature[0], signature[1], publicKeyFelt)
}
