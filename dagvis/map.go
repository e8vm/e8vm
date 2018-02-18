package dagvis

import (
	"bytes"
	"errors"
	"fmt"
	"sort"
)

// Map is a visualized DAG
type Map struct {
	Nodes map[string]*MapNode

	Nedge  int
	Ncrit  int
	Nlayer int
}

// Reverse reverses the map.
func (m *Map) Reverse() {
	for _, node := range m.Nodes {
		node.Ins, node.Outs = node.Outs, node.Ins
		node.AllIns, node.AllOuts = node.AllOuts, node.AllIns
		node.CritIns, node.CritOuts = node.CritOuts, node.CritIns
		node.layer = m.Nlayer - 1 - node.layer
	}
}

func initMap(g *Graph) (*Map, error) {
	ret := new(Map)
	ret.Nodes = make(map[string]*MapNode)

	// create the nodes
	for name := range g.Nodes {
		ret.Nodes[name] = newMapNode(name)
	}

	// connect the links
	ret.Nedge = 0
	for in, outs := range g.Nodes {
		inNode := ret.Nodes[in]
		if inNode == nil {
			panic("bug")
		}

		for _, out := range outs {
			outNode, found := ret.Nodes[out]
			if !found {
				err := fmt.Errorf("missing node %q for %q", out, in)
				return nil, err
			}

			outNode.Ins[in] = inNode
			inNode.Outs[out] = outNode

			ret.Nedge++
		}
	}

	return ret, nil
}

// NewMap creates a map from a graph where all the ins and outs
// are populated, and nodes are mapped into layers.
func NewMap(g *Graph) (*Map, error) {
	ret, err := initMap(g)
	if err != nil {
		return nil, err
	}

	// make them into layers
	layers, err := ret.makeLayers()
	if err != nil {
		return nil, err
	}

	ret.Nlayer = len(layers)

	// propogate all ins/outs
	ret.buildAlls(layers)

	// compute the critical dependencies
	ret.buildCrits()

	return ret, nil
}

func (m *Map) makeLayers() ([][]*MapNode, error) {
	var ret [][]*MapNode
	var cur []*MapNode
	left := make(map[*MapNode]struct{})

	for _, node := range m.Nodes {
		left[node] = struct{}{}
	}

	for _, node := range m.Nodes {
		if len(node.Ins) == 0 {
			cur = append(cur, node)
		}
		node.nhit = 0
	}

	n := 0

	for len(cur) > 0 {
		for _, node := range cur {
			node.layer = len(ret)
			delete(left, node)
		}

		ret = append(ret, cur)
		n += len(cur)

		var next []*MapNode
		for _, node := range cur {
			for _, out := range node.Outs {
				out.nhit++
				if out.nhit == len(out.Ins) {
					next = append(next, out)
				}
			}
		}

		cur = next
	}

	if len(left) != 0 {
		circle := minCircle(m.Nodes)
		if len(circle) == 0 {
			panic("should find a circle")
		}

		msg := new(bytes.Buffer)
		fmt.Fprintf(msg, "graph has circle: ")
		for i, node := range circle {
			if i != 0 {
				fmt.Fprintf(msg, "->")
			}
			fmt.Fprintf(msg, node.Name)
		}
		return nil, errors.New(msg.String())
	}

	return ret, nil
}

func (m *Map) buildAlls(layers [][]*MapNode) {
	for _, layer := range layers {
		for _, node := range layer {
			for _, out := range node.Outs {
				// propagate all incoming nodes into the out node
				for _, in := range node.AllIns {
					out.AllIns[in.Name] = in
					in.AllOuts[out.Name] = out
				}

				// connect this edge as well
				out.AllIns[node.Name] = node
				node.AllOuts[out.Name] = out
			}
		}
	}
}

func isCrit(from, to *MapNode) bool {
	for _, via := range from.AllOuts {
		if via == to {
			continue
		}

		if via.AllOuts[to.Name] != nil {
			return false
		}
	}

	return true
}

func (m *Map) buildCrits() {
	m.Ncrit = 0

	for _, node := range m.Nodes {
		for _, out := range node.Outs {
			if !isCrit(node, out) {
				continue
			}

			node.CritOuts[out.Name] = out
			out.CritIns[node.Name] = node
			m.Ncrit++
		}
	}
}

// SortedNodes returns a nodes that are sorted in topological order.
func (m *Map) SortedNodes() []*MapNode {
	var ret mapNodes

	for _, node := range m.Nodes {
		ret = append(ret, node)
	}

	sort.Sort(byLayer{ret})
	return ret
}

// SortedLayers returns nodes that are sorted in topological layers.
// Each node resides in its lowest possible layer.
func (m *Map) SortedLayers() [][]*MapNode {
	ret := make([][]*MapNode, m.Nlayer)

	for _, node := range m.Nodes {
		ret[node.layer] = append(ret[node.layer], node)
	}

	for _, layer := range ret {
		sort.Sort(byNcritOuts{layer})
	}

	return ret
}
