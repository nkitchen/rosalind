package phylo

import "bytes"
import "fmt"
import "rosalind/tree"

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

type edge struct {
	tail, head *tree.Node
}

func revEdge(e edge) edge {
	return edge{e.head, e.tail}
}

var nilEdge = edge{nil, nil}

func (e edge) String() string {
	return fmt.Sprintf("%v -> %v", e.tail, e.head)
}

// A set of quartets represented implicitly by subtrees.
// [[A nil] [C nil]] is all quartets A||C, i.e., ab||cd with a, b in A
// and c, d in C.
// [[A B] [C D]] is all quartets AB||CD, i.e., ab||cd with a in A, b in B,
// c in C, and d in D.
type quartetSet [2][2]edge

func (qs quartetSet) String() string {
	b := &bytes.Buffer{}
	fmt.Fprintf(b, "{%v", qs[0][0])
	if qs[0][1] != nilEdge {
		fmt.Fprintf(b, ", %v", qs[0][1])
	}
	fmt.Fprintf(b, "}||{%v", qs[1][0])
	if qs[1][1] != nilEdge {
		fmt.Fprintf(b, ", %v", qs[1][1])
	}
	fmt.Fprintf(b, "}")
	return b.String()
}

type quartetData struct {
	// Children in rooted subtrees
	children map[edge][2]edge
	claimedQuartets map[edge][]quartetSet
}

func newQuartetData() *quartetData {
	qd := quartetData{}
	qd.children = map[edge][2]edge{}
	qd.claimedQuartets = map[edge][]quartetSet{}
	return &qd
}

type intersector struct {
	children1, children2 map[edge][2]edge
	// Intersection sizes already computed
	leafMemo map[[2]edge]int
}

func newIntersector(children1, children2 map[edge][2]edge) *intersector {
	return &intersector{children1, children2, map[[2]edge]int{}}
}

// Algorithm from "Computing the quartet distance between evolutionary trees",
// Bryant, Tsang, Kearney, and Li, 2000, as described in Tsang's master's
// thesis.
func QuartetDistance(t1, t2 *tree.Node, numLeaves int) int {
	qd1 := newQuartetData()
	e := edge{nil, t1}
	qd1.collectSubtrees(e)
	qd1.claimQuartets(e)

	qd2 := newQuartetData()
	e = edge{nil, t2}
	qd2.collectSubtrees(e)
	qd2.claimQuartets(e)

	totalQuartets := binom4(numLeaves)

	i := newIntersector(qd1.children, qd2.children)
	numShared := 0
	for _, claimed1 := range qd1.claimedQuartets {
		for _, qs1 := range claimed1 {
			for _, claimed2 := range qd2.claimedQuartets {
				for _, qs2 := range claimed2 {
					numShared += i.numSharedQuartets(qs1, qs2)
				}
			}
		}
	}

	return 2 * totalQuartets - 2 * numShared
}

func (qd *quartetData) collectSubtrees(e edge) {
	r := e.head
	if len(r.Children) == 3 {
		in := [3]edge{}
		out := [3]edge{}
		for i, c := range r.Children {
			out[i] = edge{r, c.Node}
			in[i] = revEdge(out[i])
			qd.collectSubtrees(out[i])
		}
		for i, a := range in {
			b := out[(i+1)%3]
			c := out[(i+2)%3]
			qd.children[a] = [2]edge{b, c}
		}
	} else if len(r.Children) == 2 {
		in := [2]edge{}
		out := [2]edge{}
		for i, c := range r.Children {
			out[i] = edge{r, c.Node}
			in[i] = revEdge(out[i])
			qd.collectSubtrees(out[i])
		}
		qd.children[e] = out

		for i, a := range in {
			qd.children[a] = [2]edge{out[(i+1)%2], revEdge(e)}
		}
	}
}

func (qd *quartetData) claimQuartets(e edge) {
	if len(e.head.Children) == 0 {
		return
	}

	if e.tail != nil {
		sides := [2][][2]edge{}
		for i, ee := range ([]edge{revEdge(e), e}) {
			s := qd.children[ee]
			// cross of subtrees
			sides[i] = append(sides[i], s)

			// unclaimed subtrees
			for _, t := range s {
				_, claimed := qd.claimedQuartets[t]
				if !claimed {
					_, claimed = qd.claimedQuartets[revEdge(t)]
				}
				if !claimed {
					sides[i] = append(sides[i], [2]edge{t, nilEdge})
				}
			}
		}

		sets := []quartetSet{}
		for _, a := range sides[0] {
			for _, b := range sides[1] {
				 sets = append(sets, quartetSet{a, b})
			}
		}
		qd.claimedQuartets[e] = sets
	}

	r := e.head
	for _, t := range r.Children {
		qd.claimQuartets(edge{r, t.Node})
	}
}

