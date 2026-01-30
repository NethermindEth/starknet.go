package rpcv10

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/NethermindEth/juno/core/felt"
)

var ErrInvalidBlockID = errors.New("invalid blockid")

type Block struct {
	BlockHeader
	Status BlockStatus `json:"status"`
	// Transactions The transactions in this block
	Transactions []BlockTransaction `json:"transactions"`
}

type PreConfirmedBlock struct {
	PreConfirmedBlockHeader
	Transactions []BlockTransaction `json:"transactions"`
}

// encoding/json doesn't support inlining fields
type BlockWithReceipts struct {
	BlockHeader
	Status BlockStatus `json:"status"`
	BlockBodyWithReceipts
}

type BlockBodyWithReceipts struct {
	Transactions []TransactionWithReceipt `json:"transactions"`
}

type TransactionWithReceipt struct {
	Transaction Transaction        `json:"transaction"`
	Receipt     TransactionReceipt `json:"receipt"`
}

// UnmarshalJSON unmarshals the JSON representation of a TransactionWithReceipt.
func (twr *TransactionWithReceipt) UnmarshalJSON(data []byte) error {
	type Temp struct {
		Transaction json.RawMessage    `json:"transaction"`
		Receipt     TransactionReceipt `json:"receipt"`
	}
	var temp Temp
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	txn, err := unmarshalTxn(temp.Transaction)
	if err != nil {
		return err
	}
	twr.Transaction = txn
	twr.Receipt = temp.Receipt

	return nil
}

// The dynamic block being constructed by the sequencer. Note that this object
// will be deprecated upon decentralisation.
type PreConfirmedBlockWithReceipts struct {
	PreConfirmedBlockHeader
	BlockBodyWithReceipts
}

type BlockTxHashes struct {
	BlockHeader
	Status BlockStatus `json:"status"`
	// Transactions The hashes of the transactions included in this block
	Transactions []*felt.Felt `json:"transactions"`
}

type PreConfirmedBlockTxHashes struct {
	PreConfirmedBlockHeader
	Transactions []*felt.Felt `json:"transactions"`
}

type BlockHeader struct {
	// The root of Merkle Patricia trie for events in the block
	EventCommitment *felt.Felt `json:"event_commitment"`
	// The number of events in the block
	EventCount uint64 `json:"event_count"`
	// Hash The hash of this block
	Hash *felt.Felt `json:"block_hash"`
	// Specifies whether the data of this block is published via blob data or calldata
	L1DAMode L1DAMode `json:"l1_da_mode"`
	// The price of l1 data gas in the block
	L1DataGasPrice ResourcePrice `json:"l1_data_gas_price"`
	// The price of l1 gas in the block
	L1GasPrice ResourcePrice `json:"l1_gas_price"`
	// The price of l2 gas in the block
	L2GasPrice ResourcePrice `json:"l2_gas_price"`
	// NewRoot The new global state root
	NewRoot *felt.Felt `json:"new_root"`
	// Number the block number (its height)
	Number uint64 `json:"block_number"`
	// ParentHash The hash of this block's parent
	ParentHash *felt.Felt `json:"parent_hash"`
	// The root of Merkle Patricia trie for receipts in the block
	ReceiptCommitment *felt.Felt `json:"receipt_commitment"`
	// SequencerAddress the StarkNet identity of the sequencer submitting this block
	SequencerAddress *felt.Felt `json:"sequencer_address"`
	// Semver of the current Starknet protocol
	StarknetVersion string `json:"starknet_version"`
	// The state diff commitment hash in the block
	StateDiffCommitment *felt.Felt `json:"state_diff_commitment"`
	// The length of the state diff in the block
	StateDiffLength uint64 `json:"state_diff_length"`
	// Timestamp the time in which the block was created, encoded in Unix time
	Timestamp uint64 `json:"timestamp"`
	// The root of Merkle Patricia trie for transactions in the block
	TransactionCommitment *felt.Felt `json:"transaction_commitment"`
	// The number of transactions in the block
	TransactionCount uint64 `json:"transaction_count"`
}

