package typed

import (
	"bytes"
	"encoding/json"
	"fmt"
	"slices"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/utils"
)

var (
	//TODO: remove this
	// There is also an array version of each type. The array is defined like this: 'type' + '*' (e.g.: "felt*", "bool*", "string*"...)
	REVISION_0_TYPES []string = []string{
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
	REVISION_1_TYPES []string = []string{
		"u128",
		"i128",
		"ContractAddress",
		"ClassHash",
		"timestamp",
		"shortstring",
		"enum",
		"u256",
		"TokenAmount",
		"NftId",
	}
)

type TypedData struct {
	Types       map[string]TypeDefinition `json:"types"`
	PrimaryType string                    `json:"primaryType"`
	Domain      Domain                    `json:"domain"`
	Message     map[string]any            `json:"message"`
	Revision    revision                  `json:"-"`
}

type Domain struct {
	Name     string      `json:"name"`
	Version  json.Number `json:"version"`
	ChainId  json.Number `json:"chainId"`
	Revision uint8       `json:"revision,omitempty"`
}

type TypeDefinition struct {
	Name       string     `json:"-"`
	Enconding  *felt.Felt `json:"-"`
	Parameters []TypeParameter
}

type TypeParameter struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Contains string `json:"contains,omitempty"`
}

// NewTypedData initializes a new TypedData object with the given types, primary type, and domain
// for interacting and signing in accordance with https://github.com/0xs34n/starknet.js/tree/develop/src/utils/typedData
// If the primary type is invalid, it returns an error with the message "invalid primary type: {pType}".
// If there is an error encoding the type hash, it returns an error with the message "error encoding type hash: {enc.String()} {err}".
//
// Parameters:
// - types: a map[string]TypeDefinition representing the types associated with their names.
// - pType: a string representing the primary type.
// - dom: a Domain representing the domain.
// Returns:
// - td: a TypedData object
// - err: an error if any
func NewTypedData(types []TypeDefinition, primaryType string, domain Domain, message []byte) (td *TypedData, err error) {
	typesMap := make(map[string]TypeDefinition)

	for _, typeDef := range types {
		typesMap[typeDef.Name] = typeDef
	}

	if _, ok := typesMap[primaryType]; !ok {
		return td, fmt.Errorf("invalid primary type: %s", primaryType)
	}

	messageMap := make(map[string]any)
	err = json.Unmarshal(message, &messageMap)
	if err != nil {
		return td, fmt.Errorf("error unmarshalling the message: %w", err)
	}

	revision, err := GetRevision(domain.Revision)
	if err != nil {
		return td, fmt.Errorf("error getting revision: %w", err)
	}

	td = &TypedData{
		Types:       typesMap,
		PrimaryType: primaryType,
		Domain:      domain,
		Message:     messageMap,
		Revision:    revision,
	}

	for k, v := range td.Types {
		enc, err := getTypeHash(k, td.Types)
		if err != nil {
			return td, fmt.Errorf("error encoding type hash: %s %w", k, err)
		}
		v.Enconding = enc
		td.Types[k] = v
	}
	return td, nil
}

// GetMessageHash calculates the hash of a typed message for a given account using the StarkCurve.
//
// (ref: https://github.com/starknet-io/SNIPs/blob/5d5a42c654c27b377d8b7f90b453065fd19ec2eb/SNIPS/snip-12.md#specification)
//
// Parameters:
// - account: A string representing the account.
// Returns:
// - hash: A pointer to a felt.Felt representing the calculated hash.
func (td *TypedData) GetMessageHash(account string) (hash *felt.Felt, err error) {
	//signed_data = encode(PREFIX_MESSAGE, Enc[domain_separator], account, Enc[message])

	elements := []*felt.Felt{}

	//PREFIX_MESSAGE
	starknetMessage, err := utils.HexToFelt(utils.StrToHex("StarkNet Message"))
	if err != nil {
		return hash, err
	}
	elements = append(elements, starknetMessage)

	//Enc[domain_separator]
	domEnc, err := td.GetStructHash(td.Revision.Domain())
	if err != nil {
		return hash, err
	}
	elements = append(elements, domEnc)

	//account
	accountFelt, err := utils.HexToFelt(account)
	if err != nil {
		return hash, err
	}
	elements = append(elements, accountFelt)

	//Enc[message]
	msgEnc, err := td.GetStructHash(td.PrimaryType)
	if err != nil {
		return hash, err
	}
	elements = append(elements, msgEnc)

	return td.Revision.HashMethod(elements...), nil
}

