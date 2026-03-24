<div align="center">

<img src="https://img.shields.io/badge/built%20with-Go-00ADD8?style=flat-square&logo=go" />
<img src="https://img.shields.io/badge/license-MIT-green?style=flat-square" />
<img src="https://img.shields.io/badge/platform-Linux%20%7C%20macOS%20%7C%20Windows-lightgrey?style=flat-square" />

# gpad

**A terminal-native knowledge vault.**

Write notes in Markdown. Connect them with `[[wikilinks]]`. Tag them with `#tags`.  
Navigate everything in a full-screen TUI. Sync across machines with Git.  
Zero config to start. One static binary.

</div>

---

## Why gpad

Most note tools are either too simple (just files) or too heavy (Electron apps, databases, subscriptions). gpad is for people who live in the terminal and want Obsidian-style linking and tagging without leaving the shell.

- **Plain Markdown files** — open with any editor, grep with any tool, version with any Git host
- **Full-screen TUI** — file tree, live preview, panels for search/tags/graph/links all in one
- **`[[wikilinks]]`** — write them while taking notes, gpad builds the graph automatically
- **Git sync** — every save pushes in the background, or sync manually, or work offline forever
- **Fast** — built in Go, no runtime, opens instantly

---

## Installation

**Build from source (Go 1.22+):**

```bash
git clone https://github.com/Abhishek-Krishna-A-M/gpad
cd gpad
go mod tidy
go build -o gpad ./cmd/gpad/
sudo mv gpad /usr/local/bin/
```

**Linux / macOS:**

```bash
curl -fsSL https://raw.githubusercontent.com/Abhishek-Krishna-A-M/gpad/main/install.sh | sh
```

**Windows:**

```powershell
powershell -ExecutionPolicy Bypass -File install.ps1
```

---

## Quick start

```bash
gpad                      # open the TUI
gpad today                # open today's daily note
gpad open ideas.md        # create and open a note from the shell
gpad git init git@github.com:you/notes.git   # connect git sync
```

---

## The TUI

`gpad` with no arguments opens the full-screen TUI. Press `?` for the full help screen.

```
+-------------------------------------------------------------+
| gpad /                                    42 notes    v     |
+------------------+------------------------------------------+
| v daily/         | quantum-foam                             |
|   2026-03-22.md  | ─────────────                            |
|   2026-03-21.md  | ## The idea                              |
| > college/       |                                          |
|   ideas.md    *  |   Connects to -> physics and             |
|   quantum-foam   |   -> mathematics. Core idea              |
|   todo.md        |   relates to #science.                   |
|                  |                                          |
|                  | ─────────────────────────────            |
|                  | 312 words · 18 lines · 2 links           |
|                  | #science #ideas                          |
|                  | <- ideas.md                              |
+------------------+------------------------------------------+
| : cmd  / filter  F search  T tags  W graph  D daily  ? help |
+-------------------------------------------------------------+
| NORMAL    quantum-foam.md                              5/7  |
+-------------------------------------------------------------+
```

**Left panel** — file tree. Navigate with `j`/`k`, enter directories with `Enter`, go back with `h` or `-`.  
**Right panel** — live preview that updates as you move. Shows rendered markdown, tags, word count, and backlinks.  
**Bottom bar** — hints in normal mode, command input with `:`, filter input with `/`.

---

## Keybindings

Press `?` inside the TUI for the full help screen. All bindings are remappable in `~/.gpad/keybinds.json`.

### Navigation

| Key | Action |
|---|---|
| `j` / down | move down |
| `k` / up | move up |
| `g` | jump to top |
| `G` | jump to bottom |
| `Enter` | directory: cd into it · note: open in editor |
| `l` / right | expand directory in-place |
| `h` / left | collapse directory or go to parent |
| `-` | go up one directory level |
| `Space` | toggle expand directory |

Directory navigation works like lf: `Enter` replaces the tree with the directory's contents. `h` or `-` goes back up.

### File operations

| Key | Action |
|---|---|
| `o` | open in editor |
| `v` | view rendered markdown (less pager) |
| `n` | new note in current directory |
| `d` / `x` | delete (asks confirmation) |
| `r` | rename — prefills `:mv <current> ` in command bar |
| `y` | yank path |
| `p` | paste yanked note into current directory |
| `P` | toggle pin — pinned notes show `*` in the tree |

### Panels

| Key | Panel |
|---|---|
| `F` | Search — fuzzy title + full-text body, type to filter |
| `T` | Tags — full tag index, Enter drills into tag's notes |
| `W` | Graph — linked notes sorted by connection count |
| `D` | Daily — all daily notes newest first |
| `L` | Links — backlinks and outlinks for selected note |

Inside any panel: `j`/`k` navigate, `/` filter, `Enter` open, `Esc`/`q` close.

### General

| Key | Action |
|---|---|
| `s` | git sync (pull + push) |
| `t` | open today's daily note |
| `:` | command mode |
| `/` | filter tree |
| `?` | full help screen |
| `Ctrl+L` | force redraw |
| `q` / `Esc` | quit |

---

## Command mode

Press `:` to enter command mode. Every gpad feature is available here.

### Notes

