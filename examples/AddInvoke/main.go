package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/NethermindEth/juno/core/felt"
	starknetgo "github.com/NethermindEth/starknet.go"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/types"
	"github.com/NethermindEth/starknet.go/utils"
	"github.com/joho/godotenv"
)

var (
	name                string = "testnet"
	someMainnetContract string = "0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"
	contractMethod      string = "approve"
)

func main() {
	fmt.Println("Starting simpeCall example")
	godotenv.Load(fmt.Sprintf(".env.%s", name))
	// base := os.Getenv("INTEGRATION_BASE")
	// c, err := ethrpc.DialContext(context.Background(), base)
	// if err != nil {
	// 	fmt.Println("Failed to connect to the client, did you specify the url in the .env.mainnet?")
	// 	panic(err)
	// }
	// clientv02 := rpc.NewProvider(c)
	fmt.Println("Established connection with the client")

	contractAddress, err := utils.HexToFelt(someMainnetContract)
	if err != nil {
		panic(err)
	}
	accountAddress, _ := new(felt.Felt).SetString("0x043784df59268c02b716e20bf77797bd96c68c2f100b2a634e448c35e3ad363e")

	// Now build the trasaction
	invokeTx := rpc.BroadcastedInvokeV1Transaction{
		BroadcastedTxnCommonProperties: rpc.BroadcastedTxnCommonProperties{
			Nonce:   new(felt.Felt).SetUint64(2), // Likely incorrect - but the seq should tell us this
			MaxFee:  new(felt.Felt).SetUint64(1),
			Version: rpc.TransactionV1,
			Type:    rpc.TransactionType_Invoke,
		},
		SenderAddress: accountAddress,
	}
	//calldata
	// Make read contract call
	tx := rpc.FunctionCall{
		ContractAddress:    contractAddress,
		EntryPointSelector: types.GetSelectorFromNameFelt(contractMethod),
		Calldata:           []*felt.Felt{contractAddress, new(felt.Felt).SetUint64(1)},
	}
	invokeTx.Calldata = fmtCalldata([]rpc.FunctionCall{tx})

	// transaction calldata matches feeder block, and voyager, 0x6a2932d91197a5475488b2f917ebc7aeeaa066b616bc96601f46f006d27ff0c
	// => transaction hash is incorrect (not fmtCallData). Alos know signature is correct from postman.

	// sign
	fakePrivKeyFelt, _ := new(felt.Felt).SetString("0x043b7fe9d91942c98cd5fd37579bd99ec74f879c4c79d886633eecae9dad35fa")
	txHash, err := TransactionHash(
		invokeTx.Calldata,
		rpc.TxDetails{
			Nonce:   invokeTx.Nonce,
			MaxFee:  invokeTx.MaxFee,
			Version: invokeTx.Version,
		},
		accountAddress,
	)
	x, y, err := starknetgo.Curve.SignFelt(txHash, fakePrivKeyFelt)
	if err != nil {
		panic(err)
	}
	invokeTx.Signature = []*felt.Felt{x, y}

	fmt.Println("Making Call() request")
	qwe, _ := json.MarshalIndent(invokeTx, "", "")
	fmt.Println(string(qwe))

	pub, _ := new(big.Int).SetString("2090221843434510384432085791482977629840322403554658343615172301617258923551", 0)
	hash, _ := new(big.Int).SetString("2391207323525339369856503773624499041713147169482476892458076737242741151771", 0)

	fmt.Println("txHash", txHash) // It seems the txhash is incorrect. Signature is correct.
	fmt.Println("acntadr", accountAddress, accountAddress.BigInt(new(big.Int)))
	fmt.Println("pub", pub, new(felt.Felt).SetBytes(pub.Bytes()))
	fmt.Println("hash", hash, new(felt.Felt).SetBytes(hash.Bytes()))
	fmt.Println(invokeTx.Signature[0].BigInt(new(big.Int)))
	fmt.Println(invokeTx.Signature[1].BigInt(new(big.Int)))
	// callResp, err := clientv02.AddInvokeTransaction(context.Background(), invokeTx)
	// if err != nil {
	// 	fmt.Println("=======")
	// 	panic(err.Error())
	// }

	// fmt.Println(fmt.Sprintf("Response to %s():%s ", contractMethod, callResp))
}
func fmtCalldata(fnCalls []rpc.FunctionCall) []*felt.Felt {
	callArray := []*felt.Felt{}
	callData := []*felt.Felt{new(felt.Felt).SetUint64(uint64(len(fnCalls)))}

	for _, tx := range fnCalls {
		callData = append(callData, tx.ContractAddress, tx.EntryPointSelector)

		if len(tx.Calldata) == 0 {
			callData = append(callData, &felt.Zero, &felt.Zero)
			continue
		}

		callData = append(callData, new(felt.Felt).SetUint64(uint64(len(callArray))), new(felt.Felt).SetUint64(uint64(len(tx.Calldata)+1)))
		// callData = append(callData, new(felt.Felt).SetUint64(uint64(len(callArray))), new(felt.Felt).SetUint64(uint64(len(tx.Calldata))))
		for _, cd := range tx.Calldata {
			callArray = append(callArray, cd)
		}

	}
	callData = append(callData, new(felt.Felt).SetUint64(uint64(len(callArray)+1)))
	callData = append(callData, callArray...)
	callData = append(callData, new(felt.Felt).SetUint64(uint64(0)))
	return callData
}

// computeHashOnElementsFelt hashes the array of felts provided as input
func computeHashOnElementsFelt(feltArr []*felt.Felt) (*felt.Felt, error) {
	bigIntArr, err := utils.FeltArrToBigIntArr(feltArr)
	if err != nil {
		return nil, err
	}
	hash, err := starknetgo.Curve.ComputeHashOnElements(*bigIntArr)
	if err != nil {
		return nil, err
	}
	fmt.Println("utils.BigIntToFelt(hash)", hash)
	return utils.BigIntToFelt(hash)
}

// calculateTransactionHashCommon [specification] calculates the transaction hash in the StarkNet network - a unique identifier of the transaction.
// [specification]: https://github.com/starkware-libs/cairo-lang/blob/master/src/starkware/starknet/core/os/transaction_hash/transaction_hash.py#L27C5-L27C38
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

func TransactionHash(formattedcallData []*felt.Felt, txDetails rpc.TxDetails, AccountAddress *felt.Felt) (*felt.Felt, error) {

	if txDetails.Nonce == nil || txDetails.MaxFee == nil {
		return nil, errors.New("qweqweqwe")
	}

	calldataHash, err := computeHashOnElementsFelt(formattedcallData)
	if err != nil {
		return nil, err
	}
	fmt.Println("prefix", new(felt.Felt).SetBytes([]byte("invoke")))
	fmt.Println("chain id", new(felt.Felt).SetBytes([]byte("SN_GOERLI")))
	return calculateTransactionHashCommon(
		new(felt.Felt).SetBytes([]byte("invoke")),
		new(felt.Felt).SetUint64(1),
		AccountAddress,
		&felt.Zero,
		calldataHash,
		txDetails.MaxFee,
		new(felt.Felt).SetBytes([]byte("SN_GOERLI")),
		[]*felt.Felt{txDetails.Nonce},
	)
}
