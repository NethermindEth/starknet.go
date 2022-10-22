package contracts

import (
	"context"
	_ "embed"
	"fmt"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/dontpanicdao/caigo/rpcv01"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/joho/godotenv"
)

func TestInstallAccounts(t *testing.T) {
	godotenv.Load()
	if os.Getenv("INTEGRATION") != "true" {
		t.Skip("only run the test with INTEGRATION=true")
	}

	type TestCase struct {
		privateKey      string
		AccountCompiled []byte
		PluginCompiled  []byte
		ProxyCompiled   []byte
	}

	var TestCases = map[string][]TestCase{
		"devnet": {{
			privateKey:      "0x1",
			AccountCompiled: AccountV0Compiled,
		}},
	}[testEnv]
	for _, test := range TestCases {
		privateKey, _ := big.NewInt(0).SetString(test.privateKey, 0)
		ctx := context.Background()
		ctx, _ = context.WithTimeout(ctx, time.Second*60)
		client, err := ethrpc.DialContext(ctx, "http://localhost:5050/rpc")
		if err != nil {
			t.Fatalf("error connecting to devnet, %v\n", err)
		}
		provider := rpcv01.NewProvider(client)
		accountManager, err := InstallAndWaitForAccountNoWallet(ctx, provider, privateKey, test.PluginCompiled, test.AccountCompiled, test.ProxyCompiled)
		if err != nil {
			t.Fatal("should succeed, instead", err)
		}
		fmt.Println("deployment transaction", accountManager.TransactionHash)
	}
}
