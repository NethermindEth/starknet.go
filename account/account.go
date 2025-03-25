package account

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/NethermindEth/juno/core/crypto"
	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/contracts"
	"github.com/NethermindEth/starknet.go/curve"
	"github.com/NethermindEth/starknet.go/hash"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"
)

var (
	ErrNotAllParametersSet   = errors.New("not all neccessary parameters have been set")
	ErrTxnTypeUnSupported    = errors.New("unsupported transaction type")
	ErrTxnVersionUnSupported = errors.New("unsupported transaction version")
	ErrFeltToBigInt          = errors.New("felt to BigInt error")
)

var (
	PREFIX_TRANSACTION    = new(felt.Felt).SetBytes([]byte("invoke"))
	PREFIX_DECLARE        = new(felt.Felt).SetBytes([]byte("declare"))
	PREFIX_DEPLOY_ACCOUNT = new(felt.Felt).SetBytes([]byte("deploy_account"))
)

//go:generate mockgen -destination=../mocks/mock_account.go -package=mocks -source=account.go AccountInterface
type AccountInterface interface {
	BuildAndEstimateDeployAccountTxn(ctx context.Context, salt *felt.Felt, classHash *felt.Felt, constructorCalldata []*felt.Felt, multiplier float64) (*rpc.BroadcastDeployAccountTxnV3, *felt.Felt, error)
	BuildAndSendInvokeTxn(ctx context.Context, functionCalls []rpc.InvokeFunctionCall, multiplier float64) (*rpc.AddInvokeTransactionResponse, error)
	BuildAndSendDeclareTxn(ctx context.Context, casmClass *contracts.CasmClass, contractClass *contracts.ContractClass, multiplier float64) (*rpc.AddDeclareTransactionResponse, error)
	Nonce(ctx context.Context) (*felt.Felt, error)
	SendTransaction(ctx context.Context, txn rpc.BroadcastTxn) (*rpc.TransactionResponse, error)
	Sign(ctx context.Context, msg *felt.Felt) ([]*felt.Felt, error)
	SignInvokeTransaction(ctx context.Context, tx rpc.InvokeTxnType) error
	SignDeployAccountTransaction(ctx context.Context, tx rpc.DeployAccountType, precomputeAddress *felt.Felt) error
	SignDeclareTransaction(ctx context.Context, tx rpc.DeclareTxnType) error
	TransactionHashInvoke(invokeTxn rpc.InvokeTxnType) (*felt.Felt, error)
	TransactionHashDeployAccount(tx rpc.DeployAccountType, contractAddress *felt.Felt) (*felt.Felt, error)
	TransactionHashDeclare(tx rpc.DeclareTxnType) (*felt.Felt, error)
	WaitForTransactionReceipt(ctx context.Context, transactionHash *felt.Felt, pollInterval time.Duration) (*rpc.TransactionReceiptWithBlockInfo, error)
}

var _ AccountInterface = &Account{}

type Account struct {
	Provider       rpc.RpcProvider
	ChainId        *felt.Felt
	AccountAddress *felt.Felt
	publicKey      string
	CairoVersion   int
	ks             Keystore
}

// NewAccount creates a new Account instance.
//
// Parameters:
// - provider: is the provider of type rpc.RpcProvider
// - accountAddress: is the account address of type *felt.Felt
// - publicKey: is the public key of type string
// - keystore: is the keystore of type Keystore
// It returns:
// - *Account: a pointer to newly created Account
// - error: an error if any
func NewAccount(provider rpc.RpcProvider, accountAddress *felt.Felt, publicKey string, keystore Keystore, cairoVersion int) (*Account, error) {
	account := &Account{
		Provider:       provider,
		AccountAddress: accountAddress,
		publicKey:      publicKey,
		ks:             keystore,
		CairoVersion:   cairoVersion,
	}

	chainID, err := provider.ChainID(context.Background())
	if err != nil {
		return nil, err
	}
	account.ChainId = new(felt.Felt).SetBytes([]byte(chainID))

	return account, nil
}

// Nonce retrieves the nonce for the account's contract address.
func (account *Account) Nonce(ctx context.Context) (*felt.Felt, error) {
	return account.Provider.Nonce(context.Background(), rpc.WithBlockTag("pending"), account.AccountAddress)
}

