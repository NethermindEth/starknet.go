package rpc

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

const (
	TestPublicKey            = "0x783318b2cc1067e5c06d374d2bb9a0382c39aabd009b165d7a268b882971d6"
	DevNetETHAddress         = "0x49d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"
	TestNetETHAddress        = "0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"
	DevNetAccount032Address  = "0x06bb9425718d801fd06f144abb82eced725f0e81db61d2f9f4c9a26ece46a829"
	TestNetAccount032Address = "0x6ca4fdd437dffde5253ba7021ef7265c88b07789aa642eafda37791626edf00"
	DevNetAccount040Address  = "0x080dff79c6216ad300b872b73ff41e271c63f213f8a9dc2017b164befa53b9"
	TestNetAccount040Address = "0x6cbfa37f409610fee26eeb427ed854b3a4b24580d9b9ef6c3e38db7b3f7322c"
)

// testConfiguration is a type that is used to configure tests
type testConfiguration struct {
	provider *Provider
	base     string
}

var (
	// set the environment for the test, default: mock
	testEnv = "mock"

	// testConfigurations are predefined test configurations
	testConfigurations = map[string]testConfiguration{
		// Requires a Mainnet Starknet JSON-RPC compliant node (e.g. pathfinder)
		// (ref: https://github.com/eqlabs/pathfinder)
		"mainnet": {
			base: "https://free-rpc.nethermind.io/mainnet-juno",
		},
		// Requires a Testnet Starknet JSON-RPC compliant node (e.g. pathfinder)
		// (ref: https://github.com/eqlabs/pathfinder)
		"testnet": {
			base: "https://free-rpc.nethermind.io/sepolia-juno",
		},
		// Requires a Devnet configuration running locally
		// (ref: https://github.com/0xSpaceShard/starknet-devnet-rs)
		"devnet": {
			base: "http://localhost:5050/",
		},
		// Used with a mock as a standard configuration, see `mock_test.go``
		"mock":        {},
		"integration": {},
	}
)

// TestMain is a Go function that serves as the entry point for running tests.
//
// It takes a pointer to the testing.M struct as its parameter and returns nothing.
// The purpose of this function is to set up any necessary test environment
// variables before running the tests and to clean up any resources afterwards.
// It also parses command line flags and exits with the exit code returned by
// the testing.M.Run() function.
//
// Parameters:
// - m: the testing.M struct
// Returns:
//
//	none
func TestMain(m *testing.M) {
	flag.StringVar(&testEnv, "env", "mock", "set the test environment")
	flag.Parse()

	os.Exit(m.Run())
}

// beforeEach initializes the test environment configuration before running the script.
//
// Parameters:
// - t: The testing.T object for testing purposes
// Returns:
// - *testConfiguration: a pointer to the testConfiguration struct
func beforeEach(t *testing.T) *testConfiguration {
	t.Helper()
	_ = godotenv.Load(fmt.Sprintf(".env.%s", testEnv), ".env")
	testConfig, ok := testConfigurations[testEnv]
	if !ok {
		t.Fatal("env supports mock, testnet, mainnet, devnet, integration")
	}
	if testEnv == "mock" {
		testConfig.provider = &Provider{
			c: &rpcMock{},
		}
		return &testConfig
	}

	base := os.Getenv("INTEGRATION_BASE")
	if base != "" {
		testConfig.base = base
	}
	c, err := NewProvider(testConfig.base)
	if err != nil {
		t.Fatal("connect should succeed, instead:", err)
	}

	testConfig.provider = c
	t.Cleanup(func() {
		testConfig.provider.c.Close()
	})
	return &testConfig
}

