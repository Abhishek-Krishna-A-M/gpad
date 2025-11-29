package help

import "fmt"

func Markdown() {
	fmt.Println(`
Markdown Basics
===============

Headings
--------
# Heading 1
## Heading 2
### Heading 3


Text Formatting
---------------
**bold**  
*italic*  
~~strikethrough~~  
` + "`inline code`" + `


Lists
-----
- item 1
- item 2
  - nested item


Code Blocks
-----------
Use triple backticks in real Markdown:

~~~
code here
~~~


Links & Images
--------------
[OpenAI](https://openai.com)  
![Image alt](image.png)


Quotes
------
> quoted text


Tables
------
| Name | Age |
|------|-----|
| Ana  | 22  |
| Ben  | 30  |


Horizontal Rule
---------------
---


Escaping Characters
-------------------
Use backslash before characters:
\*  \_  \#  \` + "`" + `
`)
}
