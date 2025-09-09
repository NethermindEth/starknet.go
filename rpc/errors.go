package rpc

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/rpcerr"
)

//nolint:exhaustruct
var (
	_ rpcerr.RPCData = rpcerr.StringErrData("")

	_ rpcerr.RPCData = &CompilationErrData{}
	_ rpcerr.RPCData = &ContractErrData{}
	_ rpcerr.RPCData = &TransactionExecErrData{}
	_ rpcerr.RPCData = &TraceStatusErrData{}
)

//nolint:exhaustruct
var (
	ErrFailedToReceiveTxn = &rpcerr.RPCError{
		Code:    1,
		Message: "Failed to write transaction",
	}
	ErrNoTraceAvailable = &rpcerr.RPCError{
		Code:    10,
		Message: "No trace available for transaction",
		Data:    &TraceStatusErrData{},
	}
	ErrContractNotFound = &rpcerr.RPCError{
		Code:    20,
		Message: "Contract not found",
	}
	ErrEntrypointNotFound = &rpcerr.RPCError{
		Code:    21,
		Message: "Requested entrypoint does not exist in the contract",
	}
	ErrBlockNotFound = &rpcerr.RPCError{
		Code:    24,
		Message: "Block not found",
	}
	ErrInvalidTxnHash = &rpcerr.RPCError{
		Code:    25,
		Message: "Invalid transaction hash",
	}
	ErrInvalidBlockHash = &rpcerr.RPCError{
		Code:    26,
		Message: "Invalid block hash",
	}
	ErrInvalidTxnIndex = &rpcerr.RPCError{
		Code:    27,
		Message: "Invalid transaction index in a block",
	}
	ErrClassHashNotFound = &rpcerr.RPCError{
		Code:    28,
		Message: "Class hash not found",
	}
	ErrHashNotFound = &rpcerr.RPCError{
		Code:    29,
		Message: "Transaction hash not found",
	}
	ErrPageSizeTooBig = &rpcerr.RPCError{
		Code:    31,
		Message: "Requested page size is too big",
	}
	ErrNoBlocks = &rpcerr.RPCError{
		Code:    32,
		Message: "There are no blocks",
	}
	ErrInvalidContinuationToken = &rpcerr.RPCError{
		Code:    33,
		Message: "The supplied continuation token is invalid or unknown",
	}
	ErrTooManyKeysInFilter = &rpcerr.RPCError{
		Code:    34,
		Message: "Too many keys provided in a filter",
	}
	ErrContractError = &rpcerr.RPCError{
		Code:    40,
		Message: "Contract error",
		Data:    &ContractErrData{},
	}
	ErrTxnExec = &rpcerr.RPCError{
		Code:    41,
		Message: "Transaction execution error",
		Data:    &TransactionExecErrData{},
	}
	ErrStorageProofNotSupported = &rpcerr.RPCError{
		Code:    42,
		Message: "the node doesn't support storage proofs for blocks that are too far in the past",
	}
	ErrInvalidContractClass = &rpcerr.RPCError{
		Code:    50,
		Message: "Invalid contract class",
	}
	ErrClassAlreadyDeclared = &rpcerr.RPCError{
		Code:    51,
		Message: "Class already declared",
	}
	ErrInvalidTransactionNonce = &rpcerr.RPCError{
		Code:    52,
		Message: "Invalid transaction nonce",
		Data:    rpcerr.StringErrData(""),
	}
	ErrInsufficientResourcesForValidate = &rpcerr.RPCError{
		Code:    53,
		Message: "The transaction's resources don't cover validation or the minimal transaction fee",
	}
	ErrInsufficientAccountBalance = &rpcerr.RPCError{
		Code:    54,
		Message: "Account balance is smaller than the transaction's maximal fee (calculated as the sum of each resource's limit x max price)",
	}
	ErrValidationFailure = &rpcerr.RPCError{
		Code:    55,
		Message: "Account validation failed",
		Data:    rpcerr.StringErrData(""),
	}
	ErrCompilationFailed = &rpcerr.RPCError{
		Code:    56,
		Message: "Compilation failed",
		Data:    rpcerr.StringErrData(""),
	}
	ErrContractClassSizeTooLarge = &rpcerr.RPCError{
		Code:    57,
		Message: "Contract class size is too large",
	}
	ErrNonAccount = &rpcerr.RPCError{
		Code:    58,
		Message: "Sender address is not an account contract",
	}
	ErrDuplicateTx = &rpcerr.RPCError{
		Code:    59,
		Message: "A transaction with the same hash already exists in the mempool",
	}
	ErrCompiledClassHashMismatch = &rpcerr.RPCError{
		Code:    60,
		Message: "The compiled class hash did not match the one supplied in the transaction",
	}
	ErrUnsupportedTxVersion = &rpcerr.RPCError{
		Code:    61,
		Message: "The transaction version is not supported",
	}
	ErrUnsupportedContractClassVersion = &rpcerr.RPCError{
		Code:    62,
		Message: "The contract class version is not supported",
	}
	ErrUnexpectedError = &rpcerr.RPCError{
		Code:    63,
		Message: "An unexpected error occurred",
		Data:    rpcerr.StringErrData(""),
	}
	ErrReplacementTransactionUnderpriced = &rpcerr.RPCError{
		Code:    64,
		Message: "Replacement transaction is underpriced",
	}
	ErrFeeBelowMinimum = &rpcerr.RPCError{
		Code:    65,
		Message: "Transaction fee below minimum",
	}
	ErrInvalidSubscriptionID = &rpcerr.RPCError{
		Code:    66,
		Message: "Invalid subscription id",
	}
	ErrTooManyAddressesInFilter = &rpcerr.RPCError{
		Code:    67,
		Message: "Too many addresses in filter sender_address filter",
	}
	ErrTooManyBlocksBack = &rpcerr.RPCError{
		Code:    68,
		Message: "Cannot go back more than 1024 blocks",
	}
	ErrCompilationError = &rpcerr.RPCError{
		Code:    100,
		Message: "Failed to compile the contract",
		Data:    &CompilationErrData{},
	}
)