```
:new <name.md>               create note in current directory
:new <name.md> -t meeting    create from template
:open <note>                 open note by path
:view <note>                 view rendered in pager
:mkdir <n>                create directory
```

### File management

```
:mv <src> <dest>             move or rename
:cp <src> <dest>             copy
:rm <note>                   delete
```

### Panels and search

```
:find <query>                open search panel with query
:tags                        tag browser
:graph                       graph panel
:daily                       daily notes panel
:links                       links panel for selected note
:today                       open today's daily note
```

### Git

```
:git init <url>              connect vault to git remote
:git status                  show remote and autopush state
:sync                        pull then push
```

### Config

```
:config                      show current config
:config editor nvim          set editor
:config autopush on          push after every save
:config autopush off         manual sync only
:keybinds                    edit ~/.gpad/keybinds.json
```

### Templates

```
:template list
:template new <n>
:template edit <n>
:template delete <n>
```

### Other

```
:help                        full help screen
:q / :quit                   quit
```

---

## Custom keybindings

`~/.gpad/keybinds.json` is created on first run. Open it with `:keybinds`, edit any mapping.

```json
{
  "s": "panel_search",
  "S": "sync",
  "ctrl+n": "new"
}
```

**All action names:**

| Action | Default | Description |
|---|---|---|
| `move_up` | `k` | move cursor up |
| `move_down` | `j` | move cursor down |
| `move_top` | `g` | jump to top |
| `move_bottom` | `G` | jump to bottom |
| `enter` | `Enter` | cd into dir or open note |
| `expand` | `l` `Space` | expand dir in-place |
| `left` | `h` | collapse or go to parent |
| `up_dir` | `-` | go up one directory |
| `open` | `o` | open in editor |
| `view` | `v` | view in pager |
| `new` | `n` | new note |
| `delete` | `d` | delete |
| `rename` | `r` | rename |
| `yank` | `y` | copy path |
| `paste` | `p` | paste |
| `pin` | `P` | toggle pin |
| `sync` | `s` | git sync |
| `today` | `t` | open daily note |
| `command_mode` | `:` | command bar |
| `filter_mode` | `/` | filter tree |
| `help` | `?` | help screen |
| `redraw` | `Ctrl+L` | force redraw |
| `quit` | `q` | quit |
| `panel_search` | `F` | search panel |
| `panel_tags` | `T` | tag browser |
| `panel_graph` | `W` | graph panel |
| `panel_daily` | `D` | daily notes |
| `panel_links` | `L` | links panel |

---

## CLI reference

All commands work standalone for scripts and aliases.

```bash
# Notes
gpad open <note>               open or create
gpad new <note> [-t template]  create with template
gpad view <note>               render in terminal
gpad ls                        tree view

# Daily
gpad today
gpad today yesterday
gpad today list

# Search
gpad find [query]              TUI panel if no query
gpad find -f <query>           fuzzy title only
gpad find -t <query>           full-text only

# Links
gpad links <note>
gpad graph [note]

# Tags
gpad tags [tag]
gpad tag add <tag> <note>
gpad tag rm <tag> <note>

# Files
gpad mv <src> <dest>
gpad cp <src> <dest>
gpad rm [-r] [-y] <note>

# Pin
gpad pin <note>
gpad unpin <note>
gpad pinned

# Templates
gpad template list|new|edit|delete

# Git + config
gpad git init <url>
gpad sync
gpad config editor nvim
gpad config autopush on|off
gpad config show

# Shell completion
source <(gpad completion bash)   # add to ~/.bashrc
source <(gpad completion zsh)    # add to ~/.zshrc
gpad completion fish | source
```

---

## Wikilinks

Write `[[note name]]` in any note body to create a link:

```markdown
This connects to [[physics]] and [[mathematics]].
See also [[college/linear-algebra|Linear Algebra]].
```

| Syntax | Meaning |
|---|---|
| `[[note]]` | resolves to note.md anywhere in vault |
| `[[folder/note]]` | explicit path |
| `[[note\|alias]]` | display alias |
| `[[note#heading]]` | heading anchor |

The graph (`W`), links panel (`L`), and preview footer all update automatically.

---

## Tags

```markdown
---
tags: [go, cli, tools]
---

Also about #programming and #linux.
```

Both frontmatter tags and inline `#tags` are indexed.

---

## Daily notes

```bash
gpad today          # opens notes/daily/2026-03-22.md
```

Press `D` in the TUI for the daily notes panel, or `t` to jump to today.

---

## Templates

Built-in: `note`, `daily`, `meeting`, `idea`. Placeholders:

| | |
|---|---|
| `{{title}}` | note title from filename |
| `{{date}}` | today YYYY-MM-DD |
| `{{time}}` | HH:MM |
| `{{cursor}}` | stripped on write |

---

## Git sync

```bash
gpad git init git@github.com:you/notes.git
```

After connecting, saves push to Git **in the background** — the TUI returns instantly. A brief `↑ pushing...` appears in the status bar and clears automatically. Works with SSH and HTTPS. Fully offline without any Git setup.

---

## Vault layout

```
~/.gpad/
├── config.json
├── keybinds.json
├── index.json
├── templates/
└── notes/
    ├── daily/
    ├── college/
    └── ideas.md
```

---

## License

[MIT](LICENSE)
