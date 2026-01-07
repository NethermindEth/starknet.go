package account

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/client/rpcerr"
	"github.com/NethermindEth/starknet.go/contracts"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"
)

// BuildAndSendInvokeTxn builds and sends a v3 invoke transaction with the
// given function calls. It automatically calculates the nonce, formats the
// calldata, estimates fees, and signs the transaction with the account's private
// key.
//
// Parameters:
//   - ctx: The context.Context for the request.
//   - functionCalls: A slice of rpc.InvokeFunctionCall representing the function
//     calls for the transaction, allowing either single or multiple function calls
//     in the same transaction.
//   - opts: options for building/estimating the transaction. Pass `nil` to use
//     default values.
//
// Returns:
//   - rpc.AddInvokeTransactionResponse: the response of the submitted transaction.
//   - error: An error if the transaction building fails.
func (account *Account) BuildAndSendInvokeTxn(
	ctx context.Context,
	functionCalls []rpc.InvokeFunctionCall,
	opts *TxnOptions,
) (rpc.AddInvokeTransactionResponse, error) {
	var response rpc.AddInvokeTransactionResponse
	nonce, err := account.Nonce(ctx)
	if err != nil {
		return response, err
	}

	callData, err := account.FmtCalldata(utils.InvokeFuncCallsToFunctionCalls(functionCalls))
	if err != nil {
		return response, err
	}

	if opts == nil {
		opts = new(TxnOptions)
	}
	tip, err := calculateTip(ctx, account.Provider, opts)
	if err != nil {
		return response, err
	}

	// building and signing the txn, as it needs a signature to estimate the fee
	broadcastInvokeTxnV3 := utils.BuildInvokeTxn(
		account.Address,
		nonce,
		callData,
		makeResourceBoundsMapWithZeroValues(),
		&utils.TxnOptions{
			Tip:            tip,
			UseQueryBit:    opts.UseQueryBit,
			UseBlake2sHash: false,
		},
	)

	err = account.SignInvokeTransaction(ctx, broadcastInvokeTxnV3)
	if err != nil {
		return response, err
	}

	// estimate txn fee
	estimateFee, err := account.Provider.EstimateFee(
		ctx,
		[]rpc.BroadcastTxn{broadcastInvokeTxnV3},
		opts.SimulationFlags(),
		opts.BlockID(),
	)
	if err != nil {
		return response, err
	}
	txnFee := estimateFee[0]
	broadcastInvokeTxnV3.ResourceBounds = utils.FeeEstToResBoundsMap(
		txnFee,
		opts.FmtFeeMultiplier(),
	)

	// assuring the signed txn version will be rpc.TransactionV3, since queryBit
	// txn version is only used for estimation/simulation
	broadcastInvokeTxnV3.Version = rpc.TransactionV3

	// signing the txn again with the estimated fee, as the fee value is used in
	// the txn hash calculation
	err = account.SignInvokeTransaction(ctx, broadcastInvokeTxnV3)
	if err != nil {
		return response, err
	}

	response, err = account.Provider.AddInvokeTransaction(ctx, broadcastInvokeTxnV3)
	if err != nil {
		return response, err
	}

	return response, nil
}

