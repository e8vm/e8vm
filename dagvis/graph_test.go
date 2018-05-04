package dagvis

import (
	"testing"

	"reflect"
	"strings"
)

func checkGraphEqual(t *testing.T, got *Graph, nodes map[string][]string) {
	want := NewGraph(nodes)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestGraphReverse(t *testing.T) {
	g := NewGraph(map[string][]string{
		"a": {"b", "c"},
		"b": {"c"},
		"c": nil,
	})

	got := g.Reverse()
	checkGraphEqual(t, got, map[string][]string{
		"c": {"a", "b"},
		"b": {"a"},
		"a": nil,
	})
}

func TestGraphRemove(t *testing.T) {
	g := NewGraph(map[string][]string{
		"a": {"b", "c"},
		"b": {"c"},
		"c": nil,
	})

	got := g.Remove("b")
	checkGraphEqual(t, got, map[string][]string{
		"a": {"c"},
		"c": nil,
	})
}

func TestSubGraph(t *testing.T) {
	g := NewGraph(map[string][]string{
		"a": {"b", "c"},
		"b": {"c"},
		"c": nil,
	})

	got := g.SubGraph(func(n string) bool {
		return n != "b"
	})

	checkGraphEqual(t, got, map[string][]string{
		"a": {"c"},
		"c": nil,
	})
}

func TestGraphRename(t *testing.T) {
	g := NewGraph(map[string][]string{
		"a": {"b", "c"},
		"b": {"c"},
		"c": nil,
	})

	got, err := g.Rename(func(name string) (string, error) {
		return strings.ToUpper(name), nil
	})
	if err != nil {
		t.Fatal(err)
	}

	checkGraphEqual(t, got, map[string][]string{
		"A": {"B", "C"},
		"B": {"C"},
		"C": nil,
	})
}