// TestChainID is a function that tests the ChainID function in the Go test file.
//
// The function initializes a test configuration and defines a test set with different chain IDs for different environments.
// It then iterates over the test set and for each test, creates a new spy and sets the spy as the provider's client.
// The function calls the ChainID function and compares the returned chain ID with the expected chain ID.
// If there is a mismatch or an error occurs, the function logs a fatal error.
//
// Parameters:
// - t: the testing object for running the test cases
// Returns:
//
//	none
func TestChainID(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		ChainID string
	}
	testSet := map[string][]testSetType{
		"devnet":  {{ChainID: "SN_SEPOLIA"}},
		"mainnet": {{ChainID: "SN_MAIN"}},
		"mock":    {{ChainID: "SN_SEPOLIA"}},
		"testnet": {{ChainID: "SN_SEPOLIA"}},
	}[testEnv]

	for _, test := range testSet {
		spy := NewSpy(testConfig.provider.c)
		testConfig.provider.c = spy
		chain, err := testConfig.provider.ChainID(context.Background())
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

// TestSyncing tests the syncing functionality.
//
// It initializes a test configuration and sets up a test set. Then it loops
// through the test set and creates a spy object. It calls the Syncing function
// of the provider using the test configuration. It checks if there is any
// error during syncing, and if so, it fails the test. If the starting block
// hash is not nil, it compares the sync object with the spy object. It checks
// if the current block number is a positive number and if the current block
// hash starts with "0x". If the starting block hash is nil, it compares the
// sync object with the spy object and checks if the current block hash is nil.
//
// Parameters:
// - t: the testing object for running the test cases
// Returns:
//
//	none
func TestSyncing(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		ChainID string
	}

	testSet := map[string][]testSetType{
		"devnet":  {},
		"mainnet": {{ChainID: "SN_MAIN"}},
		"mock":    {{ChainID: "MOCK"}},
		"testnet": {{ChainID: "SN_SEPOLIA"}},
	}[testEnv]

	for range testSet {
		spy := NewSpy(testConfig.provider.c)
		testConfig.provider.c = spy
		sync, err := testConfig.provider.Syncing(context.Background())
		if err != nil {
			t.Fatal("Syncing error:", err)
		}
		if sync.StartingBlockHash != nil {
			if diff, err := spy.Compare(sync, false); err != nil || diff != "FullMatch" {
				if _, err := spy.Compare(sync, true); err != nil {
					log.Fatal(err)
				}
				t.Fatal("expecting to match", err)
			}
			i, ok := big.NewInt(0).SetString(string(sync.CurrentBlockNum), 0)
			if !ok || i.Cmp(big.NewInt(0)) <= 0 {
				t.Fatal("CurrentBlockNum should be positive number, instead: ", sync.CurrentBlockNum)
			}
			if !strings.HasPrefix(sync.CurrentBlockHash.String(), "0x") {
				t.Fatal("current block hash should return a string starting with 0x")
			}
		} else {
			if _, err := spy.Compare(sync, false); err != nil {
				log.Fatal(err)
			}
			require.Nil(t, sync.CurrentBlockHash)

		}

	}
}

func TestGetBlock(t *testing.T) {
	testConfig := beforeEach(t)
	type testSetType struct {
		BlockID      BlockID
		ExpectedResp *Block
		ExpectedErr  *RPCError
	}

	testSet := map[string][]testSetType{
		"devnet": {},
		"mock": {
			{
				BlockID: BlockID{Tag: "latest"},
				ExpectedResp: &Block{
					BlockHeader: BlockHeader{
						L1DAMode: L1DAModeBlob,
						L1DataGasPrice: ResourcePrice{
							PriceInWei: new(felt.Felt).SetUint64(1),
							PriceInFRI: new(felt.Felt).SetUint64(1),
						},
						L1GasPrice: ResourcePrice{
							PriceInWei: new(felt.Felt).SetUint64(1),
							PriceInFRI: new(felt.Felt).SetUint64(1),
						},
					},
				},
				ExpectedErr: nil,
			},
		},
	}[testEnv]

	for _, test := range testSet {
		block, err := testConfig.provider.BlockWithTxHashes(context.Background(), BlockID{Tag: "latest"})
		if test.ExpectedErr != nil {
			require.Equal(t, test.ExpectedErr, err)
		} else {
			blockCasted := block.(*BlockTxHashes)
			expectdBlockHeader := test.ExpectedResp.BlockHeader
			require.Equal(t, blockCasted.L1DAMode, expectdBlockHeader.L1DAMode)
			require.Equal(t, blockCasted.L1DataGasPrice.PriceInWei, expectdBlockHeader.L1DataGasPrice.PriceInWei)
			require.Equal(t, blockCasted.L1DataGasPrice.PriceInFRI, expectdBlockHeader.L1DataGasPrice.PriceInFRI)
			require.Equal(t, blockCasted.L1GasPrice.PriceInWei, expectdBlockHeader.L1GasPrice.PriceInWei)
			require.Equal(t, blockCasted.L1GasPrice.PriceInFRI, expectdBlockHeader.L1GasPrice.PriceInFRI)
		}

	}
}
func TestCookieManagement(t *testing.T) {
	// Don't return anything unless cookie is set.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := r.Cookie("session_id"); err == http.ErrNoCookie {
			http.SetCookie(w, &http.Cookie{
				Name:  "session_id",
				Value: "12345",
				Path:  "/",
			})
		} else {
			var result string
			err := mock_starknet_chainId(&result, "")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			data := map[string]interface{}{
				"jsonrpc": "2.0",
				"id":      1,
				"result":  result,
			}
			if err := json.NewEncoder(w).Encode(data); err != nil {
				log.Fatal(err)
			}
		}
	}))
	defer server.Close()

	client, err := NewProvider(server.URL)
	require.Nil(t, err)

	resp, err := client.ChainID(context.Background())
	require.NotNil(t, err)
	require.Equal(t, resp, "")

	resp, err = client.ChainID(context.Background())
	require.Nil(t, err)
	require.Equal(t, resp, "SN_SEPOLIA")

	resp, err = client.ChainID(context.Background())
	require.Nil(t, err)
	require.Equal(t, resp, "SN_SEPOLIA")
}
