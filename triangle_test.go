package triangle

import (
	"fmt"
	"os"
	"testing"
)

func Test(t *testing.T) {
	Debug = true

	// https://www.cs.cmu.edu/~quake/triangle.delaunay.html
	tcs := []struct {
		name string
		mesh Triangulation
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
				{1.45, -2.21, 4},
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
				{NodeIndexes: [2]int{0, 1}, Marker: 5},
				{NodeIndexes: [2]int{4, 6}, Marker: 0},
				{NodeIndexes: [2]int{6, 7}, Marker: 0},
				{NodeIndexes: [2]int{7, 5}, Marker: 10},
				{NodeIndexes: [2]int{5, 4}, Marker: 0},
			},
			Holes: []Node{
				{1.5, 1.5, 0},
			},
		},
	}, {
		name: "square",
		mesh: Triangulation{
			Nodes: []Node{
				{1, 0, 0},   // 0
				{0, 1, 0},   // 1
				{-1, 0, 10}, // 2
				{0, -1, 0},  // 3
				{0, 0, 0},   // 4
			},
			Segments: []Segment{
				{NodeIndexes: [2]int{0, 1}, Marker: 2},
				{NodeIndexes: [2]int{1, 2}, Marker: 0},
				{NodeIndexes: [2]int{2, 3}, Marker: 5},
				{NodeIndexes: [2]int{3, 0}, Marker: 0},
			},
			Holes: []Node{},
		},
	}, {
		name: "SquareInSquare",
		mesh: Triangulation{
			Nodes: []Node{
				{1, 0, 0},   // 0
				{0, 1, 0},   // 1
				{-1, 0, 11}, // 2
				{0, -1, 0},  // 3
				{0, 0, 0},   // 4
				{2, 0, 0},   // 5
				{0, 2, 20},  // 6
				{-2, 0, 0},  // 7
				{0, -2, 0},  // 8
			},
			Segments: []Segment{
				{NodeIndexes: [2]int{0, 1}, Marker: 0},
				{NodeIndexes: [2]int{1, 2}, Marker: 0},
				{NodeIndexes: [2]int{2, 3}, Marker: 0},
				{NodeIndexes: [2]int{3, 0}, Marker: 0},
				{NodeIndexes: [2]int{5, 6}, Marker: 5},
				{NodeIndexes: [2]int{6, 7}, Marker: 0},
				{NodeIndexes: [2]int{7, 8}, Marker: 0},
				{NodeIndexes: [2]int{8, 5}, Marker: 0},
			},
			Holes: []Node{},
			Regions: []Node{
				{0, 0, 2},
				{0, -1.8, 9},
			},
		},
	}}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			if testing.Verbose() {
				fmt.Fprintf(os.Stdout, "In :\n%s\n", tc.mesh)
			}
			err := tc.mesh.Run("")
			if err != nil {
				t.Fatal(err)
			}
			if testing.Verbose() {
				fmt.Fprintf(os.Stdout, "Out:\n%s\n", tc.mesh)
			}
		})
	}
}
