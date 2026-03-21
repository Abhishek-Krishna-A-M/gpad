// Package templates manages note templates stored in ~/.gpad/templates/.
//
// A template is a plain .md file with optional {{placeholders}}:
//
//	{{title}}   → replaced with the note title
//	{{date}}    → replaced with today's date (YYYY-MM-DD)
//	{{time}}    → replaced with current time (HH:MM)
//	{{cursor}}  → marks where the cursor should land (stripped on write)
package templates

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Abhishek-Krishna-A-M/gpad/internal/storage"
)

// List returns all template names (without .md extension).
func List() []string {
	dir := storage.TemplatesDir()
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".md") {
			names = append(names, strings.TrimSuffix(e.Name(), ".md"))
		}
	}
	return names
}

// Apply reads a template and returns the rendered content for a new note.
func Apply(templateName, noteTitle string) (string, error) {
	path := filepath.Join(storage.TemplatesDir(), templateName+".md")
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("template %q not found (run 'gpad template list')", templateName)
	}

	now := time.Now()
	content := string(data)
	content = strings.ReplaceAll(content, "{{title}}", noteTitle)
	content = strings.ReplaceAll(content, "{{date}}", now.Format("2006-01-02"))
	content = strings.ReplaceAll(content, "{{time}}", now.Format("15:04"))
	content = strings.ReplaceAll(content, "{{cursor}}", "") // strip marker
	return content, nil
}

// Save writes a new template file.
func Save(name, content string) error {
	if !strings.HasSuffix(name, ".md") {
		name += ".md"
	}
	path := filepath.Join(storage.TemplatesDir(), name)
	return os.WriteFile(path, []byte(content), 0644)
}

// Delete removes a template.
func Delete(name string) error {
	path := filepath.Join(storage.TemplatesDir(), name+".md")
	return os.Remove(path)
}

// Path returns the absolute path of a template file.
func Path(name string) string {
	return filepath.Join(storage.TemplatesDir(), name+".md")
}

// EnsureDefaults seeds built-in templates if the templates dir is empty.
func EnsureDefaults() {
	dir := storage.TemplatesDir()
	entries, _ := os.ReadDir(dir)
	if len(entries) > 0 {
		return
	}

	defaults := map[string]string{
		"note": `---
title: {{title}}
date: {{date}}
tags: []
---

# {{title}}

{{cursor}}
`,
		"daily": `---
title: {{date}}
date: {{date}}
tags: [daily]
---

# {{date}}

## Tasks

- [ ] 

## Notes

## Done

`,
		"meeting": `---
title: {{title}}
date: {{date}}
tags: [meeting]
---

# {{title}}

**Date:** {{date}}  
**Attendees:**  

## Agenda

- 

## Notes

## Action items

- [ ] 
`,
		"idea": `---
title: {{title}}
date: {{date}}
tags: [idea]
---

# {{title}}

## The idea

## Why it matters

## Next steps

- [ ] 
`,
	}

	for name, content := range defaults {
		_ = Save(name, content)
	}
}
