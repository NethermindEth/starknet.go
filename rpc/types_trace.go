package rpc

import "github.com/NethermindEth/juno/core/felt"

type SimulateTransactionInput struct {
	//a sequence of transactions to simulate, running each transaction on the state resulting from applying all the previous ones
	Txns            []Transaction    `json:"transactions"`
	BlockID         BlockID          `json:"block_id"`
	SimulationFlags []SimulationFlag `json:"simulation_flags"`
}

type SimulationFlag string

const (
	SKIP_FEE_CHARGE SimulationFlag = "SKIP_FEE_CHARGE"
	SKIP_EXECUTE    SimulationFlag = "SKIP_EXECUTE"
)

// The execution trace and consumed resources of the required transactions
type SimulateTransactionOutput struct {
	Txns []SimulatedTransaction `json:"result"`
}

type SimulatedTransaction struct {
	TxnTrace `json:"transaction_trace"`
	FeeEstimate
}

type TxnTrace interface{}

var _ TxnTrace = InvokeTxnTrace{}
var _ TxnTrace = DeclareTxnTrace{}
var _ TxnTrace = DeployAccountTxnTrace{}
var _ TxnTrace = L1HandlerTxnTrace{}

// the execution trace of an invoke transaction
type InvokeTxnTrace struct {
	ValidateInvocation FnInvocation `json:"validate_invocation"`
	//the trace of the __execute__ call or constructor call, depending on the transaction type (none for declare transactions)
	ExecuteInvocation     ExecInvocation  `json:"execute_invocation"`
	FeeTransferInvocation FnInvocation    `json:"fee_transfer_invocation"`
	StateDiff             StateDiff       `json:"state_diff"`
	Type                  TransactionType `json:"type"`
}

// the execution trace of a declare transaction
type DeclareTxnTrace struct {
	ValidateInvocation    FnInvocation    `json:"validate_invocation"`
	FeeTransferInvocation FnInvocation    `json:"fee_transfer_invocation"`
	StateDiff             StateDiff       `json:"state_diff"`
	Type                  TransactionType `json:"type"`
}

// the execution trace of a deploy account transaction
type DeployAccountTxnTrace struct {
	ValidateInvocation FnInvocation `json:"validate_invocation"`
	//the trace of the __execute__ call or constructor call, depending on the transaction type (none for declare transactions)
	ConstructorInvocation FnInvocation    `json:"constructor_invocation"`
	FeeTransferInvocation FnInvocation    `json:"fee_transfer_invocation"`
	StateDiff             StateDiff       `json:"state_diff"`
	Type                  TransactionType `json:"type"`
}

// the execution trace of an L1 handler transaction
type L1HandlerTxnTrace struct {
	//the trace of the __execute__ call or constructor call, depending on the transaction type (none for declare transactions)
	FunctionInvocation FnInvocation    `json:"function_invocation"`
	Type               TransactionType `json:"type"`
}

type EntryPointType string

const (
	External    EntryPointType = "EXTERNAL"
	L1Handler   EntryPointType = "L1_HANDLER"
	Constructor EntryPointType = "CONSTRUCTOR"
)

type CallType string

const (
	LibraryCall CallType = "LIBRARY_CALL"
	Call        CallType = "CALL"
)

type FnInvocation struct {
	FunctionCall

	//The address of the invoking contract. 0 for the root invocation
	CallerAddress *felt.Felt `json:"caller_address"`

	// The hash of the class being called
	ClassHash *felt.Felt `json:"class_hash"`

	EntryPointType EntryPointType `json:"entry_point_type"`

	CallType CallType `json:"call_type"`

	//The value returned from the function invocation
	Result []*felt.Felt `json:"result"`

	// The calls made by this invocation
	NestedCalls []FnInvocation `json:"calls"`

	// The events emitted in this invocation
	InvocationEvents []OrderedEvent `json:"events"`

	// The messages sent by this invocation to L1
	L1Messages []OrderedMsg `json:"messages"`
}

// A single pair of transaction hash and corresponding trace
type Trace struct {
	TraceRoot TxnTrace   `json:"trace_root,omitempty"`
	TxnHash   *felt.Felt `json:"transaction_hash,omitempty"`
}

type ExecInvocation struct {
	FunctionInvocation FnInvocation `json:"function_invocation,omitempty"`
	RevertReason       string       `json:"revert_reason,omitempty"`
}
