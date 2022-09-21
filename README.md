# triangle
Go interface of triangulation C Triangle
```

package triangle // import "github.com/Konstantin8105/triangle"


VARIABLES

var Debug = false

TYPES

type Node struct {
	X, Y   float64
	Marker int
}
    Node is 2D coordinate {X,Y}

type Segment struct {
	NodeIndexes [2]int
	Marker      int
}
    Segment is line between 2 points

type Triangle struct {
	NodeIndexes [3]int
	Marker      int
}
    Triangle is triangle between 3 points

type Triangulation struct {
	Nodes     []Node
	Segments  []Segment
	Holes     []Node
	Triangles []Triangle
	Regions   []Node
}
    Triangulation is input/output mesh

func (tr *Triangulation) Run(flag string) error

func (t Triangulation) String() string
    String return typical string result

```
