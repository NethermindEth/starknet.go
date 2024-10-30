package rpc

import (
	"encoding/json"
	"fmt"

	"github.com/NethermindEth/juno/core/felt"
)

type CasmCompiledContractClass struct {
	EntryPointsByType CasmEntryPointsByType `json:"entry_points_by_type"`
	ByteCode          []*felt.Felt          `json:"bytecode"`
	Prime             []*felt.Felt          `json:"prime"`
	CompilerVersion   string                `json:"compiler_version"`
	Hints             Hints                 `json:"hints"`
	// a list of sizes of segments in the bytecode, each segment is hashed invidually when computing the bytecode hash
	BytecodeSegmentLengths []int `json:"bytecode_segment_lengths,omitempty"`
}

type CasmEntryPointsByType struct {
	Constructor []CasmEntryPoint `json:"CONSTRUCTOR"`
	External    []CasmEntryPoint `json:"EXTERNAL"`
	L1Handler   []CasmEntryPoint `json:"L1_HANDLER"`
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

func (hints *Hints) GetValues() (int, []Hint) {
	return hints.Int, hints.HintArr
}

type Hint any // TODO: finish the description and create an unmarshall func to every 'any' type

type DeprecatedHint struct {
	any // is one of these types: DeprecatedHintEnum, AssertAllAccessesUsed, AssertLtAssertValidInput, Felt252DictRead, or Felt252DictWrite
}

// UnmarshalJSON unmarshals the JSON data into a DeprecatedHint object.
//
// Parameters:
// - data: The JSON data to be unmarshalled
// Returns:
// - error: An error if the unmarshalling process fails
func (depHint *DeprecatedHint) UnmarshalJSON(data []byte) error {
	var dec interface{}
	if err := json.Unmarshal(data, &dec); err != nil {
		return err
	}

	switch hint := dec.(type) {
	case string:
		*depHint = DeprecatedHint{DeprecatedHintEnum(hint)}
	case AssertAllAccessesUsed:
		*depHint = DeprecatedHint{hint}
	case AssertLtAssertValidInput:
		*depHint = DeprecatedHint{hint}
	case Felt252DictRead:
		*depHint = DeprecatedHint{hint}
	case Felt252DictWrite:
		*depHint = DeprecatedHint{hint}
	default:
		return fmt.Errorf("failed to unmarshal 'DeprecatedHint'")
	}

	return nil
}

type CoreHint struct {
	any // it can be a wide variety of types. Ref: https://github.com/starkware-libs/starknet-specs/blob/19ab9d4df4ae3acc1a52cde5a43a7cace08bcc4b/api/starknet_executables.json#L403
}

// UnmarshalJSON unmarshals the JSON data into a CoreHint object.
//
// Parameters:
// - data: The JSON data to be unmarshalled
// Returns:
// - error: An error if the unmarshalling process fails
func (coreHint *CoreHint) UnmarshalJSON(data []byte) error {
	var dec interface{}
	if err := json.Unmarshal(data, &dec); err != nil {
		return err
	}

	switch hint := dec.(type) {
	case AllocSegment:
		*coreHint = CoreHint{hint}
	case TestLessThan:
		*coreHint = CoreHint{hint}
	case TestLessThanOrEqual:
		*coreHint = CoreHint{hint}
	case TestLessThanOrEqualAddress:
		*coreHint = CoreHint{hint}
	case WideMul128:
		*coreHint = CoreHint{hint}
	case DivMod:
		*coreHint = CoreHint{hint}
	case Uint256DivMod:
		*coreHint = CoreHint{hint}
	case Uint512DivModByUint256:
		*coreHint = CoreHint{hint}
	case SquareRoot:
		*coreHint = CoreHint{hint}
	case Uint256SquareRoot:
		*coreHint = CoreHint{hint}
	case LinearSplit:
		*coreHint = CoreHint{hint}
	case AllocFelt252Dict:
		*coreHint = CoreHint{hint}
	case Felt252DictEntryInit:
		*coreHint = CoreHint{hint}
	case Felt252DictEntryUpdate:
		*coreHint = CoreHint{hint}
	case GetSegmentArenaIndex:
		*coreHint = CoreHint{hint}
	case InitSquashData:
		*coreHint = CoreHint{hint}
	case GetCurrentAccessIndex:
		*coreHint = CoreHint{hint}
	case ShouldSkipSquashLoop:
		*coreHint = CoreHint{hint}
	case GetCurrentAccessDelta:
		*coreHint = CoreHint{hint}
	case ShouldContinueSquashLoop:
		*coreHint = CoreHint{hint}
	case GetNextDictKey:
		*coreHint = CoreHint{hint}
	case AssertLeFindSmallArcs:
		*coreHint = CoreHint{hint}
	case AssertLeIsFirstArcExcluded:
		*coreHint = CoreHint{hint}
	case AssertLeIsSecondArcExcluded:
		*coreHint = CoreHint{hint}
	case RandomEcPoint:
		*coreHint = CoreHint{hint}
	case FieldSqrt:
		*coreHint = CoreHint{hint}
	case DebugPrint:
		*coreHint = CoreHint{hint}
	case AllocConstantSize:
		*coreHint = CoreHint{hint}
	case U256InvModN:
		*coreHint = CoreHint{hint}
	case EvalCircuit:
		*coreHint = CoreHint{hint}
	default:
		return fmt.Errorf("failed to unmarshal 'CoreHint'")
	}

	return nil
}

type StarknetHint struct {
	any // is one of these types: SystemCall or Cheatcode
}

// UnmarshalJSON unmarshals the JSON data into a StarknetHint object.
//
// Parameters:
// - data: The JSON data to be unmarshalled
// Returns:
// - error: An error if the unmarshalling process fails
func (strkHint *StarknetHint) UnmarshalJSON(data []byte) error {
	var dec interface{}
	if err := json.Unmarshal(data, &dec); err != nil {
		return err
	}

	switch hint := dec.(type) {
	case SystemCall:
		*strkHint = StarknetHint{hint}
	case Cheatcode:
		*strkHint = StarknetHint{hint}
	default:
		return fmt.Errorf("failed to unmarshal 'StarknetHint'")
	}

	return nil
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
	Register Register `json:"Register"`
	Offset   int      `json:"offset"`
}

type Register string

const (
	AP Register = "AP"
	FP Register = "FP"
)

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

type ResOperand any // is one of these types: Deref, DoubleDeref, Immediate, or BinOp

type Deref CellRef

// A (CellRef, offsest) tuple, but adapted to a golang struct
type DoubleDeref struct {
	CellRef CellRef
	Offset  int
}

func (dd *DoubleDeref) GetValues() (CellRef, int) {
	return dd.CellRef, dd.Offset
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
	B         any       `json:"b"` // is one of these types: Deref or Immediate
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
	FirstKeys  CellRef    `json:"first_key"`
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
	TOrR0     CellRef    `json:"t_or_k0"`
	TOrR1     CellRef    `json:"t_or_k1"`
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