// BuildAndSendInvokeTxn builds and sends a v3 invoke transaction with the given function calls.
// It automatically calculates the nonce, formats the calldata, estimates fees, and signs the transaction with the account's private key.
//
// Parameters:
//   - ctx: The context.Context for the request.
//   - functionCalls: A slice of rpc.InvokeFunctionCall representing the function calls for the transaction, allowing either single or
//     multiple function calls in the same transaction.
//   - multiplier: A safety factor for fee estimation that helps prevent transaction failures due to
//     fee fluctuations. It multiplies both the max amount and max price per unit by this value.
//     A value of 1.5 (50% buffer) is recommended to balance between transaction success rate and
//     avoiding excessive fees. Higher values provide more safety margin but may result in overpayment.
//
// Returns:
//   - *rpc.AddInvokeTransactionResponse: the response of the submitted transaction.
//   - error: An error if the transaction building fails.
func (account *Account) BuildAndSendInvokeTxn(ctx context.Context, functionCalls []rpc.InvokeFunctionCall, multiplier float64) (*rpc.AddInvokeTransactionResponse, error) {
	nonce, err := account.Provider.Nonce(ctx, rpc.WithBlockTag("pending"), account.AccountAddress)
	if err != nil {
		return nil, err
	}

	callData, err := account.FmtCalldata(utils.InvokeFuncCallsToFunctionCalls(functionCalls))
	if err != nil {
		return nil, err
	}

	// building and signing the txn, as it needs a signature to estimate the fee
	broadcastInvokeTxnV3 := utils.BuildInvokeTxn(account.AccountAddress, nonce, callData, makeResourceBoundsMapWithZeroValues())
	err = account.SignInvokeTransaction(ctx, &broadcastInvokeTxnV3.InvokeTxnV3)
	if err != nil {
		return nil, err
	}

	// estimate txn fee
	estimateFee, err := account.Provider.EstimateFee(ctx, []rpc.BroadcastTxn{broadcastInvokeTxnV3}, []rpc.SimulationFlag{rpc.SKIP_VALIDATE}, rpc.WithBlockTag("pending"))
	if err != nil {
		return nil, err
	}
	txnFee := estimateFee[0]
	broadcastInvokeTxnV3.ResourceBounds = utils.FeeEstToResBoundsMap(txnFee, multiplier)

	// signing the txn again with the estimated fee, as the fee value is used in the txn hash calculation
	err = account.SignInvokeTransaction(ctx, &broadcastInvokeTxnV3.InvokeTxnV3)
	if err != nil {
		return nil, err
	}

	txnResponse, err := account.Provider.AddInvokeTransaction(ctx, broadcastInvokeTxnV3)
	if err != nil {
		return nil, err
	}

	return txnResponse, nil
}

// BuildAndSendDeclareTxn builds and sends a v3 declare transaction.
// It automatically calculates the nonce, formats the calldata, estimates fees, and signs the transaction with the account's private key.
//
// Parameters:
//   - ctx: The context.Context for the request.
//   - casmClass: The casm class of the contract to be declared
//   - contractClass: The sierra contract class of the contract to be declared
//   - multiplier: A safety factor for fee estimation that helps prevent transaction failures due to
//     fee fluctuations. It multiplies both the max amount and max price per unit by this value.
//     A value of 1.5 (50% buffer) is recommended to balance between transaction success rate and
//     avoiding excessive fees. Higher values provide more safety margin but may result in overpayment.
//
// Returns:
//   - *rpc.AddDeclareTransactionResponse: the response of the submitted transaction.
//   - error: An error if the transaction building fails.
func (account *Account) BuildAndSendDeclareTxn(
	ctx context.Context,
	casmClass *contracts.CasmClass,
	contractClass *contracts.ContractClass,
	multiplier float64,
) (*rpc.AddDeclareTransactionResponse, error) {
	nonce, err := account.Provider.Nonce(ctx, rpc.WithBlockTag("pending"), account.AccountAddress)
	if err != nil {
		return nil, err
	}

	// building and signing the txn, as it needs a signature to estimate the fee
	broadcastDeclareTxnV3, err := utils.BuildDeclareTxn(account.AccountAddress, casmClass, contractClass, nonce, makeResourceBoundsMapWithZeroValues())
	if err != nil {
		return nil, err
	}
	err = account.SignDeclareTransaction(ctx, &broadcastDeclareTxnV3)
	if err != nil {
		return nil, err
	}

	// estimate txn fee
	estimateFee, err := account.Provider.EstimateFee(ctx, []rpc.BroadcastTxn{broadcastDeclareTxnV3}, []rpc.SimulationFlag{rpc.SKIP_VALIDATE}, rpc.WithBlockTag("pending"))
	if err != nil {
		return nil, err
	}
	txnFee := estimateFee[0]
	broadcastDeclareTxnV3.ResourceBounds = utils.FeeEstToResBoundsMap(txnFee, multiplier)

	// signing the txn again with the estimated fee, as the fee value is used in the txn hash calculation
	err = account.SignDeclareTransaction(ctx, &broadcastDeclareTxnV3)
	if err != nil {
		return nil, err
	}

	txnResponse, err := account.Provider.AddDeclareTransaction(ctx, broadcastDeclareTxnV3)
	if err != nil {
		return nil, err
	}

	return txnResponse, nil
}

