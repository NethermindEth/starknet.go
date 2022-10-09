package test

import (
	"flag"
	"os"
	"testing"
)

const (
	DevNetETHAddress         = "0x62230ea046a9a5fbc261ac77d03c8d41e5d442db2284587570ab46455fd2488"
	TestNetETHAddress        = "0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"
	DevNetAccount032Address  = "0x06bb9425718d801fd06f144abb82eced725f0e81db61d2f9f4c9a26ece46a829"
	TestNetAccount032Address = "0x32fb76dfaa8d647c1dfa28cf5123543285250e0fcee7dfd76e4b7fa1544cfad"
	DevNetAccount040Address  = "0x080dff79c6216ad300b872b73ff41e271c63f213f8a9dc2017b164befa53b9"
	TestNetAccount040Address = "0x43eb0aebc7e9a628df79fc731cdc37b581338c913839a3f67aae2309d9e88c5"
)

// testConfiguration is a type that is used to configure tests
type testConfiguration struct {
	base string
}

var (
	// set the environment for the test, default: devnet
	testEnv = "devnet"
)

// TestMain is used to trigger the tests and, in that case, check for the environment to use.
func TestMain(m *testing.M) {
	flag.StringVar(&testEnv, "env", "devnet", "set the test environment")
	flag.Parse()
	os.Exit(m.Run())
}
