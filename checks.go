package markdown_to_slate

func processChecks(prnodes *[]Node) bool {
	allChecks := true
	for i, node := range *prnodes {
		if len(node.Nodes) == 0 {
			allChecks = false
			continue
		}
		//log.Println("begin process check", node.Nodes[0])
		var leaf Leaf
		if len(node.Nodes[0].Leaves) == 0 {
			//allChecks = false
			if len(node.Nodes[0].Nodes) == 0 {
				allChecks = false
				continue
			}
			if len(node.Nodes[0].Nodes[0].Leaves) == 0 {
				allChecks = false
				continue
			} else {
				leaf = node.Nodes[0].Nodes[0].Leaves[0]
			}

		} else {
			leaf = node.Nodes[0].Leaves[0]
		}
		//log.Println("leaf:", leaf.Text)
		tx := leaf.Text
		if len(tx) < 4 {
			allChecks = false
			continue
		}
		prefix := tx[0:4]
		//log.Println("process check", tx, prefix)
		if len(node.Nodes) > 0 && (prefix == "[ ] " || prefix == "[x] ") {
			textNode := Node{
				Object: "text",
				Leaves: []Leaf{
					Leaf{
						Object: "leaf",
						Text:   tx[4:len(tx)],
						Marks:  []Mark{},
					},
				},
			}
			ci := Node{
				Object: "block",
				Type:   "check-list-item",
				Data:   map[string]interface{}{},
				Nodes:  []Node{textNode},
			}

			checked := (prefix == "[x] ")
			ci.Data["checked"] = checked

			//log.Println("checks replaces node")
			(*prnodes)[i] = ci
		} else {
			allChecks = false
		}
	}

	return allChecks
}
