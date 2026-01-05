package methods

import (
	"context"

	"github.com/NethermindEth/starknet.go/client/rpcerr"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/rpc/internal"
)

// ChainID returns the chain ID for transaction replay protection.
//
// Parameters:
//   - ctx: The context.Context object for the function
//
// Returns:
//   - string: The chain ID
//   - error: An error if any occurred during the execution
func ChainID(ctx context.Context, c rpc.Caller) (string, error) {
	var result string
	if err := internal.Do(ctx, c, "starknet_chainId", &result); err != nil {
		return "", rpcerr.UnwrapToRPCErr(err)
	}

	return internalUtils.HexToShortStr(result), nil
}
