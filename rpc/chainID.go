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
func (provider *Provider) ChainID(ctx context.Context) (string, error) {
	if provider.chainID != "" {
		return provider.chainID, nil
	}
	var result string
	if err := do(ctx, provider.c, "starknet_chainId", &result); err != nil {
		return "", rpcerr.UnwrapToRPCErr(err)
	}
	provider.chainID = internalUtils.HexToShortStr(result)

	return provider.chainID, nil
}
