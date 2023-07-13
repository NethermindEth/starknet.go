package starknetgo

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/NethermindEth/starknet.go/artifacts"
	"github.com/NethermindEth/starknet.go/gateway"
	"github.com/NethermindEth/starknet.go/rpc"
	devtest "github.com/NethermindEth/starknet.go/test"
	"github.com/NethermindEth/starknet.go/types"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/joho/godotenv"
)

const (
	TestPublicKey            = "0x783318b2cc1067e5c06d374d2bb9a0382c39aabd009b165d7a268b882971d6"
	DevNetETHAddress         = "0x62230ea046a9a5fbc261ac77d03c8d41e5d442db2284587570ab46455fd2488"
	TestNetETHAddress        = "0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"
	DevNetAccount032Address  = "0x0536244bba4dc9bb219d964b477af6d18f7096635a96284bb0e008bf137650ec"
	TestNetAccount032Address = "0x6ca4fdd437dffde5253ba7021ef7265c88b07789aa642eafda37791626edf00"
	DevNetAccount040Address  = "0x058079067104f58fd9f1ef949cd2d2b482d7bca39b793983f077edaf51d979e9"
	TestNetAccount040Address = "0x6cbfa37f409610fee26eeb427ed854b3a4b24580d9b9ef6c3e38db7b3f7322c"
	TestnetCounterAddress    = "0x51e94d515df16ecae5be4a377666121494eb54193d854fcf5baba2b0da679c6"
)

// testGatewayConfiguration is a type that is used to configure tests
type testGatewayConfiguration struct {
	client            *gateway.GatewayProvider
	base              string
	CounterAddress    string
	AccountAddress    string
	AccountPrivateKey string
}

var (
	// set the environment for the test, default: mock
	testEnv = "mock"

	// testConfigurations are predefined test configurations
	testRPCConfigurations = map[string]testRPCConfiguration{
		// Requires a Mainnet StarkNet JSON-RPC compliant node (e.g. pathfinder)
		// (ref: https://github.com/eqlabs/pathfinder)
		"mainnet": {
			base: "http://localhost:9545",
		},
		// Requires a Testnet StarkNet JSON-RPC compliant node (e.g. pathfinder)
		// (ref: https://github.com/eqlabs/pathfinder)
		"testnet": {
			base: "http://localhost:9545",
		},
		// Requires a Devnet configuration running locally
		// (ref: https://github.com/Shard-Labs/starknet-devnet)
		"devnet": {
			base: "http://localhost:5050/rpc",
		},
		// Used with a mock as a standard configuration, see `mock_test.go``
		"mock": {},
	}

	testGatewayConfigurations = map[string]testGatewayConfiguration{
		"mainnet": {
			base: "https://alpha4-mainnet.starknet.io",
		},
		// Requires a Testnet StarkNet JSON-RPC compliant node (e.g. pathfinder)
		// (ref: https://github.com/eqlabs/pathfinder)
		"testnet": {
			base:           "https://alpha4.starknet.io",
			CounterAddress: TestnetCounterAddress,
			AccountAddress: TestNetAccount040Address,
		},
		// Requires a Devnet configuration running locally
		// (ref: https://github.com/Shard-Labs/starknet-devnet)
		"devnet": {
			base: "http://localhost:5050",
		},
		// Used with a mock as a standard configuration, see `mock_test.go``
		"mock": {},
	}
)

func InstallCounterContract(provider *gateway.GatewayProvider) (string, error) {
	class := types.ContractClass{}
	if err := json.Unmarshal(artifacts.CounterCompiled, &class); err != nil {
		return "", err
	}
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()
	tx, err := provider.Deploy(context.Background(), class, types.DeployRequest{
		ContractAddressSalt: "0x0",
		ConstructorCalldata: []string{},
	})
	if err != nil {
		return "", err
	}
	fmt.Println("deploy counter txHash", tx.TransactionHash)
	_, receipt, err := provider.WaitForTransaction(ctx, tx.TransactionHash, 3, 20)
	if err != nil {
		log.Printf("contract Address: %s\n", tx.ContractAddress)
		log.Printf("transaction Hash: %s\n", tx.TransactionHash)
		return "", err
	}
	if !receipt.Status.IsTransactionFinal() ||
		receipt.Status == types.TransactionRejected {
		return "", fmt.Errorf("installation status: %s", receipt.Status)
	}
	return tx.ContractAddress, nil
}

// beforeEach checks the configuration and initializes it before running the script
func beforeGatewayEach(t *testing.T) *testGatewayConfiguration {
	t.Helper()
	godotenv.Load(fmt.Sprintf(".env.%s", testEnv), ".env")
	testConfig, ok := testGatewayConfigurations[testEnv]
	if !ok {
		t.Fatal("env supports testnet, mainnet or devnet")
	}
	switch testEnv {
	default:
		testConfig.client = gateway.NewProvider(gateway.WithBaseURL(testConfig.base))
	}
	t.Cleanup(func() {
	})
	return &testConfig
}

