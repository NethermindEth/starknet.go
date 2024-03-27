package devnet

import (
	"flag"
	"os"
	"testing"
)

const (
	DevNetETHAddress = "0x49d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"
)

// testConfiguration is a type that is used to configure tests
type testConfiguration struct { //nolint:golint,unused
	base string
}

var (
	// set the environment for the test, default: devnet
	testEnv = "devnet"
)

// TestMain is the main test function for the package, checks configuration for the environment to use.
//
// It initializes the test environment and runs the test cases.
//
// Parameters:
// - m: is the testing.M parameter
// Returns:
//
//	none
func TestMain(m *testing.M) {
	flag.StringVar(&testEnv, "env", "devnet", "set the test environment")
	flag.Parse()
	os.Exit(m.Run())
}
