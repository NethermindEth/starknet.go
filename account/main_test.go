package account

import (
	"errors"
	"fmt"
	"math/big"
	"os"
	"testing"

	"github.com/NethermindEth/starknet.go/devnet"
	"github.com/NethermindEth/starknet.go/internal/tests"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/NethermindEth/starknet.go/rpc/rpcv10"
	"github.com/stretchr/testify/require"
)

type TestSetup struct {
	// the ProviderURL url for the test
	ProviderURL string
	// the test account data
	PrivKey        string
	PubKey         string
	AccountAddress string

	Wrapper     providerWrapper
	MockWrapper *MockproviderWrapper
	Account     *Account
}

// TestMain is the main function for the account tests.
func TestMain(m *testing.M) {
	tests.LoadEnv()

	os.Exit(m.Run())
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
func BeforeEach(t *testing.T) (TestSetup, *Account) {
	t.Helper()
	// @todo make this work

	var tConfig TestSetup

	if tests.TEST_ENV == tests.MockEnv {
		// mockCtrl := gomock.NewController(t)
		// mockClient := clientmock.NewMockClient(mockCtrl)

		// spy := tests.NewRPCSpy(mockClient)
		// provider := spy

		// wsSpy := tests.NewWSSpy(mockClient)
		// wsProvider := wsSpy

		// testConfig.MockClient = mockClient
		// testConfig.Provider = provider
		// testConfig.RPCSpy = spy
		// testConfig.WsProvider = wsProvider
		// testConfig.WSSpy = wsSpy

		// return testConfig
	}

	tConfig.ProviderURL = os.Getenv("HTTP_PROVIDER_URL")
	if tConfig.ProviderURL == "" {
		panic("Failed to load HTTP_PROVIDER_URL, empty string")
	}

	// load the test account data, only required for some tests
	tConfig.PrivKey = os.Getenv("STARKNET_PRIVATE_KEY")
	tConfig.PubKey = os.Getenv("STARKNET_PUBLIC_KEY")
	tConfig.AccountAddress = os.Getenv("STARKNET_ACCOUNT_ADDRESS")

	return tConfig, &Account{}
}

// returns a new account type from the provided account data in the tConfig
func setupAcc(t *testing.T, tsetup *TestSetup) (*Account, error) {
	t.Helper()

	ks := NewMemKeystore()
	privKeyBI, ok := new(big.Int).SetString(tsetup.PrivKey, 0)
	if !ok {
		return nil, errors.New("failed to convert privKey to big.Int")
	}
	ks.Put(tsetup.PubKey, privKeyBI)

	accAddress, err := internalUtils.HexToFelt(tsetup.AccountAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to convert accountAddress to felt: %w", err)
	}

	// @todo make it work with any provider
	acc, err := NewAccount(&rpcv10.Provider{}, accAddress, tsetup.PubKey, ks, CairoV2)
	if err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	return acc, nil
}

// newDevnet creates a new devnet with the given URL.
//
// Parameters:
//   - t: The testing.T instance for running the test
//   - url: The URL of the devnet to be created
//
// Returns:
//   - *devnet.DevNet: a pointer to a devnet object
//   - []devnet.TestAccount: a slice of test accounts
//   - error: an error, if any
func newDevnet(t *testing.T, url string) (*devnet.DevNet, []devnet.TestAccount, error) {
	t.Helper()
	devnetInstance := devnet.NewDevNet(url)
	acnts, err := devnetInstance.Accounts()

	return devnetInstance, acnts, err
}

// newDevnetAccount creates a new devnet account
//
// Parameters:
//   - t: The testing.T instance for running the test
//   - provider: The RPC provider
//   - accData: The test account data
//
// Returns:
//   - *Account: The new devnet account
//   - error: An error, if any
func newDevnetAccount(
	t *testing.T,
	provider providerWrapper,
	accData devnet.TestAccount,
	cairoVersion CairoVersion,
) *Account {
	t.Helper()
	fakeUserAddr := internalUtils.TestHexToFelt(t, accData.Address)
	fakeUserPriv := internalUtils.TestHexToFelt(t, accData.PrivateKey)

	// Set up ks
	ks := NewMemKeystore()
	ks.Put(accData.PublicKey, fakeUserPriv.BigInt(new(big.Int)))

	// @todo make it work with any provider
	acnt, err := NewAccount(&rpcv10.Provider{}, fakeUserAddr, accData.PublicKey, ks, cairoVersion)
	require.NoError(t, err)

	return acnt
}
