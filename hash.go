package caigo

import (
	"fmt"
	"math/big"

	"github.com/NethermindEth/caigo/types"
)

func fmtCalldataStrings(calls []types.FunctionCall) (calldataStrings []string) {
	callArray := fmtCalldata(calls)
	for _, data := range callArray {
		calldataStrings = append(calldataStrings, fmt.Sprintf("0x%x", data))
	}
	return calldataStrings
}

/*
Formats the multicall transactions in a format which can be signed and verified by the network and OpenZeppelin account contracts
*/
func fmtCalldata(calls []types.FunctionCall) (calldataArray []*big.Int) {
	callArray := []*big.Int{big.NewInt(int64(len(calls)))}

	for _, tx := range calls {
		address := tx.ContractAddress.Big()
		callArray = append(callArray, address, types.GetSelectorFromName(tx.EntryPointSelector))

		if len(tx.Calldata) == 0 {
			callArray = append(callArray, big.NewInt(0), big.NewInt(0))

			continue
		}

		callArray = append(callArray, big.NewInt(int64(len(calldataArray))), big.NewInt(int64(len(tx.Calldata))))
		for _, cd := range tx.Calldata {
			calldataArray = append(calldataArray, types.SNValToBN(cd))
		}
	}

	callArray = append(callArray, big.NewInt(int64(len(calldataArray))))
	callArray = append(callArray, calldataArray...)
	return callArray
}
