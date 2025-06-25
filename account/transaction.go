package account

import (
	"context"
	"fmt"
	"time"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/contracts"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"
)

// BuildAndSendInvokeTxn builds and sends a v3 invoke transaction with the given function calls.
// It automatically calculates the nonce, formats the calldata, estimates fees, and signs the transaction with the account's private key.
//
// Parameters:
//   - ctx: The context.Context for the request.
//   - functionCalls: A slice of rpc.InvokeFunctionCall representing the function calls for the transaction, allowing either single or
//     multiple function calls in the same transaction.
//   - opts: options for building/estimating the transaction. See more info in the TxnOptions type description.
//
// Returns:
//   - *rpc.AddInvokeTransactionResponse: the response of the submitted transaction.
//   - error: An error if the transaction building fails.
func (account *Account) BuildAndSendInvokeTxn(
	ctx context.Context,
	functionCalls []rpc.InvokeFunctionCall,
	opts *TxnOptions,
) (*rpc.AddInvokeTransactionResponse, error) {
	nonce, err := account.Nonce(ctx)
	if err != nil {
		return nil, err
	}

	callData, err := account.FmtCalldata(utils.InvokeFuncCallsToFunctionCalls(functionCalls))
	if err != nil {
		return nil, err
	}

	if opts == nil {
		opts = new(TxnOptions)
	}
	fmtTipAndMultiplier(opts)

	// building and signing the txn, as it needs a signature to estimate the fee
	broadcastInvokeTxnV3 := utils.BuildInvokeTxn(
		account.Address,
		nonce,
		callData,
		makeResourceBoundsMapWithZeroValues(),
		&utils.TxnOptions{
			Tip:         opts.Tip,
			UseQueryBit: opts.UseQueryBit,
		},
	)

	err = account.SignInvokeTransaction(ctx, broadcastInvokeTxnV3)
	if err != nil {
		return nil, err
	}

	// estimate txn fee
	estimateFee, err := account.Provider.EstimateFee(
		ctx,
		[]rpc.BroadcastTxn{broadcastInvokeTxnV3},
		opts.SimulationFlags(),
		opts.BlockID(),
	)
	if err != nil {
		return nil, err
	}
	txnFee := estimateFee[0]
	broadcastInvokeTxnV3.ResourceBounds = utils.FeeEstToResBoundsMap(txnFee, opts.Multiplier)

	// assuring the signed txn version will be rpc.TransactionV3, since queryBit txn version is only used for estimation/simulation
	broadcastInvokeTxnV3.Version = rpc.TransactionV3

	// signing the txn again with the estimated fee, as the fee value is used in the txn hash calculation
	err = account.SignInvokeTransaction(ctx, broadcastInvokeTxnV3)
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
//   - opts: options for building/estimating the transaction. See more info in the TxnOptions type description.
//
// Returns:
//   - *rpc.AddDeclareTransactionResponse: the response of the submitted transaction.
//   - error: An error if the transaction building fails.
func (account *Account) BuildAndSendDeclareTxn(
	ctx context.Context,
	casmClass *contracts.CasmClass,
	contractClass *contracts.ContractClass,
	opts *TxnOptions,
) (*rpc.AddDeclareTransactionResponse, error) {
	nonce, err := account.Nonce(ctx)
	if err != nil {
		return nil, err
	}

	if opts == nil {
		opts = new(TxnOptions)
	}
	fmtTipAndMultiplier(opts)

	// building and signing the txn, as it needs a signature to estimate the fee
	broadcastDeclareTxnV3, err := utils.BuildDeclareTxn(
		account.Address,
		casmClass,
		contractClass,
		nonce,
		makeResourceBoundsMapWithZeroValues(),
		&utils.TxnOptions{
			Tip:         opts.Tip,
			UseQueryBit: opts.UseQueryBit,
		},
	)
	if err != nil {
		return nil, err
	}

	err = account.SignDeclareTransaction(ctx, broadcastDeclareTxnV3)
	if err != nil {
		return nil, err
	}

	// estimate txn fee
	estimateFee, err := account.Provider.EstimateFee(
		ctx,
		[]rpc.BroadcastTxn{broadcastDeclareTxnV3},
		opts.SimulationFlags(),
		opts.BlockID(),
	)
	if err != nil {
		return nil, err
	}
	txnFee := estimateFee[0]
	broadcastDeclareTxnV3.ResourceBounds = utils.FeeEstToResBoundsMap(txnFee, opts.Multiplier)

	// assuring the signed txn version will be rpc.TransactionV3, since queryBit txn version is only used for estimation/simulation
	broadcastDeclareTxnV3.Version = rpc.TransactionV3

	// signing the txn again with the estimated fee, as the fee value is used in the txn hash calculation
	err = account.SignDeclareTransaction(ctx, broadcastDeclareTxnV3)
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
//   - opts: options for building/estimating the transaction. See more info in the TxnOptions type description.
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
	opts *TxnOptions,
) (*rpc.BroadcastDeployAccountTxnV3, *felt.Felt, error) {
	if opts == nil {
		opts = new(TxnOptions)
	}
	fmtTipAndMultiplier(opts)

	// building and signing the txn, as it needs a signature to estimate the fee
	broadcastDepAccTxnV3 := utils.BuildDeployAccountTxn(
		&felt.Zero,
		salt,
		constructorCalldata,
		classHash,
		makeResourceBoundsMapWithZeroValues(),
		&utils.TxnOptions{
			Tip:         opts.Tip,
			UseQueryBit: opts.UseQueryBit,
		},
	)

	precomputedAddress := PrecomputeAccountAddress(salt, classHash, constructorCalldata)

	err := account.SignDeployAccountTransaction(ctx, broadcastDepAccTxnV3, precomputedAddress)
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
	broadcastDepAccTxnV3.ResourceBounds = utils.FeeEstToResBoundsMap(txnFee, opts.Multiplier)

	// assuring the signed txn version will be rpc.TransactionV3, since queryBit txn version is only used for estimation/simulation
	broadcastDepAccTxnV3.Version = rpc.TransactionV3

	// signing the txn again with the estimated fee, as the fee value is used in the txn hash calculation
	err = account.SignDeployAccountTransaction(ctx, broadcastDepAccTxnV3, precomputedAddress)
	if err != nil {
		return nil, nil, err
	}

	return broadcastDepAccTxnV3, precomputedAddress, nil
}

// DeployContractWithUDC deploys a contract using UDC.
//
// Parameters:
//   - ctx: the context
//   - classHash: the class hash of the contract to be deployed
//   - salt: the salt for the address of the deployed contract
//   - constructorCalldata: the parameters passed to the constructor
//   - udcAddress: the address of the UDC contract. If nil, the default address will be used.
//
// It returns:
//   - *rpc.AddInvokeTransactionResponse: the response from the provider
//   - error: an error if any
//
// TODO: improve docs
func (account *Account) DeployContractWithUDC(
	ctx context.Context,
	classHash *felt.Felt,
	constructorCalldata []*felt.Felt,
	txnOpts *TxnOptions,
	udcOpts *utils.UDCOptions,
) (*rpc.AddInvokeTransactionResponse, error) {
	// TODO: implement this

	// fromZeroFelt := new(felt.Felt).SetUint64(1)

	// calldataLen := new(felt.Felt).SetUint64(uint64(len(constructorCalldata)))
	// udcCallData := append([]*felt.Felt{classHash, salt, fromZeroFelt, calldataLen}, constructorCalldata...)

	// var finalUdcAddress *felt.Felt
	// if udcAddress != nil {
	// 	finalUdcAddress = udcAddress
	// } else {
	// 	var err error
	// 	// Default address is same for Mainnet and Sepolia testnet.
	// 	// https://docs.openzeppelin.com/contracts-cairo/0.14.0/udc
	// 	finalUdcAddress, err = new(felt.Felt).SetString("0x04a64cd09a853868621d94cae9952b106f2c36a3f81260f85de6696c6b050221")
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }

	fnCall := rpc.InvokeFunctionCall{
		// ContractAddress: finalUdcAddress,
		// FunctionName:    "deploy_contract",
		// CallData:        udcCallData,
	}

	// Setting multiplier to 1.5 for now, I think ideally the user should be able to set it.
	return account.BuildAndSendInvokeTxn(context.Background(), []rpc.InvokeFunctionCall{fnCall}, nil)
}

