package main

import (
	"encoding/json"
	"log"

	"github.com/glebtv/markdown_to_slate"
)

var Example1 = "Code tag\n" +
	"```\n" +
	`[07/Oct/2018:13:53:57 +0300] "GET /system/image.png HTTP/1.1" 304 0 "-" "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36 OPR/56.0.3051.36` + "\n```"

var Example2 = `numbered list,
1) test *bold* _em_ ~~strike~~
http://example.com/test

2) second item
text after item` + "\n and some ```inline code``` for this"

var Example3 = `[rt1622_regions_sorted.csv](/original/rt1622_regions_sorted.csv) (340.9 KiB)
[rt1622_cities_sorted.csv](/original/rt1622_cities_sorted.csv) (22.6 MiB)`

var Example4 = `36) Тип. Косяк с наверным ID

Как стало пуэбло?

[![image.png](/thumb/image.jpg)](/original/image.png "image.png 228.1 KiB")`

var Example5 = `list test:

first list:

- list is li
- test

second list:

1. list is ol
2. test
`

var Example6 = `checked list:

- [ ] unchecked
- [x] checked
`

func Run(input string) {
	data := markdown_to_slate.Parse([]byte(input))
	//s, err := json.MarshalIndent(data, "", "    ")
	s, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	log.Println(string(s))
}

func main() {
	//Run(Example1)
	//Run(Example2)
	//Run(Example3)
	Run(Example4)
	//Run(Example5)
	//Run(Example6)
}
