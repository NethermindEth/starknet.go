package rpcv10

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/NethermindEth/juno/core/felt"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
)

type BlockTransaction struct {
	Hash *felt.Felt `json:"transaction_hash"`
	Transaction
}

// @changed removed FunctionCall and replaced it with its fields
type InvokeTxnV0 struct {
	Type               TransactionType    `json:"type"`
	MaxFee             *felt.Felt         `json:"max_fee"`
	Version            TransactionVersion `json:"version"`
	Signature          []*felt.Felt       `json:"signature"`
	ContractAddress    *felt.Felt         `json:"contract_address"`
	EntryPointSelector *felt.Felt         `json:"entry_point_selector"`
	Calldata           []*felt.Felt       `json:"calldata"`
}

// @new added methods for all fields in all transaction types
func (tx InvokeTxnV0) GetType() TransactionType          { return tx.Type }
func (tx InvokeTxnV0) GetMaxFee() *felt.Felt             { return tx.MaxFee }
func (tx InvokeTxnV0) GetVersion() TransactionVersion    { return tx.Version }
func (tx InvokeTxnV0) GetSignature() []*felt.Felt        { return tx.Signature }
func (tx InvokeTxnV0) GetContractAddress() *felt.Felt    { return tx.ContractAddress }
func (tx InvokeTxnV0) GetEntryPointSelector() *felt.Felt { return tx.EntryPointSelector }
func (tx InvokeTxnV0) GetCalldata() []*felt.Felt         { return tx.Calldata }

type InvokeTxnV1 struct {
	MaxFee        *felt.Felt         `json:"max_fee"`
	Version       TransactionVersion `json:"version"`
	Signature     []*felt.Felt       `json:"signature"`
	Nonce         *felt.Felt         `json:"nonce"`
	Type          TransactionType    `json:"type"`
	SenderAddress *felt.Felt         `json:"sender_address"`
	// The data expected by the account's `execute` function (in most usecases, this includes the
	// called contract address and a function selector)
	Calldata []*felt.Felt `json:"calldata"`
}

func (tx InvokeTxnV1) GetMaxFee() *felt.Felt          { return tx.MaxFee }
func (tx InvokeTxnV1) GetVersion() TransactionVersion { return tx.Version }
func (tx InvokeTxnV1) GetSignature() []*felt.Felt     { return tx.Signature }
func (tx InvokeTxnV1) GetNonce() *felt.Felt           { return tx.Nonce }
func (tx InvokeTxnV1) GetType() TransactionType       { return tx.Type }
func (tx InvokeTxnV1) GetSenderAddress() *felt.Felt   { return tx.SenderAddress }
func (tx InvokeTxnV1) GetCalldata() []*felt.Felt      { return tx.Calldata }

type InvokeTxnV3 struct {
	Type           TransactionType        `json:"type"`
	SenderAddress  *felt.Felt             `json:"sender_address"`
	Calldata       []*felt.Felt           `json:"calldata"`
	Version        TransactionVersion     `json:"version"`
	Signature      []*felt.Felt           `json:"signature"`
	Nonce          *felt.Felt             `json:"nonce"`
	ResourceBounds *ResourceBoundsMapping `json:"resource_bounds"`
	Tip            U64                    `json:"tip"`
	// The data needed to allow the paymaster to pay for the transaction in native tokens
	PayMasterData []*felt.Felt `json:"paymaster_data"`
	// The data needed to deploy the account contract from which this tx will be initiated
	AccountDeploymentData []*felt.Felt `json:"account_deployment_data"`
	// The storage domain of the account's nonce (an account has a nonce per DA mode)
	NonceDataMode DataAvailabilityMode `json:"nonce_data_availability_mode"`
	// The storage domain of the account's balance from which fee will be charged
	FeeMode DataAvailabilityMode `json:"fee_data_availability_mode"`
}

