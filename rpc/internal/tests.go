package internal

import (
	"os"
	"testing"

	"github.com/NethermindEth/starknet.go/internal/tests"
	"github.com/NethermindEth/starknet.go/internal/tests/mocks/clientmock"
	"github.com/NethermindEth/starknet.go/rpc"
	"go.uber.org/mock/gomock"
)

func TestMain(m *testing.M) {
	tests.LoadEnv()

	os.Exit(m.Run())
}

// TestSetup is a type that is used to store setup data for the RPC tests.
type TestSetup struct {
	Base string
	// @todo rename
	Provider rpc.Caller
	RPCSpy   tests.RPCSpyer

	WsBase string
	// @todo rename
	WsProvider rpc.Subscriber
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
		provider := spy

		wsSpy := tests.NewWSSpy(mockClient)
		wsProvider := wsSpy

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

	client, err := NewHTTPClient(t.Context(), testConfig.Base)
	if err != nil {
		t.Fatalf("failed to connect to the %s provider: %v", testConfig.Base, err)
	}

	spy := tests.NewRPCSpy(client)
	testConfig.RPCSpy = spy

	testConfig.Provider = spy
	t.Cleanup(func() {
		testConfig.Provider.Close()
	})

	if tests.TEST_ENV == tests.DevnetEnv {
		return testConfig
	}

	if isWs {
		wsBase := os.Getenv("WS_PROVIDER_URL")
		if wsBase != "" {
			testConfig.WsBase = wsBase
		}

		wsClient, err := NewWSClient(t.Context(), testConfig.WsBase)
		if err != nil {
			t.Fatalf("failed to connect to the %s websocket provider: %v", testConfig.WsBase, err)
		}

		spy := tests.NewWSSpy(wsClient)
		testConfig.WSSpy = spy

		testConfig.WsProvider = spy
		t.Cleanup(func() {
			testConfig.WsProvider.Close()
		})
	}

	// load the test account data, only required for some tests
	testConfig.PrivKey = os.Getenv("STARKNET_PRIVATE_KEY")
	testConfig.PubKey = os.Getenv("STARKNET_PUBLIC_KEY")
	testConfig.AccountAddress = os.Getenv("STARKNET_ACCOUNT_ADDRESS")

	return testConfig
}
