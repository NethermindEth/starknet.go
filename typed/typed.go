package typed

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/utils"
)

type TypedData struct {
	types       map[string]TypeDefinition
	primaryType string
	domain      Domain
	message     map[string]any
	revision    *revision
}

func (td *TypedData) Types() map[string]TypeDefinition {
	copyMap := make(map[string]TypeDefinition, len(td.types))
	for k, v := range td.types {
		copyMap[k] = v
	}
	return copyMap
}

func (td *TypedData) PrimaryType() string {
	return td.primaryType
}

func (td *TypedData) Domain() Domain {
	return td.domain
}

func (td *TypedData) Message() map[string]any {
	copyMap := make(map[string]any, len(td.message))
	for k, v := range td.message {
		copyMap[k] = v
	}
	return copyMap
}

func (td *TypedData) Revision() revision {
	return *td.revision
}

type Domain struct {
	Name     string `json:"name"`
	Version  string `json:"version"`
	ChainId  string `json:"chainId"`
	Revision uint8  `json:"revision,omitempty"`
}

type TypeDefinition struct {
	Name               string     `json:"-"`
	Enconding          *felt.Felt `json:"-"`
	EncoddingString    string     `json:"-"`
	SingleEncString    string     `json:"-"`
	ReferencedTypesEnc []string   `json:"-"`
	Parameters         []TypeParameter
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
	//types
	typesMap := make(map[string]TypeDefinition)
	for _, typeDef := range types {
		typesMap[typeDef.Name] = typeDef
	}

	//primary type
	if _, ok := typesMap[primaryType]; !ok {
		return td, fmt.Errorf("invalid primary type: %s", primaryType)
	}

	//message
	messageMap := make(map[string]any)
	err = json.Unmarshal(message, &messageMap)
	if err != nil {
		return td, fmt.Errorf("error unmarshalling the message: %w", err)
	}

	//revision
	revision, err := GetRevision(domain.Revision)
	if err != nil {
		return td, fmt.Errorf("error getting revision: %w", err)
	}

	//domain type encoding
	domainTypeDef, err := encodeTypes(revision.Domain(), typesMap, revision)
	if err != nil {
		return td, err
	}
	typesMap[revision.Domain()] = domainTypeDef

	//types encoding
	primaryTypeDef, err := encodeTypes(primaryType, typesMap, revision)
	if err != nil {
		return td, err
	}
	typesMap[primaryType] = primaryTypeDef

	for _, typeDef := range typesMap {
		if typeDef.EncoddingString == "" {
			return td, fmt.Errorf("'encodeTypes' failed: type '%s' doesn't have encode value", typeDef.Name)
		}
	}

	td = &TypedData{
		types:       typesMap,
		primaryType: primaryType,
		domain:      domain,
		message:     messageMap,
		revision:    revision,
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
	domEnc, err := td.GetStructHash(td.revision.Domain())
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
	msgEnc, err := td.GetStructHash(td.primaryType)
	if err != nil {
		return hash, err
	}
	elements = append(elements, msgEnc)

	return td.revision.HashMethod(elements...), nil
}

// GetStructHash calculates the hash of a type and its respective data.
//
// Parameters:
// - typeName: the name of the type to be hashed.
// Returns:
// - hash: A pointer to a felt.Felt representing the calculated hash.
// - err: any error if any
func (td *TypedData) GetStructHash(typeName string, context ...string) (hash *felt.Felt, err error) {
	typeDef, ok := td.types[typeName]
	if !ok {
		if typeDef, ok = td.revision.Types().Preset[typeName]; !ok {
			return hash, fmt.Errorf("error getting the type definition of %s", typeName)
		}
	}
	encTypeData, err := EncodeData(&typeDef, td, context...)
	if err != nil {
		return hash, err
	}

	return td.revision.HashMethod(append([]*felt.Felt{typeDef.Enconding}, encTypeData...)...), nil
}

func shortGetStructHash(
	typeDef *TypeDefinition,
	typedData *TypedData,
	data map[string]any,
	isEnum bool,
	context ...string,
) (hash *felt.Felt, err error) {

	encTypeData, err := encodeData(typeDef, typedData, data, isEnum, context...)
	if err != nil {
		return hash, err
	}

	if isEnum {
		return typedData.revision.HashMethod(encTypeData...), nil
	}
	return typedData.revision.HashMethod(append([]*felt.Felt{typeDef.Enconding}, encTypeData...)...), nil
}

// GetTypeHash returns the hash of the given type.
//
// Parameters:
// - inType: the type to hash
// Returns:
// - ret: the hash of the given type
// - err: any error if any
func (td *TypedData) GetTypeHash(typeName string) (*felt.Felt, error) {
	//TODO: create/update methods descriptions
	typeDef, ok := td.types[typeName]
	if !ok {
		if typeDef, ok = td.revision.Types().Preset[typeName]; !ok {
			return typeDef.Enconding, fmt.Errorf("type '%s' not found", typeName)
		}
	}
	return typeDef.Enconding, nil
}

// EncodeType encodes the given inType using the TypedData struct.
//
// Parameters:
// - inType: the type to encode
// Returns:
// - enc: the encoded type
// - err: any error if any
func encodeTypes(typeName string, types map[string]TypeDefinition, revision *revision, isEnum ...bool) (newTypeDef TypeDefinition, err error) {
	getTypeEncodeString := func(typeName string, typeDef TypeDefinition, customTypesStringEnc *[]string, isEnum ...bool) (result string, err error) {
		verifyTypeName := func(param TypeParameter, isEnum ...bool) error {
			singleTypeName, _ := strings.CutSuffix(param.Type, "*")

			if isBasicType(singleTypeName) {
				if singleTypeName == "merkletree" {
					if param.Contains == "" {
						return fmt.Errorf("missing 'contains' value from '%s'", param.Name)
					}
					newTypeDef, err := encodeTypes(param.Contains, types, revision)
					if err != nil {
						return err
					}

					types[param.Contains] = newTypeDef
				}
				return nil
			}

			if isPresetType(singleTypeName) {
				typeEnc, ok := revision.Types().Preset[singleTypeName]
				if !ok {
					return fmt.Errorf("error trying to get the type definition of '%s'", singleTypeName)
				}
				*customTypesStringEnc = append(*customTypesStringEnc, append([]string{typeEnc.SingleEncString}, typeEnc.ReferencedTypesEnc...)...)

				return nil
			}

			if newTypeDef := types[singleTypeName]; newTypeDef.SingleEncString != "" {
				*customTypesStringEnc = append(*customTypesStringEnc, append([]string{newTypeDef.SingleEncString}, newTypeDef.ReferencedTypesEnc...)...)
				return nil
			}

			newTypeDef, err := encodeTypes(singleTypeName, types, revision, isEnum...)
			if err != nil {
				return err
			}

			*customTypesStringEnc = append(*customTypesStringEnc, append([]string{newTypeDef.SingleEncString}, newTypeDef.ReferencedTypesEnc...)...)
			types[singleTypeName] = newTypeDef

			return nil
		}

		var buf bytes.Buffer
		quotationMark := ""
		if revision.Version() == 1 {
			quotationMark = `"`
		}

		buf.WriteString(quotationMark + typeName + quotationMark)
		buf.WriteString("(")

		for i, param := range typeDef.Parameters {
			if len(isEnum) != 0 {
				reg, err := regexp.Compile(`[^\(\),\s]+`)
				if err != nil {
					return "", err
				}
				typesArr := reg.FindAllString(param.Type, -1)
				var fullTypeName string
				for i, typeNam := range typesArr {
					fullTypeName += `"` + typeNam + `"`
					if i < (len(typesArr) - 1) {
						fullTypeName += `,`
					}
				}
				buf.WriteString(fmt.Sprintf(quotationMark+"%s"+quotationMark+":"+`(`+"%s"+`)`, param.Name, fullTypeName))

				for _, typeNam := range typesArr {
					err = verifyTypeName(TypeParameter{Type: typeNam})
					if err != nil {
						return "", err
					}
				}
			} else {
				currentTypeName := param.Type

				if currentTypeName == "enum" {
					if param.Contains == "" {
						return "", fmt.Errorf("missing 'contains' value from '%s'", param.Name)
					}
					currentTypeName = param.Contains
					err = verifyTypeName(TypeParameter{Type: currentTypeName}, true)
					if err != nil {
						return "", err
					}
				}

				buf.WriteString(fmt.Sprintf(quotationMark+"%s"+quotationMark+":"+quotationMark+"%s"+quotationMark, param.Name, currentTypeName))

				err = verifyTypeName(param)
				if err != nil {
					return "", err
				}
			}
			if i != (len(typeDef.Parameters) - 1) {
				buf.WriteString(",")
			}
		}
		buf.WriteString(")")

		return buf.String(), nil
	}

	typeDef, ok := types[typeName]
	if !ok {
		return typeDef, fmt.Errorf("can't parse type %s from types %v", typeName, types)
	}

	if newTypeDef = types[typeName]; newTypeDef.EncoddingString != "" {
		return newTypeDef, nil
	}

	referencedTypesEnc := make([]string, 0)

	singleEncString, err := getTypeEncodeString(typeName, typeDef, &referencedTypesEnc, isEnum...)
	if err != nil {
		return typeDef, err
	}

	fullEncString := singleEncString
	// appends the custom types' encode
	if len(referencedTypesEnc) > 0 {
		// temp map just to remove duplicated items
		uniqueMap := make(map[string]bool)
		for _, typeEncStr := range referencedTypesEnc {
			uniqueMap[typeEncStr] = true
		}
		// clear the array
		referencedTypesEnc = make([]string, 0, len(uniqueMap))
		// fill it again
		for typeEncStr := range uniqueMap {
			referencedTypesEnc = append(referencedTypesEnc, typeEncStr)
		}

		slices.Sort(referencedTypesEnc)

		for _, typeEncStr := range referencedTypesEnc {
			fullEncString += typeEncStr
		}
	}

	newTypeDef = TypeDefinition{
		Name:               typeDef.Name,
		Parameters:         typeDef.Parameters,
		Enconding:          utils.GetSelectorFromNameFelt(fullEncString),
		EncoddingString:    fullEncString,
		SingleEncString:    singleEncString,
		ReferencedTypesEnc: referencedTypesEnc,
	}

	return newTypeDef, nil
}

func EncodeData(typeDef *TypeDefinition, td *TypedData, context ...string) (enc []*felt.Felt, err error) {
	if typeDef.Name == "StarkNetDomain" || typeDef.Name == "StarknetDomain" {
		domainMap := make(map[string]any)
		domainBytes, err := json.Marshal(td.domain)
		if err != nil {
			return enc, err
		}
		err = json.Unmarshal(domainBytes, &domainMap)
		if err != nil {
			return enc, err
		}

		// ref: https://community.starknet.io/t/signing-transactions-and-off-chain-messages/66
		domainMap["chain_id"] = domainMap["chainId"]

		return encodeData(typeDef, td, domainMap, false, context...)
	}

	return encodeData(typeDef, td, td.message, false, context...)
}

func encodeData(
	typeDef *TypeDefinition,
	typedData *TypedData,
	data map[string]any,
	isEnum bool,
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

	// helper functions
	verifyType := func(param TypeParameter, data any, isEnum bool) (resp *felt.Felt, err error) {
		//helper functions
		var handleStandardTypes func(param TypeParameter, data any, rev *revision) (resp *felt.Felt, err error)
		var handleObjectTypes func(typeDef *TypeDefinition, data any, isEnum ...bool) (resp *felt.Felt, err error)
		var handleArrays func(param TypeParameter, data any, rev *revision, isMerkle ...bool) (resp *felt.Felt, err error)

		handleStandardTypes = func(param TypeParameter, data any, rev *revision) (resp *felt.Felt, err error) {
			switch param.Type {
			case "merkletree":
				tempParam := TypeParameter{
					Name: param.Name,
					Type: param.Contains,
				}
				resp, err := handleArrays(tempParam, data, rev, true)
				if err != nil {
					return resp, err
				}
				return resp, nil
			case "enum":
				typeDef, ok := typedData.types[param.Contains]
				if !ok {
					return resp, fmt.Errorf("error trying to get the type definition of '%s' in contains of '%s'", param.Contains, param.Name)
				}
				resp, err := handleObjectTypes(&typeDef, data, true)
				if err != nil {
					return resp, err
				}
				return resp, nil
			case "NftId", "TokenAmount", "u256":
				typeDef, ok := rev.Types().Preset[param.Type]
				if !ok {
					return resp, fmt.Errorf("error trying to get the type definition of '%s'", param.Type)
				}
				resp, err := handleObjectTypes(&typeDef, data)
				if err != nil {
					return resp, err
				}
				return resp, nil
			default:
				resp, err := encodePieceOfData(param.Type, data, rev)
				if err != nil {
					return resp, err
				}
				return resp, nil
			}
		}

		handleObjectTypes = func(typeDef *TypeDefinition, data any, isEnum ...bool) (resp *felt.Felt, err error) {
			mapData, ok := data.(map[string]any)
			if !ok {
				return resp, fmt.Errorf("error trying to convert the value of '%s' to an map", typeDef)
			}

			if len(isEnum) != 0 && isEnum[0] {
				resp, err = shortGetStructHash(typeDef, typedData, mapData, true)
			} else {
				resp, err = shortGetStructHash(typeDef, typedData, mapData, false)
			}
			if err != nil {
				return resp, err
			}

			return resp, nil
		}

		handleArrays = func(param TypeParameter, data any, rev *revision, isMerkle ...bool) (resp *felt.Felt, err error) {
			var handleMerkleTree func(felts []*felt.Felt) *felt.Felt
			// ref https://github.com/starknet-io/starknet.js/blob/3cfdd8448538128bf9fd158d2e87be20310a69e3/src/utils/merkle.ts#L41
			handleMerkleTree = func(felts []*felt.Felt) *felt.Felt {
				if len(felts) == 1 {
					return felts[0]
				}
				var localArr []*felt.Felt

				for i := 0; i < len(felts); i += 2 {
					if i+1 == len(felts) {
						localArr = append(localArr, rev.HashMerkleMethod(felts[i], new(felt.Felt)))
					} else {
						localArr = append(localArr, rev.HashMerkleMethod(felts[i], felts[i+1]))
					}
				}

				return handleMerkleTree(localArr)
			}

			dataArray, ok := data.([]any)
			if !ok {
				return resp, fmt.Errorf("error trying to convert the value of '%s' to an array", param.Name)
			}
			localEncode := []*felt.Felt{}
			singleParamType, _ := strings.CutSuffix(param.Type, "*")

			if isBasicType(singleParamType) {
				for _, item := range dataArray {
					resp, err := handleStandardTypes(TypeParameter{Name: param.Name, Type: singleParamType, Contains: param.Contains}, item, rev)
					if err != nil {
						return resp, err
					}
					localEncode = append(localEncode, resp)
				}
				return rev.HashMethod(localEncode...), nil
			}

			var typeDef TypeDefinition
			if isPresetType(singleParamType) {
				typeDef, ok = rev.Types().Preset[singleParamType]
			} else {
				typeDef, ok = typedData.types[singleParamType]
			}
			if !ok {
				return resp, fmt.Errorf("error trying to get the type definition of '%s'", singleParamType)
			}

			for _, item := range dataArray {
				resp, err := handleObjectTypes(&typeDef, item, isEnum)
				if err != nil {
					return resp, err
				}
				localEncode = append(localEncode, resp)
			}

			if len(isMerkle) != 0 {
				return handleMerkleTree(localEncode), nil
			}
			return rev.HashMethod(localEncode...), nil
		}

		//function logic
		if strings.HasSuffix(param.Type, "*") {
			resp, err := handleArrays(param, data, typedData.revision)
			if err != nil {
				return resp, err
			}
			return resp, nil
		}

		if isStandardType(param.Type) {
			resp, err := handleStandardTypes(param, data, typedData.revision)
			if err != nil {
				return resp, err
			}
			return resp, nil
		}

		nextTypeDef, ok := typedData.types[param.Type]
		if !ok {
			return resp, fmt.Errorf("error trying to get the type definition of '%s'", param.Type)
		}
		resp, err = handleObjectTypes(&nextTypeDef, data, isEnum)
		if err != nil {
			return resp, err
		}
		return resp, nil
	}

	getData := func(key string) (any, error) {
		value, ok := data[key]
		if !ok {
			return value, fmt.Errorf("error trying to get the value of the '%s' param", key)
		}
		return value, nil
	}

	// function logic
	for paramIndex, param := range typeDef.Parameters {
		if isEnum {
			value, exists := data[param.Name]
			// check if it's the selected enum option
			if !exists {
				if paramIndex == len(typeDef.Parameters)-1 {
					return enc, fmt.Errorf("no enum option selected for '%s', the data is not valid", typeDef.Name)
				}
				continue
			}

			dataArr, ok := value.([]any)
			if !ok {
				return enc, fmt.Errorf("error trying to convert the data value of '%s' to an array", param.Name)
			}

			enc = append(enc, new(felt.Felt).SetUint64(uint64(paramIndex)))

			if len(dataArr) == 0 {
				enc = append(enc, &felt.Zero)
				break
			}

			reg := regexp.MustCompile(`[^\(\),\s]+`)
			typesArr := reg.FindAllString(param.Type, -1)

			for i, typeNam := range typesArr {
				resp, err := verifyType(TypeParameter{Type: typeNam}, dataArr[i], false)
				if err != nil {
					return enc, err
				}
				enc = append(enc, resp)
			}

			break
		}

		localData, err := getData(param.Name)
		if err != nil {
			return enc, err
		}

		resp, err := verifyType(param, localData, false)
		if err != nil {
			return enc, err
		}
		enc = append(enc, resp)
	}

	return enc, nil
}

func encodePieceOfData(typeName string, data any, rev *revision) (resp *felt.Felt, err error) {
	getFeltFromData := func() (feltValue *felt.Felt, err error) {
		strValue := func(data any) string {
			switch v := data.(type) {
			case string:
				return v
			case float64:
				// Handle floating point numbers without trailing zeros
				if float64(int64(v)) == v {
					return strconv.FormatInt(int64(v), 10)
				}
				return strconv.FormatFloat(v, 'f', -1, 64)
			case float32:
				if float32(int32(v)) == v {
					return strconv.FormatInt(int64(v), 10)
				}
				return strconv.FormatFloat(float64(v), 'f', -1, 32)
			case int:
				return strconv.Itoa(v)
			case int64:
				return strconv.FormatInt(v, 10)
			case int32:
				return strconv.FormatInt(int64(v), 10)
			case bool:
				return strconv.FormatBool(v)
			case nil:
				return ""
			default:
				return fmt.Sprintf("%v", v)
			}
		}(data)
		hexValue := utils.StrToHex(strValue)
		feltValue, err = utils.HexToFelt(hexValue)
		if err != nil {
			return feltValue, err
		}

		return feltValue, nil
	}

	switch typeName {
	case "felt", "shortstring", "u128", "ContractAddress", "ClassHash", "timestamp":
		resp, err = getFeltFromData()
		if err != nil {
			return resp, err
		}
		return resp, nil
	case "bool":
		boolVal, ok := data.(bool)
		if !ok {
			return resp, fmt.Errorf("faild to convert '%v' to 'bool'", data)
		}
		if boolVal {
			return new(felt.Felt).SetUint64(1), nil
		}
		return new(felt.Felt).SetUint64(0), nil
	case "i128":
		strValue := fmt.Sprintf("%v", data)
		bigNum, ok := new(big.Int).SetString(strValue, 0)
		if !ok {
			return resp, fmt.Errorf("faild to convert '%s' of type 'i128' to big.Int", strValue)
		}
		feltValue := new(felt.Felt).SetBigInt(bigNum)
		return feltValue, nil
	case "string":
		if rev.Version() == 0 {
			resp, err := getFeltFromData()
			if err != nil {
				return resp, err
			}
			return resp, nil
		} else {
			value := fmt.Sprintf("%v", data)
			byteArr, err := utils.StringToByteArrFelt(value)
			if err != nil {
				return resp, err
			}
			return rev.HashMethod(byteArr...), nil
		}
	case "selector":
		value := fmt.Sprintf("%v", data)
		return utils.GetSelectorFromNameFelt(value), nil
	default:
		return resp, fmt.Errorf("invalid type '%s'", typeName)
		// check revision
		// }
		// if nextTypeDef, ok := typedData.rev[param.Type]; ok {
		// 	structEnc, err := shortGetStructHash(&nextTypeDef, typedData, data, append(context, param.Name)...)
		// 	if err != nil {
		// 		return felt, err
		// 	}
		// 	enc = append(enc, structEnc)
	}
}

func (td *TypedData) UnmarshalJSON(data []byte) error {
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

	*td = *resultTypedData
	return nil
}

func (domain *Domain) UnmarshalJSON(data []byte) error {
	var dec map[string]any
	if err := json.Unmarshal(data, &dec); err != nil {
		return err
	}

	getField := func(fieldName string) (string, error) {
		value, ok := dec[fieldName]
		if !ok {
			return "", fmt.Errorf("error getting value of '%s' from 'domain' struct", fieldName)
		}
		return fmt.Sprintf("%v", value), nil
	}

	name, err := getField("name")
	if err != nil {
		return err
	}

	version, err := getField("version")
	if err != nil {
		return err
	}

	chainId, err := getField("chainId")
	if err != nil {
		var err2 error
		// ref: https://community.starknet.io/t/signing-transactions-and-off-chain-messages/66
		chainId, err2 = getField("chain_id")
		if err2 != nil {
			return err
		}
	}

	revision, err := getField("revision")
	if err != nil {
		revision = "0"
	}
	numRevision, err := strconv.ParseUint(revision, 10, 8)
	if err != nil {
		return err
	}

	*domain = Domain{
		Name:     name,
		Version:  version,
		ChainId:  chainId,
		Revision: uint8(numRevision),
	}
	return nil
}
