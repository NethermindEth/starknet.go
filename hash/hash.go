package hash

import (
	"errors"
	"slices"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/contracts"
	"github.com/NethermindEth/starknet.go/curve"
	"github.com/NethermindEth/starknet.go/rpc"
)

var (
	prefixInvoke        = new(felt.Felt).SetBytes([]byte("invoke"))
	prefixDeclare       = new(felt.Felt).SetBytes([]byte("declare"))
	prefixDeployAccount = new(felt.Felt).SetBytes([]byte("deploy_account"))
)

var (
	ErrNotAllParametersSet = errors.New("not all necessary parameters have been set")
	ErrFeltToBigInt        = errors.New("felt to BigInt error")
)

// CalculateDeprecatedTransactionHashCommon calculates the transaction hash
// common to be used in the StarkNet network - a unique identifier of the transaction.
// [specification]: https://github.com/starkware-libs/cairo-lang/blob/8276ac35830148a397e1143389f23253c8b80e93/src/starkware/starknet/core/os/transaction_hash/deprecated_transaction_hash.py#L29
//
// Parameters:
//   - txHashPrefix: The prefix of the transaction hash
//   - version: The version of the transaction
//   - contractAddress: The address of the contract
//   - entryPointSelector: The selector of the entry point
//   - calldata: The data of the transaction
//   - maxFee: The maximum fee for the transaction
//   - chainID: The ID of the blockchain
//   - additionalData: Additional data to be included in the hash
//
// Returns:
//   - *felt.Felt: the calculated transaction hash
//
//nolint:lll // The link would be unclickable if we break the line.
func CalculateDeprecatedTransactionHashCommon(
	txHashPrefix *felt.Felt,
	version *felt.Felt,
	contractAddress *felt.Felt,
	entryPointSelector *felt.Felt,
	calldata *felt.Felt,
	maxFee *felt.Felt,
	chainID *felt.Felt,
	additionalData []*felt.Felt,
) *felt.Felt {
	dataToHash := []*felt.Felt{
		txHashPrefix,
		version,
		contractAddress,
		entryPointSelector,
		calldata,
		maxFee,
		chainID,
	}
	dataToHash = append(dataToHash, additionalData...)

	return curve.PedersenArray(dataToHash...)
}

// ClassHash calculates the hash of a contract class.
//
// It takes a contract class as input and calculates the hash by combining
// various elements of the class. The hash is calculated using the
// PoseidonArray function from the Curve package. The elements used in the
// hash calculation include the contract class version, constructor entry point,
// external entry point, L1 handler entry point, ABI, and Sierra program.
// The ABI is converted to bytes and then hashed using the StarknetKeccak
// function from the Curve package.
// Finally, the ContractClassVersionHash, ExternalHash, L1HandleHash,
// ConstructorHash, ABIHash, and
// SierraProgamHash are combined using the PoseidonArray function from the
// Curve package.
//
// Parameters:
//   - contract: A contract class object of type contracts.ContractClass.
//
// Returns:
//   - *felt.Felt: a pointer to a felt.Felt object that represents the calculated hash.
//   - error: an error object if there was an error during the hash calculation.
func ClassHash(contract *contracts.ContractClass) *felt.Felt {
	// https://docs.starknet.io/architecture-and-concepts/smart-contracts/class-hash/

	Version := "CONTRACT_CLASS_V" + contract.ContractClassVersion
	ContractClassVersionHash := new(felt.Felt).SetBytes([]byte(Version))
	ConstructorHash := hashEntryPointByType(contract.EntryPointsByType.Constructor)
	ExternalHash := hashEntryPointByType(contract.EntryPointsByType.External)
	L1HandleHash := hashEntryPointByType(contract.EntryPointsByType.L1Handler)
	SierraProgamHash := curve.PoseidonArray(contract.SierraProgram...)
	ABIHash := curve.StarknetKeccak([]byte(contract.ABI))

	// https://docs.starknet.io/architecture-and-concepts/smart-contracts/class-hash/#computing_the_cairo_1_class_hash
	//nolint:lll // The link would be unclickable if we break the line.
	return curve.PoseidonArray(
		ContractClassVersionHash,
		ExternalHash,
		L1HandleHash,
		ConstructorHash,
		ABIHash,
		SierraProgamHash,
	)
}

