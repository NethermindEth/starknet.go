package rpc

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

const (
	DevNetETHAddress = "0x49d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"
)

// testConfiguration is a type that is used to configure tests
type testConfiguration struct {
	provider   *Provider
	wsProvider *WsProvider
	base       string
	wsBase     string
}

var (
	// set the environment for the test, default: mock
	testEnv = "mock"

	// testConfigurations are predefined test configurations
	testConfigurations = map[string]testConfiguration{
		// Requires a Mainnet Starknet JSON-RPC compliant node (e.g. Juno)
		// (ref: https://github.com/NethermindEth/juno)
		"mainnet": {
			base: "https://free-rpc.nethermind.io/mainnet-juno",
		},
		// Requires a Testnet Starknet JSON-RPC compliant node (e.g. Juno)
		// (ref: https://github.com/NethermindEth/juno)
		"testnet": {
			base:   "https://free-rpc.nethermind.io/sepolia-juno",
			wsBase: "ws://localhost:6061",
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
// - isWs: a boolean value to check if the test is for the websocket provider
// Returns:
// - *testConfiguration: a pointer to the testConfiguration struct
func beforeEach(t *testing.T, isWs bool) *testConfiguration {
	t.Helper()
	_ = godotenv.Load(fmt.Sprintf(".env.%s", testEnv), ".env")
	testConfig, ok := testConfigurations[testEnv]
	if !ok {
		t.Fatal("env supports mock, testnet, wsTestnet, mainnet, devnet, integration")
	}

	if testEnv == "mock" {
		testConfig.provider = &Provider{
			c: &rpcMock{},
		}
	}

	base := os.Getenv("HTTP_PROVIDER_URL")
	if base != "" {
		testConfig.base = base
	}

	client, err := NewProvider(testConfig.base)
	if err != nil {
		t.Fatalf("failed to connect to the %s provider: %v", testConfig.base, err)
	}
	testConfig.provider = client
	t.Cleanup(func() {
		testConfig.provider.c.Close()
	})

	if testEnv == "devnet" {
		return &testConfig
	}

	if isWs {
		wsBase := os.Getenv("WS_PROVIDER_URL")
		if wsBase != "" {
			testConfig.wsBase = wsBase

		}

		wsClient, err := NewWebsocketProvider(testConfig.wsBase)
		if err != nil {
			t.Fatalf("failed to connect to the %s websocket provider: %v", testConfig.wsBase, err)
		}
		testConfig.wsProvider = wsClient
		t.Cleanup(func() {
			testConfig.wsProvider.c.Close()
		})
	}

	return &testConfig
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
			var rawResp json.RawMessage
			err := mock_starknet_chainId(&rawResp)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			var result string
			if err := json.Unmarshal(rawResp, &result); err != nil {
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
