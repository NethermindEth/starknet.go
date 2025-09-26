package paymaster

import (
	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/client/rpcerr"
)

// aliases to facilitate usage

type (
	RPCError      = rpcerr.RPCError
	StringErrData = rpcerr.StringErrData
)

// Paymaster-specific errors based on SNIP-29 specification
//
//nolint:exhaustruct
var (
	ErrInvalidAddress = &RPCError{
		Code:    150,
		Message: "An error occurred (INVALID_ADDRESS)",
	}

	ErrTokenNotSupported = &RPCError{
		Code:    151,
		Message: "An error occurred (TOKEN_NOT_SUPPORTED)",
	}

	ErrInvalidSignature = &RPCError{
		Code:    153,
		Message: "An error occurred (INVALID_SIGNATURE)",
	}

	ErrMaxAmountTooLow = &RPCError{
		Code:    154,
		Message: "An error occurred (MAX_AMOUNT_TOO_LOW)",
	}

	ErrClassHashNotSupported = &RPCError{
		Code:    155,
		Message: "An error occurred (CLASS_HASH_NOT_SUPPORTED)",
	}

	ErrTransactionExecutionError = &RPCError{
		Code:    156,
		Message: "An error occurred (TRANSACTION_EXECUTION_ERROR)",
		Data:    &TxnExecutionErrData{},
	}

	ErrInvalidTimeBounds = &RPCError{
		Code:    157,
		Message: "An error occurred (INVALID_TIME_BOUNDS)",
	}

	ErrInvalidDeploymentData = &RPCError{
		Code:    158,
		Message: "An error occurred (INVALID_DEPLOYMENT_DATA)",
	}

	ErrInvalidClassHash = &RPCError{
		Code:    159,
		Message: "An error occurred (INVALID_ADDRESS)",
	}

	ErrInvalidID = &RPCError{
		Code:    160,
		Message: "An error occurred (INVALID_ID)",
	}

	ErrUnknownError = &RPCError{
		Code:    163,
		Message: "An error occurred (UNKNOWN_ERROR)",
		Data:    StringErrData(""),
	}
)

// TxnExecutionErrData represents the structured data for TRANSACTION_EXECUTION_ERROR
type TxnExecutionErrData struct {
	ExecutionError string `json:"execution_error"`
}

// ErrorMessage implements the RPCData interface
func (t TxnExecutionErrData) ErrorMessage() string {
	return t.ExecutionError
}

// Object containing data about the token: contract address, number of decimals and current price in STRK
type TokenData struct {
	// Token contract address
	TokenAddress *felt.Felt `json:"token_address"`
	// The number of decimals of the token
	Decimals uint8 `json:"decimals"`
	// Price in STRK (in FRI units)
	PriceInStrk string `json:"price_in_strk"` // u256 as a hex string
}
