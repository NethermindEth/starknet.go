package caigo

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/dontpanicdao/caigo/gateway"
	"github.com/dontpanicdao/caigo/rpcv01"
	"github.com/dontpanicdao/caigo/rpcv02"
	"github.com/dontpanicdao/caigo/types"
)

var (
	ErrUnsupportedAccount = errors.New("unsupported account implementation")
	MAX_FEE, _            = big.NewInt(0).SetString("0x20000000000000", 0)
)

const (
	TRANSACTION_PREFIX      = "invoke"
	DECLARE_PREFIX          = "declare"
	EXECUTE_SELECTOR        = "__execute__"
	CONTRACT_ADDRESS_PREFIX = "STARKNET_CONTRACT_ADDRESS"
)

type account interface {
	TransactionHash(calls []types.FunctionCall, details types.ExecuteDetails) (*big.Int, error)
	Call(ctx context.Context, call types.FunctionCall) ([]string, error)
	Nonce(ctx context.Context) (*big.Int, error)
	EstimateFee(ctx context.Context, calls []types.FunctionCall, details types.ExecuteDetails) (*types.FeeEstimate, error)
	Execute(ctx context.Context, calls []types.FunctionCall, details types.ExecuteDetails) (*types.AddInvokeTransactionOutput, error)
	Declare(ctx context.Context, classHash string, contract types.ContractClass, details types.ExecuteDetails) (types.AddDeclareResponse, error)
	Deploy(ctx context.Context, classHash string, details types.ExecuteDetails) (*types.AddDeployResponse, error)
}

var _ account = &Account{}

type AccountPlugin interface {
	PluginCall(calls []types.FunctionCall) (types.FunctionCall, error)
}

type ProviderType string

const (
	ProviderRPCv01  ProviderType = "rpcv01"
	ProviderRPCv02  ProviderType = "rpcv02"
	ProviderGateway ProviderType = "gateway"
)

type Account struct {
	rpcv01         *rpcv01.Provider
	rpcv02         *rpcv02.Provider
	sequencer      *gateway.GatewayProvider
	provider       ProviderType
	chainId        string
	AccountAddress types.Felt
	sender         types.Felt
	ks             Keystore
	version        uint64
	plugin         AccountPlugin
}

type AccountOption struct {
	AccountPlugin AccountPlugin
	version       uint64
}

type AccountOptionFunc func(types.Felt, types.Felt) (AccountOption, error)

func AccountVersion0(types.Felt, types.Felt) (AccountOption, error) {
	return AccountOption{
		version: uint64(0),
	}, nil
}

func AccountVersion1(types.Felt, types.Felt) (AccountOption, error) {
	return AccountOption{
		version: uint64(1),
	}, nil
}

func newAccount(sender, address types.Felt, ks Keystore, options ...AccountOptionFunc) (*Account, error) {
	var accountPlugin AccountPlugin
	version := uint64(0)
	for _, o := range options {
		opt, err := o(sender, address)
		if err != nil {
			return nil, err
		}
		if opt.version != 0 {
			version = opt.version
		}
		if opt.AccountPlugin != nil {
			if accountPlugin != nil {
				return nil, errors.New("multiple plugins not supported")
			}
			accountPlugin = opt.AccountPlugin
		}
	}
	return &Account{
		AccountAddress: address,
		version:        version,
		plugin:         accountPlugin,
		ks:             ks,
		sender:         sender,
	}, nil
}

func setAccountProvider(account *Account, provider interface{}) error {
	switch p := provider.(type) {
	case *rpcv01.Provider:
		chainID, err := p.ChainID(context.Background())
		if err != nil {
			return err
		}
		account.chainId = chainID
		account.provider = ProviderRPCv01
		account.rpcv01 = p
		return nil
	case *rpcv02.Provider:
		chainID, err := p.ChainID(context.Background())
		if err != nil {
			return err
		}
		account.chainId = chainID
		account.provider = ProviderRPCv02
		account.rpcv02 = p
		return nil
	}
	return errors.New("unsupported provider")
}

func NewRPCAccount[Provider *rpcv02.Provider](sender, address types.Felt, ks Keystore, provider Provider, options ...AccountOptionFunc) (*Account, error) {
	account, err := newAccount(sender, address, ks, options...)
	if err != nil {
		return nil, err
	}
	err = setAccountProvider(account, provider)
	return account, err
}

func NewGatewayAccount(sender, address types.Felt, ks Keystore, provider *gateway.GatewayProvider, options ...AccountOptionFunc) (*Account, error) {
	account, err := newAccount(sender, address, ks, options...)
	if err != nil {
		return nil, err
	}
	chainID, err := provider.ChainID(context.Background())
	if err != nil {
		return nil, err
	}
	account.chainId = chainID
	account.provider = ProviderGateway
	account.sequencer = provider
	return account, nil
}

