package types

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/dontpanicdao/caigo/felt"
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
	Encoding    felt.Felt
	Definitions []Definition
}

type Definition struct {
	Name string
	Type string
}

type TypedMessage interface {
	FmtDefinitionEncoding(string) ([]felt.Felt, error)
}

/*
encoding definition for standard StarkNet Domain messages
*/
func (dm Domain) FmtDefinitionEncoding(field string) ([]felt.Felt, error) {
	fmtEnc := []felt.Felt{}
	switch field {
	case "name":
		name, err := felt.UTF8StrToFelt(dm.Name)
		if err != nil {
			return nil, err
		}
		fmtEnc = append(fmtEnc, *name)
	case "version":
		fmtEnc = append(fmtEnc, felt.StrToFelt(dm.Version))
	case "chainId":
		fmtEnc = append(fmtEnc, felt.BigToFelt(big.NewInt(int64(dm.ChainId))))
	}
	return fmtEnc, nil
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
		v.Encoding = *enc
		td.Types[k] = v
	}
	return td, nil
}

// (ref: https://github.com/0xs34n/starknet.js/blob/767021a203ac0b9cdb282eb6d63b33bfd7614858/src/utils/typedData/index.ts#L166)
func (td TypedData) GetMessageHash(account felt.Felt, msg TypedMessage, sc felt.StarkCurve) (hash *felt.Felt, err error) {
	msgType, _ := felt.UTF8StrToFelt("StarkNet Message")
	elements := []felt.Felt{*msgType}

	domEnc, err := td.GetTypedMessageHash("StarkNetDomain", td.Domain, sc)
	if err != nil {
		return hash, fmt.Errorf("could not hash domain: %w", err)
	}
	elements = append(elements, *domEnc)
	elements = append(elements, account)
	msgEnc, err := td.GetTypedMessageHash(td.PrimaryType, msg, sc)
	if err != nil {
		return hash, fmt.Errorf("could not hash message: %w", err)
	}
	elements = append(elements, *msgEnc)
	hash, err = sc.ComputeHashOnElements(elements)
	return hash, err
}

func (td TypedData) GetTypedMessageHash(inType string, msg TypedMessage, sc felt.StarkCurve) (*felt.Felt, error) {
	prim := td.Types[inType]
	elements := []felt.Felt{prim.Encoding}

	for _, def := range prim.Definitions {
		if def.Type == "felt" {
			fmtDefinitions, err := msg.FmtDefinitionEncoding(def.Name)
			if err != nil {
				return nil, err
			}
			elements = append(elements, fmtDefinitions...)
			continue
		}

		innerElements := []felt.Felt{}
		encType := td.Types[def.Type]
		innerElements = append(innerElements, encType.Encoding)
		fmtDefinitions, err := msg.FmtDefinitionEncoding(def.Name)
		innerElements = append(innerElements, fmtDefinitions...)
		innerElements = append(innerElements, felt.BigToFelt(big.NewInt(int64(len(innerElements)))))

		innerHash, err := sc.HashElements(innerElements)
		if err != nil {
			return innerHash, fmt.Errorf("error hashing internal elements: %v %w", innerElements, err)
		}
		elements = append(elements, *innerHash)
	}
	return sc.ComputeHashOnElements(elements)
}

func (td TypedData) GetTypeHash(inType string) (ret *felt.Felt, err error) {
	enc, err := td.EncodeType(inType)
	if err != nil {
		return ret, err
	}
	sel := felt.GetSelectorFromName(enc)
	return &sel, nil
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
