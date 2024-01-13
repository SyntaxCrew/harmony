import { defineConfig } from 'vitepress'

// https://vitepress.dev/reference/site-config
export default defineConfig({
  base: '/docs',
  title: "Harmony",
  description: "High-performance Golang web framework. Get started with guides, API references, and examples for efficient and scalable web development.",
  themeConfig: {
    nav: [
      { text: 'Home', link: '/' },
      { text: 'Docs', link: '/getting-started/quick-start' }
    ],

    sidebar: [
      {
        text: 'Getting Started',
        base: '/getting-started',
        items: [
          { text: 'Quick Start', link: '/quick-start' },
        ]
      },
      {
        text: 'Guide',
        base: '/guide',
        items: [
          { text: 'Binding', link: '/binding' },
          { text: 'Context', link: '/context' },
          {
            text: 'Middlewares',
            collapsed: true,
            base: '/guide/middlewares',
            items: [
              { text: 'Gzip', link: '/gzip', },
              { text: 'Logger', link: '/logger' },
            ]
          },
          { text: 'Routing', link: '/routing' },
        ]
      },
    ],

    socialLinks: [
      { icon: 'github', link: 'https://github.com/SyntaxCrew/harmony' }
    ],

    search: {
      provider: 'local'
    },
  }
})
