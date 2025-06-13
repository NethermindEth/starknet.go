package rpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
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

// beforeEach initialises the test environment configuration before running the script.
//
// Parameters:
//   - t: The testing.T object for testing purposes
//   - isWs: a boolean value to check if the test is for the websocket provider
//
// Returns:
//   - *testConfiguration: a pointer to the testConfiguration struct
func beforeEach(t *testing.T, isWs bool) *testConfiguration {
	t.Helper()

	var testConfig testConfiguration

	if internal.TEST_ENV == internal.MockEnv {
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

	if internal.TEST_ENV == internal.DevnetEnv || internal.TEST_ENV == internal.MainnetEnv {
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
		// Handle version check request
		if r.Method == http.MethodPost {
			var request map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)

				fmt.Println("Error decoding request body:", err)

				return
			}

			// Check if this is a version request
			if method, ok := request["method"].(string); ok && method == "starknet_specVersion" {
				data := map[string]interface{}{
					"jsonrpc": "2.0",
					"id":      request["id"],
					"result":  "0.8.0",
				}
				if err := json.NewEncoder(w).Encode(data); err != nil {
					log.Fatal(err)
				}

				return
			}
		}

		// Handle cookie management
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

// TestVersionCompatibility tests that the provider correctly handles version compatibility warnings
func TestVersionCompatibility(t *testing.T) {
	const wrongVersion = "0.5.0"

	// Set up a single server that responds differently based on query parameters
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		if method, ok := request["method"].(string); ok && method == "starknet_specVersion" {
			// Get test case from query parameter
			testCase := r.URL.Query().Get("testCase")

			switch testCase {
			case "compatible":
				// Return the same version as RPCVersion
				data := map[string]interface{}{
					"jsonrpc": "2.0",
					"id":      request["id"],
					"result":  rpcVersion,
				}
				if err := json.NewEncoder(w).Encode(data); err != nil {
					log.Fatal(err)
				}
			case "incompatible":
				// Return a different version
				data := map[string]interface{}{
					"jsonrpc": "2.0",
					"id":      request["id"],
					"result":  wrongVersion, // Different version
				}
				if err := json.NewEncoder(w).Encode(data); err != nil {
					log.Fatal(err)
				}
			case "error":
				// Return an error
				data := map[string]interface{}{
					"jsonrpc": "2.0",
					"id":      request["id"],
					"error": map[string]interface{}{
						"code":    -32601,
						"message": "Method not found",
					},
				}
				if err := json.NewEncoder(w).Encode(data); err != nil {
					log.Fatal(err)
				}
			}
		}
	}))
	defer testServer.Close()

	// Test cases
	testCases := []struct {
		name            string
		queryParam      string
		expectedWarning string
	}{
		{
			name:            "Compatible version",
			queryParam:      "compatible",
			expectedWarning: "",
		},
		{
			name:            "Incompatible version",
			queryParam:      "incompatible",
			expectedWarning: fmt.Sprintf(warnVersionMismatch, rpcVersion, wrongVersion),
		},
		{
			name:            "Error fetching version",
			queryParam:      "error",
			expectedWarning: warnVersionCheckFailed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Capture stdout
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Create provider with query parameter - this will trigger the version check
			serverURL := testServer.URL + "?testCase=" + tc.queryParam
			provider, err := NewProvider(serverURL)
			require.NoError(t, err)
			require.NotNil(t, provider)

			// Read captured output
			w.Close()
			os.Stdout = old
			var buf bytes.Buffer
			_, err = io.Copy(&buf, r)
			require.NoError(t, err, "Failed to read from pipe")
			output := buf.String()

			// Check if warning is present as expected
			if tc.expectedWarning == "" {
				require.Empty(t, output, "Expected no warning")
			} else {
				require.Contains(t, output, tc.expectedWarning, "Expected warning not found")
			}
		})
	}
}
