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
	TransactionHash(calls rpc.FunctionCall, txDetails rpc.TxDetails) (*felt.Felt, error)
	Call(ctx context.Context, call rpc.FunctionCall, blockId rpc.BlockID) ([]*felt.Felt, error)
	Nonce(ctx context.Context) (*felt.Felt, error)
	Sign(ctx context.Context, msg *felt.Felt) ([]*felt.Felt, error)
	SignInvokeTransaction(ctx context.Context, invokeTx *rpc.BroadcastedInvokeV1Transaction) error
	EstimateFee(ctx context.Context, broadcastTxs []rpc.BroadcastedTransaction, blockId rpc.BlockID) ([]rpc.FeeEstimate, error)
	AddInvokeTransaction(ctx context.Context, invokeTx *rpc.BroadcastedInvokeV1Transaction) (*rpc.AddInvokeTransactionResponse, error)
	Execute(ctx context.Context, invokeTx *rpc.BroadcastedInvokeV1Transaction) (*rpc.AddInvokeTransactionResponse, error) // Todo: generalise once rpcv04 PRs are merged
}

var _ AccountInterface = &Account{}

type Account struct {
	provider       rpc.RpcProvider
	ChainId        *felt.Felt
	accountAddress *felt.Felt
	ks             starknetgo.Keystore
	version        uint64
}

func NewAccount(provider rpc.RpcProvider, version uint64, accountAddress *felt.Felt, keystore starknetgo.Keystore, setChainId bool) (*Account, error) {
	account := &Account{
		provider:       provider,
		accountAddress: accountAddress,
		ks:             keystore,
		version:        version,
	}

	if setChainId == false {
		return account, nil
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

func (account *Account) TransactionHash(call rpc.FunctionCall, txDetails rpc.TxDetails) (*felt.Felt, error) {

	if call.Calldata == nil || txDetails.Nonce == nil || txDetails.MaxFee == nil || account.accountAddress == nil {
		return nil, ErrNotAllParametersSet
	}

	calldataHash, err := computeHashOnElementsFelt(call.Calldata)
	if err != nil {
		return nil, err
	}

	return calculateTransactionHashCommon(
		new(felt.Felt).SetBytes([]byte(TRANSACTION_PREFIX)),
		new(felt.Felt).SetUint64(account.version),
		account.accountAddress,
		&felt.Zero,
		calldataHash,
		txDetails.MaxFee,
		account.ChainId,
		[]*felt.Felt{txDetails.Nonce},
	)
}

func (account *Account) TransactionHash2(callData []*felt.Felt, txDetails rpc.TxDetails) (*felt.Felt, error) {

	if len(callData) == 0 || txDetails.Nonce == nil || txDetails.MaxFee == nil || account.accountAddress == nil {
		return nil, ErrNotAllParametersSet
	}

	calldataHash, err := computeHashOnElementsFelt(callData)
	if err != nil {
		return nil, err
	}

	return calculateTransactionHashCommon(
		new(felt.Felt).SetBytes([]byte(TRANSACTION_PREFIX)),
		new(felt.Felt).SetUint64(account.version),
		account.accountAddress,
		&felt.Zero,
		calldataHash,
		txDetails.MaxFee,
		account.ChainId,
		[]*felt.Felt{txDetails.Nonce},
	)
}

func (account *Account) Nonce(ctx context.Context) (*felt.Felt, error) {
	switch account.version {
	case 1:
		// Todo: simplfy after rpc PRs are merged, return account.provider.Nonce(...)
		nonce, err := account.provider.Nonce(ctx, rpc.WithBlockTag("latest"), account.accountAddress)
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
	s1, s2, err := account.ks.Sign(ctx, account.accountAddress.String(), msgBig)
	if err != nil {
		return nil, err
	}
	s1Felt, err := utils.BigIntToFelt(s1)
	if err != nil {
		return nil, err
	}
	s2Felt, err := utils.BigIntToFelt(s2)
	if err != nil {
		return nil, err
	}
	return []*felt.Felt{s1Felt, s2Felt}, nil
}

func (account *Account) SignInvokeTransaction(ctx context.Context, invokeTx *rpc.BroadcastedInvokeV1Transaction) error {
	txHash, err := account.TransactionHash2(invokeTx.Calldata, rpc.TxDetails{Nonce: invokeTx.Nonce, MaxFee: invokeTx.MaxFee})
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

func (account *Account) EstimateFee(ctx context.Context, broadcastTxs []rpc.BroadcastedTransaction, blockId rpc.BlockID) ([]rpc.FeeEstimate, error) {
	switch account.version {
	case 1:
		return account.provider.EstimateFee(ctx, broadcastTxs, blockId)
	default:
		return nil, ErrAccountVersionNotSupported
	}
}

// AddInvokeTransaction submits an invoke transaction to the rpc provider.
func (account *Account) AddInvokeTransaction(ctx context.Context, invokeTx *rpc.BroadcastedInvokeV1Transaction) (*rpc.AddInvokeTransactionResponse, error) {
	switch account.version {
	case 1:
		return account.provider.AddInvokeTransaction(ctx, invokeTx)
	default:
		return nil, ErrAccountVersionNotSupported
	}
}

// Execute Sets maxFee to twice the estimated fee (if not already set), sets the nonce, calculates the transaction hash, signs the transaction
// and finally submits an addInvokeTransaction to the rpc provider.
func (account *Account) Execute(ctx context.Context, invokeTx *rpc.BroadcastedInvokeV1Transaction) (*rpc.AddInvokeTransactionResponse, error) {
	if account.version != 1 {
		return nil, ErrAccountVersionNotSupported
	}

	// Set max fee if not already set
	if invokeTx.MaxFee == nil {
		estimate, err := account.EstimateFee(ctx, []rpc.BroadcastedTransaction{invokeTx}, rpc.WithBlockTag("latest"))
		if err != nil {
			return nil, err
		}
		overallFee, err := new(felt.Felt).SetString(string(estimate[0].OverallFee))
		if err != nil {
			return nil, err
		}
		newMaxFee := new(felt.Felt).Mul(overallFee, new(felt.Felt).SetUint64(2))
		invokeTx.MaxFee = newMaxFee
	}

	// Get and set nonce
	nonce, err := account.Nonce(ctx)
	if err != nil {
		return nil, err
	}
	invokeTx.Nonce = nonce

	// Sign transaction
	err = account.SignInvokeTransaction(ctx, invokeTx)
	if err != nil {
		return nil, err
	}

	// Submit transaction
	switch account.version {
	case 1:
		return account.AddInvokeTransaction(ctx, invokeTx)
	default:
		return nil, ErrAccountVersionNotSupported
	}
}
