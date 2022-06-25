package rpc

import (
	"context"
	_ "embed"
	"encoding/hex"
	"math/big"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

// testConfiguration is a type that is used to configure tests
type testConfiguration struct {
	client  *Client
	base    string
	chainid string
}

// testConfiguration are predefined test configurations
var testEnvironments = map[string]testConfiguration{
	// Requires a Mainnet StarkNet JSON-RPC compliant node (e.g. pathfinder)
	// (ref: https://github.com/eqlabs/pathfinder)
	"mainnet_pathfinder": {
		base:    "http://localhost:9545",
		chainid: "SN_MAIN",
	},
	// Requires a Testnet StarkNet JSON-RPC compliant node (e.g. pathfinder)
	// (ref: https://github.com/eqlabs/pathfinder)
	"testnet_pathfinder": {
		base:    "http://localhost:9545",
		chainid: "SN_GOERLI",
	},
	// Requires a Devnet configuration running locally
	// (ref: https://github.com/Shard-Labs/starknet-devnet)
	"devnet": {
		base:    "http://localhost:5050/rpc",
		chainid: "DEVNET",
	},
	// Used with a mock as a standard configuration, see `mock_test.go``
	"mock": {
		chainid: "MOCK",
	},
}

// beforeEach checks the configuration and initializes it before running the script
func beforeEach(t *testing.T) *testConfiguration {
	godotenv.Load()
	integration := os.Getenv("INTEGRATION")
	if integration == "" {
		integration = "mock"
	}
	testConfig, ok := testEnvironments[integration]
	if !ok {
		t.Fatal("INTEGRATION supports testnet_pathfinder, mainnet_pathfinder or devnet")
	}
	if integration == "mock" {
		testConfig.client = &Client{
			c: &rpcMock{},
		}
		return &testConfig
	}
	client, err := DialContext(context.Background(), testConfig.base)
	if err != nil {
		t.Fatal("connect should succeed, instead:", err)
	}
	testConfig.client = client
	return &testConfig
}

func TestIntegrationChainID(t *testing.T) {
	testConfig := beforeEach(t)
	defer testConfig.client.Close()

	chain, err := testConfig.client.ChainID(context.Background())

	if err != nil {
		t.Fatal(err)
	}
	chainInt, ok := big.NewInt(0).SetString(chain, 0)
	if !ok {
		t.Fatal("could not load str representation of an int")
	}
	chainID, err := hex.DecodeString(chainInt.Text(16))
	if err != nil {
		t.Fatal(err)
	}
	if string(chainID) != testConfig.chainid {
		t.Fatalf("expecting %s, instead: %s", testConfig.chainid, string(chainID))
	}
}

func TestIntegrationSyncing(t *testing.T) {
	testConfig := beforeEach(t)
	defer testConfig.client.Close()

	sync, err := testConfig.client.Syncing(context.Background())

	if err != nil {
		t.Fatal(err)
	}
	if sync == nil || sync.CurrentBlockNum == "" {
		t.Fatal("should succeed")
	}
	i, ok := big.NewInt(0).SetString(sync.CurrentBlockNum, 0)
	if !ok || i.Cmp(big.NewInt(0)) <= 0 {
		t.Fatal("returned value should be a number, instead: ", sync.CurrentBlockNum)
	}
}
