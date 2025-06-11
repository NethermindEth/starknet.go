package account

import (
	"context"
	"errors"
	"fmt"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/curve"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/NethermindEth/starknet.go/rpc"
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
//   - invokeTx: the InvokeTxnV3 pointer representing the transaction to be invoked.
//
// Returns:
//   - error: an error if there was an error in the signing or invoking process
func (account *Account) SignInvokeTransaction(ctx context.Context, invokeTx rpc.InvokeTxnType) error {
	switch invoke := invokeTx.(type) {
	case *rpc.InvokeTxnV0:
		signature, err := signInvokeTransaction(ctx, account, invoke)
		if err != nil {
			return err
		}
		invoke.Signature = signature
	case *rpc.InvokeTxnV1:
		signature, err := signInvokeTransaction(ctx, account, invoke)
		if err != nil {
			return err
		}
		invoke.Signature = signature
	case *rpc.InvokeTxnV3:
		signature, err := signInvokeTransaction(ctx, account, invoke)
		if err != nil {
			return err
		}
		invoke.Signature = signature
	default:
		return fmt.Errorf("invalid invoke txn of type %T, did you pass a valid invoke txn pointer?", invoke)
	}

	return nil
}

// signInvokeTransaction is a generic helper function that signs an invoke transaction.
func signInvokeTransaction[T rpc.InvokeTxnType](ctx context.Context, account *Account, invokeTx *T) ([]*felt.Felt, error) {
	txHash, err := account.TransactionHashInvoke(*invokeTx)
	if err != nil {
		return nil, err
	}
	signature, err := account.Sign(ctx, txHash)
	if err != nil {
		return nil, err
	}

	return signature, nil
}

// SignDeployAccountTransaction signs a deploy account transaction.
//
// Parameters:
//   - ctx: the context.Context for the function execution
//   - tx: the *rpc.DeployAccountTxnV3 pointer representing the transaction to be signed
//   - precomputeAddress: the precomputed address for the transaction
//
// Returns:
//   - error: an error if any
func (account *Account) SignDeployAccountTransaction(ctx context.Context, tx rpc.DeployAccountType, precomputeAddress *felt.Felt) error {
	switch deployAcc := tx.(type) {
	case *rpc.DeployAccountTxnV1:
		signature, err := signDeployAccountTransaction(ctx, account, deployAcc, precomputeAddress)
		if err != nil {
			return err
		}
		deployAcc.Signature = signature
	case *rpc.DeployAccountTxnV3:
		signature, err := signDeployAccountTransaction(ctx, account, deployAcc, precomputeAddress)
		if err != nil {
			return err
		}
		deployAcc.Signature = signature
	default:
		return fmt.Errorf("invalid deploy account txn of type %T, did you pass a valid deploy account txn pointer?", deployAcc)
	}

	return nil
}

// signDeployAccountTransaction is a generic helper function that signs a deploy account transaction.
func signDeployAccountTransaction[T rpc.DeployAccountType](
	ctx context.Context,
	account *Account,
	tx *T,
	precomputeAddress *felt.Felt,
) ([]*felt.Felt, error) {
	txHash, err := account.TransactionHashDeployAccount(*tx, precomputeAddress)
	if err != nil {
		return nil, err
	}
	signature, err := account.Sign(ctx, txHash)
	if err != nil {
		return nil, err
	}

	return signature, nil
}

// SignDeclareTransaction signs a declare transaction using the provided Account.
//
// Parameters:
//   - ctx: the context.Context
//   - tx: the pointer to a Declare or BroadcastDeclare txn
//
// Returns:
//   - error: an error if any
func (account *Account) SignDeclareTransaction(ctx context.Context, tx rpc.DeclareTxnType) error {
	switch declare := tx.(type) {
	case *rpc.DeclareTxnV1:
		signature, err := signDeclareTransaction(ctx, account, declare)
		if err != nil {
			return err
		}
		declare.Signature = signature
	case *rpc.DeclareTxnV2:
		signature, err := signDeclareTransaction(ctx, account, declare)
		if err != nil {
			return err
		}
		declare.Signature = signature
	case *rpc.DeclareTxnV3:
		signature, err := signDeclareTransaction(ctx, account, declare)
		if err != nil {
			return err
		}
		declare.Signature = signature
	case *rpc.BroadcastDeclareTxnV3:
		signature, err := signDeclareTransaction(ctx, account, declare)
		if err != nil {
			return err
		}
		declare.Signature = signature
	default:
		return fmt.Errorf("invalid declare txn of type %T, did you pass a valid declare txn pointer?", declare)
	}

	return nil
}

// signDeclareTransaction is a generic helper function that signs a declare transaction.
func signDeclareTransaction[T rpc.DeclareTxnType](ctx context.Context, account *Account, tx *T) ([]*felt.Felt, error) {
	txHash, err := account.TransactionHashDeclare(*tx)
	if err != nil {
		return nil, err
	}
	signature, err := account.Sign(ctx, txHash)
	if err != nil {
		return nil, err
	}

	return signature, nil
}

// Verifies the validity of the signature for a given message hash using the account's public key.
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
