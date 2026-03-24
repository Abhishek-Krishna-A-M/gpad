package tui

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/Abhishek-Krishna-A-M/gpad/internal/storage"
)

// Action is a named TUI command that can be bound to any key.
type Action string

const (
	ActionQuit           Action = "quit"
	ActionMoveUp         Action = "move_up"
	ActionMoveDown       Action = "move_down"
	ActionMoveTop        Action = "move_top"
	ActionMoveBottom     Action = "move_bottom"
	ActionEnter          Action = "enter"        // cd into dir or open file
	ActionExpand         Action = "expand"        // expand dir in-place (Space)
	ActionLeft           Action = "left"          // collapse / go to parent
	ActionUp             Action = "up_dir"        // go up one directory (-)
	ActionOpen           Action = "open"          // open in editor
	ActionView           Action = "view"          // view rendered in pager
	ActionNew            Action = "new"           // new note in current dir
	ActionDelete         Action = "delete"        // delete selected
	ActionRename         Action = "rename"        // rename (prefill mv cmd)
	ActionYank           Action = "yank"          // copy path
	ActionPaste          Action = "paste"         // paste yanked
	ActionPin            Action = "pin"           // toggle pin
	ActionSync           Action = "sync"          // git sync
	ActionToday          Action = "today"         // open today's daily note
	ActionCommandMode    Action = "command_mode"  // enter : command mode
	ActionFilterMode     Action = "filter_mode"   // enter / filter mode
	ActionHelp           Action = "help"          // show help overlay
	ActionRedraw         Action = "redraw"        // force full redraw
	ActionPanelSearch    Action = "panel_search"  // open search panel
	ActionPanelTags      Action = "panel_tags"    // open tag browser
	ActionPanelGraph     Action = "panel_graph"   // open graph view
	ActionPanelDaily     Action = "panel_daily"   // open daily notes list
	ActionPanelLinks     Action = "panel_links"   // show links for selected note
)

// KeyMap maps key strings to Actions.
// Key strings: single chars ("j"), ctrl combos ("ctrl+l"), special ("enter", "space").
type KeyMap map[string]Action

// defaultKeyMap returns the built-in vim-style keybindings.
// Users can override any binding in ~/.gpad/keybinds.json.
func defaultKeyMap() KeyMap {
	return KeyMap{
		// navigation
		"j":       ActionMoveDown,
		"k":       ActionMoveUp,
		"down":    ActionMoveDown,
		"up":      ActionMoveUp,
		"g":       ActionMoveTop,
		"G":       ActionMoveBottom,
		"enter":   ActionEnter,
		"l":       ActionExpand,
		"right":   ActionExpand,
		"h":       ActionLeft,
		"left":    ActionLeft,
		"-":       ActionUp,
		"space":   ActionExpand,

		// file operations
		"o":      ActionOpen,
		"v":      ActionView,
		"n":      ActionNew,
		"d":      ActionDelete,
		"x":      ActionDelete,
		"r":      ActionRename,
		"y":      ActionYank,
		"p":      ActionPaste,
		"P":      ActionPin,
		"s":      ActionSync,
		"t":      ActionToday,

		// modes
		":":      ActionCommandMode,
		"/":      ActionFilterMode,
		"?":      ActionHelp,
		"ctrl+l": ActionRedraw,
		"q":      ActionQuit,
		"ctrl+c": ActionQuit,
		"esc":    ActionQuit,

		// panels — capital letters for panel switching (fast one-key access)
		"F":      ActionPanelSearch,
		"T":      ActionPanelTags,
		"W":      ActionPanelGraph,
		"D":      ActionPanelDaily,
		"L":      ActionPanelLinks,
	}
}

// keybindsPath returns the path to the user keybinds file.
func keybindsPath() string {
	return filepath.Join(storage.GpadDir(), "keybinds.json")
}

// LoadKeyMap loads user keybinds from ~/.gpad/keybinds.json and merges
// them over the defaults. Unknown keys in the file are ignored.
// Missing file → returns defaults silently.
func LoadKeyMap() KeyMap {
	km := defaultKeyMap()

	data, err := os.ReadFile(keybindsPath())
	if err != nil {
		return km // file doesn't exist yet — defaults only
	}

	var user KeyMap
	if err := json.Unmarshal(data, &user); err != nil {
		return km // malformed JSON — fall back to defaults
	}

	for key, action := range user {
		km[key] = action
	}
	return km
}

// SaveDefaults writes the default keybinds file if it doesn't exist yet,
// so users can open it in their editor and customise.
func SaveDefaultKeybinds() {
	path := keybindsPath()
	if _, err := os.Stat(path); err == nil {
		return // already exists
	}

	km := defaultKeyMap()
	b, err := json.MarshalIndent(km, "", "  ")
	if err != nil {
		return
	}

	header := []byte("// gpad keybinds — edit this file to remap any action\n" +
		"// Reload: restart gpad or press ctrl+l\n" +
		"// Actions: quit move_up move_down enter expand left up_dir open view\n" +
		"//          new delete rename yank paste pin sync today command_mode\n" +
		"//          filter_mode help redraw panel_search panel_tags panel_graph\n" +
		"//          panel_daily panel_links\n")

	_ = os.WriteFile(path, append(header, b...), 0644)
}

// rawKeyToString converts a raw byte sequence to a key string for KeyMap lookup.
func rawKeyToString(b []byte) string {
	if len(b) == 0 {
		return ""
	}

	// ANSI escape sequences (arrow keys)
	if b[0] == 0x1b && len(b) >= 3 && b[1] == '[' {
		switch b[2] {
		case 'A':
			return "up"
		case 'B':
			return "down"
		case 'C':
			return "right"
		case 'D':
			return "left"
		}
	}

	// Special single-byte keys
	switch b[0] {
	case 0x1b:
		return "esc"
	case '\r', '\n':
		return "enter"
	case ' ':
		return "space"
	case 0x7f, 0x08:
		return "backspace"
	case 0x09:
		return "tab"
	case 0x03:
		return "ctrl+c"
	case 0x0c:
		return "ctrl+l"
	case 0x15:
		return "ctrl+u"
	case 0x17:
		return "ctrl+w"
	}

	// Printable ASCII
	if len(b) == 1 && b[0] >= 0x20 && b[0] < 0x7f {
		return string(b[0])
	}

	return ""
}
