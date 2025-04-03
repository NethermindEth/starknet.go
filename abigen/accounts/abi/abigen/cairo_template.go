package abigen

import (
	_ "embed"
	"strings"

	cairoabi "github.com/NethermindEth/starknet.go/abigen/accounts/abi"
)

var tmplCairoSource string

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
		Name string
	}
	Const      bool
	Inputs     []tmplCairoField
	Outputs    []tmplCairoField
}

type tmplCairoEvent struct {
	Original   cairoabi.Event
	Normalized struct {
		Name string
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
	words := strings.FieldsFunc(input, func(r rune) bool {
		return r == '_' || r == ' ' || r == '-'
	})
	
	for i, word := range words {
		if i == 0 {
			words[i] = strings.Title(word)
		} else {
			words[i] = strings.Title(word)
		}
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
		normalizedStruct := struct{ Name string }{
			Name: ToCamelCase(name),
		}

		isConst := method.StateMutability == "view"

		methods[name] = tmplCairoMethod{
			Original:   method,
			Normalized: normalizedStruct,
			Const:      isConst,
		}
	}

	for name, event := range abi.Events {
		normalizedStruct := struct{ Name string }{
			Name: ToCamelCase(name),
		}

		events[name] = tmplCairoEvent{
			Original:   event,
			Normalized: normalizedStruct,
		}
	}

	return &cairoBinder{
		abi:     abi,
		methods: methods,
		events:  events,
	}
}
