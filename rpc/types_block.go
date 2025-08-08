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
	// The block which is currently being built by the block proposer in height `latest` + 1.
	BlockTagPre_confirmed BlockTag = "pre_confirmed"
	// The latest Starknet block finalised by the consensus on L2.
	BlockTagLatest BlockTag = "latest"
	// The latest Starknet block which was included in a state update on L1 and finalised by the consensus on L1.
	BlockTagL1Accepted BlockTag = "l1_accepted"
)

// BlockID is a struct that is used to choose between different
// search types.
type BlockID struct {
	Number *uint64    `json:"block_number,omitempty"`
	Hash   *felt.Felt `json:"block_hash,omitempty"`
	// A tag specifying a dynamic reference to a block. Tag `l1_accepted` refers
	// to the latest Starknet block which was included in a state update on L1 and
	// finalised by the consensus on L1. Tag `latest` refers to the latest Starknet
	// block finalised by the consensus on L2. Tag `pre_confirmed` refers to the block
	// which is currently being built by the block proposer in height `latest` + 1.
	Tag BlockTag `json:",omitempty"`
}

// Block hash, number or tag, same as BLOCK_ID, but without 'pre_confirmed' or 'l1_accepted'
type SubscriptionBlockID BlockID

// WithBlockNumber sets the block number for the SubscriptionBlockID.
func (b *SubscriptionBlockID) WithBlockNumber(number uint64) SubscriptionBlockID {
	b.Number = &number

	return *b
}

// WithBlockHash sets the block hash for the SubscriptionBlockID.
func (b *SubscriptionBlockID) WithBlockHash(hash *felt.Felt) SubscriptionBlockID {
	b.Hash = hash

	return *b
}

// WithLatestTag sets the block tag to latest for the SubscriptionBlockID.
// It's the only block tag allowed for this type.
func (b *SubscriptionBlockID) WithLatestTag() SubscriptionBlockID {
	b.Tag = BlockTagLatest

	return *b
}

func (b *BlockID) UnmarshalJSON(data []byte) error {
	var tag string

	if err := json.Unmarshal(data, &tag); err == nil {
		if tag == string(BlockTagPre_confirmed) || tag == string(BlockTagLatest) || tag == string(BlockTagL1Accepted) {
			b.Tag = BlockTag(tag)

			return nil
		}
	}

	type Alias BlockID
	var aux Alias
	if err := json.Unmarshal(data, &aux); err == nil {
		*b = BlockID(aux)

		return nil
	}

	return errors.New("invalid block ID")
}

func (b *SubscriptionBlockID) UnmarshalJSON(data []byte) error {
	var aux BlockID
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.Tag == BlockTagPre_confirmed || aux.Tag == BlockTagL1Accepted {
		return fmt.Errorf("invalid block tag for this type: %s", aux.Tag)
	}

	*b = SubscriptionBlockID(aux)

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
	if b.Tag == BlockTagPre_confirmed || b.Tag == BlockTagLatest || b.Tag == BlockTagL1Accepted {
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

func (b SubscriptionBlockID) MarshalJSON() ([]byte, error) {
	if b.Tag == BlockTagPre_confirmed || b.Tag == BlockTagL1Accepted {
		return nil, fmt.Errorf("invalid block tag for this type: %s", b.Tag)
	}

	return BlockID(b).MarshalJSON()
}

// checkForPre_confirmed checks if the block ID has the 'pre_confirmed' tag. If it does, it returns an error.
// This is used to prevent the user from using the 'pre_confirmed' tag on methods that do not support it.
func checkForPre_confirmed(b BlockID) error {
	if b.Tag == BlockTagPre_confirmed {
		return errors.Join(ErrInvalidBlockID, errors.New("'pre_confirmed' tag is not supported on this method"))
	}

	return nil
}

type BlockStatus string

const (
	BlockStatus_Pre_confirmed BlockStatus = "PRE_CONFIRMED"
	BlockStatus_AcceptedOnL2  BlockStatus = "ACCEPTED_ON_L2"
	BlockStatus_AcceptedOnL1  BlockStatus = "ACCEPTED_ON_L1"
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
	case "PRE_CONFIRMED":
		*bs = BlockStatus_Pre_confirmed
	case "ACCEPTED_ON_L2":
		*bs = BlockStatus_AcceptedOnL2
	case "ACCEPTED_ON_L1":
		*bs = BlockStatus_AcceptedOnL1
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