func (tx InvokeTxnV3) GetType() TransactionType                  { return tx.Type }
func (tx InvokeTxnV3) GetSenderAddress() *felt.Felt              { return tx.SenderAddress }
func (tx InvokeTxnV3) GetCalldata() []*felt.Felt                 { return tx.Calldata }
func (tx InvokeTxnV3) GetVersion() TransactionVersion            { return tx.Version }
func (tx InvokeTxnV3) GetSignature() []*felt.Felt                { return tx.Signature }
func (tx InvokeTxnV3) GetNonce() *felt.Felt                      { return tx.Nonce }
func (tx InvokeTxnV3) GetResourceBounds() *ResourceBoundsMapping { return tx.ResourceBounds }
func (tx InvokeTxnV3) GetTip() U64                               { return tx.Tip }
func (tx InvokeTxnV3) GetPayMasterData() []*felt.Felt            { return tx.PayMasterData }
func (tx InvokeTxnV3) GetAccountDeploymentData() []*felt.Felt    { return tx.AccountDeploymentData }
func (tx InvokeTxnV3) GetNonceDataMode() DataAvailabilityMode    { return tx.NonceDataMode }
func (tx InvokeTxnV3) GetFeeMode() DataAvailabilityMode          { return tx.FeeMode }

type DeclareTxnV0 struct {
	Type TransactionType `json:"type"`
	// SenderAddress the address of the account contract sending the declaration transaction
	SenderAddress *felt.Felt         `json:"sender_address"`
	MaxFee        *felt.Felt         `json:"max_fee"`
	Version       TransactionVersion `json:"version"`
	Signature     []*felt.Felt       `json:"signature"`
	ClassHash     *felt.Felt         `json:"class_hash"`
}

func (tx DeclareTxnV0) GetType() TransactionType       { return tx.Type }
func (tx DeclareTxnV0) GetSenderAddress() *felt.Felt   { return tx.SenderAddress }
func (tx DeclareTxnV0) GetMaxFee() *felt.Felt          { return tx.MaxFee }
func (tx DeclareTxnV0) GetVersion() TransactionVersion { return tx.Version }
func (tx DeclareTxnV0) GetSignature() []*felt.Felt     { return tx.Signature }
func (tx DeclareTxnV0) GetClassHash() *felt.Felt       { return tx.ClassHash }

type DeclareTxnV1 struct {
	Type TransactionType `json:"type"`
	// SenderAddress the address of the account contract sending the declaration transaction
	SenderAddress *felt.Felt         `json:"sender_address"`
	MaxFee        *felt.Felt         `json:"max_fee"`
	Version       TransactionVersion `json:"version"`
	Signature     []*felt.Felt       `json:"signature"`
	Nonce         *felt.Felt         `json:"nonce"`
	// ClassHash the hash of the declared class
	ClassHash *felt.Felt `json:"class_hash"`
}

func (tx DeclareTxnV1) GetType() TransactionType       { return tx.Type }
func (tx DeclareTxnV1) GetSenderAddress() *felt.Felt   { return tx.SenderAddress }
func (tx DeclareTxnV1) GetMaxFee() *felt.Felt          { return tx.MaxFee }
func (tx DeclareTxnV1) GetVersion() TransactionVersion { return tx.Version }
func (tx DeclareTxnV1) GetSignature() []*felt.Felt     { return tx.Signature }
func (tx DeclareTxnV1) GetNonce() *felt.Felt           { return tx.Nonce }
func (tx DeclareTxnV1) GetClassHash() *felt.Felt       { return tx.ClassHash }

type DeclareTxnV2 struct {
	Type TransactionType `json:"type"`
	// SenderAddress the address of the account contract sending the declaration transaction
	SenderAddress     *felt.Felt         `json:"sender_address"`
	CompiledClassHash *felt.Felt         `json:"compiled_class_hash"`
	MaxFee            *felt.Felt         `json:"max_fee"`
	Version           TransactionVersion `json:"version"`
	Signature         []*felt.Felt       `json:"signature"`
	Nonce             *felt.Felt         `json:"nonce"`
	ClassHash         *felt.Felt         `json:"class_hash"`
}

