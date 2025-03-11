package contracts

import (
	"encoding/json"
	"fmt"

	"github.com/NethermindEth/juno/core/felt"
)

type CasmCompiledContractClass struct {
	EntryPointsByType CasmEntryPointsByType `json:"entry_points_by_type"`
	ByteCode          []*felt.Felt          `json:"bytecode"`
	Prime             NumAsHex              `json:"prime"`
	CompilerVersion   string                `json:"compiler_version"`
	Hints             []Hints               `json:"hints"`
	// a list of sizes of segments in the bytecode, each segment is hashed individually when computing the bytecode hash
	BytecodeSegmentLengths *NestedUints `json:"bytecode_segment_lengths,omitempty"`
}

// Validate ensures all required fields are present and valid
func (c *CasmCompiledContractClass) Validate() error {
	if c.ByteCode == nil {
		return fmt.Errorf("bytecode is required")
	}
	if c.Prime == "" {
		return fmt.Errorf("prime is required")
	}
	if c.CompilerVersion == "" {
		return fmt.Errorf("compiler_version is required")
	}
	if c.Hints == nil {
		return fmt.Errorf("hints is required")
	}
	if err := c.EntryPointsByType.Validate(); err != nil {
		return fmt.Errorf("entry_points_by_type validation failed: %w", err)
	}
	return nil
}

// UnmarshalJSON implements json.Unmarshaler
func (c *CasmCompiledContractClass) UnmarshalJSON(data []byte) error {
	type Alias CasmCompiledContractClass
	aux := &Alias{}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	*c = CasmCompiledContractClass(*aux)
	return c.Validate()
}

// An integer number in hex format (0x...)
// TODO: duplicate of rpc.NumAsHex to avoid import cycle. Maybe move to a shared 'types' package?
type NumAsHex string

type CasmEntryPointsByType struct {
	Constructor []CasmClassEntryPoint `json:"CONSTRUCTOR"`
	External    []CasmClassEntryPoint `json:"EXTERNAL"`
	L1Handler   []CasmClassEntryPoint `json:"L1_HANDLER"`
}

// Validate ensures all required fields are present and valid
func (e *CasmEntryPointsByType) Validate() error {
	if e.Constructor == nil {
		return fmt.Errorf("CONSTRUCTOR is required")
	}
	if e.External == nil {
		return fmt.Errorf("EXTERNAL is required")
	}
	if e.L1Handler == nil {
		return fmt.Errorf("L1_HANDLER is required")
	}
	return nil
}

// 2-tuple of pc value and an array of hints to execute, but adapted to a golang struct
type Hints struct {
	Int     int
	HintArr []Hint
}

// UnmarshalJSON implements json.Unmarshaler interface
func (h *Hints) UnmarshalJSON(data []byte) error {
	var tuple []json.RawMessage
	if err := json.Unmarshal(data, &tuple); err != nil {
		return err
	}

	if len(tuple) != 2 {
		return fmt.Errorf("expected tuple of length 2, got %d", len(tuple))
	}

	// Unmarshal the first element (integer)
	if err := json.Unmarshal(tuple[0], &h.Int); err != nil {
		return fmt.Errorf("failed to unmarshal Int: %w", err)
	}

	// Unmarshal the second element (array of hints)
	if err := json.Unmarshal(tuple[1], &h.HintArr); err != nil {
		return fmt.Errorf("failed to unmarshal HintArr: %w", err)
	}

	return nil
}

// MarshalJSON implements json.Marshaler interface
func (h Hints) MarshalJSON() ([]byte, error) {
	return json.Marshal([2]interface{}{h.Int, h.HintArr})
}

func (hints *Hints) Values() (int, []Hint) {
	return hints.Int, hints.HintArr
}

func (hints *Hints) Tuple() [2]any {
	return [2]any{hints.Int, hints.HintArr}
}

// Can be one of the following hints
type Hint struct {
	Type string
	Data interface{}
}

