package rpc

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/NethermindEth/starknet.go/internal"
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

// the environment for the test, default: mock
var testEnv = ""

// TestMain is used to trigger the tests and set up the test environment.
//
// It sets up the test environment by loading environment variables using the internal.LoadEnv() function,
// which handles parsing command line flags and loading the appropriate .env files based on the
// specified environment (mock, testnet, mainnet, or devnet).
// After setting up the environment, it runs the tests and exits with the return value of the test suite.
//
// Parameters:
//   - m: The testing.M object that provides the entry point for running tests
// Returns:
//
//	none
func TestMain(m *testing.M) {
	testEnv = internal.LoadEnv()

	os.Exit(m.Run())
}

// beforeEach initializes the test environment configuration before running the script.
//
// Parameters:
//   - t: The testing.T object for testing purposes
//   - isWs: a boolean value to check if the test is for the websocket provider
// Returns:
//   - *testConfiguration: a pointer to the testConfiguration struct
func beforeEach(t *testing.T, isWs bool) *testConfiguration {
	t.Helper()

	var testConfig testConfiguration

	if testEnv == "mock" {
		testConfig.provider = &Provider{
			c: &rpcMock{},
		}
		return &testConfig
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

	if testEnv == "devnet" || testEnv == "mainnet" {
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
