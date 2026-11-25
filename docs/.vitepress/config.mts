import { defineConfig } from 'vitepress'

export default defineConfig({
  title: 'notesmd-cli',
  description: 'Interact with your Obsidian vault from the terminal — no Obsidian instance required.',
  base: '/notesmd-cli/',

  head: [
    ['link', { rel: 'icon', href: '/notesmd-cli/favicon.svg', type: 'image/svg+xml' }],
  ],

  themeConfig: {
    logo: { svg: '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 28 28"><rect width="28" height="28" rx="7" fill="#7c6af7"/><path d="M7 8h14M7 14h9M7 20h11" stroke="white" stroke-width="2.2" stroke-linecap="round"/></svg>' },

    nav: [
      { text: 'Guide', link: '/guide/installation' },
      { text: 'Commands', link: '/commands/vault' },
      {
        text: 'v0.3.9',
        items: [
          { text: 'Changelog', link: 'https://github.com/Yakitrak/notesmd-cli/releases' },
          { text: 'Contributing', link: '/contributing' },
        ],
      },
    ],

    sidebar: {
      '/guide/': [
        {
          text: 'Getting Started',
          items: [
            { text: 'Installation', link: '/guide/installation' },
            { text: 'Quick Start', link: '/guide/quick-start' },
            { text: 'Editor Flag', link: '/guide/editor' },
          ],
        },
        {
          text: 'More',
          items: [
            { text: 'Headless / No Obsidian', link: '/guide/headless' },
            { text: 'Migrating from obsidian-cli', link: '/guide/migration' },
          ],
        },
      ],
      '/commands/': [
        {
          text: 'Commands',
          items: [
            { text: 'Vault Management', link: '/commands/vault' },
            { text: 'Create & Update', link: '/commands/create' },
            { text: 'Open & Print', link: '/commands/open' },
            { text: 'Move & Delete', link: '/commands/move' },
            { text: 'Search', link: '/commands/search' },
            { text: 'Daily Notes', link: '/commands/daily' },
            { text: 'Frontmatter', link: '/commands/frontmatter' },
          ],
        },
      ],
    },

    socialLinks: [
      { icon: 'github', link: 'https://github.com/Yakitrak/notesmd-cli' },
    ],

    search: {
      provider: 'local',
    },

    footer: {
      message: 'Released under the MIT License.',
      copyright: 'Copyright © Yakitrak',
    },

    editLink: {
      pattern: 'https://github.com/Yakitrak/notesmd-cli/edit/main/docs/:path',
      text: 'Edit this page on GitHub',
    },
  },
})
