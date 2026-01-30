package rpcv9

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/NethermindEth/juno/core/felt"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/NethermindEth/starknet.go/rpc/types"
)

type BlockTransaction struct {
	Hash *felt.Felt `json:"transaction_hash"`
	Transaction
}

// @changed removed FunctionCall and replaced it with its fields
type InvokeTxnV0 struct {
	Type               types.TransactionType    `json:"type"`
	MaxFee             *felt.Felt               `json:"max_fee"`
	Version            types.TransactionVersion `json:"version"`
	Signature          []*felt.Felt             `json:"signature"`
	ContractAddress    *felt.Felt               `json:"contract_address"`
	EntryPointSelector *felt.Felt               `json:"entry_point_selector"`
	Calldata           []*felt.Felt             `json:"calldata"`
}

// @new added methods for all fields in all transaction types
func (tx InvokeTxnV0) GetType() types.TransactionType       { return tx.Type }
func (tx InvokeTxnV0) GetMaxFee() *felt.Felt                { return tx.MaxFee }
func (tx InvokeTxnV0) GetVersion() types.TransactionVersion { return tx.Version }
func (tx InvokeTxnV0) GetSignature() []*felt.Felt           { return tx.Signature }
func (tx InvokeTxnV0) GetContractAddress() *felt.Felt       { return tx.ContractAddress }
func (tx InvokeTxnV0) GetEntryPointSelector() *felt.Felt    { return tx.EntryPointSelector }
func (tx InvokeTxnV0) GetCalldata() []*felt.Felt            { return tx.Calldata }

type InvokeTxnV1 struct {
	MaxFee        *felt.Felt               `json:"max_fee"`
	Version       types.TransactionVersion `json:"version"`
	Signature     []*felt.Felt             `json:"signature"`
	Nonce         *felt.Felt               `json:"nonce"`
	Type          types.TransactionType    `json:"type"`
	SenderAddress *felt.Felt               `json:"sender_address"`
	// The data expected by the account's `execute` function (in most usecases, this includes the
	// called contract address and a function selector)
	Calldata []*felt.Felt `json:"calldata"`
}

func (tx InvokeTxnV1) GetMaxFee() *felt.Felt                { return tx.MaxFee }
func (tx InvokeTxnV1) GetVersion() types.TransactionVersion { return tx.Version }
func (tx InvokeTxnV1) GetSignature() []*felt.Felt           { return tx.Signature }
func (tx InvokeTxnV1) GetNonce() *felt.Felt                 { return tx.Nonce }
func (tx InvokeTxnV1) GetType() types.TransactionType       { return tx.Type }
func (tx InvokeTxnV1) GetSenderAddress() *felt.Felt         { return tx.SenderAddress }
func (tx InvokeTxnV1) GetCalldata() []*felt.Felt            { return tx.Calldata }

type InvokeTxnV3 struct {
	Type           types.TransactionType    `json:"type"`
	SenderAddress  *felt.Felt               `json:"sender_address"`
	Calldata       []*felt.Felt             `json:"calldata"`
	Version        types.TransactionVersion `json:"version"`
	Signature      []*felt.Felt             `json:"signature"`
	Nonce          *felt.Felt               `json:"nonce"`
	ResourceBounds *ResourceBoundsMapping   `json:"resource_bounds"`
	Tip            types.U64                `json:"tip"`
	// The data needed to allow the paymaster to pay for the transaction in native tokens
	PayMasterData []*felt.Felt `json:"paymaster_data"`
	// The data needed to deploy the account contract from which this tx will be initiated
	AccountDeploymentData []*felt.Felt `json:"account_deployment_data"`
	// The storage domain of the account's nonce (an account has a nonce per DA mode)
	NonceDataMode types.DataAvailabilityMode `json:"nonce_data_availability_mode"`
	// The storage domain of the account's balance from which fee will be charged
	FeeMode types.DataAvailabilityMode `json:"fee_data_availability_mode"`
}

