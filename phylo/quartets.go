package phylo

import "rosalind/tree"
import "fmt"

type Pair [2]int
type Quartet [2]Pair

// NewPair returns a canonical pair {a, b}.
func NewPair(a, b int) Pair {
	if a <= b {
		return Pair{a, b}
	}
	return Pair{b, a}
}

// NewQuartet returns a canonical quartet {a, b} | {c, d}.
func NewQuartet(p, q Pair) Quartet {
	if p[0] <= q[0] {
		return Quartet{p, q}
	}
	return Quartet{q, p}
}

func (p1 Pair) Less(p2 Pair) bool {
	switch {
	case p1[0] < p2[0]:
		return true
	case p1[0] > p2[0]:
		return false
	case p1[1] < p2[1]:
		return true
	}
	return false
}

func (q1 Quartet) Less(q2 Quartet) bool {
	switch {
	case q1[0].Less(q2[0]):
		return true
	case q2[0].Less(q1[0]):
		return false
	case q1[1].Less(q2[1]):
	    return true
	}
	return false
}

type QuartetSlice []Quartet

func (qs QuartetSlice) Len() int {
	return len(qs)
}

func (qs QuartetSlice) Less(i, j int) bool {
	return qs[i].Less(qs[j])
}

func (qs QuartetSlice) Swap(i, j int) {
	qs[i], qs[j] = qs[j], qs[i]
}

func (q Quartet) String() string {
	return fmt.Sprintf("%v:%v|%v:%v", q[0][0], q[0][1], q[1][0], q[1][1])
}

type treeData struct {
	// Leaves in each subtree
	leaves map[*tree.Node][]int
	splits map[tree.Edge]CharArray
	// Maps pairs to their lowest common ancestors
	pairLCAs map[Pair]*tree.Node
	// Maps pairs to the splits where they first appear
	// (constructed from inverse splits for pairs at the root)
	pairSplits map[Pair]CharArray
}

func newTreeData() *treeData {
	td := &treeData{}
	td.leaves = map[*tree.Node][]int{}
	td.splits = map[tree.Edge]CharArray{}
	td.pairLCAs = map[Pair]*tree.Node{}
	td.pairSplits = map[Pair]CharArray{}
	return td
}

// Assumes that taxa are only at the leaves.
func QuartetDistance(t1, t2 *tree.Node, taxa map[string]int) int {
	td1 := newTreeData()
	leaves := make([]int, len(taxa))
	collectLeaves(t1, taxa, leaves, td1)
	collectSplits(t1.Edge(), taxa, td1.splits)
	collectPairs(t1, td1)

	td2 := newTreeData()
	leaves = make([]int, len(taxa))
	collectLeaves(t2, taxa, leaves, td2)
	collectSplits(t2.Edge(), taxa, td2.splits)
	collectPairs(t2, td2)

	q1 := binom4(len(td1.leaves[t1]))
	q2 := binom4(len(td2.leaves[t2]))
	shared := 0
	for p, a1 := range td1.pairSplits {
		a2, ok := td2.pairSplits[p]
		if !ok {
			continue
		}
		if len(a1) != len(a2) {
			panic("Length mismatch")
		}
		
		fmt.Println("shared pair", p)
		fmt.Println("a1", a1)
		fmt.Println("a2", a2)
		sharedLeavesAbove := 0
		for i := range a1 {
			if a1[i] == 0 && a2[i] == 0 {
				sharedLeavesAbove++
			}
		}
		fmt.Println("sharedLeavesAbove", sharedLeavesAbove)
		shared += sharedLeavesAbove * (sharedLeavesAbove - 1) / 2
	}
	fmt.Println("QuartetDistance: q1", q1, "q2", q2, "shared", shared)
	// We actually find each shared quartet twice, once for each pair.
	return q1 + q2 - shared
}

// Returns n choose 4.
func binom4(n int) int {
	if n < 4 {
		return 0
	}
	b := 1
	for k := 0; k < 4; k++ {
		b *= n - k
	}
	return b / (4 * 3 * 2)
}

// The slices for all the nodes' leaves shared the same backing array,
// so the total storage of leaves requires only O(N) space.
func collectLeaves(t *tree.Node, taxa map[string]int,
                   collected []int, td *treeData) int {
	n := 0
	i, ok := taxa[t.Label]
	if ok {
		collected[0] = i
		n++
	}

	for _, child := range t.Children {
		n += collectLeaves(child.Node, taxa, collected[n:], td)
	}
	td.leaves[t] = collected[:n]
	return n
}

func collectPairs(t *tree.Node, td *treeData) {
	for _, child := range t.Children {
		collectSubtreePairs(child, td)
	}

	n := len(t.Children)
	if n != 3 {
		panic(fmt.Sprint("Unexpected number of child nodes:", n))
	}

	for i := 0; i < len(t.Children) - 1; i++ {
		for j := i + 1; j < len(t.Children); j++ {
			var k int
			for k = range t.Children {
				if k != i && k != j {
					break
				}
			}
			if k == len(t.Children) {
				panic("Third child not found")
			}

			var s CharArray
			ss, ok := td.splits[t.Children[k]]
			if !ok {
				continue
			}
			s.Not(ss)

			a := td.leaves[t.Children[i].Node]
			b := td.leaves[t.Children[j].Node]
			for _, x := range a {
				for _, y := range b {
					p := NewPair(x, y)
					td.pairLCAs[p] = t
					td.pairSplits[p] = s
				}
			}
		}
	}
}

func collectSubtreePairs(t tree.Edge, td *treeData) {
	for _, child := range t.Children {
		collectSubtreePairs(child, td)
	}

	for i := 0; i < len(t.Children) - 1; i++ {
		for j := i + 1; j < len(t.Children); j++ {
			a := td.leaves[t.Children[i].Node]
			b := td.leaves[t.Children[j].Node]
			s, ok := td.splits[t]
			for _, x := range a {
				for _, y := range b {
					if !ok {
						panic(fmt.Sprintf("No split found for node %v", t.Node))
					}
					p := NewPair(x, y)
					td.pairLCAs[p] = t.Node
					td.pairSplits[p] = s
				}
			}
		}
	}
}

var _ = fmt.Println