// testConfiguration is a type that is used to configure tests
type testRPCConfiguration struct {
	providerv02 *rpc.Provider
	base        string
}

// TestMain is used to trigger the tests and, in that case, check for the environment to use.
func TestMain(m *testing.M) {
	baseURL := ""
	flag.StringVar(&testEnv, "env", "mock", "set the test environment")
	flag.StringVar(&baseURL, "base-url", "", "change the baseUrl")
	flag.Parse()
	godotenv.Load(fmt.Sprintf(".env.%s", testEnv), ".env")
	if baseURL != "" {
		gwLocalConfig := testGatewayConfigurations[testEnv]
		gwLocalConfig.base = baseURL
		testGatewayConfigurations[testEnv] = gwLocalConfig
		rpcLocalConfig := testRPCConfigurations[testEnv]
		rpcLocalConfig.base = baseURL
		testRPCConfigurations[testEnv] = rpcLocalConfig
	}
	switch testEnv {
	case "devnet":
		provider := gateway.NewProvider(gateway.WithBaseURL(testGatewayConfigurations["devnet"].base))
		counterAddress, err := InstallCounterContract(provider)
		if err != nil {
			fmt.Println("error installing counter contract", err)
			os.Exit(1)
		}
		localConfig := testGatewayConfigurations[testEnv]
		accounts, err := devtest.NewDevNet().Accounts()
		if err != nil {
			fmt.Println("error getting devnet accounts", err)
			os.Exit(1)
		}
		localConfig.AccountAddress = accounts[0].Address
		localConfig.AccountPrivateKey = accounts[0].PrivateKey
		localConfig.CounterAddress = counterAddress
		testGatewayConfigurations[testEnv] = localConfig
	case "testnet":
		localConfig := testGatewayConfigurations[testEnv]
		localConfig.AccountPrivateKey = os.Getenv("TESTNET_ACCOUNT_PRIVATE_KEY")
		testGatewayConfigurations[testEnv] = localConfig
	}

	os.Exit(m.Run())
}

// beforeEach checks the configuration and initializes it before running the script
func beforeRPCEach(t *testing.T) *testRPCConfiguration {
	t.Helper()
	godotenv.Load(fmt.Sprintf(".env.%s", testEnv), ".env")
	testConfig, ok := testRPCConfigurations[testEnv]
	if !ok {
		t.Fatal("env supports mock, testnet, mainnet or devnet")
	}
	testConfig.base = "https://starknet-goerli.cartridge.gg"
	base := os.Getenv("INTEGRATION_BASE")
	if base != "" {
		testConfig.base = base
	}
	c, err := ethrpc.DialContext(context.Background(), testConfig.base)
	if err != nil {
		t.Fatal("connect should succeed, instead:", err)
	}
	clientv02 := rpc.NewProvider(c)
	testConfig.providerv02 = clientv02
	return &testConfig
}

// TestChainID checks the chainId matches the one for the environment
func TestGeneral_ChainID(t *testing.T) {
	testConfig := beforeRPCEach(t)

	type testSetType struct {
		ChainID string
	}
	testSet := map[string][]testSetType{
		"devnet":  {{ChainID: "SN_GOERLI"}},
		"mainnet": {{ChainID: "SN_MAIN"}},
		"mock":    {{ChainID: "MOCK"}},
		"testnet": {{ChainID: "SN_GOERLI"}},
	}[testEnv]

	fmt.Printf("----------------------------\n")
	fmt.Printf("Env: %s\n", testEnv)
	fmt.Printf("Url: %s\n", testConfig.base)
	fmt.Printf("----------------------------\n")

	for _, test := range testSet {
		chain, err := testConfig.providerv02.ChainID(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		if chain != test.ChainID {
			t.Fatalf("expecting %s, instead: %s", test.ChainID, chain)
		}
	}
}

// TestSyncing checks the values returned are consistent
func TestGeneral_Syncing(t *testing.T) {
	testConfig := beforeRPCEach(t)

	type testSetType struct {
		ChainID string
	}

	testSet := map[string][]testSetType{
		"devnet":  {},
		"mainnet": {{ChainID: "SN_MAIN"}},
		"mock":    {{ChainID: "MOCK"}},
		"testnet": {{ChainID: "SN_GOERLI"}},
	}[testEnv]

	for range testSet {
		syncv02, err := testConfig.providerv02.Syncing(context.Background())
		if err != nil {
			t.Fatal("BlockWithTxHashes match the expected error:", err)
		}
		i, ok := big.NewInt(0).SetString(string(syncv02.CurrentBlockNum), 0)
		if !ok || i.Cmp(big.NewInt(0)) <= 0 {
			t.Fatal("CurrentBlockNum should be positive number, instead: ", syncv02.CurrentBlockNum)
		}
		if !strings.HasPrefix(syncv02.CurrentBlockHash, "0x") {
			t.Fatal("current block hash should return a string starting with 0x")
		}
	}
}
