package markdown_to_slate

import (
	"log"

	blackfriday "gopkg.in/russross/blackfriday.v2"
)

type Mark struct {
	Data   map[string]interface{} `json:"data,omitempty"`
	Object string                 `json:"object,omitempty"`
	Type   string                 `json:"type,omitempty"`
}

type Leaf struct {
	Object string `json:"object"`
	Text   string `json:"text"`
	Marks  []Mark `json:"marks"`
}

type Node struct {
	Object string                 `json:"object,omitempty"`
	Type   string                 `json:"type,omitempty"`
	Data   map[string]interface{} `json:"data,omitempty"`
	Nodes  []Node                 `json:"nodes,omitempty"`
	Leaves []Leaf                 `json:"leaves,omitempty"`
}

func ProcessTextNode(node *blackfriday.Node) Leaf {
	text := string(node.Literal)
	leaf := Leaf{
		Object: "leaf",
		Text:   text,
		Marks:  []Mark{},
	}

	return leaf
}

func ProcessParagraphNode(node *blackfriday.Node) Node {
	leaves := make([]Leaf, 0)
	node.Walk(func(lnode *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		if node == lnode {
			return blackfriday.GoToNext
		}
		if !entering {
			return blackfriday.GoToNext
		}

		if lnode.Type == blackfriday.Text {
			leaf := ProcessTextNode(lnode)
			leaves = append(leaves, leaf)
			return blackfriday.SkipChildren
		}
		if lnode.Type == blackfriday.Hardbreak {
			return blackfriday.SkipChildren
		}

		if lnode.Type == blackfriday.Paragraph {
			cn := ProcessParagraphNode(lnode)
			leaves = append(leaves, cn.Leaves...)
			return blackfriday.SkipChildren
		}

		if lnode.Type == blackfriday.Link {
			// TODO
			log.Println("TODO link:", string(lnode.LinkData.Destination), string(lnode.LinkData.Title))
			return blackfriday.SkipChildren
		}

		if lnode.Type == blackfriday.Image {
			// TODO
			log.Println("TODO image:", string(lnode.LinkData.Destination), string(lnode.LinkData.Title))
			return blackfriday.SkipChildren
		}

		log.Println("not processing child node in paragraph:", lnode.Type, "::", string(node.Literal))
		return blackfriday.GoToNext
	})

	list := Node{
		Object: "block",
		Type:   "paragraph",
		Nodes: []Node{
			Node{
				Object: "text",
				Leaves: leaves,
			},
		},
	}
	return list
}

func ProcessListNode(node *blackfriday.Node) []Node {
	listType := "numbered-list"
	if node.ListData.ListFlags&blackfriday.ListTypeOrdered == blackfriday.ListTypeOrdered {
		listType = "numbered-list"
	} else {
		listType = "bulleted-list"
	}

	list := make([]Node, 0)

	node.Walk(func(cnode *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		if node == cnode {
			return blackfriday.GoToNext
		}
		if !entering {
			return blackfriday.GoToNext
		}

		if cnode.Type == blackfriday.Item {
			plist := make([]Node, 0)
			cnode.Walk(func(lnode *blackfriday.Node, entering bool) blackfriday.WalkStatus {
				if !entering {
					return blackfriday.GoToNext
				}
				tlist := make([]Node, 0)
				ProcessNode(&tlist, lnode)
				for _, n := range tlist {
					plist = append(plist, n.Nodes...)
				}
				return blackfriday.GoToNext
			})
			list = append(list, Node{Object: "block", Type: "list-item", Nodes: plist})
			return blackfriday.SkipChildren
		}
		return blackfriday.GoToNext
	})

	if listType == "bulleted-list" {
		//log.Println("list node done, children:")
		allChecks := true
		for _, node := range list {
			if len(node.Nodes) == 0 || len(node.Nodes[0].Leaves) == 0 {
				allChecks = false
				continue
			}
			tx := node.Nodes[0].Leaves[0].Text
			if len(tx) < 4 {
				allChecks = false
				continue
			}
			prefix := tx[0:4]
			if len(node.Nodes) > 0 && (prefix == "[ ] " || prefix == "[x] ") {
				node.Nodes[0].Object = "block"
				node.Nodes[0].Type = "check-list-item"
				node.Nodes[0].Data = map[string]interface{}{}

				checked := (prefix == "[x] ")
				node.Nodes[0].Data["checked"] = checked

				cnode := Node{Object: "text", Leaves: []Leaf{Leaf{
					Object: "leaf",
					Text:   tx[4:len(tx)],
					Marks:  []Mark{},
				}}}
				node.Nodes[0].Nodes = []Node{cnode}
				node.Nodes[0].Leaves = []Leaf{}
			} else {
				allChecks = false
			}
			//spew.Dump(node)
		}

		if allChecks {
			nds := make([]Node, 0)
			for _, n := range list {
				nds = append(nds, n.Nodes...)
			}
			return nds
		}
	}

	//spew.Dump(list)

	return []Node{Node{Object: "block", Type: listType, Nodes: list}}
}

func ProcessNode(ret *[]Node, node *blackfriday.Node) blackfriday.WalkStatus {
	//log.Println("process node:", node.Type, "::", string(node.Literal))
	if node.Type == blackfriday.Hardbreak {
		return blackfriday.SkipChildren
	}

	if node.Type == blackfriday.List {
		nds := ProcessListNode(node)
		//spew.Dump(nds)
		*ret = append(*ret, nds...)

		return blackfriday.SkipChildren
	}

	//if node.Type == blackfriday.Text {
	//nds := ProcessTextNode(node)
	//*ret = append(*ret, nds)
	//return blackfriday.SkipChildren
	//}

	if node.Type == blackfriday.Paragraph || node.Type == blackfriday.CodeBlock {
		//if len(node.Literal) == 0 {
		//return blackfriday.GoToNext
		//}
		nds := ProcessParagraphNode(node)
		*ret = append(*ret, nds)
		return blackfriday.SkipChildren
	}

	return blackfriday.GoToNext
}

func Parse(input []byte) []Node {
	//log.Println(string(input))

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

		//log.Println(node.Type, "::", string(node.Literal))
		return ProcessNode(&ret, node)
		//return r.renderNode(&buf, node, entering)
	})

	return ret
}
