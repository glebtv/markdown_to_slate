package markdown_to_slate

import "strings"

type Mark struct {
	Data   map[string]interface{} `json:"data,omitempty"`
	Object string                 `json:"object,omitempty"`
	Type   string                 `json:"type,omitempty"`
}

// deprecated
type Leaf struct {
	Object string `json:"object"`
	Text   string `json:"text"`
	Marks  []Mark `json:"marks"`
}

func (leaf *Leaf) Replace(from, to string) {
	leaf.Text = strings.Replace(leaf.Text, from, to, -1)
}

type Node struct {
	Object string                 `json:"object,omitempty"`
	Type   string                 `json:"type,omitempty"`
	Data   map[string]interface{} `json:"data,omitempty"`
	Nodes  []Node                 `json:"nodes,omitempty"`
	Text   *string                `json:"text,omitempty"`
	// deprecated
	Leaves []Leaf `json:"leaves,omitempty"`
}

func (node *Node) Replace(from, to string) {
	for i, _ := range node.Nodes {
		node.Nodes[i].Replace(from, to)
	}
	for i, _ := range node.Leaves {
		node.Leaves[i].Replace(from, to)
	}
}