func (tx DeclareTxnV2) GetType() TransactionType         { return tx.Type }
func (tx DeclareTxnV2) GetSenderAddress() *felt.Felt     { return tx.SenderAddress }
func (tx DeclareTxnV2) GetCompiledClassHash() *felt.Felt { return tx.CompiledClassHash }
func (tx DeclareTxnV2) GetMaxFee() *felt.Felt            { return tx.MaxFee }
func (tx DeclareTxnV2) GetVersion() TransactionVersion   { return tx.Version }
func (tx DeclareTxnV2) GetSignature() []*felt.Felt       { return tx.Signature }
func (tx DeclareTxnV2) GetNonce() *felt.Felt             { return tx.Nonce }
func (tx DeclareTxnV2) GetClassHash() *felt.Felt         { return tx.ClassHash }

type DeclareTxnV3 struct {
	Type              TransactionType        `json:"type"`
	SenderAddress     *felt.Felt             `json:"sender_address"`
	CompiledClassHash *felt.Felt             `json:"compiled_class_hash"`
	Version           TransactionVersion     `json:"version"`
	Signature         []*felt.Felt           `json:"signature"`
	Nonce             *felt.Felt             `json:"nonce"`
	ClassHash         *felt.Felt             `json:"class_hash"`
	ResourceBounds    *ResourceBoundsMapping `json:"resource_bounds"`
	Tip               U64                    `json:"tip"`
	// The data needed to allow the paymaster to pay for the transaction in native tokens
	PayMasterData []*felt.Felt `json:"paymaster_data"`
	// The data needed to deploy the account contract from which this tx will be initiated
	AccountDeploymentData []*felt.Felt `json:"account_deployment_data"`
	// The storage domain of the account's nonce (an account has a nonce per DA mode)
	NonceDataMode DataAvailabilityMode `json:"nonce_data_availability_mode"`
	// The storage domain of the account's balance from which fee will be charged
	FeeMode DataAvailabilityMode `json:"fee_data_availability_mode"`
}

func (tx DeclareTxnV3) GetType() TransactionType                  { return tx.Type }
func (tx DeclareTxnV3) GetSenderAddress() *felt.Felt              { return tx.SenderAddress }
func (tx DeclareTxnV3) GetCompiledClassHash() *felt.Felt          { return tx.CompiledClassHash }
func (tx DeclareTxnV3) GetVersion() TransactionVersion            { return tx.Version }
func (tx DeclareTxnV3) GetSignature() []*felt.Felt                { return tx.Signature }
func (tx DeclareTxnV3) GetNonce() *felt.Felt                      { return tx.Nonce }
func (tx DeclareTxnV3) GetClassHash() *felt.Felt                  { return tx.ClassHash }
func (tx DeclareTxnV3) GetResourceBounds() *ResourceBoundsMapping { return tx.ResourceBounds }
func (tx DeclareTxnV3) GetTip() U64                               { return tx.Tip }
func (tx DeclareTxnV3) GetPayMasterData() []*felt.Felt            { return tx.PayMasterData }
func (tx DeclareTxnV3) GetAccountDeploymentData() []*felt.Felt    { return tx.AccountDeploymentData }
func (tx DeclareTxnV3) GetNonceDataMode() DataAvailabilityMode    { return tx.NonceDataMode }
func (tx DeclareTxnV3) GetFeeMode() DataAvailabilityMode          { return tx.FeeMode }

// DeployAccountTxnV1 The structure of a deployAccount transaction.
type DeployAccountTxnV1 struct {
	MaxFee    *felt.Felt         `json:"max_fee"`
	Version   TransactionVersion `json:"version"`
	Signature []*felt.Felt       `json:"signature"`
	Nonce     *felt.Felt         `json:"nonce"`
	Type      TransactionType    `json:"type"`
	// ClassHash The hash of the deployed contract's class
	ClassHash *felt.Felt `json:"class_hash"`
	// ContractAddressSalt The salt for the address of the deployed contract
	ContractAddressSalt *felt.Felt `json:"contract_address_salt"`
	// ConstructorCalldata The parameters passed to the constructor
	ConstructorCalldata []*felt.Felt `json:"constructor_calldata"`
}

