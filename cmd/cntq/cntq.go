package main

import "bufio"
import "fmt"
import "os"
import "rosalind/tree"

const M = 1000000

var leaves := map[*tree.Node]int64
var pairs := map[*tree.Node]int64

func main() {
	br := bufio.NewReader(os.Stdin)

	_, _ = br.ReadString('\n')

    t, _ := tree.ReadNewick(br)

	countLeaves(t)
	countPairs(t)
	e := tree.Edge{t, 0}
	fmt.Println(countQuartets(e, 0))
}

func countLeaves(n *tree.Node) {
	if len(n.Children) == 0 {
		leaves[n] = 1
		return
	}

	k := int64(0)
	for _, e := range n.Children {
		countLeaves(e.Node)
		k += leaves[e.Node]
	}
	leaves[n] = k
}

func countPairs(t *tree.Node) {
	if leaves[t] == 1 {
		pairs[t] = 0
		return
	}

	p := int64(0)
	for _, e := range t.Children {
		countPairs(e.Node)
		p += pairs[e.Node]
	}

	for i := 0; i < len(t.Children) - 1; i++ {
		e := t.Children[i]
		for j := i + 1; j < len(t.Children); j++ {
			f := t.Children[j]
			p += leaves[e.Node] * leaves[f.Node]
		}
	}
	pairs[t] = p
}

func countQuartets(e tree.Edge, pairsAbove int64) int64 {
	if pairs[e.Node] == 0 {
		return 0
	}

	fmt.Println("countQuartets of:")
	tree.Print(e.Node)

	a := leavesAbove * (leavesAbove - 1) / 2 % M
	b := leavesInSubtree * (leavesInSubtree - 1) / 2 % M
	q := a * b % M
	for _, child := range e.Children {
		outsideChild := pairs[e.Node] - pairs[child.Node] + pairsAbove
		q += countQuartets(child, outsideChild)
	}
	fmt.Printf("leavesAbove: %v leavesInSubtree: %v q: %v\n",
	           leavesAbove, leavesInSubtree, q)
	return q
}
