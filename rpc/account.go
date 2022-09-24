package rpc

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/dontpanicdao/caigo"
	"github.com/dontpanicdao/caigo/rpc/types"
)

const (
	EXECUTE_SELECTOR   string = "__execute__"
	TRANSACTION_PREFIX string = "invoke"
)

type account interface {
	Sign(msgHash *big.Int) (*big.Int, *big.Int, error)
	TransactionHash(calls []types.FunctionCall, details types.ExecuteDetails) (*big.Int, error)
	Call(ctx context.Context, call types.FunctionCall) ([]string, error)
	Nonce(ctx context.Context) (*big.Int, error)
	EstimateFee(ctx context.Context, calls []types.FunctionCall, details types.ExecuteDetails) (*types.FeeEstimate, error)
	Execute(ctx context.Context, calls []types.FunctionCall, details types.ExecuteDetails) (*types.AddInvokeTransactionOutput, error)
}

var _ account = &Account{}

type Account struct {
	Provider *Provider
	Address  string
	private  *big.Int
	version  *big.Int
}

type accountOption struct {
	version *big.Int
}

type AccountOption func() accountOption

func AccountVersion0() accountOption {
	return accountOption{
		version: big.NewInt(0),
	}
}

func AccountVersion1() accountOption {
	return accountOption{
		version: big.NewInt(1),
	}
}

func (provider *Provider) NewAccount(private, address string, options ...AccountOption) (*Account, error) {
	version := big.NewInt(0)
	for _, o := range options {
		opt := o()
		if opt.version != nil {
			version = opt.version
		}
	}
	priv := caigo.SNValToBN(private)

	return &Account{
		Provider: provider,
		Address:  address,
		private:  priv,
		version:  version,
	}, nil
}

func (account *Account) Call(ctx context.Context, call types.FunctionCall) ([]string, error) {
	return account.Provider.Call(ctx, call, WithBlockTag("latest"))
}

func (account *Account) Sign(msgHash *big.Int) (*big.Int, *big.Int, error) {
	return caigo.Curve.Sign(msgHash, account.private)
}

func (account *Account) TransactionHash(calls []types.FunctionCall, details types.ExecuteDetails) (*big.Int, error) {
	chainID, err := account.Provider.ChainID(context.Background())
	if err != nil {
		return nil, err
	}

	callArray := fmtExecuteCalldata(details.Nonce, calls)
	cdHash, err := caigo.Curve.ComputeHashOnElements(callArray)
	if err != nil {
		return nil, err
	}

	multiHashData := []*big.Int{
		caigo.UTF8StrToBig(TRANSACTION_PREFIX),
		account.version,
		caigo.SNValToBN(account.Address),
		caigo.GetSelectorFromName(EXECUTE_SELECTOR),
		cdHash,
		details.MaxFee,
		caigo.UTF8StrToBig(chainID),
	}

	return caigo.Curve.ComputeHashOnElements(multiHashData)
}

func (account *Account) Nonce(ctx context.Context) (*big.Int, error) {
	nonce, err := account.Provider.Call(
		ctx,
		types.FunctionCall{
			ContractAddress:    types.HexToHash(account.Address),
			EntryPointSelector: "get_nonce",
			CallData:           []string{},
		},
		WithBlockTag("latest"),
	)
	if err != nil {
		return nil, err
	}
	if len(nonce) == 0 {
		return nil, errors.New("nonce error")
	}
	n, ok := big.NewInt(0).SetString(nonce[0], 0)
	if !ok {
		return nil, errors.New("nonce error")
	}
	return n, nil
}

func (account *Account) EstimateFee(ctx context.Context, calls []types.FunctionCall, details types.ExecuteDetails) (*types.FeeEstimate, error) {
	var err error
	nonce := details.Nonce
	if details.Nonce == nil {
		nonce, err = account.Nonce(ctx)
		if err != nil {
			return nil, err
		}
	}
	maxFee, _ := big.NewInt(0).SetString("0x200000000", 0)
	if details.MaxFee != nil {
		maxFee = details.MaxFee
	}
	version := big.NewInt(0)
	if account.version != nil {
		version = account.version
	}
	txHash, err := account.TransactionHash(
		calls,
		types.ExecuteDetails{
			Nonce:  nonce,
			MaxFee: maxFee,
		},
	)
	if err != nil {
		return nil, err
	}
	s1, s2, err := account.Sign(txHash)
	if err != nil {
		return nil, err
	}
	calldata := fmtExecuteCalldataStrings(nonce, calls)
	accountDefaultV0Entrypoint := "__execute__"
	call := types.Call{
		MaxFee:             fmt.Sprintf("0x%s", maxFee.Text(16)),
		Version:            types.NumAsHex(fmt.Sprintf("0x%s", version.Text(16))),
		Signature:          []string{fmt.Sprintf("0x%s", s1.Text(16)), fmt.Sprintf("0x%s", s2.Text(16))},
		ContractAddress:    types.HexToHash(account.Address),
		EntryPointSelector: &accountDefaultV0Entrypoint,
		CallData:           calldata,
	}
	return account.Provider.EstimateFee(ctx, call, WithBlockTag("latest"))
}

func (account *Account) Execute(ctx context.Context, calls []types.FunctionCall, details types.ExecuteDetails) (*types.AddInvokeTransactionOutput, error) {
	if account.version != nil && account.version.Cmp(big.NewInt(0)) != 0 {
		return nil, errors.New("only invoke v0 is implemented")
	}
	var err error
	version := big.NewInt(0)
	if account.version != nil {
		version = account.version
	}
	nonce := details.Nonce
	if details.Nonce == nil {
		nonce, err = account.Nonce(ctx)
		if err != nil {
			return nil, err
		}
	}
	maxFee := details.MaxFee
	if details.MaxFee == nil {
		estimate, err := account.EstimateFee(ctx, calls, details)
		if err != nil {
			return nil, err
		}
		v, ok := big.NewInt(0).SetString(string(estimate.OverallFee), 0)
		if !ok {
			return nil, errors.New("could not match OverallFee to big.Int")
		}
		maxFee = v.Mul(v, big.NewInt(2))
	}
	txHash, err := account.TransactionHash(
		calls,
		types.ExecuteDetails{
			Nonce:  nonce,
			MaxFee: maxFee,
		},
	)
	if err != nil {
		return nil, err
	}
	s1, s2, err := account.Sign(txHash)
	if err != nil {
		return nil, err
	}
	calldata := fmtExecuteCalldataStrings(nonce, calls)
	return account.Provider.AddInvokeTransaction(
		context.Background(),
		types.FunctionCall{
			ContractAddress:    types.HexToHash(account.Address),
			EntryPointSelector: "__execute__",
			CallData:           calldata,
		},
		[]string{fmt.Sprintf("0x%s", s1.Text(16)), fmt.Sprintf("0x%s", s2.Text(16))},
		fmt.Sprintf("0x%s", maxFee.Text(16)),
		fmt.Sprintf("0x%s", version.Text(16)),
	)
}
