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

	"github.com/NethermindEth/starknet.go/internal/tests"
	"github.com/stretchr/testify/require"
)

func TestCookieManagement(t *testing.T) {
	tests.RunTestOn(t, tests.MockEnv)

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
					"result":  rpcVersion,
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

	client, err := NewProvider(t.Context(), server.URL)
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
	tests.RunTestOn(t, tests.MockEnv)

	const diffNodeVersion = "0.5.0"

	// Set up a single server that responds differently based on query parameters
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		if method, ok := request["method"].(string); ok && method == "starknet_specVersion" {
			// get node version from query parameter
			nodeVersion := r.URL.Query().Get("nodeVersion")

			if nodeVersion != "" {
				// Return the same version as RPCVersion
				data := map[string]interface{}{
					"jsonrpc": "2.0",
					"id":      request["id"],
					"result":  nodeVersion,
				}
				if err := json.NewEncoder(w).Encode(data); err != nil {
					log.Fatal(err)
				}
			} else {
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
		nodeVersion     string
		expectedWarning string
	}{
		{
			name:            "Compatible version",
			nodeVersion:     rpcVersion.String(),
			expectedWarning: "",
		},
		{
			name:            "Incompatible version",
			nodeVersion:     diffNodeVersion,
			expectedWarning: "warning: the RPC provider version is " + diffNodeVersion + ", and is different from the version " + rpcVersion.String() + " implemented by the SDK. This may cause unexpected behaviour.",
		},
		{
			name:            "Error fetching version",
			nodeVersion:     "",
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
			serverURL := testServer.URL + "?nodeVersion=" + tc.nodeVersion
			provider, err := NewProvider(context.Background(), serverURL)
			require.NoError(t, err)
			require.NotNil(t, provider)

			// Read captured output
			require.NoError(t, w.Close())
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
