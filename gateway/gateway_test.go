package gateway

import (
	"context"
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

// testConfiguration is a type that is used to configure tests
type testConfiguration struct {
	client *Gateway
	base   string
}

var (
	// set the environment for the test, default: mock
	testEnv = "mock"

	// testConfigurations are predefined test configurations
	testConfigurations = map[string]testConfiguration{
		"mainnet": {
			base: "https://alpha-mainnet.starknet.io",
		},
		"testnet": {
			base: "https://alpha4.starknet.io",
		},
		"devnet": {
			base: "http://localhost:5000",
		},
		"mock": {},
	}
)

// TestMain is used to trigger the tests and, in that case, check for the environment to use.
func TestMain(m *testing.M) {
	flag.StringVar(&testEnv, "env", "mock", "set the test environment")
	flag.Parse()

	os.Exit(m.Run())
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
	case "mock":
		testConfig.client = &Gateway{
			client: &httpMock{},
		}
	default:
		testConfig.client = NewClient(WithChain(testEnv))
	}
	t.Cleanup(func() {
	})
	return &testConfig
}

// TestGateway checks the gateway can be accessed.
func TestGateway(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		BlockHash string
	}
	testSet := map[string][]testSetType{
		"devnet":  {},
		"mainnet": {{BlockHash: "0x4ee4c886d1767b7165a1e3a7c6ad145543988465f2bda680c16a79536f6d81f"}},
		"mock":    {{BlockHash: "0xdeadbeef"}},
		"testnet": {{BlockHash: "0x787af09f1cacdc5de1df83e8cdca3a48c1194171c742e78a9f684cb7aa4db"}},
	}[testEnv]

	for _, test := range testSet {
		block, err := testConfig.client.Block(context.Background(), &BlockOptions{BlockHash: test.BlockHash})

		if err != nil {
			t.Fatal(err)
		}
		if block.BlockHash != test.BlockHash {
			t.Fatalf("expecting %s, instead: %s", "", block.BlockHash)
		}
	}
}
