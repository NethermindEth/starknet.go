package rpc

import (
	"context"

	"github.com/NethermindEth/starknet.go/client/rpcerr"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
)

// ChainID returns the chain ID for transaction replay protection.
//
// Parameters:
//   - ctx: The context.Context object for the function
//
// Returns:
//   - string: The chain ID
//   - error: An error if any occurred during the execution
func ChainID(ctx context.Context, c callCloser) (string, error) {
	var result string
	if err := do(ctx, c, "starknet_chainId", &result); err != nil {
		return "", rpcerr.UnwrapToRPCErr(err)
	}

	return internalUtils.HexToShortStr(result), nil
}
