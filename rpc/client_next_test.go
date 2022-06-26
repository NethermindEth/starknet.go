package rpc

import (
	"context"
	"encoding/hex"
	"flag"
	"fmt"
	"math/big"
	"os"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/joho/godotenv"
)

// testConfiguration is a type that is used to configure tests
type testConfiguration struct {
	client *Client
	base   string
}

var (
	// set the environment for the test, default: mock
	testEnv = "mock"

	// testConfigurations are predefined test configurations
	testConfigurations = map[string]testConfiguration{
		// Requires a Mainnet StarkNet JSON-RPC compliant node (e.g. pathfinder)
		// (ref: https://github.com/eqlabs/pathfinder)
		"mainnet": {
			base: "http://localhost:9545",
		},
		// Requires a Testnet StarkNet JSON-RPC compliant node (e.g. pathfinder)
		// (ref: https://github.com/eqlabs/pathfinder)
		"testnet": {
			base: "http://localhost:9545",
		},
		// Requires a Devnet configuration running locally
		// (ref: https://github.com/Shard-Labs/starknet-devnet)
		"devnet": {
			base: "http://localhost:5050/rpc",
		},
		// Used with a mock as a standard configuration, see `mock_test.go``
		"mock": {},
	}
)

func TestMain(m *testing.M) {
	flag.StringVar(&testEnv, "env", "mock", "set the test environment")
	flag.Parse()

	os.Exit(m.Run())
}

// beforeEach checks the configuration and initializes it before running the script
func beforeEach(t *testing.T) *testConfiguration {
	godotenv.Load(fmt.Sprintf(".env.%s", testEnv), ".env")
	testConfig, ok := testConfigurations[testEnv]
	if !ok {
		t.Fatal("env supports mock, testnet, mainnet or devnet")
	}
	if testEnv == "mock" {
		testConfig.client = &Client{
			c: &rpcMock{},
		}
		return &testConfig
	}
	base := os.Getenv("INTEGRATION_BASE")
	if base != "" {
		testConfig.base = base
	}
	client, err := DialContext(context.Background(), testConfig.base)
	if err != nil {
		t.Fatal("connect should succeed, instead:", err)
	}
	testConfig.client = client
	return &testConfig
}

// TestChainID checks the chainId matches the one for the environment
func TestChainID(t *testing.T) {
	testConfig := beforeEach(t)
	defer testConfig.client.Close()

	type testSetType struct {
		ChainID string
	}
	testSet := map[string][]testSetType{
		"devnet":  {{ChainID: "DEVNET"}},
		"mainnet": {{ChainID: "SN_MAIN"}},
		"mock":    {{ChainID: "MOCK"}},
		"testnet": {{ChainID: "SN_GOERLI"}},
	}[testEnv]

	fmt.Printf("----------------------------\n")
	fmt.Printf("Env: %s\n", testEnv)
	fmt.Printf("Url: %s\n", testConfig.base)
	fmt.Printf("----------------------------\n")

	for _, test := range testSet {
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
		if string(chainID) != test.ChainID {
			t.Fatalf("expecting %s, instead: %s", test.ChainID, string(chainID))
		}
	}
}

// TestSyncing checks the values returned are consistent
func TestSyncing(t *testing.T) {
	testConfig := beforeEach(t)
	defer testConfig.client.Close()

	sync, err := testConfig.client.Syncing(context.Background())

	if err != nil {
		t.Fatal(err)
	}
	if sync == nil || sync.CurrentBlockNum == "" {
		t.Fatal("should succeed")
	}
	if !strings.HasPrefix(sync.CurrentBlockNum, "0x") {
		t.Fatal("CurrentBlockNum should start with 0x, instead:", sync.CurrentBlockHash)
	}
	i, ok := big.NewInt(0).SetString(sync.CurrentBlockNum, 0)
	if !ok || i.Cmp(big.NewInt(0)) <= 0 {
		t.Fatal("CurrentBlockNum should be positive number, instead: ", sync.CurrentBlockNum)
	}
}

// TestClose checks the function is called
func TestClose(t *testing.T) {
	testConfig := beforeEach(t)

	testConfig.client.Close()

	switch client := testConfig.client.c.(type) {
	case *rpc.Client:
		return
	case *rpcMock:
		if client.closed {
			return
		}
		t.Fatalf("client should have been closed")
	default:
		t.Fatalf("client unsupported type %T", testConfig.client.c)
	}
}
