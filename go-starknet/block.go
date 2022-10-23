package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dontpanicdao/caigo/gateway"
	"github.com/urfave/cli/v2"
)

type friendlyBlock gateway.Block

func (f friendlyBlock) MarshalJSON() ([]byte, error) {
	content := map[string]interface{}{}
	content["block_hash"] = f.BlockHash
	content["block_number"] = f.BlockNumber
	content["status"] = f.Status
	content["timestamp"] = time.Unix(int64(f.Timestamp), 0).Format("2006-01-02 15:04:05")
	transactions := map[string]struct {
		Total           int      `json:"total"`
		TransactionHash []string `json:"transaction_hashes"`
	}{}
	for _, v := range f.Transactions {
		num := transactions[v.Type]
		num.Total++
		num.TransactionHash = append(num.TransactionHash, v.TransactionHash)
		transactions[v.Type] = num
	}
	content["transactions"] = transactions
	return json.Marshal(content)
}

var blockFlags = []cli.Flag{
	&cli.StringFlag{
		Name:  "provider",
		Usage: "choose between the gateway and rpc",
		Value: "gateway",
	},
	&cli.StringFlag{
		Name:  "base-url",
		Usage: "change the default baseURL",
		Value: "",
	},
	&cli.StringFlag{
		Name:  "block-id",
		Usage: "provides a specific block id; for now, only latest is supported",
		Value: "latest",
	},
	&cli.StringFlag{
		Name:  "format",
		Usage: "display different display format",
		Value: "friendly",
	},
}

func block(cCtx *cli.Context) error {
	providerName := cCtx.Value("provider")
	if providerName.(string) != "gateway" {
		return fmt.Errorf("provider not supported")
	}
	num := cCtx.Value("block-id")
	if num.(string) != "latest" {
		return fmt.Errorf("not supported")
	}
	baseURL := cCtx.Value("base-url")
	if baseURL.(string) == "" {
		baseURL = "https://alpha4.starknet.io"
	}
	provider := gateway.NewProvider(gateway.WithBaseURL(baseURL.(string)))
	block, err := provider.Block(context.Background(), &gateway.BlockOptions{})
	if err != nil {
		return err
	}
	output, err := json.MarshalIndent(friendlyBlock(*block), " ", "    ")
	if err != nil {
		return err
	}
	fmt.Println(string(output))
	return nil
}
