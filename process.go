package markdown_to_slate

import (
	"log"
	"strings"

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

	if len(leaves) > 0 {
		for i, l := range leaves {
			leaves[i].Text = l.Text + "\n"
		}
	}

	return leaves
}

func unnestParagraphs(nodes *[]Node, depth int) {
	for _, node := range *nodes {
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
				nds := []Node{}
				for _, nd := range node.Nodes {
					nds = append(nds, nd.Nodes...)
				}
				//log.Println("unnset node", node.Nodes, "to", nds)
				node.Nodes = nds
			}
		}
	}

	for i, _ := range *nodes {
		if (*nodes)[i].Type == "paragraph" {
			if len((*nodes)[i].Nodes) > 0 {
				if (*nodes)[i].Nodes[0].Type == "paragraph" {
					(*nodes)[i].Type = "div"
				}
				if len((*nodes)[i].Nodes[0].Nodes) > 0 {
					if (*nodes)[i].Nodes[0].Nodes[0].Type == "paragraph" {
						(*nodes)[i].Type = "div"
					}
				}
			}
		}
	}

	for i, _ := range *nodes {
		if len((*nodes)[i].Nodes) > 0 {
			unnestParagraphs(&(*nodes)[i].Nodes, depth+1)
		}
	}
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
				//log.Println("wrapNext runs")
				//nodes = append(nodes, Node{
				//Object: "block",
				//Type:   "paragraph",
				//Nodes:  []Node{*n},
				//})
				nodes = append(nodes, *n)
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

	unnestParagraphs(&nodes, 0)

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
		return nil, true
	}
	if node.Type == blackfriday.HorizontalRule {
		return nil, true
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
				l.Text = strings.Replace(l.Text[:len(l.Text)-1], "\n\n", "\n", -1)
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
		var nds []Node
		if node.FirstChild != nil && node.FirstChild == node.LastChild && node.FirstChild.Type == blackfriday.Paragraph {
			nds = ProcessChildren(node.FirstChild, level+1)
		} else {
			nds = ProcessChildren(node, level+1)
		}
		return &Node{
			Object: "block",
			Type:   "list-item",
			Nodes:  nds,
		}, false
	}

	if node.Type == blackfriday.BlockQuote {
		nds := ProcessChildren(node, level+1)
		return &Node{
			Object: "block",
			Type:   "block-quote",
			Nodes:  nds,
		}, false
	}

	if node.Type == blackfriday.HTMLBlock {
		//log.Println("html:", node.Literal)
		if string(node.Literal) == "" {
			return nil, false
		}
		l := ProcessTextNode(node)
		return &Node{
			Object: "block",
			Type:   "code",
			Nodes: []Node{Node{
				Object: "text",
				Leaves: []Leaf{l},
			}},
		}, false

		//return &Node{
		//Object: "block",
		//Type:   "code",
		//Nodes: []Node{Node{
		//Object: "block",
		//Type:   "code_line",
		//Nodes: []Node{Node{
		//Object: "text",
		//Leaves: []Leaf{l},
		//}},
		//}},
		//}, false
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
		//scs := spew.ConfigState{DisableMethods: true, Indent: "\t"}
		//node.Parent = nil
		//node.FirstChild.Parent = nil
		//scs.Dump(node.FirstChild)

		//log.Println(node.FirstChild.Type)

		var nds []Node
		ltype := "block"
		if node.FirstChild.Next == nil && node.FirstChild.Type == blackfriday.Text {
			//log.Println("found text link")
			l := ProcessTextNode(node.FirstChild)
			if l.Text == "" {
				l.Text = "link"
			}
			nds = []Node{Node{
				Object: "text",
				Leaves: []Leaf{l},
			}}
			ltype = "inline"
		} else {
			nds = ProcessChildren(node, level+1)
		}
		//nds[0].Type == "paragraph"

		//if len(nds) == 1 && nds[0].Type == "image" {
		//nds[0].Nodes = nil
		//return &nds[0]
		//}

		//log.Println("nodes1:", len(nds), nds[0].Type)
		//log.Println("nodes2:", len(nds[0].Nodes), nds[0].Nodes[0].Type, nds[0].Nodes[0].Object)

		//if len(nds) == 1 && nds[0].Type == "paragraph" && len(nds[0].Nodes) == 1 && nds[0].Nodes[0].Object == "text" {
		//ltype = "inline"
		//nds[0] = nds[0].Nodes[0]
		//}

		return &Node{
			Object: ltype,
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
