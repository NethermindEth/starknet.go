package rpc

import (
	"encoding/json"
	"fmt"

	"github.com/NethermindEth/juno/core/felt"
)

const (
	InvalidJSON    = -32700 // Invalid JSON was received by the server.
	InvalidRequest = -32600 // The JSON sent is not a valid Request object.
	MethodNotFound = -32601 // The method does not exist / is not available.
	InvalidParams  = -32602 // Invalid method parameter(s).
	InternalError  = -32603 // Internal JSON-RPC error.
)

// Err returns an RPCError based on the given code and data.
//
// Parameters:
// - code: an integer representing the error code.
// - data: any data associated with the error.
// Returns
// - *RPCError: a pointer to an RPCError object.
func Err(code int, data *RPCData) *RPCError {
	switch code {
	case InvalidJSON:
		return &RPCError{Code: InvalidJSON, Message: "Parse error", Data: data}
	case InvalidRequest:
		return &RPCError{Code: InvalidRequest, Message: "Invalid Request", Data: data}
	case MethodNotFound:
		return &RPCError{Code: MethodNotFound, Message: "Method Not Found", Data: data}
	case InvalidParams:
		return &RPCError{Code: InvalidParams, Message: "Invalid Params", Data: data}
	default:
		return &RPCError{Code: InternalError, Message: "Internal Error", Data: data}
	}
}

// tryUnwrapToRPCErr unwraps the error and checks if it matches any of the given RPC errors.
// If a match is found, the corresponding RPC error is returned.
// If no match is found, the function returns an InternalError with the original error.
//
// Parameters:
// - err: The error to be unwrapped
// - rpcErrors: variadic list of *RPCError objects to be checked
// Returns:
// - error: the original error
func tryUnwrapToRPCErr(baseError error, rpcErrors ...*RPCError) *RPCError {
	errBytes, err := json.Marshal(baseError)
	if err != nil {
		return &RPCError{Code: InternalError, Message: err.Error()}
	}

	var nodeErr RPCError
	err = json.Unmarshal(errBytes, &nodeErr)
	if err != nil {
		return &RPCError{Code: InternalError, Message: err.Error()}
	}

	for _, rpcErr := range rpcErrors {
		if nodeErr.Code == rpcErr.Code && nodeErr.Message == rpcErr.Message {
			return &nodeErr
		}
	}

	if nodeErr.Code == 0 {
		return &RPCError{Code: InternalError, Message: "The error is not a valid RPC error", Data: &RPCData{Message: baseError.Error()}}
	}

	return Err(nodeErr.Code, nodeErr.Data)
}

type RPCError struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Data    *RPCData `json:"data,omitempty"`
}

func (e RPCError) Error() string {
	if e.Data == nil || e.Data.Message == "" {
		return e.Message
	}
	return e.Message + ": " + e.Data.Message
}

type RPCData struct {
	Message                       string                         `json:",omitempty"`
	CompilationErrorData          *CompilationErrorData          `json:",omitempty"`
	ContractErrorData             *ContractErrorData             `json:",omitempty"`
	TransactionExecutionErrorData *TransactionExecutionErrorData `json:",omitempty"`
}

func (rpcData *RPCData) UnmarshalJSON(data []byte) error {
	var message string
	if err := json.Unmarshal(data, &message); err == nil {
		rpcData.Message = message
		return nil
	}

	var compilationErrData CompilationErrorData
	if err := json.Unmarshal(data, &compilationErrData); err == nil {
		*rpcData = RPCData{
			Message:              rpcData.Message + compilationErrData.CompilationError,
			CompilationErrorData: &compilationErrData,
		}
		return nil
	}

	var contractErrData ContractErrorData
	if err := json.Unmarshal(data, &contractErrData); err == nil {
		*rpcData = RPCData{
			Message:           rpcData.Message + contractErrData.RevertError.Message,
			ContractErrorData: &contractErrData,
		}
		return nil
	}

	var txExErrData TransactionExecutionErrorData
	if err := json.Unmarshal(data, &txExErrData); err == nil {
		*rpcData = RPCData{
			Message:                       rpcData.Message + txExErrData.ExecutionError.Message,
			TransactionExecutionErrorData: &txExErrData,
		}
		return nil
	}

	return fmt.Errorf("failed to unmarshal RPCData")
}

func (rpcData *RPCData) MarshalJSON() ([]byte, error) {
	var temp any

	if rpcData.CompilationErrorData != nil {
		temp = *rpcData.CompilationErrorData
		return json.Marshal(temp)
	}

	if rpcData.ContractErrorData != nil {
		temp = *rpcData.ContractErrorData
		return json.Marshal(temp)
	}

	if rpcData.TransactionExecutionErrorData != nil {
		temp = *rpcData.TransactionExecutionErrorData
		return json.Marshal(temp)
	}

	temp = rpcData.Message

	return json.Marshal(temp)
}

