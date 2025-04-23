import { defineConfig } from 'vocs'
import { sidebar } from './sidebar'
export default defineConfig({
  title: 'starknet.go',
  titleTemplate: '%s · starknet.go',
  editLink: {
    pattern: 'https://github.com/NethermindEth/starknet.go/edit/main/www/docs/pages/:path',
    text: 'Suggest changes to this page',
  },
  sidebar,
  iconUrl: { light: '/logo.svg', dark: '/logo.svg' },
  logoUrl: { light: '/starknetgo_vertical_light.png', dark: '/starknetgo_vertical_dark.png' },
  socials: [
    {
      icon: 'github',
      link: 'https://github.com/NethermindEth/starknet.go',
    },
    {
      icon: 'x',
      link: 'https://x.com/NethermindStark',
    },
    {
      icon: 'telegram',
      link: 'https://t.me/StarknetGo',
    },
  ],
  theme: {
    accentColor: {
      light: '#ff9318',
      dark: '#ffc517',
    }
  },
  topNav: [
    { text: 'Docs', link: '/introduction/why-starknet-go', match: '/docs' }
  ],
  
})
