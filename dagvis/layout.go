package dagvis

import (
	"sort"
)

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
		sum += in.y
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
	y := n.y
	if tak[y-2] && !tak[y-1] {
		n.y--
	} else if !tak[y-1] && tak[y+2] && !tak[y+1] {
		n.y++
	}
}

func makeNodeList(m map[string]*MapNode) []string {
	var ret []string
	for name := range m {
		ret = append(ret, name)
	}
	sort.Strings(ret)
	return ret
}

// LayoutMap creates the MapView for the given Map.
func LayoutMap(m *Map) *MapView {
	pushTight(m) // push it tight

	v := &MapView{
		Nodes: make(map[string]*MapNodeView),
	}

	layers := m.SortedLayers()
	slotTaken := make([]map[int]bool, m.Nlayer)
	for i := range slotTaken {
		slotTaken[i] = make(map[int]bool)
	}

	ymin := 0
	for _, layer := range layers {
		for _, node := range layer {
			x := node.layer

			node.x = x

			tak := slotTaken[x]
			node.y = findY(node, tak)
			snapNearBy(node, tak)

			y := node.y
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
		node.y -= ymin
		if node.y > ymax {
			ymax = node.y
		}

		v.Nodes[node.Name] = &MapNodeView{
			Name:     node.Name,
			X:        node.x,
			Y:        node.y,
			CritIns:  makeNodeList(node.CritIns),
			CritOuts: makeNodeList(node.CritOuts),
		}
	}

	v.Width = m.Nlayer
	v.Height = ymax + 1

	return v
}

// Layout layouts a DAG into a map.
func Layout(g *Graph) (*Map, *MapView, error) {
	m, err := NewMap(g) // build the map
	if err != nil {
		return nil, nil, err
	}

	v := LayoutMap(m) // assign coordinates

	return m, v, nil
}

// LayoutJSON layouts a DAG into a map in json format.
func LayoutJSON(g *Graph) ([]byte, error) {
	_, v, err := Layout(g)
	if err != nil {
		return nil, err
	}

	return marshalMap(v), nil
}

// RevLayout layouts a DAG into a map from right to left
// its more suitable for top-down designed projects.
func RevLayout(g *Graph) (*Map, *MapView, error) {
	m, v, err := Layout(g.Reverse())
	if err != nil {
		return nil, nil, err
	}
	m.Reverse()
	v.Reverse()
	return m, v, nil
}

// RevLayoutJSON layouts a DAG into a map in json format.
func RevLayoutJSON(g *Graph) ([]byte, error) {
	_, v, err := RevLayout(g)
	if err != nil {
		return nil, err
	}

	return marshalMap(v), nil
}
