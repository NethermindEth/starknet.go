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

type InvokeTxnV3[
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
		Calldata              []*felt.Felt       `json:"calldata"`
		Version               TransactionVersion `json:"version"`
		Signature             []*felt.Felt       `json:"signature"`
		Nonce                 *felt.Felt         `json:"nonce"`
		ResourceBounds        *RBM               `json:"resource_bounds"`
		Tip                   u64                `json:"tip"`
		PayMasterData         []*felt.Felt       `json:"paymaster_data"`
		AccountDeploymentData []*felt.Felt       `json:"account_deployment_data"`
		NonceDataMode         DA                 `json:"nonce_data_availability_mode"`
		FeeMode               DA                 `json:"fee_data_availability_mode"`
	}
}