type PreConfirmedBlockHeader struct {
	// The block number of the block that the proposer is currently building.
	// Note that this is a local view of the node, whose accuracy depends on its
	// polling interval length.
	Number uint64 `json:"block_number"`
	// Timestamp the time in which the block was created, encoded in Unix time
	Timestamp uint64 `json:"timestamp"`
	// SequencerAddress the StarkNet identity of the sequencer submitting this block
	SequencerAddress *felt.Felt `json:"sequencer_address"`
	// The price of l1 gas in the block
	L1GasPrice ResourcePrice `json:"l1_gas_price"`
	// The price of l2 gas in the block
	L2GasPrice ResourcePrice `json:"l2_gas_price"`
	// Semver of the current Starknet protocol
	StarknetVersion string `json:"starknet_version"`
	// The price of l1 data gas in the block
	L1DataGasPrice ResourcePrice `json:"l1_data_gas_price"`
	// Specifies whether the data of this block is published via blob data or calldata
	L1DAMode L1DAMode `json:"l1_da_mode"`
}

// BlockHashAndNumberOutput is a struct that is returned by BlockHashAndNumber.
type BlockHashAndNumberOutput struct {
	Number uint64     `json:"block_number"`
	Hash   *felt.Felt `json:"block_hash"`
}

type BlockStatus string

const (
	BlockStatusPreConfirmed BlockStatus = "PRE_CONFIRMED"
	BlockStatusAcceptedOnL2 BlockStatus = "ACCEPTED_ON_L2"
	BlockStatusAcceptedOnL1 BlockStatus = "ACCEPTED_ON_L1"
)

// UnmarshalJSON unmarshals the JSON representation of a BlockStatus.
//
// It takes in a byte slice containing the JSON data to be unmarshaled.
// The function returns an error if there is an issue unmarshaling the data.
//
// Parameters:
//   - data: It takes a byte slice as a parameter, which represents the JSON data to
//     be unmarshaled
//
// Returns:
//   - error: an error if the unmarshaling fails
func (bs *BlockStatus) UnmarshalJSON(data []byte) error {
	unquoted, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}

	switch unquoted {
	case "PRE_CONFIRMED":
		*bs = BlockStatusPreConfirmed
	case "ACCEPTED_ON_L2":
		*bs = BlockStatusAcceptedOnL2
	case "ACCEPTED_ON_L1":
		*bs = BlockStatusAcceptedOnL1
	default:
		return fmt.Errorf("unsupported status: %s", data)
	}

	return nil
}

// MarshalJSON returns the JSON encoding of BlockStatus.
//
// Parameters:
//
//	none
//
// Returns:
//   - []byte: a byte slice
//   - error: an error if any
func (bs BlockStatus) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(string(bs))), nil
}

type L1DAMode int

const (
	L1DAModeBlob L1DAMode = iota
	L1DAModeCalldata
)

func (mode L1DAMode) String() string {
	switch mode {
	case L1DAModeBlob:
		return "BLOB"
	case L1DAModeCalldata:
		return "CALLDATA"
	default:
		return "Unknown L1DAMode"
	}
}

func (mode *L1DAMode) UnmarshalJSON(b []byte) error {
	str := strings.Trim(string(b), "\"")
	switch str {
	case "BLOB":
		*mode = L1DAModeBlob
	case "CALLDATA":
		*mode = L1DAModeCalldata
	default:
		return fmt.Errorf("unknown L1DAMode: %s", str)
	}

	return nil
}

func (mode L1DAMode) MarshalJSON() ([]byte, error) {
	return json.Marshal(mode.String())
}

type ResourcePrice struct {
	// the price of one unit of the given resource, denominated in fri (10^-18 strk)
	PriceInFRI *felt.Felt `json:"price_in_fri,omitempty"`
	// The price of one unit of the given resource, denominated in wei
	PriceInWei *felt.Felt `json:"price_in_wei"`
}
