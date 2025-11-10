import type { Sidebar } from "vocs";

export const sidebar: Sidebar = {
  "/docs/": [
    {
      text: "Introduction",
      items: [
        {
          text: "Why Starknet.go?",
          link: "/docs/introduction/why-starknet-go",
        },
        { text: "Installation", link: "/docs/introduction/installation" },
        { text: "Getting Started", link: "/docs/introduction/getting-started" },
        { text: "Contributing", link: "/docs/introduction/contributing" },
      ],
    },
    {
      text: "account",
      collapsed: true,
      items: [
        { text: "Overview", link: "/docs/account/" },
        {
          text: "Concepts",
          collapsed: true,
          items: [
            { text: "Transaction Signing", link: "/docs/account/concepts/transaction-signing" },
          ],
        },
        {
          text: "Functions",
          collapsed: true,
          items: [
            { text: "FmtCallDataCairo0", link: "/docs/account/functions/fmt-calldata-cairo0" },
            { text: "FmtCallDataCairo2", link: "/docs/account/functions/fmt-calldata-cairo2" },
            { text: "PrecomputeAccountAddress", link: "/docs/account/functions/precompute-account-address" },
          ],
        },
        {
          text: "Methods",
          collapsed: true,
          items: [
            { text: "NewAccount", link: "/docs/account/methods/new-account" },
            { text: "BuildAndEstimateDeployAccountTxn", link: "/docs/account/methods/build-and-estimate-deploy-account-txn" },
            { text: "BuildAndSendDeclareTxn", link: "/docs/account/methods/build-and-send-declare-txn" },
            { text: "BuildAndSendInvokeTxn", link: "/docs/account/methods/build-and-send-invoke-txn" },
            { text: "DeployContractWithUDC", link: "/docs/account/methods/deploy-contract-with-udc" },
            { text: "FmtCalldata", link: "/docs/account/methods/fmt-calldata" },
            { text: "Nonce", link: "/docs/account/methods/nonce" },
            { text: "SendTransaction", link: "/docs/account/methods/send-transaction" },
            { text: "Sign", link: "/docs/account/methods/sign" },
            { text: "SignDeclareTransaction", link: "/docs/account/methods/sign-declare-transaction" },
            { text: "SignDeployAccountTransaction", link: "/docs/account/methods/sign-deploy-account-transaction" },
            { text: "SignInvokeTransaction", link: "/docs/account/methods/sign-invoke-transaction" },
            { text: "TransactionHashDeclare", link: "/docs/account/methods/transaction-hash-declare" },
            { text: "TransactionHashDeployAccount", link: "/docs/account/methods/transaction-hash-deploy-account" },
            { text: "TransactionHashInvoke", link: "/docs/account/methods/transaction-hash-invoke" },
            { text: "Verify", link: "/docs/account/methods/verify" },
            { text: "WaitForTransactionReceipt", link: "/docs/account/methods/wait-for-transaction-receipt" },
          ],
        },
        {
          text: "Keystore",
          collapsed: true,
          items: [
            { text: "GetRandomKeys", link: "/docs/account/keystore/get-random-keys" },
            { text: "NewMemKeystore", link: "/docs/account/keystore/new-mem-keystore" },
            { text: "SetNewMemKeystore", link: "/docs/account/keystore/set-new-mem-keystore" },
            { text: "Get", link: "/docs/account/keystore/get" },
            { text: "Put", link: "/docs/account/keystore/put" },
            { text: "Sign", link: "/docs/account/keystore/sign" },
          ],
        },
        {
          text: "Types",
          link: "/docs/account/types",
        },
      ],
    },
    {
      text: "rpc",
      collapsed: true,
      items: [
        { text: "Overview", link: "/docs/rpc/" },
        { text: "Helper Functions", link: "/docs/rpc/helper-functions" },
        {
          text: "Block Methods",
          collapsed: true,
          items: [
            { text: "BlockNumber", link: "/docs/rpc/methods/block-number" },
            {
              text: "BlockHashAndNumber",
              link: "/docs/rpc/methods/block-hash-and-number",
            },
            { text: "BlockWithTxs", link: "/docs/rpc/methods/block-with-txs" },
            {
              text: "BlockWithTxHashes",
              link: "/docs/rpc/methods/block-with-tx-hashes",
            },
            {
              text: "BlockWithReceipts",
              link: "/docs/rpc/methods/block-with-receipts",
            },
            {
              text: "BlockTransactionCount",
              link: "/docs/rpc/methods/block-transaction-count",
            },
            {
              text: "StateUpdate",
              link: "/docs/rpc/methods/state-update",
            },
          ],
        },
        {
          text: "Transaction Methods",
          collapsed: true,
          items: [
            {
              text: "TransactionByHash",
              link: "/docs/rpc/methods/transaction-by-hash",
            },
            {
              text: "TransactionReceipt",
              link: "/docs/rpc/methods/transaction-receipt",
            },
            {
              text: "TransactionStatus",
              link: "/docs/rpc/methods/get-transaction-status",
            },
            {
              text: "TransactionByBlockIdAndIndex",
              link: "/docs/rpc/methods/get-transaction-by-block-id-and-index",
            },
            {
              text: "MessagesStatus",
              link: "/docs/rpc/methods/get-messages-status",
            },
          ],
        },
        { text: "Call", link: "/docs/rpc/methods/call" },
        {
          text: "Chain Methods",
          collapsed: true,
          items: [
            { text: "ChainID", link: "/docs/rpc/methods/chain-ID" },
            { text: "SpecVersion", link: "/docs/rpc/methods/spec-version" },
            { text: "Syncing", link: "/docs/rpc/methods/syncing" },
          ],
        },
        {
          text: "Contract Methods",
          collapsed: true,
          items: [
            { text: "Class", link: "/docs/rpc/methods/get-class" },
            { text: "ClassAt", link: "/docs/rpc/methods/get-class-at" },
            {
              text: "ClassHashAt",
              link: "/docs/rpc/methods/get-class-hash-at",
            },
            { text: "Nonce", link: "/docs/rpc/methods/get-nonce" },
            { text: "StorageAt", link: "/docs/rpc/methods/get-storage-at" },
            { text: "StorageProof", link: "/docs/rpc/methods/storage-proof" },
          ],
        },
        {
          text: "Fee Methods",
          collapsed: true,
          items: [
            { text: "EstimateFee", link: "/docs/rpc/methods/estimate-fee" },
            {
              text: "EstimateMessageFee",
              link: "/docs/rpc/methods/estimate-message-fee",
            },
          ],
        },
        { text: "Events", link: "/docs/rpc/methods/events" },
        {
          text: "Trace Methods",
          collapsed: true,
          items: [
            {
              text: "TraceTransaction",
              link: "/docs/rpc/methods/trace-transaction",
            },
            {
              text: "TraceBlockTransactions",
              link: "/docs/rpc/methods/trace-block-transactions",
            },
            {
              text: "SimulateTransactions",
              link: "/docs/rpc/methods/simulate-transaction",
            },
          ],
        },
        {
          text: "Write Methods",
          collapsed: true,
          items: [
            {
              text: "AddInvokeTransaction",
              link: "/docs/rpc/methods/add-invoke-transaction",
            },
            {
              text: "AddDeclareTransaction",
              link: "/docs/rpc/methods/add-declare-transaction",
            },
            {
              text: "AddDeployAccountTransaction",
              link: "/docs/rpc/methods/add-deploy-account-transaction",
            },
          ],
        },
        {
          text: "CompiledCasm",
          link: "/docs/rpc/methods/compiled-casm",
        },
      ],
    },
    {
      text: "devnet",
      collapsed: true,
      items: [
        { text: "Overview", link: "/docs/devnet/" },
        {
          text: "Methods",
          collapsed: true,
          items: [
            { text: "NewDevNet", link: "/docs/devnet/methods/new-devnet" },
            { text: "IsAlive", link: "/docs/devnet/methods/is-alive" },
            { text: "Accounts", link: "/docs/devnet/methods/accounts" },
            { text: "Mint", link: "/docs/devnet/methods/mint" },
            { text: "FeeToken", link: "/docs/devnet/methods/fee-token" },
          ],
        },
        { text: "Types", link: "/docs/devnet/types" },
      ],
    },
    {
      text: "contracts",
      collapsed: true,
      items: [
        { text: "Overview", link: "/docs/contracts/" },
        {
          text: "Functions",
          collapsed: true,
          items: [
            { text: "PrecomputeAddress", link: "/docs/contracts/functions/precompute-address" },
            { text: "UnmarshalCasmClass", link: "/docs/contracts/functions/unmarshal-casm-class" },
          ],
        },
        {
          text: "Types",
          collapsed: true,
          items: [
            { text: "ContractClass", link: "/docs/contracts/functions/contract-class" },
            { text: "CasmClass", link: "/docs/contracts/functions/casm-class" },
          ],
        },
      ],
    },
    {
      text: "utils",
      collapsed: true,
      items: [
        { text: "Overview", link: "/docs/utils/" },
        {
          text: "Type Conversions",
          collapsed: true,
          items: [
            { text: "HexToFelt", link: "/docs/utils/functions/hex-to-felt" },
            { text: "HexArrToFelt", link: "/docs/utils/functions/hex-arr-to-felt" },
            { text: "FeltToBigInt", link: "/docs/utils/functions/felt-to-big-int" },
            { text: "BigIntToFelt", link: "/docs/utils/functions/big-int-to-felt" },
            { text: "HexToBN", link: "/docs/utils/functions/hex-to-bn" },
            { text: "BigToHex", link: "/docs/utils/functions/big-to-hex" },
            { text: "Uint64ToFelt", link: "/docs/utils/functions/uint64-to-felt" },
          ],
        },
        {
          text: "Unit Conversions",
          collapsed: true,
          items: [
            { text: "ETHToWei", link: "/docs/utils/functions/eth-to-wei" },
            { text: "WeiToETH", link: "/docs/utils/functions/wei-to-eth" },
            { text: "STRKToFRI", link: "/docs/utils/functions/strk-to-fri" },
            { text: "FRIToSTRK", link: "/docs/utils/functions/fri-to-strk" },
          ],
        },
        {
          text: "Transaction Builders",
          collapsed: true,
          items: [
            { text: "BuildInvokeTxn", link: "/docs/utils/functions/build-invoke-txn" },
            { text: "BuildDeclareTxn", link: "/docs/utils/functions/build-declare-txn" },
            { text: "BuildDeployAccountTxn", link: "/docs/utils/functions/build-deploy-account-txn" },
          ],
        },
        {
          text: "String & Hex Utilities",
          collapsed: true,
          items: [
            { text: "HexToShortStr", link: "/docs/utils/functions/hex-to-short-str" },
            { text: "StrToHex", link: "/docs/utils/functions/str-to-hex" },
          ],
        },
        {
          text: "Selector & Hashing",
          collapsed: true,
          items: [
            { text: "GetSelectorFromName", link: "/docs/utils/functions/get-selector-from-name" },
            { text: "GetSelectorFromNameFelt", link: "/docs/utils/functions/get-selector-from-name-felt" },
            { text: "Keccak256", link: "/docs/utils/functions/keccak256" },
          ],
        },
      ],
    },
    {
      text: "typeddata",
      collapsed: true,
      items: [
        { text: "Overview", link: "/docs/typeddata/" },
        { text: "TypedData", link: "/docs/typeddata/typed-data" },
        { text: "Domain", link: "/docs/typeddata/domain" },
        { text: "Functions", link: "/docs/typeddata/functions" },
      ],
    },
    {
      text: "hash",
      collapsed: true,
      items: [
        { text: "Overview", link: "/docs/hash/" },
        {
          text: "Concepts",
          collapsed: true,
          items: [
            { text: "Transaction Hashing", link: "/docs/hash/concepts/transaction-hashing" },
          ],
        },
        {
          text: "Functions",
          collapsed: true,
          items: [
            { text: "ClassHash", link: "/docs/hash/functions/class-hash" },
            { text: "CompiledClassHash", link: "/docs/hash/functions/compiled-class-hash" },
            { text: "TransactionHashInvokeV0", link: "/docs/hash/functions/transaction-hash-invoke-v0" },
            { text: "TransactionHashInvokeV1", link: "/docs/hash/functions/transaction-hash-invoke-v1" },
            { text: "TransactionHashInvokeV3", link: "/docs/hash/functions/transaction-hash-invoke-v3" },
            { text: "TransactionHashDeclareV1", link: "/docs/hash/functions/transaction-hash-declare-v1" },
            { text: "TransactionHashDeclareV2", link: "/docs/hash/functions/transaction-hash-declare-v2" },
            { text: "TransactionHashDeclareV3", link: "/docs/hash/functions/transaction-hash-declare-v3" },
            { text: "TransactionHashBroadcastDeclareV3", link: "/docs/hash/functions/transaction-hash-broadcast-declare-v3" },
            { text: "TransactionHashDeployAccountV1", link: "/docs/hash/functions/transaction-hash-deploy-account-v1" },
            { text: "TransactionHashDeployAccountV3", link: "/docs/hash/functions/transaction-hash-deploy-account-v3" },
            { text: "CalculateDeprecatedTransactionHashCommon", link: "/docs/hash/functions/calculate-deprecated-transaction-hash-common" },
            { text: "TipAndResourcesHash", link: "/docs/hash/functions/tip-and-resources-hash" },
            { text: "DataAvailabilityModeConc", link: "/docs/hash/functions/data-availability-mode-conc" },
          ],
        },
      ],
    },
    {
      text: "curve",
      collapsed: true,
      items: [
        { text: "Overview", link: "/docs/curve/" },
        {
          text: "Functions",
          collapsed: true,
          items: [
            { text: "Pedersen", link: "/docs/curve/functions/pedersen" },
            { text: "PedersenArray", link: "/docs/curve/functions/pedersen-array" },
            { text: "Poseidon", link: "/docs/curve/functions/poseidon" },
            { text: "PoseidonArray", link: "/docs/curve/functions/poseidon-array" },
            { text: "StarknetKeccak", link: "/docs/curve/functions/starknet-keccak" },
            { text: "ComputeHashOnElements", link: "/docs/curve/functions/compute-hash-on-elements" },
            { text: "HashPedersenElements", link: "/docs/curve/functions/hash-pedersen-elements" },
            { text: "GetRandomKeys", link: "/docs/curve/functions/get-random-keys" },
            { text: "PrivateKeyToPoint", link: "/docs/curve/functions/private-key-to-point" },
            { text: "GetYCoordinate", link: "/docs/curve/functions/get-y-coordinate" },
            { text: "Sign", link: "/docs/curve/functions/sign" },
            { text: "SignFelts", link: "/docs/curve/functions/sign-felts" },
            { text: "Verify", link: "/docs/curve/functions/verify" },
            { text: "VerifyFelts", link: "/docs/curve/functions/verify-felts" },
          ],
        },
      ],
    },
    {
      text: "merkle",
      collapsed: true,
      items: [
        { text: "Overview", link: "/docs/merkle/" },
        {
          text: "Functions",
          collapsed: true,
          items: [
            { text: "MerkleHash", link: "/docs/merkle/functions/merkle-hash" },
            { text: "ProofMerklePath", link: "/docs/merkle/functions/proof-merkle-path" },
            { text: "NewFixedSizeMerkleTree", link: "/docs/merkle/functions/new-fixed-size-merkle-tree" },
            { text: "Proof (Method)", link: "/docs/merkle/functions/proof" },
          ],
        },
      ],
    },
    {
      text: "paymaster",
      collapsed: true,
      items: [
        { text: "Overview", link: "/docs/paymaster/" },
        {
          text: "Methods",
          collapsed: true,
          items: [
            { text: "New", link: "/docs/paymaster/methods/new" },
            { text: "IsAvailable", link: "/docs/paymaster/methods/is-available" },
            { text: "GetSupportedTokens", link: "/docs/paymaster/methods/get-supported-tokens" },
            { text: "BuildTransaction", link: "/docs/paymaster/methods/build-transaction" },
            { text: "ExecuteTransaction", link: "/docs/paymaster/methods/execute-transaction" },
            { text: "TrackingIDToLatestHash", link: "/docs/paymaster/methods/tracking-id-to-latest-hash" },
          ],
        },
      ],
    },
    {
      text: "client",
      collapsed: true,
      items: [
        { text: "Overview", link: "/docs/client/" },
        { text: "Client Functions", link: "/docs/client/functions" },
        { text: "Client Methods", link: "/docs/client/methods" },
        { text: "Types", link: "/docs/client/types" },
      ],
    },
    {
      text: "examples",
      collapsed: true,
      items: [
        { text: "Overview", link: "/docs/examples/" },
        { text: "Deploy Account", link: "/docs/examples/deploy-account" },
        { text: "Deploy Contract UDC", link: "/docs/examples/deploy-contract-udc" },
        { text: "Invoke", link: "/docs/examples/invoke" },
        { text: "Read Events", link: "/docs/examples/read-events" },
        { text: "Simple Call", link: "/docs/examples/simple-call" },
        { text: "Simple Declare", link: "/docs/examples/simple-declare" },
        { text: "Typed Data", link: "/docs/examples/typed-data" },
        { text: "WebSocket", link: "/docs/examples/websocket" },
      ],
    },
  ],
};
