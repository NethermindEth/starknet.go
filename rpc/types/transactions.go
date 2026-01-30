package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strconv"

	"github.com/NethermindEth/juno/core/felt"
)

// TransactionResponse is a generic response for all transaction types sent to the network.
type TransactionResponse struct {
	// Present for all transaction types
	Hash *felt.Felt `json:"transaction_hash"`
	// Present only for declare transactions
	ClassHash *felt.Felt `json:"class_hash,omitempty"`
	// Present only for deploy_account transactions
	ContractAddress *felt.Felt `json:"contract_address,omitempty"`
}

// DA_MODE: Specifies a storage domain in Starknet. Each domain has different
// guarantees regarding availability
type DataAvailabilityMode string

const (
	DAModeL1 DataAvailabilityMode = "L1"
	DAModeL2 DataAvailabilityMode = "L2"
)

// MarshalJSON implements the json.Marshaler interface.
// It validates that the DataAvailabilityMode is either L1 or L2 before marshalling.
func (da DataAvailabilityMode) MarshalJSON() ([]byte, error) {
	switch da {
	case DAModeL1, DAModeL2:
		return json.Marshal(string(da))
	default:
		return nil, fmt.Errorf(
			"invalid DataAvailabilityMode: %s, must be either L1 or L2",
			string(da),
		)
	}
}

func (da *DataAvailabilityMode) UInt64() (uint64, error) {
	switch *da {
	case DAModeL1:
		return uint64(0), nil
	case DAModeL2:
		return uint64(1), nil
	}

	return 0, errors.New("unknown DAMode")
}

type Resource string

// Values used in the Resource Bounds hash calculation
// Ref: https://docs.starknet.io/architecture-and-concepts/network-architecture/transactions/#v3_hash_calculation
//
//nolint:lll // The link would be unclickable if we break the line.
const (
	ResourceL1Gas     Resource = "L1_GAS"
	ResourceL2Gas     Resource = "L2_GAS"
	ResourceL1DataGas Resource = "L1_DATA"
)

// string must be NUM_AS_HEX
type TransactionVersion string

const (
	TransactionV0             TransactionVersion = "0x0"
	TransactionV0WithQueryBit TransactionVersion = "0x100000000000000000000000000000000"
	TransactionV1             TransactionVersion = "0x1"
	TransactionV1WithQueryBit TransactionVersion = "0x100000000000000000000000000000001"
	TransactionV2             TransactionVersion = "0x2"
	TransactionV2WithQueryBit TransactionVersion = "0x100000000000000000000000000000002"
	TransactionV3             TransactionVersion = "0x3"
	TransactionV3WithQueryBit TransactionVersion = "0x100000000000000000000000000000003"
)

// BigInt returns a big integer corresponding to the transaction version.
//
// Returns:
//   - *big.Int: a pointer to a big.Int
//   - error: an error if the conversion fails
func (v *TransactionVersion) BigInt() (*big.Int, error) {
	switch *v {
	case TransactionV0:
		return big.NewInt(0), nil
	case TransactionV1:
		return big.NewInt(1), nil
	case TransactionV2:
		return big.NewInt(2), nil
	case TransactionV3:
		return big.NewInt(3), nil
	}

	// Handle versions with query bit.
	// Remove the 0x prefix and convert to big.Int
	version, ok := new(big.Int).SetString(string(*v)[2:], 16) //nolint:mnd // hex base
	if !ok {
		return big.NewInt(-1), errors.New(fmt.Sprint("TransactionVersion %i not supported", *v))
	}

	return version, nil
}

// Int returns an integer corresponding to the transaction version.
// For versions with query bit, it returns the base version number (e.g.
// TransactionV2WithQueryBit returns 2).
// Returns -1 for invalid versions.
//
// Returns:
//   - int: the integer version, or -1 for invalid versions
func (v *TransactionVersion) Int() int {
	switch *v {
	case TransactionV0, TransactionV0WithQueryBit:
		return 0
	case TransactionV1, TransactionV1WithQueryBit:
		return 1
	case TransactionV2, TransactionV2WithQueryBit:
		return 2
	case TransactionV3, TransactionV3WithQueryBit:
		return 3
	}

	// Handle invalid versions
	return -1
}

type TransactionType string

const (
	TransactionTypeDeclare       TransactionType = "DECLARE"
	TransactionTypeDeployAccount TransactionType = "DEPLOY_ACCOUNT"
	TransactionTypeDeploy        TransactionType = "DEPLOY"
	TransactionTypeInvoke        TransactionType = "INVOKE"
	TransactionTypeL1Handler     TransactionType = "L1_HANDLER"
)

// UnmarshalJSON unmarshals the JSON data into a TransactionType.
func (tt *TransactionType) UnmarshalJSON(data []byte) error {
	unquoted, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}

	switch unquoted {
	case "DECLARE":
		*tt = TransactionTypeDeclare
	case "DEPLOY_ACCOUNT":
		*tt = TransactionTypeDeployAccount
	case "DEPLOY":
		*tt = TransactionTypeDeploy
	case "INVOKE":
		*tt = TransactionTypeInvoke
	case "L1_HANDLER":
		*tt = TransactionTypeL1Handler
	default:
		return fmt.Errorf("unsupported transaction type: %s", data)
	}

	return nil
}

// MarshalJSON marshals the TransactionType to JSON.
func (tt TransactionType) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(string(tt))), nil
}
