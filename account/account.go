package account

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/contracts"
	"github.com/NethermindEth/starknet.go/rpc/rpcv10"
	"github.com/NethermindEth/starknet.go/rpc/rpcv9"
	"github.com/NethermindEth/starknet.go/rpc/types"
	"github.com/NethermindEth/starknet.go/types/constraints"
)

var (
	ErrTxnTypeUnSupported    = errors.New("unsupported transaction type")
	ErrTxnVersionUnSupported = errors.New("unsupported transaction version")
)

// @changed
type AccountInterface interface {
	BuildAndEstimateDeployAccountTxn(
		ctx context.Context,
		salt *felt.Felt,
		classHash *felt.Felt,
		constructorCalldata []*felt.Felt,
		opts *TxnOptions,
	) (*types.BroadcastDeployAccountTxnV3, *felt.Felt, error)
	BuildAndSendInvokeTxn(
		ctx context.Context,
		functionCalls []types.InvokeFunctionCall,
		opts *TxnOptions,
	) (types.TransactionResponse, error)
	BuildAndSendDeclareTxn(
		ctx context.Context,
		casmClass *contracts.CasmClass,
		contractClass *contracts.ContractClass,
		opts *TxnOptions,
	) (types.TransactionResponse, error)
	DeployContractWithUDC(
		ctx context.Context,
		classHash *felt.Felt,
		constructorCalldata []*felt.Felt,
		txnOpts *TxnOptions,
		udcOpts *UDCOptions,
	) (types.TransactionResponse, *felt.Felt, error)
	Nonce(ctx context.Context) (*felt.Felt, error)
	SendTransaction(ctx context.Context, txn types.BroadcastTxn) (types.TransactionResponse, error)
	Sign(ctx context.Context, msg *felt.Felt) ([]*felt.Felt, error)
	SignInvokeTransaction(ctx context.Context, tx types.InvokeTxnType) error
	SignDeployAccountTransaction(
		ctx context.Context,
		tx types.DeployAccountType,
		precomputeAddress *felt.Felt,
	) error
	SignDeclareTransaction(ctx context.Context, tx types.DeclareTxnType) error
	TransactionHashInvoke(invokeTxn types.InvokeTxnType) (*felt.Felt, error)
	TransactionHashDeployAccount(
		tx types.DeployAccountType,
		contractAddress *felt.Felt,
	) (*felt.Felt, error)
	TransactionHashDeclare(tx types.DeclareTxnType) (*felt.Felt, error)
	Verify(msgHash *felt.Felt, signature []*felt.Felt) (bool, error)
	WaitForTransactionReceipt(
		ctx context.Context,
		transactionHash *felt.Felt,
		pollInterval time.Duration,
	) (*types.TransactionReceiptWithBlockInfo, error)
}

var _ AccountInterface = (*Account)(nil)

// @changed
type Account struct {
	// TODO: in the future, make all fields private and add getter methods
	// @changed
	provider     providerWrapper
	ChainID      *felt.Felt
	Address      *felt.Felt
	publicKey    string
	CairoVersion CairoVersion
	ks           Keystore
}

// @new
func (a *Account) ProviderAsV9() rpcv9.RPCProvider {
	return a.provider.AsV9()
}

// @new
func (a *Account) ProviderAsV10() rpcv10.RPCProvider {
	return a.provider.AsV10()
}

// CairoVersion represents the version of Cairo used by the account contract.
type CairoVersion int

// Cairo version constants
const (
	// CairoV0 represents Cairo 0 contracts
	CairoV0 CairoVersion = 0
	// CairoV2 represents Cairo 2 contracts
	CairoV2 CairoVersion = 2
)

// @changed
// NewAccount creates a new Account instance.
//
// Parameters:
//   - provider: the provider to use
//   - accountAddress: the account address
//   - publicKey: the public key of the account
//   - keystore: the keystore to use
//   - cairoVersion: the cairo version of the account (CairoVersion0 or CairoVersion2)
//
// It returns:
//   - *Account: a pointer to newly created Account
//   - error: an error if any
//
// @todo make it a generic function
func NewAccount(
	// @changed
	provider *rpcv10.Provider,
	accountAddress *felt.Felt,
	publicKey string,
	keystore Keystore,
	cairoVersion CairoVersion,
) (*Account, error) {
	chainID, err := provider.ChainID(context.Background())
	if err != nil {
		return nil, err
	}

	account := &Account{
		provider:     newWrapperFrom(provider),
		Address:      accountAddress,
		publicKey:    publicKey,
		ks:           keystore,
		CairoVersion: cairoVersion,
		ChainID:      new(felt.Felt).SetBytes([]byte(chainID)),
	}

	return account, nil
}

