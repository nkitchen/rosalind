package phylo

import "rosalind/tree"
import "fmt"

type Pair [2]int
type Quartet [2]Pair

func newPair(a, b int) Pair {
	if a <= b {
		return Pair{a, b}
	}
	return Pair{b, a}
}

func newQuartet(p, q Pair) Quartet {
	if p[0] <= q[0] {
		return Quartet{p, q}
	}
	return Quartet{q, p}
}

type treeData struct {
	// Leaves in each subtree
	leaves map[*tree.Node][]int
	splits map[tree.Edge]CharArray
	// Maps pairs to the splits where they first appear
	// (constructed from inverse splits for pairs at the root)
	pairs map[Pair]CharArray
}

func newTreeData() *treeData {
	td := &treeData{}
	td.leaves = map[*tree.Node][]int{}
	td.splits = map[tree.Edge]CharArray{}
	td.pairs = map[Pair]CharArray{}
	return td
}

// Assumes that taxa are only at the leaves.
func QuartetDistance(t1, t2 *tree.Node, taxa map[string]int) int {
	td1 := newTreeData()
	leaves := make([]int, len(taxa))
	collectLeaves(t1, taxa, leaves, td1)
	fmt.Println("leaves", td1.leaves)
	collectSplits(t1.Edge(), taxa, td1.splits)
	fmt.Println("splits", td1.splits)
	collectPairs(t1, td1)
	fmt.Println("pairs", td1.pairs)

	td2 := newTreeData()
	leaves = make([]int, len(taxa))
	collectLeaves(t2, taxa, leaves, td2)
	fmt.Println("leaves", td2.leaves)
	collectSplits(t2.Edge(), taxa, td2.splits)
	fmt.Println("splits", td2.splits)
	collectPairs(t2, td2)
	fmt.Println("pairs", td2.pairs)

	q1 := 0
	q2 := 0
	shared := 0
	for p, a1 := range td1.pairs {
		fmt.Println("Pair", p)
		a2, ok := td2.pairs[p]
		fmt.Println("a1", a1, "a2", a2)
		if !ok {
			continue
		}
		if len(a1) != len(a2) {
			panic("Length mismatch")
		}
		// Increment quartet counts
		m1 := len(a1) - a1.PopCount()
		q1 += m1 * (m1 - 1)
		m2 := len(a2) - a2.PopCount()
		q2 += m2 * (m2 - 1)
		// Find shared quartets
		sharedAbove := 0
		for i := range a1 {
			if a1[i] == 0 && a2[i] == 0 {
				sharedAbove++
			}
		}
		shared += sharedAbove * (sharedAbove - 1)
	}
	return q1 + q2 - 2 * shared
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
	fmt.Println("subtree pairs", td.pairs)

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
			fmt.Println("Third child:", t.Children[k].Node)

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
					p := newPair(x, y)
					td.pairs[p] = s
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
					p := newPair(x, y)
					td.pairs[p] = s
				}
			}
		}
	}
}

//func collectQuartets(t tree.Edge,
//					 leavesAbove []int, pairsAbove []Pair,
//					 td *treeData) {
//	for _, pa := range pairsAbove {
//		for _, pb := range td.lcaPairs[t.Node] {
//			td.quartets[newQuartet(pa, pb.pair)] = true
//		}
//	}
//
//	newPairs := map[*tree.Node][]Pair{}
//	for _, child := range t.Children {
//		a := []Pair{}
//		for _, x := range leavesAbove {
//			for _, y := range td.leaves[child.Node] {
//				a = append(a, newPair(x, y))
//			}
//		}
//		newPairs[child.Node] = a
//	}
//			
//	for _, child := range t.Children {
//		if len(td.leaves[child.Node]) < 2 {
//			continue
//		}
//
//		leavesAboveChild := []int{}
//		leavesAboveChild = append(leavesAboveChild, leavesAbove...)
//		pairsAboveChild := []Pair{}
//		pairsAboveChild = append(pairsAboveChild, pairsAbove...)
//
//		for _, other := range t.Children {
//			if other == child {
//				continue
//			}
//			leavesAboveChild = append(leavesAboveChild, td.leaves[other.Node]...)
//			pairsAboveChild = append(pairsAboveChild, newPairs[other.Node]...)
//			pairsAboveChild = append(pairsAboveChild, td.allPairs[other.Node]...)
//		}
//		for _, cp := range td.lcaPairs[t.Node] {
//			if cp.from[0] != child.Node &&
//			   cp.from[1] != child.Node {
//				pairsAboveChild = append(pairsAboveChild, cp.pair)
//		    }
//		}
//
//		collectQuartets(child, leavesAboveChild, pairsAboveChild, td)
//	}
//}

var _ = fmt.Println
