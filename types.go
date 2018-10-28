package markdown_to_slate

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