func (tx InvokeTxnV3) GetType() types.TransactionType               { return tx.Type }
func (tx InvokeTxnV3) GetSenderAddress() *felt.Felt                 { return tx.SenderAddress }
func (tx InvokeTxnV3) GetCalldata() []*felt.Felt                    { return tx.Calldata }
func (tx InvokeTxnV3) GetVersion() types.TransactionVersion         { return tx.Version }
func (tx InvokeTxnV3) GetSignature() []*felt.Felt                   { return tx.Signature }
func (tx InvokeTxnV3) GetNonce() *felt.Felt                         { return tx.Nonce }
func (tx InvokeTxnV3) GetResourceBounds() *ResourceBoundsMapping    { return tx.ResourceBounds }
func (tx InvokeTxnV3) GetTip() types.U64                            { return tx.Tip }
func (tx InvokeTxnV3) GetPayMasterData() []*felt.Felt               { return tx.PayMasterData }
func (tx InvokeTxnV3) GetAccountDeploymentData() []*felt.Felt       { return tx.AccountDeploymentData }
func (tx InvokeTxnV3) GetNonceDataMode() types.DataAvailabilityMode { return tx.NonceDataMode }
func (tx InvokeTxnV3) GetFeeMode() types.DataAvailabilityMode       { return tx.FeeMode }

type DeclareTxnV0 struct {
	Type types.TransactionType `json:"type"`
	// SenderAddress the address of the account contract sending the declaration transaction
	SenderAddress *felt.Felt               `json:"sender_address"`
	MaxFee        *felt.Felt               `json:"max_fee"`
	Version       types.TransactionVersion `json:"version"`
	Signature     []*felt.Felt             `json:"signature"`
	ClassHash     *felt.Felt               `json:"class_hash"`
}

func (tx DeclareTxnV0) GetType() types.TransactionType       { return tx.Type }
func (tx DeclareTxnV0) GetSenderAddress() *felt.Felt         { return tx.SenderAddress }
func (tx DeclareTxnV0) GetMaxFee() *felt.Felt                { return tx.MaxFee }
func (tx DeclareTxnV0) GetVersion() types.TransactionVersion { return tx.Version }
func (tx DeclareTxnV0) GetSignature() []*felt.Felt           { return tx.Signature }
func (tx DeclareTxnV0) GetClassHash() *felt.Felt             { return tx.ClassHash }

type DeclareTxnV1 struct {
	Type types.TransactionType `json:"type"`
	// SenderAddress the address of the account contract sending the declaration transaction
	SenderAddress *felt.Felt               `json:"sender_address"`
	MaxFee        *felt.Felt               `json:"max_fee"`
	Version       types.TransactionVersion `json:"version"`
	Signature     []*felt.Felt             `json:"signature"`
	Nonce         *felt.Felt               `json:"nonce"`
	// ClassHash the hash of the declared class
	ClassHash *felt.Felt `json:"class_hash"`
}

func (tx DeclareTxnV1) GetType() types.TransactionType       { return tx.Type }
func (tx DeclareTxnV1) GetSenderAddress() *felt.Felt         { return tx.SenderAddress }
func (tx DeclareTxnV1) GetMaxFee() *felt.Felt                { return tx.MaxFee }
func (tx DeclareTxnV1) GetVersion() types.TransactionVersion { return tx.Version }
func (tx DeclareTxnV1) GetSignature() []*felt.Felt           { return tx.Signature }
func (tx DeclareTxnV1) GetNonce() *felt.Felt                 { return tx.Nonce }
func (tx DeclareTxnV1) GetClassHash() *felt.Felt             { return tx.ClassHash }

type DeclareTxnV2 struct {
	Type types.TransactionType `json:"type"`
	// SenderAddress the address of the account contract sending the declaration transaction
	SenderAddress     *felt.Felt               `json:"sender_address"`
	CompiledClassHash *felt.Felt               `json:"compiled_class_hash"`
	MaxFee            *felt.Felt               `json:"max_fee"`
	Version           types.TransactionVersion `json:"version"`
	Signature         []*felt.Felt             `json:"signature"`
	Nonce             *felt.Felt               `json:"nonce"`
	ClassHash         *felt.Felt               `json:"class_hash"`
}

func (tx DeclareTxnV2) GetType() types.TransactionType       { return tx.Type }
func (tx DeclareTxnV2) GetSenderAddress() *felt.Felt         { return tx.SenderAddress }
func (tx DeclareTxnV2) GetCompiledClassHash() *felt.Felt     { return tx.CompiledClassHash }
func (tx DeclareTxnV2) GetMaxFee() *felt.Felt                { return tx.MaxFee }
func (tx DeclareTxnV2) GetVersion() types.TransactionVersion { return tx.Version }
func (tx DeclareTxnV2) GetSignature() []*felt.Felt           { return tx.Signature }
func (tx DeclareTxnV2) GetNonce() *felt.Felt                 { return tx.Nonce }
func (tx DeclareTxnV2) GetClassHash() *felt.Felt             { return tx.ClassHash }

