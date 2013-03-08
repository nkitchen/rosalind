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

// Assumes that taxa are only at the leaves.
func Quartets(t *tree.Node, taxa map[string]int) map[Quartet]bool {
	leaves := map[*tree.Node][]int{}
	collectLeaves(t, taxa, leaves)

	// Pairs by their lowest common ancestor
	lcaPairs := map[*tree.Node][]Pair{}
	collectPairs(t, leaves, lcaPairs)

	q := map[Quartet]bool{}
	collectQuartets(tree.Edge{t, 0}, lcaPairs, nil, q)
	return q
}

func collectLeaves(t *tree.Node, taxa map[string]int,
                   leaves map[*tree.Node][]int) {
	i, ok := taxa[t.Label]
	if ok {
		leaves[t] = append(leaves[t], i)
	}

	for _, child := range t.Children {
		collectLeaves(child.Node, taxa, leaves)
		leaves[t] = append(leaves[t], leaves[child.Node]...)
	}
}

func collectPairs(t *tree.Node, leaves map[*tree.Node][]int,
                  pairs map[*tree.Node][]Pair) {
	for _, child := range t.Children {
		collectPairs(child.Node, leaves, pairs)
	}

	for i := 0; i < len(t.Children) - 1; i++ {
		for j := i + 1; j < len(t.Children); j++ {
			a := leaves[t.Children[i].Node]
			b := leaves[t.Children[j].Node]
			for _, x := range a {
				for _, y := range b {
					pairs[t] = append(pairs[t], newPair(x, y))
				}
			}
		}
	}
}

func collectQuartets(t tree.Edge, lcaPairs map[*tree.Node][]Pair,
					 priorPairs []Pair, quartets map[Quartet]bool) []Pair {
	fmt.Println("collectQuartets:")
	tree.Print(t.Node)
	fmt.Printf("lcaPairs: %v\n", lcaPairs)
	fmt.Printf("priorPairs: %v\n", priorPairs)
	fmt.Printf("quartets: %v\n", quartets)
	for _, pa := range priorPairs {
		for _, pb := range lcaPairs[t.Node] {
			quartets[newQuartet(pa, pb)] = true
		}
	}

	priorPairs = append(priorPairs, lcaPairs[t.Node]...)

	for _, child := range t.Children {
		priorPairs = collectQuartets(child, lcaPairs, priorPairs, quartets)
	}

	return priorPairs
}
