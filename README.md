<div align="center">

# gpad

**A fast, cross-platform, Git-powered CLI notes manager written in Go.**

Take notes from anywhere.  
Sync across devices using GitHub (SSH or HTTPS).  
Render Markdown in the terminal.  
Zero bloat. One static binary. Works everywhere.

</div>

---

## Features

- Notes stored as Markdown inside a folder (`~/.gpad/notes`)
- Tree view of notes with nested directories (`ideas/ai.md`, etc.)
- Create & edit notes using your preferred editor (nvim, code, micro…)
- Terminal Markdown viewer (clean headings, lists, quotes, blocks)
- Offline mode OR GitHub sync mode
- Auto sync (git add → commit → push) after editing
- Manual sync (`gpad sync`, `gpad sync log`)
- HTTPS or SSH Git support (with auto-switch to SSH)
- Uninstall command
- STDIN viewer support  
  ```sh
  echo "# Hi" | gpad view -
  ```

Built for speed, simplicity, and developer workflows.

---

## Installation

### Linux / macOS

Install with the official script:

```sh
curl -fsSL https://raw.githubusercontent.com/Abhishek-Krishna-A-M/gpad/main/install.sh | sh
```

Or install manually:

```sh
go install github.com/Abhishek-Krishna-A-M/gpad/cmd/gpad@latest
```

### Windows (PowerShell)

```powershell
powershell -ExecutionPolicy Bypass -File install.ps1
```

Or download the `.exe` from [Releases](https://github.com/Abhishek-Krishna-A-M/gpad/releases) and add it to PATH.

---

## Usage

### Initialize gpad

**Offline mode (default)**

```sh
gpad init
```

**GitHub sync mode (HTTPS or SSH)**

```sh
gpad init --github https://github.com/user/notes.git
```

Recommended (no passwords required):

```sh
gpad init --github git@github.com:user/notes.git
```

If HTTPS is used, gpad will show credential helper instructions.

### Create or Edit Note

One command handles both:

```sh
gpad ideas/ai.md
```

- If the file exists → edit it.
- If not → create it.

### View Notes

```sh
gpad view ideas/ai.md
```

Or view piped Markdown:

```sh
echo "# Hi" | gpad view -
```

### List Notes (Tree View)

```sh
gpad list
```

Example output:

```
notes/
├── ideas/
│   └── ai.md
└── test.md
```

`.git/` is automatically hidden.

### Git Sync

**Sync now (pull + push)**

```sh
gpad sync
```

**Show recent sync logs**

```sh
gpad sync log
```

---

## Configuration

### Set default editor

```sh
gpad config editor nvim
gpad config editor "code -w"
```

### Toggle auto-push

```sh
gpad config autopush on
gpad config autopush off
```

---

## Uninstall

```sh
gpad uninstall
```

Removes all gpad data (`~/.gpad`) but not the binary.

---

## Markdown Help

```sh
gpad help markdown
```

Shows a beginner-friendly Markdown reference.

---

## Architecture Overview

```
cmd/gpad/           → main CLI entrypoint
internal/cli/       → argument parsing + commands
internal/notes/     → create/edit/list logic
internal/viewer/    → Markdown terminal renderer
internal/storage/   → file system paths
internal/gitx/      → git clone/pull/push/merge/SSH
internal/editor/    → editor launcher
internal/help/      → Markdown help text
internal/config/    → JSON config loader/saver
```

---

## Roadmap (v2)

- Theme system (colors & styles)
- Rclone cloud backend (Google Drive, S3, etc.)
- Full-text search (ripgrep integration)
- TUI viewer (optional)
- Note templates
- Image/PDF handling

---

## Contributing

Issues & PRs welcome.  
Please follow Go best practices and minimize dependencies.

---

## License

MIT License.

---

**gpad** — minimal, fast, developer-friendly notes.
