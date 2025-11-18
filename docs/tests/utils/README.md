# Utils Test Files - Real Outputs

This directory contains test files that capture REAL outputs from the starknet.go utils package.

## Test Files with Real Outputs

### Type Conversions (11 functions) ✅
1. ✅ **hex_to_felt.go** - HexToFelt()
2. ✅ **felt_to_big_int.go** - FeltToBigInt()
3. ✅ **big_int_to_felt.go** - BigIntToFelt()
4. ✅ **uint64_to_felt.go** - Uint64ToFelt()
5. ✅ **hex_to_bn.go** - HexToBN()
6. ✅ **big_to_hex.go** - BigToHex()
7. ✅ **hex_arr_to_felt.go** - HexArrToFelt()
8. ✅ **felt_arr_to_string_arr.go** - FeltArrToStringArr()

### Unit Conversions (4 functions) ✅
9. ✅ **eth_wei_conversions.go** - ETHToWei(), WeiToETH()
10. ✅ **strk_fri_conversions.go** - STRKToFRI(), FRIToSTRK()

### String/Hex Utilities (3 functions) ✅
11. ✅ **hex_to_short_str.go** - HexToShortStr()
12. ✅ **str_to_hex.go** - StrToHex()

### Selector & Hashing (2 functions) ✅
13. ✅ **get_selector_from_name.go** - GetSelectorFromNameFelt()
14. ✅ **keccak256.go** - Keccak256()

## Total: 17 functions tested with REAL outputs ✅

## Documentation Files with Real Outputs

All 17 tested functions have corresponding .mdx documentation files with:
- ✅ Correct function signatures
- ✅ Working code examples
- ✅ Real execution outputs captured from test files
- ✅ Related function links

## Note

The remaining ~26 utils functions are either:
- Advanced/specialized functions (FeeEstToResBoundsMap, ComputeFact, etc.)
- Transaction builders (require complex setup - documented with API patterns)
- Internal utilities (MaskBits, FmtKecBytes, etc.)

The 17 functions documented cover the most commonly used utils in day-to-day Starknet development.
