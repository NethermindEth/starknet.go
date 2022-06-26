package rpc

import (
	"context"
	"math/big"
	"testing"
	// "github.com/dontpanicdao/caigo"
	"github.com/dontpanicdao/caigo/types"
)

var (
	client      *Client
	accountAddr = "0x028105caf03e1c4eb96b1c18d39d9f03bd53e5d2affd0874792e5bf05f3e529f"
	classHash   = "0x25ec026985a3bf9d0cc1fe17326b245dfdc3ff89b8fde106542a3ea56c5a918"
)

// Requires a Mainnet StarkNet JSON-RPC compliant node (e.g. pathfinder)
// (ref: https://github.com/eqlabs/pathfinder)
func init() {
	var err error
	client, err = DialContext(context.Background(), "http://localhost:9545")
	if err != nil {
		panic(err.Error())
	}
}

func TestClient(t *testing.T) {
	id, err := client.ChainID(context.Background())
	if err != nil {
		t.Errorf("Could not retrieve block: %v\n", err)
	}
	if id != "0x534e5f4d41494e" {
		t.Errorf("Not on mainnet node: %v\n", err)
	}

	_, err = client.Syncing(context.Background())
	if err != nil {
		t.Errorf("Could not retrieve sync information: %v\n", err)
	}
}

func TestBlock(t *testing.T) {
	num, err := client.BlockNumber(context.Background())
	if err != nil {
		t.Errorf("Could not retrieve block: %v\n", err)
	}

	block, err := client.BlockByNumber(context.Background(), num, "FULL_TXN_AND_RECEIPTS")
	if err != nil {
		t.Errorf("Could not retrieve block: %v\n", err)
	}
	if block.Transactions[0].TransactionReceipt.Status == "" {
		t.Errorf("Could not retrieve transaction receipts: %v\n", block.Transactions[0])
	}

	_, err = client.BlockByNumber(context.Background(), num, "FULL_TXNS")
	if err != nil {
		t.Errorf("Could not retrieve block: %v\n", err)
	}

	block, err = client.BlockByHash(context.Background(), block.BlockHash, "FULL_TXN_AND_RECEIPTS")
	if err != nil {
		t.Errorf("Could not retrieve block: %v\n", err)
	}
	if block.Transactions[0].TransactionReceipt.Status == "" {
		t.Errorf("Could not retrieve transaction receipts: %v\n", block.Transactions[0])
	}

	_, err = client.BlockByHash(context.Background(), block.BlockHash, "FULL_TXNS")
	if err != nil {
		t.Errorf("Could not retrieve block: %v\n", err)
	}
}

func TestTransaction(t *testing.T) {
	block, err := client.BlockByHash(context.Background(), "latest", "FULL_TXNS")
	if err != nil {
		t.Errorf("Could not retrieve block: %v\n", err)
	}

	_, err = client.TransactionByHash(context.Background(), block.Transactions[0].TransactionHash)
	if err != nil {
		t.Errorf("Could not retrieve block: %v\n", err)
	}

	_, err = client.TransactionReceipt(context.Background(), block.Transactions[0].TransactionHash)
	if err != nil {
		t.Errorf("Could not retrieve block: %v\n", err)
	}
}

func TestContract(t *testing.T) {
	block, err := client.BlockByNumber(context.Background(), big.NewInt(5), "FULL_TXNS")
	if err != nil {
		t.Errorf("Could not retrieve block: %v\n", err)
	}

	code, err := client.CodeAt(context.Background(), block.Transactions[0].ContractAddress)
	if err != nil {
		t.Errorf("Could not retrieve code: %v\n", err)
	}
	if len(code.Bytecode) == 0 {
		t.Errorf("Could not retrieve bytecode: %v\n", err)
	}

	// tested against pathfinder @ 0313b14ea1fad8f73635a3002d106908813e57f1
	classHash, err := client.ClassHashAt(context.Background(), accountAddr)
	if err != nil {
		t.Errorf("Could not retrieve class hash: %v\n", err)
	}

	_, err = client.ClassAt(context.Background(), accountAddr)
	if err != nil {
		t.Errorf("Could not retrieve class: %v\n", err)
	}

	_, err = client.Class(context.Background(), classHash)
	if err != nil {
		t.Errorf("Could not retrieve class: %v\n", err)
	}
}

func TestEvents(t *testing.T) {
	p := EventParams{
		FromBlock:  800,
		ToBlock:    1701,
		PageSize:   1000,
		PageNumber: 0,
	}
	events, err := client.Events(context.Background(), p)
	if err != nil {
		t.Errorf("Could not retrieve events: %v\n", err)
	}
	if len(events.Events) == 0 {
		t.Errorf("Could not retrieve events: %v\n", err)
	}
}

func TestCall(t *testing.T) {
	_, err := client.Call(context.Background(), types.FunctionCall{
		ContractAddress:    "0x06a09ccb1caaecf3d9683efe335a667b2169a409d19c589ba1eb771cd210af75",
		EntryPointSelector: "decimals",
	}, "0x02ecb1ac7d4925714279245073eb712e13af1263eec175c7917700eafba710b6")
	if err != nil {
		t.Errorf("Could not call contract function: %v\n", err)
	}
}

func TestClientClose(t *testing.T) {
	client.Close()
}
