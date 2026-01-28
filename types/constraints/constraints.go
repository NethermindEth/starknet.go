package constraints

import (
	"math/big"

	"github.com/NethermindEth/juno/core/felt"
)

type ResourceBounds[u64 U64, u128 U128] interface {
	~struct {
		MaxAmount       u64  `json:"max_amount"`
		MaxPricePerUnit u128 `json:"max_price_per_unit"`
	}
}

type ResourceBoundsImpl[u64 U64, u128 U128] struct {
	MaxAmount       u64  `json:"max_amount"`
	MaxPricePerUnit u128 `json:"max_price_per_unit"`
}

type ResourceBoundsMapping[u64 U64, u128 U128, B ResourceBounds[u64, u128]] interface {
	~struct {
		L1Gas     B `json:"l1_gas"`
		L1DataGas B `json:"l1_data_gas"`
		L2Gas     B `json:"l2_gas"`
	}
}

type ResourceBoundsMappingImpl[u64 U64, u128 U128, B ResourceBounds[u64, u128]] struct {
	L1Gas     B `json:"l1_gas"`
	L1DataGas B `json:"l1_data_gas"`
	L2Gas     B `json:"l2_gas"`
}

type DataAvailabilityMode interface {
	~string
	UInt64() (uint64, error)
}
type U64 interface {
	~string
	ToUint64() (uint64, error)
}
type U128 interface {
	~string
	ToBigInt() (*big.Int, error)
}

type FeeEstimation[PriceUnit ~string] interface {
	~struct {
		L1GasConsumed     *felt.Felt `json:"l1_gas_consumed"`
		L1GasPrice        *felt.Felt `json:"l1_gas_price"`
		L2GasConsumed     *felt.Felt `json:"l2_gas_consumed"`
		L2GasPrice        *felt.Felt `json:"l2_gas_price"`
		L1DataGasConsumed *felt.Felt `json:"l1_data_gas_consumed"`
		L1DataGasPrice    *felt.Felt `json:"l1_data_gas_price"`
		OverallFee        *felt.Felt `json:"overall_fee"`
		Unit              PriceUnit  `json:"unit"`
	}
}

type FeeEstimationImpl[PriceUnit ~string] struct {
	L1GasConsumed     *felt.Felt `json:"l1_gas_consumed"`
	L1GasPrice        *felt.Felt `json:"l1_gas_price"`
	L2GasConsumed     *felt.Felt `json:"l2_gas_consumed"`
	L2GasPrice        *felt.Felt `json:"l2_gas_price"`
	L1DataGasConsumed *felt.Felt `json:"l1_data_gas_consumed"`
	L1DataGasPrice    *felt.Felt `json:"l1_data_gas_price"`
	OverallFee        *felt.Felt `json:"overall_fee"`
	Unit              PriceUnit  `json:"unit"`
}