// hashEntryPointByType calculates the hash of an entry point by type.
//
// Parameters:
//   - entryPoint: A slice of contracts.SierraEntryPoint objects
//
// Returns:
//   - *felt.Felt: a pointer to a felt.Felt object that represents the calculated
//     hash.
func hashEntryPointByType(entryPoint []contracts.SierraEntryPoint) *felt.Felt {
	flattened := make([]*felt.Felt, 0, len(entryPoint))
	for _, elt := range entryPoint {
		flattened = append(
			flattened,
			elt.Selector,
			new(felt.Felt).SetUint64(uint64(elt.FunctionIdx)),
		)
	}

	return curve.PoseidonArray(flattened...)
}

// getByteCodeSegmentHasher calculates hasher function for byte code array from
// casm file. This code is adaptation of:
// https://github.com/starkware-libs/cairo-lang/blob/efa9648f57568aad8f8a13fbf027d2de7c63c2c0/src/starkware/starknet/core/os/contract_class/compiled_class_hash.py
//
// Parameters:
//   - bytecode: Array of compiled bytecode values from casm file
//   - bytecodeSegmentLengths: Nested datastructure of bytecode_segment_lengths
//     values from casm file
//   - visitedPcs: array pointer for tracking which bytecode bits were already
//     processed, needed for recursive processing
//   - bytecodeOffset: pointer at current offset in bytecode array, needed for
//     recursive processing organisation
//
// Returns:
//   - hasherFunc: closure that calculates hash for given bytecode array, or nil
//     in case of error
//   - uint64: size of the current processed bytecode array, or nil in case of error
//   - error: error if any happened or nil if everything fine
//
//nolint:lll // The link would be unclickable if we break the line.
func getByteCodeSegmentHasher(
	bytecode []*felt.Felt,
	bytecodeSegmentLengths contracts.NestedUints,
	visitedPcs *[]uint64,
	bytecodeOffset uint64,
	hashFunc func(...*felt.Felt) *felt.Felt,
) (hash *felt.Felt, size uint64, err error) {
	if !bytecodeSegmentLengths.IsArray {
		segmentValue := *bytecodeSegmentLengths.Value
		segmentEnd := bytecodeOffset + segmentValue

		for {
			visitedPcsData := *visitedPcs

			if len(visitedPcsData) == 0 {
				break
			}

			lastVisitedPcs := visitedPcsData[len(visitedPcsData)-1]

			if (bytecodeOffset > lastVisitedPcs) || (lastVisitedPcs >= segmentEnd) {
				break
			}

			*visitedPcs = visitedPcsData[:len(visitedPcsData)-1]
		}

		bytecodePart := bytecode[bytecodeOffset:segmentEnd]

		return hashFunc(bytecodePart...), segmentValue, nil
	}

	type bytecodeSegment struct {
		Value *felt.Felt
		Size  uint64
	}

	segments := []bytecodeSegment{}
	totalLen := uint64(0)

	for _, item := range bytecodeSegmentLengths.Values {
		visitedPcsData := *visitedPcs
		var visitedPcBefore *uint64

		if len(visitedPcsData) > 0 {
			visitedPcBefore = &visitedPcsData[len(visitedPcsData)-1]
		}

		segmentHash, segmentLen, err := getByteCodeSegmentHasher(
			bytecode,
			item,
			visitedPcs,
			bytecodeOffset,
			hashFunc,
		)
		if err != nil {
			return nil, 0, err
		}

		var visitedPcAfter *uint64
		if len(visitedPcsData) > 0 {
			visitedPcAfter = &visitedPcsData[len(visitedPcsData)-1]
		}

		isUsed := visitedPcAfter != visitedPcBefore

		if isUsed && *visitedPcBefore != bytecodeOffset {
			return nil, 0, errors.New(
				"invalid segment structure: PC {visited_pc_before} was visited, " +
					"but the beginning of the segment ({bytecode_offset}) was not",
			)
		}

		segments = append(segments, bytecodeSegment{
			Value: segmentHash,
			Size:  segmentLen,
		})
		bytecodeOffset += segmentLen
		totalLen += segmentLen
	}

	components := make([]*felt.Felt, len(segments)*2)

	for i, val := range segments {
		components[i*2] = felt.NewFromUint64[felt.Felt](val.Size)
		components[i*2+1] = val.Value
	}

	return new(felt.Felt).Add(
		felt.NewFromUint64[felt.Felt](1),
		hashFunc(components...),
	), totalLen, nil
}

