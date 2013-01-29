package main

import "bufio"
import "fmt"
import "os"
import "rosalind/tree"
import "sort"
import "strings"

func main() {
	br := bufio.NewReader(os.Stdin)

	line, _ := br.ReadString(';')
	t, err := tree.ReadNewick(strings.NewReader(line))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	taxa := findTaxa(tree.Edge{t,0}, nil)
	sort.StringSlice(taxa).Sort()

	taxaInv := map[string]int{}
	for i, taxon := range taxa {
		taxaInv[taxon] = i
	}

    ctab, _ := charTable(tree.Edge{t,0}, taxaInv, nil)
	for _, array := range ctab {
		for _, c := range array {
			fmt.Print(c)
		}
		fmt.Println()
	}
}

func findTaxa(t tree.Edge, taxa []string) []string {
	if t.Label != "" {
		taxa = append(taxa, t.Label)
	}
	for _, u := range t.Children {
		taxa = findTaxa(u, taxa)
	}
	return taxa
}

// charTable returns the character table and the set of taxa for the subtree.
func charTable(t tree.Edge, taxaInv map[string]int, table [][]int) ([][]int, map[string]bool) {
	s := make(map[string]bool)

	if t.Label != "" {
		s[t.Label] = true
	}

	added := false
	for _, u := range t.Children {
		var r map[string] bool
		table, r = charTable(u, taxaInv, table)
		for label := range r {
			if !s[label] {
				added = true
			}
			s[label] = true
		}
	}

	if 1 < len(s) && len(s) < len(taxaInv) - 1 && added {
		a := make([]int, len(taxaInv))
		for label := range s {
			a[taxaInv[label]] = 1
		}
		table = append(table, a)
	}

	return table, s
}
