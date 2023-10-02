package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/NethermindEth/starknet.go/gateway"
	"github.com/urfave/cli/v2"
)

var transactionCommand = cli.Command{
	Name:    "get_transaction",
	Aliases: []string{"t"},
	Usage:   "get a transaction",
	Flags:   transactionFlags,
	Action:  transaction,
}

type friendlyTransaction gateway.StarknetTransaction

// MarshalJSON marshals the friendlyTransaction object to JSON.
//
// It returns a byte slice and an error. The byte slice contains the JSON representation of the object,
// while the error indicates any error encountered during the marshaling process.
func (f friendlyTransaction) MarshalJSON() ([]byte, error) {
	content := map[string]interface{}{}
	content["treansaction_hash"] = f.Transaction.TransactionHash
	content["contract_address"] = f.Transaction.ContractAddress
	content["class_hash"] = f.Transaction.ClassHash
	content["type"] = f.Transaction.Type
	content["version"] = f.Transaction.Version
	return json.Marshal(content)
}

var transactionFlags = []cli.Flag{
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
		Name:     "hash",
		Usage:    "provides a specific transaction hash",
		Required: true,
	},
	&cli.StringFlag{
		Name:  "format",
		Usage: "display different display format",
		Value: "friendly",
	},
}

// transaction is a function that processes a transaction using the given cli.Context.
//
// It takes the following parameters:
// - cCtx: a pointer to a cli.Context object that contains information about the transaction.
//
// It returns an error if the transaction is not supported by the provider, if the transaction hash is invalid,
// or if there is an error retrieving the transaction.
func transaction(cCtx *cli.Context) error {
	providerName := cCtx.Value("provider")
	if providerName.(string) != "gateway" {
		return fmt.Errorf("provider not supported")
	}
	format := cCtx.Value("format").(string)

	opts := gateway.TransactionOptions{}
	transactionHash := cCtx.Value("hash").(string)
	if !strings.HasPrefix(transactionHash, "0x") {
		return fmt.Errorf("invalid tx %s", transactionHash)
	}
	opts.TransactionHash = transactionHash
	baseURL := cCtx.Value("base-url")
	if baseURL.(string) == "" {
		baseURL = "https://alpha4.starknet.io"
	}
	provider := gateway.NewProvider(gateway.WithBaseURL(baseURL.(string)))
	transaction, err := provider.Transaction(context.Background(), opts)
	if err != nil {
		return err
	}
	var output []byte
	switch format {
	case "friendly":
		output, err = json.MarshalIndent(friendlyTransaction(*transaction), " ", "    ")
	default:
		output, err = json.MarshalIndent(*transaction, " ", "    ")
	}
	if err != nil {
		return err
	}
	fmt.Println(string(output))
	return nil
}
