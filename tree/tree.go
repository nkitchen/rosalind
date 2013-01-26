package tree

import "fmt"

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
		printSubtree(t.Children[i], p)
	}
}

// find returns the path from a node with the given label to the root.
func (t *Node) find(label string) []*Node {
	if t.Label == label {
		return []*Node{t}
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
	p := t.find(a)
	q := t.find(b)
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

	if len(q) == 0 {
		// q was a prefix of p
		return len(p)
	}

	return len(q) + len(p)
}

func (t *Node) String() string {
	return t.Label
}
