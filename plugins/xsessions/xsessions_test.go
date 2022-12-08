package xsessions

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/dontpanicdao/caigo/rpcv01"
	"github.com/dontpanicdao/caigo/rpcv02"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/joho/godotenv"
)

// beforeEach checks the configuration and initializes it before running the script
func beforeEachRPCv01(t *testing.T) *rpcv01.Provider {
	t.Helper()
	godotenv.Load(".env.devnet")
	url := os.Getenv("STARKNET_NODE_URL")
	if url == "" {
		t.Fatalf("could not find url, check .env exists and contains STARKNET_NODE_URL")
	}
	c, err := ethrpc.DialContext(context.Background(), fmt.Sprintf("%s/rpc", url))
	if err != nil {
		t.Fatal("connect should succeed, instead:", err)
	}
	provider := rpcv01.NewProvider(c)
	t.Cleanup(func() {
		c.Close()
	})
	return provider
}

// beforeEach checks the configuration and initializes it before running the script
func beforeEachRPCv02(t *testing.T) *rpcv02.Provider {
	t.Helper()
	godotenv.Load(".env.devnet")
	url := os.Getenv("STARKNET_NODE_URL")
	if url == "" {
		t.Fatalf("could not find url, check .env exists and contains STARKNET_NODE_URL")
	}
	c, err := ethrpc.DialContext(context.Background(), fmt.Sprintf("%s/rpc", url))
	if err != nil {
		t.Fatal("connect should succeed, instead:", err)
	}
	provider := rpcv02.NewProvider(c)
	t.Cleanup(func() {
		c.Close()
	})
	return provider
}
