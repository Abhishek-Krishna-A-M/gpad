<div align="center">
  
  # gpad
  
> A fast, cross-platform, Git-powered CLI notes manager written in Go.

 📝 Take notes from anywhere.
 🔄 Sync across devices using GitHub (SSH or HTTPS).
 🖥️ Render Markdown directly in your terminal.
 ⚡ Zero bloat. One static binary. Works everywhere.
</div>
---

## Features

- **Structured Storage** — Notes are stored as Markdown inside `~/.gpad/notes`.
- **Flexible Editing** — Create and edit notes using your preferred editor (`nvim`, `code`, `micro`, etc.).
- **Terminal Viewer** — Clean headings, lists, and code blocks rendered with ANSI colors.
- **Smart Pager** — Automatically uses `less -R` so you can scroll through long notes.
- **Auto Sync** — Optional background `git add → commit → push` after every edit.
- **Manual Sync** — Dedicated `gpad sync` command for manual pulls and pushes.
- **Safety First** — Deletion is restricted to the notes directory to prevent accidental data loss.

---

## Installation

**Linux / macOS:**

```bash
curl -fsSL https://raw.githubusercontent.com/Abhishek-Krishna-A-M/gpad/main/install.sh | sh
```

**Windows (PowerShell):**

```powershell
powershell -ExecutionPolicy Bypass -File install.ps1
```

---

## Usage

### 1. Initialize gpad

gpad auto-initializes for offline use. For Git sync, run:

```bash
# GitHub sync mode (SSH recommended)
gpad git init git@github.com:user/notes.git
```

### 2. Create or Edit a Note

One command handles both — if the file doesn't exist, gpad creates it for you:

```bash
gpad open college/math.md
```

> Running `gpad open` without arguments lists your notes and folders interactively.

### 3. View a Note

View rendered Markdown in the terminal. After viewing, gpad will ask if you'd like to open the file for editing:

```bash
gpad view college/math.md
```

### 4. Git Syncing

Manage your repository manually or toggle automation:

```bash
# Manual Pull/Push
gpad sync

# Toggle Auto-Push on/off
gpad config autopush on
```

### 5. Management

```bash
# Delete a note (with safety checks)
gpad rm college/math.md

# Show Markdown syntax guide
gpad help markdown
```

---

## Configuration

| Command | Description |
|---|---|
| `gpad config editor nvim` | Set your preferred editor |
| `gpad config autopush on` | Enable automatic Git push after edits |

---

## Contributing

Issues and PRs are welcome! This project follows a modular Go structure:

- **`internal/cmd`** — CLI command definitions (Cobra).
- **`internal/core`** — Business logic (Sync, Git operations).
- **`internal/viewer`** — Markdown rendering logic.

---

## License

[MIT License](LICENSE)
