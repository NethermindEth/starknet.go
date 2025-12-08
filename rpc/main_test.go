package rpc

import (
	"os"
	"testing"

	"github.com/NethermindEth/starknet.go/internal/tests"
	"github.com/NethermindEth/starknet.go/internal/tests/mocks/clientmock"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestMain(m *testing.M) {
	tests.LoadEnv()

	os.Exit(m.Run())
}

// TestSetup is a type that is used to store setup data for the RPC tests.
type TestSetup struct {
	Provider   *Provider
	WsProvider *WsProvider
	Base       string
	WsBase     string
	Spy        tests.Spyer
	// Only present in mock environment
	MockClient *clientmock.MockClient

	AccountAddress string
	PrivKey        string
	PubKey         string
}

// BeforeEach initialises the environment setup before running the tests.
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

		spy := tests.NewJSONRPCSpy(mockClient)

		provider := &Provider{
			c: spy,
		}

		testConfig.MockClient = mockClient
		testConfig.Provider = provider
		testConfig.Spy = spy

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

	spy := tests.NewJSONRPCSpy(provider.c)
	testConfig.Spy = spy
	provider.c = spy

	testConfig.Provider = provider
	t.Cleanup(func() {
		testConfig.Provider.c.Close()
	})

	if tests.TEST_ENV == tests.DevnetEnv || tests.TEST_ENV == tests.MainnetEnv {
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
		testConfig.WsProvider = wsClient
		t.Cleanup(func() {
			testConfig.WsProvider.c.Close()
		})
	}

	// load the test account data, only required for some tests
	testConfig.PrivKey = os.Getenv("STARKNET_PRIVATE_KEY")
	testConfig.PubKey = os.Getenv("STARKNET_PUBLIC_KEY")
	testConfig.AccountAddress = os.Getenv("STARKNET_ACCOUNT_ADDRESS")

	return testConfig
}

// GetCommonBlockIDs returns a list of common block IDs to use in the RPC tests.
// It includes all block tags, a range of block numbers and the latest block hash.
func GetCommonBlockIDs(t *testing.T, provider *Provider) []BlockID {
	t.Helper()

	// *** all valid block tags ***
	commonBlockIDs := []BlockID{
		WithBlockTag(BlockTagLatest),
		WithBlockTag(BlockTagPreConfirmed),
		WithBlockTag(BlockTagL1Accepted),
	}

	// *** getting the common block number range ***

	// 5 blocks from the first 1M blocks of the network
	// (a lot of changes in the first blocks)
	commonBlockIDs = append(commonBlockIDs, []BlockID{
		WithBlockNumber(0),
		WithBlockNumber(200_000),
		WithBlockNumber(400_000),
		WithBlockNumber(600_000),
		WithBlockNumber(800_000),
		WithBlockNumber(1_000_000),
	}...)

	// get the latest block number of the network
	blockHashAndNumber, err := provider.BlockHashAndNumber(t.Context())
	require.NoError(t, err, "failed to get the block number")

	// after the block 1_000_000, we add one block every 500_000 blocks
	// until the latest block
	for i := uint64(1_500_000); i < blockHashAndNumber.Number; i += 500_000 {
		commonBlockIDs = append(commonBlockIDs, WithBlockNumber(i))
	}

	// add the latest block hash
	commonBlockIDs = append(commonBlockIDs, WithBlockHash(blockHashAndNumber.Hash))

	return commonBlockIDs
}