// BuildAndEstimateDeployAccountTxn builds and signs a v3 deploy account transaction, estimates the fee, and computes the address.
//
// This function doesn't send the transaction because the precomputed account address requires funding first. This address is calculated
// deterministically and returned by this function, and must be funded with the appropriate amount of STRK tokens. Without sufficient
// funds, the transaction will fail. See the 'examples/deployAccount/' for more details on how to do this.
//
// Parameters:
//   - ctx: The context.Context for the request.
//   - salt: the salt for the address of the deployed contract
//   - classHash: the class hash of the contract to be deployed
//   - constructorCalldata: the parameters passed to the constructor
//   - multiplier: A safety factor for fee estimation that helps prevent transaction failures due to
//     fee fluctuations. It multiplies both the max amount and max price per unit by this value.
//     A value of 1.5 (50% buffer) is recommended to balance between transaction success rate and
//     avoiding excessive fees. Higher values provide more safety margin but may result in overpayment.
//
// Returns:
//   - *rpc.BroadcastDeployAccountTxnV3: the transaction to be broadcasted, signed and with the estimated fee based on the multiplier
//   - *felt.Felt: the precomputed account address as a *felt.Felt, it needs to be funded with appropriate amount of tokens
//   - error: an error if any
func (account *Account) BuildAndEstimateDeployAccountTxn(
	ctx context.Context,
	salt *felt.Felt,
	classHash *felt.Felt,
	constructorCalldata []*felt.Felt,
	multiplier float64,
) (*rpc.BroadcastDeployAccountTxnV3, *felt.Felt, error) {
	// building and signing the txn, as it needs a signature to estimate the fee
	broadcastDepAccTxnV3 := utils.BuildDeployAccountTxn(&felt.Zero, salt, constructorCalldata, classHash, makeResourceBoundsMapWithZeroValues())

	precomputedAddress := PrecomputeAccountAddress(salt, classHash, constructorCalldata)

	// signing the txn, as it needs a signature to estimate the fee
	err := account.SignDeployAccountTransaction(ctx, &broadcastDepAccTxnV3.DeployAccountTxnV3, precomputedAddress)
	if err != nil {
		return nil, nil, err
	}

	// estimate txn fee
	estimateFee, err := account.Provider.EstimateFee(ctx, []rpc.BroadcastTxn{broadcastDepAccTxnV3}, []rpc.SimulationFlag{rpc.SKIP_VALIDATE}, rpc.WithBlockTag("pending"))
	if err != nil {
		return nil, nil, err
	}
	txnFee := estimateFee[0]
	broadcastDepAccTxnV3.ResourceBounds = utils.FeeEstToResBoundsMap(txnFee, multiplier)

	// signing the txn again with the estimated fee, as the fee value is used in the txn hash calculation
	err = account.SignDeployAccountTransaction(ctx, &broadcastDepAccTxnV3.DeployAccountTxnV3, precomputedAddress)
	if err != nil {
		return nil, nil, err
	}

	return &broadcastDepAccTxnV3, precomputedAddress, nil
}

// Sign signs the given felt message using the account's private key.
//
// Parameters:
// - ctx: is the context used for the signing operation
// - msg: is the felt message to be signed
// Returns:
// - []*felt.Felt: an array of signed felt messages
// - error: an error, if any
func (account *Account) Sign(ctx context.Context, msg *felt.Felt) ([]*felt.Felt, error) {
	msgBig := internalUtils.FeltToBigInt(msg)

	s1, s2, err := account.ks.Sign(ctx, account.publicKey, msgBig)
	if err != nil {
		return nil, err
	}
	s1Felt := internalUtils.BigIntToFelt(s1)
	s2Felt := internalUtils.BigIntToFelt(s2)

	return []*felt.Felt{s1Felt, s2Felt}, nil
}

// SignInvokeTransaction signs and invokes a transaction.
//
// Parameters:
//   - ctx: the context.Context for the function execution.
//   - invokeTx: the InvokeTxnV3 pointer representing the transaction to be invoked.
//
// Returns:
//   - error: an error if there was an error in the signing or invoking process
func (account *Account) SignInvokeTransaction(ctx context.Context, invokeTx rpc.InvokeTxnType) error {
	switch invoke := invokeTx.(type) {
	case *rpc.InvokeTxnV0:
		signature, err := signInvokeTransaction(ctx, account, invoke)
		if err != nil {
			return err
		}
		invoke.Signature = signature
	case *rpc.InvokeTxnV1:
		signature, err := signInvokeTransaction(ctx, account, invoke)
		if err != nil {
			return err
		}
		invoke.Signature = signature
	case *rpc.InvokeTxnV3:
		signature, err := signInvokeTransaction(ctx, account, invoke)
		if err != nil {
			return err
		}
		invoke.Signature = signature
	default:
		return fmt.Errorf("invalid invoke txn of type %T, did you pass a valid invoke txn pointer?", invoke)
	}

	return nil
}

// TODO: make func description
func signInvokeTransaction[T any](ctx context.Context, account *Account, invokeTx *T) ([]*felt.Felt, error) {
	txHash, err := account.TransactionHashInvoke(invokeTx)
	if err != nil {
		return nil, err
	}
	signature, err := account.Sign(ctx, txHash)
	if err != nil {
		return nil, err
	}

	return signature, nil
}

// SignDeployAccountTransaction signs a deploy account transaction.
//
// Parameters:
// - ctx: the context.Context for the function execution
// - tx: the *rpc.DeployAccountTxnV3 pointer representing the transaction to be signed
// - precomputeAddress: the precomputed address for the transaction
// Returns:
// - error: an error if any
func (account *Account) SignDeployAccountTransaction(ctx context.Context, tx rpc.DeployAccountType, precomputeAddress *felt.Felt) error {
	switch deployAcc := tx.(type) {
	case *rpc.DeployAccountTxn:
		signature, err := signDeployAccountTransaction(ctx, account, deployAcc, precomputeAddress)
		if err != nil {
			return err
		}
		deployAcc.Signature = signature
	case *rpc.DeployAccountTxnV3:
		signature, err := signDeployAccountTransaction(ctx, account, deployAcc, precomputeAddress)
		if err != nil {
			return err
		}
		deployAcc.Signature = signature
	default:
		return fmt.Errorf("invalid deploy account txn of type %T, did you pass a valid deploy account txn pointer?", deployAcc)
	}

	return nil
}

// TODO: make func description
func signDeployAccountTransaction[T any](ctx context.Context, account *Account, tx *T, precomputeAddress *felt.Felt) ([]*felt.Felt, error) {
	txHash, err := account.TransactionHashDeployAccount(*tx, precomputeAddress)
	if err != nil {
		return nil, err
	}
	signature, err := account.Sign(ctx, txHash)
	if err != nil {
		return nil, err
	}

	return signature, nil
}

