package dagvis

func critOutMaxLayer(n *MapNode) int {
	ret := n.layer
	for _, out := range n.CritOuts {
		if out.layer > ret {
			ret = out.layer
		}
	}

	return ret
}

func avgCritInY(n *MapNode) int {
	nIn := len(n.CritIns)
	if nIn == 0 {
		return 0
	}

	sum := 0
	for _, in := range n.CritIns {
		sum += in.Y
	}

	return (sum + nIn/2) / nIn // round up
}

func findY(n *MapNode, tak map[int]bool) int {
	yavg := avgCritInY(n)

	offset := 0
	for {
		if !tak[yavg+offset] {
			return yavg + offset
		}
		if !tak[yavg-offset] {
			return yavg - offset
		}
		offset++
	}
}

func snapNearBy(n *MapNode, tak map[int]bool) {
	y := n.Y
	if tak[y-2] && !tak[y-1] {
		n.Y--
	} else if !tak[y-1] && tak[y+2] && !tak[y+1] {
		n.Y++
	}
}

func layout(m *Map) {
	layers := m.SortedLayers()
	slotTaken := make([]map[int]bool, m.Nlayer)
	for i := range slotTaken {
		slotTaken[i] = make(map[int]bool)
	}

	ymin := 0
	for _, layer := range layers {
		for _, node := range layer {
			x := node.layer

			node.X = x

			tak := slotTaken[x]
			node.Y = findY(node, tak)
			snapNearBy(node, tak)

			y := node.Y
			tak[y-1] = true
			tak[y] = true
			tak[y+1] = true

			xmax := critOutMaxLayer(node)
			for i := x + 1; i < xmax; i++ {
				slotTaken[i][y] = true
			}

			if y < ymin {
				ymin = y
			}
		}
	}

	ymax := 0
	for _, node := range m.Nodes {
		node.Y -= ymin
		if node.Y > ymax {
			ymax = node.Y
		}
	}

	m.Width = m.Nlayer
	m.Height = ymax + 1
	m.IsTopDown = false
}

// Layout layouts a DAG into a map.
func Layout(g *Graph) (*Map, error) {
	m, e := NewMap(g) // build the map
	if e != nil {
		return nil, e
	}

	pushTight(m) // push it tight
	layout(m)    // assign coordinates

	return m, nil
}

// LayoutJSON layouts a DAG into a map in json format.
func LayoutJSON(g *Graph) ([]byte, error) {
	m, e := Layout(g)
	if e != nil {
		return nil, e
	}

	return marshalMap(m), nil
}

// RevLayout layouts a DAG into a map from right to left
// its more suitable for top-down designed projects.
func RevLayout(g *Graph) (*Map, error) {
	m, e := Layout(g.Reverse())
	if e != nil {
		return nil, e
	}
	m.Reverse()
	return m, nil
}

// RevLayoutJSON layouts a DAG into a map in json format.
func RevLayoutJSON(g *Graph) ([]byte, error) {
	m, e := RevLayout(g)
	if e != nil {
		return nil, e
	}

	return marshalMap(m), nil
}