// getByteCodeSegmentHasher calculates hash for byte code array from casm file
//
// Parameters:
//   - bytecode: Array of compiled bytecode values from casm file
//   - bytecodeSegmentLengths: Nested datastructure of bytecode_segment_lengths
//     values from casm file
//
// Returns:
//   - *felt.Felt: Hash value
//   - error: Error message
func hashCasmClassByteCode(
	bytecode []*felt.Felt,
	bytecodeSegmentLengths contracts.NestedUints,
	hashFunc func(...*felt.Felt) *felt.Felt,
) (*felt.Felt, error) {
	visited := make([]uint64, len(bytecode))

	for i := range bytecode {
		visited[i] = uint64(i) //nolint:gosec // Never overflows
	}

	slices.Reverse(visited)

	hash, _, err := getByteCodeSegmentHasher(
		bytecode,
		bytecodeSegmentLengths,
		&visited,
		uint64(0),
		hashFunc,
	)
	if err != nil {
		return nil, err
	}

	return hash, nil
}

// compiledClassHash calculates the hash of a compiled class in the Casm format
// using the provided hash function.
//
// Parameters:
//   - casmClass: A `contracts.CasmClass` object
//
// Returns:
//   - *felt.Felt: a pointer to a felt.Felt object that represents the calculated hash.
func compiledClassHash(
	casmClass *contracts.CasmClass,
	hashFunc func(...*felt.Felt) *felt.Felt,
) (*felt.Felt, error) {
	ContractClassVersionHash := new(felt.Felt).SetBytes([]byte("COMPILED_CLASS_V1"))
	ExternalHash := hashCasmEntryPoints(casmClass.EntryPointsByType.External, hashFunc)
	L1HandleHash := hashCasmEntryPoints(casmClass.EntryPointsByType.L1Handler, hashFunc)
	ConstructorHash := hashCasmEntryPoints(casmClass.EntryPointsByType.Constructor, hashFunc)

	var ByteCodeHasH *felt.Felt
	var err error

	if casmClass.BytecodeSegmentLengths != nil {
		ByteCodeHasH, err = hashCasmClassByteCode(
			casmClass.ByteCode,
			*casmClass.BytecodeSegmentLengths,
			hashFunc,
		)
		if err != nil {
			return nil, err
		}
	} else {
		ByteCodeHasH = hashFunc(casmClass.ByteCode...)
	}

	//nolint:lll // The link would be unclickable if we break the line.
	// https://github.com/software-mansion/starknet.py/blob/39af414389984efbc6edc48b0fe1f914ea5b9a77/starknet_py/hash/casm_class_hash.py#L18
	return hashFunc(
		ContractClassVersionHash,
		ExternalHash,
		L1HandleHash,
		ConstructorHash,
		ByteCodeHasH,
	), nil
}

// CompiledClassHash calculates the hash of a compiled class in the Casm format
// using the Poseidon hash function.
// This function will be deprecated in Starknet v0.14.1 onwards.
//
// Parameters:
//   - casmClass: A `contracts.CasmClass` object
//
// Returns:
//   - *felt.Felt: a pointer to a felt.Felt object that represents the calculated hash.
func CompiledClassHash(casmClass *contracts.CasmClass) (*felt.Felt, error) {
	return compiledClassHash(casmClass, curve.PoseidonArray)
}

// CompiledClassHashV2 calculates the hash of a compiled class in the Casm format
// using the Blake2s hash function. This is correct hash function to calculate the
// compiled class hash for Starknet v0.14.1 onwards.
//
// Parameters:
//   - casmClass: A `contracts.CasmClass` object
//
// Returns:
//   - *felt.Felt: a pointer to a felt.Felt object that represents the calculated hash.
func CompiledClassHashV2(casmClass *contracts.CasmClass) (*felt.Felt, error) {
	return compiledClassHash(casmClass, curve.Blake2sArray)
}