// BuildAndSendDeclareTxn builds and sends a v3 declare transaction. It
// automatically calculates the nonce, formats the calldata, estimates fees, and
// signs the transaction with the account's private key.
//
// Parameters:
//   - ctx: The context.Context for the request.
//   - casmClass: The casm class of the contract to be declared
//   - contractClass: The sierra contract class of the contract to be declared
//   - opts: options for building/estimating the transaction. Pass `nil` to use
//     default values.
//
// Returns:
//   - rpc.AddDeclareTransactionResponse: the response of the submitted transaction.
//   - error: An error if the transaction building fails.
func (account *Account) BuildAndSendDeclareTxn(
	ctx context.Context,
	casmClass *contracts.CasmClass,
	contractClass *contracts.ContractClass,
	opts *TxnOptions,
) (rpc.AddDeclareTransactionResponse, error) {
	var response rpc.AddDeclareTransactionResponse
	nonce, err := account.Nonce(ctx)
	if err != nil {
		return response, err
	}

	if opts == nil {
		opts = new(TxnOptions)
	}
	tip, err := calculateTip(ctx, account.Provider, opts)
	if err != nil {
		return response, err
	}

	var useBlake2sHash bool
	if opts.UseBlake2sHash == nil {
		useBlake2sHash, err = shouldUseBlake2sHash(ctx, account.Provider)
		if err != nil {
			return response, fmt.Errorf("failed to check whether to use Blake2s hash: %w", err)
		}
	} else {
		useBlake2sHash = *opts.UseBlake2sHash
	}

	// building and signing the txn, as it needs a signature to estimate the fee
	broadcastDeclareTxnV3, err := utils.BuildDeclareTxn(
		account.Address,
		casmClass,
		contractClass,
		nonce,
		makeResourceBoundsMapWithZeroValues(),
		&utils.TxnOptions{
			Tip:            tip,
			UseQueryBit:    opts.UseQueryBit,
			UseBlake2sHash: useBlake2sHash,
		},
	)
	if err != nil {
		return response, err
	}

	err = account.SignDeclareTransaction(ctx, broadcastDeclareTxnV3)
	if err != nil {
		return response, err
	}

	// estimate txn fee
	estimateFee, err := account.Provider.EstimateFee(
		ctx,
		[]rpc.BroadcastTxn{broadcastDeclareTxnV3},
		opts.SimulationFlags(),
		opts.BlockID(),
	)
	if err != nil {
		return response, err
	}
	txnFee := estimateFee[0]
	broadcastDeclareTxnV3.ResourceBounds = utils.FeeEstToResBoundsMap(
		txnFee,
		opts.FmtFeeMultiplier(),
	)

	// assuring the signed txn version will be rpc.TransactionV3, since queryBit
	// txn version is only used for estimation/simulation
	broadcastDeclareTxnV3.Version = rpc.TransactionV3

	// signing the txn again with the estimated fee, as the fee value is used in
	// the txn hash calculation
	err = account.SignDeclareTransaction(ctx, broadcastDeclareTxnV3)
	if err != nil {
		return response, err
	}

	response, err = account.Provider.AddDeclareTransaction(ctx, broadcastDeclareTxnV3)
	if err != nil {
		return response, err
	}

	return response, nil
}

// BuildAndEstimateDeployAccountTxn builds and signs a v3 deploy account
// transaction, estimates the fee, and computes the address.
//
// This function doesn't send the transaction because the precomputed account
// address requires funding first. This address is calculated deterministically
// and returned by this function, and must be funded with the appropriate amount
// of STRK tokens. Without sufficient funds, the transaction will fail. See the
// 'examples/deployAccount/' for more details on how to do this.
//
// Parameters:
//   - ctx: The context.Context for the request.
//   - salt: the salt for the address of the deployed contract
//   - classHash: the class hash of the contract to be deployed
//   - constructorCalldata: the parameters passed to the constructor
//   - opts: options for building/estimating the transaction. Pass `nil` to use
//     default values.
//
// Returns:
//   - *rpc.BroadcastDeployAccountTxnV3: the transaction to be broadcasted, signed
//     and with the estimated fee based on the multiplier
//   - *felt.Felt: the precomputed account address as a *felt.Felt, it needs to be
//     funded with appropriate amount of tokens
//   - error: an error if any
func (account *Account) BuildAndEstimateDeployAccountTxn(
	ctx context.Context,
	salt *felt.Felt,
	classHash *felt.Felt,
	constructorCalldata []*felt.Felt,
	opts *TxnOptions,
) (*rpc.BroadcastDeployAccountTxnV3, *felt.Felt, error) {
	if opts == nil {
		opts = new(TxnOptions)
	}
	tip, err := calculateTip(ctx, account.Provider, opts)
	if err != nil {
		return nil, nil, err
	}
	// building and signing the txn, as it needs a signature to estimate the fee
	broadcastDepAccTxnV3 := utils.BuildDeployAccountTxn(
		&felt.Zero,
		salt,
		constructorCalldata,
		classHash,
		makeResourceBoundsMapWithZeroValues(),
		&utils.TxnOptions{
			Tip:            tip,
			UseQueryBit:    opts.UseQueryBit,
			UseBlake2sHash: false,
		},
	)

	precomputedAddress := PrecomputeAccountAddress(salt, classHash, constructorCalldata)

	err = account.SignDeployAccountTransaction(ctx, broadcastDepAccTxnV3, precomputedAddress)
	if err != nil {
		return nil, nil, err
	}

	// estimate txn fee
	estimateFee, err := account.Provider.EstimateFee(
		ctx,
		[]rpc.BroadcastTxn{broadcastDepAccTxnV3},
		opts.SimulationFlags(),
		opts.BlockID(),
	)
	if err != nil {
		return nil, nil, err
	}
	txnFee := estimateFee[0]
	broadcastDepAccTxnV3.ResourceBounds = utils.FeeEstToResBoundsMap(
		txnFee,
		opts.FmtFeeMultiplier(),
	)

	// assuring the signed txn version will be rpc.TransactionV3, since queryBit
	// txn version is only used for estimation/simulation
	broadcastDepAccTxnV3.Version = rpc.TransactionV3

	// signing the txn again with the estimated fee, as the fee value is used in
	// the txn hash calculation
	err = account.SignDeployAccountTransaction(ctx, broadcastDepAccTxnV3, precomputedAddress)
	if err != nil {
		return nil, nil, err
	}

	return broadcastDepAccTxnV3, precomputedAddress, nil
}

