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

type crossPair struct {
	pair Pair
	// The subtrees containing the elements of the pair
	from [2]*tree.Node
}

type treeData struct {
	// Leaves in each subtree
	leaves map[*tree.Node][]int
	// Pairs in each subtree
	allPairs map[*tree.Node][]Pair
	// Pairs by their lowest common ancestor
	lcaPairs map[*tree.Node][]crossPair

	quartets map[Quartet]bool
}

func newTreeData() *treeData {
	td := &treeData{}
	td.leaves = map[*tree.Node][]int{}
	td.allPairs = map[*tree.Node][]Pair{}
	td.lcaPairs = map[*tree.Node][]crossPair{}
	td.quartets = map[Quartet]bool{}
	return td
}

// Assumes that taxa are only at the leaves.
func Quartets(t *tree.Node, taxa map[string]int) map[Quartet]bool {
	td := newTreeData()
	collectLeaves(t, taxa, td)
	collectPairs(t, td)

	collectQuartets(tree.Edge{t, 0}, nil, nil, td)
	return td.quartets
}

func collectLeaves(t *tree.Node, taxa map[string]int, td *treeData) {
	a := []int{}
	i, ok := taxa[t.Label]
	if ok {
		a = append(a, i)
	}

	for _, child := range t.Children {
		collectLeaves(child.Node, taxa, td)
		a = append(a, td.leaves[child.Node]...)
	}
	td.leaves[t] = a
}

func collectPairs(t *tree.Node, td *treeData) {
	all := []Pair{}
	for _, child := range t.Children {
		collectPairs(child.Node, td)
		all = append(all, td.allPairs[child.Node]...)
	}

	cross := []crossPair{}
	for i := 0; i < len(t.Children) - 1; i++ {
		for j := i + 1; j < len(t.Children); j++ {
			a := td.leaves[t.Children[i].Node]
			b := td.leaves[t.Children[j].Node]
			for _, x := range a {
				for _, y := range b {
					var cp crossPair
					cp.pair = newPair(x, y)
					cp.from[0] = t.Children[i].Node
					cp.from[1] = t.Children[j].Node
					cross = append(cross, cp)
					all = append(all, cp.pair)
				}
			}
		}
	}
	td.lcaPairs[t] = cross
	td.allPairs[t] = all
}

func collectQuartets(t tree.Edge,
					 leavesAbove []int, pairsAbove []Pair,
					 td *treeData) {
	for _, pa := range pairsAbove {
		for _, pb := range td.lcaPairs[t.Node] {
			td.quartets[newQuartet(pa, pb.pair)] = true
		}
	}

	newPairs := map[*tree.Node][]Pair{}
	for _, child := range t.Children {
		a := []Pair{}
		for _, x := range leavesAbove {
			for _, y := range td.leaves[child.Node] {
				a = append(a, newPair(x, y))
			}
		}
		newPairs[child.Node] = a
	}
			
	for _, child := range t.Children {
		if len(td.leaves[child.Node]) < 2 {
			continue
		}

		leavesAboveChild := []int{}
		leavesAboveChild = append(leavesAboveChild, leavesAbove...)
		pairsAboveChild := []Pair{}
		pairsAboveChild = append(pairsAboveChild, pairsAbove...)

		for _, other := range t.Children {
			if other == child {
				continue
			}
			leavesAboveChild = append(leavesAboveChild, td.leaves[other.Node]...)
			pairsAboveChild = append(pairsAboveChild, newPairs[other.Node]...)
			pairsAboveChild = append(pairsAboveChild, td.allPairs[other.Node]...)
		}
		for _, cp := range td.lcaPairs[t.Node] {
			if cp.from[0] != child.Node &&
			   cp.from[1] != child.Node {
				pairsAboveChild = append(pairsAboveChild, cp.pair)
		    }
		}

		collectQuartets(child, leavesAboveChild, pairsAboveChild, td)
	}
}

var _ = fmt.Println