// Structured type for the ErrCompilationError data
type CompilationErrData struct {
	CompilationError string `json:"compilation_error"`
}

func (c *CompilationErrData) ErrorMessage() string {
	return c.CompilationError
}

// Structured type for the ErrContractError data
type ContractErrData struct {
	RevertError ContractExecutionError `json:"revert_error"`
}

func (c *ContractErrData) ErrorMessage() string {
	return c.RevertError.Message
}

// Structured type for the ErrTransactionExecError data
type TransactionExecErrData struct {
	TransactionIndex int                    `json:"transaction_index"`
	ExecutionError   ContractExecutionError `json:"execution_error"`
}

func (t *TransactionExecErrData) ErrorMessage() string {
	return t.ExecutionError.Message
}

// Structured type for the ErrTraceStatusError data
type TraceStatusErrData struct {
	Status TraceStatus `json:"status"`
}

func (t *TraceStatusErrData) ErrorMessage() string {
	return string(t.Status)
}

// structured error that can later be processed by wallets or sdks
type ContractExecutionError struct {
	// the error raised during execution
	Message              string
	ContractExecErrInner *ContractExecutionErrorInner
}

func (contractEx *ContractExecutionError) UnmarshalJSON(data []byte) error {
	var message string

	if err := json.Unmarshal(data, &message); err == nil {
		*contractEx = ContractExecutionError{
			Message:              message,
			ContractExecErrInner: nil,
		}

		return nil
	}

	var contractErrStruct ContractExecutionErrorInner

	if err := json.Unmarshal(data, &contractErrStruct); err == nil {
		message := fmt.Sprintf("Contract address= %s, Class hash= %s, Selector= %s, Nested error: ",
			contractErrStruct.ContractAddress,
			contractErrStruct.ClassHash,
			contractErrStruct.Selector,
		)

		*contractEx = ContractExecutionError{
			Message:              message + contractErrStruct.Error.Message,
			ContractExecErrInner: &contractErrStruct,
		}

		return nil
	}

	return errors.New("failed to unmarshal ContractExecutionError")
}

func (contractEx *ContractExecutionError) MarshalJSON() ([]byte, error) {
	var temp any

	if contractEx.ContractExecErrInner != nil {
		temp = contractEx.ContractExecErrInner

		return json.Marshal(temp)
	}

	temp = contractEx.Message

	return json.Marshal(temp)
}

// can be either this struct or a string. The parent type will handle the unmarshalling
type ContractExecutionErrorInner struct {
	ContractAddress *felt.Felt              `json:"contract_address"`
	ClassHash       *felt.Felt              `json:"class_hash"`
	Selector        *felt.Felt              `json:"selector"`
	Error           *ContractExecutionError `json:"error"`
}

type TraceStatus string

const (
	TraceStatusReceived TraceStatus = "RECEIVED"
	TraceStatusRejected TraceStatus = "REJECTED"
)

func (s *TraceStatus) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	switch TraceStatus(str) {
	case TraceStatusReceived, TraceStatusRejected:
		*s = TraceStatus(str)

		return nil
	default:
		return fmt.Errorf("invalid trace status: %s", str)
	}
}
