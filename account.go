package caigo

import (
	"context"
	"fmt"
	"math/big"

	"github.com/dontpanicdao/caigo/types"
)

const (
	EXECUTE_SELECTOR    string = "__execute__"
	TRANSACTION_PREFIX  string = "invoke"
	TRANSACTION_VERSION int64  = 0
)

type Account struct {
	Curve    StarkCurve
	Provider types.Provider
	Address  string
	private  *big.Int
	PublicX  *big.Int
	PublicY  *big.Int
}

/*
	Instantiate a new StarkNet Account which includes structures for calling the network and signing transactions:
	- private signing key
	- stark curve definition
	- full provider definition
	- public key pair for signature verifications
*/
func (sc StarkCurve) NewAccount(private, address string, provider types.Provider) (*Account, error) {
	if len(sc.ConstantPoints) == 0 {
		return nil, fmt.Errorf("must initiate precomputed constant points")
	}
	priv := SNValToBN(private)
	x, y, err := sc.PrivateToPoint(priv)
	if err != nil {
		return nil, err
	}

	return &Account{
		Curve:    sc,
		Provider: provider,
		Address:  address,
		private:  priv,
		PublicX:  x,
		PublicY:  y,
	}, nil
}

/*
	invocation wrapper for StarkNet account calls to '__execute__' contact calls through an account abstraction
	- implementation has been tested against OpenZeppelin Account contract as of: https://github.com/OpenZeppelin/cairo-contracts/blob/4116c1ecbed9f821a2aa714c993a35c1682c946e/src/openzeppelin/account/Account.cairo
	- accepts a multicall
*/
func (account *Account) Execute(ctx context.Context, call types.Transaction) (*types.AddTxResponse, error) {
	calls := []types.Transaction{call}

	nonce, err := account.Provider.AccountNonce(ctx, account.Address)
	if err != nil {
		return nil, err
	}

	req := types.Transaction{
		ContractAddress:    account.Address,
		EntryPointSelector: EXECUTE_SELECTOR,
		Calldata:           FmtExecuteCalldataStrings(nonce, calls),
	}

	// provide good signature so we can get estimate for ECDSA signing
	fee, err := account.EstimateFeeMultiCall(ctx, req, nonce, calls)
	if err != nil {
		return nil, err
	}
	req.MaxFee = fee

	r, s, err := account.SignMultiCall(req.MaxFee, nonce, calls)
	if err != nil {
		return nil, err
	}
	req.Signature = []string{r.String(), s.String()}

	return account.Provider.Invoke(ctx, req)
}

func (account *Account) ExecuteMultiCall(ctx context.Context, maxFee string, calls []types.Transaction) (*types.AddTxResponse, error) {
	nonce, err := account.Provider.AccountNonce(ctx, account.Address)
	if err != nil {
		return nil, err
	}

	req := types.Transaction{
		ContractAddress:    account.Address,
		EntryPointSelector: EXECUTE_SELECTOR,
		MaxFee:             maxFee,
		Calldata:           FmtExecuteCalldataStrings(nonce, calls),
	}

	r, s, err := account.SignMultiCall(maxFee, nonce, calls)
	if err != nil {
		return nil, err
	}

	req.Signature = []string{r.String(), s.String()}

	return account.Provider.Invoke(ctx, req)
}

func (account *Account) SignMultiCall(maxFee string, nonce *big.Int, calls []types.Transaction) (*big.Int, *big.Int, error) {
	hash, err := account.HashMultiCall(maxFee, nonce, calls)
	if err != nil {
		return nil, nil, err
	}

	return account.Curve.Sign(hash, account.private)
}

func (account *Account) HashMultiCall(fee string, nonce *big.Int, calls []types.Transaction) (hash *big.Int, err error) {
	chainID, err := account.Provider.ChainID(context.Background())
	if err != nil {
		return nil, err
	}

	callArray := FmtExecuteCalldata(nonce, calls)
	callArray = append(callArray, big.NewInt(int64(len(callArray))))
	cdHash, err := sc.HashElements(callArray)
	if err != nil {
		return hash, err
	}

	multiHashData := []*big.Int{
		UTF8StrToBig(TRANSACTION_PREFIX),
		big.NewInt(TRANSACTION_VERSION),
		SNValToBN(account.Address),
		GetSelectorFromName(EXECUTE_SELECTOR),
		cdHash,
		SNValToBN(fee),
		UTF8StrToBig(chainID),
	}

	multiHashData = append(multiHashData, big.NewInt(int64(len(multiHashData))))
	hash, err = account.Curve.HashElements(multiHashData)
	return hash, err
}

func (account *Account) EstimateFeeMultiCall(ctx context.Context, req types.Transaction, nonce *big.Int, calls []types.Transaction) (string, error) {
	r, s, err := account.SignMultiCall("0", nonce, calls)
	if err != nil {
		return "", err
	}

	req.Signature = []string{r.String(), s.String()}
	fee, err := account.Provider.EstimateFee(ctx, req)

	return BigToHex(fee.Amount), err
}

func FmtExecuteCalldataStrings(nonce *big.Int, calls []types.Transaction) (calldataStrings []string) {
	callArray := FmtExecuteCalldata(nonce, calls)
	for _, data := range callArray {
		calldataStrings = append(calldataStrings, data.String())
	}
	return calldataStrings
}

/*
	Formats the multicall transactions in a format which can be signed and verified by the network and OpenZeppelin account contracts
*/
func FmtExecuteCalldata(nonce *big.Int, calls []types.Transaction) (calldataArray []*big.Int) {
	callArray := []*big.Int{big.NewInt(int64(len(calls)))}

	for _, tx := range calls {
		callArray = append(callArray, SNValToBN(tx.ContractAddress), GetSelectorFromName(tx.EntryPointSelector))

		if len(tx.Calldata) == 0 {
			callArray = append(callArray, big.NewInt(0), big.NewInt(0))

			continue
		}

		callArray = append(callArray, big.NewInt(int64(len(calldataArray))), big.NewInt(int64(len(tx.Calldata))))
		for _, cd := range tx.Calldata {
			calldataArray = append(calldataArray, SNValToBN(cd))
		}
	}

	callArray = append(callArray, big.NewInt(int64(len(calldataArray))))
	callArray = append(callArray, calldataArray...)
	callArray = append(callArray, nonce)
	return callArray
}
