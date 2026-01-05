package rpcv10

// @todo remove this file later

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

// TestSetup is a type that is used to store setup data for the RPC tests.
type TestSetup struct {
	Base     string
	Provider *Provider
	RPCSpy   tests.RPCSpyer

	WsBase     string
	WsProvider *WsProvider
	WSSpy      tests.WSSpyer

	// Only present in mock environment
	MockClient *clientmock.MockClient

	AccountAddress string
	PrivKey        string
	PubKey         string
}

// BeforeEach initialises the environment setup before running the tests.
// It must be called inside subtests if that's the case.
//
// Parameters:
//   - t: The testing.T object
//   - isWs: a boolean value to check if the test will use the websocket provider
//
// Returns:
//   - TestSetup: the TestSetup struct containing the setup data
func BeforeEach(t *testing.T, isWs bool) TestSetup {
	t.Helper()

	var testConfig TestSetup

	if tests.TEST_ENV == tests.MockEnv {
		mockCtrl := gomock.NewController(t)
		mockClient := clientmock.NewMockClient(mockCtrl)

		spy := tests.NewRPCSpy(mockClient)
		provider := &Provider{
			c: spy,
		}

		wsSpy := tests.NewWSSpy(mockClient)
		wsProvider := &WsProvider{
			s: wsSpy,
		}

		testConfig.MockClient = mockClient
		testConfig.Provider = provider
		testConfig.RPCSpy = spy
		testConfig.WsProvider = wsProvider
		testConfig.WSSpy = wsSpy

		return testConfig
	}

	base := os.Getenv("HTTP_PROVIDER_URL")
	if base != "" {
		testConfig.Base = base
	}

	provider, err := NewProvider(t.Context(), testConfig.Base)
	if err != nil {
		t.Fatalf("failed to connect to the %s provider: %v", testConfig.Base, err)
	}

	spy := tests.NewRPCSpy(provider.c)
	testConfig.RPCSpy = spy
	provider.c = spy

	testConfig.Provider = provider
	t.Cleanup(func() {
		testConfig.Provider.c.Close()
	})

	if tests.TEST_ENV == tests.DevnetEnv {
		return testConfig
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

		spy := tests.NewWSSpy(wsClient.s)
		testConfig.WSSpy = spy
		wsClient.s = spy

		testConfig.WsProvider = wsClient
		t.Cleanup(func() {
			testConfig.WsProvider.s.Close()
		})
	}

	// load the test account data, only required for some tests
	testConfig.PrivKey = os.Getenv("STARKNET_PRIVATE_KEY")
	testConfig.PubKey = os.Getenv("STARKNET_PUBLIC_KEY")
	testConfig.AccountAddress = os.Getenv("STARKNET_ACCOUNT_ADDRESS")

	return testConfig
}
