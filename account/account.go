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
	TransactionHash2(callData []*felt.Felt, nonce *felt.Felt, maxFee *felt.Felt, accountAddress *felt.Felt) (*felt.Felt, error)
	Call(ctx context.Context, call rpc.FunctionCall, blockId rpc.BlockID) ([]*felt.Felt, error)
	Nonce(ctx context.Context) (*felt.Felt, error)
	Sign(ctx context.Context, msg *felt.Felt) ([]*felt.Felt, error)
	SignInvokeTransaction(ctx context.Context, invokeTx *rpc.BroadcastedInvokeV1Transaction) error
	EstimateFee(ctx context.Context, broadcastTxs []rpc.BroadcastedTransaction, blockId rpc.BlockID) ([]rpc.FeeEstimate, error)
	AddInvokeTransaction(ctx context.Context, invokeTx *rpc.BroadcastedInvokeV1Transaction) (*rpc.AddInvokeTransactionResponse, error)
	BuildInvokeTx(ctx context.Context, invokeTx *rpc.BroadcastedInvokeV1Transaction, fnCall *[]rpc.FunctionCall) error
}

var _ AccountInterface = &Account{}

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

func (account *Account) Call(ctx context.Context, call rpc.FunctionCall, blockId rpc.BlockID) ([]*felt.Felt, error) {
	return account.provider.Call(ctx,
		rpc.FunctionCall{
			ContractAddress:    call.ContractAddress,
			EntryPointSelector: call.EntryPointSelector,
			Calldata:           call.Calldata},
		blockId)
}

// TransactionHash2 requires the callData to be compiled beforehand
func (account *Account) TransactionHash2(callData []*felt.Felt, nonce *felt.Felt, maxFee *felt.Felt, accountAddress *felt.Felt) (*felt.Felt, error) {

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

func (account *Account) Nonce(ctx context.Context) (*felt.Felt, error) {
	switch account.version {
	case 1:
		// Todo: simplfy after rpc PRs are merged, return account.provider.Nonce(...)
		nonce, err := account.provider.Nonce(ctx, rpc.WithBlockTag("latest"), account.AccountAddress)
		if err != nil {
			return nil, err
		}
		return new(felt.Felt).SetString(*nonce)
	default:
		return nil, ErrAccountVersionNotSupported
	}
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

	txHash, err := account.TransactionHash2(invokeTx.Calldata, invokeTx.Nonce, invokeTx.MaxFee, account.AccountAddress)
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

// BuildInvokeTx Sets maxFee to twice the estimated fee (if not already set), compiles and sets the CallData, calculates the transaction hash, signs the transaction.
func (account *Account) BuildInvokeTx(ctx context.Context, invokeTx *rpc.BroadcastedInvokeV1Transaction, fnCall *[]rpc.FunctionCall) error {
	if account.version != 1 {
		return ErrAccountVersionNotSupported
	}

	// Set max fee if not already set
	// if invokeTx.MaxFee == nil {
	// 	estimate, err := account.EstimateFee(ctx, []rpc.BroadcastedTransaction{invokeTx}, rpc.WithBlockTag("latest"))
	// 	if err != nil {
	// 		return err
	// 	}
	// 	overallFee, err := new(felt.Felt).SetString(string(estimate[0].OverallFee))
	// 	if err != nil {
	// 		return err
	// 	}
	// 	newMaxFee := new(felt.Felt).Mul(overallFee, new(felt.Felt).SetUint64(2))
	// 	invokeTx.MaxFee = newMaxFee
	// }
	// Compile callData
	// invokeTx.Calldata = fmtCalldata(*fnCall)
	// Get and set nonce
	// nonce, err := account.Nonce(ctx)
	// if err != nil {
	// 	return nil, err
	// }
	// invokeTx.Nonce = nonce

	// Sign transaction
	err := account.SignInvokeTransaction(ctx, invokeTx)
	if err != nil {
		return err
	}
	return nil
}

func (account *Account) EstimateFee(ctx context.Context, broadcastTxs []rpc.BroadcastedTransaction, blockId rpc.BlockID) ([]rpc.FeeEstimate, error) {
	switch account.version {
	case 1:
		return account.provider.EstimateFee(ctx, broadcastTxs, blockId)
	default:
		return nil, ErrAccountVersionNotSupported
	}
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
