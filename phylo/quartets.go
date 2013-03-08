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
	leaves := map[*tree.Node]map[int]bool{}
	collectLeaves(t, taxa, leaves)

	// Pairs by their lowest common ancestor
	lcaPairs := map[*tree.Node][]Pair{}
	// All pairs in subtree
	allPairs := map[*tree.Node][]Pair{}
	collectPairs(t, leaves, lcaPairs, allPairs)

	q := map[Quartet]bool{}
	collectQuartets(tree.Edge{t, 0}, leaves, lcaPairs, allPairs, nil, q)
	return q
}

func collectLeaves(t *tree.Node, taxa map[string]int,
                   leaves map[*tree.Node]map[int]bool) {
	leaves[t] = map[int]bool{}
	i, ok := taxa[t.Label]
	if ok {
		leaves[t][i] = true
	}

	for _, child := range t.Children {
		collectLeaves(child.Node, taxa, leaves)
		for k := range leaves[child.Node] {
			leaves[t][k] = true
		}
	}
}

func collectPairs(t *tree.Node, leaves map[*tree.Node]map[int]bool,
                  lcaPairs, allPairs map[*tree.Node][]Pair) {
	for _, child := range t.Children {
		collectPairs(child.Node, leaves, lcaPairs, allPairs)
		allPairs[t] = append(allPairs[t], allPairs[child.Node]...)
	}

	for i := 0; i < len(t.Children) - 1; i++ {
		for j := i + 1; j < len(t.Children); j++ {
			a := leaves[t.Children[i].Node]
			b := leaves[t.Children[j].Node]
			for x := range a {
				for y := range b {
					lcaPairs[t] = append(lcaPairs[t], newPair(x, y))
				}
			}
		}
	}
	fmt.Printf("lcaPairs(%v): %v\n", t, lcaPairs[t])
	allPairs[t] = append(allPairs[t], lcaPairs[t]...)
}

func collectQuartets(t tree.Edge,
					 leaves map[*tree.Node]map[int]bool,
                     lcaPairs, allPairs map[*tree.Node][]Pair,
					 pairsAbove []Pair,
                     quartets map[Quartet]bool) {
	fmt.Printf("collectQuartets(%v, lcaPairs=%v, allPairs=%v, pairsAbove=%v)\n", t.Node, lcaPairs[t.Node], allPairs[t.Node], pairsAbove)
	for _, child := range t.Children {
		pairsAboveChild := []Pair{}
		pairsAboveChild = append(pairsAboveChild, pairsAbove...)
		for _, p := range allPairs[t.Node] {
			inChild := leaves[child.Node][p[0]] ||
			   leaves[child.Node][p[1]]
			if !inChild {
				pairsAboveChild = append(pairsAboveChild, p)
			}
		}

		fmt.Printf("pairsAboveChild(%v): %v\n", child.Node, pairsAboveChild)
		fmt.Printf("x lcaPairs: %v\n", lcaPairs[child.Node])
		for _, pa := range pairsAboveChild {
			for _, pb := range lcaPairs[child.Node] {
				quartets[newQuartet(pa, pb)] = true
			}
		}
		fmt.Printf("quartets: %v\n", quartets)
		
		collectQuartets(child, leaves, lcaPairs, allPairs, pairsAboveChild,
		                quartets)
		fmt.Printf("quartets after child %v: %v\n", child.Node, quartets)
	}
	fmt.Printf("collectQuartets(%v): %v\n", t.Node, quartets)
}

var _ = fmt.Println