// GetStructHash calculates the hash of a type and its respective data.
//
// Parameters:
// - typeName: the name of the type to be hashed.
// Returns:
// - hash: A pointer to a felt.Felt representing the calculated hash.
// - err: any error if any
func (td *TypedData) GetStructHash(typeName string, context ...string) (hash *felt.Felt, err error) {
	typeDef, ok := td.Types[typeName]
	if !ok {
		return hash, fmt.Errorf("error getting the type definition of %s", typeName)
	}
	encTypeData, err := EncodeData(&typeDef, td, context...)
	if err != nil {
		return hash, err
	}

	return td.Revision.HashMethod(append([]*felt.Felt{typeDef.Enconding}, encTypeData...)...), nil
}

func shortGetStructHash(
	typeDef *TypeDefinition,
	typedData *TypedData,
	data map[string]any,
	context ...string,
) (hash *felt.Felt, err error) {

	encTypeData, err := encodeData(typeDef, typedData, data, context...)
	if err != nil {
		return hash, err
	}

	return typedData.Revision.HashMethod(append([]*felt.Felt{typeDef.Enconding}, encTypeData...)...), nil
}

// GetTypeHash returns the hash of the given type.
//
// Parameters:
// - inType: the type to hash
// Returns:
// - ret: the hash of the given type
// - err: any error if any
func (td *TypedData) GetTypeHash(typeName string) (ret *felt.Felt, err error) {
	//TODO: create/update methods descriptions
	return getTypeHash(typeName, td.Types)
}

func getTypeHash(typeName string, types map[string]TypeDefinition) (ret *felt.Felt, err error) {
	enc, err := encodeType(typeName, types)
	if err != nil {
		return ret, err
	}
	return utils.GetSelectorFromNameFelt(enc), nil
}

// EncodeType encodes the given inType using the TypedData struct.
//
// Parameters:
// - inType: the type to encode
// Returns:
// - enc: the encoded type
// - err: any error if any
func encodeType(typeName string, types map[string]TypeDefinition) (enc string, err error) {
	customTypesEncodeResp := make(map[string]string)

	var getEncodeType func(typeName string, typeDef TypeDefinition) (result string, err error)
	getEncodeType = func(typeName string, typeDef TypeDefinition) (result string, err error) {
		var buf bytes.Buffer

		buf.WriteString(typeName)
		buf.WriteString("(")

		var ok bool

		for i, param := range typeDef.Parameters {
			buf.WriteString(fmt.Sprintf("%s:%s", param.Name, param.Type))
			if i != (len(typeDef.Parameters) - 1) {
				buf.WriteString(",")
			}
			// e.g.: "felt" or "felt*"
			if slices.Contains(REVISION_0_TYPES, param.Type) || slices.Contains(REVISION_0_TYPES, fmt.Sprintf("%s*", param.Type)) {
				continue
			} else if _, ok = customTypesEncodeResp[param.Type]; !ok {
				var customTypeDef TypeDefinition
				if customTypeDef, ok = types[param.Type]; !ok { //OBS: this is wrong on V1
					return "", fmt.Errorf("can't parse type %s from types %v", param.Type, types)
				}
				customTypesEncodeResp[param.Type], err = getEncodeType(param.Type, customTypeDef)
				if err != nil {
					return "", err
				}
			}
		}
		buf.WriteString(")")

		return buf.String(), nil
	}

	var typeDef TypeDefinition
	var ok bool
	if typeDef, ok = types[typeName]; !ok {
		return "", fmt.Errorf("can't parse type %s from types %v", typeName, types)
	}
	enc, err = getEncodeType(typeName, typeDef)
	if err != nil {
		return "", err
	}

	// appends the custom types' encode
	if len(customTypesEncodeResp) > 0 {
		// sort the types
		keys := make([]string, 0, len(customTypesEncodeResp))
		for key := range customTypesEncodeResp {
			keys = append(keys, key)
		}
		slices.Sort(keys)

		for _, key := range keys {
			enc = enc + customTypesEncodeResp[key]
		}
	}

	return enc, nil
}

