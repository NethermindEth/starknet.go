package rpc

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
	ErrInvalidMessageSelector = &RPCError{
		code:    21,
		message: "Invalid message selector",
	}
	ErrInvalidCallData = &RPCError{
		code:    22,
		message: "Invalid call data",
	}
	ErrBlockNotFound = &RPCError{
		code:    24,
		message: "Block not found",
	}
	ErrHashNotFound = &RPCError{
		code:    25,
		message: "Transaction hash not found",
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
)