func (tx DeployAccountTxnV1) GetMaxFee() *felt.Felt                { return tx.MaxFee }
func (tx DeployAccountTxnV1) GetVersion() TransactionVersion       { return tx.Version }
func (tx DeployAccountTxnV1) GetSignature() []*felt.Felt           { return tx.Signature }
func (tx DeployAccountTxnV1) GetNonce() *felt.Felt                 { return tx.Nonce }
func (tx DeployAccountTxnV1) GetType() TransactionType             { return tx.Type }
func (tx DeployAccountTxnV1) GetClassHash() *felt.Felt             { return tx.ClassHash }
func (tx DeployAccountTxnV1) GetContractAddressSalt() *felt.Felt   { return tx.ContractAddressSalt }
func (tx DeployAccountTxnV1) GetConstructorCalldata() []*felt.Felt { return tx.ConstructorCalldata }

type DeployAccountTxnV3 struct {
	Type                TransactionType        `json:"type"`
	Version             TransactionVersion     `json:"version"`
	Signature           []*felt.Felt           `json:"signature"`
	Nonce               *felt.Felt             `json:"nonce"`
	ContractAddressSalt *felt.Felt             `json:"contract_address_salt"`
	ConstructorCalldata []*felt.Felt           `json:"constructor_calldata"`
	ClassHash           *felt.Felt             `json:"class_hash"`
	ResourceBounds      *ResourceBoundsMapping `json:"resource_bounds"`
	Tip                 U64                    `json:"tip"`
	// The data needed to allow the paymaster to pay for the transaction in native tokens
	PayMasterData []*felt.Felt `json:"paymaster_data"`
	// The storage domain of the account's nonce (an account has a nonce per DA mode)
	NonceDataMode DataAvailabilityMode `json:"nonce_data_availability_mode"`
	// The storage domain of the account's balance from which fee will be charged
	FeeMode DataAvailabilityMode `json:"fee_data_availability_mode"`
}

func (tx DeployAccountTxnV3) GetType() TransactionType                  { return tx.Type }
func (tx DeployAccountTxnV3) GetVersion() TransactionVersion            { return tx.Version }
func (tx DeployAccountTxnV3) GetSignature() []*felt.Felt                { return tx.Signature }
func (tx DeployAccountTxnV3) GetNonce() *felt.Felt                      { return tx.Nonce }
func (tx DeployAccountTxnV3) GetContractAddressSalt() *felt.Felt        { return tx.ContractAddressSalt }
func (tx DeployAccountTxnV3) GetConstructorCalldata() []*felt.Felt      { return tx.ConstructorCalldata }
func (tx DeployAccountTxnV3) GetClassHash() *felt.Felt                  { return tx.ClassHash }
func (tx DeployAccountTxnV3) GetResourceBounds() *ResourceBoundsMapping { return tx.ResourceBounds }
func (tx DeployAccountTxnV3) GetTip() U64                               { return tx.Tip }
func (tx DeployAccountTxnV3) GetPayMasterData() []*felt.Felt            { return tx.PayMasterData }
func (tx DeployAccountTxnV3) GetNonceDataMode() DataAvailabilityMode    { return tx.NonceDataMode }
func (tx DeployAccountTxnV3) GetFeeMode() DataAvailabilityMode          { return tx.FeeMode }

// DeployTxn The structure of a deploy transaction. Note that this transaction type
// is deprecated and will no longer be supported in future versions
type DeployTxn struct {
	// ClassHash The hash of the deployed contract's class
	ClassHash           *felt.Felt         `json:"class_hash"`
	Version             TransactionVersion `json:"version"`
	Type                TransactionType    `json:"type"`
	ContractAddressSalt *felt.Felt         `json:"contract_address_salt"`
	ConstructorCalldata []*felt.Felt       `json:"constructor_calldata"`
}

func (tx DeployTxn) GetClassHash() *felt.Felt             { return tx.ClassHash }
func (tx DeployTxn) GetVersion() TransactionVersion       { return tx.Version }
func (tx DeployTxn) GetType() TransactionType             { return tx.Type }
func (tx DeployTxn) GetContractAddressSalt() *felt.Felt   { return tx.ContractAddressSalt }
func (tx DeployTxn) GetConstructorCalldata() []*felt.Felt { return tx.ConstructorCalldata }

