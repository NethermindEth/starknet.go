package gateway

import (
	"context"
	"log"
	"net/http"
	"testing"
)

func Test_Block_Devnet(t *testing.T) {
    gw := NewClient()

    // I believe setup devnet is the function thats incrementing the block by 1
    err := setupDevnet(context.Background())
		if err != nil {
			log.Fatal("error starting test", err)
		}
    
    // not sure if I should be using gw.Block to check for the current block number

        block, err := gw.Block(context.Background())
    if block.BlockNumber != 1 {
        t.Fatal("block number should be 1")
    }

    // or be making a newRequest for the block?
    // Think I should make a request for the blocknumber and check if its equal to 1
    req, err := gw.newRequest(context.Background(), http.MethodGet, "/get_block")
    if err != nil {
        t.Fatal()
    }

}


func Test_Block(t *testing.T) {
    testConfig := beforeEach(t)

    type testSetType struct {
        BlockHash string
        opts      *BlockOptions
    }
    testSet := map[string][]testSetType{
        "devnet": {
            // testConfig.client.do()
            // TODO/2: check if devnet starts with 1 for instance and we can have a test
            // If not, we would have to create a block
        },
        "testnet": {{
            // TODO/3: instead of testing just the blockHash we should
            // (1) Marshal the block back into a []byte and (2) compare it with the original message
            // check "github.com/nsf/jsondiff" for ideas how to do that.
            BlockHash: "0x57f5102f7c61826926a4d76e544d2272cad091aa4e4b12e8e3e2120a220bd11",
            opts:      &BlockOptions{BlockNumber: 159179}}},
            // TODO/1: add a test for mainnet
		"mainnet": {{
			BlockHash: "0x3bb30a6d1a3b6dcbc935b18c976126ab8d1fea60ef055be3c78530624824d50",
            opts: &BlockOptions{BlockNumber: 5879},
		}},
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