func (account *Account) Call(ctx context.Context, call types.FunctionCall) ([]string, error) {
	switch account.provider {
	case ProviderRPCv01:
		if account.rpcv01 == nil {
			return nil, ErrUnsupportedAccount
		}
		return account.rpcv01.Call(ctx, call, rpcv01.WithBlockTag("latest"))
	case ProviderRPCv02:
		if account.rpcv02 == nil {
			return nil, ErrUnsupportedAccount
		}
		return account.rpcv02.Call(ctx, call, rpcv02.WithBlockTag("latest"))
	case ProviderGateway:
		if account.sequencer == nil {
			return nil, ErrUnsupportedAccount
		}
		return account.sequencer.Call(ctx, call, "latest")
	}
	return nil, ErrUnsupportedAccount
}

func (account *Account) TransactionHash(calls []types.FunctionCall, details types.ExecuteDetails) (*big.Int, error) {

	var callArray []*big.Int
	switch account.version {
	case 1:
		callArray = fmtCalldata(calls)
	default:
		return nil, fmt.Errorf("version %d unsupported", account.version)
	}
	cdHash, err := Curve.ComputeHashOnElements(callArray)
	if err != nil {
		return nil, err
	}

	var multiHashData []*big.Int
	switch account.version {
	case 1:
		multiHashData = []*big.Int{
			types.UTF8StrToBig(TRANSACTION_PREFIX),
			big.NewInt(int64(account.version)),
			account.AccountAddress.Big(),
			big.NewInt(0),
			cdHash,
			details.MaxFee,
			types.UTF8StrToBig(account.chainId),
			details.Nonce,
		}
	default:
		return nil, fmt.Errorf("version %d unsupported", account.version)
	}
	return Curve.ComputeHashOnElements(multiHashData)
}

func (account *Account) estimateFeeHash(calls []types.FunctionCall, details types.ExecuteDetails, version *big.Int) (*big.Int, error) {
	var callArray []*big.Int
	switch account.version {
	case 1:
		callArray = fmtCalldata(calls)
	default:
		return nil, fmt.Errorf("version %d unsupported", account.version)
	}
	cdHash, err := Curve.ComputeHashOnElements(callArray)
	if err != nil {
		return nil, err
	}
	var multiHashData []*big.Int
	switch account.version {
	case 1:
		multiHashData = []*big.Int{
			types.UTF8StrToBig(TRANSACTION_PREFIX),
			version,
			account.AccountAddress.Big(),
			big.NewInt(0),
			cdHash,
			details.MaxFee,
			types.UTF8StrToBig(account.chainId),
			details.Nonce,
		}
	default:
		return nil, fmt.Errorf("version %d unsupported", account.version)
	}
	return Curve.ComputeHashOnElements(multiHashData)
}

func (account *Account) Nonce(ctx context.Context) (*big.Int, error) {
	switch account.version {
	case 1:
		switch account.provider {
		case ProviderRPCv01:
			nonce, err := account.rpcv01.Nonce(
				ctx,
				types.HexToHash(account.AccountAddress),
			)
			if err != nil {
				return nil, err
			}
			n, ok := big.NewInt(0).SetString(*nonce, 0)
			if !ok {
				return nil, errors.New("nonce error")
			}
			return n, nil
		case ProviderRPCv02:
			nonce, err := account.rpcv02.Nonce(
				ctx,
				rpcv02.WithBlockTag("latest"),
				account.AccountAddress,
			)
			if err != nil {
				return nil, err
			}
			n, ok := big.NewInt(0).SetString(*nonce, 0)
			if !ok {
				return nil, errors.New("nonce error")
			}
			return n, nil
		case ProviderGateway:
			return account.sequencer.Nonce(ctx, account.AccountAddress.String(), "latest")
		}
	}
	return nil, fmt.Errorf("version %d unsupported", account.version)
}

