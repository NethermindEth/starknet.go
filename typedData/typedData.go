package typedData

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"regexp"
	"slices"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/curve"
	"github.com/NethermindEth/starknet.go/utils"
)

var (
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
	Types       map[string]TypeDefinition
	PrimaryType string
	Domain      Domain
	Message     map[string]any
}

type Domain struct {
	Name     string
	Version  json.Number
	ChainId  json.Number
	Revision uint8 `json:"contains,omitempty"`
}

type TypeDefinition struct {
	Name       string `json:"-"`
	Encoding   *big.Int
	Parameters []TypeParameter
}

type TypeParameter struct {
	Name     string
	Type     string
	Contains string `json:"contains,omitempty"`
}

type TypedMessage interface {
	FmtDefinitionEncoding(string) []*big.Int
}

// FmtDefinitionEncoding formats the definition (standard Starknet Domain) encoding.
//
// Parameters:
// - field: the field to format the encoding for
// Returns:
// - fmtEnc: a slice of big integers
func (dm Domain) FmtDefinitionEncoding(field string) (fmtEnc []*big.Int) {
	processStrToBig := func(fieldVal string) {
		feltVal := strToFelt(fieldVal)
		bigInt := utils.FeltToBigInt(feltVal)
		fmtEnc = append(fmtEnc, bigInt)
	}

	switch field {
	case "name":
		processStrToBig(dm.Name)
	case "version":
		processStrToBig(dm.Version.String())
	case "chainId":
		processStrToBig(dm.ChainId.String())
	}
	return fmtEnc
}

