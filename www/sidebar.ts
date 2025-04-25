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
            text: 'Guides',
            items: [
                { text: 'How to use Starknet.go', link: '/docs/guides/how-to-use-starknet-go' },
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
            text: 'RPC',
            collapsed: true,
            items: [
                { text: 'Overview', link: '/docs/rpc/' },
                { text: 'RPC Methods', link: '/docs/rpc/methods' },
                { text: 'RPC Examples', link: '/docs/rpc/examples' },
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
                { text: 'ABI Types', link: '/docs/abi/types' },
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
