package rpc

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/NethermindEth/juno/core/felt"
)

type CasmCompiledContractClass struct {
	EntryPointsByType CasmEntryPointsByType `json:"entry_points_by_type"`
	ByteCode          []*felt.Felt          `json:"bytecode"`
	Prime             NumAsHex              `json:"prime"`
	CompilerVersion   string                `json:"compiler_version"`
	Hints             []Hints               `json:"hints"`
	// a list of sizes of segments in the bytecode, each segment is hashed individually when computing the bytecode hash
	BytecodeSegmentLengths []int `json:"bytecode_segment_lengths,omitempty"`
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

type CasmEntryPointsByType struct {
	Constructor []CasmEntryPoint `json:"CONSTRUCTOR"`
	External    []CasmEntryPoint `json:"EXTERNAL"`
	L1Handler   []CasmEntryPoint `json:"L1_HANDLER"`
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

type CasmEntryPoint struct {
	DeprecatedCairoEntryPoint
	// the hash of the right child
	Builtin []string `json:"builtins"`
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
func (h *Hints) MarshalJSON() ([]byte, error) {
	return json.Marshal([2]interface{}{h.Int, h.HintArr})
}

func (hints *Hints) Values() (int, []Hint) {
	return hints.Int, hints.HintArr
}

func (hints *Hints) Tuple() [2]any {
	return [2]any{hints.Int, hints.HintArr}
}

// Can have only one of the following hints
type Hint struct {
	DeprecatedHint DeprecatedHint `json:",omitempty"`
	CoreHint       CoreHint       `json:",omitempty"`
	StarknetHint   StarknetHint   `json:",omitempty"`
}

// Validate ensures only one hint type is set
func (h *Hint) Validate() error {
	count := 0
	if !reflect.ValueOf(h.DeprecatedHint).IsZero() {
		count++
	}
	if !reflect.ValueOf(h.CoreHint).IsZero() {
		count++
	}
	if !reflect.ValueOf(h.StarknetHint).IsZero() {
		count++
	}
	if count != 1 {
		return fmt.Errorf("exactly one hint type must be set, got %d", count)
	}
	return nil
}

// UnmarshalJSON implements json.Unmarshaler
func (h *Hint) UnmarshalJSON(data []byte) error {
	type HintAlias Hint
	aux := &HintAlias{}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	*h = Hint(*aux)
	return h.Validate()
}

// MarshalJSON implements json.Marshaler interface
func (h *Hint) MarshalJSON() ([]byte, error) {
	if err := h.Validate(); err != nil {
		return nil, err
	}
	return json.Marshal(*h)
}

type DeprecatedHint struct {
	DeprecatedHintEnum
	AssertAllAccessesUsed    AssertAllAccessesUsed    `json:",omitempty"`
	AssertLtAssertValidInput AssertLtAssertValidInput `json:",omitempty"`
	Felt252DictRead          Felt252DictRead          `json:",omitempty"`
	Felt252DictWrite         Felt252DictWrite         `json:",omitempty"`
}

// Can have only one of the following hints
type CoreHint struct {
	AllocConstantSize           AllocConstantSize           `json:",omitempty"`
	AllocFelt252Dict            AllocFelt252Dict            `json:",omitempty"`
	AllocSegment                AllocSegment                `json:",omitempty"`
	AssertLeFindSmallArcs       AssertLeFindSmallArcs       `json:",omitempty"`
	AssertLeIsFirstArcExcluded  AssertLeIsFirstArcExcluded  `json:",omitempty"`
	AssertLeIsSecondArcExcluded AssertLeIsSecondArcExcluded `json:",omitempty"`
	DebugPrint                  DebugPrint                  `json:",omitempty"`
	DivMod                      DivMod                      `json:",omitempty"`
	EvalCircuit                 EvalCircuit                 `json:",omitempty"`
	Felt252DictEntryInit        Felt252DictEntryInit        `json:",omitempty"`
	Felt252DictEntryUpdate      Felt252DictEntryUpdate      `json:",omitempty"`
	FieldSqrt                   FieldSqrt                   `json:",omitempty"`
	GetCurrentAccessDelta       GetCurrentAccessDelta       `json:",omitempty"`
	GetCurrentAccessIndex       GetCurrentAccessIndex       `json:",omitempty"`
	GetNextDictKey              GetNextDictKey              `json:",omitempty"`
	GetSegmentArenaIndex        GetSegmentArenaIndex        `json:",omitempty"`
	InitSquashData              InitSquashData              `json:",omitempty"`
	LinearSplit                 LinearSplit                 `json:",omitempty"`
	RandomEcPoint               RandomEcPoint               `json:",omitempty"`
	ShouldContinueSquashLoop    ShouldContinueSquashLoop    `json:",omitempty"`
	ShouldSkipSquashLoop        ShouldSkipSquashLoop        `json:",omitempty"`
	SquareRoot                  SquareRoot                  `json:",omitempty"`
	TestLessThan                TestLessThan                `json:",omitempty"`
	TestLessThanOrEqual         TestLessThanOrEqual         `json:",omitempty"`
	TestLessThanOrEqualAddress  TestLessThanOrEqualAddress  `json:",omitempty"`
	U256InvModN                 U256InvModN                 `json:",omitempty"`
	Uint256DivMod               Uint256DivMod               `json:",omitempty"`
	Uint256SquareRoot           Uint256SquareRoot           `json:",omitempty"`
	Uint512DivModByUint256      Uint512DivModByUint256      `json:",omitempty"`
	WideMul128                  WideMul128                  `json:",omitempty"`
}

type StarknetHint struct {
	Cheatcode  Cheatcode  `json:",omitempty"`
	SystemCall SystemCall `json:",omitempty"`
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

type ResOperand struct {
	BinOp       BinOp       `json:",omitempty"`
	Deref       Deref       `json:",omitempty"`
	DoubleDeref DoubleDeref `json:",omitempty"`
	Immediate   Immediate   `json:",omitempty"`
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
func (dd *DoubleDeref) MarshalJSON() ([]byte, error) {
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

type B struct {
	Deref     Deref     `json:",omitempty"`
	Immediate Immediate `json:",omitempty"`
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
	DictAccess ResOperand `json:"dict_access"`
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
