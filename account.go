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

type Signer struct {
	Curve    StarkCurve
	Provider types.Provider
	Address  string
	private  *big.Int
	PublicX  *big.Int
	PublicY  *big.Int
}

/*
	Instantiate a new StarkNet Signer which includes structures for calling the network and signing transactions:
	- private signing key
	- stark curve definition
	- full provider definition
	- public key pair for signature verifications
*/
func (sc StarkCurve) NewSigner(private *big.Int, address string, provider types.Provider) (*Signer, error) {
	if len(sc.ConstantPoints) == 0 {
		return nil, fmt.Errorf("must initiate precomputed constant points")
	}
	x, y, err := sc.PrivateToPoint(private)
	if err != nil {
		return nil, err
	}

	return &Signer{
		Curve:    sc,
		Provider: provider,
		Address:  address,
		private:  private,
		PublicX:  x,
		PublicY:  y,
 	}, nil
}


/*
	invocation wrapper for StarkNet account calls to '__execute__' contact calls through an account abstraction
	- implementation has been tested against OpenZeppelin Account contract as of: https://github.com/OpenZeppelin/cairo-contracts/blob/4116c1ecbed9f821a2aa714c993a35c1682c946e/src/openzeppelin/account/Account.cairo
	- accepts a multicall
*/
func (signer *Signer) Execute(ctx context.Context, maxFee string, txs []types.Transaction) (*types.AddTxResponse, error) {
	chainID, err := signer.Provider.ChainID(ctx)
	if err != nil {
		return nil, err
	}

	nonce, err := signer.Provider.AccountNonce(ctx, signer.Address)
	if err != nil {
		return nil, err
	}

	req := types.Transaction{
		ContractAddress:    signer.Address,
		EntryPointSelector: EXECUTE_SELECTOR,
		MaxFee:				maxFee,
		Calldata:           FmtExecuteCalldataStrings(nonce, txs),
	}

	r, s, err := signer.SignMulticall(signer.Address, chainID, maxFee, nonce, txs)
	if err != nil {
		return nil, err
	}

	req.Signature = []string{r.String(), s.String()}

	return signer.Provider.Invoke(ctx, req)
}

func (signer *Signer) ExecuteSingle(ctx context.Context, tx types.Transaction) (*types.AddTxResponse, error) {
	txs := []types.Transaction{tx}

	chainID, err := signer.Provider.ChainID(ctx)
	if err != nil {
		return nil, err
	}

	nonce, err := signer.Provider.AccountNonce(ctx, signer.Address)
	if err != nil {
		return nil, err
	}

	req := types.Transaction{
		ContractAddress:    signer.Address,
		EntryPointSelector: EXECUTE_SELECTOR,
		Calldata:           FmtExecuteCalldataStrings(nonce, txs),
	}

	feeR, feeS, err := signer.SignMulticall(signer.Address, chainID, "0", nonce, txs)
	req.Signature = []string{feeR.String(), feeS.String()}

	fee, err := signer.Provider.EstimateFee(ctx, req)
	req.MaxFee = BigToHex(fee.Amount)
	r, s, err := signer.SignMulticall(signer.Address, chainID, req.MaxFee, nonce, txs)
	if err != nil {
		return nil, err
	}

	req.Signature = []string{r.String(), s.String()}

	return signer.Provider.Invoke(ctx, req)
}

func (signer *Signer) SignMulticall(address, chainID, maxFee string, nonce *big.Int, txs []types.Transaction) (*big.Int, *big.Int, error) {
	hash, err := signer.Curve.HashMulticall(address, chainID, maxFee, nonce, txs)
	if err != nil {
		return nil, nil, err
	}

	return signer.Curve.Sign(hash, signer.private)
}

func (sc StarkCurve) HashMulticall(addr, chainId, fee string, nonce *big.Int, txs []types.Transaction) (hash *big.Int, err error) {
	callArray := FmtExecuteCalldata(nonce, txs)
	callArray = append(callArray, big.NewInt(int64(len(callArray))))
	cdHash, err := sc.HashElements(callArray)
	if err != nil {
		return hash, err
	}

	multiHashData := []*big.Int{
		UTF8StrToBig(TRANSACTION_PREFIX),
		big.NewInt(TRANSACTION_VERSION),
		SNValToBN(addr),
		GetSelectorFromName(EXECUTE_SELECTOR),
		cdHash,
		SNValToBN(fee),
		UTF8StrToBig(chainId),
	}

	multiHashData = append(multiHashData, big.NewInt(int64(len(multiHashData))))
	hash, err = sc.HashElements(multiHashData)
	return hash, err
}

func FmtExecuteCalldataStrings(nonce *big.Int, txs []types.Transaction) (calldataStrings []string) {
	callArray := FmtExecuteCalldata(nonce, txs)
	for _, data := range callArray {
		calldataStrings = append(calldataStrings, data.String())
	}
	return calldataStrings
}

/*
	Formats the multicall transactions in a format which can be signed and verified by the network and OpenZeppelin account contracts
*/
func FmtExecuteCalldata(nonce *big.Int, txs []types.Transaction) (calldataArray []*big.Int) {
	callArray := []*big.Int{big.NewInt(int64(len(txs)))}

	for _, tx := range txs {
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
