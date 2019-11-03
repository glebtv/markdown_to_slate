package markdown_to_slate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func RunStringify(t *testing.T, input string) string {
	s, err := Stringify([]byte(input))
	if err != nil {
		t.Fatal(err)
	}
	return string(s)
}

const newLineExample = `[{"object":"block","type":"paragraph","data":{},"nodes":[{"object":"text","text":"тест","marks":[]}]},{"object":"block","type":"paragraph","data":{},"nodes":[{"object":"text","text":"test","marks":[]}]}]`
const newLineResult = "тест\ntest\n"

func TestNewLine(t *testing.T) {
	s := RunStringify(t, newLineExample)
	assert.Equal(t, s, newLineResult, "")
}
