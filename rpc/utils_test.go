package rpc

import (
	"testing"

	"github.com/NethermindEth/starknet.go/rpc/rpcv10"
	"github.com/NethermindEth/starknet.go/rpc/rpcv9"
)

// @todo add tests
func TestTypes(t *testing.T) {
	_, _ = WaitForTransactionReceipt(&rpcv10.Provider{}, nil, nil, 0)
	_, _ = WaitForTransactionReceipt(&rpcv9.Provider{}, nil, nil, 0)
}