// hashCasmEntryPoints calculates the hash of a CasmClassEntryPoint array
// using the provided hash function.
//
// Parameters:
//   - entryPoint: An array of CasmClassEntryPoint objects
//   - hashFunc: A function that takes a variadic number of pointers to felt.Felt
//     and returns a pointer to a felt.Felt
//
// Returns:
//   - *felt.Felt: a pointer to a Felt type
func hashCasmEntryPoints(
	entryPoint []contracts.CasmEntryPoint,
	hashFunc func(...*felt.Felt) *felt.Felt,
) *felt.Felt {
	flattened := make([]*felt.Felt, 0, len(entryPoint))
	for _, elt := range entryPoint {
		builtInFlat := []*felt.Felt{}
		for _, builtIn := range elt.Builtins {
			builtInFlat = append(builtInFlat, new(felt.Felt).SetBytes([]byte(builtIn)))
		}
		builtInHash := hashFunc(builtInFlat...)
		flattened = append(
			flattened,
			elt.Selector,
			new(felt.Felt).SetUint64(uint64(elt.Offset)),
			builtInHash,
		)
	}

	return hashFunc(flattened...)
}

// TransactionHashInvokeV0 calculates the transaction hash for a invoke V0 transaction.
//
// Parameters:
//   - txn: The invoke V0 transaction to calculate the hash for
//   - chainID: The chain ID as a *felt.Felt
//
// Returns:
//   - *felt.Felt: the calculated transaction hash
//   - error: an error if any
func TransactionHashInvokeV0(txn *rpc.InvokeTxnV0, chainID *felt.Felt) (*felt.Felt, error) {
	//nolint:lll // The link would be unclickable if we break the line.
	// https://docs.starknet.io/architecture-and-concepts/network-architecture/transactions/#v0_deprecated_hash_calculation
	if txn.Version == "" || len(txn.Calldata) == 0 || txn.MaxFee == nil ||
		txn.EntryPointSelector == nil {
		return nil, ErrNotAllParametersSet
	}

	calldataHash := curve.PedersenArray(txn.Calldata...)
	txnVersionFelt, err := new(felt.Felt).SetString(string(txn.Version))
	if err != nil {
		return nil, err
	}

	return CalculateDeprecatedTransactionHashCommon(
		prefixInvoke,
		txnVersionFelt,
		txn.ContractAddress,
		txn.EntryPointSelector,
		calldataHash,
		txn.MaxFee,
		chainID,
		[]*felt.Felt{},
	), nil
}

// TransactionHashInvokeV1 calculates the transaction hash for a invoke V1 transaction.
//
// Parameters:
//   - txn: The invoke V1 transaction to calculate the hash for
//   - chainID: The chain ID as a *felt.Felt
//
// Returns:
//   - *felt.Felt: the calculated transaction hash
//   - error: an error if any
func TransactionHashInvokeV1(txn *rpc.InvokeTxnV1, chainID *felt.Felt) (*felt.Felt, error) {
	//nolint:lll // The link would be unclickable if we break the line.
	// https://docs.starknet.io/architecture-and-concepts/network-architecture/transactions/#v1_deprecated_hash_calculation
	if txn.Version == "" || len(txn.Calldata) == 0 || txn.Nonce == nil || txn.MaxFee == nil ||
		txn.SenderAddress == nil {
		return nil, ErrNotAllParametersSet
	}

	calldataHash := curve.PedersenArray(txn.Calldata...)
	txnVersionFelt, err := new(felt.Felt).SetString(string(txn.Version))
	if err != nil {
		return nil, err
	}

	return CalculateDeprecatedTransactionHashCommon(
		prefixInvoke,
		txnVersionFelt,
		txn.SenderAddress,
		&felt.Zero,
		calldataHash,
		txn.MaxFee,
		chainID,
		[]*felt.Felt{txn.Nonce},
	), nil
}

