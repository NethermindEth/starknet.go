package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NethermindEth/starknet.go/internal/tests"
	"github.com/stretchr/testify/assert"
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
			data := map[string]interface{}{
				"jsonrpc": "2.0",
				"id":      1,
				"result":  "0x534e5f5345504f4c4941",
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

			data := map[string]interface{}{
				"jsonrpc": "2.0",
				"id":      request["id"],
				"result":  nodeVersion,
			}
			if err := json.NewEncoder(w).Encode(data); err != nil {
				log.Fatal(err)
			}
		}
	}))
	defer testServer.Close()

	// Test cases
	testCases := []struct {
		name        string
		nodeVersion string
		expectedErr error
	}{
		{
			name:        "Compatible version",
			nodeVersion: rpcVersion.String(),
			expectedErr: nil,
		},
		{
			name:        "Incompatible version",
			nodeVersion: diffNodeVersion,
			expectedErr: errors.Join(
				ErrIncompatibleVersion,
				fmt.Errorf("expected version: %s, got: %s", rpcVersion, diffNodeVersion),
			),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create provider with query parameter - this will trigger the version check
			serverURL := testServer.URL + "?nodeVersion=" + tc.nodeVersion
			provider, err := NewProvider(context.Background(), serverURL)

			if tc.expectedErr == nil {
				assert.NoError(t, err)
				assert.NotNil(t, provider)
			} else {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedErr.Error())
				assert.NotNil(t, provider)
			}
		})
	}
}