type DeclareTxnV3 struct {
	Type              types.TransactionType    `json:"type"`
	SenderAddress     *felt.Felt               `json:"sender_address"`
	CompiledClassHash *felt.Felt               `json:"compiled_class_hash"`
	Version           types.TransactionVersion `json:"version"`
	Signature         []*felt.Felt             `json:"signature"`
	Nonce             *felt.Felt               `json:"nonce"`
	ClassHash         *felt.Felt               `json:"class_hash"`
	ResourceBounds    *ResourceBoundsMapping   `json:"resource_bounds"`
	Tip               types.U64                `json:"tip"`
	// The data needed to allow the paymaster to pay for the transaction in native tokens
	PayMasterData []*felt.Felt `json:"paymaster_data"`
	// The data needed to deploy the account contract from which this tx will be initiated
	AccountDeploymentData []*felt.Felt `json:"account_deployment_data"`
	// The storage domain of the account's nonce (an account has a nonce per DA mode)
	NonceDataMode types.DataAvailabilityMode `json:"nonce_data_availability_mode"`
	// The storage domain of the account's balance from which fee will be charged
	FeeMode types.DataAvailabilityMode `json:"fee_data_availability_mode"`
}

func (tx DeclareTxnV3) GetType() types.TransactionType               { return tx.Type }
func (tx DeclareTxnV3) GetSenderAddress() *felt.Felt                 { return tx.SenderAddress }
func (tx DeclareTxnV3) GetCompiledClassHash() *felt.Felt             { return tx.CompiledClassHash }
func (tx DeclareTxnV3) GetVersion() types.TransactionVersion         { return tx.Version }
func (tx DeclareTxnV3) GetSignature() []*felt.Felt                   { return tx.Signature }
func (tx DeclareTxnV3) GetNonce() *felt.Felt                         { return tx.Nonce }
func (tx DeclareTxnV3) GetClassHash() *felt.Felt                     { return tx.ClassHash }
func (tx DeclareTxnV3) GetResourceBounds() *ResourceBoundsMapping    { return tx.ResourceBounds }
func (tx DeclareTxnV3) GetTip() types.U64                            { return tx.Tip }
func (tx DeclareTxnV3) GetPayMasterData() []*felt.Felt               { return tx.PayMasterData }
func (tx DeclareTxnV3) GetAccountDeploymentData() []*felt.Felt       { return tx.AccountDeploymentData }
func (tx DeclareTxnV3) GetNonceDataMode() types.DataAvailabilityMode { return tx.NonceDataMode }
func (tx DeclareTxnV3) GetFeeMode() types.DataAvailabilityMode       { return tx.FeeMode }

// DeployAccountTxnV1 The structure of a deployAccount transaction.
type DeployAccountTxnV1 struct {
	MaxFee    *felt.Felt               `json:"max_fee"`
	Version   types.TransactionVersion `json:"version"`
	Signature []*felt.Felt             `json:"signature"`
	Nonce     *felt.Felt               `json:"nonce"`
	Type      types.TransactionType    `json:"type"`
	// ClassHash The hash of the deployed contract's class
	ClassHash *felt.Felt `json:"class_hash"`
	// ContractAddressSalt The salt for the address of the deployed contract
	ContractAddressSalt *felt.Felt `json:"contract_address_salt"`
	// ConstructorCalldata The parameters passed to the constructor
	ConstructorCalldata []*felt.Felt `json:"constructor_calldata"`
}

func (tx DeployAccountTxnV1) GetMaxFee() *felt.Felt                { return tx.MaxFee }
func (tx DeployAccountTxnV1) GetVersion() types.TransactionVersion { return tx.Version }
func (tx DeployAccountTxnV1) GetSignature() []*felt.Felt           { return tx.Signature }
func (tx DeployAccountTxnV1) GetNonce() *felt.Felt                 { return tx.Nonce }
func (tx DeployAccountTxnV1) GetType() types.TransactionType       { return tx.Type }
func (tx DeployAccountTxnV1) GetClassHash() *felt.Felt             { return tx.ClassHash }
func (tx DeployAccountTxnV1) GetContractAddressSalt() *felt.Felt   { return tx.ContractAddressSalt }
func (tx DeployAccountTxnV1) GetConstructorCalldata() []*felt.Felt { return tx.ConstructorCalldata }

