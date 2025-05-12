package rpc

import (
	"context"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/contracts"
)

// Get the CASM code resulting from compiling a given class
//
// Parameters:
//   - ctx: The context.Context used for the request
//   - classHash: The hash of the contract class whose CASM will be returned
//
// Returns:
//   - CasmCompiledContractClass: The compiled contract class
//   - error: An error if any occurred during the execution
func (provider *Provider) CompiledCasm(ctx context.Context, classHash *felt.Felt) (*contracts.CasmClass, error) {
	var result contracts.CasmClass
	if err := do(ctx, provider.c, "starknet_getCompiledCasm", &result, classHash); err != nil {
		return nil, tryUnwrapToRPCErr(err, ErrClassHashNotFound, ErrCompilationError)
	}

	return &result, nil
}
