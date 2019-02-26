package markdown_to_slate

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	jd "github.com/josephburnett/jd/lib"
	"github.com/tidwall/pretty"
)

func Run(input string) string {
	data := Parse([]byte(input))
	s, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return string(s)
}

func RunReverse(input string) string {
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

func assertEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if a == b {
		return
	}
	if len(message) == 0 {
		aj, err := jd.ReadJsonString(a.(string))
		if err != nil {
			panic(err)
		}
		bj, err := jd.ReadJsonString(b.(string))
		if err != nil {
			panic(err)
		}
		message = fmt.Sprintf(
			"\ngot:\n%v\nexpected:\n%v",
			printJSON(a.(string)),
			printJSON(b.(string)),
		)
		message += "\n" + aj.Diff(bj).Render()
	}
	t.Fatal(message)
}

const paragraphExample = "test paragraph\n"
const paragraphResult = `[{"object":"block","type":"paragraph","nodes":[{"object":"text","leaves":[{"object":"leaf","text":"test paragraph","marks":[]}]}]}]`

func TestParagraph(t *testing.T) {
	s := Run(paragraphExample)

	assertEqual(t, s, paragraphResult, "")
}

func TestParagraphReverse(t *testing.T) {
	s := RunReverse(paragraphExample)

	assertEqual(t, s, paragraphExample, "")
}

const listExample = "- list\n- items"
const listResult = `[{"object":"block","type":"bulleted-list","nodes":[{"object":"block","type":"list-item","nodes":[{"object":"text","leaves":[{"object":"leaf","text":"list","marks":[]}]}]},{"object":"block","type":"list-item","nodes":[{"object":"text","leaves":[{"object":"leaf","text":"items","marks":[]}]}]}]}]`

func TestList(t *testing.T) {
	s := Run(listExample)

	assertEqual(t, s, listResult, "")
}

const checklistResult = `[{"object":"block","type":"paragraph","nodes":[{"object":"text","leaves":[{"object":"leaf","text":"check list:","marks":[]},{"object":"leaf","text":"- [ ] unchecked","marks":[]},{"object":"leaf","text":"- [x] checked","marks":[]}]}]}]`

func TestChecklist(t *testing.T) {
	s := Run("check list:\n- [ ] unchecked\n- [x] checked")

	assertEqual(t, s, checklistResult, "")
}

const codeResult = `[{"object":"block","type":"paragraph","nodes":[{"object":"text","leaves":[{"object":"leaf","text":"code ","marks":[]},{"object":"leaf","text":":","marks":[]}]}]},{"object":"block","type":"paragraph"}]`

func TestCode(t *testing.T) {
	s := Run("code ```inline```:\n```\nblock code\n```\n")

	assertEqual(t, s, codeResult, "")
}

const linkExample = `[rt1622_regions_sorted.csv](/original/rt1622_regions_sorted.csv) (340.9 KiB)
[rt1622_cities_sorted.csv](/original/rt1622_cities_sorted.csv) (22.6 MiB)`
const linkResult = `[{"object":"block","type":"paragraph","nodes":[{"object":"text","leaves":[{"object":"leaf","text":"","marks":[]}]},{"object":"inline","type":"link","data":{"href":"/original/rt1622_regions_sorted.csv","title":""}},{"object":"text","leaves":[{"object":"leaf","text":" (340.9 KiB)","marks":[]},{"object":"leaf","text":"","marks":[]}]},{"object":"inline","type":"link","data":{"href":"/original/rt1622_cities_sorted.csv","title":""}},{"object":"text","leaves":[{"object":"leaf","text":" (22.6 MiB)","marks":[]}]}]}]`

func TestLink(t *testing.T) {
	s := Run(linkExample)

	assertEqual(t, s, linkResult, "")
}

const imageExample = `[![i4.png](/system/i4.jpg)](/system/i4.png?1451419607 "i4.png 22.6 KiB")`
const imageResult = `[{"object":"block","type":"paragraph","nodes":[{"object":"text","leaves":[{"object":"leaf","text":"","marks":[]}]},{"object":"inline","type":"link","data":{"href":"/system/i4.png?1451419607","title":"i4.png 22.6 KiB"}}]}]`

func TestImage(t *testing.T) {
	s := Run(imageExample)

	assertEqual(t, s, imageResult, "")
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
		s := Run(string(content))
		assertEqual(t, s, string(slate), "")
	}
}
