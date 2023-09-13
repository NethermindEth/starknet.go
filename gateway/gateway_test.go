package gateway_test

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	starknetgo "github.com/NethermindEth/starknet.go"
	"github.com/NethermindEth/starknet.go/gateway"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/test"
	"github.com/NethermindEth/starknet.go/types"
	"github.com/NethermindEth/starknet.go/utils"
	"github.com/joho/godotenv"
)

// testConfiguration is a type that is used to configure tests
type testConfiguration struct {
	client         *gateway.Gateway
	base           string
	privateKey     string
	accountAddress string
	publicKey      string
}

var (
	// set the environment for the test, default: mock
	testEnv = "mock"

	// testConfigurations are predefined test configurations
	testConfigurations = map[string]testConfiguration{
		"mainnet": {
			base: "https://alpha-mainnet.starknet.io",
		},
		"testnet": {
			base: "https://alpha4.starknet.io",
		},
		"devnet": {
			base: "http://localhost:5050",
		},
		"mock": {},
	}
)

func setupDevnet(ctx context.Context) error {
	provider := gateway.NewProvider(gateway.WithBaseURL(testConfigurations["devnet"].base))

	v, err := test.NewDevNet().Accounts()
	if err != nil {
		return fmt.Errorf("could not connect to devnet: %v", err)
	}

	contract := rpc.DeprecatedContractClass{}
	if err := json.Unmarshal(counterCompiled, &contract); err != nil {
		return err
	}
	ks := starknetgo.NewMemKeystore()

	v0PrivKey, err := utils.HexToFelt(v[0].PrivateKey)
	if err != nil {
		return err
	}
	v0Address, err := utils.HexToFelt(v[0].Address)
	if err != nil {
		return err
	}
	account, err := starknetgo.NewGatewayAccount(
		v0PrivKey,
		v0Address,
		ks,
		provider,
		starknetgo.AccountVersion1,
	)
	if err != nil {
		return err
	}

	// starknet-class-hash --deprecated counter.json
	classHash := "0x36a03c54ac060f8083f1d1254ca824ae36bc73222a81a851a91ac0c36d852d6"

	// declare
	declare, err := account.Declare(ctx, classHash, contract, types.ExecuteDetails{})
	if err != nil {
		return err
	}
	_, receipt, err := provider.WaitForTransaction(ctx, declare.TransactionHash, 3, 10)
	if err != nil {
		log.Printf("transaction Hash: %s\n", declare.TransactionHash)
		return err
	}
	if receipt.Status == types.TransactionRejected {
		log.Printf("transaction Hash: %s\n", declare.TransactionHash)
		return fmt.Errorf("declare rejected: %+v", receipt.TransactionFailureReason)
	}

	// deploy
	tx, err := account.Deploy(ctx, classHash, types.ExecuteDetails{})
	if err != nil {
		return err
	}
	counterAddress = tx.ContractAddress
	_, receipt, err = provider.WaitForTransaction(ctx, tx.TransactionHash, 3, 10)
	if err != nil {
		// log.Printf("contract address: %s\n", tx.ContractAddress)
		log.Printf("transaction Hash: %s\n", tx.TransactionHash)
		return err
	}
	fmt.Printf("receipt: %+v\n", receipt.Events[0])
	fmt.Printf("receipt: %+v\n", receipt)
	if receipt.Status == types.TransactionRejected {
		log.Printf("transaction Hash: %s\n", tx.TransactionHash)
		return fmt.Errorf("deployed rejected: %+v", receipt.TransactionFailureReason)
	}
	return nil
}

// TestMain is used to trigger the tests and, in that case, check for the environment to use.
func TestMain(m *testing.M) {
	flag.StringVar(&testEnv, "env", "mock", "set the test environment")
	flag.Parse()
	if testEnv == "devnet" {
		err := setupDevnet(context.Background())
		if err != nil {
			log.Fatal("error starting test", err)
		}
	}
	os.Exit(m.Run())
}

// beforeEach checks the configuration and initializes it before running the script
func beforeEach(t *testing.T) *testConfiguration {
	t.Helper()
	godotenv.Load(fmt.Sprintf(".env.%s", testEnv), ".env")
	testConfig, ok := testConfigurations[testEnv]
	if !ok {
		t.Fatal("env supports testnet, mainnet or devnet")
	}
	switch testEnv {
	case "mock":
		testConfig.client = &gateway.Gateway{
			Client: &httpMock{},
		}
	case "devnet":
		v, err := test.NewDevNet().Accounts()
		if err != nil {
			t.Fatal("could not connect to devnet", err)
		}
		testConfig.privateKey = v[0].PrivateKey
		testConfig.publicKey = v[0].PublicKey
		testConfig.accountAddress = v[0].Address
		testConfig.client = gateway.NewClient(gateway.WithChain(testEnv))
	default:
		testConfig.client = gateway.NewClient(gateway.WithChain(testEnv))
	}
	t.Cleanup(func() {
	})
	return &testConfig
}

// TestGateway checks the gateway can be accessed.
func TestGateway(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		BlockHash string
	}
	testSet := map[string][]testSetType{
		"devnet":  {},
		"mainnet": {{BlockHash: "0x4ee4c886d1767b7165a1e3a7c6ad145543988465f2bda680c16a79536f6d81f"}},
		"mock":    {{BlockHash: "0xdeadbeef"}},
		"testnet": {{BlockHash: "0x787af09f1cacdc5de1df83e8cdca3a48c1194171c742e78a9f684cb7aa4db"}},
	}[testEnv]

	for _, test := range testSet {
		block, err := testConfig.client.Block(context.Background(), &gateway.BlockOptions{BlockHash: test.BlockHash})

		if err != nil {
			t.Fatal(err)
		}
		if block.BlockHash != test.BlockHash {
			t.Fatalf("expecting %s, instead: %s", "", block.BlockHash)
		}
	}
}

func TestContractAddresses(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		Starknet             string
		GpsStatementVerifier string
	}
	testSet := map[string][]testSetType{
		"devnet": {},
		"mock":   {},
		"testnet": {
			{
				Starknet:             "0xde29d060d45901fb19ed6c6e959eb22d8626708e",
				GpsStatementVerifier: "0xab43ba48c9edf4c2c4bb01237348d1d7b28ef168",
			},
		},
		"mainnet": {
			{
				Starknet:             "0xc662c410c0ecf747543f5ba90660f6abebd9c8c4",
				GpsStatementVerifier: "0x47312450b3ac8b5b8e247a6bb6d523e7605bdb60",
			},
		},
	}[testEnv]

	for _, test := range testSet {
		addresses, err := testConfig.client.ContractAddresses(context.Background())
		if err != nil {
			t.Errorf("could not get starknet addresses - \n%v\n", err)
		}

		if strings.ToLower(addresses.Starknet) != test.Starknet {
			t.Errorf("fetched incorrect addresses - \n%s %s\n", strings.ToLower(addresses.Starknet), test.Starknet)
		}
	}
}
