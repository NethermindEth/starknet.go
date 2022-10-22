package contracts

import (
	"context"
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/dontpanicdao/caigo/rpcv01"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/joho/godotenv"
)

var (
	// set the environment for the test, default: mock
	testEnv = "devnet"
)

// TestMain is used to trigger the tests and, in that case, check for the environment to use.
func TestMain(m *testing.M) {
	flag.StringVar(&testEnv, "env", "devnet", "set the test environment")
	flag.Parse()
	os.Exit(m.Run())
}

type testConfiguration struct {
	baseURL string
	client  *rpcv01.Provider
}

var testConfigurations = map[string]testConfiguration{
	"devnet": {
		baseURL: "http://localhost:5050/rpc",
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
		c, err := ethrpc.DialContext(context.Background(), testConfig.baseURL)
		if err != nil {
			t.Fatal("connect should succeed, instead:", err)
		}
		client := rpcv01.NewProvider(c)
		testConfig.client = client
	}
	t.Cleanup(func() {
	})
	return &testConfig
}
