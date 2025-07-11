package rpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
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

type Pre_confirmedBlock struct {
	Pre_confirmedBlockHeader
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
	Transaction BlockTransaction   `json:"transaction"`
	Receipt     TransactionReceipt `json:"receipt"`
}

// The dynamic block being constructed by the sequencer. Note that this object will be deprecated upon decentralisation.
type Pre_confirmedBlockWithReceipts struct {
	Pre_confirmedBlockHeader
	BlockBodyWithReceipts
}

type BlockTxHashes struct {
	BlockHeader
	Status BlockStatus `json:"status"`
	// Transactions The hashes of the transactions included in this block
	Transactions []*felt.Felt `json:"transactions"`
}

type Pre_confirmedBlockTxHashes struct {
	Pre_confirmedBlockHeader
	Transactions []*felt.Felt `json:"transactions"`
}

type BlockHeader struct {
	// Hash The hash of this block
	Hash *felt.Felt `json:"block_hash"`
	// ParentHash The hash of this block's parent
	ParentHash *felt.Felt `json:"parent_hash"`
	// Number the block number (its height)
	Number uint64 `json:"block_number"`
	// NewRoot The new global state root
	NewRoot *felt.Felt `json:"new_root"`
	// Timestamp the time in which the block was created, encoded in Unix time
	Timestamp uint64 `json:"timestamp"`
	// SequencerAddress the StarkNet identity of the sequencer submitting this block
	SequencerAddress *felt.Felt `json:"sequencer_address"`
	// The price of l1 gas in the block
	L1GasPrice ResourcePrice `json:"l1_gas_price"`
	// The price of l2 gas in the block
	L2GasPrice ResourcePrice `json:"l2_gas_price"`
	// The price of l1 data gas in the block
	L1DataGasPrice ResourcePrice `json:"l1_data_gas_price"`
	// Specifies whether the data of this block is published via blob data or calldata
	L1DAMode L1DAMode `json:"l1_da_mode"`
	// Semver of the current Starknet protocol
	StarknetVersion string `json:"starknet_version"`
}

type Pre_confirmedBlockHeader struct {
	// The block number of the block that the proposer is currently building.
	// Note that this is a local view of the node, whose accuracy depends on its polling interval length.
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
	Number uint64     `json:"block_number,omitempty"`
	Hash   *felt.Felt `json:"block_hash,omitempty"`
}

// BlockTag represents the possible values for a block tag.
type BlockTag string

const (
	// BlockTagLatest represents the latest confirmed block.
	BlockTagLatest BlockTag = "latest"
	// BlockTagPre_confirmed represents the pre_confirmed block that is yet to be confirmed.
	BlockTagPre_confirmed BlockTag = "pre_confirmed"
)

// BlockID is a struct that is used to choose between different
// search types.
type BlockID struct {
	Number *uint64    `json:"block_number,omitempty"`
	Hash   *felt.Felt `json:"block_hash,omitempty"`
	Tag    BlockTag   `json:"block_tag,omitempty"`
}

// checkForPre_confirmed checks if the block ID has the 'pre_confirmed' tag. If it does, it returns an error.
// This is used to prevent the user from using the 'pre_confirmed' tag on methods that do not support it.
func checkForPre_confirmed(b BlockID) error {
	if b.Tag == BlockTagPre_confirmed {
		return errors.Join(ErrInvalidBlockID, errors.New("'pre_confirmed' tag is not supported on this method"))
	}

	return nil
}

// MarshalJSON marshals the BlockID to JSON format.
//
// It returns a byte slice and an error. The byte slice contains the JSON representation of the BlockID,
// while the error indicates any error that occurred during the marshalling process.
//
// Parameters:
//
//	none
//
// Returns:
//   - []byte: the JSON representation of the BlockID
//   - error: any error that occurred during the marshalling process
func (b BlockID) MarshalJSON() ([]byte, error) {
	if b.Tag == BlockTagPre_confirmed || b.Tag == BlockTagLatest {
		return []byte(strconv.Quote(string(b.Tag))), nil
	}

	if b.Tag != "" {
		return nil, ErrInvalidBlockID
	}

	if b.Number != nil {
		return []byte(fmt.Sprintf(`{"block_number":%d}`, *b.Number)), nil
	}

	if b.Hash != nil && b.Hash.BigInt(big.NewInt(0)).BitLen() != 0 {
		return []byte(fmt.Sprintf(`{"block_hash":%q}`, b.Hash.String())), nil
	}

	return json.Marshal(nil)
}

type BlockStatus string

const (
	BlockStatus_Pre_confirmed BlockStatus = "PENDING"
	BlockStatus_AcceptedOnL2  BlockStatus = "ACCEPTED_ON_L2"
	BlockStatus_AcceptedOnL1  BlockStatus = "ACCEPTED_ON_L1"
	BlockStatus_Rejected      BlockStatus = "REJECTED"
)

// UnmarshalJSON unmarshals the JSON representation of a BlockStatus.
//
// It takes in a byte slice containing the JSON data to be unmarshaled.
// The function returns an error if there is an issue unmarshaling the data.
//
// Parameters:
//   - data: It takes a byte slice as a parameter, which represents the JSON data to be unmarshaled
//
// Returns:
//   - error: an error if the unmarshaling fails
func (bs *BlockStatus) UnmarshalJSON(data []byte) error {
	unquoted, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}

	switch unquoted {
	case "PENDING":
		*bs = BlockStatus_Pre_confirmed
	case "ACCEPTED_ON_L2":
		*bs = BlockStatus_AcceptedOnL2
	case "ACCEPTED_ON_L1":
		*bs = BlockStatus_AcceptedOnL1
	case "REJECTED":
		*bs = BlockStatus_Rejected
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
