package devnet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"github.com/NethermindEth/juno/core/felt"
)

type DevNet struct {
	baseURL   string
	idCounter atomic.Uint64
}

// jsonRPCRequest represents a JSON-RPC request.
type jsonRPCRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      uint64      `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

// jsonRPCResponse represents a JSON-RPC response.
type jsonRPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      uint64          `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *jsonRPCError   `json:"error,omitempty"`
}

// jsonRPCError represents a JSON-RPC error.
type jsonRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *jsonRPCError) Error() string {
	return fmt.Sprintf("JSON-RPC error %d: %s", e.Code, e.Message)
}

// rpcCall makes a JSON-RPC call to the DevNet server.
func (devnet *DevNet) rpcCall(method string, params, result interface{}) error {
	reqBody := jsonRPCRequest{
		JSONRPC: "2.0",
		ID:      devnet.idCounter.Add(1),
		Method:  method,
		Params:  params,
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, devnet.api("/"), bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer closeBody(resp)

	var rpcResp jsonRPCResponse
	if err := json.NewDecoder(resp.Body).Decode(&rpcResp); err != nil {
		return err
	}

	if rpcResp.Error != nil {
		return rpcResp.Error
	}

	if result != nil && len(rpcResp.Result) > 0 {
		return json.Unmarshal(rpcResp.Result, result)
	}

	return nil
}

type TestAccount struct {
	PrivateKey string `json:"private_key"`
	PublicKey  string `json:"public_key"`
	Address    string `json:"address"`
}

// NewDevNet creates a new DevNet instance.
//
// It accepts an optional baseURL parameter, which is a string representing
// the base URL of the DevNet server. If no baseURL is provided, the default
// value of "http://localhost:5050" is used.
//
// Parameters:
//   - baseURL: a string representing the base URL of the DevNet server
//
// Returns:
//   - *DevNet: a pointer to the newly created DevNet instance
func NewDevNet(baseURL ...string) *DevNet {
	if len(baseURL) == 0 {
		return &DevNet{
			baseURL: "http://localhost:5050",
		}
	}

	return &DevNet{
		baseURL: strings.TrimSuffix(baseURL[0], "/"),
	}
}

// api returns the full URL for a given URI.
//
// Parameter:
//   - uri: a string which represents the URI path
//
// Returns:
//   - string which is the full URL constructed using the `devnet.baseURL` and `uri`
func (devnet *DevNet) api(uri string) string {
	uri = strings.TrimPrefix(uri, "/")

	return fmt.Sprintf("%s/%s", devnet.baseURL, uri)
}

// Accounts retrieves a list of test accounts from the DevNet API.
//
// It makes a JSON-RPC call to "devnet_getPredeployedAccounts" and
// decodes the response into a slice of TestAccount structs. It returns
// the list of accounts and any error that occurred during the process.
//
// Parameters:
//
//	none
//
// Returns:
//   - []TestAccount: a slice of TestAccount structs
//   - error: an error if any
func (devnet *DevNet) Accounts() ([]TestAccount, error) {
	var accounts []TestAccount
	err := devnet.rpcCall("devnet_getPredeployedAccounts", nil, &accounts)

	return accounts, err
}

// IsAlive checks if the DevNet is alive.
//
// It sends a GET request to the "/is_alive" endpoint of the DevNet API.
// It returns true if the response status code is 200 (http.StatusOK),
// and false otherwise.
//
// Parameters:
//
//	none
//
// Returns:
//
//	bool: true if the DevNet is alive, false otherwise
func (devnet *DevNet) IsAlive() bool {
	req, err := http.NewRequest(http.MethodGet, devnet.api("/is_alive"), http.NoBody)
	if err != nil {
		return false
	}
	client := &http.Client{
		Timeout: 3 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer closeBody(resp)

	return resp.StatusCode == http.StatusOK
}

type MintResponse struct {
	NewBalance      string `json:"new_balance"`
	Unit            string `json:"unit"`
	TransactionHash string `json:"tx_hash"`
}

// mintParams represents the parameters for the devnet_mint RPC call.
type mintParams struct {
	Address string `json:"address"`
	Amount  uint64 `json:"amount"`
	Unit    string `json:"unit"`
}

// Mint mints a certain amount of tokens for a given address.
//
// It makes a JSON-RPC call to "devnet_mint" to mint tokens.
// The unit defaults to "WEI" for backward compatibility.
//
// Parameters:
//   - address: is the address to mint tokens for
//   - amount: is the amount of tokens to mint
//
// Returns:
//   - *MintResponse: a MintResponse
//   - error: an error if any
func (devnet *DevNet) Mint(address *felt.Felt, amount *big.Int) (*MintResponse, error) {
	params := mintParams{
		Address: address.String(),
		Amount:  amount.Uint64(),
		Unit:    "WEI",
	}

	var mint MintResponse
	err := devnet.rpcCall("devnet_mint", params, &mint)

	return &mint, err
}

// closeBody closes the response body and logs any errors.
func closeBody(resp *http.Response) {
	if err := resp.Body.Close(); err != nil {
		log.Printf("Error closing response body: %v", err)
	}
}