func (account *Account) prepFunctionInvoke(ctx context.Context, messageType string, calls []types.FunctionCall, details types.ExecuteDetails) (*types.FunctionInvoke, error) {
	if messageType != "invoke" && messageType != "estimate" {
		return nil, errors.New("unsupported message type")
	}
	nonce := details.Nonce
	var err error
	if details.Nonce == nil {
		nonce, err = account.Nonce(ctx)
		if err != nil {
			return nil, err
		}
	}
	maxFee := MAX_FEE
	if details.MaxFee != nil {
		maxFee = details.MaxFee
	}
	if account.plugin != nil {
		call, err := account.plugin.PluginCall(calls)
		if err != nil {
			return nil, err
		}
		calls = append([]types.FunctionCall{call}, calls...)
	}
	// version, _ := big.NewInt(0).SetString("0x100000000000000000000000000000000", 0)
	version, _ := big.NewInt(0).SetString("0x0", 0)
	var txHash *big.Int
	switch messageType {
	case "invoke":
		version = big.NewInt(int64(account.version))
		txHash, err = account.TransactionHash(
			calls,
			types.ExecuteDetails{
				Nonce:  nonce,
				MaxFee: maxFee,
			},
		)
		if err != nil {
			return nil, err
		}
	case "estimate":
		if account.version == 1 {
			// version, _ = big.NewInt(0).SetString("0x100000000000000000000000000000001", 0)
			version, _ = big.NewInt(0).SetString("0x1", 0)
		}
		txHash, err = account.estimateFeeHash(
			calls,
			types.ExecuteDetails{
				Nonce:  nonce,
				MaxFee: maxFee,
			},
			version,
		)
		if err != nil {
			return nil, err
		}
	}
	s1, s2, err := account.ks.Sign(ctx, account.sender.String(), txHash)
	if err != nil {
		return nil, err
	}

	switch account.version {
	case 1:
		calldata := fmtCalldataStrings(calls)
		return &types.FunctionInvoke{
			MaxFee:        maxFee,
			Version:       version,
			Signature:     types.Signature{s1, s2},
			SenderAddress: account.AccountAddress,
			Calldata:      calldata,
			Nonce:         nonce,
		}, nil
	}
	return nil, ErrUnsupportedAccount
}

func (account *Account) EstimateFee(ctx context.Context, calls []types.FunctionCall, details types.ExecuteDetails) (*types.FeeEstimate, error) {
	call, err := account.prepFunctionInvoke(ctx, "estimate", calls, details)
	if err != nil {
		return nil, err
	}
	switch account.provider {
	case ProviderRPCv01:
		return account.rpcv01.EstimateFee(ctx, *call, rpcv01.WithBlockTag("latest"))
	case ProviderRPCv02:
		signature := []string{}
		for _, v := range call.Signature {
			signature = append(signature, fmt.Sprintf("0x%x", v))
		}
		switch account.version {
		case 1:
			return account.rpcv02.EstimateFee(ctx, rpcv02.BroadcastedInvokeV1Transaction{
				BroadcastedTxnCommonProperties: rpcv02.BroadcastedTxnCommonProperties{
					MaxFee:    call.MaxFee,
					Version:   rpcv02.TransactionV1,
					Signature: signature,
					Nonce:     call.Nonce,
					Type:      "INVOKE",
				},
				Calldata:      call.Calldata,
				SenderAddress: account.AccountAddress,
			}, rpcv02.WithBlockTag("latest"))
		}
	case ProviderGateway:
		return account.sequencer.EstimateFee(ctx, *call, "latest")
	}
	return nil, ErrUnsupportedAccount
}

func (account *Account) Execute(ctx context.Context, calls []types.FunctionCall, details types.ExecuteDetails) (*types.AddInvokeTransactionOutput, error) {
	maxFee := details.MaxFee
	if maxFee == nil {
		estimate, err := account.EstimateFee(ctx, calls, details)
		if err != nil {
			return nil, err
		}
		fmt.Printf("fee %+v\n", estimate)
		v, ok := big.NewInt(0).SetString(string(estimate.OverallFee), 0)
		if !ok {
			return nil, errors.New("could not match OverallFee to big.Int")
		}
		maxFee = v.Mul(v, big.NewInt(2))
	}
	details.MaxFee = maxFee
	call, err := account.prepFunctionInvoke(ctx, "invoke", calls, details)
	if err != nil {
		return nil, err
	}
	switch account.provider {
	case ProviderRPCv01:
		signature := []string{}
		for _, k := range call.Signature {
			signature = append(signature, fmt.Sprintf("0x%s", k.Text(16)))
		}
		return account.rpcv01.AddInvokeTransaction(
			context.Background(),
			call.FunctionCall,
			signature,
			fmt.Sprintf("0x%s", maxFee.Text(16)),
			fmt.Sprintf("0x%d", account.version),
			call.Nonce,
		)
	case ProviderRPCv02:
		signature := []string{}
		for _, v := range call.Signature {
			signature = append(signature, fmt.Sprintf("0x%x", v))
		}
		switch account.version {
		case 1:
			return account.rpcv02.AddInvokeTransaction(ctx, rpcv02.BroadcastedInvokeV1Transaction{
				BroadcastedTxnCommonProperties: rpcv02.BroadcastedTxnCommonProperties{
					MaxFee:    call.MaxFee,
					Version:   rpcv02.TransactionV1,
					Signature: signature,
					Nonce:     call.Nonce,
					Type:      "INVOKE",
				},
				SenderAddress: account.AccountAddress,
				Calldata:      call.Calldata,
			})
		}
	case ProviderGateway:
		return account.sequencer.Invoke(
			context.Background(),
			*call,
		)
	}
	return nil, ErrUnsupportedAccount
}

