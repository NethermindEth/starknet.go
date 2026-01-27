package constraints

import "github.com/NethermindEth/juno/core/felt"

// @todo add comments for everything in the pkg
type DeclareTxnV0Interface[
	TransactionType ~string,
	TransactionVersion ~string,
] interface {
	GetType() TransactionType
	GetSenderAddress() *felt.Felt
	GetMaxFee() *felt.Felt
	GetVersion() TransactionVersion
	GetSignature() []*felt.Felt
	GetClassHash() *felt.Felt
}

type DeclareTxnV1Interface[
	TransactionType ~string,
	TransactionVersion ~string,
] interface {
	GetType() TransactionType
	GetSenderAddress() *felt.Felt
	GetMaxFee() *felt.Felt
	GetVersion() TransactionVersion
	GetSignature() []*felt.Felt
	GetNonce() *felt.Felt
	GetClassHash() *felt.Felt
}

type DeclareTxnV2Interface[
	TransactionType ~string,
	TransactionVersion ~string,
] interface {
	GetType() TransactionType
	GetSenderAddress() *felt.Felt
	GetCompiledClassHash() *felt.Felt
	GetMaxFee() *felt.Felt
	GetVersion() TransactionVersion
	GetSignature() []*felt.Felt
	GetNonce() *felt.Felt
	GetClassHash() *felt.Felt
}

type DeclareTxnV3Interface[
	TransactionType ~string,
	TransactionVersion ~string,
	u64 U64,
	u128 U128,
	RB ResourceBounds[u64, u128],
	RBM ResourceBoundsMapping[u64, u128, RB],
	DA DataAvailabilityMode,
] interface {
	GetType() TransactionType
	GetSenderAddress() *felt.Felt
	GetCompiledClassHash() *felt.Felt
	GetVersion() TransactionVersion
	GetSignature() []*felt.Felt
	GetNonce() *felt.Felt
	GetResourceBounds() *RBM
	GetTip() u64
	GetPayMasterData() []*felt.Felt
	GetAccountDeploymentData() []*felt.Felt
	GetNonceDataMode() DA
	GetFeeMode() DA

	// The "ClassHash" is not included in the interface because the BroadcastDeclareTxnV3
	// does not include it, it contains the "ContractClass" instead. So, we only kept
	// the common fields for both types.

	// GetClassHash() *felt.Felt
}