// UnmarshalJSON implements json.Unmarshaler
func (h *Hint) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as string enum first
	var enumVal string
	if err := json.Unmarshal(data, &enumVal); err == nil {
		switch DeprecatedHintEnum(enumVal) {
		case AssertCurrentAccessIndicesIsEmpty, AssertAllKeysUsed, AssertLeAssertThirdArcExcluded:
			h.Type = "enum"
			h.Data = DeprecatedHintEnum(enumVal)
			return nil
		}
	}

	// If not an enum, try to unmarshal as an object
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if len(raw) != 1 {
		return fmt.Errorf("hint must have exactly one field, got %d", len(raw))
	}

	// Get the single key and value
	var hintType string
	var hintData json.RawMessage
	for k := range raw {
		hintType = k
		hintData = raw[k]
	}

	// Deprecated Hint
	switch hintType {
	case "AssertAllAccessesUsed":
		return unmarshalJSON[AssertAllAccessesUsed](hintData, hintType, h)
	case "AssertLtAssertValidInput":
		return unmarshalJSON[AssertLtAssertValidInput](hintData, hintType, h)
	case "Felt252DictRead":
		return unmarshalJSON[Felt252DictRead](hintData, hintType, h)
	case "Felt252DictWrite":
		return unmarshalJSON[Felt252DictWrite](hintData, hintType, h)
	}

	// Core Hint
	switch hintType {
	case "AllocConstantSize":
		return unmarshalJSON[AllocConstantSize](hintData, hintType, h)
	case "AllocFelt252Dict":
		return unmarshalJSON[AllocFelt252Dict](hintData, hintType, h)
	case "AllocSegment":
		return unmarshalJSON[AllocSegment](hintData, hintType, h)
	case "AssertLeFindSmallArcs":
		return unmarshalJSON[AssertLeFindSmallArcs](hintData, hintType, h)
	case "AssertLeIsFirstArcExcluded":
		return unmarshalJSON[AssertLeIsFirstArcExcluded](hintData, hintType, h)
	case "AssertLeIsSecondArcExcluded":
		return unmarshalJSON[AssertLeIsSecondArcExcluded](hintData, hintType, h)
	case "DebugPrint":
		return unmarshalJSON[DebugPrint](hintData, hintType, h)
	case "DivMod":
		return unmarshalJSON[DivMod](hintData, hintType, h)
	case "EvalCircuit":
		return unmarshalJSON[EvalCircuit](hintData, hintType, h)
	case "Felt252DictEntryInit":
		return unmarshalJSON[Felt252DictEntryInit](hintData, hintType, h)
	case "Felt252DictEntryUpdate":
		return unmarshalJSON[Felt252DictEntryUpdate](hintData, hintType, h)
	case "FieldSqrt":
		return unmarshalJSON[FieldSqrt](hintData, hintType, h)
	case "GetCurrentAccessDelta":
		return unmarshalJSON[GetCurrentAccessDelta](hintData, hintType, h)
	case "GetCurrentAccessIndex":
		return unmarshalJSON[GetCurrentAccessIndex](hintData, hintType, h)
	case "GetNextDictKey":
		return unmarshalJSON[GetNextDictKey](hintData, hintType, h)
	case "GetSegmentArenaIndex":
		return unmarshalJSON[GetSegmentArenaIndex](hintData, hintType, h)
	case "InitSquashData":
		return unmarshalJSON[InitSquashData](hintData, hintType, h)
	case "LinearSplit":
		return unmarshalJSON[LinearSplit](hintData, hintType, h)
	case "RandomEcPoint":
		return unmarshalJSON[RandomEcPoint](hintData, hintType, h)
	case "ShouldContinueSquashLoop":
		return unmarshalJSON[ShouldContinueSquashLoop](hintData, hintType, h)
	case "ShouldSkipSquashLoop":
		return unmarshalJSON[ShouldSkipSquashLoop](hintData, hintType, h)
	case "SquareRoot":
		return unmarshalJSON[SquareRoot](hintData, hintType, h)
	case "TestLessThan":
		return unmarshalJSON[TestLessThan](hintData, hintType, h)
	case "TestLessThanOrEqual":
		return unmarshalJSON[TestLessThanOrEqual](hintData, hintType, h)
	case "TestLessThanOrEqualAddress":
		return unmarshalJSON[TestLessThanOrEqualAddress](hintData, hintType, h)
	case "U256InvModN":
		return unmarshalJSON[U256InvModN](hintData, hintType, h)
	case "Uint256DivMod":
		return unmarshalJSON[Uint256DivMod](hintData, hintType, h)
	case "Uint256SquareRoot":
		return unmarshalJSON[Uint256SquareRoot](hintData, hintType, h)
	case "Uint512DivModByUint256":
		return unmarshalJSON[Uint512DivModByUint256](hintData, hintType, h)
	case "WideMul128":
		return unmarshalJSON[WideMul128](hintData, hintType, h)
	}

	// Starknet Hint
	switch hintType {
	case "Cheatcode":
		return unmarshalJSON[Cheatcode](hintData, hintType, h)
	case "SystemCall":
		return unmarshalJSON[SystemCall](hintData, hintType, h)
	}

	return fmt.Errorf("unknown hint type: %s", hintType)
}