// @changed removed FunctionCall and replaced it with its fields
type L1HandlerTxn struct {
	Type               TransactionType    `json:"type"`
	Version            TransactionVersion `json:"version"`
	Nonce              string             `json:"nonce"`
	ContractAddress    *felt.Felt         `json:"contract_address"`
	EntryPointSelector *felt.Felt         `json:"entry_point_selector"`
	Calldata           []*felt.Felt       `json:"calldata"`
}

func (tx L1HandlerTxn) GetType() TransactionType          { return tx.Type }
func (tx L1HandlerTxn) GetVersion() TransactionVersion    { return tx.Version }
func (tx L1HandlerTxn) GetNonce() string                  { return tx.Nonce }
func (tx L1HandlerTxn) GetContractAddress() *felt.Felt    { return tx.ContractAddress }
func (tx L1HandlerTxn) GetEntryPointSelector() *felt.Felt { return tx.EntryPointSelector }
func (tx L1HandlerTxn) GetCalldata() []*felt.Felt         { return tx.Calldata }

type ResourceBoundsMapping struct {
	// The max amount and max price per unit of L1 gas used in this tx
	L1Gas ResourceBounds `json:"l1_gas"`
	// The max amount and max price per unit of L1 blob gas used in this tx
	L1DataGas ResourceBounds `json:"l1_data_gas"`
	// The max amount and max price per unit of L2 gas used in this tx
	L2Gas ResourceBounds `json:"l2_gas"`
}

// DA_MODE: Specifies a storage domain in Starknet. Each domain has different
// guarantees regarding availability
type DataAvailabilityMode string

const (
	DAModeL1 DataAvailabilityMode = "L1"
	DAModeL2 DataAvailabilityMode = "L2"
)

// MarshalJSON implements the json.Marshaler interface.
// It validates that the DataAvailabilityMode is either L1 or L2 before marshalling.
func (da DataAvailabilityMode) MarshalJSON() ([]byte, error) {
	switch da {
	case DAModeL1, DAModeL2:
		return json.Marshal(string(da))
	default:
		return nil, fmt.Errorf(
			"invalid DataAvailabilityMode: %s, must be either L1 or L2",
			string(da),
		)
	}
}

// @changed now is a value receiver instead of a pointer receiver
func (da DataAvailabilityMode) UInt64() (uint64, error) {
	switch da {
	case DAModeL1:
		return uint64(0), nil
	case DAModeL2:
		return uint64(1), nil
	}

	return 0, errors.New("unknown DAMode")
}

type Resource string

// Values used in the Resource Bounds hash calculation
// Ref: https://docs.starknet.io/architecture-and-concepts/network-architecture/transactions/#v3_hash_calculation
//
//nolint:lll // The link would be unclickable if we break the line.
const (
	ResourceL1Gas     Resource = "L1_GAS"
	ResourceL2Gas     Resource = "L2_GAS"
	ResourceL1DataGas Resource = "L1_DATA"
)

type ResourceBounds struct {
	// The max amount of the resource that can be used in the tx
	MaxAmount U64 `json:"max_amount"`
	// The max price per unit of this resource for this tx
	MaxPricePerUnit U128 `json:"max_price_per_unit"`
}

