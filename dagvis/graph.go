// Package dagvis visualizes a DAG graph into
// a structured, layered planer map.
package dagvis

import (
	"sort"
)

// Graph is a directed graph
type Graph struct {
	Nodes map[string][]string
}

// NewGraph create a new graph with the given nodes.
func NewGraph(nodes map[string][]string) *Graph {
	return &Graph{Nodes: nodes}
}

// Reverse the graph
func (g *Graph) Reverse() *Graph {
	ret := make(map[string][]string)

	for n := range g.Nodes {
		ret[n] = nil // touch every node first
	}

	for n, lst := range g.Nodes {
		for _, m := range lst {
			ret[m] = append(ret[m], n)
		}
	}

	for _, list := range ret {
		sort.Strings(list)
	}

	return &Graph{Nodes: ret}
}

// Remove removes a node in a graph and returns the new graph.
func (g *Graph) Remove(node string) *Graph {
	ret := make(map[string][]string)
	for k, vs := range g.Nodes {
		if k == node {
			continue
		}

		var outs []string
		for _, v := range vs {
			if v == node {
				continue
			}
			outs = append(outs, v)
		}
		ret[k] = outs
	}
	return &Graph{Nodes: ret}
}

// SubGraph returns a sub graph of the current graph, which
// only contains nodes that passes the filter.
func (g *Graph) SubGraph(f func(string) bool) *Graph {
	hits := make(map[string]bool)
	for k := range g.Nodes {
		if f(k) {
			hits[k] = true
		}
	}

	ret := make(map[string][]string)
	for k, vs := range g.Nodes {
		if !hits[k] {
			continue
		}
		var outs []string
		for _, v := range vs {
			if !hits[v] {
				continue
			}
			outs = append(outs, v)
		}
		ret[k] = outs
	}
	return &Graph{Nodes: ret}
}

// Rename renames the name of each node in the graph
func (g *Graph) Rename(f func(string) (string, error)) (*Graph, error) {
	if f == nil {
		panic("rename function is nil")
	}

	ret := new(Graph)
	ret.Nodes = make(map[string][]string)

	for k, vs := range g.Nodes {
		newK, e := f(k)
		if e != nil {
			return nil, e
		}

		newVs := make([]string, 0, len(vs))
		for _, v := range vs {
			newV, e := f(v)
			if e != nil {
				return nil, e
			}

			newVs = append(newVs, newV)
		}

		ret.Nodes[newK] = newVs
	}

	return ret, nil
}