type DeployAccountTxnV3 struct {
	Type                types.TransactionType    `json:"type"`
	Version             types.TransactionVersion `json:"version"`
	Signature           []*felt.Felt             `json:"signature"`
	Nonce               *felt.Felt               `json:"nonce"`
	ContractAddressSalt *felt.Felt               `json:"contract_address_salt"`
	ConstructorCalldata []*felt.Felt             `json:"constructor_calldata"`
	ClassHash           *felt.Felt               `json:"class_hash"`
	ResourceBounds      *ResourceBoundsMapping   `json:"resource_bounds"`
	Tip                 types.U64                `json:"tip"`
	// The data needed to allow the paymaster to pay for the transaction in native tokens
	PayMasterData []*felt.Felt `json:"paymaster_data"`
	// The storage domain of the account's nonce (an account has a nonce per DA mode)
	NonceDataMode types.DataAvailabilityMode `json:"nonce_data_availability_mode"`
	// The storage domain of the account's balance from which fee will be charged
	FeeMode types.DataAvailabilityMode `json:"fee_data_availability_mode"`
}

func (tx DeployAccountTxnV3) GetType() types.TransactionType               { return tx.Type }
func (tx DeployAccountTxnV3) GetVersion() types.TransactionVersion         { return tx.Version }
func (tx DeployAccountTxnV3) GetSignature() []*felt.Felt                   { return tx.Signature }
func (tx DeployAccountTxnV3) GetNonce() *felt.Felt                         { return tx.Nonce }
func (tx DeployAccountTxnV3) GetContractAddressSalt() *felt.Felt           { return tx.ContractAddressSalt }
func (tx DeployAccountTxnV3) GetConstructorCalldata() []*felt.Felt         { return tx.ConstructorCalldata }
func (tx DeployAccountTxnV3) GetClassHash() *felt.Felt                     { return tx.ClassHash }
func (tx DeployAccountTxnV3) GetResourceBounds() *ResourceBoundsMapping    { return tx.ResourceBounds }
func (tx DeployAccountTxnV3) GetTip() types.U64                            { return tx.Tip }
func (tx DeployAccountTxnV3) GetPayMasterData() []*felt.Felt               { return tx.PayMasterData }
func (tx DeployAccountTxnV3) GetNonceDataMode() types.DataAvailabilityMode { return tx.NonceDataMode }
func (tx DeployAccountTxnV3) GetFeeMode() types.DataAvailabilityMode       { return tx.FeeMode }

// DeployTxn The structure of a deploy transaction. Note that this transaction type
// is deprecated and will no longer be supported in future versions
type DeployTxn struct {
	// ClassHash The hash of the deployed contract's class
	ClassHash           *felt.Felt               `json:"class_hash"`
	Version             types.TransactionVersion `json:"version"`
	Type                types.TransactionType    `json:"type"`
	ContractAddressSalt *felt.Felt               `json:"contract_address_salt"`
	ConstructorCalldata []*felt.Felt             `json:"constructor_calldata"`
}

func (tx DeployTxn) GetClassHash() *felt.Felt             { return tx.ClassHash }
func (tx DeployTxn) GetVersion() types.TransactionVersion { return tx.Version }
func (tx DeployTxn) GetType() types.TransactionType       { return tx.Type }
func (tx DeployTxn) GetContractAddressSalt() *felt.Felt   { return tx.ContractAddressSalt }
func (tx DeployTxn) GetConstructorCalldata() []*felt.Felt { return tx.ConstructorCalldata }

// @changed removed FunctionCall and replaced it with its fields
type L1HandlerTxn struct {
	Type               types.TransactionType    `json:"type"`
	Version            types.TransactionVersion `json:"version"`
	Nonce              string                   `json:"nonce"`
	ContractAddress    *felt.Felt               `json:"contract_address"`
	EntryPointSelector *felt.Felt               `json:"entry_point_selector"`
	Calldata           []*felt.Felt             `json:"calldata"`
}

