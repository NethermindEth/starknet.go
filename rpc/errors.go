package rpc

import "errors"

func tryUnwrapToRPCErr(err error, rpcErrors ...*RPCError) error {
	for _, rpcErr := range rpcErrors {
		if errors.Is(err, rpcErr) {
			return rpcErr
		}
	}

	return err
}

type RPCError struct {
	code    int
	message string
}

func (e *RPCError) Error() string {
	return e.message
}

func (e *RPCError) Code() int {
	return e.code
}

var (
	ErrFailedToReceiveTxn = &RPCError{
		code:    1,
		message: "Failed to write transaction",
	}
	ErrContractNotFound = &RPCError{
		code:    20,
		message: "Contract not found",
	}
	ErrBlockNotFound = &RPCError{
		code:    24,
		message: "Block not found",
	}
	ErrHashNotFound = &RPCError{
		code:    25,
		message: "Transaction hash not found",
	}
	ErrInvalidBlockHash = &RPCError{
		code:    24,
		message: "Invalid block hash",
	}
	ErrInvalidTxnIndex = &RPCError{
		code:    27,
		message: "Invalid transaction index in a block",
	}
	ErrClassHashNotFound = &RPCError{
		code:    28,
		message: "Class hash not found",
	}
	ErrPageSizeTooBig = &RPCError{
		code:    31,
		message: "Requested page size is too big",
	}
	ErrNoBlocks = &RPCError{
		code:    32,
		message: "There are no blocks",
	}
	ErrInvalidContinuationToken = &RPCError{
		code:    33,
		message: "The supplied continuation token is invalid or unknown",
	}
	ErrTooManyKeysInFilter = &RPCError{
		code:    34,
		message: "Too many keys provided in a filter",
	}
	ErrContractError = &RPCError{
		code:    40,
		message: "Contract error",
	}
	ErrInvalidContractClass = &RPCError{
		code:    50,
		message: "Invalid contract class",
	}
	ErrClassAlreadyDeclared = &RPCError{
		code:    51,
		message: "Class already declared",
	}
	ErrInvalidTransactionNonce = &RPCError{
		code:    52,
		message: "Invalid transaction nonce",
	}
	ErrInsufficientMaxFee = &RPCError{
		code:    53,
		message: "Max fee is smaller than the minimal transaction cost (validation plus fee transfer)",
	}
	ErrInsufficientAccountBalance = &RPCError{
		code:    54,
		message: "Account balance is smaller than the transaction's max_fee",
	}
	ErrValidationFailure = &RPCError{
		code:    55,
		message: "Account validation failed",
	}
	ErrCompilationFailed = &RPCError{
		code:    56,
		message: "Compilation failed",
	}
	ErrContractClassSizeTooLarge = &RPCError{
		code:    57,
		message: "Contract class size is too large",
	}
	ErrNonAccount = &RPCError{
		code:    58,
		message: "Sender address is not an account contract",
	}
	ErrDuplicateTx = &RPCError{
		code:    59,
		message: "A transaction with the same hash already exists in the mempool",
	}
	ErrCompiledClassHashMismatch = &RPCError{
		code:    60,
		message: "The compiled class hash did not match the one supplied in the transaction",
	}
	ErrUnsupportedTxVersion = &RPCError{
		code:    61,
		message: "The transaction version is not supported",
	}
	ErrUnsupportedContractClassVersion = &RPCError{
		code:    62,
		message: "The contract class version is not supported",
	}
	ErrUnexpectedError = &RPCError{
		code:    63,
		message: "An unexpected error occurred",
	}
)
