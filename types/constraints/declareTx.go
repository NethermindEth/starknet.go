package constraints

import (
	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/contracts"
	"github.com/NethermindEth/starknet.go/rpc/types"
)

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
	TxType, TxVersion ~string,
	RB ResourceBounds,
	RBM ResourceBoundsMapping[RB],
	DA DataAvailabilityMode,
] interface {
	GetType() TxType
	GetSenderAddress() *felt.Felt
	GetCompiledClassHash() *felt.Felt
	GetVersion() TxVersion
	GetSignature() []*felt.Felt
	GetNonce() *felt.Felt
	GetResourceBounds() *RBM
	GetTip() types.U64
	GetPayMasterData() []*felt.Felt
	GetAccountDeploymentData() []*felt.Felt
	GetNonceDataMode() DA
	GetFeeMode() DA

	// The "ClassHash" is not included in the interface because the BroadcastDeclareTxnV3
	// does not include it, it contains the "ContractClass" instead. So, we only kept
	// the common fields for both types.

	// GetClassHash() *felt.Felt
}

type BroadcastDeclareTxnV3[
	TxType, TxVersion ~string,
	RB ResourceBounds,
	RBM ResourceBoundsMapping[RB],
	DA DataAvailabilityMode,
] interface {
	~struct {
		Type                  TxType                   `json:"type"`
		SenderAddress         *felt.Felt               `json:"sender_address"`
		CompiledClassHash     *felt.Felt               `json:"compiled_class_hash"`
		Version               TxVersion                `json:"version"`
		Signature             []*felt.Felt             `json:"signature"`
		Nonce                 *felt.Felt               `json:"nonce"`
		ContractClass         *contracts.ContractClass `json:"contract_class"`
		ResourceBounds        *RBM                     `json:"resource_bounds"`
		Tip                   types.U64                `json:"tip"`
		PayMasterData         []*felt.Felt             `json:"paymaster_data"`
		AccountDeploymentData []*felt.Felt             `json:"account_deployment_data"`
		NonceDataMode         DA                       `json:"nonce_data_availability_mode"`
		FeeMode               DA                       `json:"fee_data_availability_mode"`
	}
}
