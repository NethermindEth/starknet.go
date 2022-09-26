package rpc

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/dontpanicdao/caigo"
	"github.com/dontpanicdao/caigo/rpc/types"
)

type Policy struct {
	ContractAddress string `json:"contractAddress"`
	Selector        string `json:"selector"`
}

type XSession struct {
	Key      string   `json:"key"`
	Expires  uint64   `json:"expires"`
	Policies []Policy `json:"policies"`
}

var _ account = &XSessionAccount{}

// H(Policy(contractAddress:felt,selector:selector))
var POLICY_TYPE_HASH, _ = big.NewInt(0).SetString("0x2f0026e78543f036f33e26a8f5891b88c58dc1e20cbbfaf0bb53274da6fa568", 0)

type XSessionAccount struct {
	Provider *Provider
	Address  string
	plugin   *big.Int
	private  *big.Int
	version  *big.Int
	xsession XSession
	mt       caigo.FixedSizeMerkleTree
}

func (provider *Provider) NewXSessionAccount(private, address string, pluginClassHash string, xsession XSession, options ...AccountOption) (*XSessionAccount, error) {
	version := big.NewInt(0)
	for _, o := range options {
		opt := o()
		if opt.version != nil {
			version = opt.version
		}
	}
	if version.Cmp(big.NewInt(0)) != 0 {
		return nil, errors.New("account v1 not yet supported")
	}
	priv, ok := big.NewInt(0).SetString(private, 0)
	if !ok {
		return nil, errors.New("could not convert private key")
	}
	plugin, ok := big.NewInt(0).SetString(pluginClassHash, 0)
	if !ok {
		return nil, errors.New("could not convert plugin class hash")
	}
	leaves := []*big.Int{}
	for _, policy := range xsession.Policies {
		contract, ok := big.NewInt(0).SetString(policy.ContractAddress, 0)
		if !ok {
			return nil, errors.New("could not convert contract address")
		}
		leaf, _ := caigo.Curve.HashElements(append(
			[]*big.Int{},
			POLICY_TYPE_HASH,
			contract,
			caigo.GetSelectorFromName(policy.Selector),
		))
		leaves = append(leaves, leaf)
	}
	mt, err := caigo.NewFixedSizeMerkleTree(leaves...)
	if err != nil {
		return nil, fmt.Errorf("could not create merkle tree, error: %v", err)
	}
	return &XSessionAccount{
		Provider: provider,
		Address:  address,
		plugin:   plugin,
		private:  priv,
		version:  version,
		xsession: xsession,
		mt:       *mt,
	}, nil
}

func (account *XSessionAccount) pluginCall(calls []types.FunctionCall) (types.FunctionCall, error) {
	data := []string{
		fmt.Sprintf("0x%s", account.plugin.Text(16)),
		fmt.Sprintf("0x%s", big.NewInt(int64(account.xsession.Expires))),
		fmt.Sprintf("0x%s", account.mt.Root.Text(16)),
	}
	proofs := []*big.Int{}
	proofSize := int64(0)
	for _, call := range calls {
		leaf, _ := caigo.Curve.HashElements(append(
			[]*big.Int{},
			POLICY_TYPE_HASH,
			call.ContractAddress.Big(),
			caigo.GetSelectorFromName(call.EntryPointSelector),
		))
		p, err := account.mt.GetProof(leaf, 0, []*big.Int{})
		if err != nil {
			return types.FunctionCall{}, err
		}
		if proofSize == 0 {
			proofSize = int64(len(p))
		}
		if proofSize != int64(len(p)) {
			return types.FunctionCall{}, errors.New("proof does not match proofsize")
		}
		proofs = append(proofs, p...)
	}
	data = append(data, fmt.Sprintf("0x%s", big.NewInt(proofSize).Text(16)))
	for _, proof := range proofs {
		data = append(data, fmt.Sprintf("0x%s", proof.Text(16)))
	}
	return types.FunctionCall{
		ContractAddress:    types.HexToHash(account.Address),
		EntryPointSelector: "use_plugin",
		CallData:           data,
	}, nil
}

func (account *XSessionAccount) Call(ctx context.Context, call types.FunctionCall) ([]string, error) {
	return account.Provider.Call(ctx, call, WithBlockTag("latest"))
}

