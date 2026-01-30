package types

import (
	"fmt"
	"strconv"
)

type TxnExecutionStatus string

const (
	TxnExecutionStatusSUCCEEDED TxnExecutionStatus = "SUCCEEDED"
	TxnExecutionStatusREVERTED  TxnExecutionStatus = "REVERTED"
)

// UnmarshalJSON unmarshals the JSON data into a TxnExecutionStatus struct.
func (ex *TxnExecutionStatus) UnmarshalJSON(data []byte) error {
	unquoted, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}
	switch unquoted {
	case "SUCCEEDED":
		*ex = TxnExecutionStatusSUCCEEDED
	case "REVERTED":
		*ex = TxnExecutionStatusREVERTED
	default:
		return fmt.Errorf("unsupported execution status: %s", data)
	}

	return nil
}

// MarshalJSON returns the JSON encoding of the TxnExecutionStatus.
func (ex TxnExecutionStatus) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(string(ex))), nil
}

// String returns the string representation of the TxnExecutionStatus.
func (ex TxnExecutionStatus) String() string {
	return string(ex)
}

type TxnFinalityStatus string

const (
	TxnFinalityStatusPreConfirmed TxnFinalityStatus = "PRE_CONFIRMED"
	TxnFinalityStatusAcceptedOnL2 TxnFinalityStatus = "ACCEPTED_ON_L2"
	TxnFinalityStatusAcceptedOnL1 TxnFinalityStatus = "ACCEPTED_ON_L1"
)

// UnmarshalJSON unmarshals the JSON data into a TxnFinalityStatus.
func (fs *TxnFinalityStatus) UnmarshalJSON(data []byte) error {
	unquoted, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}
	switch unquoted {
	case "PRE_CONFIRMED":
		*fs = TxnFinalityStatusPreConfirmed
	case "ACCEPTED_ON_L2":
		*fs = TxnFinalityStatusAcceptedOnL2
	case "ACCEPTED_ON_L1":
		*fs = TxnFinalityStatusAcceptedOnL1
	default:
		return fmt.Errorf("unsupported finality status: %s", data)
	}

	return nil
}

// MarshalJSON marshals the TxnFinalityStatus into JSON.
func (fs TxnFinalityStatus) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(string(fs))), nil
}

// String returns the string representation of the TxnFinalityStatus.
func (fs TxnFinalityStatus) String() string {
	return string(fs)
}

type TxnStatus string

const (
	TxnStatusReceived     TxnStatus = "RECEIVED"
	TxnStatusCandidate    TxnStatus = "CANDIDATE"
	TxnStatusPreConfirmed TxnStatus = "PRE_CONFIRMED"
	TxnStatusAcceptedOnL2 TxnStatus = "ACCEPTED_ON_L2"
	TxnStatusAcceptedOnL1 TxnStatus = "ACCEPTED_ON_L1"
)

// Transaction status result, including finality status and execution status
type TxnStatusResult struct {
	FinalityStatus  TxnStatus          `json:"finality_status"`
	ExecutionStatus TxnExecutionStatus `json:"execution_status,omitempty"`
	// the failure reason, only appears if execution_status is REVERTED
	FailureReason string `json:"failure_reason,omitempty"`
}
