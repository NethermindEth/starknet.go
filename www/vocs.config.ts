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
  ogImageUrl: {
    '/': '/og-image.png',
  },
  iconUrl: { light: '/favicons/light.png', dark: '/favicons/dark.png' },
  logoUrl: { 
    light: '/Starknet.Go_Light_Powered_by_Nethermind.png', 
    dark: '/Starknet.Go_Dark_Powered_by_Nethermind.png' 
  },
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