// TransactionHashInvokeV3 calculates the transaction hash for a invoke V3 transaction.
//
// Parameters:
//   - txn: The invoke V3 transaction to calculate the hash for
//   - chainID: The chain ID as a *felt.Felt
//
// Returns:
//   - *felt.Felt: the calculated transaction hash
//   - error: an error if any
func TransactionHashInvokeV3(txn *rpc.InvokeTxnV3, chainID *felt.Felt) (*felt.Felt, error) {
	//nolint:lll // The links would be unclickable if we break the line.
	// https://github.com/starknet-io/SNIPs/blob/main/SNIPS/snip-8.md#protocol-changes
	// https://docs.starknet.io/architecture-and-concepts/network-architecture/transactions/#v3_hash_calculation
	if txn.Version == "" || txn.ResourceBounds == nil || len(txn.Calldata) == 0 ||
		txn.Nonce == nil ||
		txn.SenderAddress == nil ||
		txn.PayMasterData == nil ||
		txn.AccountDeploymentData == nil {
		return nil, ErrNotAllParametersSet
	}

	txnVersionFelt, err := new(felt.Felt).SetString(string(txn.Version))
	if err != nil {
		return nil, err
	}
	DAUint64, err := DataAvailabilityModeConc(txn.FeeMode, txn.NonceDataMode)
	if err != nil {
		return nil, err
	}
	tipUint64, err := txn.Tip.ToUint64()
	if err != nil {
		return nil, err
	}
	tipAndResourceHash, err := TipAndResourcesHash(tipUint64, txn.ResourceBounds)
	if err != nil {
		return nil, err
	}

	return curve.PoseidonArray(
		prefixInvoke,
		txnVersionFelt,
		txn.SenderAddress,
		tipAndResourceHash,
		curve.PoseidonArray(txn.PayMasterData...),
		chainID,
		txn.Nonce,
		felt.NewFromUint64[felt.Felt](DAUint64),
		curve.PoseidonArray(txn.AccountDeploymentData...),
		curve.PoseidonArray(txn.Calldata...),
	), nil
}

// TransactionHashDeclareV1 calculates the transaction hash for a declare V1 transaction.
//
// Parameters:
//   - txn: The declare V1 transaction to calculate the hash for
//   - chainID: The chain ID as a *felt.Felt
//
// Returns:
//   - *felt.Felt: the calculated transaction hash
//   - error: an error if any
func TransactionHashDeclareV1(txn *rpc.DeclareTxnV1, chainID *felt.Felt) (*felt.Felt, error) {
	//nolint:lll // The link would be unclickable if we break the line.
	// https://docs.starknet.io/architecture-and-concepts/network-architecture/transactions/#v1_deprecated_hash_calculation_2
	if txn.SenderAddress == nil || txn.Version == "" || txn.ClassHash == nil ||
		txn.MaxFee == nil || txn.Nonce == nil {
		return nil, ErrNotAllParametersSet
	}

	calldataHash := curve.PedersenArray(txn.ClassHash)

	txnVersionFelt, err := new(felt.Felt).SetString(string(txn.Version))
	if err != nil {
		return nil, err
	}

	return CalculateDeprecatedTransactionHashCommon(
		prefixDeclare,
		txnVersionFelt,
		txn.SenderAddress,
		&felt.Zero,
		calldataHash,
		txn.MaxFee,
		chainID,
		[]*felt.Felt{txn.Nonce},
	), nil
}

// TransactionHashDeclareV2 calculates the transaction hash for a declare V2
// transaction.
//
// Parameters:
//   - txn: The declare V2 transaction to calculate the hash for
//   - chainID: The chain ID as a *felt.Felt
//
// Returns:
//   - *felt.Felt: the calculated transaction hash
//   - error: an error if any
func TransactionHashDeclareV2(
	txn *rpc.DeclareTxnV2,
	chainID *felt.Felt,
) (*felt.Felt, error) {
	//nolint:lll // The link would be unclickable if we break the line.
	// https://docs.starknet.io/architecture-and-concepts/network-architecture/transactions/#v2_deprecated_hash_calculation
	if txn.CompiledClassHash == nil || txn.SenderAddress == nil ||
		txn.Version == "" ||
		txn.ClassHash == nil ||
		txn.MaxFee == nil ||
		txn.Nonce == nil {
		return nil, ErrNotAllParametersSet
	}

	calldataHash := curve.PedersenArray(txn.ClassHash)

	txnVersionFelt, err := new(felt.Felt).SetString(string(txn.Version))
	if err != nil {
		return nil, err
	}

	return CalculateDeprecatedTransactionHashCommon(
		prefixDeclare,
		txnVersionFelt,
		txn.SenderAddress,
		&felt.Zero,
		calldataHash,
		txn.MaxFee,
		chainID,
		[]*felt.Felt{txn.Nonce, txn.CompiledClassHash},
	), nil
}

