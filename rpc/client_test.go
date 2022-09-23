package rpc

import (
	"context"
	"flag"
	"fmt"
	"math/big"
	"os"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/joho/godotenv"
)

const (
	TestPublicKey         = "0x783318b2cc1067e5c06d374d2bb9a0382c39aabd009b165d7a268b882971d6"
	DevNetETHAddress      = "0x62230ea046a9a5fbc261ac77d03c8d41e5d442db2284587570ab46455fd2488"
	DevNetAccountAddress  = "0x06bb9425718d801fd06f144abb82eced725f0e81db61d2f9f4c9a26ece46a829"
	TestNetAccountAddress = "0x19e63006d7df131737f5222283da28de2d9e2f0ee92fdc4c4c712d1659826b0"
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

// TestMain is used to trigger the tests and, in that case, check for the environment to use.
func TestMain(m *testing.M) {
	flag.StringVar(&testEnv, "env", "mock", "set the test environment")
	flag.Parse()

	os.Exit(m.Run())
}

// beforeEach checks the configuration and initializes it before running the script
func beforeEach(t *testing.T) *testConfiguration {
	t.Helper()
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

	testConfig.base = "https://starknet-goerli.cartridge.gg"
	base := os.Getenv("INTEGRATION_BASE")
	if base != "" {
		testConfig.base = base
	}
	client, err := DialContext(context.Background(), testConfig.base)
	if err != nil {
		t.Fatal("connect should succeed, instead:", err)
	}
	testConfig.client = client
	t.Cleanup(func() {
		testConfig.client.Close()
	})
	return &testConfig
}

// TestChainID checks the chainId matches the one for the environment
func TestChainID(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		ChainID string
	}
	testSet := map[string][]testSetType{
		"devnet":  {{ChainID: "SN_GOERLI"}},
		"mainnet": {{ChainID: "SN_MAIN"}},
		"mock":    {{ChainID: "MOCK"}},
		"testnet": {{ChainID: "SN_GOERLI"}},
	}[testEnv]

	fmt.Printf("----------------------------\n")
	fmt.Printf("Env: %s\n", testEnv)
	fmt.Printf("Url: %s\n", testConfig.base)
	fmt.Printf("----------------------------\n")

	for _, test := range testSet {
		spy := NewSpy(testConfig.client.c)
		testConfig.client.c = spy
		chain, err := testConfig.client.ChainID(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		if _, err := spy.Compare(chain, false); err != nil {
			t.Fatal("expecting to match", err)
		}
		if chain != test.ChainID {
			t.Fatalf("expecting %s, instead: %s", test.ChainID, chain)
		}
	}
}

// TestSyncing checks the values returned are consistent
func TestSyncing(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		ChainID string
	}

	testSet := map[string][]testSetType{
		"devnet":  {},
		"mainnet": {{ChainID: "SN_MAIN"}},
		"mock":    {{ChainID: "MOCK"}},
		"testnet": {{ChainID: "SN_GOERLI"}},
	}[testEnv]

	for range testSet {
		spy := NewSpy(testConfig.client.c)
		testConfig.client.c = spy
		sync, err := testConfig.client.Syncing(context.Background())
		if err != nil {
			t.Fatal("BlockWithTxHashes match the expected error:", err)
		}
		if diff, err := spy.Compare(sync, false); err != nil || diff != "FullMatch" {
			spy.Compare(sync, true)
			t.Fatal("expecting to match", err)
		}
		i, ok := big.NewInt(0).SetString(sync.CurrentBlockNum, 0)
		if !ok || i.Cmp(big.NewInt(0)) <= 0 {
			t.Fatal("CurrentBlockNum should be positive number, instead: ", sync.CurrentBlockNum)
		}
		if !strings.HasPrefix(sync.CurrentBlockHash, "0x") {
			t.Fatal("current block hash should return a string starting with 0x")
		}
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