func (rb ResourceBounds) Bytes(resource Resource) ([]byte, error) {
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

// string must be NUM_AS_HEX
type TransactionVersion string

const (
	TransactionV0             TransactionVersion = "0x0"
	TransactionV0WithQueryBit TransactionVersion = "0x100000000000000000000000000000000"
	TransactionV1             TransactionVersion = "0x1"
	TransactionV1WithQueryBit TransactionVersion = "0x100000000000000000000000000000001"
	TransactionV2             TransactionVersion = "0x2"
	TransactionV2WithQueryBit TransactionVersion = "0x100000000000000000000000000000002"
	TransactionV3             TransactionVersion = "0x3"
	TransactionV3WithQueryBit TransactionVersion = "0x100000000000000000000000000000003"
)

// BigInt returns a big integer corresponding to the transaction version.
//
// Returns:
//   - *big.Int: a pointer to a big.Int
//   - error: an error if the conversion fails
func (v *TransactionVersion) BigInt() (*big.Int, error) {
	switch *v {
	case TransactionV0:
		return big.NewInt(0), nil
	case TransactionV1:
		return big.NewInt(1), nil
	case TransactionV2:
		return big.NewInt(2), nil
	case TransactionV3:
		return big.NewInt(3), nil
	}

	// Handle versions with query bit.
	// Remove the 0x prefix and convert to big.Int
	version, ok := new(big.Int).SetString(string(*v)[2:], 16) //nolint:mnd // hex base
	if !ok {
		return big.NewInt(-1), errors.New(fmt.Sprint("TransactionVersion %i not supported", *v))
	}

	return version, nil
}

// Int returns an integer corresponding to the transaction version.
// For versions with query bit, it returns the base version number (e.g.
// TransactionV2WithQueryBit returns 2).
// Returns -1 for invalid versions.
//
// Returns:
//   - int: the integer version, or -1 for invalid versions
func (v *TransactionVersion) Int() int {
	switch *v {
	case TransactionV0, TransactionV0WithQueryBit:
		return 0
	case TransactionV1, TransactionV1WithQueryBit:
		return 1
	case TransactionV2, TransactionV2WithQueryBit:
		return 2
	case TransactionV3, TransactionV3WithQueryBit:
		return 3
	}

	// Handle invalid versions
	return -1
}

// SubPendingTxnsInput is the optional input of the
// starknet_subscribeNewTransactionReceipts subscription.
type SubNewTxnReceiptsInput struct {
	// Optional: A vector of finality statuses to receive updates for.
	// Only `PRE_CONFIRMED` and `ACCEPTED_ON_L2` are supported. Default is
	// `ACCEPTED_ON_L2`.
	FinalityStatus []TxnFinalityStatus `json:"finality_status,omitempty"`
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
	FinalityStatus []TxnStatus `json:"finality_status,omitempty"`
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
	FinalityStatus TxnStatus `json:"finality_status"`
}

func (txn *TxnWithHashAndStatus) UnmarshalJSON(data []byte) error {
	// type alias TxnWithHashAndStatus
	var aux BlockTransaction

	err := json.Unmarshal(data, &aux)
	if err != nil {
		return err
	}

	var aux2 struct {
		FinalityStatus TxnStatus `json:"finality_status"`
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

	switch TransactionType(txnAsMap["type"].(string)) {
	case TransactionTypeDeclare:
		switch TransactionVersion(txnAsMap["version"].(string)) {
		case TransactionV0:
			return unmarshalTxnToType[DeclareTxnV0](data)
		case TransactionV1:
			return unmarshalTxnToType[DeclareTxnV1](data)
		case TransactionV2:
			return unmarshalTxnToType[DeclareTxnV2](data)
		case TransactionV3:
			return unmarshalTxnToType[DeclareTxnV3](data)
		default:
			return nil, errors.New(
				"internal error with Declare transaction version and unmarshalTxn()",
			)
		}
	case TransactionTypeDeploy:
		return unmarshalTxnToType[DeployTxn](data)
	case TransactionTypeDeployAccount:
		switch TransactionVersion(txnAsMap["version"].(string)) {
		case TransactionV1:
			return unmarshalTxnToType[DeployAccountTxnV1](data)
		case TransactionV3:
			return unmarshalTxnToType[DeployAccountTxnV3](data)
		}
	case TransactionTypeInvoke:
		switch TransactionVersion(txnAsMap["version"].(string)) {
		case TransactionV0:
			return unmarshalTxnToType[InvokeTxnV0](data)
		case TransactionV1:
			return unmarshalTxnToType[InvokeTxnV1](data)
		case TransactionV3:
			return unmarshalTxnToType[InvokeTxnV3](data)
		}
	case TransactionTypeL1Handler:
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
