package devnet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/NethermindEth/juno/core/felt"
)

type DevNet struct {
	baseURL string
}

type TestAccount struct {
	PrivateKey string `json:"private_key"`
	PublicKey  string `json:"public_key"`
	Address    string `json:"address"`
}

// NewDevNet creates a new DevNet instance.
//
// It accepts an optional baseURL parameter, which is a string representing the base URL of the DevNet server.
// If no baseURL is provided, the default value of "http://localhost:5050" is used.
//
// Parameters:
// - baseURL: a string representing the base URL of the DevNet server
// Returns:
// - *DevNet: a pointer to the newly created DevNet instance
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
// - uri: a string which represents the URI path
// Returns:
// - string which is the full URL constructed using the `devnet.baseURL` and `uri`
func (devnet *DevNet) api(uri string) string {
	uri = strings.TrimPrefix(uri, "/")
	return fmt.Sprintf("%s/%s", devnet.baseURL, uri)
}

// Accounts retrieves a list of test accounts from the DevNet API.
//
// It does an HTTP GET request to the "/predeployed_accounts" endpoint and
// decodes the response body into a slice of TestAccount structs. It returns
// the list of accounts and any error that occurred during the process.
//
// Parameters:
//
//	none
//
// Returns:
// - []TestAccount: a slice of TestAccount structs
func (devnet *DevNet) Accounts() ([]TestAccount, error) {
	req, err := http.NewRequest(http.MethodGet, devnet.api("/predeployed_accounts"), nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var accounts []TestAccount
	err = json.NewDecoder(resp.Body).Decode(&accounts)
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
	req, err := http.NewRequest(http.MethodGet, devnet.api("/is_alive"), nil)
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
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

type MintResponse struct {
	NewBalance      string `json:"new_balance"`
	Unit            string `json:"unit"`
	TransactionHash string `json:"tx_hash"`
}

// Mint mints a certain amount of tokens for a given address.
//
// Parameters:
// - address: is the address to mint tokens for
// - amount: is the amount of tokens to mint
// Returns:
// - *MintResponse: a MintResponse
// - error: an error if any
func (devnet *DevNet) Mint(address *felt.Felt, amount *big.Int) (*MintResponse, error) {
	data := struct {
		Address *felt.Felt `json:"address"`
		Amount  *big.Int   `json:"amount"`
	}{
		Address: address,
		Amount:  amount,
	}
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	payloadReader := bytes.NewBuffer(payload)
	req, err := http.NewRequest(http.MethodPost, devnet.api("/mint"), payloadReader)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var mint MintResponse
	err = json.NewDecoder(resp.Body).Decode(&mint)
	return &mint, err
}

type FeeToken struct {
	Symbol  string
	Address *felt.Felt
}

// FeeToken retrieves the fee token from the DevNet API.
//
// This function does a GET request to the "/fee_token" endpoint of the DevNet API
// to retrieve the fee token.
//
// Parameters:
//
//	none
//
// Returns:
//   - *FeeToken: a pointer to a FeeToken object
//   - error: an error, if any
func (devnet *DevNet) FeeToken() (*FeeToken, error) {
	req, err := http.NewRequest("GET", devnet.api("/fee_token"), nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var token FeeToken
	err = json.NewDecoder(resp.Body).Decode(&token)
	return &token, err
}
