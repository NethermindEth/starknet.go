package types

import (
	"fmt"
	"strconv"

	"github.com/NethermindEth/juno/core/felt"
)

// Fee estimation common fields
type FeeEstimationCommon struct {
	// The Ethereum gas consumption of the transaction, charged for L1->L2
	// messages and, depending on the block's DA_MODE, state diffs
	L1GasConsumed *felt.Felt `json:"l1_gas_consumed"`

	// The gas price (in wei or fri, depending on the tx version) that was
	// used in the cost estimation.
	L1GasPrice *felt.Felt `json:"l1_gas_price"`

	// The L2 gas consumption of the transaction
	L2GasConsumed *felt.Felt `json:"l2_gas_consumed"`

	// The L2 gas price (in wei or fri, depending on the tx version) that
	// was used in the cost estimation.
	L2GasPrice *felt.Felt `json:"l2_gas_price"`

	// The Ethereum data gas consumption of the transaction.
	L1DataGasConsumed *felt.Felt `json:"l1_data_gas_consumed"`

	// The data gas price (in wei or fri, depending on the tx version) that
	// was used in the cost estimation.
	L1DataGasPrice *felt.Felt `json:"l1_data_gas_price"`

	// The estimated fee for the transaction (in wei or fri, depending on the
	// tx version), equals to gas_consumed*gas_price + data_gas_consumed*data_gas_price.
	OverallFee *felt.Felt `json:"overall_fee"`
}

type FeeEstimation struct {
	FeeEstimationCommon
	// Units in which the fee is given, can only be FRI
	Unit PriceUnitFri `json:"unit"`
}

type MessageFeeEstimation struct {
	FeeEstimationCommon
	// Units in which the fee is given, can only be WEI
	Unit PriceUnitWei `json:"unit"`
}

type FeePayment struct {
	Amount *felt.Felt `json:"amount"`
	Unit   PriceUnit  `json:"unit"`
}

// Units in which the fee is given
type PriceUnit string

const (
	UnitWei PriceUnit = "WEI"
	UnitFri PriceUnit = "FRI"
)

// Representation of the unit WEI
type PriceUnitWei string

const (
	WeiUnit PriceUnitWei = "WEI"
)

// Representation of the unit FRI
type PriceUnitFri string

const (
	FriUnit PriceUnitFri = "FRI"
)

// Unmarshals the JSON data into a PriceUnit.
func (f *PriceUnit) UnmarshalJSON(data []byte) error {
	unquoted, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}

	switch unquoted {
	case "WEI":
		*f = UnitWei
	case "FRI":
		*f = UnitFri
	default:
		return fmt.Errorf("unsupported price unit: %s", data)
	}

	return nil
}

// Unmarshals the JSON data into a PriceUnitWei.
func (f *PriceUnitWei) UnmarshalJSON(data []byte) error {
	unquoted, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}

	if unquoted != string(WeiUnit) {
		return fmt.Errorf("price unit should be WEI, got: %s", data)
	}

	*f = WeiUnit

	return nil
}

// Unmarshals the JSON data into a PriceUnitFri.
func (f *PriceUnitFri) UnmarshalJSON(data []byte) error {
	unquoted, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}

	if unquoted != string(FriUnit) {
		return fmt.Errorf("price unit should be FRI, got: %s", data)
	}

	*f = FriUnit

	return nil
}
