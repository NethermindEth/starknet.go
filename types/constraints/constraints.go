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

type DeclareTxnV3[
	TransactionType ~string,
	TransactionVersion ~string,
	u64 U64,
	u128 U128,
	RB ResourceBounds[u64, u128],
	RBM ResourceBoundsMapping[u64, u128, RB],
	DA DataAvailabilityMode,
] interface {
	~struct {
		Type                  TransactionType    `json:"type"`
		SenderAddress         *felt.Felt         `json:"sender_address"`
		CompiledClassHash     *felt.Felt         `json:"compiled_class_hash"`
		Version               TransactionVersion `json:"version"`
		Signature             []*felt.Felt       `json:"signature"`
		Nonce                 *felt.Felt         `json:"nonce"`
		ClassHash             *felt.Felt         `json:"class_hash"`
		ResourceBounds        *RBM               `json:"resource_bounds"`
		Tip                   u64                `json:"tip"`
		PayMasterData         []*felt.Felt       `json:"paymaster_data"`
		AccountDeploymentData []*felt.Felt       `json:"account_deployment_data"`
		NonceDataMode         DA                 `json:"nonce_data_availability_mode"`
		FeeMode               DA                 `json:"fee_data_availability_mode"`
	}
}

type DeclareTxnV3Impl[
	TransactionType ~string,
	TransactionVersion ~string,
	u64 U64,
	u128 U128,
	RB ResourceBounds[u64, u128],
	RBM ResourceBoundsMapping[u64, u128, RB],
	DA DataAvailabilityMode,
] struct {
	Type                  TransactionType    `json:"type"`
	SenderAddress         *felt.Felt         `json:"sender_address"`
	CompiledClassHash     *felt.Felt         `json:"compiled_class_hash"`
	Version               TransactionVersion `json:"version"`
	Signature             []*felt.Felt       `json:"signature"`
	Nonce                 *felt.Felt         `json:"nonce"`
	ClassHash             *felt.Felt         `json:"class_hash"`
	ResourceBounds        *RBM               `json:"resource_bounds"`
	Tip                   u64                `json:"tip"`
	PayMasterData         []*felt.Felt       `json:"paymaster_data"`
	AccountDeploymentData []*felt.Felt       `json:"account_deployment_data"`
	NonceDataMode         DA                 `json:"nonce_data_availability_mode"`
	FeeMode               DA                 `json:"fee_data_availability_mode"`
}

// THIS WORKSSS****************************
// func foo[
// 	TransactionType ~string,
// 	TransactionVersion ~string,
// 	u64 internal.U64,
// 	u128 internal.U128,
// 	RB internal.ResourceBounds[u64, u128],
// 	RBM internal.ResourceBoundsMapping[u64, u128, RB],
// 	DA internal.DataAvailabilityMode,
// 	D internal.DeclareTxnV3[TransactionType, TransactionVersion, u64, u128, RB, RBM, DA],
// ](d D) {}
