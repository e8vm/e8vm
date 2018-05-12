package dagvis

// CheckDAG checks if a graph is a valid DAG.  It returns true
// when all the graph, links are valid and has no circular
// dependency.
func CheckDAG(g *Graph) error {
	m, err := initMap(g)
	if err != nil {
		return err
	}

	_, err = m.makeLayers()
	return err
}
