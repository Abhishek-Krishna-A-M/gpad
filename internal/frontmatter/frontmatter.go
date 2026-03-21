// Package frontmatter reads and writes YAML frontmatter blocks in markdown files.
//
// Supported format:
//
//	---
//	title: My Note
//	date: 2026-03-22
//	tags: [go, cli, programming]
//	pinned: true
//	---
package frontmatter

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"
)

type Meta struct {
	Title  string
	Date   time.Time
	Tags   []string
	Pinned bool
	Extra  map[string]string
}

// Parse reads frontmatter and body from a file.
func Parse(path string) (Meta, string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Meta{}, "", err
	}
	meta, body := parseBytes(data)
	return meta, body, nil
}

func parseBytes(data []byte) (Meta, string) {
	meta := Meta{Extra: map[string]string{}}
	s := string(data)

	if !strings.HasPrefix(s, "---\n") {
		return meta, s
	}
	rest := s[4:]
	end := strings.Index(rest, "\n---")
	if end == -1 {
		return meta, s
	}
	yamlBlock := rest[:end]
	body := strings.TrimPrefix(rest[end+4:], "\n")
	// consume optional trailing newline after closing ---
	body = strings.TrimPrefix(body, "\n")

	scanner := bufio.NewScanner(strings.NewReader(yamlBlock))
	for scanner.Scan() {
		line := scanner.Text()
		idx := strings.Index(line, ":")
		if idx == -1 {
			continue
		}
		key := strings.TrimSpace(line[:idx])
		val := strings.TrimSpace(line[idx+1:])
		switch key {
		case "title":
			meta.Title = val
		case "date":
			if t, e := time.Parse("2006-01-02", val); e == nil {
				meta.Date = t
			}
		case "tags":
			meta.Tags = parseTags(val)
		case "pinned":
			meta.Pinned = val == "true"
		default:
			meta.Extra[key] = val
		}
	}
	return meta, body
}

func parseTags(val string) []string {
	val = strings.Trim(val, "[] ")
	var tags []string
	for _, t := range strings.Split(val, ",") {
		t = strings.TrimSpace(t)
		t = strings.Trim(t, `"'`)
		t = strings.TrimPrefix(t, "#")
		if t != "" {
			tags = append(tags, t)
		}
	}
	return tags
}

// Write serializes frontmatter + body back to a file.
func Write(path string, meta Meta, body string) error {
	var buf bytes.Buffer
	buf.WriteString("---\n")
	if meta.Title != "" {
		fmt.Fprintf(&buf, "title: %s\n", meta.Title)
	}
	if !meta.Date.IsZero() {
		fmt.Fprintf(&buf, "date: %s\n", meta.Date.Format("2006-01-02"))
	}
	if len(meta.Tags) > 0 {
		fmt.Fprintf(&buf, "tags: [%s]\n", strings.Join(meta.Tags, ", "))
	}
	if meta.Pinned {
		buf.WriteString("pinned: true\n")
	}
	for k, v := range meta.Extra {
		fmt.Fprintf(&buf, "%s: %s\n", k, v)
	}
	buf.WriteString("---\n\n")
	buf.WriteString(body)
	return os.WriteFile(path, buf.Bytes(), 0644)
}

// EnsureFrontmatter adds frontmatter to a file that doesn't have it yet.
func EnsureFrontmatter(path, title string) error {
	meta, body, err := Parse(path)
	if err != nil {
		return err
	}
	if meta.Title != "" {
		return nil // already has frontmatter
	}
	meta.Title = title
	meta.Date = time.Now()
	return Write(path, meta, body)
}

// AddTag adds a tag to a note's frontmatter without duplicating.
func AddTag(path, tag string) error {
	meta, body, err := Parse(path)
	if err != nil {
		return err
	}
	tag = strings.TrimPrefix(strings.ToLower(tag), "#")
	for _, t := range meta.Tags {
		if t == tag {
			return nil
		}
	}
	meta.Tags = append(meta.Tags, tag)
	if meta.Date.IsZero() {
		meta.Date = time.Now()
	}
	return Write(path, meta, body)
}

// RemoveTag removes a tag from a note's frontmatter.
func RemoveTag(path, tag string) error {
	meta, body, err := Parse(path)
	if err != nil {
		return err
	}
	tag = strings.TrimPrefix(strings.ToLower(tag), "#")
	updated := meta.Tags[:0]
	for _, t := range meta.Tags {
		if t != tag {
			updated = append(updated, t)
		}
	}
	meta.Tags = updated
	return Write(path, meta, body)
}

// InlineTags scans body text for #word style inline tags.
func InlineTags(body string) []string {
	seen := map[string]bool{}
	var tags []string
	for _, word := range strings.Fields(body) {
		if strings.HasPrefix(word, "#") && len(word) > 1 {
			tag := strings.ToLower(strings.Trim(word[1:], ".,!?;:\"'()"))
			if tag != "" && !seen[tag] {
				seen[tag] = true
				tags = append(tags, tag)
			}
		}
	}
	return tags
}
