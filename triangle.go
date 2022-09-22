package triangle

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"strings"
)

const (
	firstNodeIndex = 1
)

var Debug = false

// Triangulation is input/output mesh
type Triangulation struct {
	Nodes     []Node
	Segments  []Segment
	Holes     []Node
	Triangles []Triangle
	Regions   []Node
}

// String return typical string result
func (t Triangulation) String() string {
	var out string
	for i, n := range t.Nodes {
		out += fmt.Sprintf("Node     %03d: {%+12e %+12e} %3d\n",
			i, n.X, n.Y, n.Marker)
	}
	for i, s := range t.Segments {
		out += fmt.Sprintf("Segment  %03d: {%03d %03d} %3d\n",
			i, s.NodeIndexes[0], s.NodeIndexes[1], s.Marker)
	}
	for i, n := range t.Holes {
		out += fmt.Sprintf("Hole     %03d: {%+12e %+12e} %3d\n",
			i, n.X, n.Y, n.Marker)
	}
	for i, t := range t.Triangles {
		out += fmt.Sprintf("Triangle %03d: {%03d %03d %03d} %3d\n",
			i, t.NodeIndexes[0], t.NodeIndexes[1], t.NodeIndexes[2], t.Marker)
	}
	for i, n := range t.Regions {
		out += fmt.Sprintf("Region   %03d: {%+12e %+12e} %3d\n",
			i, n.X, n.Y, n.Marker)
	}
	return out
}

// Triangle is triangle between 3 points
type Triangle struct {
	NodeIndexes [3]int
	Marker      int
}

// Node is 2D coordinate {X,Y}
type Node struct {
	X, Y   float64
	Marker int
}

// Segment is line between 2 points
type Segment struct {
	NodeIndexes [2]int
	Marker      int
}

// .node files
// First line:
//		<# of vertices> <dimension (must be 2)> <# of attributes> <# of boundary markers (0 or 1)>
// Remaining lines:
//		<vertex #> <x> <y> [attributes] [boundary marker]
//
// See: https://www.cs.cmu.edu/~quake/triangle.node.html
func (tr *Triangulation) createNodefile() (body string) {
	body += fmt.Sprintf("%d 2 0 1\n", len(tr.Nodes))
	for i := range tr.Nodes {
		body += fmt.Sprintf("%d %14.6e %14.6e %d\n",
			i+firstNodeIndex, tr.Nodes[i].X, tr.Nodes[i].Y, tr.Nodes[i].Marker)
	}
	return
}

// .poly files
//
// First line: <# of vertices> <dimension (must be 2)> <# of attributes> <# of boundary markers (0 or 1)>
// Following lines: <vertex #> <x> <y> [attributes] [boundary marker]
// One line: <# of segments> <# of boundary markers (0 or 1)>
// Following lines: <segment #> <endpoint> <endpoint> [boundary marker]
// One line: <# of holes>
// Following lines: <hole #> <x> <y>
// Optional line: <# of regional attributes and/or area constraints>
// Optional following lines: <region #> <x> <y> <attribute> <maximum area>
func (tr *Triangulation) createPolyfile() (body string) {
	body += tr.createNodefile() + "\n"

	body += fmt.Sprintf("%d 1\n", len(tr.Segments))
	for i := range tr.Segments {
		body += fmt.Sprintf("%d %d %d %d\n",
			i+firstNodeIndex,
			tr.Segments[i].NodeIndexes[0]+firstNodeIndex,
			tr.Segments[i].NodeIndexes[1]+firstNodeIndex,
			tr.Segments[i].Marker)
	}
	body += "\n"

	body += fmt.Sprintf("%d\n", len(tr.Holes))
	for i := range tr.Holes {
		body += fmt.Sprintf("%d %14.6e %14.6e\n",
			i+firstNodeIndex, tr.Holes[i].X, tr.Holes[i].Y)
	}

	body += fmt.Sprintf("%d\n", len(tr.Regions))
	for i := range tr.Regions {
		body += fmt.Sprintf("%d %14.6e %14.6e %d\n",
			i+firstNodeIndex, tr.Regions[i].X, tr.Regions[i].Y, tr.Regions[i].Marker)
	}

	return
}

func (tr *Triangulation) readPolyfile(filename string) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("cannot read .poly file: %v", err)
		}
	}()
	content, err := cleanAndRead(filename)
	if err != nil {
		return err
	}

	// TODO err = fmt.Errorf("not implemented")
	fmt.Println(string(content))

	return
}

