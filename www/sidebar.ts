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
                { text: 'Platform Compatibility', link: '/docs/introduction/platform-compatibility' },
                { text: 'FAQ', link: '/docs/introduction/faq' },
            ],
        },
        {
            text: 'Guides',
            items: [
                { text: 'How to use Starknet.go', link: '/docs/guides/how-to-use-starknet-go' },
            ],
        },
        {
            text: 'Examples',
            items: [
                { text: 'Example', link: '/docs/example/example' },
            ],
        },
    ],
}
