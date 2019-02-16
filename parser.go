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
			blackfriday.CommonExtensions |
				blackfriday.HardLineBreak |
				blackfriday.AutoHeadingIDs |
				blackfriday.Autolink,
		),
	)

	return parser
}

func (p *Parser) Parse(input []byte) []Node {
	data := p.processor.Parse(input)

	//scs := spew.ConfigState{DisableMethods: true, Indent: "\t"}
	//scs.Dump(data)

	return ProcessChildren(data, 0)

	//if data.FirstChild != nil {
	//nodes := []Node{}
	//child := data.FirstChild
	//for {
	//if child == nil {
	//break
	//}
	//nds := ProcessChildren(child, 1)
	//nodes = append(nodes, nds...)

	//child := data.FirstChild
	//}
	//return nodes
	//}

	//return []Node{}
	//scs.Dump(nodes)

}
