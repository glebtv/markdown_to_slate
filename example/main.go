package main

import (
	"encoding/json"
	"log"

	"github.com/glebtv/markdown_to_slate"
)

var Example1 = "Code ```inline``` tag\n" +
	"Code2 block:\n" +
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
var Example6 = "test:\n```ruby\ntest code\n```"

var Example7 = "Code ```inline``` tag\n" +
	"Code block:\n" +
	"```\nblock code\ntest\n```\nafter code"

var Example8 = "Code ```inline``` tag\nline two"

var Example9 = `
- [ ] 1) ссылка http://www.test.ru

- [ ] 2) ссылка

- [ ] 3) ссылки
`

var Example = `http://travelask.ru/admin/travel_comments
Пример ответа, у которого есть коммент http://travelask.ru/questions/1345272-kakoy-samyy-luchshiy-plyazh-v-suhume


Добавляем в админке фильтр "коммент удален". В этот фильтр добавляем комменты, которые получили статус "частично удален". Можно посмотреть тут
https://drive.google.com/file/d/1Kbbc2slKXKIVZSaD5Vxh6ce6SodAv2pm/view?usp=drivesdk`

func Run(input string) {
	data := markdown_to_slate.Parse([]byte(input))
	//s, err := json.Marshal(data)
	s, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		panic(err)
	}
	log.Println(string(s))
}

func main() {
	Run(Example)
}
