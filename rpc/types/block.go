package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strconv"

	"github.com/NethermindEth/juno/core/felt"
)

var ErrInvalidBlockID = errors.New("invalid blockid")

// BlockTag represents the possible values for a block tag.
type BlockTag string

const (
	// The block which is currently being built by the block proposer in height `latest` + 1.
	BlockTagPreConfirmed BlockTag = "pre_confirmed"
	// The latest Starknet block finalised by the consensus on L2.
	BlockTagLatest BlockTag = "latest"
	// The latest Starknet block which was included in a state update on L1 and
	// finalised by the consensus on L1.
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

// UnmarshalJSON unmarshals the JSON representation of a BlockID.
func (b *BlockID) UnmarshalJSON(data []byte) error {
	var tag string

	if err := json.Unmarshal(data, &tag); err == nil {
		if tag == string(BlockTagPreConfirmed) || tag == string(BlockTagLatest) ||
			tag == string(BlockTagL1Accepted) {
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

// MarshalJSON marshals the BlockID to JSON format.
func (b BlockID) MarshalJSON() ([]byte, error) {
	if b.Tag == BlockTagPreConfirmed || b.Tag == BlockTagLatest || b.Tag == BlockTagL1Accepted {
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

// TODO: rename it to SubBlockID

// Block hash, number or tag, same as BLOCK_ID, but without 'pre_confirmed' or 'l1_accepted'
type SubscriptionBlockID BlockID

// BlockID returns a BlockID from a SubscriptionBlockID.
func (b *SubscriptionBlockID) BlockID() BlockID {
	return BlockID{
		Number: b.Number,
		Hash:   b.Hash,
		Tag:    b.Tag,
	}
}

func (b *SubscriptionBlockID) UnmarshalJSON(data []byte) error {
	var aux BlockID
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.Tag == BlockTagPreConfirmed || aux.Tag == BlockTagL1Accepted {
		return fmt.Errorf("invalid block tag for this type: %s", aux.Tag)
	}

	*b = SubscriptionBlockID(aux)

	return nil
}

func (b SubscriptionBlockID) MarshalJSON() ([]byte, error) {
	if b.Tag == BlockTagPreConfirmed || b.Tag == BlockTagL1Accepted {
		return nil, fmt.Errorf("invalid block tag for this type: %s", b.Tag)
	}

	return BlockID(b).MarshalJSON()
}

// TODO: remove methods and make tem WithSubBlockNumber, ...

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

// WithBlockNumber returns a BlockID with the given block number.
//
// Parameters:
//   - n: The block number to use for the BlockID.
//
// Returns:
//   - BlockID: A BlockID struct with the specified block number
func WithBlockNumber(n uint64) BlockID {
	var blockID BlockID
	blockID.Number = &n

	return blockID
}

// WithBlockHash returns a BlockID with the given hash.
//
// Parameters:
//   - h: The hash to use for the BlockID.
//
// Returns:
//   - BlockID: A BlockID struct with the specified hash
func WithBlockHash(h *felt.Felt) BlockID {
	var blockID BlockID
	blockID.Hash = h

	return blockID
}

// WithBlockTag creates a new BlockID with the specified tag.
//
// Parameters:
//   - tag: The tag for the BlockID
//
// Returns:
//   - BlockID: A BlockID struct with the specified tag
func WithBlockTag(tag BlockTag) BlockID {
	var blockID BlockID
	blockID.Tag = tag

	return blockID
}