func unmarshalJSON[T any](hintData []byte, hintType string, h *Hint) error {
	var hint T
	if err := json.Unmarshal(hintData, &hint); err != nil {
		return fmt.Errorf("failed to unmarshal %T: %w", hint, err)
	}
	h.Type = hintType
	h.Data = hint
	return nil
}

// MarshalJSON implements json.Marshaler
func (h Hint) MarshalJSON() ([]byte, error) {
	if h.Type == "" || h.Data == nil {
		return nil, fmt.Errorf("hint type and data must be set")
	}

	// For enum types, marshal directly as string
	if h.Type == "enum" {
		return json.Marshal(h.Data)
	}

	// Handle types that contain ResOperand or B fields which require custom JSON marshaling.
	// These types need special handling due to their polymorphic field types.
	// For all other types, default JSON marshaling will be used.
	switch typedHint := h.Data.(type) {
	// Deprecated Hint
	case AssertLtAssertValidInput, Felt252DictRead, Felt252DictWrite:
		return marshalJSON(typedHint, &h)
	// Core Hint
	case TestLessThan, TestLessThanOrEqual, TestLessThanOrEqualAddress,
		WideMul128, DivMod, LinearSplit, AllocFelt252Dict,
		Felt252DictEntryInit, Felt252DictEntryUpdate, GetSegmentArenaIndex,
		InitSquashData, GetCurrentAccessIndex, AssertLeFindSmallArcs,
		FieldSqrt, DebugPrint, AllocConstantSize, U256InvModN, EvalCircuit:
		return marshalJSON(typedHint, &h)
	// Starknet Hint
	case Cheatcode, SystemCall:
		return marshalJSON(typedHint, &h)
	default:
		// For all other types, use default marshaling
		return json.Marshal(map[string]interface{}{
			h.Type: h.Data,
		})
	}
}

func marshalJSON[T any](typedHint T, h *Hint) ([]byte, error) {
	rawHint, err := json.Marshal(typedHint)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal %T: %w", typedHint, err)
	}
	return json.Marshal(map[string]json.RawMessage{h.Type: rawHint})
}

type DeprecatedHintEnum string

const (
	AssertCurrentAccessIndicesIsEmpty DeprecatedHintEnum = "AssertCurrentAccessIndicesIsEmpty"
	AssertAllKeysUsed                 DeprecatedHintEnum = "AssertAllKeysUsed"
	AssertLeAssertThirdArcExcluded    DeprecatedHintEnum = "AssertLeAssertThirdArcExcluded"
)

type AssertAllAccessesUsed struct {
	NotUsedAccesses CellRef `json:"n_used_accesses"`
}

