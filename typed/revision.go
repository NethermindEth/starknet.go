package typed

import (
	"fmt"
	"slices"
	"strings"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/curve"
)

var (
	// There is also an array version of each type. The array is defined like this: 'type' + '*' (e.g.: "felt*", "bool*", "string*"...)
	revision_0_basic_types []string = []string{
		"felt",
		"bool",
		"string", //up to 31 ASCII characters
		"selector",
		"merkletree",
	}

	// Revision 1 includes all types from Revision 0 plus these. The only difference is that for Revision 1 "string" represents an
	// arbitrary size string instead of having a 31 ASCII characters limit in Revision 0; for this limit, use the new type "shortstring" instead.
	//
	// There is also an array version of each type. The array is defined like this: 'type' + '*' (e.g.: "ClassHash*", "timestamp*", "shortstring*"...)
	revision_1_basic_types []string = []string{
		"enum",
		"u128",
		"i128",
		"ContractAddress",
		"ClassHash",
		"timestamp",
		"shortstring",
	}

	//lint:ignore U1000 Variable used to check Preset types in other pieces of code
	revision_1_preset_types []string = []string{
		"NftId",
		"TokenAmount",
		"u256",
	}
)

var RevisionV0 revision
var RevisionV1 revision

func init() {
	preset := make(map[string]TypeDefinition)

	RevisionV0 = revision{
		version:          0,
		domain:           "StarkNetDomain",
		hashMethod:       curve.PedersenArray,
		hashMerkleMethod: curve.Pedersen,
		types: RevisionTypes{
			Basic:  revision_0_basic_types,
			Preset: preset,
		},
	}

	preset = getRevisionV1PresetTypes()
	RevisionV1 = revision{
		version:          1,
		domain:           "StarknetDomain",
		hashMethod:       curve.PoseidonArray,
		hashMerkleMethod: curve.Poseidon,
		types: RevisionTypes{
			Basic:  append(revision_1_basic_types, revision_0_basic_types...),
			Preset: preset,
		},
	}
}

type revision struct {
	//TODO: create a enum
	version          uint8
	domain           string
	hashMethod       func(felts ...*felt.Felt) *felt.Felt
	hashMerkleMethod func(a, b *felt.Felt) *felt.Felt
	types            RevisionTypes
}

type RevisionTypes struct {
	Basic  []string
	Preset map[string]TypeDefinition
}

func (rev *revision) Version() uint8 {
	return rev.version
}

func (rev *revision) Domain() string {
	return rev.domain
}

func (rev *revision) HashMethod(felts ...*felt.Felt) *felt.Felt {
	return rev.hashMethod(felts...)
}

func (rev *revision) HashMerkleMethod(a *felt.Felt, b *felt.Felt) *felt.Felt {
	var first, second *felt.Felt
	if a.Cmp(b) > 0 {
		first = b
		second = a
	} else {
		first = a
		second = b
	}
	return rev.hashMerkleMethod(first, second)
}

func (rev *revision) Types() RevisionTypes {
	return rev.types
}

func GetRevision(version uint8) (rev *revision, err error) {
	switch version {
	case 0:
		return &RevisionV0, nil
	case 1:
		return &RevisionV1, nil
	default:
		return rev, fmt.Errorf("invalid revision version")
	}
}

func getRevisionV1PresetTypes() map[string]TypeDefinition {
	//NftId
	//TokenAmount
	//u256
	presetTypes := []TypeDefinition{
		{
			Name: "NftId",
			Parameters: []TypeParameter{
				{
					Name: "collection_address",
					Type: "ContractAddress",
				},
				{
					Name: "token_id",
					Type: "u256",
				},
			},
		},
		{
			Name: "TokenAmount",
			Parameters: []TypeParameter{
				{
					Name: "token_address",
					Type: "ContractAddress",
				},
				{
					Name: "amount",
					Type: "u256",
				},
			},
		},
		{
			Name: "u256",
			Parameters: []TypeParameter{
				{
					Name: "low",
					Type: "u128",
				},
				{
					Name: "high",
					Type: "u128",
				},
			},
		},
	}

	result := make(map[string]TypeDefinition)

	for _, typeDef := range presetTypes {
		result[typeDef.Name] = typeDef
	}

	return result
}

// Check if the provided type name is a standard type defined at the SNIP 12, also validates arrays
func isStandardType(typeName string) bool {
	typeName, _ = strings.CutSuffix(typeName, "*")

	if slices.Contains(revision_0_basic_types, typeName) ||
		slices.Contains(revision_1_basic_types, typeName) ||
		slices.Contains(revision_1_preset_types, typeName) {
		return true
	}

	return false
}

// Check if the provided type name is a basic type defined at the SNIP 12, also validates arrays
func isBasicType(typeName string) bool {
	typeName, _ = strings.CutSuffix(typeName, "*")

	if slices.Contains(revision_0_basic_types, typeName) ||
		slices.Contains(revision_1_basic_types, typeName) {
		return true
	}

	return false
}

// Check if the provided type name is a preset type defined at the SNIP 12, also validates arrays
func isPresetType(typeName string) bool {
	typeName, _ = strings.CutSuffix(typeName, "*")

	return slices.Contains(revision_1_preset_types, typeName)
}
