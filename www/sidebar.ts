import type { Sidebar } from 'vocs'

export const sidebar: Sidebar = {
    '/': [
        {
            text: 'Introduction',
            items: [
                { text: 'Why Starknet.go?', link: '/introduction/why-starknet-go' },
                { text: 'Installation', link: '/introduction/installation' },
                { text: 'Getting Started', link: '/introduction/getting-started' },
                { text: 'Contributing', link: '/introduction/contributing' },
            ],
        },
        {
            text: 'Guides',
            items: [
                { text: 'How to use Starknet.go', link: '/introduction/guides/how-to-use-starknet-go' },
            ],
        },
        { 
            text: 'Account',
            collapsed: true,
            items: [
                { text: 'Overview', link: '/account/' },
                { text: 'Account Management', link: '/account/account-management' },
                { text: 'Transaction Handling', link: '/account/transaction-handling' },
                { text: 'Signature Verification', link: '/account/signature-verification' },
                { text: 'Account Utilities', link: '/account/account-utilities' },
            ]
        },
        {
            text: 'Client',
            collapsed: true,
            items: [
                { text: 'Overview', link: '/client/' },
                { text: 'Client Configuration', link: '/client/configuration' },
                { text: 'Client Methods', link: '/client/methods' },
                { text: 'Client Examples', link: '/client/examples' },
            ]
        },
        {
            text: 'RPC',
            collapsed: true,
            items: [
                { text: 'Overview', link: '/rpc/' },
                { text: 'RPC Methods', link: '/rpc/methods' },
                { text: 'RPC Examples', link: '/rpc/examples' },
            ]
        },
        {
            text: 'Curves',
            collapsed: true,
            items: [
                { text: 'Overview', link: '/curves/' },
                { text: 'Curve Types', link: '/curves/types' },
                { text: 'Curve Operations', link: '/curves/operations' },
            ]
        },
        {
            text: 'Utilities',
            collapsed: true,
            items: [
                { text: 'Overview', link: '/utilities/' },
                { text: 'Common Utilities', link: '/utilities/common' },
                { text: 'Type Utilities', link: '/utilities/types' },
            ]
        },
        { 
            text: 'Developer Tools',
            collapsed: true,
            items: [
                { text: 'Overview', link: '/devnet/' },
                { text: 'Devnet Setup', link: '/devnet/setup' },
                { text: 'Devnet Usage', link: '/devnet/usage' },
            ]
        },
        {
            text: 'ABI',
            collapsed: true,
            items: [
                { text: 'Overview', link: '/abi/' },
                { text: 'ABI Types', link: '/abi/types' },
                { text: 'ABI Encoding', link: '/abi/encoding' },
                { text: 'ABI Decoding', link: '/abi/decoding' },

            ],
        },
    ],
}