// SendTransaction can send Invoke, Declare, and Deploy transactions. It provides a unified way to send different transactions.
// It can only send v3 transactions.
//
// Parameters:
//   - ctx: the context.Context object for the transaction.
//   - txn: the Broadcast V3 Transaction to be sent.
//
// Returns:
//   - *rpc.TransactionResponse: the transaction response.
//   - error: an error if any.
func (account *Account) SendTransaction(ctx context.Context, txn rpc.BroadcastTxn) (*rpc.TransactionResponse, error) {
	switch tx := txn.(type) {
	// broadcast invoke v3, pointer and struct
	case *rpc.BroadcastInvokeTxnV3:
		resp, err := account.Provider.AddInvokeTransaction(ctx, tx)
		if err != nil {
			return nil, err
		}

		return &rpc.TransactionResponse{Hash: resp.Hash}, nil //nolint:exhaustruct
	case rpc.BroadcastInvokeTxnV3:
		resp, err := account.Provider.AddInvokeTransaction(ctx, &tx)
		if err != nil {
			return nil, err
		}

		return &rpc.TransactionResponse{Hash: resp.Hash}, nil //nolint:exhaustruct
	// broadcast declare v3, pointer and struct
	case *rpc.BroadcastDeclareTxnV3:
		resp, err := account.Provider.AddDeclareTransaction(ctx, tx)
		if err != nil {
			return nil, err
		}

		return &rpc.TransactionResponse{Hash: resp.Hash, ClassHash: resp.ClassHash}, nil //nolint:exhaustruct
	case rpc.BroadcastDeclareTxnV3:
		resp, err := account.Provider.AddDeclareTransaction(ctx, &tx)
		if err != nil {
			return nil, err
		}

		return &rpc.TransactionResponse{Hash: resp.Hash, ClassHash: resp.ClassHash}, nil //nolint:exhaustruct
	// broadcast deploy account v3, pointer and struct
	case *rpc.BroadcastDeployAccountTxnV3:
		resp, err := account.Provider.AddDeployAccountTransaction(ctx, tx)
		if err != nil {
			return nil, err
		}

		return &rpc.TransactionResponse{Hash: resp.Hash, ContractAddress: resp.ContractAddress}, nil //nolint:exhaustruct
	case rpc.BroadcastDeployAccountTxnV3:
		resp, err := account.Provider.AddDeployAccountTransaction(ctx, &tx)
		if err != nil {
			return nil, err
		}

		return &rpc.TransactionResponse{Hash: resp.Hash, ContractAddress: resp.ContractAddress}, nil //nolint:exhaustruct
	default:
		return nil, fmt.Errorf("unsupported transaction type: should be a v3 transaction, instead got %T", tx)
	}
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
func (account *Account) WaitForTransactionReceipt(
	ctx context.Context,
	transactionHash *felt.Felt,
	pollInterval time.Duration,
) (*rpc.TransactionReceiptWithBlockInfo, error) {
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