// TransactionHashDeclareV3 calculates the transaction hash for a declare V3 transaction.
//
// Parameters:
//   - txn: The declare V3 transaction to calculate the hash for
//   - chainID: The chain ID as a *felt.Felt
//
// Returns:
//   - *felt.Felt: the calculated transaction hash
//   - error: an error if any
func TransactionHashDeclareV3(
	txn *rpc.DeclareTxnV3,
	chainID *felt.Felt,
) (*felt.Felt, error) {
	//nolint:lll // The links would be unclickable if we break the line.
	// https://docs.starknet.io/architecture-and-concepts/network-architecture/transactions/#v3_hash_calculation_2
	// https://github.com/starknet-io/SNIPs/blob/main/SNIPS/snip-8.md#protocol-changes
	if txn.Version == "" || txn.ResourceBounds == nil || txn.Nonce == nil ||
		txn.SenderAddress == nil ||
		txn.PayMasterData == nil ||
		txn.AccountDeploymentData == nil ||
		txn.ClassHash == nil ||
		txn.CompiledClassHash == nil {
		return nil, ErrNotAllParametersSet
	}

	txnVersionFelt, err := new(felt.Felt).SetString(string(txn.Version))
	if err != nil {
		return nil, err
	}
	DAUint64, err := DataAvailabilityModeConc(txn.FeeMode, txn.NonceDataMode)
	if err != nil {
		return nil, err
	}
	tipUint64, err := txn.Tip.ToUint64()
	if err != nil {
		return nil, err
	}

	tipAndResourceHash, err := TipAndResourcesHash(tipUint64, txn.ResourceBounds)
	if err != nil {
		return nil, err
	}

	return curve.PoseidonArray(
		prefixDeclare,
		txnVersionFelt,
		txn.SenderAddress,
		tipAndResourceHash,
		curve.PoseidonArray(txn.PayMasterData...),
		chainID,
		txn.Nonce,
		felt.NewFromUint64[felt.Felt](DAUint64),
		curve.PoseidonArray(txn.AccountDeploymentData...),
		txn.ClassHash,
		txn.CompiledClassHash,
	), nil
}

// TransactionHashBroadcastDeclareV3 calculates the transaction hash for a
// broadcast declare V3 transaction.
//
// Parameters:
//   - txn: The broadcast declare V3 transaction to calculate the hash for
//   - chainID: The chain ID as a *felt.Felt
//
// Returns:
//   - *felt.Felt: the calculated transaction hash
//   - error: an error if any
func TransactionHashBroadcastDeclareV3(
	txn *rpc.BroadcastDeclareTxnV3,
	chainID *felt.Felt,
) (*felt.Felt, error) {
	//nolint:lll // The links would be unclickable if we break the line.
	// https://docs.starknet.io/architecture-and-concepts/network-architecture/transactions/#v3_hash_calculation_2
	// https://github.com/starknet-io/SNIPs/blob/main/SNIPS/snip-8.md#protocol-changes
	if txn.Version == "" || txn.ResourceBounds == nil || txn.Nonce == nil ||
		txn.SenderAddress == nil ||
		txn.PayMasterData == nil ||
		txn.AccountDeploymentData == nil ||
		txn.ContractClass == nil ||
		txn.CompiledClassHash == nil {
		return nil, ErrNotAllParametersSet
	}

	txnVersionFelt, err := new(felt.Felt).SetString(string(txn.Version))
	if err != nil {
		return nil, err
	}
	DAUint64, err := DataAvailabilityModeConc(txn.FeeMode, txn.NonceDataMode)
	if err != nil {
		return nil, err
	}
	tipUint64, err := txn.Tip.ToUint64()
	if err != nil {
		return nil, err
	}

	tipAndResourceHash, err := TipAndResourcesHash(tipUint64, txn.ResourceBounds)
	if err != nil {
		return nil, err
	}

	return curve.PoseidonArray(
		prefixDeclare,
		txnVersionFelt,
		txn.SenderAddress,
		tipAndResourceHash,
		curve.PoseidonArray(txn.PayMasterData...),
		chainID,
		txn.Nonce,
		felt.NewFromUint64[felt.Felt](DAUint64),
		curve.PoseidonArray(txn.AccountDeploymentData...),
		ClassHash(txn.ContractClass),
		txn.CompiledClassHash,
	), nil
}

