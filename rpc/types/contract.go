package types

type EntryPoint struct {
	// The offset of the entry point in the program
	Offset NumAsHex `json:"offset"`
	// A unique identifier of the entry point (function) in the program
	Selector string `json:"selector"`
}

type ABI []ABIEntry

type EntryPointsByType struct {
	Constructor []EntryPoint `json:"CONSTRUCTOR"`
	External    []EntryPoint `json:"EXTERNAL"`
	L1Handler   []EntryPoint `json:"L1_HANDLER"`
}

type ContractClass struct {
	// Program A base64 representation of the compressed program code
	Program string `json:"program"`

	EntryPointsByType EntryPointsByType `json:"entry_points_by_type"`

	Abi *ABI `json:"abi,omitempty"`
}

type ABIEntry interface {
	IsType() string
}

type StructABIType string

const (
	StructABITypeEvent StructABIType = "struct"
)

type EventABIType string

const (
	EventABITypeEvent EventABIType = "event"
)

type FunctionABIType string

const (
	FunctionABITypeFunction  FunctionABIType = "function"
	FunctionABITypeL1Handler FunctionABIType = "l1_handler"
)

type StructABIEntry struct {
	// The event type
	Type StructABIType `json:"type"`

	// The event name
	Name string `json:"name"`

	Size uint64 `json:"size"`

	Members []StructMember `json:"members"`
}

type StructMember struct {
	TypedParameter
	Offset uint64 `json:"offset"`
}

type EventABIEntry struct {
	// The event type
	Type EventABIType `json:"type"`

	// The event name
	Name string `json:"name"`

	Keys []TypedParameter `json:"keys"`

	Data TypedParameter `json:"data"`
}

type FunctionABIEntry struct {
	// The function type
	Type FunctionABIType `json:"type"`

	// The function name
	Name string `json:"name"`

	Inputs []TypedParameter `json:"inputs"`

	Outputs []TypedParameter `json:"outputs"`
}

type TypedParameter struct {
	// The parameter's name
	Name string `json:"name"`

	// The parameter's type
	Type string `json:"type"`
}
