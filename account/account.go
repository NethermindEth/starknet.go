package account

import (
	"context"
	"errors"

	"github.com/NethermindEth/juno/core/felt"
	starknetgo "github.com/NethermindEth/starknet.go"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"
)

var (
	ErrAccountVersionNotSupported = errors.New("Account version not supported")
	ErrNotAllParametersSet        = errors.New("Not all neccessary parameters have been set")
	ErrTxnTypeUnSupported         = errors.New("Unsupported transction type")
	ErrFeltToBigInt               = errors.New("Felt to BigInt error")
)

const (
	TRANSACTION_PREFIX      = "invoke"
	DECLARE_PREFIX          = "declare"
	EXECUTE_SELECTOR        = "__execute__"
	CONTRACT_ADDRESS_PREFIX = "STARKNET_CONTRACT_ADDRESS"
)

//go:generate mockgen -destination=../mocks/mock_account.go -package=mocks -source=account.go AccountInterface
type AccountInterface interface {
	Sign(ctx context.Context, msg *felt.Felt) ([]*felt.Felt, error)
	BuildInvokeTx(ctx context.Context, invokeTx *rpc.BroadcastedInvokeV1Transaction, fnCall *[]rpc.FunctionCall) error
	TransactionHashInvoke(callData []*felt.Felt, nonce *felt.Felt, maxFee *felt.Felt, accountAddress *felt.Felt) (*felt.Felt, error)
	TransactionHashDeployAccount(tx rpc.BroadcastedDeployAccountTransaction, contractAddress *felt.Felt) (*felt.Felt, error)
	SignInvokeTransaction(ctx context.Context, tx *rpc.BroadcastedInvokeV1Transaction) error
	SignDeployAccountTransaction(ctx context.Context, tx *rpc.BroadcastedDeployAccountTransaction, precomputeAddress *felt.Felt) error
	AddInvokeTransaction(ctx context.Context, invokeTx *rpc.BroadcastedInvokeV1Transaction) (*rpc.AddInvokeTransactionResponse, error) //todo: post rpcv04 merge
}

var _ AccountInterface = &Account{}

// var _ rpc.RpcProvider = &Account{} //todo: post rpcv04 merge

type Account struct {
	provider       rpc.RpcProvider
	ChainId        *felt.Felt
	AccountAddress *felt.Felt
	publicKey      string
	ks             starknetgo.Keystore
	version        uint64
}

func NewAccount(provider rpc.RpcProvider, version uint64, accountAddress *felt.Felt, publicKey string, keystore starknetgo.Keystore) (*Account, error) {
	account := &Account{
		provider:       provider,
		AccountAddress: accountAddress,
		publicKey:      publicKey,
		ks:             keystore,
		version:        version,
	}

	chainID, err := provider.ChainID(context.Background())
	if err != nil {
		return nil, err
	}
	account.ChainId = new(felt.Felt).SetBytes([]byte(chainID))

	return account, nil
}

// TransactionHash2 requires the callData to be compiled beforehand
func (account *Account) TransactionHashInvoke(callData []*felt.Felt, nonce *felt.Felt, maxFee *felt.Felt, accountAddress *felt.Felt) (*felt.Felt, error) {

	if len(callData) == 0 || nonce == nil || maxFee == nil || accountAddress == nil {
		return nil, ErrNotAllParametersSet
	}
	calldataHash, err := computeHashOnElementsFelt(callData)
	if err != nil {
		return nil, err
	}
	return calculateTransactionHashCommon(
		new(felt.Felt).SetBytes([]byte(TRANSACTION_PREFIX)),
		new(felt.Felt).SetUint64(account.version),
		accountAddress,
		&felt.Zero,
		calldataHash,
		maxFee,
		account.ChainId,
		[]*felt.Felt{nonce},
	)
}

func (account *Account) Sign(ctx context.Context, msg *felt.Felt) ([]*felt.Felt, error) {

	msgBig, ok := utils.FeltToBigInt(msg)
	if ok != true {
		return nil, ErrFeltToBigInt
	}
	s1, s2, err := account.ks.Sign(ctx, account.publicKey, msgBig)
	if err != nil {
		return nil, err
	}
	s1Felt, _ := utils.BigIntToFelt(s1)
	s2Felt, _ := utils.BigIntToFelt(s2)

	return []*felt.Felt{s1Felt, s2Felt}, nil
}

