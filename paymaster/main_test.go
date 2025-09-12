package paymaster

import (
	"os"
	"testing"

	"github.com/NethermindEth/starknet.go/client"
	"github.com/NethermindEth/starknet.go/internal/tests"
	"github.com/NethermindEth/starknet.go/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var avnuPaymasterURL = "https://sepolia.paymaster.avnu.fi"

func TestMain(m *testing.M) {
	tests.LoadEnv()
	os.Exit(m.Run())
}

type MockPaymaster struct {
	*Paymaster
	// this should be a pointer to the mock client used in the Paymaster struct.
	// This is intended to have an easy access to the mock client, without having to
	// type cast it from the `callCloser` interface every time.
	c *mocks.MockClient
}

// Creates a real Sepolia paymaster client.
func SetupPaymaster(t *testing.T) *Paymaster {
	t.Helper()

	var pm *Paymaster
	var err error

	if tests.TEST_ENV == tests.IntegrationEnv {
		apiKey := os.Getenv("AVNU_API_KEY")
		require.NotEmpty(t, apiKey, "AVNU_API_KEY is not set")
		apiHeader := client.WithHeader("x-paymaster-api-key", apiKey)
		pm, err = NewPaymasterClient(avnuPaymasterURL, apiHeader)
	} else {
		pm, err = NewPaymasterClient(avnuPaymasterURL)
	}

	require.NoError(t, err, "failed to create paymaster client")

	return pm
}

// Creates a mock paymaster client.
func SetupMockPaymaster(t *testing.T) *MockPaymaster {
	t.Helper()

	client := mocks.NewMockClient(gomock.NewController(t))
	mpm := &MockPaymaster{
		Paymaster: &Paymaster{c: client},
		c:         client,
	}

	return mpm
}
