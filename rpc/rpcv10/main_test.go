package rpcv10

// @todo remove this file later

import (
	"os"
	"testing"

	"github.com/NethermindEth/starknet.go/internal/tests"
	"github.com/NethermindEth/starknet.go/internal/tests/mocks/clientmock"
	"github.com/NethermindEth/starknet.go/rpc/internal"
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

// BeforeEach forwards the call to the internal.BeforeEach function.“
func BeforeEach(t *testing.T, isWs bool) TestSetup {
	baseSetup := internal.BeforeEach(t, isWs)

	provider := Provider{c: baseSetup.Provider}
	wsProvider := WsProvider{s: baseSetup.WsProvider}

	return TestSetup{
		Base:           baseSetup.Base,
		Provider:       &provider,
		RPCSpy:         baseSetup.RPCSpy,
		WsBase:         baseSetup.WsBase,
		WsProvider:     &wsProvider,
		WSSpy:          baseSetup.WSSpy,
		MockClient:     baseSetup.MockClient,
		AccountAddress: baseSetup.AccountAddress,
		PrivKey:        baseSetup.PrivKey,
		PubKey:         baseSetup.PubKey,
	}
}
