package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dontpanicdao/caigo/types"
)

type DevNet struct {
	baseURL string
}

type TestAccount struct {
	PrivateKey string `json:"private_key"`
	PublicKey  string `json:"public_key"`
	Address    string `json:"address"`
}

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

func (devnet *DevNet) api(uri string) string {
	uri = strings.TrimPrefix(uri, "/")
	return fmt.Sprintf("%s/%s", devnet.baseURL, uri)
}

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
	NewBalance uint64 `json:"new_balance"`
	Unit       string `json:"unit"`
}

func (devnet *DevNet) Mint(address types.Hash, amount uint64) (*MintResponse, error) {
	data := struct {
		Address types.Hash `json:"address"`
		Amount  uint64     `json:"amount"`
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
	Address types.Hash
}

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
