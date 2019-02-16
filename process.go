package markdown_to_slate

import (
	"log"
	"strconv"

	blackfriday "gopkg.in/russross/blackfriday.v2"
)

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

func ProcessChildren(parent *blackfriday.Node) []Node {
	nodes := make([]Node, 0)

	if parent.FirstChild == nil {
		return nodes
	}
	//log.Println(len(parent.Literal))
	//log.Println(parent.Text)

	node := parent.FirstChild
	for {
		n := ProcessNode(node)
		if n != nil {
			nodes = append(nodes, *n)
		}

		node = node.Next
		if node == nil {
			break
		}
	}

	return nodes
}

func ProcessNode(node *blackfriday.Node) *Node {
	if node.Type == blackfriday.Hardbreak {
		return nil
	}

	if node.Type == blackfriday.Text {
		l := ProcessTextNode(node)
		if l.Text == "" {
			return nil
		}
		return &Node{
			Object: "inline",
			Type:   "inline",
			Leaves: []Leaf{l},
		}
	}

	if node.Type == blackfriday.HTMLSpan {
		l := ProcessTextNode(node)
		if l.Text == "" {
			return nil
		}
		return &Node{
			Object: "inline",
			Type:   "inline",
			Leaves: []Leaf{l},
		}
	}

	if node.Type == blackfriday.Emph {
		lvs := ProcessTextChildren(node, []Mark{Mark{
			Object: "mark",
			Type:   "italic",
		}})
		return &Node{
			Object: "inline",
			Type:   "inline",
			Leaves: lvs,
		}
	}

	if node.Type == blackfriday.Strong {
		lvs := ProcessTextChildren(node, []Mark{Mark{
			Object: "mark",
			Type:   "bold",
		}})
		return &Node{
			Object: "inline",
			Type:   "inline",
			Leaves: lvs,
		}
	}

	if node.Type == blackfriday.Code {
		lvs := ProcessTextChildren(node, []Mark{Mark{
			Object: "mark",
			Type:   "code",
		}})
		return &Node{
			Object: "inline",
			Type:   "inline",
			Leaves: lvs,
		}
	}

	if node.Type == blackfriday.Del {
		lvs := ProcessTextChildren(node, []Mark{Mark{
			Object: "mark",
			Type:   "del",
		}})
		return &Node{
			Object: "inline",
			Type:   "inline",
			Leaves: lvs,
		}
	}

	if node.Type == blackfriday.List {
		listType := "numbered-list"
		if node.ListData.ListFlags&blackfriday.ListTypeOrdered == blackfriday.ListTypeOrdered {
			listType = "numbered-list"
		} else {
			listType = "bulleted-list"
		}

		nds := ProcessChildren(node)

		if listType == "bulleted-list" {
			allChecks := processChecks(&nds)
			if allChecks {
				return &Node{
					Object: "block",
					Type:   "paragraph",
					Nodes:  nds,
				}
			}
		}

		return &Node{
			Object: "block",
			Type:   listType,
			Nodes:  nds,
		}
	}

	if node.Type == blackfriday.Heading {
		nds := ProcessChildren(node)
		return &Node{
			Object: "block",
			Type:   "heading-" + strconv.Itoa(node.HeadingData.Level),
			Nodes:  nds,
		}
	}

	if node.Type == blackfriday.CodeBlock {
		//scs := spew.ConfigState{DisableMethods: true, Indent: "\t"}
		//scs.Dump(node)
		nds := ProcessChildren(node)
		//spew.Dump(nds)
		if string(node.Literal) != "" {
			nds = append(nds, Node{
				Object: "block",
				Type:   "code_line",
				Leaves: []Leaf{ProcessTextNode(node)},
			})
		}
		return &Node{
			Object: "block",
			Type:   "code",
			Nodes:  nds,
		}
	}

	if node.Type == blackfriday.Item {
		nds := ProcessChildren(node)
		return &Node{
			Object: "block",
			Type:   "list-item",
			Nodes:  nds,
		}
	}

	if node.Type == blackfriday.Paragraph {
		nds := ProcessChildren(node)
		return &Node{
			Object: "block",
			Type:   "paragraph",
			Nodes:  nds,
		}
	}

	if node.Type == blackfriday.Link {
		nds := ProcessChildren(node)
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
		}
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
		}
	}

	log.Println("not processing child node in paragraph:", node.Type, "::", string(node.Literal))
	return nil
}
