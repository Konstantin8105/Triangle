package triangle

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	// https://www.cs.cmu.edu/~quake/triangle.delaunay.html
	tcs := []struct {
		name string
		mesh    Triangulation
	}{{
		// https://www.cs.cmu.edu/~quake/spiral.node
		name: "spiral",
		mesh: Triangulation{
			Nodes: []Node{
				{0, 0, 0},
				{-0.416, 0.909, 0},
				{-1.35, 0.436, 0},
				{-1.64, -0.549, 0},
				{-1.31, -1.51, 0},
				{-0.532, -2.17, 0},
				{0.454, -2.41, 0},
				{1.45, -2.21, 0},
				{2.29, -1.66, 0},
				{2.88, -0.838, 0},
				{3.16, 0.131, 0},
				{3.12, 1.14, 0},
				{2.77, 2.08, 0},
				{2.16, 2.89, 0},
				{1.36, 3.49, 0},
			},
		},
	}, {
		// https://www.cs.cmu.edu/~quake/box.poly
		name: "box",
		mesh: Triangulation{
			Nodes: []Node{
				{0, 0, 0},
				{0, 3, 0},
				{3, 0, 0},
				{3, 3, 33},
				{1, 1, 0},
				{1, 2, 0},
				{2, 1, 0},
				{2, 2, 0},
			},
			Segments: []Segment{
				{0, 1, 5},
				{4, 6, 0},
				{6, 7, 0},
				{7, 5, 10},
				{5, 4, 0},
			},
			Holes: []Node{
				{1.5, 1.5, 0},
			},
		},
	}, {
		name: "square",
		mesh: Triangulation{
			Nodes: []Node{
				{1, 0, 0},  // 0
				{0, 1, 0},  // 1
				{-1, 0, 1}, // 2
				{0, -1, 0}, // 3
				{0, 0, 0},  // 4
			},
			Segments: []Segment{
				{0, 1, 2},
				{1, 2, 0},
				{2, 3, 5},
				{3, 0, 0},
			},
			Holes: []Node{},
		},
	}}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Println(tc.mesh)
			err := Triangulate(&tc.mesh)
			if err != nil {
				t.Fatal(err)
			}
			fmt.Println(tc.mesh)
		})
	}
}
