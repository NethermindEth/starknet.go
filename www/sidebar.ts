import type { Sidebar } from 'vocs'

export const sidebar: Sidebar = {
    '/docs/': [
        {
            text: 'Introduction',
            items: [
                { text: 'Why Starknet.go?', link: '/docs/introduction/why-starknet-go' },
                { text: 'Installation', link: '/docs/introduction/installation' },
                { text: 'Getting Started', link: '/docs/introduction/getting-started' },
                { text: 'Contributing', link: '/docs/introduction/contributing' },
            ],
        },
        {
            text: 'RPC',
            collapsed: true,
            items: [
                { text: 'Overview', link: '/docs/rpc/' },
                { text: 'Block',
                    items: [
                        { text: 'BlockNumber', link: '/docs/rpc/methods/block-number' },
                        { text: 'BlockHashAndNumber', link: '/docs/rpc/methods/block-hash-and-number' },
                        { text: 'GetBlockWithTxHashes', link: '/docs/rpc/methods/get-block-with-tx-hashes' },
                        { text: 'GetStateUpdate', link: '/docs/rpc/methods/get-state-update' },
                        { text: 'GetBlockTransactionCount', link: '/docs/rpc/methods/get-block-transaction-count' },
                        { text: 'GetBlockWithTxs', link: '/docs/rpc/methods/get-block-with-txs' },
                        { text: 'GetBlockWithReceipts', link: '/docs/rpc/methods/get-block-with-receipts' },
                    ]
                },
                { text: 'Call', link: '/docs/rpc/methods/call' },
                    {
                        text: 'Chain',
                        items: [
                            { text: 'ChainId', link: '/docs/rpc/methods/chain-id' },
                            { text: 'Syncing', link: '/docs/rpc/methods/syncing' },
                        ]
                    },
                    {
                        text: 'Contract',
                        items: [
                            { text: 'GetClass', link: '/docs/rpc/methods/get-class' },
                            { text: 'GetClassAt', link: '/docs/rpc/methods/get-class-at' },
                            { text: 'GetClassHashAt', link: '/docs/rpc/methods/get-class-hash-at' },
                            { text: 'GetNonce', link: '/docs/rpc/methods/get-nonce' },
                            { text: 'EstimateFee', link: '/docs/rpc/methods/estimate-fee' },
                            { text: 'EstimateMessageFee', link: '/docs/rpc/methods/estimate-message-fee' },
                            { text: 'GetStorageProof', link: '/docs/rpc/methods/get-storage-proof' },
                        ]
                    },
                    { text: 'GetEvents', link: '/docs/rpc/methods/get-events' },
                    {
                        text: 'Transaction',
                        items: [
                            { text: 'GetTransactionByHash', link: '/docs/rpc/methods/get-transaction-by-hash' },
                            { text: 'GetTransactionByBlockIdAndIndex', link: '/docs/rpc/methods/get-transaction-by-block-id-and-indes' },
                            { text: 'GetTransactionReceipt', link: '/docs/rpc/methods/get-transaction-receipt' },
                            { text: 'GetTransactionStatus', link: '/docs/rpc/methods/get-transaction-status' },
                            { text: 'GetMessagesStatus', link: '/docs/rpc/methods/get-messages-status' },
                        ]
                    },
                    {
                        text: 'Trace',
                        items: [
                            { text: 'TraceTransaction', link: '/docs/rpc/methods/trace-transaction' },
                            { text: 'TraceBlockTransactions', link: '/docs/rpc/methods/trace-block-transactions' },
                            { text: 'SimulateTransaction', link: '/docs/rpc/methods/simulate-transaction' },

                        ]
                    },
                    {
                        text: 'Write',
                        items: [
                            { text: 'AddInvokeTransaction', link: '/docs/rpc/methods/add-invoke-transaction' },
                            { text: 'AddDeclareTransaction', link: '/docs/rpc/methods/add-declare-transaction' },
                            { text: 'AddDeployAccountTransaction', link: '/docs/rpc/methods/add-deploy-account-transaction' },

                        ]
                    },
                    { text: 'SpecVersion', link: '/docs/rpc/methods/spec-version' },
                    { text: 'GetCompiledCasm', link: '/docs/rpc/methods/get-compiled-casm' },
                    { text: 'RPC Examples', link: '/docs/rpc/examples' },
                ],
            
        },
        { 
            text: 'Account',
            collapsed: true,
            items: [
                { text: 'Overview', link: '/docs/account/' },
                { text: 'Account Management', link: '/docs/account/account-management' },
                { text: 'Transaction Handling', link: '/docs/account/transaction-handling' },
                { text: 'Signature Verification', link: '/docs/account/signature-verification' },
                { text: 'Account Utilities', link: '/docs/account/account-utilities' },
            ]
        },
        {
            text: 'Client',
            collapsed: true,
            items: [
                { text: 'Overview', link: '/docs/client/' },
                { text: 'Client Configuration', link: '/docs/client/configuration' },
                { text: 'Client Methods', link: '/docs/client/methods' },
                { text: 'Client Examples', link: '/docs/client/examples' },
            ]
        },
        {
            text: 'Curves',
            collapsed: true,
            items: [
                { text: 'Overview', link: '/docs/curves/' },
                { text: 'Curve Types', link: '/docs/curves/types' },
                { text: 'Curve Operations', link: '/docs/curves/operations' },
            ]
        },
        {
            text: 'Utilities',
            collapsed: true,
            items: [
                { text: 'Overview', link: '/docs/utilities/' },
                { text: 'Common Utilities', link: '/docs/utilities/common' },
                { text: 'Type Utilities', link: '/docs/utilities/types' },
            ]
        },
        { 
            text: 'Developer Tools',
            collapsed: true,
            items: [
                { text: 'Overview', link: '/docs/devnet/' },
                { text: 'Devnet Setup', link: '/docs/devnet/setup' },
                { text: 'Devnet Usage', link: '/docs/devnet/usage' },
            ]
        },
        {
            text: 'ABI',
            collapsed: true,
            items: [
                { text: 'Overview', link: '/docs/abi/' },
                { text: 'ABI Encoding', link: '/docs/abi/encoding' },
                { text: 'ABI Decoding', link: '/docs/abi/decoding' },

            ],
        },
        {
            text: 'Examples',
            collapsed: true,
            items: [
                { text: 'Overview', link: '/docs/examples/' },
                { text: 'Deploy Account', link: '/docs/examples/deploy-account' },
                { text: 'Deploy Contract', link: '/docs/examples/deploy-contract' },
                { text: 'Internal Transactions', link: '/docs/examples/internal-transactions' },
                { text: 'Invoke Contract', link: '/docs/examples/invoke-contract' },
                { text: 'Read Events', link: '/docs/examples/read-events' },
                { text: 'Call', link: '/docs/examples/simple-call' },
                { text: 'Declare', link: '/docs/examples/declare' },
                { text: 'Typed Data', link: '/docs/examples/typed-data' },
                { text: 'WebSocket', link: '/docs/examples/websocket' },
            ]
        }
    ],
}
