package types

import (
	"fmt"
	"strconv"
)

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

// Flags that indicate how to simulate a given transaction. By default, the
// sequencer behaviour is replicated locally (enough funds are expected to be
// in the account, and fee will be deducted from the balance before the
// simulation of the next transaction). To skip the fee charge, use
// the SKIP_FEE_CHARGE flag.
type SimulationFlag string

const (
	SkipFeeCharge SimulationFlag = "SKIP_FEE_CHARGE"
	SkipValidate  SimulationFlag = "SKIP_VALIDATE"
)
