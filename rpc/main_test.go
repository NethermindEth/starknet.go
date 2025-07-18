package rpc

import (
	"os"
	"testing"

	"github.com/NethermindEth/starknet.go/internal/tests"
)

func TestMain(m *testing.M) {
	tests.LoadEnv()

	os.Exit(m.Run())
}

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

	if tests.TEST_ENV == tests.MockEnv {
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

	if tests.TEST_ENV == tests.DevnetEnv || tests.TEST_ENV == tests.MainnetEnv {
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