type CellRef struct {
	Register Register `json:"register"`
	Offset   int      `json:"offset"`
}

type Register string

const (
	AP Register = "AP"
	FP Register = "FP"
)

// UnmarshalJSON implements json.Unmarshaler
func (r *Register) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	switch Register(s) {
	case AP, FP:
		*r = Register(s)
		return nil
	default:
		return fmt.Errorf("invalid register value: %s, must be either AP or FP", s)
	}
}

// MarshalJSON implements json.Marshaler
func (r Register) MarshalJSON() ([]byte, error) {
	if r != AP && r != FP {
		return nil, fmt.Errorf("invalid register value: %s, must be either AP or FP", r)
	}
	return json.Marshal(string(r))
}

type AssertLtAssertValidInput struct {
	A ResOperand `json:"a"`
	B ResOperand `json:"b"`
}

type felt252Dict struct {
	DictPtr ResOperand `json:"dict_ptr"`
	Key     ResOperand `json:"key"`
}

type Felt252DictRead struct {
	felt252Dict
	ValueDst CellRef `json:"value_dst"`
}

type Felt252DictWrite struct {
	felt252Dict
	Value ResOperand `json:"value"`
}

type Felt252DictEntryInit felt252Dict

type Felt252DictEntryUpdate struct {
	DictPtr ResOperand `json:"dict_ptr"`
	Value   ResOperand `json:"value"`
}

// Can be one of the following values
type ResOperand struct {
	Type string
	Data interface{}
}

// UnmarshalJSON implements json.Unmarshaler
func (r *ResOperand) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if len(raw) != 1 {
		return fmt.Errorf("res operand must have exactly one field, got %d", len(raw))
	}

	// Get the single key and value
	var hintType string
	var hintData json.RawMessage
	for k := range raw {
		hintType = k
		hintData = raw[k]
	}

	var err error

	switch hintType {
	case "BinOp":
		var binOp BinOp
		err = json.Unmarshal(hintData, &binOp)
		r.Data = binOp
	case "Deref":
		var deref Deref
		err = json.Unmarshal(hintData, &deref)
		r.Data = deref
	case "DoubleDeref":
		var doubleDeref DoubleDeref
		err = json.Unmarshal(hintData, &doubleDeref)
		r.Data = doubleDeref
	case "Immediate":
		var immediate Immediate
		err = json.Unmarshal(hintData, &immediate)
		r.Data = immediate
	default:
		return fmt.Errorf("unknown res operand type: %s", hintType)
	}

	if err != nil {
		return fmt.Errorf("failed to unmarshal res operand as any known type: %w", err)
	}

	r.Type = hintType
	return nil
}

// MarshalJSON implements json.Marshaler
func (r ResOperand) MarshalJSON() ([]byte, error) {
	if r.Type == "" || r.Data == nil {
		return nil, fmt.Errorf("res operand type and data must be set")
	}

	return json.Marshal(map[string]interface{}{
		r.Type: r.Data,
	})
}

type Deref CellRef

// A (CellRef, offset) tuple, but adapted to a golang struct
type DoubleDeref struct {
	CellRef CellRef
	Offset  int
}

// UnmarshalJSON implements json.Unmarshaler
func (dd *DoubleDeref) UnmarshalJSON(data []byte) error {
	var tuple []json.RawMessage
	if err := json.Unmarshal(data, &tuple); err != nil {
		return err
	}

	if len(tuple) != 2 {
		return fmt.Errorf("expected tuple of length 2, got %d", len(tuple))
	}

	// Unmarshal CellRef
	if err := json.Unmarshal(tuple[0], &dd.CellRef); err != nil {
		return fmt.Errorf("failed to unmarshal CellRef: %w", err)
	}

	// Unmarshal offset
	if err := json.Unmarshal(tuple[1], &dd.Offset); err != nil {
		return fmt.Errorf("failed to unmarshal Offset: %w", err)
	}

	return nil
}

