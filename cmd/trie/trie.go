package main

import "bufio"
import "fmt"
import "io"
import "os"

var nextId = 1

type Node struct {
	id int
	child map[byte]*Node
}

func NewNode() *Node {
	n := &Node{nextId, make(map[byte]*Node)}
	nextId++
	return n
}

func main() {
	root := NewNode()

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		root.Insert(s.Text())
	}

	root.DumpAdjacencies(os.Stdout)
}

func (n *Node) Insert(s string) {
	if len(s) == 0 {
		return
	}

	b := s[0]
	c, ok := n.child[b]
	if !ok {
		c = NewNode()
		n.child[b] = c
	}
	c.Insert(s[1:])
}

func (n *Node) DumpAdjacencies(w io.Writer) {
	for b, c := range n.child {
		fmt.Printf("%d %d %c\n", n.id, c.id, b)
		c.DumpAdjacencies(w)
	}
}