func (tx L1HandlerTxn) GetType() types.TransactionType       { return tx.Type }
func (tx L1HandlerTxn) GetVersion() types.TransactionVersion { return tx.Version }
func (tx L1HandlerTxn) GetNonce() string                     { return tx.Nonce }
func (tx L1HandlerTxn) GetContractAddress() *felt.Felt       { return tx.ContractAddress }
func (tx L1HandlerTxn) GetEntryPointSelector() *felt.Felt    { return tx.EntryPointSelector }
func (tx L1HandlerTxn) GetCalldata() []*felt.Felt            { return tx.Calldata }

type ResourceBoundsMapping struct {
	// The max amount and max price per unit of L1 gas used in this tx
	L1Gas ResourceBounds `json:"l1_gas"`
	// The max amount and max price per unit of L1 blob gas used in this tx
	L1DataGas ResourceBounds `json:"l1_data_gas"`
	// The max amount and max price per unit of L2 gas used in this tx
	L2Gas ResourceBounds `json:"l2_gas"`
}

type ResourceBounds struct {
	// The max amount of the resource that can be used in the tx
	MaxAmount types.U64 `json:"max_amount"`
	// The max price per unit of this resource for this tx
	MaxPricePerUnit types.U128 `json:"max_price_per_unit"`
}

func (rb ResourceBounds) Bytes(resource types.Resource) ([]byte, error) {
	const eight = 8
	maxAmountBytes := make([]byte, eight)
	maxAmountUint64, err := rb.MaxAmount.ToUint64()
	if err != nil {
		return nil, err
	}
	binary.BigEndian.PutUint64(maxAmountBytes, maxAmountUint64)
	maxPricePerUnitFelt, err := new(felt.Felt).SetString(string(rb.MaxPricePerUnit))
	if err != nil {
		return nil, err
	}
	maxPriceBytes := maxPricePerUnitFelt.Bytes()

	return internalUtils.Flatten(
		[]byte{0},
		[]byte(resource),
		maxAmountBytes,
		maxPriceBytes[16:], // uint128.
	), nil
}

// SubPendingTxnsInput is the optional input of the
// starknet_subscribeNewTransactionReceipts subscription.
type SubNewTxnReceiptsInput struct {
	// Optional: A vector of finality statuses to receive updates for.
	// Only `PRE_CONFIRMED` and `ACCEPTED_ON_L2` are supported. Default is
	// `ACCEPTED_ON_L2`.
	FinalityStatus []types.TxnFinalityStatus `json:"finality_status,omitempty"`
	// Optional: Filter transaction receipts to only include transactions
	// sent by the specified addresses
	SenderAddress []*felt.Felt `json:"sender_address,omitempty"`
}

// SubNewTxnsInput is the optional input of the
// starknet_subscribeNewTransactions subscription.
type SubNewTxnsInput struct {
	// Optional: A vector of finality statuses to receive updates for.
	// Support all transaction statuses, except `ACCEPTED_ON_L1`. Default is
	// `ACCEPTED_ON_L2`.
	FinalityStatus []types.TxnStatus `json:"finality_status,omitempty"`
	// Optional: Filter transaction receipts to only include transactions sent
	// by the specified addresses
	SenderAddress []*felt.Felt `json:"sender_address,omitempty"`
}

// TxnWithHashAndStatus is the response of the
// starknet_subscribeNewTransactions subscription.
type TxnWithHashAndStatus struct {
	// Transaction with hash and status
	BlockTransaction
	// Finality status of the transaction, except `ACCEPTED_ON_L1`.
	FinalityStatus types.TxnStatus `json:"finality_status"`
}

func (txn *TxnWithHashAndStatus) UnmarshalJSON(data []byte) error {
	// type alias TxnWithHashAndStatus
	var aux BlockTransaction

	err := json.Unmarshal(data, &aux)
	if err != nil {
		return err
	}

	var aux2 struct {
		FinalityStatus types.TxnStatus `json:"finality_status"`
	}

	err = json.Unmarshal(data, &aux2)
	if err != nil {
		return err
	}

	txn.BlockTransaction = aux
	txn.FinalityStatus = aux2.FinalityStatus

	return nil
}

// MarshalJSON marshals the TxnWithHashAndStatus object into a JSON byte slice.
func (txn *TxnWithHashAndStatus) MarshalJSON() ([]byte, error) {
	blockTxnData, err := json.Marshal(&txn.BlockTransaction)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(blockTxnData, &result); err != nil {
		return nil, err
	}

	result["finality_status"] = txn.FinalityStatus

	return json.Marshal(result)
}

