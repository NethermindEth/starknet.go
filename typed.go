package caigo

import (
	"bytes"
	"fmt"
	"math/big"
	"strconv"
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
	ChainId int
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
	encoding definition for standard StarkNet Domain messages
*/
func (dm Domain) FmtDefinitionEncoding(field string) (fmtEnc []*big.Int) {
	switch field {
	case "name":
		fmtEnc = append(fmtEnc, UTF8StrToBig(dm.Name))
	case "version":
		n, _ := strconv.ParseInt(dm.Version, 10, 64)
		fmtEnc = append(fmtEnc, big.NewInt(int64(n)))
	case "chainId":
		fmtEnc = append(fmtEnc, big.NewInt(int64(dm.ChainId)))
	}
	return fmtEnc
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
		return td, fmt.Errorf("invalid primary type: %v\n", pType)
	}

	for k, v := range td.Types {
		enc, err := td.GetTypeHash(k)
		if err != nil {
			return td, fmt.Errorf("error encoding type hash: %v %v\n", enc, err)
		}
		v.Encoding = enc
		td.Types[k] = v
	}
	return td, nil
}

// (ref: https://github.com/0xs34n/starknet.js/blob/767021a203ac0b9cdb282eb6d63b33bfd7614858/src/utils/typedData/index.ts#L166)
func (td TypedData) GetMessageHash(account *big.Int, msg TypedMessage, sc StarkCurve) (hash *big.Int, err error) {
	elements := []*big.Int{UTF8StrToBig("StarkNet Message")}

	domEnc, err := td.GetTypedMessageHash("StarkNetDomain", td.Domain, sc)
	if err != nil {
		return hash, fmt.Errorf("could not hash domain: %v\n", err)
	}
	elements = append(elements, domEnc)
	elements = append(elements, account)

	msgEnc, err := td.GetTypedMessageHash(td.PrimaryType, msg, sc)
	if err != nil {
		return hash, fmt.Errorf("could not hash message: %v\n", err)
	}

	elements = append(elements, msgEnc)
	elements = append(elements, big.NewInt(int64(len(elements))))
	hash, err = sc.HashElements(elements)
	return hash, err
}

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
			return hash, fmt.Errorf("error hashing internal elements: %v %v\n", innerElements, err)
		}
		elements = append(elements, innerHash)
	}

	elements = append(elements, big.NewInt(int64(len(elements))))
	hash, err = sc.HashElements(elements)
	return hash, err
}

func (td TypedData) GetTypeHash(inType string) (ret *big.Int, err error) {
	enc, err := td.EncodeType(inType)
	if err != nil {
		return ret, err
	}
	sel := GetSelectorFromName(enc)
	return sel, nil
}

func (td TypedData) EncodeType(inType string) (enc string, err error) {
	var typeDefs TypeDef
	var ok bool
	if typeDefs, ok = td.Types[inType]; !ok {
		return enc, fmt.Errorf("can't parse type %v from types %v\n", inType, td.Types)
	}
	var buf bytes.Buffer
	customTypes := make(map[string]TypeDef)
	buf.WriteString(inType)
	buf.WriteString("(")
	for i, def := range typeDefs.Definitions {
		if def.Type != "felt" {
			var customTypeDef TypeDef
			if customTypeDef, ok = td.Types[def.Type]; !ok {
				return enc, fmt.Errorf("can't parse type %v from types %v\n", def.Type, td.Types)
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
