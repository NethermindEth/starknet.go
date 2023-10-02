package xsessions

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/NethermindEth/starknet.go/rpc"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/joho/godotenv"
)

// beforeEachRPC initializes and returns a new rpc.Provider instance for each test case, with the necessary setup.
//
// The function takes a testing.T instance as a parameter, which is used to mark the function as a helper function.
// The method godotenv.Load(".env.devnet") loads the environment variables from the .env.devnet file.
// The URL is retrieved from the STARKNET_NODE_URL environment variable using os.Getenv("STARKNET_NODE_URL").
// If the URL is empty, the function fails with a fatal error message.
// The function then establishes a connection to the StarkNet node using ethrpc.DialContext.
// If the connection fails, the function fails with a fatal error message.
// The rpc.Provider is created using the established connection.
// The t.Cleanup function is called to ensure that the connection is closed after the test case finishes executing.
// Finally, the rpc.Provider is returned.
func beforeEachRPC(t *testing.T) *rpc.Provider {
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
	provider := rpc.NewProvider(c)
	t.Cleanup(func() {
		c.Close()
	})
	return provider
}
