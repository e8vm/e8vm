package dagvis

// MapNodeView is the position of a map node.
type MapNodeView struct {
	Name        string
	DisplayName string
	X, Y        int
	CritIns     []string
	CritOuts    []string
}

// MapView is a layout view of a Map.
type MapView struct {
	Height    int
	Width     int
	Nodes     map[string]*MapNodeView
	IsTopDown bool
}

// AssignDisplayName assigns display names in a map.
func (v *MapView) AssignDisplayName(f func(s string) string) {
	for _, n := range v.Nodes {
		n.DisplayName = f(n.Name)
	}
}

// Reverse reverse a map view.
func (v *MapView) Reverse() {
	for _, node := range v.Nodes {
		node.X = v.Width - 1 - node.X
		node.CritIns, node.CritOuts = node.CritOuts, node.CritIns
	}

	v.IsTopDown = !v.IsTopDown
}
