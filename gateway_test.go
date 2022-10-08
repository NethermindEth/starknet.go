package caigo

import (
	"testing"

	"github.com/dontpanicdao/caigo/types"
)

const (
	SEED             int    = 100000000
	CONTRACT_ADDRESS string = "0x02b795d8c5e38c45da3b89c91174c66a3c77845bbeb87a36038f19c521dbe87e"
)

type TestAccountType struct {
	PrivateKey   string               `json:"private_key"`
	PublicKey    string               `json:"public_key"`
	Address      string               `json:"address"`
	Transactions []types.FunctionCall `json:"transactions,omitempty"`
}

func TestGatewayAccount_Execute(t *testing.T) {
	testConfig := beforeGatewayEach(t)

	type testSetType struct {
		Calls []types.FunctionCall `json:"transactions,omitempty"`
	}

	testSet := map[string][]testSetType{
		"devnet":  {},
		"testnet": {},
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		_ = test.Calls
		_ = testConfig.base
	}
}
