package dagvis

// IsDAG checks if a graph is a valid DAG.
// It returns true when all the graph, links are valid and
// has no circular dependency.
func IsDAG(g *Graph) (bool, error) {
	m, e := initMap(g)
	if e != nil {
		return false, e
	}

	_, e = m.makeLayers()
	if e != nil {
		return false, e
	}

	return true, nil
}
