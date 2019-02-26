package markdown_to_slate

func processChecks(prnodes *[]Node) bool {
	allChecks := true
	for i, node := range *prnodes {
		if len(node.Nodes) == 0 {
			allChecks = false
			continue
		}

		var leaf *Leaf
		leafDepth := 0
		if len(node.Nodes[0].Leaves) == 0 {
			if len(node.Nodes[0].Nodes) == 0 {
				allChecks = false
				continue
			}
			if len(node.Nodes[0].Nodes[0].Leaves) == 0 {
				allChecks = false
				continue
			} else {
				leaf = &node.Nodes[0].Nodes[0].Leaves[0]
				leafDepth = 2
			}
		} else {
			leaf = &node.Nodes[0].Leaves[0]
			leafDepth = 1
		}

		tx := leaf.Text
		if len(tx) < 4 {
			allChecks = false
			continue
		}
		prefix := tx[0:4]
		//log.Println("process check", tx, prefix)
		if len(node.Nodes) > 0 && (prefix == "[ ] " || prefix == "[x] ") {
			leaf.Text = tx[4:len(tx)]
			if leaf.Text == "" {
				if leafDepth == 1 {
					if len(node.Nodes[0].Leaves) == 1 {
						node.Nodes[0].Leaves = nil
					}
				} else if leafDepth == 2 {
					if len(node.Nodes[0].Nodes[0].Leaves) == 1 {
						node.Nodes[0].Nodes[0].Leaves = nil
					}
				}
			}

			node.Type = "check-list-item"

			if node.Data == nil {
				node.Data = make(map[string]interface{}, 0)
			}
			checked := (prefix == "[x] ")
			node.Data["checked"] = checked

			//log.Println("checks replaces node")
			(*prnodes)[i] = node
		} else {
			allChecks = false
		}
	}

	return allChecks
}