func (account *Account) SignInvokeTransaction(ctx context.Context, invokeTx *rpc.BroadcastedInvokeV1Transaction) error {

	txHash, err := account.TransactionHashInvoke(invokeTx.Calldata, invokeTx.Nonce, invokeTx.MaxFee, account.AccountAddress)
	if err != nil {
		return err
	}
	signature, err := account.Sign(ctx, txHash)
	if err != nil {
		return err
	}
	invokeTx.Signature = signature
	return nil
}

func (account *Account) SignDeployAccountTransaction(ctx context.Context, tx *rpc.BroadcastedDeployAccountTransaction, precomputeAddress *felt.Felt) error {

	hash, err := account.TransactionHashDeployAccount(*tx, precomputeAddress)
	if err != nil {
		return err
	}
	signature, err := account.Sign(ctx, hash)
	if err != nil {
		return err
	}
	tx.Signature = signature
	return nil
}

// TransactionHashDeployAccount computes the transaction hash for deployAccount transactions
func (account *Account) TransactionHashDeployAccount(tx rpc.BroadcastedDeployAccountTransaction, contractAddress *felt.Felt) (*felt.Felt, error) {
	Prefix_DEPLOY_ACCOUNT := new(felt.Felt).SetBytes([]byte("deploy_account"))

	calldata := []*felt.Felt{tx.ClassHash, tx.ContractAddressSalt}
	calldata = append(calldata, tx.ConstructorCalldata...)
	calldataHash, err := computeHashOnElementsFelt(calldata)
	if err != nil {
		return nil, err
	}

	versionFelt, err := new(felt.Felt).SetString(string(tx.Version))
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
		account.ChainId,
		[]*felt.Felt{tx.Nonce},
	)
}

// BuildInvokeTx Sets maxFee to twice the estimated fee (if not already set), compiles and sets the CallData, calculates the transaction hash, signs the transaction.
func (account *Account) BuildInvokeTx(ctx context.Context, invokeTx *rpc.BroadcastedInvokeV1Transaction, fnCall *[]rpc.FunctionCall) error {
	if account.version != 1 {
		return ErrAccountVersionNotSupported
	}

	invokeTx.Calldata = FmtCalldata(*fnCall)

	return account.SignInvokeTransaction(ctx, invokeTx)
}

// AddInvokeTransaction submits a complete (ie signed, and calldata has been formatted etc) invoke transaction to the rpc provider.
func (account *Account) AddInvokeTransaction(ctx context.Context, invokeTx *rpc.BroadcastedInvokeV1Transaction) (*rpc.AddInvokeTransactionResponse, error) {
	switch account.version {
	case 1:
		return account.provider.AddInvokeTransaction(ctx, invokeTx)
	default:
		return nil, ErrAccountVersionNotSupported
	}
}

/*
Formats the multicall transactions in a format which can be signed and verified by the network and OpenZeppelin account contracts
*/
func FmtCalldata(fnCalls []rpc.FunctionCall) []*felt.Felt {
	callArray := []*felt.Felt{}
	callData := []*felt.Felt{new(felt.Felt).SetUint64(uint64(len(fnCalls)))}

	for _, tx := range fnCalls {
		callData = append(callData, tx.ContractAddress, tx.EntryPointSelector)

		if len(tx.Calldata) == 0 {
			callData = append(callData, &felt.Zero, &felt.Zero)
			continue
		}

		callData = append(callData, new(felt.Felt).SetUint64(uint64(len(callArray))), new(felt.Felt).SetUint64(uint64(len(tx.Calldata))+1))
		for _, cd := range tx.Calldata {
			callArray = append(callArray, cd)
		}
	}
	callData = append(callData, new(felt.Felt).SetUint64(uint64(len(callArray)+1)))
	callData = append(callData, callArray...)
	callData = append(callData, new(felt.Felt).SetUint64(0))
	return callData
}

func (account *Account) BlockHashAndNumber(ctx context.Context) (*rpc.BlockHashAndNumberOutput, error) {
	return account.provider.BlockHashAndNumber(ctx)
}

