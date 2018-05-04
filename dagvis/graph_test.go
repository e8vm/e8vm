package dagvis

import (
	"testing"

	"reflect"
)

func TestGraphReverse(t *testing.T) {
	g := NewGraph(map[string][]string{
		"a": {"b", "c"},
		"b": {"c"},
		"c": nil,
	})

	got := g.Reverse()
	want := NewGraph(map[string][]string{
		"c": {"a", "b"},
		"b": {"a"},
		"a": nil,
	})

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestGraphRemove(t *testing.T) {
	g := NewGraph(map[string][]string{
		"a": {"b", "c"},
		"b": {"c"},
		"c": nil,
	})

	got := g.Remove("b")
	want := NewGraph(map[string][]string{
		"a": {"c"},
		"c": nil,
	})

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