// SignDeclareTransaction signs a DeclareTxnV2 transaction using the provided Account.
//
// Parameters:
// - ctx: the context.Context
// - tx: the pointer to a Declare or BroadcastDeclare txn
// Returns:
// - error: an error if any
func (account *Account) SignDeclareTransaction(ctx context.Context, tx rpc.DeclareTxnType) error {
	switch declare := tx.(type) {
	case *rpc.DeclareTxnV1:
		signature, err := signDeclareTransaction(ctx, account, declare)
		if err != nil {
			return err
		}
		declare.Signature = signature
	case *rpc.DeclareTxnV2:
		signature, err := signDeclareTransaction(ctx, account, declare)
		if err != nil {
			return err
		}
		declare.Signature = signature
	case *rpc.DeclareTxnV3:
		signature, err := signDeclareTransaction(ctx, account, declare)
		if err != nil {
			return err
		}
		declare.Signature = signature
	case *rpc.BroadcastDeclareTxnV3:
		signature, err := signDeclareTransaction(ctx, account, declare)
		if err != nil {
			return err
		}
		declare.Signature = signature
	default:
		return fmt.Errorf("invalid declare txn of type %T, did you pass a valid declare txn pointer?", declare)
	}

	return nil
}

// TODO: make func description
func signDeclareTransaction[T any](ctx context.Context, account *Account, tx *T) ([]*felt.Felt, error) {
	txHash, err := account.TransactionHashDeclare(*tx)
	if err != nil {
		return nil, err
	}
	signature, err := account.Sign(ctx, txHash)
	if err != nil {
		return nil, err
	}

	return signature, nil
}

// TransactionHashDeployAccount calculates the transaction hash for a deploy account transaction.
//
// Parameters:
// - tx: The deploy account transaction to calculate the hash for
// - contractAddress: The contract address as parameters as a *felt.Felt
// Returns:
// - *felt.Felt: the calculated transaction hash
// - error: an error if any
func (account *Account) TransactionHashDeployAccount(tx rpc.DeployAccountType, contractAddress *felt.Felt) (*felt.Felt, error) {

	// https://docs.starknet.io/architecture-and-concepts/network-architecture/transactions/#deploy_account_transaction
	switch txn := tx.(type) {
	// deployAccTxn v0, pointer and struct
	case *rpc.DeployAccountTxn:
		return TransactionHashDeployAccountV1(txn, contractAddress, account.ChainId)
	case rpc.DeployAccountTxn:
		return TransactionHashDeployAccountV1(&txn, contractAddress, account.ChainId)
	// deployAccTxn v3, pointer and struct
	case *rpc.DeployAccountTxnV3:
		return TransactionHashDeployAccountV3(txn, contractAddress, account.ChainId)
	case rpc.DeployAccountTxnV3:
		return TransactionHashDeployAccountV3(&txn, contractAddress, account.ChainId)
	default:
		return nil, fmt.Errorf("%w: got '%T' instead of a valid invoke txn type", ErrTxnTypeUnSupported, txn)
	}
}

// TODO: descriptions for all these functions
func TransactionHashDeployAccountV1(txn *rpc.DeployAccountTxn, contractAddress, chainId *felt.Felt) (*felt.Felt, error) {
	// https://docs.starknet.io/architecture-and-concepts/network-architecture/transactions/#v1_deprecated_hash_calculation_3
	calldata := []*felt.Felt{txn.ClassHash, txn.ContractAddressSalt}
	calldata = append(calldata, txn.ConstructorCalldata...)
	calldataHash := curve.PedersenArray(calldata...)

	versionFelt, err := new(felt.Felt).SetString(string(txn.Version))
	if err != nil {
		return nil, err
	}

	return hash.CalculateDeprecatedTransactionHashCommon(
		PREFIX_DEPLOY_ACCOUNT,
		versionFelt,
		contractAddress,
		&felt.Zero,
		calldataHash,
		txn.MaxFee,
		chainId,
		[]*felt.Felt{txn.Nonce},
	), nil
}

func TransactionHashDeployAccountV3(txn *rpc.DeployAccountTxnV3, contractAddress, chainId *felt.Felt) (*felt.Felt, error) {
	// https://docs.starknet.io/architecture-and-concepts/network-architecture/transactions/#v3_hash_calculation_3
	if txn.Version == "" || txn.ResourceBounds == (rpc.ResourceBoundsMapping{}) || txn.Nonce == nil || txn.PayMasterData == nil {
		return nil, ErrNotAllParametersSet
	}

	txnVersionFelt, err := new(felt.Felt).SetString(string(txn.Version))
	if err != nil {
		return nil, err
	}
	DAUint64, err := dataAvailabilityMode(txn.FeeMode, txn.NonceDataMode)
	if err != nil {
		return nil, err
	}
	tipUint64, err := txn.Tip.ToUint64()
	if err != nil {
		return nil, err
	}
	tipAndResourceHash, err := tipAndResourcesHash(tipUint64, txn.ResourceBounds)
	if err != nil {
		return nil, err
	}
	return crypto.PoseidonArray(
		PREFIX_DEPLOY_ACCOUNT,
		txnVersionFelt,
		contractAddress,
		tipAndResourceHash,
		crypto.PoseidonArray(txn.PayMasterData...),
		chainId,
		txn.Nonce,
		new(felt.Felt).SetUint64(DAUint64),
		crypto.PoseidonArray(txn.ConstructorCalldata...),
		txn.ClassHash,
		txn.ContractAddressSalt,
	), nil
}

