package help

func GetMarkdownGuide() string {
	return `
# Markdown Basics
================

## Headings
# Heading 1
## Heading 2
### Heading 3

## Text Formatting
**bold** *italic* ~~strikethrough~~  
` + "`inline code`" + `

## Lists
- item 1
- item 2
  - nested item

## Code Blocks
` + "```" + `
code here
` + "```" + `

## Quotes
> quoted text
`
}
