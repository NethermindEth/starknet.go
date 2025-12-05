package paymaster

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/client"
	"github.com/NethermindEth/starknet.go/internal/tests"
	"github.com/NethermindEth/starknet.go/internal/tests/mocks/clientmock"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const avnuPaymasterURL = "https://sepolia.paymaster.avnu.fi"

func TestMain(m *testing.M) {
	tests.LoadEnv()
	os.Exit(m.Run())
}

type MockPaymaster struct {
	*Paymaster
	// this should be a pointer to the mock client used in the Paymaster struct.
	// This is intended to have an easy access to the mock client, without having to
	// type cast it from the `callCloser` interface every time.
	c *clientmock.MockClient
}

// Creates a real Sepolia paymaster client and a spy for integration tests.
func SetupPaymaster(t *testing.T, debug ...bool) (*Paymaster, tests.Spyer) {
	t.Helper()

	apiKey := os.Getenv("AVNU_API_KEY")
	require.NotEmpty(t, apiKey, "AVNU_API_KEY is not set")
	apiHeader := client.WithHeader("x-paymaster-api-key", apiKey)

	pm, err := New(context.Background(), avnuPaymasterURL, apiHeader)
	require.NoError(t, err, "failed to create paymaster client")

	spy := tests.NewJSONRPCSpy(pm.c, debug...)
	pm.c = spy

	return pm, spy
}

// Creates a mock paymaster client.
func SetupMockPaymaster(t *testing.T) *MockPaymaster {
	t.Helper()

	pmClient := clientmock.NewMockClient(gomock.NewController(t))
	mpm := &MockPaymaster{
		Paymaster: &Paymaster{c: pmClient},
		c:         pmClient,
	}

	return mpm
}

// GetStrkAccountData returns the STRK account data from the environment variables.
// This is used for integration tests, where we need a real testnet account with STRK tokens.
func GetStrkAccountData(t *testing.T) (privKey, pubKey, accountAddress felt.Felt) {
	t.Helper()

	strkPrivKey := os.Getenv("STARKNET_PRIVATE_KEY")
	strkPubKey := os.Getenv("STARKNET_PUBLIC_KEY")
	strkAccountAddress := os.Getenv("STARKNET_ACCOUNT_ADDRESS")

	tempPrivKey := internalUtils.TestHexToFelt(t, strkPrivKey)
	tempPubKey := internalUtils.TestHexToFelt(t, strkPubKey)
	tempAccountAddress := internalUtils.TestHexToFelt(t, strkAccountAddress)

	return *tempPrivKey, *tempPubKey, *tempAccountAddress
}

// CompareEnumsHelper compares an enum type with the expected value and error expected.
func CompareEnumsHelper[T any](t *testing.T, input string, expected T, errorExpected bool) {
	t.Helper()

	var actual T
	err := json.Unmarshal([]byte(input), &actual)
	if errorExpected {
		assert.Error(t, err)
	} else {
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		marshalled, err := json.Marshal(actual)
		assert.NoError(t, err)
		assert.Equal(t, input, string(marshalled))
	}
}
