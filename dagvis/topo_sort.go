package dagvis

// TopoSort topologically sorts the graph.
func TopoSort(g *Graph) ([]string, error) {
	m, err := NewMap(g)
	if err != nil {
		return nil, err
	}

	nodes := m.SortedNodes()
	ret := make([]string, 0, len(nodes))
	for _, node := range nodes {
		ret = append(ret, node.Name)
	}

	return ret, nil
}
