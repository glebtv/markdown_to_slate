package markdown_to_slate

import (
	blackfriday "gopkg.in/russross/blackfriday.v2"
)

func Parse(input []byte) []Node {
	processor := blackfriday.New(
		blackfriday.WithExtensions(
			blackfriday.CommonExtensions | blackfriday.HardLineBreak | blackfriday.AutoHeadingIDs | blackfriday.Autolink,
		))

	data := processor.Parse(input)

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