// strToFelt converts a string (decimal, hexadecimal or UTF8 charset) to a *felt.Felt.
//
// Parameters:
// - str: the string to convert to a *felt.Felt
// Returns:
// - *felt.Felt: a *felt.Felt with the value of str
func strToFelt(str string) *felt.Felt {
	var f = new(felt.Felt)
	asciiRegexp := regexp.MustCompile(`^([[:graph:]]|[[:space:]]){1,31}$`)

	if b, ok := new(big.Int).SetString(str, 0); ok {
		f.SetBytes(b.Bytes())
		return f
	}
	// TODO: revisit conversation on seperate 'ShortString' conversion
	if asciiRegexp.MatchString(str) {
		hexStr := hex.EncodeToString([]byte(str))
		if b, ok := new(big.Int).SetString(hexStr, 16); ok {
			f.SetBytes(b.Bytes())
			return f
		}
	}

	return f
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
func NewTypedData(types []TypeDefinition, primaryType string, domain Domain, message []byte) (td TypedData, err error) {
	typesMap := make(map[string]TypeDefinition)

	for _, typeDef := range types {
		typesMap[typeDef.Name] = typeDef
	}

	messageMap := make(map[string]any)
	err = json.Unmarshal(message, &messageMap)
	if err != nil {
		return td, fmt.Errorf("error unmarshalling the message: %w", err)
	}

	td = TypedData{
		Types:       typesMap,
		PrimaryType: primaryType,
		Domain:      domain,
		Message:     messageMap,
	}
	if _, ok := td.Types[primaryType]; !ok {
		return td, fmt.Errorf("invalid primary type: %s", primaryType)
	}

	for k, v := range td.Types {
		enc, err := td.GetTypeHash(k)
		if err != nil {
			return td, fmt.Errorf("error encoding type hash: %s %w", enc.String(), err)
		}
		v.Encoding = enc
		td.Types[k] = v
	}
	return td, nil
}

// GetMessageHash calculates the hash of a typed message for a given account using the StarkCurve.
// (ref: https://github.com/0xs34n/starknet.js/blob/767021a203ac0b9cdb282eb6d63b33bfd7614858/src/utils/typedData/index.ts#L166)
//
// Parameters:
// - account: A pointer to a big.Int representing the account.
// - msg: A TypedMessage object representing the message.
// Returns:
// - hash: A pointer to a big.Int representing the calculated hash.
func (td TypedData) GetMessageHash(account *big.Int, msg TypedMessage) (hash *big.Int) {
	elements := []*big.Int{utils.UTF8StrToBig("StarkNet Message")}

	domEnc := td.GetTypedMessageHash("StarkNetDomain", td.Domain)

	elements = append(elements, domEnc)
	elements = append(elements, account)

	msgEnc := td.GetTypedMessageHash(td.PrimaryType, msg)

	elements = append(elements, msgEnc)

	return curve.ComputeHashOnElements(elements)
}

// GetTypedMessageHash calculates the hash of a typed message using the provided StarkCurve.
//
// Parameters:
//   - inType: the type of the message
//   - msg: the typed message
//
// Returns:
//   - hash: the calculated hash
func (td TypedData) GetTypedMessageHash(inType string, msg TypedMessage) (hash *big.Int) {
	prim := td.Types[inType]
	elements := []*big.Int{prim.Encoding}

	for _, def := range prim.Parameters {
		if def.Type == "felt" {
			fmtDefinitions := msg.FmtDefinitionEncoding(def.Name)
			elements = append(elements, fmtDefinitions...)
			continue
		}

		innerElements := []*big.Int{}
		encType := td.Types[def.Type]
		innerElements = append(innerElements, encType.Encoding)
		fmtDefinitions := msg.FmtDefinitionEncoding(def.Name)
		innerElements = append(innerElements, fmtDefinitions...)
		innerElements = append(innerElements, big.NewInt(int64(len(innerElements))))

		innerHash := curve.HashPedersenElements(innerElements)
		elements = append(elements, innerHash)
	}

	return curve.ComputeHashOnElements(elements)
}

// GetTypeHash returns the hash of the given type.
//
// Parameters:
// - inType: the type to hash
// Returns:
// - ret: the hash of the given type
// - err: any error if any
func (td TypedData) GetTypeHash(inType string) (ret *big.Int, err error) {
	enc, err := td.EncodeType(inType)
	if err != nil {
		return ret, err
	}
	return utils.GetSelectorFromName(enc), nil
}

// EncodeType encodes the given inType using the TypedData struct.
//
// Parameters:
// - inType: the type to encode
// Returns:
// - enc: the encoded type
// - err: any error if any
func (td TypedData) EncodeType(typeName string) (enc string, err error) {
	customTypesEncodeResp := make(map[string]string)

	var encodeType func(typeName string, typeDef TypeDefinition) (result string, err error)
	encodeType = func(typeName string, typeDef TypeDefinition) (result string, err error) {
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
				if customTypeDef, ok = td.Types[param.Type]; !ok { //OBS: this is wrong on V1
					return "", fmt.Errorf("can't parse type %s from types %v", param.Type, td.Types)
				}
				customTypesEncodeResp[param.Type], err = encodeType(param.Type, customTypeDef)
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
	if typeDef, ok = td.Types[typeName]; !ok {
		return "", fmt.Errorf("can't parse type %s from types %v", typeName, td.Types)
	}
	enc, err = encodeType(typeName, typeDef)
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

func (typedData *TypedData) UnmarshalJSON(data []byte) error {
	var dec map[string]interface{}
	if err := json.Unmarshal(data, &dec); err != nil {
		return err
	}

	// primaryType
	rawPrimaryType, ok := dec["primaryType"]
	if !ok {
		return fmt.Errorf("invalid typedData json: missing field 'primaryType'")
	}
	primaryType, ok := rawPrimaryType.(string)
	if !ok {
		return fmt.Errorf("failed to unmarshal 'primaryType', it's not a string")
	}

	// domain
	rawDomain, ok := dec["domain"]
	if !ok {
		return fmt.Errorf("invalid typedData json: missing field 'domain'")
	}
	bytesDomain, err := json.Marshal(rawDomain)
	if err != nil {
		return err
	}
	var domain Domain
	if err := json.Unmarshal(bytesDomain, &domain); err != nil {
		return err
	}

	// types
	rawTypes, err := utils.UnwrapJSON(dec, "types")
	if err != nil {
		return err
	}
	var types []TypeDefinition
	for key, value := range rawTypes {
		bytesValue, err := json.Marshal(value)
		if err != nil {
			return err
		}

		var params []TypeParameter
		if err := json.Unmarshal(bytesValue, &params); err != nil {
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

	*typedData = resultTypedData
	return nil

	// TODO: implement typedMessage unmarshal
}