func (account *Account) BlockNumber(ctx context.Context) (uint64, error) {
	return account.provider.BlockNumber(ctx)
}

func (account *Account) BlockTransactionCount(ctx context.Context, blockID rpc.BlockID) (uint64, error) {
	return account.provider.BlockTransactionCount(ctx, blockID)
}

func (account *Account) BlockWithTxHashes(ctx context.Context, blockID rpc.BlockID) (interface{}, error) {
	return account.provider.BlockWithTxHashes(ctx, blockID)
}

func (account *Account) BlockWithTxs(ctx context.Context, blockID rpc.BlockID) (interface{}, error) {
	return account.provider.BlockWithTxs(ctx, blockID)
}

func (account *Account) Call(ctx context.Context, call rpc.FunctionCall, blockId rpc.BlockID) ([]*felt.Felt, error) {
	return account.provider.Call(ctx, call, blockId)
}

func (account *Account) ChainID(ctx context.Context) (string, error) {
	return account.provider.ChainID(ctx)
}
func (account *Account) Class(ctx context.Context, blockID rpc.BlockID, classHash *felt.Felt) (rpc.ClassOutput, error) {
	return account.provider.Class(ctx, blockID, classHash)
}
func (account *Account) ClassAt(ctx context.Context, blockID rpc.BlockID, contractAddress *felt.Felt) (rpc.ClassOutput, error) {
	return account.provider.ClassAt(ctx, blockID, contractAddress)
}

func (account *Account) ClassHashAt(ctx context.Context, blockID rpc.BlockID, contractAddress *felt.Felt) (*felt.Felt, error) {
	return account.provider.ClassHashAt(ctx, blockID, contractAddress)
}

func (account *Account) EstimateFee(ctx context.Context, requests []rpc.BroadcastedTransaction, blockID rpc.BlockID) ([]rpc.FeeEstimate, error) {
	return account.provider.EstimateFee(ctx, requests, blockID)
}

func (account *Account) Events(ctx context.Context, input rpc.EventsInput) (*rpc.EventChunk, error) {
	return account.provider.Events(ctx, input)
}
func (account *Account) Nonce(ctx context.Context, blockID rpc.BlockID, contractAddress *felt.Felt) (*string, error) {
	return account.provider.Nonce(ctx, blockID, contractAddress)
}

func (account *Account) StateUpdate(ctx context.Context, blockID rpc.BlockID) (*rpc.StateUpdateOutput, error) {
	return account.provider.StateUpdate(ctx, blockID)
}
func (account *Account) Syncing(ctx context.Context) (*rpc.SyncStatus, error) {
	return account.provider.Syncing(ctx)
}
func (account *Account) TransactionByBlockIdAndIndex(ctx context.Context, blockID rpc.BlockID, index uint64) (rpc.Transaction, error) {
	return account.provider.TransactionByBlockIdAndIndex(ctx, blockID, index)
}
func (account *Account) TransactionByHash(ctx context.Context, hash *felt.Felt) (rpc.Transaction, error) {
	return account.provider.TransactionReceipt(ctx, hash)
}

func (account *Account) AddDeclareTransaction(ctx context.Context, declareTransaction rpc.BroadcastedDeclareTransaction) (*rpc.AddDeclareTransactionResponse, error) {
	switch account.version {
	case 1:
		return account.provider.AddDeclareTransaction(ctx, declareTransaction)
	default:
		return nil, ErrAccountVersionNotSupported
	}
}

func (account *Account) AddDeployAccountTransaction(ctx context.Context, deployAccountTransaction rpc.BroadcastedDeployAccountTransaction) (*rpc.AddDeployAccountTransactionResponse, error) {
	switch account.version {
	case 1:
		return account.provider.AddDeployAccountTransaction(ctx, deployAccountTransaction)
	default:
		return nil, ErrAccountVersionNotSupported
	}
}

// precomputeAddress computes the address by hashing the relevant data.
// ref: https://github.com/starkware-libs/cairo-lang/blob/master/src/starkware/starknet/core/os/contract_address/contract_address.py
func (account *Account) PrecomputeAddress(deployerAddress *felt.Felt, salt *felt.Felt, classHash *felt.Felt, constructorCalldata []*felt.Felt) (*felt.Felt, error) {
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
