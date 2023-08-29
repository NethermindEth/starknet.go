package account

import (
	"context"

	"github.com/NethermindEth/juno/core/felt"
	starknetgo "github.com/NethermindEth/starknet.go"
	"github.com/NethermindEth/starknet.go/rpc"
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
	// Nonce(ctx context.Context) (*felt.Felt, error)
	// EstimateFee(ctx context.Context, calls []rpc.FunctionCall) (*rpc.FeeEstimate, error)
	// Execute(ctx context.Context, calls []rpc.FunctionCall) (*rpc.AddInvokeTransactionResponse, error)
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
