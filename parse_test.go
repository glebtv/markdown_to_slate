package markdown_to_slate

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tidwall/pretty"
)

func Run(t *testing.T, input string) string {
	data := Parse([]byte(input))
	s, err := json.Marshal(data)
	t.Log(string(s))
	if err != nil {
		t.Fatal(err)
	}
	return string(s)
}

func RunReverse(t *testing.T, input string) string {
	data := Parse([]byte(input))
	s, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	st, err := Stringify(s)
	if err != nil {
		panic(err)
	}
	return st
}

func printJSON(j string) string {
	return string(pretty.Color(pretty.Pretty([]byte(j)), nil))
}

const paragraphExample = "test paragraph\n"
const paragraphResult = `[{"object":"block","type":"paragraph","nodes":[{"object":"text","leaves":[{"object":"leaf","text":"test paragraph","marks":[]}]}]}]`

func TestParagraph(t *testing.T) {
	s := Run(t, paragraphExample)

	assert.Equal(t, s, paragraphResult, "")
}

func TestParagraphReverse(t *testing.T) {
	s := RunReverse(t, paragraphExample)

	assert.Equal(t, s, paragraphExample, "")
}

const listExample = "- list\n- items"
const listResult = `[{"object":"block","type":"bulleted-list","nodes":[{"object":"block","type":"list-item","nodes":[{"object":"text","leaves":[{"object":"leaf","text":"list","marks":[]}]}]},{"object":"block","type":"list-item","nodes":[{"object":"text","leaves":[{"object":"leaf","text":"items","marks":[]}]}]}]}]`

func TestList(t *testing.T) {
	s := Run(t, listExample)

	assert.Equal(t, s, listResult, "")
}

const checklistResult = `[{"object":"block","type":"paragraph","nodes":[{"object":"text","leaves":[{"object":"leaf","text":"check list:","marks":[]}]}]},{"object":"block","type":"paragraph","nodes":[{"object":"block","type":"check-list-item","data":{"checked":false},"nodes":[{"object":"text","leaves":[{"object":"leaf","text":"unchecked","marks":[]}]}]},{"object":"block","type":"check-list-item","data":{"checked":true},"nodes":[{"object":"text","leaves":[{"object":"leaf","text":"checked","marks":[]}]}]}]}]`

func TestChecklist(t *testing.T) {
	s := Run(t, "check list:\n- [ ] unchecked\n- [x] checked")

	assert.Equal(t, s, checklistResult, "")
}

const codeResult = `[{"object":"block","type":"paragraph","nodes":[{"object":"text","leaves":[{"object":"leaf","text":"code ","marks":[]}]},{"object":"text","leaves":[{"object":"leaf","text":"inline","marks":[{"object":"mark","type":"code"}]}]},{"object":"text","leaves":[{"object":"leaf","text":":","marks":[]}]}]},{"object":"block","type":"code","nodes":[{"object":"block","type":"code_line","nodes":[{"object":"text","leaves":[{"object":"leaf","text":"\nblock code\n","marks":[]}]}]}]}]`

func TestCode(t *testing.T) {
	s := Run(t, "code ```inline```:\n```\nblock code\n```\n")

	assert.Equal(t, s, codeResult, "")
}

const linkExample = `[rt1622_regions_sorted.csv](/original/rt1622_regions_sorted.csv) (340.9 KiB)
[rt1622_cities_sorted.csv](/original/rt1622_cities_sorted.csv) (22.6 MiB)`
const linkResult = `[{"object":"block","type":"paragraph","nodes":[{"object":"inline","type":"link","data":{"href":"/original/rt1622_regions_sorted.csv","title":""},"nodes":[{"object":"text","leaves":[{"object":"leaf","text":"rt1622_regions_sorted.csv","marks":[]}]}]},{"object":"text","leaves":[{"object":"leaf","text":" (340.9 KiB)","marks":[]}]}]},{"object":"block","type":"paragraph","nodes":[{"object":"inline","type":"link","data":{"href":"/original/rt1622_cities_sorted.csv","title":""},"nodes":[{"object":"text","leaves":[{"object":"leaf","text":"rt1622_cities_sorted.csv","marks":[]}]}]},{"object":"text","leaves":[{"object":"leaf","text":" (22.6 MiB)","marks":[]}]}]}]`

func TestLink(t *testing.T) {
	s := Run(t, linkExample)

	assert.Equal(t, s, linkResult, "")
}

const imageExample = `[![i4.png](/system/i4.jpg)](/system/i4.png?1451419607 "i4.png 22.6 KiB")`
const imageResult = `[{"object":"block","type":"paragraph","nodes":[{"object":"block","type":"image","data":{"src":"/system/i4.jpg","title":""}}]}]`

func TestImage(t *testing.T) {
	s := Run(t, imageExample)

	assert.Equal(t, s, imageResult, "")
}

func TestFiles(t *testing.T) {
	matches, err := filepath.Glob("./files/*/*.md")
	if err != nil {
		t.Fatal(err)
	}
	for _, fn := range matches {
		content, err := ioutil.ReadFile(fn)
		if err != nil {
			t.Fatal(fmt.Sprintf("read File %v with error %v", fn, err.Error()))
		}
		sfn := fn + ".slate"
		slate, err := ioutil.ReadFile(sfn)
		if err != nil {
			t.Fatal(fmt.Sprintf("read File %v with error %v", sfn, err.Error()))
		}
		s := Run(t, string(content))
		assert.Equal(t, s, string(slate), "")
	}
}
