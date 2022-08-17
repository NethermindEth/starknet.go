package caigo

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/dontpanicdao/caigo/felt"
	"github.com/dontpanicdao/caigo/types"
)

var (
	EXECUTE_SELECTOR      = felt.GetSelectorFromName("__execute__")
	TRANSACTION_PREFIX, _ = felt.UTF8StrToFelt("invoke")
	FEE_MARGIN            = felt.BigToFelt(big.NewInt(115))
	TRANSACTION_VERSION   = felt.BigToFelt(big.NewInt(0))
)

var errNonceEmpty = errors.New("emptynonce")

type Account struct {
	Provider types.Provider
	Address  felt.Felt
	PublicX  *big.Int
	PublicY  *big.Int
	private  *big.Int
}

type ExecuteDetails struct {
	MaxFee  *felt.Felt
	Nonce   *felt.Felt
	Version *felt.Felt
}

/*
Instantiate a new StarkNet Account which includes structures for calling the network and signing transactions:
- private signing key
- stark curve definition
- full provider definition
- public key pair for signature verifications
*/
func NewAccount(private *big.Int, address felt.Felt, provider types.Provider) (*Account, error) {
	if private == nil {
		return nil, fmt.Errorf("wrongPrivate")
	}
	x, y, err := felt.GetCurve().PrivateToPoint(private)
	if err != nil {
		return nil, err
	}
	return &Account{
		Provider: provider,
		Address:  address,
		PublicX:  x,
		PublicY:  y,
		private:  private,
	}, nil
}

func (account *Account) Sign(msgHash felt.Felt) (*felt.Signature, error) {
	return felt.GetCurve().Sign(msgHash, account.private)
}

/*
invocation wrapper for StarkNet account calls to '__execute__' contact calls through an account abstraction
- implementation has been tested against OpenZeppelin Account contract as of: https://github.com/OpenZeppelin/cairo-contracts/blob/4116c1ecbed9f821a2aa714c993a35c1682c946e/src/openzeppelin/account/Account.cairo
- accepts a multicall
*/
func (account *Account) Execute(ctx context.Context, calls []types.Transaction, details ExecuteDetails) (*types.AddTxResponse, error) {
	if details.Nonce == nil {
		nonce, err := account.Provider.AccountNonce(ctx, account.Address)
		if err != nil {
			return nil, err
		}
		details.Nonce = nonce
	}

	if details.MaxFee == nil {
		fee, err := account.EstimateFee(ctx, calls, details)
		if err != nil {
			return nil, err
		}
		if fee.OverallFee == nil {
			return nil, errors.New("overallFeeUnknown")
		}
		margin := felt.NewFelt().Mul(*fee.OverallFee, FEE_MARGIN)
		if margin.IsNil() {
			return nil, errors.New("marginIsNil")
		}
		details.MaxFee = margin.Div(*margin, felt.BigToFelt(big.NewInt(100)))
	}

	req, err := account.fmtExecute(ctx, calls, details)
	if err != nil {
		return nil, err
	}

	return account.Provider.Invoke(ctx, *req)
}

func (account *Account) HashMultiCall(fee *felt.Felt, nonce felt.Felt, calls []types.Transaction) (*felt.Felt, error) {
	chainID, err := account.Provider.ChainID(context.Background())
	if err != nil {
		return nil, err
	}
	callArray := ExecuteCalldata(nonce, calls)

	// convert callArray into a BigInt array
	callArrayFelt := make([]felt.Felt, 0)
	callArrayFelt = append(callArrayFelt, callArray...)

	cdHash, err := felt.GetCurve().ComputeHashOnElements(callArrayFelt)
	if err != nil {
		return nil, err
	}
	chain, _ := felt.UTF8StrToFelt(chainID)
	multiHashData := []felt.Felt{
		*TRANSACTION_PREFIX,
		TRANSACTION_VERSION,
		account.Address,
		EXECUTE_SELECTOR,
		*cdHash,
		felt.BigToFelt(fee.Int),
		*chain,
	}

	return felt.GetCurve().ComputeHashOnElements(multiHashData)
}

func (account *Account) EstimateFee(ctx context.Context, calls []types.Transaction, details ExecuteDetails) (*types.FeeEstimate, error) {
	if details.Nonce == nil {
		nonce, err := account.Provider.AccountNonce(ctx, account.Address)
		if err != nil {
			return nil, err
		}
		details.Nonce = nonce
	}

	if details.MaxFee == nil {
		maxFee := felt.BigToFelt(big.NewInt(0))
		details.MaxFee = &maxFee
	}

	req, err := account.fmtExecute(ctx, calls, details)
	if err != nil {
		return nil, err
	}

	return account.Provider.EstimateFee(ctx, *req, "")
}

func (account *Account) fmtExecute(ctx context.Context, calls []types.Transaction, details ExecuteDetails) (*types.FunctionInvoke, error) {
	if details.Nonce == nil {
		return nil, errNonceEmpty
	}
	req := types.FunctionInvoke{
		FunctionCall: types.FunctionCall{
			ContractAddress:    account.Address,
			EntryPointSelector: &EXECUTE_SELECTOR,
			Calldata:           ExecuteCalldata(*details.Nonce, calls),
		},
		MaxFee: details.MaxFee,
	}

	hash, err := account.HashMultiCall(details.MaxFee, *details.Nonce, calls)
	if err != nil {
		return nil, err
	}

	signature, err := account.Sign(*hash)
	if err != nil {
		return nil, err
	}
	req.Signature = *signature
	return &req, nil
}

/*
Formats the multicall transactions in a format which can be signed and verified by the network and OpenZeppelin account contracts
*/
func ExecuteCalldata(nonce felt.Felt, calls []types.Transaction) (calldataArray []felt.Felt) {
	callArray := []felt.Felt{felt.BigToFelt(big.NewInt(int64(len(calls))))}

	for _, tx := range calls {
		callArray = append(callArray, tx.ContractAddress)
		if tx.EntryPointSelector != nil {
			callArray = append(callArray, *tx.EntryPointSelector)
		}

		if len(tx.Calldata) == 0 {
			callArray = append(callArray, felt.BigToFelt(big.NewInt(0)), felt.BigToFelt(big.NewInt(0)))
			continue
		}

		callArray = append(callArray, felt.BigToFelt(big.NewInt(int64(len(calldataArray)))), felt.BigToFelt(big.NewInt(int64(len(tx.Calldata)))))
		calldataArray = append(calldataArray, tx.Calldata...)
	}

	callArray = append(callArray, felt.BigToFelt(big.NewInt(int64(len(calldataArray)))))
	callArray = append(callArray, calldataArray...)
	callArray = append(callArray, nonce)
	return callArray
}
