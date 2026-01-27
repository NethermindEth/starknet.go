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
