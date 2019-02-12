package markdown_to_slate

import (
	"bytes"
	"encoding/json"
)

func StringifyNodes(buf *bytes.Buffer, nodes []Node) {
	for _, node := range nodes {
		//if node.Text {
		//buf.WriteString(node.Text)
		//}
		if node.Nodes != nil {
			StringifyNodes(buf, node.Nodes)
		}
		if node.Leaves != nil {
			for _, leaf := range node.Leaves {
				if len(leaf.Text) > 0 {
					if leaf.Text != "" {
						buf.WriteString(leaf.Text + "\n")
					}
				}
			}
		}
	}
}

func Stringify(input []byte) (string, error) {
	data := make([]Node, 0)
	err := json.Unmarshal(input, &data)
	if err != nil {
		return "", err
	}

	var ret bytes.Buffer
	StringifyNodes(&ret, data)
	return ret.String(), nil
}
