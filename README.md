<div align="center">

# gpad 2.0

**A terminal-native knowledge vault for people who live in the shell.**

Wikilinks. Backlinks. Tags. Daily notes. Full-text search. ASCII graphs.  
Git-synced across machines — or fully offline. One static binary. Zero bloat.

</div>

---

## What's new in 2.0

| Feature | Command |
|---|---|
| `[[Wikilinks]]` between notes | auto-resolved, bidirectional |
| Backlinks panel | shown at the bottom of every `view` |
| Tag system (`#tag` + frontmatter) | `gpad tags`, `gpad tag add` |
| Daily notes | `gpad today` |
| Full-text + fuzzy search | `gpad find` |
| Note templates | `gpad new -t meeting` |
| Pinned notes (★ in tree) | `gpad pin` |
| ASCII link graph | `gpad graph` |
| Word count + stats in viewer | shown on every `view` |
| Frontmatter on every note | title, date, tags — auto-added |

---

## Installation

**Linux / macOS (pre-built binary):**

```bash
curl -fsSL https://raw.githubusercontent.com/Abhishek-Krishna-A-M/gpad/main/install.sh | sh
```

**Build from source (requires Go 1.22+):**

```bash
git clone https://github.com/Abhishek-Krishna-A-M/gpad
cd gpad
go build -o gpad ./cmd/gpad/
sudo mv gpad /usr/local/bin/
```

**Windows (PowerShell):**

```powershell
powershell -ExecutionPolicy Bypass -File install.ps1
```

---

## Quick start

```bash
# gpad works immediately with no setup — notes live in ~/.gpad/notes/
gpad today                        # open today's daily note
gpad open ideas/my-first-note.md  # open or create any note
gpad ls                           # browse your vault as a tree
```

**Connect git sync (optional):**

```bash
gpad git init git@github.com:you/notes.git
# every save now pushes automatically
```

---

## Full command reference

### Notes

```bash
gpad open <note>              # open or create a note (syncs before + after)
gpad new <note>               # create with template picker
gpad new <note> -t meeting    # create with a specific template
gpad view <note>              # render markdown + backlinks + stats in terminal
gpad ls                       # tree view of vault (★ = pinned)
```

### Daily notes

```bash
gpad today                    # open today's note (creates if missing)
gpad today yesterday          # open yesterday's note
gpad today list               # show last 14 daily notes
```

### Search

```bash
gpad find <query>             # full-text body search + fuzzy title match
gpad find -f <query>          # fuzzy title only
gpad find -t <query>          # full-text body only
```

### Wikilinks & graph

```bash
gpad links <note>             # show backlinks and outlinks for a note
gpad graph                    # ASCII graph of the full vault
gpad graph <note>             # ego graph: note + immediate neighbours
```

In any note, use `[[note name]]` to link. Aliases and anchors work too:

```markdown
See [[college/math]] for the proof.
This relates to [[quantum-foam|quantum mechanics]].
The key equation is in [[physics#energy]].
```

### Tags

```bash
gpad tags                     # full tag index with counts
gpad tags <tag>               # list notes with this tag
gpad tag add <tag> <note>     # add tag to frontmatter
gpad tag rm <tag> <note>      # remove tag from frontmatter
```

Tags can live in frontmatter or inline — both are indexed:

```markdown
---
tags: [go, cli, tools]
---

This is also #programming related.
```

### Pinned notes

```bash
gpad pin <note>               # pin a note (shows ★ in ls)
gpad unpin <note>             # unpin
gpad pinned                   # list all pinned notes
```

### File management

```bash
gpad mv <src> <dest>          # move or rename (updates H1 title too)
gpad cp <src> <dest>          # copy a note
gpad rm <note>                # delete (prompts for confirmation)
gpad rm -r <folder>           # delete folder recursively
gpad rm -y <note>             # delete without prompt
```

### Templates