// TransactionHashInvoke calculates the transaction hash for the given invoke transaction.
//
// Parameters:
// - tx: The invoke transaction to calculate the hash for. Can be of type InvokeTxnV0, InvokeTxnV1, or InvokeTxnV3.
// Returns:
// - *felt.Felt: The calculated transaction hash as a *felt.Felt
// - error: an error, if any

// If the transaction type is unsupported, the function returns an error.
func (account *Account) TransactionHashInvoke(tx rpc.InvokeTxnType) (*felt.Felt, error) {
	switch txn := tx.(type) {
	// invoke v0, pointer and struct
	case *rpc.InvokeTxnV0:
		return TransactionHashInvokeV0(txn, account.ChainId)
	case rpc.InvokeTxnV0:
		return TransactionHashInvokeV0(&txn, account.ChainId)
	// invoke v1, pointer and struct
	case *rpc.InvokeTxnV1:
		return TransactionHashInvokeV1(txn, account.ChainId)
	case rpc.InvokeTxnV1:
		return TransactionHashInvokeV1(&txn, account.ChainId)
	// invoke v3, pointer and struct
	case *rpc.InvokeTxnV3:
		return TransactionHashInvokeV3(txn, account.ChainId)
	case rpc.InvokeTxnV3:
		return TransactionHashInvokeV3(&txn, account.ChainId)
	default:
		return nil, fmt.Errorf("%w: got '%T' instead of a valid invoke txn type", ErrTxnTypeUnSupported, txn)
	}
}

// TODO: descriptions for all these functions
func TransactionHashInvokeV0(txn *rpc.InvokeTxnV0, chainId *felt.Felt) (*felt.Felt, error) {
	// https://docs.starknet.io/architecture-and-concepts/network-architecture/transactions/#v0_deprecated_hash_calculation
	if txn.Version == "" || len(txn.Calldata) == 0 || txn.MaxFee == nil || txn.EntryPointSelector == nil {
		return nil, ErrNotAllParametersSet
	}

	calldataHash := curve.PedersenArray(txn.Calldata...)
	txnVersionFelt, err := new(felt.Felt).SetString(string(txn.Version))
	if err != nil {
		return nil, err
	}
	return hash.CalculateDeprecatedTransactionHashCommon(
		PREFIX_TRANSACTION,
		txnVersionFelt,
		txn.ContractAddress,
		txn.EntryPointSelector,
		calldataHash,
		txn.MaxFee,
		chainId,
		[]*felt.Felt{},
	), nil
}

func TransactionHashInvokeV1(txn *rpc.InvokeTxnV1, chainId *felt.Felt) (*felt.Felt, error) {
	// https://docs.starknet.io/architecture-and-concepts/network-architecture/transactions/#v1_deprecated_hash_calculation
	if txn.Version == "" || len(txn.Calldata) == 0 || txn.Nonce == nil || txn.MaxFee == nil || txn.SenderAddress == nil {
		return nil, ErrNotAllParametersSet
	}

	calldataHash := curve.PedersenArray(txn.Calldata...)
	txnVersionFelt, err := new(felt.Felt).SetString(string(txn.Version))
	if err != nil {
		return nil, err
	}
	return hash.CalculateDeprecatedTransactionHashCommon(
		PREFIX_TRANSACTION,
		txnVersionFelt,
		txn.SenderAddress,
		&felt.Zero,
		calldataHash,
		txn.MaxFee,
		chainId,
		[]*felt.Felt{txn.Nonce},
	), nil
}

func TransactionHashInvokeV3(txn *rpc.InvokeTxnV3, chainId *felt.Felt) (*felt.Felt, error) {
	// https://github.com/starknet-io/SNIPs/blob/main/SNIPS/snip-8.md#protocol-changes
	// https://docs.starknet.io/architecture-and-concepts/network-architecture/transactions/#v3_hash_calculation
	if txn.Version == "" || txn.ResourceBounds == (rpc.ResourceBoundsMapping{}) || len(txn.Calldata) == 0 || txn.Nonce == nil || txn.SenderAddress == nil || txn.PayMasterData == nil || txn.AccountDeploymentData == nil {
		return nil, ErrNotAllParametersSet
	}

	txnVersionFelt, err := new(felt.Felt).SetString(string(txn.Version))
	if err != nil {
		return nil, err
	}
	DAUint64, err := dataAvailabilityMode(txn.FeeMode, txn.NonceDataMode)
	if err != nil {
		return nil, err
	}
	tipUint64, err := txn.Tip.ToUint64()
	if err != nil {
		return nil, err
	}
	tipAndResourceHash, err := tipAndResourcesHash(tipUint64, txn.ResourceBounds)
	if err != nil {
		return nil, err
	}
	return crypto.PoseidonArray(
		PREFIX_TRANSACTION,
		txnVersionFelt,
		txn.SenderAddress,
		tipAndResourceHash,
		crypto.PoseidonArray(txn.PayMasterData...),
		chainId,
		txn.Nonce,
		new(felt.Felt).SetUint64(DAUint64),
		crypto.PoseidonArray(txn.AccountDeploymentData...),
		crypto.PoseidonArray(txn.Calldata...),
	), nil
}