func (account *Account) Declare(ctx context.Context, classHash string, contract types.ContractClass, details types.ExecuteDetails) (types.AddDeclareResponse, error) {
	switch account.provider {
	case ProviderRPCv02:
		panic("unsupported")
	case ProviderGateway:
		version := big.NewInt(1)
		nonce := details.Nonce
		var err error
		if details.Nonce == nil {
			nonce, err = account.Nonce(ctx)
			if err != nil {
				return types.AddDeclareResponse{}, err
			}
		}
		// TODO: use max fee estimation instead
		maxFee := MAX_FEE
		if details.MaxFee != nil {
			maxFee = details.MaxFee
		}

		// TODO: extract as declareHash
		hash, _ := big.NewInt(0).SetString(classHash, 0)
		calldataHash, err := Curve.ComputeHashOnElements([]*big.Int{hash})
		if err != nil {
			return types.AddDeclareResponse{}, err
		}
		var multiHashData []*big.Int
		switch account.version {
		case 1:
			multiHashData = []*big.Int{
				types.UTF8StrToBig(DECLARE_PREFIX),
				version,
				account.AccountAddress.Big(),
				big.NewInt(0),
				calldataHash,
				maxFee,
				types.UTF8StrToBig(account.chainId),
				nonce, // TODO: also include compiledClassHash for cairo 1.0?
			}
		default:
			return types.AddDeclareResponse{}, fmt.Errorf("version %d unsupported", account.version)
		}
		txHash, err := Curve.ComputeHashOnElements(multiHashData)
		// TODO: end extract as declareHash
		if err != nil {
			return types.AddDeclareResponse{}, err
		}

		s1, s2, err := account.ks.Sign(ctx, account.sender.String(), txHash)
		if err != nil {
			return types.AddDeclareResponse{}, err
		}
		signature := []string{}
		signature = append(signature, s1.String())
		signature = append(signature, s2.String())

		request := gateway.DeclareRequest{
			SenderAddress: account.AccountAddress,
			Version:       fmt.Sprintf("0x%x", version),
			MaxFee:        fmt.Sprintf("0x%x", maxFee),
			Nonce:         fmt.Sprintf("0x%x", nonce),
			Signature:     signature,
			ContractClass: contract,
			Type:          "DECLARE",
		}
		return account.sequencer.Declare(ctx, contract, request)
	}
	return types.AddDeclareResponse{}, ErrUnsupportedAccount
}

// Deploys a declared contract using the UDC.
// TODO: use types.DeployRequest{} as input for salt + calldata (remove contract_definition)
func (account *Account) Deploy(ctx context.Context, classHash string, details types.ExecuteDetails) (*types.AddDeployResponse, error) {
	// TODO: allow passing salt in
	salt, err := Curve.GetRandomPrivateKey()
	if err != nil {
		return nil, err
	}

	unique := true
	calldata := []string{}

	uniqueInt := big.NewInt(0)
	if unique {
		uniqueInt = big.NewInt(1)
	}

	deployerAddress := types.StrToFelt("0x41a78e741e5af2fec34b695679bc6891742439f7afb8484ecd7766661ad02bf") // UDC
	tx, err := account.Execute(ctx, []types.FunctionCall{
		{
			ContractAddress:    deployerAddress,
			EntryPointSelector: "deployContract",
			Calldata: append([]string{
				classHash,
				fmt.Sprintf("0x%x", salt),
				fmt.Sprintf("0x%x", uniqueInt), // unique
				fmt.Sprintf("0x%x", len(calldata)),
			}, calldata...),
		},
	}, details)
	if err != nil {
		return nil, err
	}

	// Calculate the resulting contract address
	constructorCalldata := []*big.Int{}
	for _, value := range calldata {
		constructorCalldata = append(constructorCalldata, types.SNValToBN(value))
	}
	constructorCalldataHash, err := Curve.ComputeHashOnElements(constructorCalldata)
	if err != nil {
		return nil, err
	}

	if unique {
		salt, err = Curve.PedersenHash([]*big.Int{
			account.AccountAddress.Big(),
			salt,
		})
		if err != nil {
			return nil, err
		}
	}

	prefix := types.HexToBN("0x535441524b4e45545f434f4e54524143545f41444452455353")

	contractAddress, err := Curve.ComputeHashOnElements([]*big.Int{
		prefix,                // CONTRACT_ADDRESS_PREFIX
		deployerAddress.Big(), // TODO: 0 if !unique
		salt,
		types.HexToBN(classHash),
		constructorCalldataHash,
	})
	if err != nil {
		return nil, err
	}

	return &types.AddDeployResponse{
		TransactionHash: tx.TransactionHash,
		ContractAddress: fmt.Sprintf("0x%x", contractAddress),
	}, nil
}
