package tests

import (
	"flag"
	"log"
	"path/filepath"
	"runtime"
	"slices"
	"testing"

	"github.com/joho/godotenv"
)

// The environment for the test defined by the `-env` flag. If not set, default: mock
var TEST_ENV TestEnv //nolint:staticcheck //Only used in tests

// An enum representing the environments for the test.
type TestEnv string

//nolint:staticcheck // Only used in tests
const (
	MockEnv           TestEnv = "mock"
	IntegrationEnv    TestEnv = "integration"
	TestnetEnv        TestEnv = "testnet"
	MainnetEnv        TestEnv = "mainnet"
	DevnetEnv         TestEnv = "devnet"
	Devnet_TestnetEnv TestEnv = "devnet-testnet"
)

func loadEnvFlag() {
	var testEnvStr string
	// set the environment for the test, default: mock
	flag.StringVar(&testEnvStr, "env", string(MockEnv), "set the test environment")
	flag.Parse()

	TEST_ENV = TestEnv(testEnvStr)
}

// Loads the environment for the tests. It must be called before `m.Run` in the TestMain function
// of each package.
// It looks for a `.env.<testEnv>` file in the `internal/tests` directory, where `<testEnv>` is
// the value of the `-env` flag.
// If the file is not found, a warning is logged.
func LoadEnv() {
	loadEnvFlag()

	switch TEST_ENV {
	case MockEnv, IntegrationEnv, TestnetEnv, MainnetEnv, DevnetEnv, Devnet_TestnetEnv:
		break
	default:
		log.Fatalf(
			`invalid test environment '%s', supports: mock, integration, testnet,
			mainnet, devnet, devnet-testnet`,
			TEST_ENV,
		)
	}

	// get the source file path
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("Failed to get source file path")
	}

	// get the directory containing the source file
	sourceDir := filepath.Dir(filename)

	customEnv := ".env." + string(TEST_ENV)
	err := godotenv.Load(filepath.Join(sourceDir, customEnv))
	if err != nil {
		log.Printf("Warning: failed to load %s, err: %s", customEnv, err)
	} else {
		log.Printf("Successfully loaded %s", customEnv)
	}
}

// Runs the test only if the environment is in the list of environments provided.
// If no environment is provided, the test fails.
// If the environment is not in the list of environments, the test is skipped.
//
// Packages like `account` and `rpc` are run on multiple environments, so in order to
// easily skip tests that are not supported on a specific environment, we use this helper.
// If a test is not supported on a specific environment, it should be skipped; otherwise, it could
// pass (since the environment is not being tested) and give a false positive.
func RunTestOn(t *testing.T, envs ...TestEnv) {
	t.Helper()

	if len(envs) == 0 {
		t.Fatal("No test environment provided")
	}

	if !slices.Contains(envs, TEST_ENV) {
		t.Skipf("Test not defined to run on '%s' environment, skipping...", TEST_ENV)
	}
}
