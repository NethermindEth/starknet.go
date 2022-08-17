package gateway

import (
	"context"
	"testing"
)

func Test_Block(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		BlockHash string
		opts *BlockOptions 
		}
	testSet := map[string][]testSetType{
		"devnet":  {},
		"testnet": {{
			BlockHash: "0x57f5102f7c61826926a4d76e544d2272cad091aa4e4b12e8e3e2120a220bd11",
			opts: &BlockOptions{BlockNumber: 159179}}},
	}[testEnv]

	for _, test := range testSet {
		block, err := testConfig.client.Block(context.Background(), test.opts)

		if err != nil {
			t.Fatal(err)
		}
		if block.BlockHash != test.BlockHash {
			t.Fatalf("expecting %s, instead: %s", block.BlockHash, test.BlockHash)
		}
	}
}

func Test_BlockIDByHash(t *testing.T) {
	gw := NewClient()

	id, err := gw.BlockIDByHash(context.Background(), "0x5239614da0a08b53fa8cbdbdcb2d852e153027ae26a2ae3d43f7ceceb28551e")
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
	if err != nil || id == "" {
		t.Errorf("Getting Block ID by Hash: %v", err)
	}

	if id != "0x5239614da0a08b53fa8cbdbdcb2d852e153027ae26a2ae3d43f7ceceb28551e" {
		t.Errorf("Wrong Block ID from Hash: %v", err)
	}
}