// Nonce retrieves the nonce for the account's contract address.
func (account *Account) Nonce(ctx context.Context) (*felt.Felt, error) {
	return account.provider.Nonce(ctx, types.WithBlockTag("pre_confirmed"), account.Address)
}

// PrecomputeAccountAddress calculates the precomputed address for an account.
// ref: https://docs.starknet.io/architecture-and-concepts/smart-contracts/contract-address/
//
// Parameters:
//   - salt: the salt for the address of the deployed contract
//   - classHash: the class hash of the contract to be deployed
//   - constructorCalldata: the parameters passed to the constructor
//
// Returns:
//   - *felt.Felt: the precomputed address as a *felt.Felt
//   - error: an error if any
func PrecomputeAccountAddress(
	salt, classHash *felt.Felt,
	constructorCalldata []*felt.Felt,
) *felt.Felt {
	return contracts.PrecomputeAddress(&felt.Zero, salt, classHash, constructorCalldata)
}

// FmtCalldata generates the formatted calldata for the given function calls and Cairo version.
//
// Parameters:
//   - fnCalls: a slice of types.FunctionCall representing the function calls.
//
// Returns:
//   - a slice of *felt.Felt representing the formatted calldata.
//   - an error if Cairo version is not supported.
func (account *Account) FmtCalldata(fnCalls []types.FunctionCall) ([]*felt.Felt, error) {
	switch account.CairoVersion {
	case CairoV0:
		return FmtCallDataCairo0(fnCalls), nil
	case CairoV2:
		return FmtCallDataCairo2(fnCalls), nil
	default:
		return nil, fmt.Errorf("account cairo version '%d' not supported", account.CairoVersion)
	}
}

// FmtCallDataCairo0 generates a slice of *felt.Felt that represents the
// calldata for the given function calls in Cairo 0 format.
//
// Parameters:
//   - fnCalls: a slice of types.FunctionCall containing the function calls.
//
// Returns:
//   - a slice of *felt.Felt representing the generated calldata.
//
// https://github.com/project3fusion/StarkSharp/blob/main/StarkSharp/StarkSharp.Rpc/Modules/Transactions/Hash/TransactionHash.cs#L27
//
//nolint:lll // The link would be unclickable if we break the line.
func FmtCallDataCairo0(callArray []types.FunctionCall) []*felt.Felt {
	calldata := make([]*felt.Felt, 0, 10) //nolint:mnd // Randomly chosen
	calls := make([]*felt.Felt, 0, 10)    //nolint:mnd // Randomly chosen

	calldata = append(calldata, new(felt.Felt).SetUint64(uint64(len(callArray))))

	offset := uint64(0)
	for _, call := range callArray {
		calldata = append(
			calldata,
			call.ContractAddress,
			call.EntryPointSelector,
			new(felt.Felt).SetUint64(offset),
		)
		callDataLen := uint64(len(call.Calldata))
		calldata = append(calldata, new(felt.Felt).SetUint64(callDataLen))
		offset += callDataLen

		calls = append(calls, call.Calldata...)
	}

	calldata = append(calldata, new(felt.Felt).SetUint64(offset))
	calldata = append(calldata, calls...)

	return calldata
}

// FmtCallDataCairo2 generates the calldata for the given function calls for
// Cairo 2 contracts.
//
// Parameters:
//   - fnCalls: a slice of types.FunctionCall containing the function calls.
//
// Returns:
//   - a slice of *felt.Felt representing the generated calldata.
//
// https://github.com/project3fusion/StarkSharp/blob/main/StarkSharp/StarkSharp.Rpc/Modules/Transactions/Hash/TransactionHash.cs#L22
//
//nolint:lll // The link would be unclickable if we break the line.
func FmtCallDataCairo2(callArray []types.FunctionCall) []*felt.Felt {
	result := make([]*felt.Felt, 0, 10) //nolint:mnd // Randomly chosen

	result = append(result, new(felt.Felt).SetUint64(uint64(len(callArray))))

	for _, call := range callArray {
		result = append(result, call.ContractAddress, call.EntryPointSelector)

		callDataLen := uint64(len(call.Calldata))
		result = append(result, new(felt.Felt).SetUint64(callDataLen))

		result = append(result, call.Calldata...)
	}

	return result
}

func makeEmptyResourceBM[
	u64 constraints.U64,
	u128 constraints.U128,
	RB constraints.ResourceBounds[u64, u128],
	RBM constraints.ResourceBoundsMapping[u64, u128, RB],
](resourceBounds *RBM) *RBM {
	return &RBM{
		L1Gas: RB{
			MaxAmount:       "0x0",
			MaxPricePerUnit: "0x0",
		},
		L1DataGas: RB{
			MaxAmount:       "0x0",
			MaxPricePerUnit: "0x0",
		},
		L2Gas: RB{
			MaxAmount:       "0x0",
			MaxPricePerUnit: "0x0",
		},
	}
}
