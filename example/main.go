package main

import (
	"encoding/json"
	"log"

	"github.com/glebtv/markdown_to_slate"
)

var Example1 = "Code ```inline``` tag\n" +
	"Code block:\n" +
	"```\nblock code\ntest\n```\n" +
	`numbered list,
1) test *em* _em_ **bold** ~~strike~~
http://example.com/test

2) second item
text after item` + "\n and some ```inline code``` for this" +
	`
[rt1622_regions_sorted.csv](/original/rt1622_regions_sorted.csv) (340.9 KiB)
[rt1622_cities_sorted.csv](/original/rt1622_cities_sorted.csv) (22.6 MiB)

[![image.png](/thumb/image.jpg)](/original/image.png "image.png 228.1 KiB")

list ul:

- list is li
- test

list ol:

1. list is ol
2. test

checked list:

- [ ] unchecked
- [x] checked

## h2

### h3

paragraph

` + "```ruby\ntest code\n```"

var Example2 = "## test\n### test 3"
var Example3 = "**test** *test*"
var Example4 = `
- [ ] unchecked
- [x] checked
`

var Example5 = `
[![image.png](/thumb/image.jpg)](/original/image.png "image.png 228.1 KiB")
`
var Example = "test:\n```ruby\ntest code\n```"

func Run(input string) {
	data := markdown_to_slate.Parse([]byte(input))
	//s, err := json.MarshalIndent(data, "", "    ")
	s, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		panic(err)
	}
	log.Println(string(s))
}

func main() {
	Run(Example)
}