func tipAndResourcesHash(tip uint64, resourceBounds rpc.ResourceBoundsMapping) (*felt.Felt, error) {
	l1Bytes, err := resourceBounds.L1Gas.Bytes(rpc.ResourceL1Gas)
	if err != nil {
		return nil, err
	}
	l2Bytes, err := resourceBounds.L2Gas.Bytes(rpc.ResourceL2Gas)
	if err != nil {
		return nil, err
	}
	l1DataGasBytes, err := resourceBounds.L1DataGas.Bytes(rpc.ResourceL1DataGas)
	if err != nil {
		return nil, err
	}
	l1Bounds := new(felt.Felt).SetBytes(l1Bytes)
	l2Bounds := new(felt.Felt).SetBytes(l2Bytes)
	l1DataGasBounds := new(felt.Felt).SetBytes(l1DataGasBytes)
	return crypto.PoseidonArray(new(felt.Felt).SetUint64(tip), l1Bounds, l2Bounds, l1DataGasBounds), nil
}

func dataAvailabilityMode(feeDAMode, nonceDAMode rpc.DataAvailabilityMode) (uint64, error) {
	const dataAvailabilityModeBits = 32
	fee64, err := feeDAMode.UInt64()
	if err != nil {
		return 0, err
	}
	nonce64, err := nonceDAMode.UInt64()
	if err != nil {
		return 0, err
	}
	return fee64 + nonce64<<dataAvailabilityModeBits, nil
}

// TransactionHashDeclare calculates the transaction hash for declaring a transaction type.
//
// Parameters:
// - tx: The `tx` parameter of type `rpc.DeclareTxnType`
// Can be one of the following types:
//   - `rpc.DeclareTxnV1`
//   - `rpc.DeclareTxnV2`
//   - `rpc.DeclareTxnV3`
//   - `rpc.BroadcastDeclareTxnV3`
//
// Returns:
// - *felt.Felt: the calculated transaction hash as `*felt.Felt` value
// - error: an error, if any
//
// If the `tx` parameter is not one of the supported types, the function returns an error `ErrTxnTypeUnSupported`.
func (account *Account) TransactionHashDeclare(tx rpc.DeclareTxnType) (*felt.Felt, error) {
	switch txn := tx.(type) {
	// Due to inconsistencies in version 0 hash calculation we don't calculate the hash
	case *rpc.DeclareTxnV0, rpc.DeclareTxnV0:
		return nil, ErrTxnVersionUnSupported
	// declare v1, pointer and struct
	case *rpc.DeclareTxnV1:
		return TransactionHashDeclareV1(txn, account.ChainId)
	case rpc.DeclareTxnV1:
		return TransactionHashDeclareV1(&txn, account.ChainId)
	// declare v2, pointer and struct
	case *rpc.DeclareTxnV2:
		return TransactionHashDeclareV2(txn, account.ChainId)
	case rpc.DeclareTxnV2:
		return TransactionHashDeclareV2(&txn, account.ChainId)
	// declare v3, pointer and struct
	case *rpc.DeclareTxnV3:
		return TransactionHashDeclareV3(txn, account.ChainId)
	case rpc.DeclareTxnV3:
		return TransactionHashDeclareV3(&txn, account.ChainId)
	// broadcast declare v3, pointer and struct
	case *rpc.BroadcastDeclareTxnV3:
		return TransactionHashBroadcastDeclareV3(txn, account.ChainId)
	case rpc.BroadcastDeclareTxnV3:
		return TransactionHashBroadcastDeclareV3(&txn, account.ChainId)
	default:
		return nil, fmt.Errorf("%w: got '%T' instead of a valid declare txn type", ErrTxnTypeUnSupported, txn)
	}
}

// TODO: descriptions for all these functions
func TransactionHashDeclareV1(txn *rpc.DeclareTxnV1, chainId *felt.Felt) (*felt.Felt, error) {
	// https://docs.starknet.io/architecture-and-concepts/network-architecture/transactions/#v1_deprecated_hash_calculation_2
	if txn.SenderAddress == nil || txn.Version == "" || txn.ClassHash == nil || txn.MaxFee == nil || txn.Nonce == nil {
		return nil, ErrNotAllParametersSet
	}

	calldataHash := curve.PedersenArray(txn.ClassHash)

	txnVersionFelt, err := new(felt.Felt).SetString(string(txn.Version))
	if err != nil {
		return nil, err
	}
	return hash.CalculateDeprecatedTransactionHashCommon(
		PREFIX_DECLARE,
		txnVersionFelt,
		txn.SenderAddress,
		&felt.Zero,
		calldataHash,
		txn.MaxFee,
		chainId,
		[]*felt.Felt{txn.Nonce},
	), nil
}

func TransactionHashDeclareV2(txn *rpc.DeclareTxnV2, chainId *felt.Felt) (*felt.Felt, error) {
	// https://docs.starknet.io/architecture-and-concepts/network-architecture/transactions/#v2_deprecated_hash_calculation
	if txn.CompiledClassHash == nil || txn.SenderAddress == nil || txn.Version == "" || txn.ClassHash == nil || txn.MaxFee == nil || txn.Nonce == nil {
		return nil, ErrNotAllParametersSet
	}

	calldataHash := curve.PedersenArray(txn.ClassHash)

	txnVersionFelt, err := new(felt.Felt).SetString(string(txn.Version))
	if err != nil {
		return nil, err
	}

	return hash.CalculateDeprecatedTransactionHashCommon(
		PREFIX_DECLARE,
		txnVersionFelt,
		txn.SenderAddress,
		&felt.Zero,
		calldataHash,
		txn.MaxFee,
		chainId,
		[]*felt.Felt{txn.Nonce, txn.CompiledClassHash},
	), nil
}

