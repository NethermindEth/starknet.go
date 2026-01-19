package rpc_test

import (
	"testing"

	"github.com/NethermindEth/starknet.go/internal/tests"
)

func TestSubscribeTransactionStatus(t *testing.T) {
	tests.RunTestOn(t, tests.TestnetEnv)

	t.Skip("flaky test. It will be fixed in the next starknet.go release")
}
