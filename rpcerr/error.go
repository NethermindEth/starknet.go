package rpcerr

import (
	"encoding/json"
	"fmt"
	"strings"
)

const (
	InvalidJSON    = -32700 // Invalid JSON was received by the server.
	InvalidRequest = -32600 // The JSON sent is not a valid Request object.
	MethodNotFound = -32601 // The method does not exist / is not available.
	InvalidParams  = -32602 // Invalid method parameter(s).
	InternalError  = -32603 // Internal JSON-RPC error.
)

// RPCError represents an error response from a JSON-RPC server.
// It contains a code, message, and optional data.
type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	// Data is optional and can be any type that implements the RPCData interface.
	// It will be nil if there is no data.
	Data RPCData `json:"data,omitempty"`
}

// UnmarshalJSON implements the json.Unmarshaler interface for RPCError.
// It handles the deserialization of JSON into an RPCError struct,
// and if there is data, it stores it as a string in the Data field.
func (e *RPCError) UnmarshalJSON(data []byte) error {
	// First try to unmarshal into a temporary struct without the RPCData interface
	var temp struct {
		Code    int             `json:"code"`
		Message string          `json:"message"`
		Data    json.RawMessage `json:"data"`
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

	// If there is data, it stores it as a string in the Data field.
	rawData, err := temp.Data.MarshalJSON()
	if err != nil {
		return err
	}
	e.Data = StringErrData(string(rawData))

	return nil
}

// Error returns the error message
func (e RPCError) Error() string {
	if e.Data == nil || e.Data.ErrorMessage() == "" {
		return fmt.Sprintf("%d %s", e.Code, e.Message)
	}

	return fmt.Sprintf("%d %s: %s", e.Code, e.Message, e.Data.ErrorMessage())
}

// RPCData is the interface that all error data types must implement
type RPCData interface {
	ErrorMessage() string
}

// StringErrData handles plain string data messages
type StringErrData string

// ErrorMessage returns the error message of the StringErrData as a string
func (s StringErrData) ErrorMessage() string {
	return string(s)
}

//nolint:exhaustruct
var (
	_ RPCData = StringErrData("")
)

// Err returns an RPCError based on the given code and data.
//
// Parameters:
//   - code: an integer representing the error code.
//   - data: any data associated with the error.
//
// Returns
//   - *RPCError: a pointer to an RPCError object.
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
		data = StringErrData(fmt.Sprintf("%d %s", code, data))

		return &RPCError{Code: InternalError, Message: "Internal Error", Data: data}
	}
}

// UnwrapToRPCErr unwraps the error and checks if it matches any of the given RPC errors.
// If a match is found, the corresponding RPC error is returned.
// If no match is found, the function returns an InternalError with the original error.
//
// Parameters:
//   - err: The error to be unwrapped
//   - rpcErrors: variadic list of *RPCError objects to be checked
//
// Returns:
//   - error: the original error
func UnwrapToRPCErr(baseError error, rpcErrors ...*RPCError) *RPCError {
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
		if nodeErr.Code == rpcErr.Code && strings.EqualFold(nodeErr.Message, rpcErr.Message) {
			return &nodeErr
		}
	}

	if nodeErr.Code == 0 {
		return &RPCError{Code: InternalError, Message: "The error is not a valid RPC error", Data: StringErrData(baseError.Error())}
	}

	// return many data as possible
	if nodeErr.Data != nil {
		return Err(nodeErr.Code, StringErrData(fmt.Sprintf("%s %s", nodeErr.Message, nodeErr.Data.ErrorMessage())))
	}

	return Err(nodeErr.Code, StringErrData(nodeErr.Message))
}
