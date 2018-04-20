package dagvis

import (
	"math/rand"
	"strconv"
	"testing"
)

func BenchmarkFindCircle(b *testing.B) {
	for n := 0; n < b.N; n++ {
		benchmarkFindCircle()
	}
}

func benchmarkFindCircle() []*MapNode {
	ret := make(map[string][]string)

	for i := 0; i < 1000; i++ {
		var edge []string
		for j := 0; j < 1000; j++ {
			if i == j {
				continue
			}
			q := rand.Int31n(500)
			if q < 1 {
				var flag bool
				for _, ele := range ret[strconv.Itoa(j)] {
					if ele == strconv.Itoa(i) {
						flag = true
					}
				}
				if !flag {
					edge = append(edge, strconv.Itoa(j))
				}
			}
		}
		ret[strconv.Itoa(i)] = edge
	}

	g := &Graph{Nodes: ret}
	nodes, _ := initMap(g)

	return minCircle(nodes.Nodes)
}
