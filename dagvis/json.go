package dagvis

import (
	"encoding/json"
)

// N is a node in the minified DAG visualization result.
type N struct {
	N    string   `json:"n"`
	X    int      `json:"x"`
	Y    int      `json:"y"`
	Ins  []string `json:"i"`
	Outs []string `json:"o"`
}

// M is a node in the minified DAG visualization result.
type M struct {
	Height int           `json:"h"`
	Width  int           `json:"w"`
	Nodes  map[string]*N `json:"n"`
}

// Output returns a json'able object of a map.
func Output(v *MapView) *M {
	res := &M{
		Height: v.Height,
		Width:  v.Width,
		Nodes:  make(map[string]*N),
	}

	for name, node := range v.Nodes {
		display := node.DisplayName
		if display == "" {
			display = name
		}

		n := &N{
			N:    display,
			X:    node.X,
			Y:    node.Y,
			Ins:  node.CritIns,
			Outs: node.CritOuts,
		}

		res.Nodes[name] = n
	}

	return res
}

func marshalMap(m *MapView) []byte {
	res := Output(m)
	ret, e := json.MarshalIndent(res, "", "    ")
	if e != nil {
		panic(e)
	}

	return ret
}
