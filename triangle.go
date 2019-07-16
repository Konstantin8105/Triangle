package triangle

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	firstNodeIndex = 1
)

type Triangulation struct {
	Nodes    []Node
	Segments []Segment
	Holes    []Node
	Triangle [][]int
}

type Node struct {
	X, Y   float64
	Marker int
}

type Segment struct {
	N1, N2 int
	Marker int
}

// .node files
// First line:
//		<# of vertices> <dimension (must be 2)> <# of attributes> <# of boundary markers (0 or 1)>
// Remaining lines:
//		<vertex #> <x> <y> [attributes] [boundary marker]
//
// See: https://www.cs.cmu.edu/~quake/triangle.node.html
func (tr *Triangulation) createNodefile() (body string) {
	body += fmt.Sprintf("%d 2 0 0\n", len(tr.Nodes))
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
			tr.Segments[i].N1+firstNodeIndex,
			tr.Segments[i].N2+firstNodeIndex,
			tr.Segments[i].Marker)
	}
	body += "\n"

	body += fmt.Sprintf("%d\n", len(tr.Holes))
	for i := range tr.Holes {
		body += fmt.Sprintf("%d %14.6e %14.6e\n",
			i+firstNodeIndex, tr.Holes[i].X, tr.Holes[i].Y)
	}

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
			err = fmt.Errorf("cannot read .ele file: %v", err)
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

	var amountPoints int
	for i := range lines {
		if len(strings.TrimSpace(lines[i])) == 0 {
			continue
		}
		// First line
		if i == 0 {
			var size int
			var attributes int
			_, err = fmt.Sscanf(lines[i], "%d %d %d", &size, &amountPoints, &attributes)
			if err != nil && err != io.EOF {
				return err
			}
			tr.Triangle = make([][]int, size)
			for i := 0; i < size; i++ {
				tr.Triangle[i] = make([]int, amountPoints)
			}
			continue
		}
		// Next lines
		var position int
		items := strings.Split(lines[i], " ")
		var counter int
		for j := range items {
			item := strings.TrimSpace(items[j])
			if len(item) == 0 {
				continue
			}
			counter++
			var value int
			n, err := fmt.Sscanf(item, "%d", &value)
			if err != nil && err != io.EOF {
				return err
			}
			if n == 0 {
				return fmt.Errorf("n == 0")
			}
			if counter == 1 {
				position = value
				continue
			}
			tr.Triangle[position-1][counter-2] = value
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

func Triangulate(tr *Triangulation) error {
	// create temp directory
	dir, err := ioutil.TempDir("", "triangle")
	if err != nil {
		return err
	}
	// TODO : defer os.RemoveAll(dir) // clean up

	var flag string

	if len(tr.Segments) == 0 {
		var (
			nodefile = filepath.Join(dir, "mesh.node")
			content  = []byte(tr.createNodefile())
		)
		if err := ioutil.WriteFile(nodefile, content, 0666); err != nil {
			return err
		}
	} else {
		// Polyline
		var (
			polyfile = filepath.Join(dir, "mesh.poly")
			content  = []byte(tr.createPolyfile())
		)
		if err := ioutil.WriteFile(polyfile, content, 0666); err != nil {
			return err
		}
		// 		flag = "-pq0L"
		flag = "-pq0La.25"
	}

	fmt.Println(flag) // TODO: remove

	// execute Triangle
	cmd := exec.Command("triangle", flag, filepath.Join(dir, "mesh"))
	if err := cmd.Run(); err != nil {
		return err
	}

	fmt.Println(dir) // TODO: remove

	// read .node file
	nodefile := filepath.Join(dir, "mesh.1.node")
	err = tr.readNodefile(nodefile)
	if err != nil {
		return err
	}

	// read .ele file
	elefile := filepath.Join(dir, "mesh.1.ele")
	err = tr.readElefile(elefile)
	if err != nil {
		return err
	}

	return nil
}