// calculateTip returns the tip to be used in the transaction. If a custom tip is
// provided, it returns it. Otherwise, it estimates the tip using the provider
// based on the tip multiplier.
func calculateTip(
	ctx context.Context,
	provider rpc.RPCProvider,
	opts *TxnOptions,
) (rpc.U64, error) {
	if opts.CustomTip != "" {
		return opts.CustomTip, nil
	}

	tip, err := rpc.EstimateTip(ctx, provider, opts.FmtTipMultiplier())
	if err != nil {
		return "", fmt.Errorf("failed to estimate tip: %w", err)
	}

	return tip, nil
}

// A helper to deploy a contract from an existing class using UDC.
//
// Parameters:
//   - ctx: The context.Context for the request.
//   - classHash: The class hash of the contract to be deployed.
//   - constructorCalldata: The parameters passed to the constructor. Pass `nil` if
//     the constructor has no arguments.
//   - txnOpts: The options for building/estimating the transaction. Pass `nil` to
//     use default values.
//   - udcOpts: The options for building the UDC calldata. Pass `nil` to use
//     default values.
//
// Returns:
//   - *rpc.AddInvokeTransactionResponse: the response of the submitted UDC
//     transaction.
//   - *felt.Felt: the salt used for the UDC deployment (either the provided one or
//     the random one)
//   - error: An error if any.
func (account *Account) DeployContractWithUDC(
	ctx context.Context,
	classHash *felt.Felt,
	constructorCalldata []*felt.Felt,
	txnOpts *TxnOptions,
	udcOpts *UDCOptions,
) (rpc.AddInvokeTransactionResponse, *felt.Felt, error) {
	var response rpc.AddInvokeTransactionResponse
	udcCallData, salt, err := utils.BuildUDCCalldata(classHash, constructorCalldata, udcOpts)
	if err != nil {
		return response, nil, err
	}

	response, err = account.BuildAndSendInvokeTxn(
		context.Background(),
		[]rpc.InvokeFunctionCall{udcCallData},
		txnOpts,
	)
	if err != nil {
		return response, nil, err
	}

	return response, salt, nil
}

