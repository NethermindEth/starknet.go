package methods

import (
	"context"

	"github.com/NethermindEth/starknet.go/client/rpcerr"
	"github.com/NethermindEth/starknet.go/rpc/callers"
	"github.com/NethermindEth/starknet.go/rpc/internal"
	"github.com/NethermindEth/starknet.go/rpc/rpcv10"
)

// SpecVersion returns the version of the Starknet JSON-RPC specification being
// implemented by the node.
//
// Parameters:
//   - ctx: The context for the function.
//
// Returns:
//   - string: The version of the Starknet JSON-RPC specification
//     implemented by the node.
//   - error: An error if the request fails.
func SpecVersion(ctx context.Context, c callers.Caller) (string, error) {
	var result string
	err := internal.Do(ctx, c, "starknet_specVersion", &result)
	if err != nil {
		return "", rpcerr.Err(rpcerr.InternalError, rpcv10.StringErrData(err.Error()))
	}

	return result, nil
}
