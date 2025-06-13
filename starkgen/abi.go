package starkgen

// Item represents a contract item ABI
// @todo make it a any type and use a switch statement to handle the different types
type Item struct {
	Type string `json:"type"`
	// The actual content is one of the following types based on the Type field
	Function    *Function    `json:"function,omitempty"`
	Constructor *Constructor `json:"constructor,omitempty"`
	L1Handler   *L1Handler   `json:"l1_handler,omitempty"`
	Event       *Event       `json:"event,omitempty"`
	Struct      *Struct      `json:"struct,omitempty"`
	Enum        *Enum        `json:"enum,omitempty"`
	Interface   *Interface   `json:"interface,omitempty"`
	Impl        *Imp         `json:"impl,omitempty"`
}

// Imp represents a contract implementation ABI
type Imp struct {
	Type          string `json:"type"`
	Name          string `json:"name"`
	InterfaceName string `json:"interface_name"`
}

// Interface represents a contract interface ABI
type Interface struct {
	Type  string `json:"type"`
	Name  string `json:"name"`
	Items []Item `json:"items"`
}

// StateMutability represents the state mutability of a function
type StateMutability string

const (
	External StateMutability = "external"
	View     StateMutability = "view"
)

// Function represents a contract function ABI
type Function struct {
	Type            string          `json:"type"`
	Name            string          `json:"name"`
	Inputs          []Input         `json:"inputs"`
	Outputs         []Output        `json:"outputs"`
	StateMutability StateMutability `json:"state_mutability"`
}

// Constructor represents a contract constructor ABI
type Constructor struct {
	Type   string  `json:"type"`
	Name   string  `json:"name"`
	Inputs []Input `json:"inputs"`
}

// L1Handler represents a contract L1 handler ABI
type L1Handler struct {
	Type            string          `json:"type"`
	Name            string          `json:"name"`
	Inputs          []Input         `json:"inputs"`
	Outputs         []Output        `json:"outputs"`
	StateMutability StateMutability `json:"state_mutability"`
}

// Event represents a contract event ABI
type Event struct {
	Type     string       `json:"type"`
	Name     string       `json:"name"`
	Kind     EventKind    `json:"kind"`
	Members  []EventField `json:"members,omitempty"`
	Variants []EventField `json:"variants,omitempty"`
}

// EventKind is a string enum representing the kind of an event
type EventKind string

const (
	EventKindStruct EventKind = "struct"
	EventKindEnum   EventKind = "enum"
)

// EventField represents a field in an event
type EventField struct {
	Name string         `json:"name"`
	Type string         `json:"type"`
	Kind EventFieldKind `json:"kind"`
}

// EventFieldKind represents how to serialize the event's field
type EventFieldKind string

const (
	KeySerde  EventFieldKind = "key"
	DataSerde EventFieldKind = "data"
	Nested    EventFieldKind = "nested"
	Flat      EventFieldKind = "flat"
)

// Input represents a function input ABI
type Input struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// Output represents a function output ABI
type Output struct {
	Type string `json:"type"`
}

// Struct represents a struct ABI
type Struct struct {
	Type    string         `json:"type"`
	Name    string         `json:"name"`
	Members []StructMember `json:"members"`
}

// StructMember represents a struct member
type StructMember struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// Enum represents an enum ABI
type Enum struct {
	Type     string        `json:"type"`
	Name     string        `json:"name"`
	Variants []EnumVariant `json:"variants"`
}

// EnumVariant represents an enum variant
type EnumVariant struct {
	Name string `json:"name"`
	Type string `json:"type"`
}
