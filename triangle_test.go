package triangle

import (
	"fmt"
	"testing"
)

func TestSpiral(t *testing.T) {
	// https://www.cs.cmu.edu/~quake/triangle.delaunay.html
	//
	// https://www.cs.cmu.edu/~quake/spiral.node
	//
	// # spiral.node
	// #
	// # A set of fifteen points in 2D, no attributes, no boundary markers.
	// 15  2  0  0
	// # And here are the fifteen points.
	//  1      0       0
	//  2     -0.416   0.909
	//  3     -1.35    0.436
	//  4     -1.64   -0.549
	//  5     -1.31   -1.51
	//  6     -0.532  -2.17
	//  7      0.454  -2.41
	//  8      1.45   -2.21
	//  9      2.29   -1.66
	// 10      2.88   -0.838
	// 11      3.16    0.131
	// 12      3.12    1.14
	// 13      2.77    2.08
	// 14      2.16    2.89
	// 15      1.36    3.49

	mesh := Triangulation{
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
	}

	fmt.Println(mesh)
	err := Triangulate(&mesh)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(mesh)
}

func TestPoly(t *testing.T) {
	// https://www.cs.cmu.edu/~quake/triangle.delaunay.html
	//
	// https://www.cs.cmu.edu/~quake/box.poly
	//
	// # A box with eight points in 2D, no attributes, one boundary marker.
	// 8 2 0 1
	// # Outer box has these vertices:
	//  1   0 0   0
	//  2   0 3   0
	//  3   3 0   0
	//  4   3 3   33     # A special marker for this point.
	// # Inner square has these vertices:
	//  5   1 1   0
	//  6   1 2   0
	//  7   2 1   0
	//  8   2 2   0
	// # Five segments with boundary markers.
	// 5 1
	//  1   1 2   5      # Left side of outer box.
	//  2   5 7   0      # These four segments enclose the hole.
	//  3   7 8   0
	//  4   8 6   10
	//  5   6 5   0
	// # One hole in the middle of the inner square.
	// 1
	//  1   1.5 1.5

	mesh := Triangulation{
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
	}

	fmt.Println(mesh)
	err := Triangulate(&mesh)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(mesh)
}

func TestPolySquare(t *testing.T) {
	mesh := Triangulation{
		Nodes: []Node{
			{1, 0, 0},
			{0, 1, 0},
			{-1, 0, 0},
			{0, -1, 0},
			{0, 0, 0},
		},
		Segments: []Segment{
			{0, 1, 0},
			{1, 2, 0},
			{2, 3, 0},
			{3, 0, 0},
		},
		Holes: []Node{},
	}

	fmt.Println(mesh)
	err := Triangulate(&mesh)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(mesh)
}
