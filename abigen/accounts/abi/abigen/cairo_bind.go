package abigen

import (
	"bytes"
	"fmt"
	"go/format"
	"regexp"
	"strings"
	"text/template"
	"unicode"

	cairoabi "github.com/NethermindEth/starknet.go/abigen/accounts/abi"
)

func BindCairo(types []string, abis []string, bytecodes []string, pkg string) (string, error) {
	contracts := make(map[string]*tmplCairoContract)
	
	structs := make(map[string]*tmplCairoStruct)

	for i := 0; i < len(types); i++ {
		cairoABI, err := cairoabi.JSON(strings.NewReader(abis[i]))
		if err != nil {
			return "", err
		}
		
		strippedABI := strings.Map(func(r rune) rune {
			if unicode.IsSpace(r) {
				return -1
			}
			return r
		}, abis[i])

		for name, event := range cairoABI.Events {
			sanitizedName := strings.ReplaceAll(name, "::", "_")
			if sanitizedName != name {
				event.Name = sanitizedName
				cairoABI.Events[sanitizedName] = event
				delete(cairoABI.Events, name)
			}
		}

		binder := newCairoBinder(cairoABI)
		methods := make(map[string]*tmplCairoMethod)
		events := make(map[string]*tmplCairoEvent)
		
		for name, method := range binder.methods {
			methodCopy := method
			methods[name] = &methodCopy
		}
		
		for name, event := range binder.events {
			eventCopy := event
			events[name] = &eventCopy
		}
		
		contract := &tmplCairoContract{
			Type:        ToCamelCase(types[i]),
			InputABI:    strippedABI,
			InputBin:    bytecodes[i],
			Constructor: cairoABI.Constructor,
			Methods:     methods,
			Events:      events,
		}
		contracts[types[i]] = contract
	}

	data := &tmplCairoData{
		Package:   pkg,
		Contracts: contracts,
		Structs:   structs,
	}
	buffer := new(bytes.Buffer)

	funcs := map[string]interface{}{
		"bindtype": bindCairoType,
	}
	fmt.Printf("Template content length: %d\n", len(tmplCairoSource))
	if len(tmplCairoSource) == 0 {
		return "", fmt.Errorf("template content is empty")
	}
	tmpl := template.Must(template.New("").Funcs(funcs).Parse(tmplCairoSource))
	if err := tmpl.Execute(buffer, data); err != nil {
		return "", fmt.Errorf("template execution error: %v", err)
	}
	
	// fmt.Printf("Generated code: %s\n", buffer.String())
	
	code, err := format.Source(buffer.Bytes())
	if err != nil {
		return "", fmt.Errorf("%v\n%s", err, buffer)
	}
	return string(code), nil
}

func bindCairoType(cairoType string) string {
	if strings.Contains(cairoType, "Array<") {
		arrayRegex := regexp.MustCompile(`Array<(.+)>`)
		if matches := arrayRegex.FindStringSubmatch(cairoType); len(matches) > 0 {
			elementType := matches[1]
			return "[]" + bindCairoType(elementType)
		}
	}

	parts := strings.Split(cairoType, "::")
	baseType := parts[len(parts)-1]

	switch baseType {
	case "felt252":
		return "*felt.Felt"
	case "u8", "u16", "u32":
		return "uint32"
	case "u64":
		return "uint64"
	case "u128":
		return "uint64" // Go doesn't have uint128, use uint64 or big.Int
	case "u256":
		return "*big.Int"
	case "bool":
		return "bool"
	case "ContractAddress":
		return "*felt.Felt"
	}
	
	arrayRegex := regexp.MustCompile(`^(.+)\[(\d*)\]$`)
	if matches := arrayRegex.FindStringSubmatch(baseType); len(matches) > 0 {
		elementType := matches[1]
		if matches[2] == "" {
			return "[]" + bindCairoType(elementType)
		} else {
			return fmt.Sprintf("[%s]%s", matches[2], bindCairoType(elementType))
		}
	}

	return "*felt.Felt"
}

func IsCairoStruct(t string) bool {
	if strings.HasPrefix(t, "Array<") || strings.Contains(t, "[") {
		return true
	}
	
	switch t {
	case "felt252", "u8", "u16", "u32", "u64", "u128", "u256", "bool", "ContractAddress":
		return false
	default:
		if strings.Contains(t, "::") {
			parts := strings.Split(t, "::")
			baseType := parts[len(parts)-1]
			switch baseType {
			case "felt252", "u8", "u16", "u32", "u64", "u128", "u256", "bool", "ContractAddress":
				return false
			default:
				return true
			}
		}
		return true
	}
}

func BindCairoStructType(typeName string, name string, members []cairoabi.Argument, structs map[string]*tmplCairoStruct) string {
	id := typeName
	if s, exist := structs[id]; exist {
		return s.Name
	}
	
	var fields []*tmplCairoField
	for _, member := range members {
		fields = append(fields, &tmplCairoField{
			Type:      bindCairoType(member.Type),
			Name:      ToCamelCase(member.Name),
			CairoType: member.Type,
		})
	}
	
	structName := ToCamelCase(name)
	if structName == "" {
		structName = fmt.Sprintf("Struct%d", len(structs))
	}
	
	structs[id] = &tmplCairoStruct{
		Name:   structName,
		Fields: fields,
	}
	return structName
}