// TransactionHashDeployAccountV1 calculates the transaction hash for a deploy
// account V1 transaction.
//
// Parameters:
//   - txn: The deploy account V1 transaction to calculate the hash for
//   - contractAddress: The contract address as parameters as a *felt.Felt
//   - chainID: The chain ID as a *felt.Felt
//
// Returns:
//   - *felt.Felt: the calculated transaction hash
//   - error: an error if any
func TransactionHashDeployAccountV1(
	txn *rpc.DeployAccountTxnV1,
	contractAddress, chainID *felt.Felt,
) (*felt.Felt, error) {
	//nolint:lll // The link would be unclickable if we break the line.
	// https://docs.starknet.io/architecture-and-concepts/network-architecture/transactions/#v1_deprecated_hash_calculation_3
	calldata := []*felt.Felt{txn.ClassHash, txn.ContractAddressSalt}
	calldata = append(calldata, txn.ConstructorCalldata...)
	calldataHash := curve.PedersenArray(calldata...)

	versionFelt, err := new(felt.Felt).SetString(string(txn.Version))
	if err != nil {
		return nil, err
	}

	return CalculateDeprecatedTransactionHashCommon(
		prefixDeployAccount,
		versionFelt,
		contractAddress,
		&felt.Zero,
		calldataHash,
		txn.MaxFee,
		chainID,
		[]*felt.Felt{txn.Nonce},
	), nil
}

// TransactionHashDeployAccountV3 calculates the transaction hash for a deploy
// account V3 transaction.
//
// Parameters:
//   - txn: The deploy account V3 transaction to calculate the hash for
//   - contractAddress: The contract address as parameters as a *felt.Felt
//   - chainID: The chain ID as a *felt.Felt
//
// Returns:
//   - *felt.Felt: the calculated transaction hash
//   - error: an error if any
func TransactionHashDeployAccountV3(
	txn *rpc.DeployAccountTxnV3,
	contractAddress, chainID *felt.Felt,
) (*felt.Felt, error) {
	//nolint:lll // The link would be unclickable if we break the line.
	// https://docs.starknet.io/architecture-and-concepts/network-architecture/transactions/#v3_hash_calculation_3
	if txn.Version == "" || txn.ResourceBounds == nil || txn.Nonce == nil ||
		txn.PayMasterData == nil {
		return nil, ErrNotAllParametersSet
	}

	txnVersionFelt, err := new(felt.Felt).SetString(string(txn.Version))
	if err != nil {
		return nil, err
	}
	DAUint64, err := DataAvailabilityModeConc(txn.FeeMode, txn.NonceDataMode)
	if err != nil {
		return nil, err
	}
	tipUint64, err := txn.Tip.ToUint64()
	if err != nil {
		return nil, err
	}
	tipAndResourceHash, err := TipAndResourcesHash(tipUint64, txn.ResourceBounds)
	if err != nil {
		return nil, err
	}

	return curve.PoseidonArray(
		prefixDeployAccount,
		txnVersionFelt,
		contractAddress,
		tipAndResourceHash,
		curve.PoseidonArray(txn.PayMasterData...),
		chainID,
		txn.Nonce,
		felt.NewFromUint64[felt.Felt](DAUint64),
		curve.PoseidonArray(txn.ConstructorCalldata...),
		txn.ClassHash,
		txn.ContractAddressSalt,
	), nil
}

func TipAndResourcesHash(
	tip uint64,
	resourceBounds *rpc.ResourceBoundsMapping,
) (*felt.Felt, error) {
	if resourceBounds == nil {
		return nil, errors.New("resource bounds are nil")
	}
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

	return curve.PoseidonArray(
		felt.NewFromUint64[felt.Felt](tip),
		l1Bounds,
		l2Bounds,
		l1DataGasBounds,
	), nil
}

func DataAvailabilityModeConc(feeDAMode, nonceDAMode rpc.DataAvailabilityMode) (uint64, error) {
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
