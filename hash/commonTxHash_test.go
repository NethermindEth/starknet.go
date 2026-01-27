package hash

import (
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/rpc/rpcv10"
	"github.com/NethermindEth/starknet.go/rpc/rpcv9"
)

func TestTypes(t *testing.T) {
	// Declare transactions
	_, _ = TransactionHashDeclareV0(rpcv9.DeclareTxnV0{}, &felt.Felt{})
	_, _ = TransactionHashDeclareV1(rpcv9.DeclareTxnV1{}, &felt.Felt{})
	_, _ = TransactionHashDeclareV2(rpcv9.DeclareTxnV2{}, &felt.Felt{})
	_, _ = TransactionHashDeclareV3(rpcv9.DeclareTxnV3{}, &felt.Felt{}, &felt.Felt{})
	_, _ = TransactionHashDeclareV3(rpcv9.BroadcastDeclareTxnV3{}, &felt.Felt{}, &felt.Felt{})

	_, _ = TransactionHashDeclareV0(rpcv10.DeclareTxnV0{}, &felt.Felt{})
	_, _ = TransactionHashDeclareV1(rpcv10.DeclareTxnV1{}, &felt.Felt{})
	_, _ = TransactionHashDeclareV2(rpcv10.DeclareTxnV2{}, &felt.Felt{})
	_, _ = TransactionHashDeclareV3(rpcv10.DeclareTxnV3{}, &felt.Felt{}, &felt.Felt{})
	_, _ = TransactionHashDeclareV3(rpcv10.BroadcastDeclareTxnV3{}, &felt.Felt{}, &felt.Felt{})

	// Deploy transactions
	_, _ = TransactionHashDeployAccountV1(rpcv9.DeployAccountTxnV1{}, &felt.Felt{}, &felt.Felt{})
	_, _ = TransactionHashDeployAccountV3(rpcv9.DeployAccountTxnV3{}, &felt.Felt{}, &felt.Felt{})
	_, _ = TransactionHashDeployAccountV3(rpcv9.BroadcastDeployAccountTxnV3{}, &felt.Felt{}, &felt.Felt{})

	_, _ = TransactionHashDeployAccountV1(rpcv10.DeployAccountTxnV1{}, &felt.Felt{}, &felt.Felt{})
	_, _ = TransactionHashDeployAccountV3(rpcv10.DeployAccountTxnV3{}, &felt.Felt{}, &felt.Felt{})
	_, _ = TransactionHashDeployAccountV3(rpcv10.BroadcastDeployAccountTxnV3{}, &felt.Felt{}, &felt.Felt{})

	// Invoke transactions
	_, _ = TransactionHashInvokeV0(rpcv9.InvokeTxnV0{}, &felt.Felt{})
	_, _ = TransactionHashInvokeV1(rpcv9.InvokeTxnV1{}, &felt.Felt{})
	_, _ = TransactionHashInvokeV3(rpcv9.InvokeTxnV3{}, &felt.Felt{})
	_, _ = TransactionHashInvokeV3(rpcv9.BroadcastInvokeTxnV3{}, &felt.Felt{})

	_, _ = TransactionHashInvokeV0(rpcv10.InvokeTxnV0{}, &felt.Felt{})
	_, _ = TransactionHashInvokeV1(rpcv10.InvokeTxnV1{}, &felt.Felt{})
	_, _ = TransactionHashInvokeV3(rpcv10.InvokeTxnV3{}, &felt.Felt{})
	_, _ = TransactionHashInvokeV3(rpcv10.BroadcastInvokeTxnV3{}, &felt.Felt{})
}