func TransactionHashDeclareV3(txn *rpc.DeclareTxnV3, chainId *felt.Felt) (*felt.Felt, error) {
	// https://docs.starknet.io/architecture-and-concepts/network-architecture/transactions/#v3_hash_calculation_2
	// https://github.com/starknet-io/SNIPs/blob/main/SNIPS/snip-8.md#protocol-changes
	if txn.Version == "" || txn.ResourceBounds == (rpc.ResourceBoundsMapping{}) || txn.Nonce == nil || txn.SenderAddress == nil || txn.PayMasterData == nil || txn.AccountDeploymentData == nil ||
		txn.ClassHash == nil || txn.CompiledClassHash == nil {
		return nil, ErrNotAllParametersSet
	}

	txnVersionFelt, err := new(felt.Felt).SetString(string(txn.Version))
	if err != nil {
		return nil, err
	}
	DAUint64, err := dataAvailabilityMode(txn.FeeMode, txn.NonceDataMode)
	if err != nil {
		return nil, err
	}
	tipUint64, err := txn.Tip.ToUint64()
	if err != nil {
		return nil, err
	}

	tipAndResourceHash, err := tipAndResourcesHash(tipUint64, txn.ResourceBounds)
	if err != nil {
		return nil, err
	}
	return crypto.PoseidonArray(
		PREFIX_DECLARE,
		txnVersionFelt,
		txn.SenderAddress,
		tipAndResourceHash,
		crypto.PoseidonArray(txn.PayMasterData...),
		chainId,
		txn.Nonce,
		new(felt.Felt).SetUint64(DAUint64),
		crypto.PoseidonArray(txn.AccountDeploymentData...),
		txn.ClassHash,
		txn.CompiledClassHash,
	), nil
}

func TransactionHashBroadcastDeclareV3(txn *rpc.BroadcastDeclareTxnV3, chainId *felt.Felt) (*felt.Felt, error) {
	// https://docs.starknet.io/architecture-and-concepts/network-architecture/transactions/#v3_hash_calculation_2
	// https://github.com/starknet-io/SNIPs/blob/main/SNIPS/snip-8.md#protocol-changes
	if txn.Version == "" || txn.ResourceBounds == (rpc.ResourceBoundsMapping{}) || txn.Nonce == nil || txn.SenderAddress == nil || txn.PayMasterData == nil || txn.AccountDeploymentData == nil ||
		txn.ContractClass == nil || txn.CompiledClassHash == nil {
		return nil, ErrNotAllParametersSet
	}

	txnVersionFelt, err := new(felt.Felt).SetString(string(txn.Version))
	if err != nil {
		return nil, err
	}
	DAUint64, err := dataAvailabilityMode(txn.FeeMode, txn.NonceDataMode)
	if err != nil {
		return nil, err
	}
	tipUint64, err := txn.Tip.ToUint64()
	if err != nil {
		return nil, err
	}

	tipAndResourceHash, err := tipAndResourcesHash(tipUint64, txn.ResourceBounds)
	if err != nil {
		return nil, err
	}
	return crypto.PoseidonArray(
		PREFIX_DECLARE,
		txnVersionFelt,
		txn.SenderAddress,
		tipAndResourceHash,
		crypto.PoseidonArray(txn.PayMasterData...),
		chainId,
		txn.Nonce,
		new(felt.Felt).SetUint64(DAUint64),
		crypto.PoseidonArray(txn.AccountDeploymentData...),
		hash.ClassHash(txn.ContractClass),
		txn.CompiledClassHash,
	), nil
}

// PrecomputeAccountAddress calculates the precomputed address for an account.
// ref: https://docs.starknet.io/architecture-and-concepts/smart-contracts/contract-address/
//
// Parameters:
// - salt: the salt for the address of the deployed contract
// - classHash: the class hash of the contract to be deployed
// - constructorCalldata: the parameters passed to the constructor
// Returns:
// - *felt.Felt: the precomputed address as a *felt.Felt
// - error: an error if any
func PrecomputeAccountAddress(salt *felt.Felt, classHash *felt.Felt, constructorCalldata []*felt.Felt) *felt.Felt {
	return contracts.PrecomputeAddress(&felt.Zero, salt, classHash, constructorCalldata)
}

// WaitForTransactionReceipt waits for the transaction receipt of the given transaction hash to succeed or fail.
//
// Parameters:
//   - ctx: The context
//   - transactionHash: The hash
//   - pollInterval: The time interval to poll the transaction receipt
//
// It returns:
//   - *rpc.TransactionReceipt: the transaction receipt
//   - error: an error
func (account *Account) WaitForTransactionReceipt(ctx context.Context, transactionHash *felt.Felt, pollInterval time.Duration) (*rpc.TransactionReceiptWithBlockInfo, error) {
	t := time.NewTicker(pollInterval)
	for {
		select {
		case <-ctx.Done():
			return nil, rpc.Err(rpc.InternalError, rpc.StringErrData(ctx.Err().Error()))
		case <-t.C:
			receiptWithBlockInfo, err := account.Provider.TransactionReceipt(ctx, transactionHash)
			if err != nil {
				rpcErr := err.(*rpc.RPCError)
				if rpcErr.Code == rpc.ErrHashNotFound.Code && rpcErr.Message == rpc.ErrHashNotFound.Message {
					continue
				} else {
					return nil, err
				}
			}
			return receiptWithBlockInfo, nil
		}
	}
}

