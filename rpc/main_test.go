package rpc

import (
	"os"
	"testing"

	"github.com/NethermindEth/starknet.go/internal/tests"
	"github.com/NethermindEth/starknet.go/internal/tests/mocks/clientmock"
	"go.uber.org/mock/gomock"
)

func TestMain(m *testing.M) {
	tests.LoadEnv()

	os.Exit(m.Run())
}

const (
	DevNetETHAddress = "0x49d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"
)

// testConfiguration is a type that is used to configure tests
type TestConfiguration struct {
	Provider   *Provider
	WsProvider *WsProvider
	Base       string
	WsBase     string
	Spy        *tests.Spyer
	// Only present in mock environment
	MockClient *clientmock.MockClient

	AccountAddress string
	PrivKey        string
	PubKey         string
}

// BeforeEach initialises the test environment configuration before running the script.
//
// Parameters:
//   - t: The testing.T object for testing purposes
//   - isWs: a boolean value to check if the test is for the websocket provider
//
// Returns:
//   - *testConfiguration: a pointer to the testConfiguration struct
func BeforeEach(t *testing.T, isWs bool) *TestConfiguration {
	t.Helper()

	var testConfig TestConfiguration

	if tests.TEST_ENV == tests.MockEnv {
		mockCtrl := gomock.NewController(t)
		mockClient := clientmock.NewMockClient(mockCtrl)

		spy := tests.NewJSONRPCSpy(mockClient)

		provider := &Provider{
			c: spy,
		}

		testConfig.MockClient = mockClient
		testConfig.Provider = provider
		testConfig.Spy = &spy

		return &testConfig
	}

	base := os.Getenv("HTTP_PROVIDER_URL")
	if base != "" {
		testConfig.Base = base
	}

	client, err := NewProvider(t.Context(), testConfig.Base)
	if err != nil {
		t.Fatalf("failed to connect to the %s provider: %v", testConfig.Base, err)
	}
	testConfig.Provider = client
	t.Cleanup(func() {
		testConfig.Provider.c.Close()
	})

	if tests.TEST_ENV == tests.DevnetEnv || tests.TEST_ENV == tests.MainnetEnv {
		return &testConfig
	}

	if isWs {
		wsBase := os.Getenv("WS_PROVIDER_URL")
		if wsBase != "" {
			testConfig.WsBase = wsBase
		}

		wsClient, err := NewWebsocketProvider(t.Context(), testConfig.WsBase)
		if err != nil {
			t.Fatalf("failed to connect to the %s websocket provider: %v", testConfig.WsBase, err)
		}
		testConfig.WsProvider = wsClient
		t.Cleanup(func() {
			testConfig.WsProvider.c.Close()
		})
	}

	// load the test account data, only required for some tests
	testConfig.PrivKey = os.Getenv("STARKNET_PRIVATE_KEY")
	testConfig.PubKey = os.Getenv("STARKNET_PUBLIC_KEY")
	testConfig.AccountAddress = os.Getenv("STARKNET_ACCOUNT_ADDRESS")

	return &testConfig
}
