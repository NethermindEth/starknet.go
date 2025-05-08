package rpc

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/NethermindEth/juno/core/felt"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
)

type SimulateTransactionInput struct {
	// a sequence of transactions to simulate, running each transaction on the state resulting from applying all the previous ones
	Txns            []BroadcastTxn   `json:"transactions"`
	BlockID         BlockID          `json:"block_id"`
	SimulationFlags []SimulationFlag `json:"simulation_flags"`
}

type SimulationFlag string

const (
	SKIP_FEE_CHARGE SimulationFlag = "SKIP_FEE_CHARGE"
	SKIP_EXECUTE    SimulationFlag = "SKIP_EXECUTE"
	// Flags that indicate how to simulate a given transaction. By default, the sequencer behaviour is replicated locally
	SKIP_VALIDATE SimulationFlag = "SKIP_VALIDATE"
)

type SimulatedTransaction struct {
	TxnTrace      `json:"transaction_trace"`
	FeeEstimation `json:"fee_estimation"`
}

type TxnTrace interface{}

var (
	_ TxnTrace = InvokeTxnTrace{}
	_ TxnTrace = DeclareTxnTrace{}
	_ TxnTrace = DeployAccountTxnTrace{}
	_ TxnTrace = L1HandlerTxnTrace{}
)

// the execution trace of an invoke transaction
type InvokeTxnTrace struct {
	ValidateInvocation *FnInvocation `json:"validate_invocation,omitempty"`
	// the trace of the __execute__ call or constructor call, depending on the transaction type (none for declare transactions)
	ExecuteInvocation     ExecInvocation     `json:"execute_invocation"`
	FeeTransferInvocation *FnInvocation      `json:"fee_transfer_invocation,omitempty"`
	StateDiff             *StateDiff         `json:"state_diff,omitempty"`
	Type                  TransactionType    `json:"type"`
	ExecutionResources    ExecutionResources `json:"execution_resources"`
}

// the execution trace of a declare transaction
type DeclareTxnTrace struct {
	ValidateInvocation    *FnInvocation      `json:"validate_invocation,omitempty"`
	FeeTransferInvocation *FnInvocation      `json:"fee_transfer_invocation,omitempty"`
	StateDiff             *StateDiff         `json:"state_diff,omitempty"`
	Type                  TransactionType    `json:"type"`
	ExecutionResources    ExecutionResources `json:"execution_resources"`
}

// the execution trace of a deploy account transaction
type DeployAccountTxnTrace struct {
	ValidateInvocation *FnInvocation `json:"validate_invocation,omitempty"`
	// the trace of the __execute__ call or constructor call, depending on the transaction type (none for declare transactions)
	ConstructorInvocation FnInvocation       `json:"constructor_invocation"`
	FeeTransferInvocation *FnInvocation      `json:"fee_transfer_invocation,omitempty"`
	StateDiff             *StateDiff         `json:"state_diff,omitempty"`
	Type                  TransactionType    `json:"type"`
	ExecutionResources    ExecutionResources `json:"execution_resources"`
}

// the execution trace of an L1 handler transaction
type L1HandlerTxnTrace struct {
	// the trace of the __execute__ call or constructor call, depending on the transaction type (none for declare transactions)
	FunctionInvocation FnInvocation       `json:"function_invocation"`
	StateDiff          *StateDiff         `json:"state_diff,omitempty"`
	ExecutionResources ExecutionResources `json:"execution_resources"`
	Type               TransactionType    `json:"type"`
}

type EntryPointType string

const (
	External    EntryPointType = "EXTERNAL"
	L1Handler   EntryPointType = "L1_HANDLER"
	Constructor EntryPointType = "CONSTRUCTOR"
)

type CallType string

const (
	CallTypeLibraryCall CallType = "LIBRARY_CALL"
	CallTypeCall        CallType = "CALL"
	CallTypeDelegate    CallType = "DELEGATE"
)

type FnInvocation struct {
	FunctionCall

	// The address of the invoking contract. 0 for the root invocation
	CallerAddress *felt.Felt `json:"caller_address"`

	// The hash of the class being called
	ClassHash *felt.Felt `json:"class_hash"`

	EntryPointType EntryPointType `json:"entry_point_type"`

	CallType CallType `json:"call_type"`

	// The value returned from the function invocation
	Result []*felt.Felt `json:"result"`

	// The calls made by this invocation
	NestedCalls []FnInvocation `json:"calls"`

	// The events emitted in this invocation
	InvocationEvents []OrderedEvent `json:"events"`

	// The messages sent by this invocation to L1
	L1Messages []OrderedMsg `json:"messages"`

	// Resources consumed by the internal call
	ExecutionResources InnerCallExecutionResources `json:"execution_resources"`

	// True if this inner call panicked
	IsReverted bool `json:"is_reverted"`
}

