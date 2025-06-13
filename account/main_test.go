package account_test

import (
	"errors"
	"fmt"
	"math/big"
	"os"
	"testing"

	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/devnet"
	"github.com/NethermindEth/starknet.go/internal"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/stretchr/testify/require"
)

type testConfig struct {
	// the providerURL url for the test
	providerURL string
	// the test account data
	privKey        string
	pubKey         string
	accountAddress string
}

var (
	tConfig testConfig
)

// TestMain is used to trigger the tests and, in that case, check for the environment to use.
func TestMain(m *testing.M) {
	if internal.TEST_ENV == internal.MockEnv {
		os.Exit(m.Run())
	}
	tConfig.providerURL = os.Getenv("HTTP_PROVIDER_URL")
	if tConfig.providerURL == "" {
		panic("Failed to load HTTP_PROVIDER_URL, empty string")
	}

	// load the test account data, only required for some tests
	tConfig.privKey = os.Getenv("STARKNET_PRIVATE_KEY")
	tConfig.pubKey = os.Getenv("STARKNET_PUBLIC_KEY")
	tConfig.accountAddress = os.Getenv("STARKNET_ACCOUNT_ADDRESS")

	os.Exit(m.Run())
}

// returns a new account type from the provided account data in the tConfig
func setupAcc(t *testing.T, provider rpc.RpcProvider) (*account.Account, error) {
	t.Helper()

	ks := account.NewMemKeystore()
	privKeyBI, ok := new(big.Int).SetString(tConfig.privKey, 0)
	if !ok {
		return nil, errors.New("failed to convert privKey to big.Int")
	}
	ks.Put(tConfig.pubKey, privKeyBI)

	accAddress, err := internalUtils.HexToFelt(tConfig.accountAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to convert accountAddress to felt: %w", err)
	}

	acc, err := account.NewAccount(provider, accAddress, tConfig.pubKey, ks, account.CairoV2)
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

// newDevnetAccount creates a new devnet account from a test account.
//
// Parameters:
//   - t: The testing.T instance for running the test
//   - provider: The RPC provider
//   - accData: The test account data
//
// Returns:
//   - *account.Account: The new devnet account
//   - error: An error, if any
func newDevnetAccount(
	t *testing.T,
	provider *rpc.Provider,
	accData devnet.TestAccount,
	cairoVersion account.CairoVersion,
) *account.Account {
	t.Helper()
	fakeUserAddr := internalUtils.TestHexToFelt(t, accData.Address)
	fakeUserPriv := internalUtils.TestHexToFelt(t, accData.PrivateKey)

	// Set up ks
	ks := account.NewMemKeystore()
	ks.Put(accData.PublicKey, fakeUserPriv.BigInt(new(big.Int)))

	acnt, err := account.NewAccount(provider, fakeUserAddr, accData.PublicKey, ks, cairoVersion)
	require.NoError(t, err)

	return acnt
}
