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
	// The leaves in each subtree (subtree with respect to the data structure)
	subtreeLeaves map[*tree.Node][]int
	// The leaves in the unrooted subtree on the other side of each incident edge
	edgeLeaves map[*tree.Node][][]int
	// Maps pairs to splits where they appear together
	pairSplits map[Pair][]CharArray
}

func newTreeData() *treeData {
	td := &treeData{}
	td.subtreeLeaves = map[*tree.Node][]int{}
	td.edgeLeaves = map[*tree.Node][][]int{}
	td.pairSplits = map[Pair][]CharArray{}
	return td
}

// Assumes that taxa are only at the leaves.
func QuartetDistance(t1, t2 *tree.Node, taxa map[string]int) int {
	td1 := newTreeData()
	collectSubtreeLeaves(t1, taxa, nil, td1)
	collectEdgeLeavesBelow(td1)
	collectEdgeLeavesAbove(t1, nil, td1)
	collectPairs(td1)

	td2 := newTreeData()
	collectSubtreeLeaves(t2, taxa, nil, td2)
	collectEdgeLeavesBelow(td2)
	collectEdgeLeavesAbove(t2, nil, td2)
	collectPairs(td2)

	q1 := binom4(len(td1.subtreeLeaves[t1]))
	q2 := binom4(len(td2.subtreeLeaves[t2]))
	shared := 0
	for p, a1 := range td1.pairSplits {
		a2, ok := td2.pairSplits[p]
		if !ok {
			continue
		}

		fmt.Println("shared pair", p)
		fmt.Println("a1", a1)
		fmt.Println("a2", a2)

		for _, s1 := range a1 {
			for _, s2 := range a2 {
				if len(s1) != len(s2) {
					panic("Length mismatch")
				}
		
				sharedLeavesAbove := 0
				for i := range s1 {
					if s1[i] == 0 && s2[i] == 0 {
						sharedLeavesAbove++
					}
				}
				fmt.Println("sharedLeavesAbove", sharedLeavesAbove)
				shared += sharedLeavesAbove * (sharedLeavesAbove - 1) / 2
			}
		}
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
func collectSubtreeLeaves(t *tree.Node, taxa map[string]int,
                          collected []int, td *treeData) int {
	if collected == nil {
		collected = make([]int, len(taxa))
	}

	n := 0
	i, ok := taxa[t.Label]
	if ok {
		collected[0] = i
		n++
	}

	for _, child := range t.Children {
		n += collectSubtreeLeaves(child.Node, taxa, collected[n:], td)
	}
	td.subtreeLeaves[t] = collected[:n]
	return n
}

func collectEdgeLeavesBelow(td *treeData) {
	for t := range td.subtreeLeaves {
		for _, child := range t.Children {
			a := td.subtreeLeaves[child.Node]
			td.edgeLeaves[t] = append(td.edgeLeaves[t], a)
		}
	}
}

func collectEdgeLeavesAbove(t *tree.Node, leaves []int, td *treeData) {
	if len(leaves) > 0 {
		td.edgeLeaves[t] = append(td.edgeLeaves[t], leaves)
	}

	for _, child := range t.Children {
		a := append([]int{}, leaves...)
		for _, other := range t.Children {
			if other == child {
				continue
			}
			a = append(a, td.subtreeLeaves[other.Node]...)
		}
		collectEdgeLeavesAbove(child.Node, a, td)
	}
}

func collectPairs(td *treeData) {
	fmt.Println("edgeLeaves", td.edgeLeaves)
	numTaxa := 0
	for _, leaves := range td.subtreeLeaves {
		if len(leaves) > numTaxa {
			numTaxa = len(leaves)
		}
	}

	for _, e := range td.edgeLeaves {
		if len(e) != 3 {
			continue
		}

		for i := range e {
			j := (i + 1) % len(e)
			s := make(CharArray, numTaxa)
			for _, leaf := range e[i] {
				s[leaf] = 1
			}
			for _, leaf := range e[j] {
				s[leaf] = 1
			}
			if s.PopCount() == numTaxa - 1 {
				continue
			}
			for _, x := range e[i] {
				for _, y := range e[j] {
					p := NewPair(x, y)
					td.pairSplits[p] = append(td.pairSplits[p], s)
				}
			}
		}
	}
	fmt.Println("pairSplits", td.pairSplits)
}
