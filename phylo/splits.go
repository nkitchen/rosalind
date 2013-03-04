package phylo

import "rosalind/tree"

// Splits returns the nontrivial splits contained in a tree.
func Splits(t *tree.Node, taxa map[string]int) []CharArray {
	_, s := collectSplits(tree.Edge{t, 0}, taxa, nil)
	return s
}

// collectSplits returns a character array for all the taxa in a subtree
// and appends the nontrivial splits for the subtree to collected.
func collectSplits(e tree.Edge, taxa map[string]int,
                   collected []CharArray) (CharArray, []CharArray) {
	a := make(CharArray, len(taxa))
	i, ok := taxa[e.Label]
	if ok {
		a[i] = 1
	}
    for _, child := range e.Children {
		b, c := collectSplits(child, taxa, collected)
		a.Or(a, b)
		collected = c
	}

	n := a.PopCount()
	if 1 < n && n < len(taxa) - 1 {
		return a, append(collected, a)
	}
	return a, collected
}
