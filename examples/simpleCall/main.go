package main

import (
	"context"
	"fmt"
	"os"

	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/types"
	"github.com/NethermindEth/starknet.go/utils"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/joho/godotenv"
)

var (
	name                string = "mainnet"
	someMainnetContract string = "0x024dE48Fb640DB135B3dc85ef0FE2789e032FbCA2fca54E58aB8dB93ca22F767"
	contractMethod      string = "getName"
)

func main() {
	fmt.Println("Starting simpeCall example")
	godotenv.Load(fmt.Sprintf(".env.%s", name))
	base := os.Getenv("INTEGRATION_BASE")
	c, err := ethrpc.DialContext(context.Background(), base)
	if err != nil {
		fmt.Println("Failed to connect to the client, did you specify the url in the .env.mainnet?")
		panic(err)
	}
	clientv02 := rpc.NewProvider(c)
	fmt.Println("Established connection with the client")

	contractAddress, err := utils.HexToFelt(someMainnetContract)
	if err != nil {
		panic(err)
	}

	// Make read contract call
	tx := rpc.FunctionCall{
		ContractAddress:    contractAddress,
		EntryPointSelector: types.GetSelectorFromNameFelt(contractMethod),
	}

	fmt.Println("Making Call() request")
	callResp, err := clientv02.Call(context.Background(), tx, rpc.BlockID{Tag: "latest"})
	if err != nil {
		panic(err.Error())
	}

	fmt.Println(fmt.Sprintf("Response to %s():%s ", contractMethod, callResp[0]))
}
