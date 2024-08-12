package typed

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"
	"regexp"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/curve"
	"github.com/NethermindEth/starknet.go/utils"
)

type TypedData struct {
	Types       map[string]TypeDef
	PrimaryType string
	Domain      Domain
	Message     TypedMessage
}

type Domain struct {
	Name    string
	Version string
	ChainId string
}

type TypeDef struct {
	Encoding    *big.Int
	Definitions []Definition
}

type Definition struct {
	Name string
	Type string
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
		processStrToBig(dm.Version)
	case "chainId":
		processStrToBig(dm.ChainId)
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
// - types: a map[string]TypeDef representing the types associated with their names.
// - pType: a string representing the primary type.
// - dom: a Domain representing the domain.
// Returns:
// - td: a TypedData object
// - err: an error if any
func NewTypedData(types map[string]TypeDef, pType string, dom Domain) (td TypedData, err error) {
	td = TypedData{
		Types:       types,
		PrimaryType: pType,
		Domain:      dom,
	}
	if _, ok := td.Types[pType]; !ok {
		return td, fmt.Errorf("invalid primary type: %s", pType)
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

	for _, def := range prim.Definitions {
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
func (td TypedData) EncodeType(inType string) (enc string, err error) {
	var typeDefs TypeDef
	var ok bool
	if typeDefs, ok = td.Types[inType]; !ok {
		return enc, fmt.Errorf("can't parse type %s from types %v", inType, td.Types)
	}
	var buf bytes.Buffer
	customTypes := make(map[string]TypeDef)
	buf.WriteString(inType)
	buf.WriteString("(")
	for i, def := range typeDefs.Definitions {
		if def.Type != "felt" {
			var customTypeDef TypeDef
			if customTypeDef, ok = td.Types[def.Type]; !ok {
				return enc, fmt.Errorf("can't parse type %s from types %v", def.Type, td.Types)
			}
			customTypes[def.Type] = customTypeDef
		}
		buf.WriteString(fmt.Sprintf("%s:%s", def.Name, def.Type))
		if i != (len(typeDefs.Definitions) - 1) {
			buf.WriteString(",")
		}
	}
	buf.WriteString(")")

	for customTypeName, customType := range customTypes {
		buf.WriteString(fmt.Sprintf("%s(", customTypeName))
		for i, def := range customType.Definitions {
			buf.WriteString(fmt.Sprintf("%s:%s", def.Name, def.Type))
			if i != (len(customType.Definitions) - 1) {
				buf.WriteString(",")
			}
		}
		buf.WriteString(")")
	}
	return buf.String(), nil
}
