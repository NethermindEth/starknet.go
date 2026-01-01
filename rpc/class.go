package rpc

import (
	"context"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/client/rpcerr"
)

// Class retrieves the class information from the Provider with the given hash.
//
// Parameters:
//   - ctx: The context.Context object
//   - blockID: The BlockID object
//   - classHash: The *felt.Felt object
//
// Returns:
//   - ClassOutput: The output of the class.
//   - error: An error if any occurred during the execution.
func (provider *Provider) Class(
	ctx context.Context,
	blockID BlockID,
	classHash *felt.Felt,
) (ClassOutput, error) {
	var rawClass map[string]any
	if err := do(ctx, provider.c, "starknet_getClass", &rawClass, blockID, classHash); err != nil {
		return nil, rpcerr.UnwrapToRPCErr(err, ErrClassHashNotFound, ErrBlockNotFound)
	}

	return typecastClassOutput(rawClass)
}
