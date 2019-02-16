package markdown_to_slate

import (
	"log"

	blackfriday "gopkg.in/russross/blackfriday.v2"
)

var hLevels = map[int]string{
	1: "one",
	2: "two",
	3: "three",
	4: "four",
	5: "five",
	6: "six",
}

func ProcessTextNode(node *blackfriday.Node) Leaf {
	return Leaf{
		Object: "leaf",
		Text:   string(node.Literal),
		Marks:  []Mark{},
	}
}

func ProcessTextChildren(parent *blackfriday.Node, marks []Mark) []Leaf {
	leaves := make([]Leaf, 0)

	node := parent.FirstChild
	for {
		if node == nil {
			break
		}
		//log.Println("process text node", node)
		l := ProcessTextNode(node)

		//log.Println(l.Text)
		if l.Text != "" {
			l.Marks = marks
			leaves = append(leaves, l)
		}

		node = node.Next
	}

	return leaves
}

func ProcessChildren(parent *blackfriday.Node, level int) []Node {
	nodes := make([]Node, 0)

	if parent.FirstChild == nil {
		return nodes
	}
	//log.Println(len(parent.Literal))
	//log.Println(parent.Text)

	//log.Println("process level: ", level)

	wrapNext := true
	node := parent.FirstChild
	for {
		n, _wrapNext := ProcessNode(node, level)
		//log.Println("processed node:", node, wrapNext, _wrapNext)
		//log.Println("")
		if n != nil {
			if wrapNext && n.Object == "text" {
				nodes = append(nodes, Node{
					Object: "block",
					Type:   "paragraph",
					Nodes:  []Node{*n},
				})
			} else {
				//if n.Object == "text" {
				//log.Println("!!!!!!!!!!!!!!")
				//spew.Dump(n)
				//}
				if n.Object == "text" && len(nodes) > 0 && nodes[len(nodes)-1].Type == "paragraph" {
					nodes[len(nodes)-1].Nodes = append(nodes[len(nodes)-1].Nodes, *n)
				} else {
					nodes = append(nodes, *n)
				}
			}
		}

		wrapNext = false
		if _wrapNext {
			wrapNext = true
		}

		node = node.Next
		if node == nil {
			break
		}
	}

	for i, node := range nodes {
		//log.Println("after process", i, node.Type)
		//spew.Dump(nodes)
		if node.Type == "paragraph" {
			allPara := true
			for _, nd := range node.Nodes {
				if nd.Type != "paragraph" {
					allPara = false
				}
			}
			if allPara {
				nodes = append([]Node{}, nodes[:i]...)
				nodes = append(nodes, node.Nodes...)
				nodes = append(nodes, nodes[i:len(nodes)-1]...)
			}
		}
	}

	//log.Println("dump level: ", level)
	//for i, node := range nodes {
	//spew.Dump(i, node)
	//}
	//log.Println("done level: ", level)

	return nodes
}

