package abi

import (
	"fmt"
	"math/big"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/utils"
)

type Type struct {
	Type string

	Components []Type

	Size int

	IsArray bool

	IsStruct bool
}

func NewType(t string) (Type, error) {
	parsedType := Type{
		Type: t,
	}

	arrayRegex := regexp.MustCompile(`^(.+)\[(\d*)\]$`)
	if matches := arrayRegex.FindStringSubmatch(t); len(matches) > 0 {
		parsedType.IsArray = true
		baseType, err := NewType(matches[1])
		if err != nil {
			return Type{}, err
		}
		parsedType.Components = []Type{baseType}

		if matches[2] != "" {
			size, err := strconv.Atoi(matches[2])
			if err != nil {
				return Type{}, fmt.Errorf("invalid array size: %v", err)
			}
			parsedType.Size = size
		}
	}

	return parsedType, nil
}

func ParseCairoType(cairoType string) (string, error) {
	parts := strings.Split(cairoType, "::")
	baseType := parts[len(parts)-1]

	switch baseType {
	case "felt252":
		return "*felt.Felt", nil
	case "u8", "u16", "u32":
		return "uint32", nil
	case "u64":
		return "uint64", nil
	case "u128":
		return "uint64", nil // Go doesn't have uint128, use uint64 or big.Int
	case "u256":
		return "*big.Int", nil
	case "bool":
		return "bool", nil
	case "ContractAddress":
		return "*felt.Felt", nil
	}

	if strings.HasSuffix(baseType, "]") {
		arrayRegex := regexp.MustCompile(`^(.+)\[(\d*)\]$`)
		if matches := arrayRegex.FindStringSubmatch(baseType); len(matches) > 0 {
			elementType, err := ParseCairoType(matches[1])
			if err != nil {
				return "", err
			}
			return "[]" + elementType, nil
		}
	}

	return "*felt.Felt", nil
}

func PackArguments(args []Argument, values []interface{}) ([]*felt.Felt, error) {
	if len(args) != len(values) {
		return nil, fmt.Errorf("argument count mismatch: %d args, %d values", len(args), len(values))
	}

	var result []*felt.Felt

	for i, arg := range args {
		packed, err := packArgument(arg, values[i])
		if err != nil {
			return nil, fmt.Errorf("failed to pack argument %s: %v", arg.Name, err)
		}
		result = append(result, packed...)
	}

	return result, nil
}

func packArgument(arg Argument, value interface{}) ([]*felt.Felt, error) {
	switch {
	case strings.Contains(arg.Type, "felt252") || strings.Contains(arg.Type, "ContractAddress"):
		switch v := value.(type) {
		case *felt.Felt:
			return []*felt.Felt{v}, nil
		case string:
			f, err := utils.HexToFelt(v)
			if err != nil {
				return nil, err
			}
			return []*felt.Felt{f}, nil
		default:
			return nil, fmt.Errorf("cannot convert %T to felt.Felt", value)
		}

	case strings.Contains(arg.Type, "u256"):
		switch v := value.(type) {
		case *big.Int:
			u256Felt, err := utils.HexToU256Felt("0x" + v.Text(16))
			if err != nil {
				return nil, err
			}
			return u256Felt, nil
		case string:
			bigInt, ok := new(big.Int).SetString(v, 0)
			if !ok {
				return nil, fmt.Errorf("invalid big.Int string: %s", v)
			}
			u256Felt, err := utils.HexToU256Felt("0x" + bigInt.Text(16))
			if err != nil {
				return nil, err
			}
			return u256Felt, nil
		default:
			return nil, fmt.Errorf("cannot convert %T to u256", value)
		}

	case strings.Contains(arg.Type, "bool"):
		switch v := value.(type) {
		case bool:
			if v {
				return []*felt.Felt{utils.Uint64ToFelt(1)}, nil
			}
			return []*felt.Felt{utils.Uint64ToFelt(0)}, nil
		default:
			return nil, fmt.Errorf("cannot convert %T to bool", value)
		}

	case strings.Contains(arg.Type, "u8"), strings.Contains(arg.Type, "u16"), strings.Contains(arg.Type, "u32"), strings.Contains(arg.Type, "u64"), strings.Contains(arg.Type, "u128"):
		var num uint64
		switch v := value.(type) {
		case uint8:
			num = uint64(v)
		case uint16:
			num = uint64(v)
		case uint32:
			num = uint64(v)
		case uint64:
			num = v
		case int:
			if v < 0 {
				return nil, fmt.Errorf("negative value %d for unsigned type", v)
			}
			num = uint64(v)
		default:
			return nil, fmt.Errorf("cannot convert %T to uint", value)
		}
		return []*felt.Felt{utils.Uint64ToFelt(num)}, nil

	case strings.HasSuffix(arg.Type, "]"): // Array type
		arrayValue := reflect.ValueOf(value)
		if arrayValue.Kind() != reflect.Slice && arrayValue.Kind() != reflect.Array {
			return nil, fmt.Errorf("expected slice or array for %s, got %T", arg.Type, value)
		}

		var result []*felt.Felt
		result = append(result, utils.Uint64ToFelt(uint64(arrayValue.Len())))

		elementType := strings.TrimSuffix(strings.TrimSuffix(arg.Type, "]"), "[")
		for i := 0; i < arrayValue.Len(); i++ {
			element := arrayValue.Index(i).Interface()
			elementArg := Argument{Type: elementType}
			packed, err := packArgument(elementArg, element)
			if err != nil {
				return nil, fmt.Errorf("failed to pack array element %d: %v", i, err)
			}
			result = append(result, packed...)
		}
		return result, nil

	default:
		switch v := value.(type) {
		case *felt.Felt:
			return []*felt.Felt{v}, nil
		case string:
			f, err := utils.HexToFelt(v)
			if err != nil {
				return nil, err
			}
			return []*felt.Felt{f}, nil
		default:
			return nil, fmt.Errorf("cannot convert %T to felt.Felt for type %s", value, arg.Type)
		}
	}
}

