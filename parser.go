package markdown_to_slate

import (
	"github.com/gernest/mention"
	blackfriday "gopkg.in/russross/blackfriday.v2"

	"strings"
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

func (p *Parser) ParseWithoutMentions(input []byte) []Node {
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

func (p *Parser) Parse(input []byte) []Node {
	issues := mention.GetTags('#', string(input))

	for i, _ := range issues {
		index := len(issues) - i - 1
		issue := issues[index]
		//log.Println("replace issue", index, issue)
		//spew.Dump(issue)
		replace := "[#" + issue.Tag + "]" + "(#" + issue.Tag + ")"
		input = []byte(string(input[:issue.Index]) + replace + string(input[issue.Index+len(issue.Tag)+1:]))
	}

	mentions := mention.GetTags('@', string(input))

	for i, _ := range mentions {
		index := len(mentions) - i - 1
		mention := mentions[index]
		replace := "[@" + mention.Tag + "]" + "(@" + mention.Tag + ")"
		input = []byte(string(input[:mention.Index]) + replace + string(input[mention.Index+len(mention.Tag)+1:]))
	}
	input = []byte(strings.Replace(string(input), "+", "♀", -1))
	input = []byte(strings.Replace(string(input), "\n", "\n\n", -1))
	//log.Println("replaced mentions:", string(input))

	nodes := p.ParseWithoutMentions(input)
	for i, _ := range nodes {
		nodes[i].Replace("♀", "+")
	}
	return nodes
}