func (i *intersector) numSharedQuartets(qs1, qs2 quartetSet) int {
	nq := 0
	np := i.numSharedPairs(qs1[0], qs2[0])
	if np != 0 {
		nq += np * i.numSharedPairs(qs1[1], qs2[1])
	}
	np = i.numSharedPairs(qs1[0], qs2[1])
	if np != 0 {
		nq += np * i.numSharedPairs(qs1[1], qs2[0])
	}
	return nq
}

func (i *intersector) numSharedPairs(a, b [2]edge) int {
	np := 0
	nl := i.numSharedLeaves(a[0], b[0])
	switch {
	case a[1] == nilEdge && b[1] == nilEdge:
		np += nl * (nl - 1) / 2
	case a[1] == nilEdge && b[1] != nilEdge:
		if nl != 0 {
			np += nl * i.numSharedLeaves(a[0], b[1])
		}
	case a[1] != nilEdge && b[1] == nilEdge:
		if nl != 0 {
			np += nl * i.numSharedLeaves(a[1], b[0])
		}
	case a[1] != nilEdge && b[1] != nilEdge:
	    if nl != 0 {
			np += nl * i.numSharedLeaves(a[1], b[1])
		}
		nl = i.numSharedLeaves(a[0], b[1])
		if nl != 0 {
			np += nl * i.numSharedLeaves(a[1], b[0])
		}
	}
	return np
}

func (i *intersector) numSharedLeaves(e1, e2 edge) int {
	if len(e1.head.Children) == 0 && len(e2.head.Children) == 0 {
		if e1.head.Label == e2.head.Label {
			return 1
		} else {
			return 0
		}
	}

	k := [2]edge{e1, e2}
	if n, ok := i.leafMemo[k]; ok {
		return n
	}

	var result int
	if len(e1.head.Children) > 0 {
		c := i.children1[e1]
		result =  i.numSharedLeaves(c[0], e2) +
			i.numSharedLeaves(c[1], e2)
	} else {
		c := i.children2[e2]
		result = i.numSharedLeaves(e1, c[0]) +
			i.numSharedLeaves(e1, c[1])
	}
	i.leafMemo[k] = result
	return result
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

//// The slices for all the nodes' leaves shared the same backing array,
//// so the total storage of leaves requires only O(N) space.
//func collectSubtreeLeaves(t *tree.Node, taxa map[string]int,
//                          collected []int, td *treeData) int {
//	if collected == nil {
//		collected = make([]int, len(taxa))
//	}
//
//	n := 0
//	i, ok := taxa[t.Label]
//	if ok {
//		collected[0] = i
//		n++
//	}
//
//	for _, child := range t.Children {
//		n += collectSubtreeLeaves(child.Node, taxa, collected[n:], td)
//	}
//	td.subtreeLeaves[t] = collected[:n]
//	return n
//}
//
//func collectEdgeLeavesBelow(td *treeData) {
//	for t := range td.subtreeLeaves {
//		for _, child := range t.Children {
//			a := td.subtreeLeaves[child.Node]
//			td.edgeLeaves[t] = append(td.edgeLeaves[t], a)
//		}
//	}
//}
//
//func collectEdgeLeavesAbove(t *tree.Node, leaves []int, td *treeData) {
//	if len(leaves) > 0 {
//		td.edgeLeaves[t] = append(td.edgeLeaves[t], leaves)
//	}
//
//	for _, child := range t.Children {
//		a := append([]int{}, leaves...)
//		for _, other := range t.Children {
//			if other == child {
//				continue
//			}
//			a = append(a, td.subtreeLeaves[other.Node]...)
//		}
//		collectEdgeLeavesAbove(child.Node, a, td)
//	}
//}
//
//func collectPairs(td *treeData) {
//	numTaxa := 0
//	for _, leaves := range td.subtreeLeaves {
//		if len(leaves) > numTaxa {
//			numTaxa = len(leaves)
//		}
//	}
//
//	for _, e := range td.edgeLeaves {
//		if len(e) != 3 {
//			continue
//		}
//
//		for i := range e {
//			j := (i + 1) % len(e)
//			s := NewCharArray(numTaxa)
//			for _, leaf := range e[i] {
//				s.Set(leaf, 1)
//			}
//			for _, leaf := range e[j] {
//				s.Set(leaf, 1)
//			}
//			if s.PopCount() == numTaxa - 1 {
//				continue
//			}
//			for _, x := range e[i] {
//				for _, y := range e[j] {
//					p := NewPair(x, y)
//					td.pairSplits[p] = append(td.pairSplits[p], s)
//				}
//			}
//		}
//	}
//}

func printClaimed(qd *quartetData) {
	for e, a := range qd.claimedQuartets {
		fmt.Println(e)
		for _, s := range a {
			fmt.Println("   ", s)
		}
	}
}
