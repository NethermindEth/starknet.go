package main

import (
	"context"
	"fmt"
	"os"

	"github.com/NethermindEth/juno/core/felt"
	starknetgo "github.com/NethermindEth/starknet.go"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/joho/godotenv"
)

var (
	network              string = "testnet"
	predeployedClassHash        = "0x2794ce20e5f2ff0d40e632cb53845b9f4e526ebd8471983f7dbd355b721d5a"
)

func main() {
	// Initialise the client.
	godotenv.Load(fmt.Sprintf(".env.%s", network))
	base := os.Getenv("INTEGRATION_BASE")
	c, err := ethrpc.DialContext(context.Background(), base)
	if err != nil {
		panic("You need to specify the testnet url in .env.testnet")
	}
	clientv02 := rpc.NewProvider(c)

	// Get keys
	pub, priv := getRandomKeys()

	classHash, err := utils.HexToFelt(predeployedClassHash)
	if err != nil {
		panic(err)
	}

	tx := rpc.BroadcastedDeployAccountTransaction{
		BroadcastedTxnCommonProperties: rpc.BroadcastedTxnCommonProperties{
			Nonce:     &felt.Zero, // Contract accounts start with nonce zero.
			MaxFee:    new(felt.Felt).SetUint64(4724395326064),
			Type:      rpc.TransactionType_DeployAccount,
			Version:   rpc.TransactionV1,
			Signature: []*felt.Felt{},
		},
		ClassHash:           classHash,
		ContractAddressSalt: pub,
		ConstructorCalldata: []*felt.Felt{pub},
	}

	precomputedAddress, err := precomputeAddress(&felt.Zero, pub, classHash, tx.ConstructorCalldata)
	fmt.Println("precomputedAddress:", precomputedAddress)

	// At this point you need to add funds to precomputed address to use it.
	var input string

	fmt.Println("The `precomputedAddress` account needs to have enough ETH to perform a transaction.")
	fmt.Println("Use the starknet faucet to send ETH to your `precomputedAddress`")
	fmt.Println("When your account has been funded by the faucet, press any key, then `enter` to continue : ")
	fmt.Scan(&input)

	// Get the chainID to sign the transaction
	chainId, err := clientv02.ChainID(context.Background())
	if err != nil {
		panic(err)
	}

	// Calculate and sign the transaction hash
	hash, err := calculateDeployAccountTransactionHash(tx, precomputedAddress, chainId)
	if err != nil {
		panic(err)
	}
	fmt.Println("Transaction hash:", hash)
	x, y, err := starknetgo.Curve.SignFelt(hash, priv)
	if err != nil {
		panic(err)
	}
	tx.Signature = []*felt.Felt{x, y}

	// Send transaction to the network
	resp, err := clientv02.AddDeployAccountTransaction(context.Background(), tx)
	if err != nil {
		panic(fmt.Sprintf("Error returned from AddDeployAccountTransaction: %+v", err))
	}

	fmt.Println("AddDeployAccountTransaction response:", resp)

}

func getRandomKeys() (*felt.Felt, *felt.Felt) {
	privateKey, err := starknetgo.Curve.GetRandomPrivateKey()
	if err != nil {
		fmt.Println("can't get random private key:", err)
		os.Exit(1)
	}
	pubX, _, err := starknetgo.Curve.PrivateToPoint(privateKey)
	if err != nil {
		fmt.Println("can't generate public key:", err)
		os.Exit(1)
	}
	privFelt, err := utils.BigIntToFelt(privateKey)
	if err != nil {
		panic(err)
	}
	pubFelt, err := utils.BigIntToFelt(pubX)
	if err != nil {
		panic(err)
	}
	return pubFelt, privFelt
}

// precomputeAddress computes the address by hashing the relevant data.
// ref: https://github.com/starkware-libs/cairo-lang/blob/master/src/starkware/starknet/core/os/contract_address/contract_address.py
// TODO: Move to contract / utils package
func precomputeAddress(deployerAddress *felt.Felt, salt *felt.Felt, classHash *felt.Felt, constructorCalldata []*felt.Felt) (*felt.Felt, error) {
	CONTRACT_ADDRESS_PREFIX := new(felt.Felt).SetBytes([]byte("STARKNET_CONTRACT_ADDRESS"))

	bigIntArr, err := utils.FeltArrToBigIntArr([]*felt.Felt{
		CONTRACT_ADDRESS_PREFIX,
		deployerAddress,
		salt,
		classHash,
	})
	if err != nil {
		return nil, err
	}

	constructorCalldataBigIntArr, err := utils.FeltArrToBigIntArr(constructorCalldata)
	constructorCallDataHashInt, _ := starknetgo.Curve.ComputeHashOnElements(*constructorCalldataBigIntArr)
	*bigIntArr = append(*bigIntArr, constructorCallDataHashInt)

	preBigInt, err := starknetgo.Curve.ComputeHashOnElements(*bigIntArr)
	if err != nil {
		return nil, err
	}
	return utils.BigIntToFelt(preBigInt)

}

func computeHashOnElementsFelt(feltArr []*felt.Felt) (*felt.Felt, error) {
	bigIntArr, err := utils.FeltArrToBigIntArr(feltArr)
	if err != nil {
		return nil, err
	}
	hash, err := starknetgo.Curve.ComputeHashOnElements(*bigIntArr)
	if err != nil {
		return nil, err
	}
	return utils.BigIntToFelt(hash)
}

// calculateDeployAccountTransactionHash computes the transaction hash for deployAccount transactions
func calculateDeployAccountTransactionHash(tx rpc.BroadcastedDeployAccountTransaction, contractAddress *felt.Felt, chainID string) (*felt.Felt, error) {
	Prefix_DEPLOY_ACCOUNT := new(felt.Felt).SetBytes([]byte("deploy_account"))
	chainIdFelt := new(felt.Felt).SetBytes([]byte(chainID))

	calldata := []*felt.Felt{tx.ClassHash, tx.ContractAddressSalt}
	calldata = append(calldata, tx.ConstructorCalldata...)
	calldataHash, err := computeHashOnElementsFelt(calldata)
	if err != nil {
		return nil, err
	}

	versionFelt, err := utils.HexToFelt(string(tx.Version))
	if err != nil {
		return nil, err
	}

	return calculateTransactionHashCommon(
		Prefix_DEPLOY_ACCOUNT,
		versionFelt,
		contractAddress,
		&felt.Zero,
		calldataHash,
		tx.MaxFee,
		chainIdFelt,
		[]*felt.Felt{tx.Nonce},
	)
}

func calculateTransactionHashCommon(
	txHashPrefix *felt.Felt,
	version *felt.Felt,
	contractAddress *felt.Felt,
	entryPointSelector *felt.Felt,
	calldata *felt.Felt,
	maxFee *felt.Felt,
	chainId *felt.Felt,
	additionalData []*felt.Felt) (*felt.Felt, error) {

	dataToHash := []*felt.Felt{
		txHashPrefix,
		version,
		contractAddress,
		entryPointSelector,
		calldata,
		maxFee,
		chainId,
	}
	dataToHash = append(dataToHash, additionalData...)

	return computeHashOnElementsFelt(dataToHash)
}