var (
	ErrFailedToReceiveTxn = &RPCError{
		Code:    1,
		Message: "Failed to write transaction",
	}
	ErrNoTraceAvailable = &RPCError{
		Code:    10,
		Message: "No trace available for transaction",
	}
	ErrContractNotFound = &RPCError{
		Code:    20,
		Message: "Contract not found",
	}
	ErrBlockNotFound = &RPCError{
		Code:    24,
		Message: "Block not found",
	}
	ErrInvalidTxnHash = &RPCError{
		Code:    25,
		Message: "Invalid transaction hash",
	}
	ErrInvalidBlockHash = &RPCError{
		Code:    26,
		Message: "Invalid block hash",
	}
	ErrInvalidTxnIndex = &RPCError{
		Code:    27,
		Message: "Invalid transaction index in a block",
	}
	ErrClassHashNotFound = &RPCError{
		Code:    28,
		Message: "Class hash not found",
	}
	ErrHashNotFound = &RPCError{
		Code:    29,
		Message: "Transaction hash not found",
	}
	ErrPageSizeTooBig = &RPCError{
		Code:    31,
		Message: "Requested page size is too big",
	}
	ErrNoBlocks = &RPCError{
		Code:    32,
		Message: "There are no blocks",
	}
	ErrInvalidContinuationToken = &RPCError{
		Code:    33,
		Message: "The supplied continuation token is invalid or unknown",
	}
	ErrTooManyKeysInFilter = &RPCError{
		Code:    34,
		Message: "Too many keys provided in a filter",
	}
	ErrContractError = &RPCError{
		Code:    40,
		Message: "Contract error",
	}
	ErrTxnExec = &RPCError{
		Code:    41,
		Message: "Transaction execution error",
	}
	ErrStorageProofNotSupported = &RPCError{
		Code:    42,
		Message: "the node doesn't support storage proofs for blocks that are too far in the past",
	}
	ErrInvalidContractClass = &RPCError{
		Code:    50,
		Message: "Invalid contract class",
	}
	ErrClassAlreadyDeclared = &RPCError{
		Code:    51,
		Message: "Class already declared",
	}
	ErrInvalidTransactionNonce = &RPCError{
		Code:    52,
		Message: "Invalid transaction nonce",
	}
	ErrInsufficientResourcesForValidate = &RPCError{
		Code:    53,
		Message: "The transaction's resources don't cover validation or the minimal transaction fee",
	}
	ErrInsufficientAccountBalance = &RPCError{
		Code:    54,
		Message: "Account balance is smaller than the transaction's max_fee",
	}
	ErrValidationFailure = &RPCError{
		Code:    55,
		Message: "Account validation failed",
	}
	ErrCompilationFailed = &RPCError{
		Code:    56,
		Message: "Compilation failed",
	}
	ErrContractClassSizeTooLarge = &RPCError{
		Code:    57,
		Message: "Contract class size is too large",
	}
	ErrNonAccount = &RPCError{
		Code:    58,
		Message: "Sender address is not an account contract",
	}
	ErrDuplicateTx = &RPCError{
		Code:    59,
		Message: "A transaction with the same hash already exists in the mempool",
	}
	ErrCompiledClassHashMismatch = &RPCError{
		Code:    60,
		Message: "The compiled class hash did not match the one supplied in the transaction",
	}
	ErrUnsupportedTxVersion = &RPCError{
		Code:    61,
		Message: "The transaction version is not supported",
	}
	ErrUnsupportedContractClassVersion = &RPCError{
		Code:    62,
		Message: "The contract class version is not supported",
	}
	ErrUnexpectedError = &RPCError{
		Code:    63,
		Message: "An unexpected error occurred",
	}
	ErrInvalidSubscriptionID = &RPCError{
		Code:    66,
		Message: "Invalid subscription id",
	}
	ErrTooManyAddressesInFilter = &RPCError{
		Code:    67,
		Message: "Too many addresses in filter sender_address filter",
	}
	ErrTooManyBlocksBack = &RPCError{
		Code:    68,
		Message: "Cannot go back more than 1024 blocks",
	}
	ErrCallOnPending = &RPCError{
		Code:    69,
		Message: "This method does not support being called on the pending block",
	}
	ErrCompilationError = &RPCError{
		Code:    100,
		Message: "Failed to compile the contract",
	}
)

type CompilationErrorData struct {
	// More data about the compilation failure
	CompilationError string `json:"compilation_error,omitempty"`
}

type ContractErrorData struct {
	// the execution trace up to the point of failure
	RevertError ContractExecutionError `json:"revert_error,omitempty"`
}

type TransactionExecutionErrorData struct {
	// The index of the first transaction failing in a sequence of given transactions
	TransactionIndex int `json:"transaction_index,omitempty"`
	// the execution trace up to the point of failure
	ExecutionError ContractExecutionError `json:"execution_error,omitempty"`
}

type ContractExecutionError struct {
	// the error raised during execution
	Message                      string `json:",omitempty"`
	*ContractExecutionErrorInner `json:",omitempty"`
}

func (contractEx *ContractExecutionError) UnmarshalJSON(data []byte) error {
	var contractErrStruct ContractExecutionErrorInner
	var message string

	if err := json.Unmarshal(data, &message); err == nil {
		*contractEx = ContractExecutionError{
			Message:                     message,
			ContractExecutionErrorInner: &contractErrStruct,
		}
		return nil
	}

	if err := json.Unmarshal(data, &contractErrStruct); err == nil {
		*contractEx = ContractExecutionError{
			Message:                     "",
			ContractExecutionErrorInner: &contractErrStruct,
		}
		return nil
	}

	return fmt.Errorf("failed to unmarshal ContractExecutionError")
}

func (contractEx *ContractExecutionError) MarshalJSON() ([]byte, error) {
	var temp any

	if contractEx.ContractExecutionErrorInner != nil {
		temp = contractEx.ContractExecutionErrorInner
		return json.Marshal(temp)
	}

	temp = contractEx.Message

	return json.Marshal(temp)
}

type ContractExecutionErrorInner struct {
	ContractAddress *felt.Felt              `json:"contract_address"`
	ClassHash       *felt.Felt              `json:"class_hash"`
	Selector        *felt.Felt              `json:"selector"`
	Error           *ContractExecutionError `json:"error"`
}
