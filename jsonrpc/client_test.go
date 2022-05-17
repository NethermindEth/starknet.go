package jsonrpc

import (
	"context"
	"math/big"
	"testing"
)

// Requires a StarkNet JSON-RPC compliant node (e.g. pathfinder)
// (ref: https://github.com/eqlabs/pathfinder)
func TestJsonRpcClient(t *testing.T) {
	client, err := DialContext(context.Background(), "http://localhost:9545")
	if err != nil {
		t.Errorf("Could not connect to local StarkNet node: %v\n", err)
	}
	defer client.Close()

	// _, err = client.BlockByHash(context.Background(), "0x14b05305c69bcfa91945cd2a1a0cd4d9e8879b96e57a1688843a0719afce7c2", "TXN_HASH")
	// if err != nil {
	// 	t.Errorf("Could not retrieve block: %v\n", err)
	// }

	_, err = client.BlockByHash(context.Background(), "0x14b05305c69bcfa91945cd2a1a0cd4d9e8879b96e57a1688843a0719afce7c2", "FULL_TXNS")
	if err != nil {
		t.Errorf("Could not retrieve block: %v\n", err)
	}

	_, err = client.BlockByHash(context.Background(), "0x14b05305c69bcfa91945cd2a1a0cd4d9e8879b96e57a1688843a0719afce7c2", "FULL_TXN_AND_RECEIPTS")
	if err != nil {
		t.Errorf("Could not retrieve block: %v\n", err)
	}

	_, err = client.BlockByNumber(context.Background(), big.NewInt(1), "FULL_TXNS")
	if err != nil {
		t.Errorf("Could not retrieve block: %v\n", err)
	}

	_, err = client.BlockByNumber(context.Background(), big.NewInt(1), "FULL_TXN_AND_RECEIPTS")
	if err != nil {
		t.Errorf("Could not retrieve block: %v\n", err)
	}

	_, err = client.BlockNumber(context.Background())
	if err != nil {
		t.Errorf("Could not retrieve block: %v\n", err)
	}

	_, err = client.CodeAt(context.Background(), "0x050c47150563ff7cf60dd60f7d0bd4d62a9cc5331441916e5099e905bdd8c4bc")
	if err != nil {
		t.Errorf("Could not retrieve code: %v\n", err)
	}

	_, err = client.TransactionReceipt(context.Background(), "0x6509ffcbec029528df550c19e3ddde62ca7b23eeb4967ed183a6a31ed106dfc")
	if err != nil {
		t.Errorf("Could not retrieve block: %v\n", err)
	}

	_, err = client.Call(context.Background(), FunctionCall{
		ContractAddress:    "0x06a09ccb1caaecf3d9683efe335a667b2169a409d19c589ba1eb771cd210af75",
		EntryPointSelector: "decimals",
	}, "0x02ecb1ac7d4925714279245073eb712e13af1263eec175c7917700eafba710b6")
	if err != nil {
		t.Errorf("Could not call contract function: %v\n", err)
	}
}