Four built-in templates are seeded on first run: `note`, `daily`, `meeting`, `idea`.

```bash
gpad template list            # list available templates
gpad template new <n>         # create a new template (opens in editor)
gpad template edit <n>        # edit existing template
gpad template delete <n>      # delete a template
```

Template placeholders:

```
{{title}}   → note title (derived from filename)
{{date}}    → today's date (YYYY-MM-DD)
{{time}}    → current time (HH:MM)
{{cursor}}  → stripped on write (marks where cursor lands)
```

### Git sync

```bash
gpad git init <url>           # connect remote (SSH or HTTPS)
gpad git status               # show remote + autopush state
gpad sync                     # manual pull + push
gpad config autopush on       # push automatically on every save
gpad config autopush off      # manual sync only
```

gpad works fully **offline** with no git setup. Add git any time.

### Configuration

```bash
gpad config editor nvim       # set preferred editor
gpad config autopush on       # enable auto-push
gpad config show              # print current config
```

Priority order for editor: `config.json` → `$EDITOR` → `$VISUAL` → nvim/vim/micro/nano.

### Shell completion

```bash
# Zsh
source <(gpad completion zsh)

# Bash
source <(gpad completion bash)

# Fish
gpad completion fish | source
```

### Syntax guide

```bash
gpad markdown                 # full markdown + gpad syntax reference
```

---

## Vault layout

```
~/.gpad/
├── config.json               # editor, git settings, pinned list
├── index.json                # link/tag index cache (auto-built)
├── templates/
│   ├── note.md
│   ├── daily.md
│   ├── meeting.md
│   └── idea.md
└── notes/
    ├── daily/
    │   ├── 2026-03-22.md
    │   └── 2026-03-21.md
    ├── college/
    │   └── math.md
    └── ideas.md
```

---

## Frontmatter

Every note gets a frontmatter block (auto-added on first open if absent):

```yaml
---
title: Quantum Foam
date: 2026-03-22
tags: [physics, ideas]
pinned: false
---
```

---

## Project structure

```
.
├── cmd/gpad/main.go          # entry point
└── internal/
    ├── cmd/                  # CLI (Cobra commands)
    │   ├── root.go           # root command + version
    │   ├── open.go           # gpad open
    │   ├── new.go            # gpad new
    │   ├── today.go          # gpad today
    │   ├── find.go           # gpad find
    │   ├── links.go          # gpad links
    │   ├── tags.go           # gpad tags / tag
    │   ├── pin.go            # gpad pin / unpin / pinned
    │   ├── graph.go          # gpad graph
    │   ├── template.go       # gpad template
    │   ├── view.go           # gpad view
    │   ├── ls.go             # gpad ls
    │   ├── git.go            # gpad git
    │   ├── sync.go           # gpad sync
    │   ├── config.go         # gpad config
    │   └── misc.go           # mv, cp, rm, completion, markdown
    ├── config/               # config.json load/save, pin list
    ├── core/                 # move, delete, copy, sync logic
    ├── daily/                # daily note open/create/list
    ├── editor/               # editor detection + open
    ├── frontmatter/          # YAML frontmatter parse/write
    ├── gitrepo/              # git init, add/commit/push, merge
    ├── help/                 # markdown syntax guide text
    ├── links/                # [[wikilink]] parse, link graph
    ├── notes/                # note open/create/stats, tree list
    ├── search/               # full-text + fuzzy search
    ├── storage/              # path helpers (~/.gpad/...)
    ├── tags/                 # tag index across vault
    ├── templates/            # template apply/save/list
    ├── ui/                   # note path listing for completion
    └── viewer/               # ANSI markdown renderer + pager
```

---

## Contributing

Issues and PRs welcome. The architecture is modular — each package has one job and no circular dependencies. Core packages (`frontmatter`, `links`, `tags`, `search`, `storage`) never import CLI packages.

---

## License

[MIT](LICENSE)