// SendTransaction can send Invoke, Declare, and Deploy transactions. It provides a unified way to send different transactions.
//
// Parameters:
//   - ctx: the context.Context object for the transaction.
//   - txn: the Broadcast Transaction to be sent.
//
// Returns:
//   - *rpc.TransactionResponse: the transaction response.
//   - error: an error if any.
func (account *Account) SendTransaction(ctx context.Context, txn rpc.BroadcastTxn) (*rpc.TransactionResponse, error) {
	switch tx := txn.(type) {
	case rpc.BroadcastInvokeTxnType:
		resp, err := account.Provider.AddInvokeTransaction(ctx, tx)
		if err != nil {
			return nil, err
		}
		return &rpc.TransactionResponse{TransactionHash: resp.TransactionHash}, nil
	case rpc.BroadcastDeclareTxnType:
		resp, err := account.Provider.AddDeclareTransaction(ctx, tx)
		if err != nil {
			return nil, err
		}
		return &rpc.TransactionResponse{TransactionHash: resp.TransactionHash, ClassHash: resp.ClassHash}, nil
	case rpc.BroadcastAddDeployTxnType:
		resp, err := account.Provider.AddDeployAccountTransaction(ctx, tx)
		if err != nil {
			return nil, err
		}
		return &rpc.TransactionResponse{TransactionHash: resp.TransactionHash, ContractAddress: resp.ContractAddress}, nil
	default:
		return nil, errors.New("unsupported transaction type")
	}
}

// FmtCalldata generates the formatted calldata for the given function calls and Cairo version.
//
// Parameters:
// - fnCalls: a slice of rpc.FunctionCall representing the function calls.
// Returns:
// - a slice of *felt.Felt representing the formatted calldata.
// - an error if Cairo version is not supported.
func (account *Account) FmtCalldata(fnCalls []rpc.FunctionCall) ([]*felt.Felt, error) {
	switch account.CairoVersion {
	case 0:
		return FmtCallDataCairo0(fnCalls), nil
	case 2:
		return FmtCallDataCairo2(fnCalls), nil
	default:
		return nil, fmt.Errorf("account cairo version '%d' not supported", account.CairoVersion)
	}
}

// FmtCallDataCairo0 generates a slice of *felt.Felt that represents the calldata for the given function calls in Cairo 0 format.
//
// Parameters:
// - fnCalls: a slice of rpc.FunctionCall containing the function calls.
//
// Returns:
// - a slice of *felt.Felt representing the generated calldata.
// https://github.com/project3fusion/StarkSharp/blob/main/StarkSharp/StarkSharp.Rpc/Modules/Transactions/Hash/TransactionHash.cs#L27
func FmtCallDataCairo0(callArray []rpc.FunctionCall) []*felt.Felt {
	var calldata []*felt.Felt
	var calls []*felt.Felt

	calldata = append(calldata, new(felt.Felt).SetUint64(uint64(len(callArray))))

	offset := uint64(0)
	for _, call := range callArray {
		calldata = append(calldata, call.ContractAddress)
		calldata = append(calldata, call.EntryPointSelector)
		calldata = append(calldata, new(felt.Felt).SetUint64(uint64(offset)))
		callDataLen := uint64(len(call.Calldata))
		calldata = append(calldata, new(felt.Felt).SetUint64(callDataLen))
		offset += callDataLen

		calls = append(calls, call.Calldata...)
	}

	calldata = append(calldata, new(felt.Felt).SetUint64(offset))
	calldata = append(calldata, calls...)

	return calldata
}

// FmtCallDataCairo2 generates the calldata for the given function calls for Cairo 2 contracs.
//
// Parameters:
// - fnCalls: a slice of rpc.FunctionCall containing the function calls.
// Returns:
// - a slice of *felt.Felt representing the generated calldata.
// https://github.com/project3fusion/StarkSharp/blob/main/StarkSharp/StarkSharp.Rpc/Modules/Transactions/Hash/TransactionHash.cs#L22
func FmtCallDataCairo2(callArray []rpc.FunctionCall) []*felt.Felt {
	var result []*felt.Felt

	result = append(result, new(felt.Felt).SetUint64(uint64(len(callArray))))

	for _, call := range callArray {
		result = append(result, call.ContractAddress)
		result = append(result, call.EntryPointSelector)

		callDataLen := uint64(len(call.Calldata))
		result = append(result, new(felt.Felt).SetUint64(callDataLen))

		result = append(result, call.Calldata...)
	}

	return result
}

func makeResourceBoundsMapWithZeroValues() rpc.ResourceBoundsMapping {
	return rpc.ResourceBoundsMapping{
		L1Gas: rpc.ResourceBounds{
			MaxAmount:       "0x0",
			MaxPricePerUnit: "0x0",
		},
		L1DataGas: rpc.ResourceBounds{
			MaxAmount:       "0x0",
			MaxPricePerUnit: "0x0",
		},
		L2Gas: rpc.ResourceBounds{
			MaxAmount:       "0x0",
			MaxPricePerUnit: "0x0",
		},
	}
}
