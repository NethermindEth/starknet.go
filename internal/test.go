package internal

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
)

// the environment for the test, default: mock
var testEnv = ""

func loadFlags() {
	// set the environment for the test, default: mock
	flag.StringVar(&testEnv, "env", "mock", "set the test environment")
	flag.Parse()
}

func LoadEnv() string {
	loadFlags()

	switch testEnv {
	case "mock", "mainnet", "testnet", "devnet":
		break
	default:
		log.Fatalf("invalid test environment '%s', supports mock, testnet, mainnet, devnet", testEnv)
	}

	// get the source file path
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("Failed to get source file path")
	}

	// get the directory containing the source file
	sourceDir := filepath.Dir(filename)

	err := godotenv.Load(filepath.Join(sourceDir, fmt.Sprintf(".env.%s", testEnv)), filepath.Join(sourceDir, ".env"))
	if err != nil {
		log.Printf("failed to load env, err: %s", err)
	}

	return testEnv
}
