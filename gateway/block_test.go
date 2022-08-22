package gateway

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/dontpanicdao/caigo/felt"
	"github.com/google/go-querystring/query"
)

func Test_BlockIDByHash(t *testing.T) {
	gw := NewClient()

	id, err := gw.BlockIDByHash(context.Background(), felt.StrToFelt("0x5239614da0a08b53fa8cbdbdcb2d852e153027ae26a2ae3d43f7ceceb28551e"))
	if err != nil || id == 0 {
		t.Errorf("Getting Block ID by Hash: %v", err)
	}

	if id != 159179 {
		t.Errorf("Wrong Block ID from Hash: %v", err)
	}
}

func Test_BlockHashByID(t *testing.T) {
	gw := NewClient()

	id, err := gw.BlockHashByID(context.Background(), 159179)
	if err != nil {
		t.Errorf("Getting Block ID by Hash: %v", err)
	}

	if id.String() != "0x5239614da0a08b53fa8cbdbdcb2d852e153027ae26a2ae3d43f7ceceb28551e" {
		t.Errorf("Wrong Block ID from Hash: %v", err)
	}
}

func TestValueWithFeltWithPrepare(t *testing.T) {
	blockHash := felt.BigToFelt(big.NewInt(1))
	v := &BlockOptions{
		BlockNumber: 0,
		BlockHash:   blockHash,
	}

	type tempBlockOptions struct {
		BlockNumber uint64 `url:"blockNumber,omitempty"`
		BlockHash   string `url:"blockHash,omitempty"`
	}
	out := tempBlockOptions{
		BlockNumber: v.BlockNumber,
	}
	if !v.BlockHash.IsNil() {
		out.BlockHash = fmt.Sprintf("0x%s", v.BlockHash.Int.Text(16))
	}
	output, err := query.Values(out)
	if err != nil {
		t.Error(err)
	}
	x := output.Get("blockHash")
	if x != "0x1" {
		t.Errorf("Blockhash should be 1 (or 0x1), instead: \"%s\"", x)
	}
}

// TestGateway checks the gateway can be accessed.
func TestBlock(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		BlockHash felt.Felt
	}
	testSet := map[string][]testSetType{
		"devnet":  {},
		"mainnet": {{BlockHash: felt.StrToFelt("0x4ee4c886d1767b7165a1e3a7c6ad145543988465f2bda680c16a79536f6d81f")}},
		"mock":    {{BlockHash: felt.StrToFelt("0xdeadbeef")}},
		"testnet": {{BlockHash: felt.StrToFelt("0x787af09f1cacdc5de1df83e8cdca3a48c1194171c742e78a9f684cb7aa4db")}},
	}[testEnv]

	for _, test := range testSet {
		block, err := testConfig.client.Block(context.Background(), &BlockOptions{BlockHash: test.BlockHash})

		if err != nil {
			t.Fatal(err)
		}

		if block == nil || block.BlockHash.Cmp(test.BlockHash.Int) != 0 {
			t.Fatalf("expecting %v, instead: %v", test.BlockHash, block.BlockHash)
		}
	}
}
