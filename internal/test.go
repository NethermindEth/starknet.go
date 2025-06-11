package internal

import (
	"flag"
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

	customEnv := ".env." + testEnv
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

	return testEnv
}
