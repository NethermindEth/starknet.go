package rpc

import (
	"context"

	"github.com/NethermindEth/juno/core/felt"
)

// Get the contract class definition in the given block associated with the given hash
//
// Parameters:
// - ctx: The context.Context used for the request
// - classHash: The hash of the contract class whose CASM will be returned
// Returns:
// - *felt.Felt: The compiled contract class
// - error: An error if any occurred during the execution
func (provider *Provider) CompiledCasm(ctx context.Context, classHash *felt.Felt) (*felt.Felt, error) {
	var result *felt.Felt
	if err := do(ctx, provider.c, "starknet_getCompiledCasm", &result, classHash); err != nil {

		return nil, tryUnwrapToRPCErr(err, ErrContractNotFound, ErrBlockNotFound)
	}
	return result, nil
}
