package abigen

import (
	"fmt"
	"strings"

	cairoabi "github.com/NethermindEth/starknet.go/abigen/accounts/abi"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type tmplCairoData struct {
	Package   string
	Contracts map[string]*tmplCairoContract
	Structs   map[string]*tmplCairoStruct
}

type tmplCairoContract struct {
	Type        string
	InputABI    string
	InputBin    string
	Constructor cairoabi.Method
	Methods     map[string]*tmplCairoMethod
	Events      map[string]*tmplCairoEvent
}

func newTmplCairoContract(name, abi, bin string, constructor cairoabi.Method, binder *cairoBinder) *tmplCairoContract {
	name = ToCamelCase(name)
	
	methods := make(map[string]*tmplCairoMethod)
	events := make(map[string]*tmplCairoEvent)
	
	for methodName, method := range binder.methods {
		methodCopy := method
		methods[methodName] = &methodCopy
	}
	
	for eventName, event := range binder.events {
		eventCopy := event
		events[eventName] = &eventCopy
	}

	return &tmplCairoContract{
		Type:        name,
		InputABI:    abi,
		InputBin:    bin,
		Constructor: constructor,
		Methods:     methods,
		Events:      events,
	}
}

type tmplCairoMethod struct {
	Original   cairoabi.Method
	Normalized struct {
		Name    string
		Inputs  []tmplCairoField
		Outputs []tmplCairoField
	}
	Const      bool
	Structured bool
	Inputs     []tmplCairoField
	Outputs    []tmplCairoField
}

type tmplCairoEvent struct {
	Original   cairoabi.Event
	Normalized struct {
		Name string
		Keys []tmplCairoField
		Data []tmplCairoField
	}
	Keys       []tmplCairoField
	Data       []tmplCairoField
}

type tmplCairoField struct {
	Name      string
	Type      string
	CairoType string
}

type tmplCairoStruct struct {
	Name   string
	Fields []*tmplCairoField
}

func ToCamelCase(input string) string {
	input = strings.ReplaceAll(input, "::", "_")
	
	if strings.Contains(input, "_") && strings.HasSuffix(input, "_Event") {
		parts := strings.Split(input, "_")
		if len(parts) > 2 {
			input = parts[len(parts)-2] + "_Event"
		}
	}
	
	words := strings.FieldsFunc(input, func(r rune) bool {
		return r == '_' || r == ' ' || r == '-'
	})
	
	caser := cases.Title(language.English)
	for i, word := range words {
		words[i] = caser.String(word)
	}
	return strings.Join(words, "")
}

type cairoBinder struct {
	abi     cairoabi.ABI
	methods map[string]tmplCairoMethod
	events  map[string]tmplCairoEvent
}

func newCairoBinder(abi cairoabi.ABI) *cairoBinder {
	methods := make(map[string]tmplCairoMethod)
	events := make(map[string]tmplCairoEvent)

	for name, method := range abi.Methods {
		inputs := make([]tmplCairoField, len(method.Inputs))
		for i, input := range method.Inputs {
			inputs[i] = tmplCairoField{
				Name:      input.Name,
				Type:      input.Type,
				CairoType: input.Type,
			}
		}
		
		outputs := make([]tmplCairoField, len(method.Outputs))
		for i, output := range method.Outputs {
			outputs[i] = tmplCairoField{
				Name:      fmt.Sprintf("ret%d", i),
				Type:      output.Type,
				CairoType: output.Type,
			}
		}

		normalizedStruct := struct{ 
			Name string
			Inputs []tmplCairoField
			Outputs []tmplCairoField
		}{
			Name: ToCamelCase(name),
			Inputs: inputs,
			Outputs: outputs,
		}

		isConst := method.StateMutability == "view"

		methods[name] = tmplCairoMethod{
			Original:   method,
			Normalized: normalizedStruct,
			Const:      isConst,
			Inputs:     inputs,
			Outputs:    outputs,
		}
	}

	for name, event := range abi.Events {
		keys := make([]tmplCairoField, len(event.Keys))
		for i, key := range event.Keys {
			keys[i] = tmplCairoField{
				Name:      key.Name,
				Type:      key.Type,
				CairoType: key.Type,
			}
		}
		
		data := make([]tmplCairoField, len(event.Data))
		for i, d := range event.Data {
			data[i] = tmplCairoField{
				Name:      d.Name,
				Type:      d.Type,
				CairoType: d.Type,
			}
		}

		normalizedStruct := struct{ 
			Name string
			Keys []tmplCairoField
			Data []tmplCairoField
		}{
			Name: ToCamelCase(name),
			Keys: keys,
			Data: data,
		}

		events[name] = tmplCairoEvent{
			Original:   event,
			Normalized: normalizedStruct,
			Keys:       keys,
			Data:       data,
		}
	}

	return &cairoBinder{
		abi:     abi,
		methods: methods,
		events:  events,
	}
}
