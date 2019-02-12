package markdown_to_slate

import (
	blackfriday "gopkg.in/russross/blackfriday.v2"
)

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
	nodes := make([]Node, 0)
	node.Walk(func(lnode *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		if node == lnode {
			return blackfriday.GoToNext
		}
		if !entering {
			return blackfriday.GoToNext
		}

		if lnode.Type == blackfriday.Text {
			leaf := ProcessTextNode(lnode)
			if leaf.Text != "" {
				leaves = append(leaves, leaf)
			}
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

		if lnode.Type == blackfriday.HTMLSpan {
			// TODO
			cn := ProcessParagraphNode(lnode)
			leaves = append(leaves, cn.Leaves...)
			return blackfriday.SkipChildren
		}

		if lnode.Type == blackfriday.Emph {
			// TODO
			cn := ProcessParagraphNode(lnode)
			leaves = append(leaves, cn.Leaves...)
			return blackfriday.SkipChildren
		}

		if lnode.Type == blackfriday.Strong {
			// TODO
			cn := ProcessParagraphNode(lnode)
			leaves = append(leaves, cn.Leaves...)
			return blackfriday.SkipChildren
		}

		if lnode.Type == blackfriday.Code {
			// TODO
			cn := ProcessParagraphNode(lnode)
			leaves = append(leaves, cn.Leaves...)
			return blackfriday.SkipChildren
		}

		if lnode.Type == blackfriday.Del {
			// TODO
			cn := ProcessParagraphNode(lnode)
			leaves = append(leaves, cn.Leaves...)
			return blackfriday.SkipChildren
		}

		if lnode.Type == blackfriday.Link {
			if len(leaves) > 0 {
				nodes = append(nodes, Node{
					Object: "text",
					Leaves: leaves,
				})
				leaves = []Leaf{}
			}
			//spew.Dump(lnode)
			cn := ProcessParagraphNode(lnode)
			if len(cn.Nodes) == 1 && cn.Nodes[0].Type == "image" {
				cn.Nodes[0].Nodes = nil
				nodes = append(nodes, cn.Nodes[0])
				return blackfriday.SkipChildren
			}
			nodes = append(nodes, Node{
				Object: "inline",
				Type:   "link",
				Data: map[string]interface{}{
					"href":  string(lnode.LinkData.Destination),
					"title": string(lnode.LinkData.Title),
				},
				Leaves: cn.Leaves,
				Nodes:  cn.Nodes,
			})
			return blackfriday.SkipChildren
			//return blackfriday.GoToNext
		}

		if lnode.Type == blackfriday.Image {
			if len(leaves) > 0 {
				nodes = append(nodes, Node{
					Object: "text",
					Leaves: leaves,
				})
				leaves = []Leaf{}
			}
			cn := ProcessParagraphNode(lnode)

			nodes = append(nodes, Node{
				Object: "block",
				Type:   "image",
				Data: map[string]interface{}{
					"src":   string(lnode.LinkData.Destination),
					"title": string(lnode.LinkData.Title),
				},
				Nodes:  cn.Nodes,
				Leaves: cn.Leaves,
			})
			return blackfriday.SkipChildren
		}

		//log.Println("not processing child node in paragraph:", lnode.Type, "::", string(node.Literal))
		return blackfriday.GoToNext
	})

	if len(leaves) > 0 {
		nodes = append(nodes, Node{
			Object: "text",
			Leaves: leaves,
		})
	}

	list := Node{
		Object: "block",
		Type:   "paragraph",
		Nodes:  nodes,
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
		}

		if allChecks {
			nds := make([]Node, 0)
			for _, n := range list {
				nds = append(nds, n.Nodes...)
			}
			return nds
		}
	}

	return []Node{Node{Object: "block", Type: listType, Nodes: list}}
}

func ProcessNode(ret *[]Node, node *blackfriday.Node) blackfriday.WalkStatus {
	if node.Type == blackfriday.Hardbreak {
		return blackfriday.SkipChildren
	}

	if node.Type == blackfriday.List {
		nds := ProcessListNode(node)
		*ret = append(*ret, nds...)

		return blackfriday.SkipChildren
	}

	if node.Type == blackfriday.Paragraph || node.Type == blackfriday.CodeBlock {
		nds := ProcessParagraphNode(node)
		*ret = append(*ret, nds)
		return blackfriday.SkipChildren
	}

	return blackfriday.GoToNext
}