func ProcessNode(node *blackfriday.Node, level int) (*Node, bool) {
	//log.Println("process node:", node)
	if node.Type == blackfriday.Hardbreak {
		//return &Node{
		//Object: "block",
		//Type:   "paragraph",
		//Nodes: []Node{Node{
		//Object: "text",
		//Leaves: []Leaf{Leaf{
		//Object: "leaf",
		//Text:   string(""),
		//Marks:  []Mark{},
		//}},
		//}},
		//}, true
		return nil, true
		//return nil
	}

	if node.Type == blackfriday.Text {
		l := ProcessTextNode(node)
		if l.Text == "" {
			return nil, false
		}
		return &Node{
			Object: "text",
			Leaves: []Leaf{l},
		}, false
	}

	if node.Type == blackfriday.HTMLSpan {
		l := ProcessTextNode(node)
		if l.Text == "" {
			return nil, false
		}
		return &Node{
			Object: "text",
			Leaves: []Leaf{l},
		}, false
	}

	if node.Type == blackfriday.Emph {
		lvs := ProcessTextChildren(node, []Mark{Mark{
			Object: "mark",
			Type:   "italic",
		}})
		return &Node{
			Object: "text",
			Leaves: lvs,
		}, false
	}

	if node.Type == blackfriday.Strong {
		lvs := ProcessTextChildren(node, []Mark{Mark{
			Object: "mark",
			Type:   "bold",
		}})
		return &Node{
			Object: "text",
			Leaves: lvs,
		}, false
	}

	if node.Type == blackfriday.Code {
		leaf := ProcessTextNode(node)
		leaf.Marks = []Mark{Mark{
			Object: "mark",
			Type:   "code",
		}}
		return &Node{
			Object: "text",
			Leaves: []Leaf{leaf},
		}, false
	}

	if node.Type == blackfriday.Del {
		lvs := ProcessTextChildren(node, []Mark{Mark{
			Object: "mark",
			Type:   "del",
		}})
		return &Node{
			Object: "text",
			Leaves: lvs,
		}, false
	}

	if node.Type == blackfriday.List {
		listType := "numbered-list"
		if node.ListData.ListFlags&blackfriday.ListTypeOrdered == blackfriday.ListTypeOrdered {
			listType = "numbered-list"
		} else {
			listType = "bulleted-list"
		}

		nds := ProcessChildren(node, level+1)

		if listType == "bulleted-list" {
			allChecks := processChecks(&nds)
			if allChecks {
				return &Node{
					Object: "block",
					Type:   "paragraph",
					Nodes:  nds,
				}, false
			}
		}

		return &Node{
			Object: "block",
			Type:   listType,
			Nodes:  nds,
		}, false
	}

	if node.Type == blackfriday.Heading {
		nds := ProcessChildren(node, level+1)
		return &Node{
			Object: "block",
			Type:   "heading-" + hLevels[node.HeadingData.Level],
			Nodes:  nds,
		}, false
	}

	if node.Type == blackfriday.CodeBlock {
		//scs := spew.ConfigState{DisableMethods: true, Indent: "\t"}
		//scs.Dump(node.Literal)
		nds := ProcessChildren(node, level+1)
		//spew.Dump(nds)
		if string(node.Literal) != "" {
			l := ProcessTextNode(node)
			// remove last newline of code block
			if l.Text[len(l.Text)-1] == byte('\n') {
				l.Text = l.Text[:len(l.Text)-1]
			}
			nds = append(nds, Node{
				Object: "block",
				Type:   "code_line",
				Nodes: []Node{Node{
					Object: "text",
					Leaves: []Leaf{l},
				}},
			})
		}
		return &Node{
			Object: "block",
			Type:   "code",
			Nodes:  nds,
		}, false
	}

	if node.Type == blackfriday.Item {
		nds := ProcessChildren(node, level+1)
		return &Node{
			Object: "block",
			Type:   "list-item",
			Nodes:  nds,
		}, false
	}

	if node.Type == blackfriday.Paragraph {
		//log.Println("para", node.Literal)
		nds := ProcessChildren(node, level+1)
		if len(nds) == 1 && nds[0].Type == "paragraph" {
			return &nds[0], false
		}
		return &Node{
			Object: "block",
			Type:   "paragraph",
			Nodes:  nds,
		}, false
	}

	if node.Type == blackfriday.Link {
		nds := ProcessChildren(node, level+1)
		//if len(nds) == 1 && nds[0].Type == "image" {
		//nds[0].Nodes = nil
		//return &nds[0]
		//}

		return &Node{
			Object: "inline",
			Type:   "link",
			Data: map[string]interface{}{
				"href":  string(node.LinkData.Destination),
				"title": string(node.LinkData.Title),
			},
			Nodes: nds,
		}, false
		//return blackfriday.GoToNext
	}

	if node.Type == blackfriday.Image {
		return &Node{
			Object: "block",
			Type:   "image",
			Data: map[string]interface{}{
				"href":  string(node.LinkData.Destination),
				"title": string(node.LinkData.Title),
			},
		}, false
	}

	log.Println("not processing child node in paragraph:", node.Type, "::", string(node.Literal))
	return nil, false
}
