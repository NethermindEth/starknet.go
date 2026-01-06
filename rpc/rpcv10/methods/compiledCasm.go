package methods

import (
	"context"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/client/rpcerr"
	"github.com/NethermindEth/starknet.go/contracts"
	"github.com/NethermindEth/starknet.go/rpc/callers"
	"github.com/NethermindEth/starknet.go/rpc/internal"
	"github.com/NethermindEth/starknet.go/rpc/rpcv10"
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
func CompiledCasm(
	ctx context.Context,
	c callers.Caller,
	classHash *felt.Felt,
) (*contracts.CasmClass, error) {
	var result contracts.CasmClass
	if err := internal.Do(ctx, c, "starknet_getCompiledCasm", &result, classHash); err != nil {
		return nil, rpcerr.UnwrapToRPCErr(
			err,
			rpcv10.ErrClassHashNotFound,
			rpcv10.ErrCompilationError,
		)
	}

	return &result, nil
}