func (account *XSessionAccount) Sign(msgHash *big.Int) (*big.Int, *big.Int, error) {
	return caigo.Curve.Sign(msgHash, account.private)
}

func (account *XSessionAccount) TransactionHash(calls []types.FunctionCall, details types.ExecuteDetails) (*big.Int, error) {
	chainID, err := account.Provider.ChainID(context.Background())
	if err != nil {
		return nil, err
	}
	var callArray []*big.Int
	switch {
	case account.version.Cmp(big.NewInt(0)) == 0:
		callArray = fmtV0Calldata(details.Nonce, calls)
	case account.version.Cmp(big.NewInt(1)) == 0:
		callArray = fmtCalldata(calls)
	default:
		return nil, fmt.Errorf("version %s unsupported", account.version.Text(10))
	}
	cdHash, err := caigo.Curve.ComputeHashOnElements(callArray)
	if err != nil {
		return nil, err
	}

	var multiHashData []*big.Int
	switch {
	case account.version.Cmp(big.NewInt(0)) == 0:
		multiHashData = []*big.Int{
			caigo.UTF8StrToBig(TRANSACTION_PREFIX),
			account.version,
			caigo.SNValToBN(account.Address),
			caigo.GetSelectorFromName(EXECUTE_SELECTOR),
			cdHash,
			details.MaxFee,
			caigo.UTF8StrToBig(chainID),
		}
	case account.version.Cmp(big.NewInt(1)) == 0:
		multiHashData = []*big.Int{
			caigo.UTF8StrToBig(TRANSACTION_PREFIX),
			account.version,
			caigo.SNValToBN(account.Address),
			cdHash,
			details.MaxFee,
			details.Nonce,
			caigo.UTF8StrToBig(chainID),
		}
	default:
		return nil, fmt.Errorf("version %s unsupported", account.version.Text(10))
	}
	return caigo.Curve.ComputeHashOnElements(multiHashData)
}

func (account *XSessionAccount) Nonce(ctx context.Context) (*big.Int, error) {
	switch {
	case account.version.Cmp(big.NewInt(0)) == 0:
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
	case account.version.Cmp(big.NewInt(1)) == 0:
		nonce, err := account.Provider.Nonce(
			ctx,
			types.HexToHash(account.Address),
		)
		if err != nil {
			return nil, err
		}
		if nonce == nil {
			return nil, errors.New("nonce is nil")
		}
		n, ok := big.NewInt(0).SetString(*nonce, 0)
		if !ok {
			return nil, errors.New("nonce error")
		}
		return n, nil
	}
	return nil, fmt.Errorf("version %s unsupported", account.version.Text(10))
}

func (account *XSessionAccount) EstimateFee(ctx context.Context, calls []types.FunctionCall, details types.ExecuteDetails) (*types.FeeEstimate, error) {
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
	pluginCall, err := account.pluginCall(calls)
	if err != nil {
		return nil, err
	}
	txHash, err := account.TransactionHash(
		append([]types.FunctionCall{pluginCall}, calls...),
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
	var calldata []string
	switch {
	case account.version.Cmp(big.NewInt(0)) == 0:
		calldata = fmtV0CalldataStrings(nonce, append([]types.FunctionCall{pluginCall}, calls...))
	case account.version.Cmp(big.NewInt(1)) == 0:
		calldata = fmtCalldataStrings(append([]types.FunctionCall{pluginCall}, calls...))
	default:
		return nil, fmt.Errorf("version %s unsupported", account.version.Text(10))
	}
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

func (account *XSessionAccount) Execute(ctx context.Context, calls []types.FunctionCall, details types.ExecuteDetails) (*types.AddInvokeTransactionOutput, error) {
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
	pluginCall, err := account.pluginCall(calls)
	if err != nil {
		return nil, err
	}
	txHash, err := account.TransactionHash(
		append([]types.FunctionCall{pluginCall}, calls...),
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
	var calldata []string
	switch {
	case account.version.Cmp(big.NewInt(0)) == 0:
		calldata = fmtV0CalldataStrings(nonce, append([]types.FunctionCall{pluginCall}, calls...))
	case account.version.Cmp(big.NewInt(1)) == 0:
		calldata = fmtCalldataStrings(append([]types.FunctionCall{pluginCall}, calls...))
	default:
		return nil, fmt.Errorf("version %s unsupported", account.version.Text(10))
	}
	// TODO: change this payload to manage both V0 and V1
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
