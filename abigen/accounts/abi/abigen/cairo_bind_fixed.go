package abigen

import (
	"bytes"
	"fmt"
	"go/format"
	"strings"
	"text/template"
	
	cairoabi "github.com/NethermindEth/starknet.go/abigen/accounts/abi"
)

// BindCairoFixed generates Go bindings for Cairo contracts
func BindCairoFixed(types []string, abis []string, bytecodes []string, pkg string) (string, error) {
	
	contracts := make(map[string]*tmplCairoContract)
	structs := make(map[string]*tmplCairoStruct)
	
	for i := 0; i < len(types); i++ {
		cairoABI, err := cairoabi.JSON(strings.NewReader(abis[i]))
		if err != nil {
			return "", err
		}
		
		strippedABI := abis[i]
		
	for name, event := range cairoABI.Events {
		sanitizedName := strings.ReplaceAll(name, "::", "_")
		if sanitizedName != name {
			event.Name = sanitizedName
			cairoABI.Events[sanitizedName] = event
			delete(cairoABI.Events, name)
		}
	}
	
	binder := newCairoBinder(cairoABI)
	contract := newTmplCairoContract(types[i], strippedABI, bytecodes[i], cairoABI.Constructor, binder)
	contracts[types[i]] = contract
	}
	
	data := &tmplCairoData{
		Package:   pkg,
		Contracts: contracts,
		Structs:   structs,
	}
	
	buffer := new(bytes.Buffer)
	
	funcs := map[string]interface{}{
		"bindtype": bindCairoTypeFixed,
	}
	
	tmpl := template.Must(template.New("cairo").Funcs(funcs).Parse(tmplCairoSimple))
	if err := tmpl.Execute(buffer, data); err != nil {
		return "", fmt.Errorf("template execution error: %v", err)
	}
	
	fmt.Printf("Generated code length: %d\n", buffer.Len())
	
	code, err := format.Source(buffer.Bytes())
	if err != nil {
		return "", fmt.Errorf("%v\n%s", err, buffer)
	}
	
	return string(code), nil
}

func toCamelCase(s string) string {
	s = strings.ReplaceAll(s, "_", " ")
	s = strings.Title(s)
	return strings.ReplaceAll(s, " ", "")
}

func bindCairoTypeFixed(cairoType string) string {
	switch {
	case cairoType == "core::felt252":
		return "*felt.Felt"
	case cairoType == "core::integer::u8", cairoType == "core::integer::u16", cairoType == "core::integer::u32":
		return "uint32"
	case cairoType == "core::integer::u64", cairoType == "core::integer::u128":
		return "uint64"
	case cairoType == "core::integer::u256":
		return "*big.Int"
	case cairoType == "core::bool":
		return "bool"
	case cairoType == "core::starknet::ContractAddress":
		return "*felt.Felt"
	case strings.HasPrefix(cairoType, "core::array::Array<"):
		innerType := cairoType[len("core::array::Array<"):len(cairoType)-1]
		return "[]" + bindCairoTypeFixed(innerType)
	case strings.Contains(cairoType, "::Event"):
		return "*felt.Felt"
	default:
		return "interface{}"
	}
}
