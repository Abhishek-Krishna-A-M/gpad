# Changelog

All notable changes to **gpad** will be documented here.

This project follows [Semantic Versioning](https://semver.org/).

## [v1.5.0] - 2026-02-28
### Added
- New Modular Architecture: Separated CLI, Core logic, and Git plumbing.
- Improved `internal/viewer`: Enhanced ANSI rendering for headers and code blocks.
- `internal/ui`: Centralized note discovery for faster searching.
- Enhanced Error Handling: Better reporting for file operations and Git sync.

### Changed
- Refactored `internal/cli` into `internal/cmd` for better Cobra integration.
- Optimized background sync logic in `internal/core`.

---

## [v1.0.1] - 2025-12-01
### Added
- `gpad mv` rename command
- Automatic heading update when renaming files
- Safer sync with `git pull --no-rebase`

### Fixed
- Push failures when remote was ahead
- Cleaner delete error messages

---

## [v1.0.0] - 2025-11-30
### Initial Release
- Create/edit Markdown notes
- Viewer with Markdown formatting
- GitHub sync (SSH/HTTPS)
- Tree view
- Auto-push
- Configurable editor
- Markdown help viewer
- Uninstall command

