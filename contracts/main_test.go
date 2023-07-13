package contracts

import (
	"context"
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/NethermindEth/starknet.go/gateway"
	"github.com/NethermindEth/starknet.go/rpc"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/joho/godotenv"
)

var (
	// set the environment for the test, default: mock
	testEnv = "mock"
)

// TestMain is used to trigger the tests and, in that case, check for the environment to use.
func TestMain(m *testing.M) {
	flag.StringVar(&testEnv, "env", "devnet", "set the test environment")
	flag.Parse()
	os.Exit(m.Run())
}

type testConfiguration struct {
	rpc     *rpc.Provider
	gateway    *gateway.GatewayProvider
	RPCBaseURL string
	GWBaseURL  string
}

var testConfigurations = map[string]testConfiguration{
	"devnet": {
		RPCBaseURL: "http://localhost:5050/rpc",
		GWBaseURL:  "http://localhost:5050",
	},
}

// beforeEach checks the configuration and initializes it before running the script
func beforeEach(t *testing.T) *testConfiguration {
	t.Helper()
	godotenv.Load(fmt.Sprintf(".env.%s", testEnv), ".env")
	testConfig, ok := testConfigurations[testEnv]
	if !ok {
		t.Fatal("env supports testnet, mainnet or devnet")
	}
	switch testEnv {
	case "devnet":
		c, err := ethrpc.DialContext(context.Background(), testConfig.RPCBaseURL)
		if err != nil {
			t.Fatal("connect should succeed, instead:", err)
		}
		clientv02 := rpc.NewProvider(c)
		testConfig.rpc = clientv02
		gw := gateway.NewProvider(gateway.WithBaseURL(testConfig.GWBaseURL))
		testConfig.gateway = gw
	}
	t.Cleanup(func() {
	})
	return &testConfig
}
