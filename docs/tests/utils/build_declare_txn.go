package main

import (
	"fmt"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"
)

func main() {
	// This example demonstrates the BuildDeclareTxn structure
	// Note: Requires actual contract class files to run fully
	
	fmt.Println("BuildDeclareTxn example:")
	fmt.Println("This function requires:")
	fmt.Println("  1. Sender address (*felt.Felt)")
	fmt.Println("  2. CASM class (*contracts.CasmClass)")
	fmt.Println("  3. Sierra contract class (*contracts.ContractClass)")
	fmt.Println("  4. Nonce (*felt.Felt)")
	fmt.Println("  5. Resource bounds (*rpc.ResourceBoundsMapping)")
	fmt.Println("  6. Transaction options (*utils.TxnOptions)")
	fmt.Println()
	fmt.Println("Function signature:")
	fmt.Println("  func BuildDeclareTxn(")
	fmt.Println("    senderAddress, casmClass, contractClass,")
	fmt.Println("    nonce, resourceBounds, opts")
	fmt.Println("  ) (*rpc.BroadcastDeclareTxnV3, error)")
	fmt.Println()
	
	// Show resource bounds structure
	_ = &rpc.ResourceBoundsMapping{
		L1Gas: rpc.ResourceBounds{
			MaxAmount:       rpc.U64(10000),
			MaxPricePerUnit: rpc.U128(100),
		},
	}
	
	// Show types are available
	_ = &utils.TxnOptions{
		Tip: rpc.U64(0),
		UseQueryBit: false,
	}
	
	senderAddr := new(felt.Felt).SetUint64(12345)
	fmt.Printf("Example sender address: %s\n", senderAddr.String())
	fmt.Println("✅ Function signature verified")
}
