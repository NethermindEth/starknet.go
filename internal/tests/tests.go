package tests

import (
	"flag"
	"log"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
)

func init() {
	TEST_ENV = loadEnv()
}

// the environment for the test, default: mock
var TEST_ENV TestEnv

type TestEnv string

const (
	MockEnv           TestEnv = "mock"
	IntegrationEnv    TestEnv = "integration"
	TestnetEnv        TestEnv = "testnet"
	MainnetEnv        TestEnv = "mainnet"
	DevnetEnv         TestEnv = "devnet"
	Devnet_TestnetEnv TestEnv = "devnet-testnet"
)

func loadFlags() {
	var testEnvStr string
	// set the environment for the test, default: mock
	flag.StringVar(&testEnvStr, "env", string(MockEnv), "set the test environment")
	flag.Parse()

	TEST_ENV = TestEnv(testEnvStr)
}

func loadEnv() TestEnv {
	loadFlags()

	switch TEST_ENV {
	case MockEnv, IntegrationEnv, TestnetEnv, MainnetEnv, DevnetEnv, Devnet_TestnetEnv:
		break
	default:
		log.Fatalf("invalid test environment '%s', supports: mock, integration, testnet, mainnet, devnet, devnet-testnet", TEST_ENV)
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
		log.Printf("successfully loaded %s", customEnv)
	}

	err = godotenv.Load(filepath.Join(sourceDir, ".env"))
	if err != nil {
		log.Printf("Warning: failed to load .env, err: %s", err)
	} else {
		log.Printf("successfully loaded .env")
	}

	return TEST_ENV
}
