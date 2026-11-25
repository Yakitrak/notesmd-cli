---
layout: home

hero:
  name: notesmd-cli
  text: Your Obsidian vault, from the terminal.
  tagline: Create, search, move, and open notes without leaving your workflow. No Obsidian instance required.
  image:
    src: /usage.png
    alt: notesmd-cli in action
  actions:
    - theme: brand
      text: Get Started
      link: /guide/installation
    - theme: alt
      text: View on GitHub
      link: https://github.com/Yakitrak/notesmd-cli

features:
  - icon: 📝
    title: Create & update notes
    details: Spin up new notes with content, overwrite or append to existing ones. Intermediate folders are created automatically.
  - icon: 🔍
    title: Search by name or content
    details: Fuzzy search by filename or full-text content search across your entire vault, with an interactive picker or JSON output for scripting.
  - icon: 🔀
    title: Move with link updates
    details: Rename or relocate notes and every wikilink and markdown link in the vault — including percent-encoded paths — updates automatically.
  - icon: 📅
    title: Daily notes
    details: Create or open today's daily note using your Obsidian config, complete with template variable expansion just like the Obsidian plugin.
  - icon: 🖊️
    title: Any editor
    details: Open notes in Obsidian or your terminal editor via the --editor flag. Set a default once and never pass the flag again.
  - icon: 🖥️
    title: Headless friendly
    details: Works on servers, containers, and CI environments with no GUI or Obsidian installation required.
---