// MarshalJSON implements json.Marshaler
func (dd DoubleDeref) MarshalJSON() ([]byte, error) {
	return json.Marshal([2]interface{}{dd.CellRef, dd.Offset})
}

func (dd *DoubleDeref) Values() (CellRef, int) {
	return dd.CellRef, dd.Offset
}

func (dd *DoubleDeref) Tuple() [2]any {
	return [2]any{dd.CellRef, dd.Offset}
}

type Immediate NumAsHex

type Operation string

const (
	Add Operation = "Add"
	Mul Operation = "Mul"
)

type BinOp struct {
	Operation Operation `json:"op"`
	A         CellRef   `json:"a"`
	B         B         `json:"b"`
}

// Can be one of the following values
type B struct {
	Type string
	Data interface{}
}

// UnmarshalJSON implements json.Unmarshaler
func (b *B) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if len(raw) != 1 {
		return fmt.Errorf("B must have exactly one field, got %d", len(raw))
	}

	// Get the single key and value
	var hintType string
	var hintData json.RawMessage
	for k := range raw {
		hintType = k
		hintData = raw[k]
	}

	switch hintType {
	case "Deref":
		var deref Deref
		if err := json.Unmarshal(hintData, &deref); err != nil {
			return err
		}
		b.Data = deref
	case "Immediate":
		var immediate Immediate
		if err := json.Unmarshal(hintData, &immediate); err != nil {
			return err
		}
		b.Data = immediate
	default:
		return fmt.Errorf("unknown B type: %s", hintType)
	}

	b.Type = hintType
	return nil
}

// MarshalJSON implements json.Marshaler
func (b B) MarshalJSON() ([]byte, error) {
	if b.Type == "" || b.Data == nil {
		return nil, fmt.Errorf("B type and data must be set")
	}

	return json.Marshal(map[string]interface{}{
		b.Type: b.Data,
	})
}

type AllocSegment struct {
	Dst CellRef `json:"dst"`
}

type baseLhsRhs struct {
	Lhs ResOperand `json:"lhs"`
	Rhs ResOperand `json:"rhs"`
}

type TestLessThan struct {
	baseLhsRhs
	Dst CellRef `json:"dst"`
}

type TestLessThanOrEqual TestLessThan

type TestLessThanOrEqualAddress TestLessThan

type WideMul128 struct {
	baseLhsRhs
	High CellRef `json:"high"`
	Low  CellRef `json:"low"`
}

type DivMod struct {
	baseLhsRhs
	Quotient  CellRef `json:"quotient"`
	Remainder CellRef `json:"remainder"`
}

type Uint256DivMod struct {
	Dividend0  ResOperand `json:"dividend0"`
	Dividend1  ResOperand `json:"dividend1"`
	Divisor0   ResOperand `json:"divisor0"`
	Divisor1   ResOperand `json:"divisor1"`
	Quotient0  CellRef    `json:"quotient0"`
	Quotient1  CellRef    `json:"quotient1"`
	Remainder0 CellRef    `json:"remainder0"`
	Remainder1 CellRef    `json:"remainder1"`
}

type Uint512DivModByUint256 struct {
	Uint256DivMod
	Dividend2 ResOperand `json:"dividend2"`
	Dividend3 ResOperand `json:"dividend3"`
	Quotient2 CellRef    `json:"quotient2"`
	Quotient3 CellRef    `json:"quotient3"`
}

type SquareRoot struct {
	Value ResOperand `json:"value"`
	Dst   CellRef    `json:"dst"`
}

type Uint256SquareRoot struct {
	ValueLow                     ResOperand `json:"value_low"`
	ValueHigh                    ResOperand `json:"value_high"`
	Sqrt0                        CellRef    `json:"sqrt0"`
	Sqrt1                        CellRef    `json:"sqrt1"`
	RemainderLow                 CellRef    `json:"remainder_low"`
	RemainderHigh                CellRef    `json:"remainder_high"`
	SqrtMul2MinusRemainderGeU128 CellRef    `json:"sqrt_mul_2_minus_remainder_ge_u128"`
}