func UnpackValues(args []Argument, data []*felt.Felt) ([]interface{}, error) {
	if len(args) == 0 {
		return []interface{}{}, nil
	}

	var result []interface{}
	offset := 0

	for _, arg := range args {
		value, consumed, err := unpackValue(arg, data[offset:])
		if err != nil {
			return nil, fmt.Errorf("failed to unpack value for %s: %v", arg.Name, err)
		}
		result = append(result, value)
		offset += consumed
	}

	return result, nil
}

func unpackValue(arg Argument, data []*felt.Felt) (interface{}, int, error) {
	if len(data) == 0 {
		return nil, 0, fmt.Errorf("no data to unpack")
	}

	switch {
	case strings.Contains(arg.Type, "felt252") || strings.Contains(arg.Type, "ContractAddress"):
		return data[0], 1, nil

	case strings.Contains(arg.Type, "u256"):
		if len(data) < 2 {
			return nil, 0, fmt.Errorf("insufficient data for u256")
		}
		low := utils.FeltToBigInt(data[0])
		high := utils.FeltToBigInt(data[1])
		
		highShifted := new(big.Int).Lsh(high, 128)
		result := new(big.Int).Add(highShifted, low)
		
		return result, 2, nil

	case strings.Contains(arg.Type, "bool"):
		val := data[0].BigInt(new(big.Int))
		return val.Cmp(big.NewInt(0)) != 0, 1, nil

	case strings.Contains(arg.Type, "u8"), strings.Contains(arg.Type, "u16"), strings.Contains(arg.Type, "u32"):
		val := data[0].BigInt(new(big.Int))
		if !val.IsUint64() {
			return nil, 0, fmt.Errorf("value too large for uint32")
		}
		return uint32(val.Uint64()), 1, nil

	case strings.Contains(arg.Type, "u64"):
		val := data[0].BigInt(new(big.Int))
		if !val.IsUint64() {
			return nil, 0, fmt.Errorf("value too large for uint64")
		}
		return val.Uint64(), 1, nil

	case strings.Contains(arg.Type, "u128"):
		val := data[0].BigInt(new(big.Int))
		return val, 1, nil

	case strings.HasSuffix(arg.Type, "]"): // Array type
		if len(data) < 1 {
			return nil, 0, fmt.Errorf("insufficient data for array")
		}

		length := data[0].BigInt(new(big.Int))
		if !length.IsUint64() {
			return nil, 0, fmt.Errorf("invalid array length")
		}
		arrayLen := int(length.Uint64())

		elementType := strings.TrimSuffix(strings.TrimSuffix(arg.Type, "]"), "[")
		elementArg := Argument{Type: elementType}

		var result []interface{}
		consumed := 1 // Start after the length element
		for i := 0; i < arrayLen; i++ {
			if consumed >= len(data) {
				return nil, 0, fmt.Errorf("insufficient data for array elements")
			}
			element, elementConsumed, err := unpackValue(elementArg, data[consumed:])
			if err != nil {
				return nil, 0, fmt.Errorf("failed to unpack array element %d: %v", i, err)
			}
			result = append(result, element)
			consumed += elementConsumed
		}
		return result, consumed, nil

	default:
		return data[0], 1, nil
	}
}