func EncodeData(typeDef *TypeDefinition, td *TypedData, context ...string) (enc []*felt.Felt, err error) {
	if typeDef.Name == "StarkNetDomain" || typeDef.Name == "StarknetDomain" {
		domainMap := make(map[string]any)
		domainBytes, err := json.Marshal(td.Domain)
		if err != nil {
			return enc, err
		}
		err = json.Unmarshal(domainBytes, &domainMap)
		if err != nil {
			return enc, err
		}

		return encodeData(typeDef, td, domainMap, context...)
	}

	return encodeData(typeDef, td, td.Message, context...)
}

func encodeData(
	typeDef *TypeDefinition,
	typedData *TypedData,
	data map[string]any,
	context ...string,
) (enc []*felt.Felt, err error) {
	if len(context) != 0 {
		for _, paramName := range context {
			value, ok := data[paramName]
			if !ok {
				return enc, fmt.Errorf("context error: parameter '%s' not found in the data map", paramName)
			}
			newData, ok := value.(map[string]any)
			if !ok {
				return enc, fmt.Errorf("context error: error generating the new data map")
			}
			data = newData
		}
	}

	getStringFromData := func(key string) (resp string, err error) {
		value, ok := data[key]
		if !ok {
			return resp, fmt.Errorf("error trying to get the value of the %s type", key)
		}
		resp = fmt.Sprintf("%v", value)
		return resp, nil
	}

	getFeltFromData := func(key string) (feltValue *felt.Felt, err error) {
		strValue, err := getStringFromData(key)
		if err != nil {
			return feltValue, err
		}
		hexValue := utils.StrToHex(strValue)
		feltValue, err = utils.HexToFelt(hexValue)
		if err != nil {
			return feltValue, err
		}

		return feltValue, nil
	}

	for _, param := range typeDef.Parameters {
		switch param.Type {
		case "felt", "bool":
			value, err := getFeltFromData(param.Name)
			if err != nil {
				return enc, err
			}
			enc = append(enc, value)
		case "string":
			if typedData.Revision.version == 0 {
				value, err := getFeltFromData(param.Name)
				if err != nil {
					return enc, err
				}
				enc = append(enc, value)
			} else {
				value, err := getStringFromData(param.Name)
				if err != nil {
					return enc, err
				}
				byteArr, err := utils.StringToByteArrFelt(value)
				if err != nil {
					return enc, err
				}
				enc = append(enc, typedData.Revision.HashMethod(byteArr...))
			}
		default:
			if nextTypeDef, ok := typedData.Types[param.Type]; ok {
				structEnc, err := shortGetStructHash(&nextTypeDef, typedData, data, append(context, param.Name)...)
				if err != nil {
					return enc, err
				}
				enc = append(enc, structEnc)
				// check revision
				// }
				// if nextTypeDef, ok := typedData.rev[param.Type]; ok {
				// 	structEnc, err := shortGetStructHash(&nextTypeDef, typedData, data, append(context, param.Name)...)
				// 	if err != nil {
				// 		return enc, err
				// 	}
				// 	enc = append(enc, structEnc)
			}
		}
	}

	return enc, nil
}

func (typedData *TypedData) UnmarshalJSON(data []byte) error {
	var dec map[string]json.RawMessage
	if err := json.Unmarshal(data, &dec); err != nil {
		return err
	}

	// primaryType
	primaryType, err := utils.GetAndUnmarshalJSONFromMap[string](dec, "primaryType")
	if err != nil {
		return err
	}

	// domain
	domain, err := utils.GetAndUnmarshalJSONFromMap[Domain](dec, "domain")
	if err != nil {
		return err
	}

	// types
	rawTypes, err := utils.GetAndUnmarshalJSONFromMap[map[string]json.RawMessage](dec, "types")
	if err != nil {
		return err
	}
	var types []TypeDefinition
	for key, value := range rawTypes {
		var params []TypeParameter
		if err := json.Unmarshal(value, &params); err != nil {
			return err
		}

		typeDef := TypeDefinition{
			Name:       key,
			Parameters: params,
		}

		types = append(types, typeDef)
	}

	// message
	rawMessage, ok := dec["message"]
	if !ok {
		return fmt.Errorf("invalid typedData json: missing field 'message'")
	}
	bytesMessage, err := json.Marshal(rawMessage)
	if err != nil {
		return err
	}

	// result
	resultTypedData, err := NewTypedData(types, primaryType, domain, bytesMessage)
	if err != nil {
		return err
	}

	*typedData = *resultTypedData
	return nil
}