type LinearSplit struct {
	Value  ResOperand `json:"value"`
	Scalar ResOperand `json:"scalar"`
	MaxX   ResOperand `json:"max_x"`
	X      CellRef    `json:"x"`
	Y      CellRef    `json:"y"`
}

type AllocFelt252Dict struct {
	SegmentArenaPtr ResOperand `json:"segment_arena_ptr"`
}

type GetSegmentArenaIndex struct {
	DictEndPtr ResOperand `json:"dict_end_ptr"`
	DictIndex  CellRef    `json:"dict_index"`
}

type InitSquashData struct {
	DictAccess ResOperand `json:"dict_accesses"`
	PtrDiff    ResOperand `json:"ptr_diff"`
	NAccesses  ResOperand `json:"n_accesses"`
	BigKeys    CellRef    `json:"big_keys"`
	FirstKey   CellRef    `json:"first_key"`
}

type GetCurrentAccessIndex struct {
	RangeCheckPtr ResOperand `json:"range_check_ptr"`
}

type GetCurrentAccessDelta struct {
	IndexDeltaMinus1 CellRef `json:"index_delta_minus1"`
}

type GetNextDictKey struct {
	NextKey CellRef `json:"next_key"`
}

type ShouldSkipSquashLoop struct {
	ShouldSkipLoop CellRef `json:"should_skip_loop"`
}

type ShouldContinueSquashLoop struct {
	ShouldContinue CellRef `json:"should_continue"`
}

type AssertLeFindSmallArcs struct {
	RangeCheckPtr ResOperand `json:"range_check_ptr"`
	A             ResOperand `json:"a"`
	B             ResOperand `json:"b"`
}

type AssertLeIsFirstArcExcluded struct {
	SkipExcludeAFlag CellRef `json:"skip_exclude_a_flag"`
}

type AssertLeIsSecondArcExcluded struct {
	SkipExcludeBMinusA CellRef `json:"skip_exclude_b_minus_a"`
}

type RandomEcPoint struct {
	X CellRef `json:"x"`
	Y CellRef `json:"y"`
}

type FieldSqrt struct {
	Val  ResOperand `json:"val"`
	Sqrt CellRef    `json:"sqrt"`
}

type DebugPrint struct {
	Start ResOperand `json:"start"`
	End   ResOperand `json:"end"`
}

type AllocConstantSize struct {
	Size ResOperand `json:"size"`
	Dst  CellRef    `json:"dst"`
}

type U256InvModN struct {
	B0        ResOperand `json:"b0"`
	B1        ResOperand `json:"b1"`
	N0        ResOperand `json:"n0"`
	N1        ResOperand `json:"n1"`
	G0OrNoInv CellRef    `json:"g0_or_no_inv"`
	G1Option  CellRef    `json:"g1_option"`
	SOrR0     CellRef    `json:"s_or_r0"`
	SOrR1     CellRef    `json:"s_or_r1"`
	TOrK0     CellRef    `json:"t_or_k0"`
	TOrK1     CellRef    `json:"t_or_k1"`
}

type EvalCircuit struct {
	NAddMods      ResOperand `json:"n_add_mods"`
	AddModBuiltin ResOperand `json:"add_mod_builtin"`
	NMulMods      ResOperand `json:"n_mul_mods"`
	MulModBuiltin ResOperand `json:"mul_mod_builtin"`
}

type SystemCall struct {
	System ResOperand `json:"system"`
}

type Cheatcode struct {
	Selector    NumAsHex   `json:"selector"`
	InputStart  ResOperand `json:"input_start"`
	InputEnd    ResOperand `json:"input_end"`
	OutputStart CellRef    `json:"output_start"`
	OutputEnd   CellRef    `json:"output_end"`
}
