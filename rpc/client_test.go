package rpc

import (
	"context"
	"testing"

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