func (tr *Triangulation) readNodefile(filename string) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("cannot read .node file: %v", err)
		}
	}()
	content, err := cleanAndRead(filename)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")

	if len(lines) == 0 {
		return fmt.Errorf("file is empty")
	}

	for i := range lines {
		if len(strings.TrimSpace(lines[i])) == 0 {
			continue
		}
		// First line
		if i == 0 {
			var size int
			var dimension, attributes, marker int
			_, err = fmt.Sscanf(lines[i], "%d %d %d %d", &size, &dimension, &attributes, &marker)
			if err != nil && err != io.EOF {
				return err
			}
			tr.Nodes = make([]Node, size)
			continue
		}
		// Next lines
		var position int
		var x, y float64
		var marker, attributes int
		n, err := fmt.Sscanf(lines[i], "%d %f %f %d %d", &position, &x, &y, &attributes, &marker)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 4 {
			marker = attributes
			attributes = 0
		}
		tr.Nodes[position-1] = Node{
			X:      x,
			Y:      y,
			Marker: marker,
		}
	}

	return nil
}

// .ele file
//
// First line:
//		<# of triangles> <nodes per triangle> <# of attributes>
// Remaining lines:
//		<triangle #> <node> <node> <node> ... [attributes]
func (tr *Triangulation) readElefile(filename string) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("cannot read .ele file `%s`: %v",
				filename, err)
		}
	}()
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s", string(debug.Stack()))
		}
	}()
	content, err := cleanAndRead(filename)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")

	if len(lines) == 0 {
		return fmt.Errorf("file is empty")
	}

	integer := func(value string) (int, error) {
		i64, err := strconv.ParseInt(value, 10, 64)
		return int(i64), err
	}

	var amountPoints int
	for i := range lines {
		if len(strings.TrimSpace(lines[i])) == 0 {
			continue
		}
		// first line
		if i == 0 {
			var size int
			var attributes int
			_, err = fmt.Sscanf(lines[i], "%d %d %d", &size, &amountPoints, &attributes)
			if err != nil && err != io.EOF {
				return err
			}
			tr.Triangles = make([]Triangle, size)
			continue
		}

		// next lines
		fs := strings.Fields(lines[i])
		var vs []int
		for i := range fs {
			var t int
			t, err = integer(fs[i])
			if err != nil {
				return
			}
			vs = append(vs, t)
		}
		tr.Triangles[vs[0]-1] = Triangle{
			NodeIndexes: [3]int{vs[1], vs[2], vs[3]},
		}
		if 4 < len(vs) {
			tr.Triangles[vs[0]-1].Marker = vs[4]
		}
	}
	return nil
}

func cleanAndRead(filename string) (content []byte, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("cannot clean file `%s`: %v", filename, err)
		}
	}()
	// read data
	content, err = ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	// remove comments from file
	lines := bytes.Split(content, []byte("\n"))
	content = content[:0] // clear slice
	for i := range lines {
		if len(bytes.TrimSpace(lines[i])) == 0 {
			// empty line
			continue
		}
		index := bytes.Index(lines[i], []byte("#"))
		if index >= 0 {
			lines[i] = lines[i][:index]
		}
		content = append(content, lines[i]...)
		if i != len(lines)-1 {
			content = append(content, '\n')
		}
	}

	return
}

func (tr *Triangulation) Run(flag string) error {
	// create temp directory
	dir, err := ioutil.TempDir("", "triangle")
	if err != nil {
		return err
	}
	if Debug {
		fmt.Fprintf(os.Stdout, "templorary dir = `%s`\n", dir)
	} else {
		defer os.RemoveAll(dir) // clean up
	}

	if len(tr.Segments) == 0 {
		var (
			nodefile = filepath.Join(dir, "mesh.node")
			content  = []byte(tr.createNodefile())
		)
		if err := ioutil.WriteFile(nodefile, content, 0666); err != nil {
			return err
		}
	}

	// polyfile
	var (
		polyfile = filepath.Join(dir, "mesh.poly")
		content  = []byte(tr.createPolyfile())
	)
	if err := ioutil.WriteFile(polyfile, content, 0666); err != nil {
		return err
	}
	if flag == "" {
		// flag = "-pq32.5a0.2ABPYXs"
		// flag = "-pq32.5AX"
		// flag = "-pA"
		// flag = "-pqa0.2AYs"
		flag = "-pqa0.2AYs"
		// flag = "-pcABeq0L"// a.05"
	}
	// execute Triangle
	cmd := exec.Command("triangle", flag, filepath.Join(dir, "mesh"))
	if err := cmd.Run(); err != nil {
		return err
	}
	// reading result files
	for _, err := range []error{
		tr.readNodefile(filepath.Join(dir, "mesh.1.node")),
		tr.readPolyfile(filepath.Join(dir, "mesh.1.poly")),
		tr.readElefile(filepath.Join(dir, "mesh.1.ele")),
	} {
		if err != nil {
			return err
		}
	}
	return nil
}
