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

/*
encoding definition for standard Starknet Domain messages
*/
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

// strToFelt converts a string containing a decimal, hexadecimal or UTF8 charset into a Felt.
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

/*
'typedData' interface for interacting and signing typed data in accordance with https://github.com/0xs34n/starknet.js/tree/develop/src/utils/typedData
*/
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

// (ref: https://github.com/starknet-io/starknet.js/blob/d7bfc37ede85448e0a55ee4efe65200ff2c45f91/src/utils/typedData.ts#L249)
func (td TypedData) GetMessageHash(account *big.Int, msg TypedMessage, sc curve.StarkCurve) (hash *big.Int, err error) {
	elements := []*big.Int{utils.UTF8StrToBig("StarkNet Message")}

	domEnc, err := td.GetTypedMessageHash("StarkNetDomain", td.Domain, sc)
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

func (td TypedData) GetTypedMessageHash(inType string, msg TypedMessage, sc curve.StarkCurve) (hash *big.Int, err error) {
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

func (td TypedData) GetTypeHash(inType string) (ret *big.Int, err error) {
	enc, err := td.EncodeType(inType)
	if err != nil {
		return ret, err
	}
	sel := utils.GetSelectorFromName(enc)
	return sel, nil
}

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
