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
	Call(ctx context.Context, call rpc.FunctionCall) ([]*felt.Felt, error)
	Execute(ctx context.Context, calls rpc.FunctionCall, details rpc.TxDetails) (*rpc.AddInvokeTransactionResponse, error)
	Nonce(ctx context.Context) (*felt.Felt, error)
	Sign(ctx context.Context, msg *felt.Felt) ([]*felt.Felt, error)
	// EstimateFee(ctx context.Context, calls []rpc.FunctionCall) (*rpc.FeeEstimate, error)
	// Declare(ctx context.Context, classHash string, contract rpc.ContractClass) (rpc.AddDeclareTransactionResponse, error)
	// DeployAccount(ctx context.Context, classHash string) (*rpc.AddDeployTransactionResponse, error) // ToDo: Should be AddDeployAccountTransactionResponse - waiting for PR to be merged
}

var _ AccountInterface = &Account{}

type Account struct {
	provider       *rpc.Provider
	chainId        string
	AccountAddress *felt.Felt
	ks             starknetgo.Keystore
	version        uint64
}

func NewAccount(provider *rpc.Provider, version uint64, accountAddress *felt.Felt, keystore starknetgo.Keystore) (*Account, error) {
	account := &Account{
		provider:       provider,
		AccountAddress: accountAddress,
		ks:             keystore,
		version:        version,
	}

	chainID, err := provider.ChainID(context.Background())
	if err != nil {
		return nil, err
	}
	account.chainId = chainID

	return account, nil
}

func (account *Account) Call(ctx context.Context, call rpc.FunctionCall) ([]*felt.Felt, error) {
	return account.provider.Call(ctx,
		rpc.FunctionCall{
			ContractAddress:    call.ContractAddress,
			EntryPointSelector: call.EntryPointSelector,
			Calldata:           call.Calldata},
		rpc.WithBlockTag("latest"))
}

func (account *Account) TransactionHash(call rpc.FunctionCall, txDetails rpc.TxDetails) (*felt.Felt, error) {
	calldataHash, err := computeHashOnElementsFelt(call.Calldata)
	if err != nil {
		return nil, err
	}
	return calculateTransactionHashCommon(
		new(felt.Felt).SetBytes([]byte(TRANSACTION_PREFIX)),
		new(felt.Felt).SetUint64(account.version),
		account.AccountAddress,
		&felt.Zero,
		calldataHash,
		txDetails.MaxFee,
		new(felt.Felt).SetBytes([]byte(account.chainId)),
		[]*felt.Felt{txDetails.Nonce},
	)
}

func (account *Account) Nonce(ctx context.Context) (*felt.Felt, error) {
	switch account.version {
	case 1:
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
	s1, s2, err := account.ks.Sign(ctx, account.AccountAddress.String(), msgBig)
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

func (account *Account) Execute(ctx context.Context, calls rpc.FunctionCall, details rpc.TxDetails) (*rpc.AddInvokeTransactionResponse, error) {
	switch account.version {
	case 1:
		txHash, err := account.TransactionHash(calls, details)
		if err != nil {
			return nil, err
		}
		signature, err := account.Sign(ctx, txHash)
		if err != nil {
			return nil, err
		}
		resp, err := account.provider.AddInvokeTransaction(
			ctx,
			rpc.BroadcastedInvokeV1Transaction{
				BroadcastedTxnCommonProperties: rpc.BroadcastedTxnCommonProperties{
					MaxFee:    details.MaxFee,
					Version:   rpc.TransactionV1,
					Signature: signature,
					Nonce:     details.Nonce,
					Type:      "INVOKE",
				},
				SenderAddress: account.AccountAddress,
				Calldata:      calls.Calldata,
			})
		if err != nil {
			return nil, err
		}
		return &rpc.AddInvokeTransactionResponse{TransactionHash: resp.TransactionHash}, nil
	default:
		return nil, ErrAccountVersionNotSupported
	}
}
