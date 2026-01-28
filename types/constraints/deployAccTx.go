package constraints

import "github.com/NethermindEth/juno/core/felt"

type DeployAccountTxnV1Interface[
	TransactionType ~string,
	TransactionVersion ~string,
] interface {
	GetMaxFee() *felt.Felt
	GetVersion() TransactionVersion
	GetSignature() []*felt.Felt
	GetNonce() *felt.Felt
	GetType() TransactionType
	GetClassHash() *felt.Felt
	GetContractAddressSalt() *felt.Felt
	GetConstructorCalldata() []*felt.Felt
}

type DeployAccountTxnV3Interface[
	TransactionType ~string,
	TransactionVersion ~string,
	u64 U64,
	u128 U128,
	RB ResourceBounds[u64, u128],
	RBM ResourceBoundsMapping[u64, u128, RB],
	DA DataAvailabilityMode,
] interface {
	GetType() TransactionType
	GetVersion() TransactionVersion
	GetSignature() []*felt.Felt
	GetNonce() *felt.Felt
	GetContractAddressSalt() *felt.Felt
	GetConstructorCalldata() []*felt.Felt
	GetClassHash() *felt.Felt
	GetResourceBounds() *RBM
	GetTip() u64
	GetPayMasterData() []*felt.Felt
	GetNonceDataMode() DA
	GetFeeMode() DA
}

type DeployAccountTxnV3[
	TransactionType ~string,
	TransactionVersion ~string,
	u64 U64,
	u128 U128,
	RB ResourceBounds[u64, u128],
	RBM ResourceBoundsMapping[u64, u128, RB],
	DA DataAvailabilityMode,
] interface {
	~struct {
		Type                TransactionType    `json:"type"`
		Version             TransactionVersion `json:"version"`
		Signature           []*felt.Felt       `json:"signature"`
		Nonce               *felt.Felt         `json:"nonce"`
		ContractAddressSalt *felt.Felt         `json:"contract_address_salt"`
		ConstructorCalldata []*felt.Felt       `json:"constructor_calldata"`
		ClassHash           *felt.Felt         `json:"class_hash"`
		ResourceBounds      *RBM               `json:"resource_bounds"`
		Tip                 u64                `json:"tip"`
		PayMasterData       []*felt.Felt       `json:"paymaster_data"`
		NonceDataMode       DA                 `json:"nonce_data_availability_mode"`
		FeeMode             DA                 `json:"fee_data_availability_mode"`
	}
}
