package markdown_to_slate

import (
	"log"
	"strings"

	"github.com/jinzhu/copier"
	"github.com/russross/blackfriday"
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

	ret := make([]Node, 0)
	for _, n := range *nodes {
		if n.Object == "block" && len(n.Nodes) > 0 {
			nds := fixInlineSiblings(&n)
			ret = append(ret, nds...)
		} else {
			ret = append(ret, n)
		}
	}
	*nodes = ret

	for i, _ := range *nodes {
		if len((*nodes)[i].Nodes) > 0 {
			unnestParagraphs(&(*nodes)[i].Nodes, depth+1)
		}
	}
}

func fixInlineSiblings(node *Node) []Node {
	// wrap inlines in divs if there are both inline and block children
	// this is due to slate core schema rule which forbids such a case
	hasInline := false
	hasBlock := false
	for _, n := range node.Nodes {
		if n.Object == "text" || n.Object == "inline" {
			hasInline = true
		}
		if n.Object == "block" {
			hasBlock = true
		}
	}

	//spew.Dump(node, hasInline, hasBlock)

	if hasBlock && hasInline {
		ret := make([]Node, 0)
		inlines := make([]Node, 0)
		baseUsed := false
		for _, n := range node.Nodes {
			if n.Object == "text" || n.Object == "inline" {
				inlines = append(inlines, n)
			} else {
				if len(inlines) > 0 {
					var nd Node
					if baseUsed {
						nd = Node{
							Object: "block",
							Type:   "div",
							Nodes:  inlines,
						}
						baseUsed = true
					} else {
						copier.Copy(&nd, &node)
						nd.Nodes = inlines
					}
					ret = append(ret, nd)
					inlines = []Node{}
				}
				ret = append(ret, n)
			}
		}
		if len(inlines) > 0 {
			var nd Node
			if baseUsed {
				nd = Node{
					Object: "block",
					Type:   "div",
					Nodes:  inlines,
				}
				baseUsed = true
			} else {
				copier.Copy(&nd, &node)
				nd.Nodes = inlines
			}
			ret = append(ret, nd)
		}
		return ret
	} else {
		return []Node{*node}
	}
}

func ProcessChildren(parent *blackfriday.Node, level int) []Node {
	nodes := make([]Node, 0)

	if parent.FirstChild == nil {
		return nodes
	}

	node := parent.FirstChild
	for {
		n := ProcessNode(node, level)
		//log.Println("processed node:", node)
		if n != nil {
			nodes = append(nodes, *n)
		}

		node = node.Next
		if node == nil {
			break
		}
	}

	if level == 0 {
		unnestParagraphs(&nodes, 0)
	}

	return nodes
}

func ProcessNode(node *blackfriday.Node, level int) *Node {
	//log.Println("process node:", node)
	if node.Type == blackfriday.Hardbreak {
		return nil
	}
	if node.Type == blackfriday.HorizontalRule {
		return nil
	}

	if node.Type == blackfriday.Text {
		l := ProcessTextNode(node)
		if l.Text == "" {
			return nil
		}
		return &Node{
			Object: "text",
			Leaves: []Leaf{l},
		}
	}

	if node.Type == blackfriday.HTMLSpan {
		l := ProcessTextNode(node)
		if l.Text == "" {
			return nil
		}
		return &Node{
			Object: "text",
			Leaves: []Leaf{l},
		}
	}

	if node.Type == blackfriday.Emph {
		lvs := ProcessTextChildren(node, []Mark{Mark{
			Object: "mark",
			Type:   "italic",
		}})
		return &Node{
			Object: "text",
			Leaves: lvs,
		}
	}

	if node.Type == blackfriday.Strong {
		lvs := ProcessTextChildren(node, []Mark{Mark{
			Object: "mark",
			Type:   "bold",
		}})
		return &Node{
			Object: "text",
			Leaves: lvs,
		}
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
		}
	}

	if node.Type == blackfriday.Del {
		lvs := ProcessTextChildren(node, []Mark{Mark{
			Object: "mark",
			Type:   "del",
		}})
		return &Node{
			Object: "text",
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

		nds := ProcessChildren(node, level+1)

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
		nds := ProcessChildren(node, level+1)
		return &Node{
			Object: "block",
			Type:   "heading-" + hLevels[node.HeadingData.Level],
			Nodes:  nds,
		}
	}

	if node.Type == blackfriday.CodeBlock {
		nds := ProcessChildren(node, level+1)
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
		}
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
		}
	}

	if node.Type == blackfriday.BlockQuote {
		nds := ProcessChildren(node, level+1)
		return &Node{
			Object: "block",
			Type:   "block-quote",
			Nodes:  nds,
		}
	}

	if node.Type == blackfriday.HTMLBlock {
		//log.Println("html:", node.Literal)
		if string(node.Literal) == "" {
			return nil
		}
		l := ProcessTextNode(node)
		return &Node{
			Object: "block",
			Type:   "code",
			Nodes: []Node{Node{
				Object: "text",
				Leaves: []Leaf{l},
			}},
		}
	}

	if node.Type == blackfriday.Paragraph {
		//log.Println("para", node.Literal)
		nds := ProcessChildren(node, level+1)
		if len(nds) == 1 && nds[0].Type == "paragraph" {
			return &nds[0]
		}
		return &Node{
			Object: "block",
			Type:   "paragraph",
			Nodes:  nds,
		}
	}

	if node.Type == blackfriday.Link {
		var nds []Node
		ltype := "block"
		if node.FirstChild.Next == nil && node.FirstChild.Type == blackfriday.Text {
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
			// TODO
			// slate не нравятся картинки в ссылках
			if len(nds) == 1 && nds[0].Type == "image" {
				return &nds[0]
			}
		}

		return &Node{
			Object: ltype,
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
				"src":   string(node.LinkData.Destination),
				"title": string(node.LinkData.Title),
			},
		}
	}

	log.Println("not processing child node in paragraph:", node.Type, "::", string(node.Literal))
	return nil
}
