package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/glebtv/markdown_to_slate"
)

var Example = "Code ```inline``` tag\n" +
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

var Example10 = `http://test.ru
Пример ответа, у которого есть коммент http://test.ru


Добавляем в админке фильтр
https://drive.google.com/
`

var Example11 = `test issue: #5 #6

test mention: @test @gleb
`

var Example12 = `
<!-- Yandex.Metrika counter -->
<script type="text/javascript">
    (function (d, w, c) {
        (w[c] = w[c] || []).push(function() {
            try {
                w.yaCounter25757501 = new Ya.Metrika({
                    id:25757501,
                    clickmap:true,
                    trackLinks:true,
                    accurateTrackBounce:true
                });
            } catch(e) { }
        });

        var n = d.getElementsByTagName("script")[0],
            s = d.createElement("script"),
            f = function () { n.parentNode.insertBefore(s, n); };
        s.type = "text/javascript";
        s.async = true;
        s.src = "https://mc.yandex.ru/metrika/watch.js";

        if (w.opera == "[object Opera]") {
            d.addEventListener("DOMContentLoaded", f, false);
        } else { f(); }
    })(document, window, "yandex_metrika_callbacks");
</script>
<noscript><div><img src="https://mc.yandex.ru/watch/25757501" style="position:absolute; left:-9999px;" alt="" /></div></noscript>
<!-- /Yandex.Metrika counter -->
`

var Example13 = `Вот тут в макете дбавили рамку для области ввода текста
https://drive.google.com/file/d/fileid/view?usp=drivesdk

Макет
https://www.figma.com/file/fileid?node-id=0%3A1`

func AllowHeader(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Next()
}

type MD2SlateParams struct {
	Body string `json:"body"`
}

func InitServer() {
	r := gin.Default()
	r.Use(cors.Default())
	r.Use(AllowHeader)
	r.Use(static.Serve("/files", static.LocalFile("../files", false)))
	r.GET("/examples", func(c *gin.Context) {
		matches, err := filepath.Glob("../files/*/*.md")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		for i, match := range matches {
			matches[i] = "/" + strings.Join(strings.Split(match, "/")[1:], "/")
		}
		c.JSON(http.StatusOK, matches)
	})
	r.POST("/md2slate", func(c *gin.Context) {
		var params MD2SlateParams
		if err := c.ShouldBindJSON(&params); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
		}
		slate := markdown_to_slate.Parse([]byte(params.Body))
		s, err := json.MarshalIndent(slate, "", "    ")
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}
		c.Data(http.StatusOK, "application/json", s)
	})
	r.Run()
}
func Run(input string) {
	data := markdown_to_slate.Parse([]byte(input))
	//s, err := json.Marshal(data)
	s, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(s))
}

func main() {
	//println(Example)
	InitServer()
}
