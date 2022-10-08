package caigo

import (
	"context"
	"math/big"

	"github.com/dontpanicdao/caigo/gateway"
	"github.com/dontpanicdao/caigo/types"
)

const (
	EXECUTE_SELECTOR    string = "__execute__"
	TRANSACTION_PREFIX  string = "invoke"
	TRANSACTION_VERSION int64  = 0
	FEE_MARGIN          uint64 = 115
)

type Account struct {
	Provider types.Provider
	Address  types.Hash
	PublicX  *big.Int
	PublicY  *big.Int
	private  *big.Int
}

/*
Instantiate a new StarkNet Account which includes structures for calling the network and signing transactions:
- private signing key
- stark curve definition
- full provider definition
- public key pair for signature verifications
*/
func NewGatewayAccount(private string, address types.Hash, provider *gateway.Gateway) (*Account, error) {
	priv := types.SNValToBN(private)
	x, y, err := Curve.PrivateToPoint(priv)
	if err != nil {
		return nil, err
	}

	return &Account{
		Provider: provider,
		Address:  address,
		PublicX:  x,
		PublicY:  y,
		private:  priv,
	}, nil
}

var _ account = &Account{}

func (account *Account) Sign(msgHash *big.Int) (*big.Int, *big.Int, error) {
	return Curve.Sign(msgHash, account.private)
}

func (account *Account) Nonce(ctx context.Context) (*big.Int, error) {
	return account.Provider.AccountNonce(ctx, account.Address)
}

func (account *Account) Call(ctx context.Context, call types.FunctionCall) ([]string, error) {
	return account.Provider.Call(ctx, call, "latest")
}

/*
invocation wrapper for StarkNet account calls to '__execute__' contact calls through an account abstraction
- implementation has been tested against OpenZeppelin Account contract as of: https://github.com/OpenZeppelin/cairo-contracts/blob/4116c1ecbed9f821a2aa714c993a35c1682c946e/src/openzeppelin/account/Account.cairo
- accepts a multicall
*/
func (account *Account) Execute(ctx context.Context, calls []types.FunctionCall, details types.ExecuteDetails) (*types.AddInvokeTransactionOutput, error) {
	if details.Nonce == nil {
		nonce, err := account.Provider.AccountNonce(ctx, account.Address)
		if err != nil {
			return nil, err
		}
		details.Nonce = nonce
	}

	if details.MaxFee == nil {
		feeEstimate, err := account.EstimateFee(ctx, calls, details)
		if err != nil {
			return nil, err
		}
		fee, _ := big.NewInt(0).SetString(string(feeEstimate.OverallFee), 0)
		expandedFee := big.NewInt(0).Mul(fee, big.NewInt(int64(FEE_MARGIN)))
		max := big.NewInt(0).Div(expandedFee, big.NewInt(100))
		details.MaxFee = max
	}

	req, err := account.fmtExecute(ctx, calls, details)
	if err != nil {
		return nil, err
	}

	return account.Provider.Invoke(ctx, *req)
}

func (account *Account) TransactionHash(calls []types.FunctionCall, details types.ExecuteDetails) (*big.Int, error) {
	chainID, err := account.Provider.ChainID(context.Background())
	if err != nil {
		return nil, err
	}

	callArray := fmtV0Calldata(details.Nonce, calls)
	cdHash, err := Curve.ComputeHashOnElements(callArray)
	if err != nil {
		return nil, err
	}

	multiHashData := []*big.Int{
		types.UTF8StrToBig(TRANSACTION_PREFIX),
		big.NewInt(TRANSACTION_VERSION),
		account.Address.Big(),
		types.GetSelectorFromName(EXECUTE_SELECTOR),
		cdHash,
		details.MaxFee,
		types.UTF8StrToBig(chainID),
	}

	return Curve.ComputeHashOnElements(multiHashData)
}

func (account *Account) EstimateFee(ctx context.Context, calls []types.FunctionCall, details types.ExecuteDetails) (*types.FeeEstimate, error) {
	if details.Nonce == nil {
		nonce, err := account.Provider.AccountNonce(ctx, account.Address)
		if err != nil {
			return nil, err
		}
		details.Nonce = nonce
	}

	if details.MaxFee == nil {
		details.MaxFee = big.NewInt(0)
	}

	req, err := account.fmtExecute(ctx, calls, details)
	if err != nil {
		return nil, err
	}

	return account.Provider.EstimateFee(ctx, *req, "")
}

func (account *Account) fmtExecute(ctx context.Context, calls []types.FunctionCall, details types.ExecuteDetails) (*types.FunctionInvoke, error) {
	req := types.FunctionInvoke{
		FunctionCall: types.FunctionCall{
			ContractAddress:    account.Address,
			EntryPointSelector: EXECUTE_SELECTOR,
			Calldata:           fmtV0CalldataStrings(details.Nonce, calls),
		},
		MaxFee: details.MaxFee,
	}

	hash, err := account.TransactionHash(calls, types.ExecuteDetails{
		MaxFee: details.MaxFee, Nonce: details.Nonce,
	})
	if err != nil {
		return nil, err
	}

	r, s, err := account.Sign(hash)
	if err != nil {
		return nil, err
	}
	req.Signature = types.Signature{r, s}

	return &req, nil
}
