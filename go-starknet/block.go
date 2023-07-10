package main

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/NethermindEth/starknet.go/gateway"
	"github.com/urfave/cli/v2"
)

var blockCommand = cli.Command{
	Name:    "get_block",
	Aliases: []string{"b"},
	Usage:   "get a block",
	Flags:   blockFlags,
	Action:  block,
}

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
	format := cCtx.Value("format").(string)
	blockOptions := &gateway.BlockOptions{}
	blockID := cCtx.Value("block-id").(string)
	switch blockID {
	case "":
		blockOptions.Tag = "pending"
	case "pending", "latest":
		blockOptions.Tag = blockID
	default:
		match := false
		if strings.HasPrefix(blockID, "0x") {
			match = true
			blockOptions.BlockHash = blockID
		}
		if ok, _ := regexp.MatchString("^[0-9]+$", blockID); ok {
			match = true
			block, err := strconv.ParseUint(blockID, 10, 64)
			if err != nil {
				return fmt.Errorf("invalid block %s", blockID)
			}
			blockOptions.BlockNumber = &block
		}
		if !match {
			return fmt.Errorf("invalid block %s", blockID)
		}
	}
	baseURL := cCtx.Value("base-url")
	if baseURL.(string) == "" {
		baseURL = "https://alpha4.starknet.io"
	}
	provider := gateway.NewProvider(gateway.WithBaseURL(baseURL.(string)))
	block, err := provider.Block(context.Background(), blockOptions)
	if err != nil {
		return err
	}
	var output []byte
	switch format {
	case "friendly":
		output, err = json.MarshalIndent(friendlyBlock(*block), " ", "    ")
	default:
		output, err = json.MarshalIndent(*block, " ", "    ")
	}
	if err != nil {
		return err
	}
	fmt.Println(string(output))
	return nil
}
