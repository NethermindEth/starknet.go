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
func Err(code int, data RPCData) *RPCError {
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
		return &RPCError{Code: InternalError, Message: err.Error(), Data: StringErrData(baseError.Error())}
	}

	var nodeErr RPCError
	err = json.Unmarshal(errBytes, &nodeErr)
	if err != nil {
		return &RPCError{Code: InternalError, Message: err.Error(), Data: StringErrData(baseError.Error())}
	}

	for _, rpcErr := range rpcErrors {
		if nodeErr.Code == rpcErr.Code && nodeErr.Message == rpcErr.Message {
			return &nodeErr
		}
	}

	if nodeErr.Code == 0 {
		return &RPCError{Code: InternalError, Message: "The error is not a valid RPC error", Data: StringErrData(baseError.Error())}
	}

	return Err(nodeErr.Code, nodeErr.Data)
}

type RPCError struct {
	Code    int     `json:"code"`
	Message string  `json:"message"`
	Data    RPCData `json:"data,omitempty"`
}

func (e RPCError) Error() string {
	if e.Data == nil || e.Data.ErrorMessage() == "" {
		return fmt.Sprintf("%d %s", e.Code, e.Message)
	}
	return fmt.Sprintf("%d %s: %s", e.Code, e.Message, e.Data.ErrorMessage())
}

// UnmarshalJSON implements the json.Unmarshaler interface for RPCError.
// It handles the deserialization of JSON into an RPCError struct,
// with special handling for the Data field.
func (e *RPCError) UnmarshalJSON(data []byte) error {
	// First try to unmarshal into a temporary struct without the RPCData interface
	var temp struct {
		Code    int             `json:"code"`
		Message string          `json:"message"`
		Data    json.RawMessage `json:"data,omitempty"`
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	e.Code = temp.Code
	e.Message = temp.Message

	// If there's no Data field, we're done
	if len(temp.Data) == 0 {
		e.Data = nil
		return nil
	}

	// Try to determine the concrete type of Data based on the RPCError code
	switch e.Code {
	case 10: // ErrNoTraceAvailable
		var data TraceStatusErrData
		if err := json.Unmarshal(temp.Data, &data); err != nil {
			return err
		}
		e.Data = &data
	case 40: // ErrContractError
		var data ContractErrData
		if err := json.Unmarshal(temp.Data, &data); err != nil {
			return err
		}
		e.Data = &data
	case 41: // ErrTxnExec
		var data TransactionExecErrData
		if err := json.Unmarshal(temp.Data, &data); err != nil {
			return err
		}
		e.Data = &data
	case 55, 56, 63: // ErrValidationFailure, ErrCompilationFailed, ErrUnexpectedError
		var strData string
		if err := json.Unmarshal(temp.Data, &strData); err != nil {
			return err
		}
		e.Data = StringErrData(strData)
	case 100: // ErrCompilationError
		var data CompilationErrData
		if err := json.Unmarshal(temp.Data, &data); err != nil {
			return err
		}
		e.Data = &data
	default:
		// For unknown error codes, try to unmarshal as string
		var strData string
		if err := json.Unmarshal(temp.Data, &strData); err == nil {
			e.Data = StringErrData(strData)
			return nil
		}

		// If not a string, set Data to nil and ignore the data field
		e.Data = nil
	}

	return nil
}

// RPCData is the interface that all error data types must implement
type RPCData interface {
	ErrorMessage() string
}

var _ RPCData = StringErrData("")
var _ RPCData = &CompilationErrData{}
var _ RPCData = &ContractErrData{}
var _ RPCData = &TransactionExecErrData{}
var _ RPCData = &TraceStatusErrData{}

var (
	ErrFailedToReceiveTxn = &RPCError{
		Code:    1,
		Message: "Failed to write transaction",
	}
	ErrNoTraceAvailable = &RPCError{
		Code:    10,
		Message: "No trace available for transaction",
		Data:    &TraceStatusErrData{},
	}
	ErrContractNotFound = &RPCError{
		Code:    20,
		Message: "Contract not found",
	}
	ErrEntrypointNotFound = &RPCError{
		Code:    21,
		Message: "Requested entrypoint does not exist in the contract",
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
		Data:    &ContractErrData{},
	}
	ErrTxnExec = &RPCError{
		Code:    41,
		Message: "Transaction execution error",
		Data:    &TransactionExecErrData{},
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
		Message: "Account balance is smaller than the transaction's maximal fee (calculated as the sum of each resource's limit x max price)",
	}
	ErrValidationFailure = &RPCError{
		Code:    55,
		Message: "Account validation failed",
		Data:    StringErrData(""),
	}
	ErrCompilationFailed = &RPCError{
		Code:    56,
		Message: "Compilation failed",
		Data:    StringErrData(""),
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
		Data:    StringErrData(""),
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
	ErrCompilationError = &RPCError{
		Code:    100,
		Message: "Failed to compile the contract",
		Data:    &CompilationErrData{},
	}
)

// StringErrData handles plain string data messages
type StringErrData string

func (s StringErrData) ErrorMessage() string {
	return string(s)
}

// Structured type for the ErrCompilationError data
type CompilationErrData struct {
	CompilationError string `json:"compilation_error,omitempty"`
}

func (c *CompilationErrData) ErrorMessage() string {
	return c.CompilationError
}

// Structured type for the ErrContractError data
type ContractErrData struct {
	RevertError ContractExecutionError `json:"revert_error,omitempty"`
}

func (c *ContractErrData) ErrorMessage() string {
	return c.RevertError.Message
}

// Structured type for the ErrTransactionExecError data
type TransactionExecErrData struct {
	TransactionIndex int                    `json:"transaction_index,omitempty"`
	ExecutionError   ContractExecutionError `json:"execution_error,omitempty"`
}

func (t *TransactionExecErrData) ErrorMessage() string {
	return t.ExecutionError.Message
}

// Structured type for the ErrTraceStatusError data
type TraceStatusErrData struct {
	Status TraceStatus `json:"status,omitempty"`
}

func (t *TraceStatusErrData) ErrorMessage() string {
	return string(t.Status)
}

// structured error that can later be processed by wallets or sdks
type ContractExecutionError struct {
	// the error raised during execution
	Message              string                       `json:",omitempty"`
	ContractExecErrInner *ContractExecutionErrorInner `json:",omitempty"`
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

	return fmt.Errorf("failed to unmarshal ContractExecutionError")
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