// the resources consumed by an inner call (does not account for state diffs since data is squashed across the transaction)
type InnerCallExecutionResources struct {
	// l1 gas consumed by this transaction, used for l2-->l1 messages and state updates if blobs are not used
	L1Gas uint `json:"l1_gas"`
	// l2 gas consumed by this transaction, used for computation and calldata
	L2Gas uint `json:"l2_gas"`
}

// A single pair of transaction hash and corresponding trace
type Trace struct {
	TraceRoot TxnTrace   `json:"trace_root,omitempty"`
	TxnHash   *felt.Felt `json:"transaction_hash,omitempty"`
}

type ExecInvocation struct {
	*FnInvocation
	RevertReason string `json:"revert_reason,omitempty"`
}

// UnmarshalJSON unmarshals the data into a SimulatedTransaction object.
//
// It takes a byte slice as the parameter, representing the JSON data to be unmarshalled.
// The function returns an error if the unmarshalling process fails.
//
// Parameters:
//   - data: The JSON data to be unmarshalled
//
// Returns:
//   - error: An error if the unmarshalling process fails
func (txn *SimulatedTransaction) UnmarshalJSON(data []byte) error {
	var dec map[string]interface{}
	if err := json.Unmarshal(data, &dec); err != nil {
		return err
	}

	// SimulatedTransaction wraps transactions in the TxnTrace field.
	rawTxnTrace, err := internalUtils.UnwrapJSON(dec, "transaction_trace")
	if err != nil {
		return err
	}

	trace, err := unmarshalTraceTxn(rawTxnTrace)
	if err != nil {
		return err
	}

	var feeEstimation FeeEstimation

	if feeEstimationData, ok := dec["fee_estimation"]; ok {
		err = remarshal(feeEstimationData, &feeEstimation)
		if err != nil {
			return err
		}
	} else {
		return errors.New("fee estimate not found")
	}

	*txn = SimulatedTransaction{
		TxnTrace:      trace,
		FeeEstimation: feeEstimation,
	}

	return nil
}

// UnmarshalJSON unmarshals the data into a Trace object.
//
// It takes a byte slice as the parameter, representing the JSON data to be unmarshalled.
// The function returns an error if the unmarshalling process fails.
//
// Parameters:
//   - data: The JSON data to be unmarshalled
//
// Returns:
//   - error: An error if the unmarshalling process fails
func (txn *Trace) UnmarshalJSON(data []byte) error {
	var dec map[string]interface{}
	if err := json.Unmarshal(data, &dec); err != nil {
		return err
	}

	// Trace wrap trace transactions in the TraceRoot field.
	rawTraceTx, err := internalUtils.UnwrapJSON(dec, "trace_root")
	if err != nil {
		return err
	}

	t, err := unmarshalTraceTxn(rawTraceTx)
	if err != nil {
		return err
	}

	var txHash *felt.Felt
	if txHashData, ok := dec["transaction_hash"]; ok {
		txHashString, ok := txHashData.(string)
		if !ok {
			return errors.New("failed to unmarshal transaction hash, transaction_hash is not a string")
		}
		txHash, err = internalUtils.HexToFelt(txHashString)
		if err != nil {
			return err
		}
	} else {
		return errors.New("failed to unmarshal transaction hash, transaction_hash not found")
	}

	*txn = Trace{
		TraceRoot: t,
		TxnHash:   txHash,
	}

	return nil
}

// unmarshalTraceTxn unmarshals a given interface and returns a TxnTrace.
//
// Parameter:
//   - t: The interface{} to be unmarshalled
//
// Returns:
//   - TxnTrace: a TxnTrace
//   - error: an error if the unmarshaling process fails
func unmarshalTraceTxn(t interface{}) (TxnTrace, error) {
	if casted, ok := t.(map[string]interface{}); ok {
		switch TransactionType(casted["type"].(string)) {
		case TransactionType_Declare:
			var txn DeclareTxnTrace
			err := remarshal(casted, &txn)

			return txn, err
		case TransactionType_DeployAccount:
			var txn DeployAccountTxnTrace
			err := remarshal(casted, &txn)

			return txn, err
		case TransactionType_Invoke:
			var txn InvokeTxnTrace
			err := remarshal(casted, &txn)

			return txn, err
		case TransactionType_L1Handler:
			var txn L1HandlerTxnTrace
			err := remarshal(casted, &txn)

			return txn, err
		}
	}

	return nil, fmt.Errorf("unknown transaction type: %v", t)
}
