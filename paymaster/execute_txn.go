package paymaster

import (
	"context"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/client/rpcerr"
	"github.com/NethermindEth/starknet.go/typedata"
)

// ExecuteTransaction sends the signed typed data to the paymaster service for execution
//
// Parameters:
//   - ctx: The context.Context object for controlling the function call
//   - request: The signed typed data of the transaction to be executed by the paymaster service
//
// Returns:
//   - *ExecuteTransactionResponse: The hash of the transaction broadcasted by the paymaster and
//     the tracking ID corresponding to the user `execute` request
//   - error: An error if any error occurs
func (p *Paymaster) ExecuteTransaction(
	ctx context.Context,
	request *ExecuteTransactionRequest,
) (*ExecuteTransactionResponse, error) {
	var response ExecuteTransactionResponse
	if err := p.c.CallContextWithSliceArgs(ctx, &response, "paymaster_executeTransaction", request); err != nil {
		return nil, rpcerr.UnwrapToRPCErr(
			err,
			ErrInvalidAddress,
			ErrClassHashNotSupported,
			ErrInvalidDeploymentData,
			ErrInvalidSignature,
			ErrUnknownError,
			ErrMaxAmountTooLow,
			ErrTransactionExecutionError,
		)
	}

	return &response, nil
}

// ExecuteTransactionRequest is the request to execute a transaction via the paymaster (transaction + parameters).
type ExecuteTransactionRequest struct {
	// Typed data build by calling paymaster_buildTransaction signed by the
	// user to be executed by the paymaster service
	Transaction *ExecutableUserTransaction `json:"transaction"`
	// Execution parameters to be used when executing the transaction
	Parameters *UserParameters `json:"parameters"`
}

// ExecutableUserTransaction is a user transaction ready for execution (deploy, invoke, or both).
type ExecutableUserTransaction struct {
	// The type of the transaction to be executed by the paymaster
	Type UserTxnType `json:"type"`
	// The deployment data for the transaction, used for `deploy` and `deploy_and_invoke` transaction types.
	// Should be `nil` for `invoke` transaction types.
	Deployment *AccDeploymentData `json:"deployment,omitempty"`
	// Invoke data signed by the user to be executed by the paymaster service, used for`invoke` and
	// `deploy_and_invoke` transaction types.
	// Should be `nil` for `deploy` transaction types.
	Invoke *ExecutableUserInvoke `json:"invoke,omitempty"`
}

// ExecutableUserInvoke is a signed typed data of an invoke transaction ready to be executed by the paymaster service.
type ExecutableUserInvoke struct {
	// The address of the user account
	UserAddress *felt.Felt `json:"user_address"`
	// Typed data returned by the endpoint paymaster_buildTransaction
	TypedData *typedata.TypedData `json:"typed_data"`
	// Signature of the associated Typed Data
	Signature []*felt.Felt `json:"signature"`
}

// ExecuteTransactionResponse is the response from executing a transaction (tracking ID and transaction hash).
type ExecuteTransactionResponse struct {
	// A unique identifier used to track an execution request of a user. Its purpose is to track
	// possibly different transactions sent by the paymaster and which are associated with a same
	// user request. Such cases can happen during congestion, where a fee or tip bump may be needed
	// in order for a transaction to enter a block
	TrackingId      *felt.Felt `json:"tracking_id"`
	TransactionHash *felt.Felt `json:"transaction_hash"`
}
