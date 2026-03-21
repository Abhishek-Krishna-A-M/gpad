package help

// GetMarkdownGuide returns the full markdown + gpad-syntax reference.
func GetMarkdownGuide() string {
	return `
# Markdown & gpad Syntax Guide

## Headings

# H1 — note title
## H2 — section
### H3 — subsection
#### H4 — detail

## Text Formatting

**bold**  *italic*  ~~strikethrough~~
` + "`inline code`" + `

## Lists

- bullet item
- another item
  - nested item

1. ordered item
2. second item

## Task Lists

- [ ] todo item
- [x] done item

## Code Blocks

` + "```" + `go
fmt.Println("hello, gpad")
` + "```" + `

## Quotes

> important thought

## Horizontal Rule

---

## gpad-Specific Syntax

### Wikilinks

[[note name]]             link to any note by name
[[folder/note]]           explicit path
[[note|display text]]     aliased link
[[note#heading]]          link to heading

### Tags

Add tags in frontmatter:

` + "```" + `yaml
---
title: My Note
date: 2026-01-01
tags: [go, cli, ideas]
---
` + "```" + `

Or inline anywhere in the body:

This is about #programming and #cli tools.

## Commands Quick Reference

` + "```" + `
gpad today              open today's daily note
gpad open <note>        open or create a note
gpad new <note>         create with template
gpad view <note>        render note in terminal
gpad find <query>       full-text + fuzzy search
gpad links <note>       show backlinks & outlinks
gpad tags               list all tags in vault
gpad tags <tag>         notes with this tag
gpad tag add <tag> <note>   add tag to note
gpad pin <note>         pin a note (shows ★ in tree)
gpad graph              ASCII link graph
gpad template list      list templates
gpad template new <name> create a template
gpad ls                 tree view of vault
gpad mv <src> <dest>    move/rename note
gpad cp <src> <dest>    copy note
gpad rm <note>          delete note
gpad sync               pull/push git remote
gpad git init <url>     connect git remote
gpad config editor vim  set preferred editor
gpad config autopush on enable auto-push
` + "```" + `
`
}
