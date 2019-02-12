package markdown_to_slate

import (
	blackfriday "gopkg.in/russross/blackfriday.v2"
)

type Parser struct {
	processor *blackfriday.Markdown
}

func NewParser() *Parser {
	parser := &Parser{}
	parser.processor = blackfriday.New(
		blackfriday.WithExtensions(
			blackfriday.CommonExtensions | blackfriday.HardLineBreak | blackfriday.AutoHeadingIDs | blackfriday.Autolink,
		))

	return parser
}

func (p *Parser) Parse(input []byte) []Node {
	data := p.processor.Parse(input)

	ret := make([]Node, 0)

	data.Walk(func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		if node == data {
			return blackfriday.GoToNext
		}
		if !entering {
			return blackfriday.GoToNext
		}

		return ProcessNode(&ret, node)
	})

	return ret
}
