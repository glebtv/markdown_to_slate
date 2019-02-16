package markdown_to_slate

func Parse(input []byte) []Node {
	return NewParser().Parse(input)
}
