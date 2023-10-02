package starknetgo

import (
	"fmt"
	"math/big"

	"github.com/NethermindEth/starknet.go/types"
)

// fmtCalldataStrings formats the given list of function calls into a list of calldata strings.
//
// It takes in a slice of types.FunctionCall called 'calls' and returns a slice of strings called 'calldataStrings'.
func fmtCalldataStrings(calls []types.FunctionCall) (calldataStrings []string) {
	callArray := fmtCalldata(calls)
	for _, data := range callArray {
		calldataStrings = append(calldataStrings, fmt.Sprintf("0x%x", data))
	}
	return calldataStrings
}

// fmtCalldata formats the given list of function calls into a list of calldata (which can be signed and verified by the network and OpenZeppelin account contracts).
//
// The function takes in a slice of types.FunctionCall called 'calls', which represents the list of function calls.
// It returns a slice of *big.Int called 'calldataArray', which contains the generated calldata.
func fmtCalldata(calls []types.FunctionCall) (calldataArray []*big.Int) {
	callArray := []*big.Int{big.NewInt(int64(len(calls)))}

	for _, tx := range calls {
		address := tx.ContractAddress.BigInt(big.NewInt(0))
		callArray = append(callArray, address, types.GetSelectorFromName(tx.EntryPointSelector.String()))

		if len(tx.Calldata) == 0 {
			callArray = append(callArray, big.NewInt(0), big.NewInt(0))

			continue
		}

		callArray = append(callArray, big.NewInt(int64(len(calldataArray))), big.NewInt(int64(len(tx.Calldata))))
		for _, cd := range tx.Calldata {
			calldataArray = append(calldataArray, types.SNValToBN(cd.String()))
		}
	}

	callArray = append(callArray, big.NewInt(int64(len(calldataArray))))
	callArray = append(callArray, calldataArray...)
	return callArray
}
