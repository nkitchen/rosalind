package tree

import "bytes"
import "fmt"
import "io"
import "math"

type Node struct {
	Label string
	Children []Edge
}

type Edge struct {
	*Node
	Weight float64
}

func Print(t *Node) {
	printSubtree(Edge{t, 0}, "")
}

func printSubtree(t Edge, prefix string) {
	_ = `
	[a]-+-[b]
	    +-[c]-+-[d]-+-[d1]
		|     |     +-[d2]
		|     |     \-[d3]
		|     +-[e]-+-[e1]
		|     |     +-[e2]
		|     |     \-[e3]
		|     \-[f]
		\-[g]---[h]
	`
	fmt.Print("[", t.Label, "]")

	buf := make([]byte, len(t.Label))
	for i := range buf {
		buf[i] = ' '
	}
	labelSpace := string(buf)

	switch {
	case len(t.Children) == 0:
	    fmt.Println()
	case len(t.Children) == 1:
	    fmt.Print("---")
	    printSubtree(t.Children[0], prefix + " " + labelSpace + "    ")
    default:
		p := prefix + " " + labelSpace + "  | "
	    fmt.Print("-+-")
	    printSubtree(t.Children[0], p)
		var i int
		for i = 1; i < len(t.Children) - 1; i++ {
			fmt.Print(prefix, " ", labelSpace, "  +-")
			printSubtree(t.Children[i], p)
		}
		fmt.Print(prefix, " ", labelSpace, `  \-`)
		p = prefix + " " + labelSpace + "    "
		printSubtree(t.Children[i], p)
	}
}

// find returns the path from a node with the given label to the root.
func (t Edge) find(label string) []Edge {
	if t.Label == label {
		return []Edge{t}
	}

	for _, c := range t.Children {
		p := c.find(label)
		if len(p) > 0 {
			return append(p, t)
		}
	}
	return nil
}

// Distance returns the distance between two nodes in the tree
// that have the given labels.
// It returns -1 if either of the labels cannot be found.
func (t *Node) Distance(a, b string) int {
	e := Edge{t, 0}
	p := e.find(a)
	q := e.find(b)
	if len(p) == 0 || len(q) == 0 {
		return -1
	}

	if len(p) < len(q) {
		p, q = q, p
	}

	for len(q) > 0 && p[len(p) - 1] == q[len(q) - 1] {
		p = p[:len(p) - 1]
		q = q[:len(q) - 1]
	}

	return len(q) + len(p)
}

// WeightedDistance returns the sum of the weights on the edges
// between two nodes in the tree that have the given labels.
// It returns +infinity if either of the labels cannot be found.
func (t *Node) WeightedDistance(a, b string) float64 {
	e := Edge{t, 0}
	p := e.find(a)
	q := e.find(b)
	if len(p) == 0 || len(q) == 0 {
		return math.Inf(1)
	}

	if len(p) < len(q) {
		p, q = q, p
	}

	for len(q) > 0 && p[len(p) - 1] == q[len(q) - 1] {
		p = p[:len(p) - 1]
		q = q[:len(q) - 1]
	}

	d := float64(0)
	for _, e := range p {
		d += e.Weight
	}
	for _, e := range q {
		d += e.Weight
	}
	return d
}

func (t *Node) String() string {
	b := &bytes.Buffer{}
	t.WriteNewick(b)
	s := b.String()
	// Omit the final semicolon.
	return s[:len(s) - 1]
}

// Edge returns a zero-weight edge for a node.
func (t *Node) Edge() Edge {
	return Edge{t, 0}
}

func (t *Node) WriteNewick(w io.Writer) error {
	err := t.writeNewickSubtree(w)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w, ";")
	return err
}

func (t *Node) writeNewickSubtree(w io.Writer) error {
	if len(t.Children) > 0 {
		_, err := fmt.Fprintf(w, "(")
		if err != nil {
			return err
		}
		for i, e := range t.Children {
			if i != 0 {
				_, err = fmt.Fprintf(w, ",")
				if err != nil {
					return err
				}
			}
			err := e.Node.writeNewickSubtree(w)
			if err != nil {
				return err
			}
			if e.Weight != 0 {
				_, err = fmt.Fprintf(w, ":%v", e.Weight)
				if err != nil {
					return err
				}
			}
		}
		_, err = fmt.Fprintf(w, ")")
		if err != nil {
			return err
		}
	}

	_, err := fmt.Fprintf(w, "%s", t.Label)
	return err
}

// SubtreeLeaves returns a map from each node to the indices of the leaves
// contained under it.
// The slices for all the nodes' leaves shared the same backing array,
// so the total storage of leaves requires only O(N) space.
func (t *Node) SubtreeLeaves(leafIndices map[string]int) map[*Node][]int {
	leaves := make([]int, len(leafIndices))
	result := map[*Node][]int{}
	collectLeaves(t, leafIndices, leaves, result)
	return result
}

func collectLeaves(t *Node, leafIndices map[string]int,
                   collected []int, result map[*Node][]int) int {
	n := 0
	i, ok := leafIndices[t.Label]
	if ok {
		collected[0] = i
		n++
	}

	for _, child := range t.Children {
		n += collectLeaves(child.Node, leafIndices, collected[n:], result)
	}
	result[t] = collected[:n]
	return n
}

// Collapse returns a new tree with leaves with the labels in leavesToKeep
// and the same topology between these as in t, but with t's other leaves omitted.
func Collapse(t *Node, leavesToKeep map[string]bool) *Node {
	if len(t.Children) == 0 {
		if leavesToKeep[t.Label] {
			return &Node{t.Label, nil}
		}
		return nil
	}

	a := make([]Edge, 0, len(t.Children))
	for _, child := range t.Children {
		s := Collapse(child.Node, leavesToKeep)
		if s != nil {
			a = append(a, Edge{s, 0})
		}
	}

	var r *Node
	switch len(a) {
	case 0:
		r = nil
	case 1:
		r = a[0].Node
	default:
		r = &Node{"", a}
	}
	return r
}

// NormalUnrootedBinary returns a tree whose root has the expected structure
// for unrooted binary trees: It's an internal node with three children
// (unless the tree has only two nodes).
// The tree returned shares as many nodes with t as possible without changing
// the data in t's nodes.
func NormalUnrootedBinary(t *Node) *Node {
	if t.Label != "" {
		return t
	}

	if len(t.Children) == 2 {
		a := t.Children[0]
		b := t.Children[1]
		leaf1 := len(a.Children) == 0
		leaf2 := len(b.Children) == 0
		switch {
		case leaf1 && leaf2:
			return &Node{a.Label, []Edge{b}}
		case leaf1 && !leaf2:
			return &Node{"", append([]Edge{a}, b.Children...)}
		case !leaf1:
			c := append([]Edge{}, a.Children...)
			return &Node{"", append(c, b)}
		}
	}

	return t
}

// CollapseUnrootedBinary returns a new unrooted binary tree
// with leaves with the labels in leavesToKeep
// and the same topology between these as in t,
// but with t's other leaves omitted.
func CollapseUnrootedBinary(t *Node, leavesToKeep map[string]bool) *Node {
	return NormalUnrootedBinary(Collapse(t, leavesToKeep))
}
