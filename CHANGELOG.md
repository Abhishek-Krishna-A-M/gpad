# Changelog

## [2.0.0] — 2026-03-22

### New features

- **`[[Wikilinks]]`** — link any note to any other by name. Aliases (`[[note|text]]`) and heading anchors (`[[note#section]]`) supported. Links are resolved case-insensitively across the full vault.
- **Backlinks** — every `gpad view` shows an outlinks + backlinks panel at the bottom. `gpad links <note>` shows the full link graph for any note.
- **`gpad graph`** — ASCII graph of the entire vault link structure. `gpad graph <note>` shows an ego graph (note + immediate neighbours).
- **Tag system** — tags in YAML frontmatter (`tags: [go, cli]`) and inline `#hashtags` in body text are both indexed. `gpad tags` shows the full tag index. `gpad tags <tag>` lists notes. `gpad tag add/rm` manages frontmatter tags.
- **Daily notes** — `gpad today` opens today's note (`notes/daily/YYYY-MM-DD.md`), creating it from the `daily` template if needed. `gpad today yesterday` and `gpad today list` also available.
- **Templates** — `gpad new <note> -t <template>` creates from a named template. Four built-ins seeded on first run: `note`, `daily`, `meeting`, `idea`. `gpad template new/edit/delete/list` for management. Placeholders: `{{title}}`, `{{date}}`, `{{time}}`, `{{cursor}}`.
- **Full-text search** — `gpad find <query>` searches note bodies with match-count scoring and an inline excerpt. `-f` for fuzzy title-only, `-t` for body-only.
- **Fuzzy search** — subsequence matching on note paths with consecutive-match and prefix bonuses.
- **Pinned notes** — `gpad pin <note>` marks a note with ★ in `gpad ls`. Stored in `config.json`. `gpad pinned` lists all pinned notes.
- **Word count + stats** — every `gpad view` shows word count, line count, and link count in a footer bar.
- **Frontmatter on every note** — `EnsureFrontmatter` auto-adds `title`, `date`, `tags` to notes that don't have it yet, on first open.
- **`gpad config show`** — print the full current configuration.
- **`gpad git status`** — show connected remote and autopush state.
- **`gpad today yesterday`** / **`gpad today list`** — navigate recent daily notes.
- **`gpad new`** — dedicated create command separate from `open`, with template flag.
- **`gpad markdown`** — updated syntax guide covers all 2.0 features.

### Improvements

- `gpad view` now renders wikilinks as `→ note` in cyan and inline `#tags` in yellow.
- `gpad view` renders frontmatter as a styled metadata block instead of raw YAML.
- `gpad view` renders task lists (`- [ ]` / `- [x]`) with ○ / ✓ symbols.
- `gpad view` renders nested lists, ordered lists, H4.
- `gpad ls` tree shows pinned notes with ★ and uses bold/colour for directories.
- `gpad mv` updates frontmatter `title` field in addition to the H1 header.
- `gpad rm` prompts for confirmation by default (`-y` to skip).
- `gpad git init` is more robust: idempotent, handles empty remotes, creates `.gitignore`.
- `gpad open` / `gpad new` pull before editing and push after (when autopush is on).
- Editor detection order: `config.json` → `$EDITOR` → `$VISUAL` → nvim/vim/micro/nano.
- All new packages have zero external dependencies beyond the standard library.

### Architecture

- New packages: `frontmatter`, `links`, `tags`, `search`, `daily`, `templates`.
- Sync logic moved from `notes` package up to CLI layer — core packages are pure.
- No circular dependencies. Core packages never import CLI packages.
- `storage` package gains `DailyDir()`, `TemplatesDir()`, `IndexPath()`.
- `config` package gains `Pin`, `Unpin`, `IsPinned` and `Pinned []string` field.

---

## [1.0.0] — initial release

- Structured storage in `~/.gpad/notes`
- Create/edit notes with any editor
- ANSI markdown rendering with `less` pager
- `gpad sync` — git pull/push
- Auto-push on save
- `gpad mv`, `gpad cp`, `gpad rm`
- Shell completion (bash, zsh, fish, powershell)
