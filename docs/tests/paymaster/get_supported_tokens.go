package main

import (
	"context"
	"fmt"

	"github.com/NethermindEth/starknet.go/paymaster"
)

func main() {
	fmt.Println("GetSupportedTokens Demo:")
	fmt.Println("  Method: pm.GetSupportedTokens(ctx)")
	fmt.Println()
	fmt.Println("Returns: []TokenData")
	fmt.Println()
	fmt.Println("TokenData structure:")
	fmt.Println("  - TokenAddress: Address of the fee token")
	fmt.Println("  - TokenSymbol: Symbol (e.g., 'ETH', 'STRK')")
	fmt.Println("  - TokenDecimals: Number of decimals")
	fmt.Println()
	fmt.Println("Example response:")
	fmt.Println("  [{")
	fmt.Println("    TokenAddress: 0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7")
	fmt.Println("    TokenSymbol: ETH")
	fmt.Println("    TokenDecimals: 18")
	fmt.Println("  }, {")
	fmt.Println("    TokenAddress: 0x04718f5a0fc34cc1af16a1cdee98ffb20c31f5cd61d6ab07201858f4287c938d")
	fmt.Println("    TokenSymbol: STRK")
	fmt.Println("    TokenDecimals: 18")
	fmt.Println("  }]")
	
	_ = context.Background()
	_ = paymaster.TokenData{}
}
