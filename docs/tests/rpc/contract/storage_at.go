package main
 
import (
    "context"
    "encoding/hex"
    "fmt"
    "log"
    "math/big"
    "os"
    "strings"
 
    "github.com/NethermindEth/starknet.go/rpc"
    "github.com/NethermindEth/starknet.go/utils"
    "github.com/joho/godotenv"
)
 
func main() {
    // Load environment variables from .env file
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }
 
    // Get RPC URL from environment variable
    rpcURL := os.Getenv("STARKNET_RPC_URL")
    if rpcURL == "" {
        log.Fatal("STARKNET_RPC_URL not found in .env file")
    }
 
    // Initialize provider
    provider, err := rpc.NewProvider(context.Background(), rpcURL)
    if err != nil {
        log.Fatal(err)
    }
 
    ctx := context.Background()
 
    // ETH contract address on Starknet (SPELOIA)
    ethContract, err := utils.HexToFelt("0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7")
    if err != nil {
        log.Fatal(err)
    }
 
    // Read token name from storage using variable name
    // The method automatically converts "ERC20_name" to its storage key
    nameValue, err := provider.StorageAt(ctx, ethContract, "ERC20_name", rpc.WithBlockTag("latest"))
    if err != nil {
        log.Fatal(err)
    }
 
    // Decode the hex-encoded string to readable text
    tokenName := decodeShortString(nameValue)
    fmt.Printf("Token Name: %s (raw: %s)\n", tokenName, nameValue)
 
    // Read token symbol from storage
    symbolValue, err := provider.StorageAt(ctx, ethContract, "ERC20_symbol", rpc.WithBlockTag("latest"))
    if err != nil {
        log.Fatal(err)
    }
    tokenSymbol := decodeShortString(symbolValue)
    fmt.Printf("Token Symbol: %s (raw: %s)\n", tokenSymbol, symbolValue)
 
    // Read decimals value and convert to integer
    decimalsValue, err := provider.StorageAt(ctx, ethContract, "ERC20_decimals", rpc.WithBlockTag("latest"))
    if err != nil {
        log.Fatal(err)
    }
    decimals := new(big.Int)
    decimals.SetString(strings.TrimPrefix(decimalsValue, "0x"), 16)
    fmt.Printf("Token Decimals: %s (raw: %s)\n", decimals.String(), decimalsValue)
 
    // Query historical total supply at a specific block number
    totalSupplyValue, err := provider.StorageAt(ctx, ethContract, "ERC20_total_supply", rpc.WithBlockNumber(100000))
    if err != nil {
        log.Fatal(err)
    }
    totalSupply := new(big.Int)
    totalSupply.SetString(strings.TrimPrefix(totalSupplyValue, "0x"), 16)
    fmt.Printf("Total Supply at block 100000: %s (raw: %s)\n", totalSupply.String(), totalSupplyValue)
}
 
// decodeShortString decodes a hex-encoded short string (Cairo felt string)
func decodeShortString(hexStr string) string {
    hexStr = strings.TrimPrefix(hexStr, "0x")
    if len(hexStr)%2 != 0 {
        hexStr = "0" + hexStr
    }
 
    bytes, err := hex.DecodeString(hexStr)
    if err != nil {
        return ""
    }
 
    return string(bytes)
}