// SendTransaction can send Invoke, Declare, and Deploy transactions. It
// provides a unified way to send different transactions. It can only send v3
// transactions.
//
// Parameters:
//   - ctx: the context.Context object for the transaction.
//   - txn: the Broadcast V3 Transaction to be sent
//
// Returns:
//   - rpc.TransactionResponse: the transaction response for each TransactionResponse
//   - error: an error if any
//
//nolint:exhaustruct // Setting only the correct fields
func (account *Account) SendTransaction(
	ctx context.Context,
	txn rpc.BroadcastTxn,
) (rpc.TransactionResponse, error) {
	var response rpc.TransactionResponse
	switch tx := txn.(type) {
	// broadcast invoke v3, pointer and struct
	case *rpc.BroadcastInvokeTxnV3:
		resp, err := account.Provider.AddInvokeTransaction(ctx, tx)
		if err != nil {
			return response, err
		}

		return rpc.TransactionResponse{Hash: resp.Hash}, nil
	case rpc.BroadcastInvokeTxnV3:
		resp, err := account.Provider.AddInvokeTransaction(ctx, &tx)
		if err != nil {
			return response, err
		}

		return rpc.TransactionResponse{Hash: resp.Hash}, nil
	// broadcast declare v3, pointer and struct
	case *rpc.BroadcastDeclareTxnV3:
		resp, err := account.Provider.AddDeclareTransaction(ctx, tx)
		if err != nil {
			return response, err
		}

		return rpc.TransactionResponse{
			Hash:      resp.Hash,
			ClassHash: resp.ClassHash,
		}, nil
	case rpc.BroadcastDeclareTxnV3:
		resp, err := account.Provider.AddDeclareTransaction(ctx, &tx)
		if err != nil {
			return response, err
		}

		return rpc.TransactionResponse{
			Hash:      resp.Hash,
			ClassHash: resp.ClassHash,
		}, nil
	// broadcast deploy account v3, pointer and struct
	case *rpc.BroadcastDeployAccountTxnV3:
		resp, err := account.Provider.AddDeployAccountTransaction(ctx, tx)
		if err != nil {
			return response, err
		}

		return rpc.TransactionResponse{
			Hash:            resp.Hash,
			ContractAddress: resp.ContractAddress,
		}, nil
	case rpc.BroadcastDeployAccountTxnV3:
		resp, err := account.Provider.AddDeployAccountTransaction(ctx, &tx)
		if err != nil {
			return response, err
		}

		return rpc.TransactionResponse{
			Hash:            resp.Hash,
			ContractAddress: resp.ContractAddress,
		}, nil
	default:
		return response, fmt.Errorf(
			"unsupported transaction type: should be a v3 transaction, instead got %T",
			tx,
		)
	}
}

// WaitForTransactionReceipt waits for the transaction receipt of the given
// transaction hash to succeed or fail.
//
// Parameters:
//   - ctx: The context
//   - transactionHash: The hash
//   - pollInterval: The time interval to poll the transaction receipt
//
// Returns:
//   - *rpc.TransactionReceiptWithBlockInfo: the transaction receipt
//   - error: an error if any
func (account *Account) WaitForTransactionReceipt(
	ctx context.Context,
	transactionHash *felt.Felt,
	pollInterval time.Duration,
) (*rpc.TransactionReceiptWithBlockInfo, error) {
	t := time.NewTicker(pollInterval)
	for {
		select {
		case <-ctx.Done():
			return nil, rpcerr.Err(rpcerr.InternalError, rpc.StringErrData(ctx.Err().Error()))
		case <-t.C:
			receiptWithBlockInfo, err := account.Provider.TransactionReceipt(ctx, transactionHash)
			if err != nil {
				rpcErr := err.(*rpc.RPCError)
				if rpcErr.Code == rpc.ErrHashNotFound.Code &&
					rpcErr.Message == rpc.ErrHashNotFound.Message {
					continue
				} else {
					return nil, err
				}
			}

			return receiptWithBlockInfo, nil
		}
	}
}

// shouldUseBlake2sHash determines whether to use the Blake2s hash function for the
// compiled class hash.
// Starknet v0.14.1 upgrade will deprecate the Poseidon hash function for the compiled class hash.
//
// Parameters:
//   - ctx: The context
//   - provider: The provider
//
// Returns:
//   - bool: whether to use the Blake2s hash function for the compiled class hash
//   - error: an error if any
func shouldUseBlake2sHash(ctx context.Context, provider rpc.RPCProvider) (bool, error) {
	block, err := provider.BlockWithTxHashes(ctx, rpc.WithBlockTag(rpc.BlockTagLatest))
	if err != nil {
		return false, fmt.Errorf("failed to get block with tx hashes: %w", err)
	}

	var starknetVersion string
	switch {
	case block.Block != nil:
		starknetVersion = block.Block.StarknetVersion
	case block.PreConfirmed != nil:
		starknetVersion = block.PreConfirmed.StarknetVersion
	default:
		return false, fmt.Errorf("unexpected block output: %+v", block)
	}

	upgradeVersion := semver.MustParse("0.14.1")

	currentVersion, err := semver.NewVersion(starknetVersion)
	if err != nil {
		return false, fmt.Errorf("failed to parse block's starknet version: %w", err)
	}

	return currentVersion.Compare(upgradeVersion) >= 0, nil
}
