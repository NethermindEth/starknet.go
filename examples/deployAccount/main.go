package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/joho/godotenv"
)

var (
	network              string = "testnet"
	predeployedClassHash        = "0x2794ce20e5f2ff0d40e632cb53845b9f4e526ebd8471983f7dbd355b721d5a"
	accountAddress              = "0xdeadbeef"
)

// main initializes the client, sets up the account, deploys a contract, and sends a transaction to the network.
//
// It loads environment variables, dials the Ethereum RPC, creates a new account, casts the account address to a felt type,
// sets up the account using the client, converts the predeployed class hash to a felt type, creates transaction data,
// precomputes an address, prompts the user to add funds to the precomputed address, signs the transaction,
// and finally sends the transaction to the network.
//
// Parameters:
//   none
// Returns:
//  none
func main() {
	// Initialise the client.
	godotenv.Load(fmt.Sprintf(".env.%s", network))
	base := os.Getenv("INTEGRATION_BASE")
	c, err := ethrpc.DialContext(context.Background(), base)
	if err != nil {
		panic("You need to specify the testnet url in .env.testnet")
	}
	clientv02 := rpc.NewProvider(c)

	// Get random keys for test purposes
	ks, pub, _ := account.GetRandomKeys()

	accountAddressFelt, err := new(felt.Felt).SetString(accountAddress)
	if err != nil {
		panic("Error casting accountAddress to felt")
	}

	// Set up account
	acnt, err := account.NewAccount(clientv02, accountAddressFelt, pub.String(), ks)
	if err != nil {
		panic(err)
	}

	classHash, err := utils.HexToFelt(predeployedClassHash)
	if err != nil {
		panic(err)
	}

	precomputedAddress, err := acnt.PrecomputeAddress(&felt.Zero, pub, classHash, []*felt.Felt{pub})

	fmt.Printf("\nIn order to deploy your account (address %s), you need to fund the acccount (using a faucet), and then press `enter` to continue : \n", precomputedAddress.String())

	reader := bufio.NewReader(os.Stdin)
	_, err = reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("Waiting for deployment")
	deployOptions := account.DeployOptions{
		ClassHash:       classHash,
		MaxFee:          new(felt.Felt).SetUint64(4724395326064),
		DeploytWaitTime: 2 * time.Second,
	}

	// Deploy the account
	resp, err := acnt.DeployAccount(deployOptions)

	if err != nil {
		panic(fmt.Sprint("Error returned from DeployAccount: ", err))
	}
	fmt.Println("Deployed with response response:", resp)

}
