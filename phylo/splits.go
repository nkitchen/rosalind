package phylo

import "rosalind/tree"

// Splits returns the nontrivial splits contained in a tree.
func Splits(t *tree.Node, taxa map[string]int) []CharArray {
	m := map[tree.Edge]CharArray{}
	collectSplits(tree.Edge{t, 0}, taxa, m)
	s := make([]CharArray, 0, len(m))
	for _, a := range m {
		s = append(s, a)
	}
	return s
}

// collectSplits returns a character array for all the taxa in a subtree
// and appends the nontrivial splits for the subtree to collected.
func collectSplits(e tree.Edge, taxa map[string]int,
                   collected map[tree.Edge]CharArray) CharArray {
	a := make(CharArray, len(taxa))
	i, ok := taxa[e.Label]
	if ok {
		a[i] = 1
	}
    for _, child := range e.Children {
		b := collectSplits(child, taxa, collected)
		a.Or(a, b)
	}

	n := a.PopCount()
	if 1 < n && n < len(taxa) - 1 {
		collected[e] = a
	}
	return a
}
