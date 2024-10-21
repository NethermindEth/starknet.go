package typedData

import (
	"encoding/json"
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
	Domain     String
	HashMethod func(felts ...*felt.Felt) *felt.Felt
	//TODO: hashMerkleMethod ?
	Types RevisionTypes
}

type RevisionTypes struct {
	Basic  []string
	Preset map[string]any
}

func NewRevision(version uint8) (rev Revision, err error) {
	preset := make(map[string]any)

	switch version {
	case 0:
		rev = Revision{
			Domain:     "StarkNetDomain",
			HashMethod: curve.PedersenArray,
			Types: RevisionTypes{
				Basic:  revision_0_basic_types,
				Preset: preset,
			},
		}
		return rev, nil
	case 1:
		preset, err = getRevisionV1PresetTypes()
		if err != nil {
			return rev, fmt.Errorf("error getting revision 1 preset types: %w", err)
		}
		rev = Revision{
			Domain:     "StarknetDomain",
			HashMethod: curve.PoseidonArray,
			Types: RevisionTypes{
				Basic:  append(revision_1_basic_types, revision_0_basic_types...),
				Preset: preset,
			},
		}
		return rev, nil
	default:
		return rev, fmt.Errorf("invalid revision version")
	}
}

func getRevisionV1PresetTypes() (result map[string]any, err error) {
	type RevV1PresetTypes struct {
		NftId       NftId
		TokenAmount TokenAmount
		U256        U256
	}

	var preset RevV1PresetTypes

	bytes, err := json.Marshal(preset)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return result, err
	}

	return result, err
}