// UnmarshalJSON unmarshals the data into a BlockTransaction object.
//
// It takes a byte slice as the parameter, representing the JSON data to be
// unmarshalled.
// The function returns an error if the unmarshalling process fails.
//
// Parameters:
//   - data: The JSON data to be unmarshalled
//
// Returns:
//   - error: An error if the unmarshalling process fails
func (blockTxn *BlockTransaction) UnmarshalJSON(data []byte) error {
	type alias BlockTransaction
	var aux alias

	err := json.Unmarshal(data, &aux)
	if err != nil {
		return err
	}

	txn, err := unmarshalTxn(data)
	if err != nil {
		return err
	}

	blockTxn.Hash = aux.Hash
	blockTxn.Transaction = txn

	return nil
}

// MarshalJSON marshals the BlockTransaction object into a JSON byte slice.
//
// It takes a pointer to a BlockTransaction object as the parameter.
// The function returns a byte slice representing the JSON data and an error if
// the marshalling process fails.
func (blockTxn *BlockTransaction) MarshalJSON() ([]byte, error) {
	// First marshal the transaction to get all its fields
	txnData, err := json.Marshal(blockTxn.Transaction)
	if err != nil {
		return nil, err
	}

	// Unmarshal into a map to add the hash field
	var result map[string]interface{}
	if err := json.Unmarshal(txnData, &result); err != nil {
		return nil, err
	}

	// Add the hash field
	result["transaction_hash"] = blockTxn.Hash

	return json.Marshal(result)
}

// unmarshalTxn unmarshals a given txn as a byte slice and returns a concrete
// transaction type wrapped in the Transaction interface.
//
// Parameters:
//   - data: The transaction to be unmarshaled
//
// Returns:
//   - Transaction: a concrete transaction type wrapped in the Transaction interface
//   - error: an error if the unmarshaling process fails
//
//nolint:gocyclo // Inevitable due to many switch cases
func unmarshalTxn(data []byte) (Transaction, error) {
	var txnAsMap map[string]interface{}
	if err := json.Unmarshal(data, &txnAsMap); err != nil {
		return nil, err
	}

	switch types.TransactionType(txnAsMap["type"].(string)) {
	case types.TransactionTypeDeclare:
		switch types.TransactionVersion(txnAsMap["version"].(string)) {
		case types.TransactionV0:
			return unmarshalTxnToType[DeclareTxnV0](data)
		case types.TransactionV1:
			return unmarshalTxnToType[DeclareTxnV1](data)
		case types.TransactionV2:
			return unmarshalTxnToType[DeclareTxnV2](data)
		case types.TransactionV3:
			return unmarshalTxnToType[DeclareTxnV3](data)
		default:
			return nil, errors.New(
				"internal error with Declare transaction version and unmarshalTxn()",
			)
		}
	case types.TransactionTypeDeploy:
		return unmarshalTxnToType[DeployTxn](data)
	case types.TransactionTypeDeployAccount:
		switch types.TransactionVersion(txnAsMap["version"].(string)) {
		case types.TransactionV1:
			return unmarshalTxnToType[DeployAccountTxnV1](data)
		case types.TransactionV3:
			return unmarshalTxnToType[DeployAccountTxnV3](data)
		}
	case types.TransactionTypeInvoke:
		switch types.TransactionVersion(txnAsMap["version"].(string)) {
		case types.TransactionV0:
			return unmarshalTxnToType[InvokeTxnV0](data)
		case types.TransactionV1:
			return unmarshalTxnToType[InvokeTxnV1](data)
		case types.TransactionV3:
			return unmarshalTxnToType[InvokeTxnV3](data)
		}
	case types.TransactionTypeL1Handler:
		return unmarshalTxnToType[L1HandlerTxn](data)
	}

	return nil, fmt.Errorf("unknown transaction type: %v", txnAsMap["type"])
}

// unmarshalTxnToType is a generic function that takes in a byte slice 'data',
// unmarshals it to a concrete transaction of type T, and returns the concrete
// transaction wrapped in the Transaction interface.
func unmarshalTxnToType[T Transaction](data []byte) (T, error) {
	var resp T

	if err := json.Unmarshal(data, &resp); err != nil {
		return resp, err
	}

	return resp, nil
}
