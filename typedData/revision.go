package typedData

import (
	"fmt"

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
		//TODO: enum?
		"u128",
		"i128",
		"ContractAddress",
		"ClassHash",
		"timestamp",
		"shortstring",
	}
)

type Revision struct {
	version    uint8
	domain     string
	hashMethod func(felts ...*felt.Felt) *felt.Felt
	//TODO: hashMerkleMethod ?
	types RevisionTypes
}

type RevisionTypes struct {
	Basic  []string
	Preset map[string]TypeDefinition
}

func (rev *Revision) Version() uint8 {
	return rev.version
}

func (rev *Revision) Domain() string {
	return rev.domain
}

func (rev *Revision) HashMethod(felts ...*felt.Felt) *felt.Felt {
	return rev.hashMethod(felts...)
}

func (rev *Revision) Types() RevisionTypes {
	return rev.types
}

func NewRevision(version uint8) (rev Revision, err error) {
	preset := make(map[string]TypeDefinition)

	switch version {
	case 0:
		rev = Revision{
			version:    0,
			domain:     "StarkNetDomain",
			hashMethod: curve.PedersenArray,
			types: RevisionTypes{
				Basic:  revision_0_basic_types,
				Preset: preset,
			},
		}
		return rev, nil
	case 1:
		preset = getRevisionV1PresetTypes()
		rev = Revision{
			version:    1,
			domain:     "StarknetDomain",
			hashMethod: curve.PoseidonArray,
			types: RevisionTypes{
				Basic:  append(revision_1_basic_types, revision_0_basic_types...),
				Preset: preset,
			},
		}
		return rev, nil
	default:
		return rev, fmt.Errorf("invalid revision version")
	}
}

func getRevisionV1PresetTypes() map[string]TypeDefinition {
	//NftId
	//TokenAmount
	//U256
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
					Type: "U256",
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
					Type: "U256",
				},
			},
		},
		{
			Name: "U256",
			Parameters: []TypeParameter{
				{
					Name: "low",
					Type: "U128",
				},
				{
					Name: "high",
					Type: "U128",
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
