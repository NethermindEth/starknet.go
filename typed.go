package starknetgo

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"
	"regexp"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/types"
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

// FmtDefinitionEncoding formats the given field value(s) into a slice of big integers. 
// This is the encoding definition for standard Starknet Domain messages.
//
// Parameter:
// - field: the field to be formatted (name, version, or chainId).
//
// Return:
// - fmtEnc: a slice of big integers containing the formatted field values.
func (dm Domain) FmtDefinitionEncoding(field string) (fmtEnc []*big.Int) {
	processStrToBig := func(fieldVal string) {
		felt := strToFelt(fieldVal)
		bigInt, ok := feltToBig(felt)
		if ok {
			fmtEnc = append(fmtEnc, bigInt)
		}
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

// strToFelt converts a string containing a decimal, hexadecimal or UTF8 charset into a Felt (*felt.Felt).
//
// The function takes a string as input and attempts to convert it to a big integer using the
// `big.Int.SetString` method. If the conversion is successful, the resulting big integer is
// set as the value of the *felt.Felt pointer `f` and returned. If the conversion fails, the
// function moves on to the next conversion method.
//
// The function also checks if the input string matches the regular expression `^([[:graph:]]|[[:space:]]){1,31}$`.
// If the string matches the regular expression, it is first encoded as a hex string using `hex.EncodeToString`,
// and then converted to a big integer using `big.Int.SetString` with base 16. If the conversion is successful,
// the resulting big integer is set as the value of the *felt.Felt pointer `f` and returned. If the conversion fails,
// the function returns the initial value of `f`, which is `&felt.Zero`.
//
// The function always returns a *felt.Felt pointer.
func strToFelt(str string) *felt.Felt {
	var f = &felt.Zero
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

// feltToBig converts a *felt.Felt to a *big.Int.
//
// It takes a pointer to a *felt.Felt as a parameter and returns a *big.Int and a boolean value.
func feltToBig(feltNum *felt.Felt) (*big.Int, bool) {
	return new(big.Int).SetString(feltNum.String(), 0)

}

// NewTypedData initializes a new TypedData struct for interacting and signing typed data in accordance with https://github.com/0xs34n/starknet.js/tree/develop/src/utils/typedData
//
// It takes in a map of types, a primary type string, and a domain.
// It returns a TypedData struct and an error.
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

// GetMessageHash calculates the hash of a TypedMessage for a given account using the Stark elliptic curve.
// (ref: https://github.com/0xs34n/starknet.js/blob/767021a203ac0b9cdb282eb6d63b33bfd7614858/src/utils/typedData/index.ts#L166)
//
// Parameters:
// - account: The account for which the hash is calculated.
// - msg: The TypedMessage to hash.
// - sc: The StarkCurve used for hashing.
//
// Returns:
// - hash: The calculated hash.
// - err: An error if the hash calculation fails.
func (td TypedData) GetMessageHash(account *big.Int, msg TypedMessage, sc StarkCurve) (hash *big.Int, err error) {
	elements := []*big.Int{types.UTF8StrToBig("Starknet Message")}

	domEnc, err := td.GetTypedMessageHash("StarknetDomain", td.Domain, sc)
	if err != nil {
		return hash, fmt.Errorf("could not hash domain: %w", err)
	}
	elements = append(elements, domEnc)
	elements = append(elements, account)

	msgEnc, err := td.GetTypedMessageHash(td.PrimaryType, msg, sc)
	if err != nil {
		return hash, fmt.Errorf("could not hash message: %w", err)
	}

	elements = append(elements, msgEnc)
	hash, err = sc.ComputeHashOnElements(elements)
	return hash, err
}

// GetTypedMessageHash calculates the hash of a typed message using a Stark elliptic curve.
//
// Parameters:
// - inType: the type of the input message.
// - msg: the typed message to calculate the hash for.
// - sc: the StarkCurve to use for the hash calculation.
//
// Returns:
// - hash: the calculated hash as a big.Int pointer.
// - err: an error if the hash calculation fails.
func (td TypedData) GetTypedMessageHash(inType string, msg TypedMessage, sc StarkCurve) (hash *big.Int, err error) {
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

		innerHash, err := sc.HashElements(innerElements)
		if err != nil {
			return hash, fmt.Errorf("error hashing internal elements: %v %w", innerElements, err)
		}
		elements = append(elements, innerHash)
	}

	hash, err = sc.ComputeHashOnElements(elements)
	return hash, err
}

// GetTypeHash returns the type hash of a given input type.
//
// It takes in a string representing the input type and returns the corresponding
// type hash as a big.Int pointer and an error if any.
func (td TypedData) GetTypeHash(inType string) (ret *big.Int, err error) {
	enc, err := td.EncodeType(inType)
	if err != nil {
		return ret, err
	}
	sel := types.GetSelectorFromName(enc)
	return sel, nil
}

// EncodeType encodes the given input type to a string representation.
//
// Parameter(s):
// - inType: the input type to encode.
//
// Return type(s):
// - enc: the encoded string representation of the input type.
// - err: an error if the encoding fails.
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
