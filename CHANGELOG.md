# Changelog

## [3.0.0] — 2026-03-25

### 🚀 The TUI Revolution
This is a major architectural shift. `gpad` has evolved from a CLI tool into a high-performance **Terminal User Interface (TUI)** designed for maximum speed and keyboard-driven efficiency.

### New Features
- **Full-Screen TUI** — Launch into a persistent workspace with a dual-pane layout: a file tree on the left and a live-rendering preview on the right.
- **Vim-Motion Navigation** — Navigate your entire vault using standard Vim keys (`h`, `j`, `k`, `l`, `g`, `G`). Stay on the home row for all operations.
- **Interactive Panels** — New dedicated interfaces for advanced workflows:
    - **Search (`F`)**: Fuzzy title and full-text body search with live filtering.
    - **Tags (`T`)**: A dedicated browser for your full tag index.
    - **Graph (`W`)**: A visual link-map of your notes sorted by connection density.
    - **Daily (`D`)**: A chronological timeline of your daily notes.
    - **Links (`L`)**: Detailed panel for inspecting backlinks and outlinks.
- **Command & Filter Modes** — Use `:` for a powerful internal command line and `/` to instantly filter your file tree.
- **Custom Keybindings** — Fully remappable controls via `~/.gpad/keybinds.json`. Change any TUI action to suit your personal workflow.

### Improvements & Integration
- **Live Preview Footer** — The TUI preview now displays real-time metadata including word count, line count, and link count.
- **Integrated Git Sync** — Background syncing with an `↑ pushing...` status indicator. No more waiting for the CLI to finish a push before you can keep writing.
- **Visual Pining** — Mark important notes with `P`. Pinned notes are now visually highlighted with a `★` in the tree.
- **Smart Directory Handling** — Navigate directories like `lf` or `ranger` using `Enter` to drill down and `h` or `-` to move back up.

### Breaking Changes
- **Default Behavior**: Running `gpad` without arguments now launches the TUI instead of the help menu.
- **CLI Mode**: Legacy CLI outputs (like `ls` or `view`) are still available but are now secondary to the TUI experience.
- **Config Migration**: The configuration has moved to a structured `~/.gpad/` directory containing `config.json`, `keybinds.json`, and `index.json`.

---

## [2.0.0] — 2026-03-22
*Initial implementation of the knowledge-vault engine.*

- **Wikilinks**: Support for `[[note]]`, aliases, and heading anchors.
- **Backlinks & Graphing**: Added `gpad links` and `gpad graph`.
- **Tagging**: Support for YAML frontmatter and inline `#hashtags`.
- **Templates**: Introduced `gpad new -t` with custom placeholders.
- **Daily Notes**: Introduced `gpad today` functionality.

---

## [1.0.0] — Initial Release
- Basic note creation and editing.
- Git pull/push integration.
- Markdown rendering via pager.
