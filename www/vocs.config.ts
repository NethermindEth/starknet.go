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
      light: '#FF4A04',
      dark: '#FFE9CF',
    },
  },
  topNav: [
    { text: 'Docs', link: '/docs/introduction/why-starknet-go', match: '/docs' }
  ],
  
})
