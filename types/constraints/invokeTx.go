package constraints

import "github.com/NethermindEth/juno/core/felt"

type InvokeTxnV0Interface[
	TransactionType ~string,
	TransactionVersion ~string,
] interface {
	GetType() TransactionType
	GetMaxFee() *felt.Felt
	GetVersion() TransactionVersion
	GetSignature() []*felt.Felt
	GetContractAddress() *felt.Felt
	GetEntryPointSelector() *felt.Felt
	GetCalldata() []*felt.Felt
}

type InvokeTxnV1Interface[
	TransactionType ~string,
	TransactionVersion ~string,
] interface {
	GetMaxFee() *felt.Felt
	GetVersion() TransactionVersion
	GetSignature() []*felt.Felt
	GetNonce() *felt.Felt
	GetType() TransactionType
	GetSenderAddress() *felt.Felt
	GetCalldata() []*felt.Felt
}

type InvokeTxnV3Interface[
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
	GetCalldata() []*felt.Felt
	GetVersion() TransactionVersion
	GetSignature() []*felt.Felt
	GetNonce() *felt.Felt
	GetResourceBounds() *RBM
	GetTip() u64
	GetPayMasterData() []*felt.Felt
	GetAccountDeploymentData() []*felt.Felt
	GetNonceDataMode() DA
	GetFeeMode() DA
}
