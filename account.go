package caigo

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
)

type Signer struct {
	private *big.Int
	Curve   StarkCurve
	Gateway *StarknetGateway
	PublicX *big.Int
	PublicY *big.Int
}

type FeeEstimate struct {
	Amount *big.Int `json:"amount"`
	Unit   string   `json:"unit"`
}

/*
	Instantiate a new StarkNet Signer which includes structures for calling the network and signing transactions:
	- private signing key
	- stark curve definition
	- full StarknetGateway definition
	- public key pair for signature verifications
*/
func (sc StarkCurve) NewSigner(private, pubX, pubY *big.Int, gw *StarknetGateway) (*Signer, error) {
	if len(sc.ConstantPoints) == 0 {
		return nil, fmt.Errorf("must initiate precomputed constant points")
	}

	return &Signer{
		private: private,
		Curve:   sc,
		Gateway: gw,
		PublicX: pubX,
		PublicY: pubY,
	}, nil
}

/*
	invocation wrapper for StarkNet account calls to '__execute__' contact calls through an account abstraction
	- implementation has been tested against OpenZeppelin Account contract as of: https://github.com/OpenZeppelin/cairo-contracts/blob/4116c1ecbed9f821a2aa714c993a35c1682c946e/src/openzeppelin/account/Account.cairo
	- accepts a multicall
*/
func (signer *Signer) Execute(ctx context.Context, address string, txs []Transaction) (addResp AddTxResponse, err error) {
	nonce, err := signer.Gateway.AccountNonce(ctx, address)
	if err != nil {
		return addResp, err
	}

	maxFee := big.NewInt(0)

	hash, err := signer.Curve.HashMulticall(address, nonce, maxFee, UTF8StrToBig(signer.Gateway.ChainId), txs)
	if err != nil {
		return addResp, err
	}

	r, s, err := signer.Curve.Sign(hash, signer.private)
	if err != nil {
		return addResp, err
	}

	req := Transaction{
		ContractAddress:    address,
		EntryPointSelector: BigToHex(GetSelectorFromName(EXECUTE_SELECTOR)),
		Calldata:           FmtExecuteCalldataStrings(nonce, txs),
		Signature:          []string{r.String(), s.String()},
	}

	return signer.Gateway.Invoke(ctx, req)
}

func (sg *StarknetGateway) AccountNonce(ctx context.Context, address string) (nonce *big.Int, err error) {
	resp, err := sg.Call(ctx, Transaction{
		ContractAddress:    address,
		EntryPointSelector: "get_nonce",
	}, nil)
	if err != nil {
		return nonce, err
	}
	if len(resp) == 0 {
		return nonce, fmt.Errorf("no resp in contract call 'get_nonce' %v", address)
	}

	return HexToBN(resp[0]), nil
}

func (sg *StarknetGateway) EstimateFee(ctx context.Context, tx Transaction) (fee FeeEstimate, err error) {
	req, err := sg.newRequest(ctx, http.MethodPost, "/estimate_fee", tx)
	if err != nil {
		return fee, err
	}

	return fee, sg.do(req, &fee)
}

func (sc StarkCurve) HashMulticall(addr string, nonce, maxFee, chainId *big.Int, txs []Transaction) (hash *big.Int, err error) {
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
		maxFee,
		chainId,
	}

	multiHashData = append(multiHashData, big.NewInt(int64(len(multiHashData))))
	hash, err = sc.HashElements(multiHashData)
	return hash, err
}

func FmtExecuteCalldataStrings(nonce *big.Int, txs []Transaction) (calldataStrings []string) {
	callArray := FmtExecuteCalldata(nonce, txs)
	for _, data := range callArray {
		calldataStrings = append(calldataStrings, data.String())
	}
	return calldataStrings
}

/*
	Formats the multicall transactions in a format which can be signed and verified by the network and OpenZeppelin account contracts
*/
func FmtExecuteCalldata(nonce *big.Int, txs []Transaction) (calldataArray []*big.Int) {
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
