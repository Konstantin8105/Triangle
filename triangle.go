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

type Triangulation struct {
	Nodes    []Node
	Triangle [][]int
}

type Node struct {
	X, Y   float64
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
			i+1, tr.Nodes[i].X, tr.Nodes[i].Y, tr.Nodes[i].Marker)
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

	var (
		nodefile = filepath.Join(dir, "mesh.node")
		content  = []byte(tr.createNodefile())
	)
	if err := ioutil.WriteFile(nodefile, content, 0666); err != nil {
		return err
	}

	// execute Triangle
	cmd := exec.Command("triangle", filepath.Join(dir, "mesh"))
	if err := cmd.Run(); err != nil {
		return err
	}

	fmt.Println(dir) // TODO: remove

	// read .node file
	nodefile = filepath.Join(dir, "mesh.1.node")
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
