package account

import (
	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/rpc"
)

/*
Formats the multicall transactions in a format which can be signed and verified by the network and OpenZeppelin account contracts
*/
func FmtCalldata(fnCalls []rpc.FunctionCall) []*felt.Felt {
	callArray := []*felt.Felt{}
	callData := []*felt.Felt{new(felt.Felt).SetUint64(uint64(len(fnCalls)))}

	for _, tx := range fnCalls {
		callData = append(callData, tx.ContractAddress, tx.EntryPointSelector)

		if len(tx.Calldata) == 0 {
			callData = append(callData, &felt.Zero, &felt.Zero)
			continue
		}

		callData = append(callData, new(felt.Felt).SetUint64(uint64(len(callArray))), new(felt.Felt).SetUint64(uint64(len(tx.Calldata))+1))
		for _, cd := range tx.Calldata {
			callArray = append(callArray, cd)
		}
	}
	callData = append(callData, new(felt.Felt).SetUint64(uint64(len(callArray)+1)))
	callData = append(callData, callArray...)
	callData = append(callData, new(felt.Felt).SetUint64(0))
	return callData
}
