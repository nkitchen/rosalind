package main

import "bufio"
import "fmt"
import "math/big"
import "os"
import "rosalind/tree"
import "sort"
import "strings"

type CharacterTable struct {
	Data []*big.Int
	NumTaxa int
}

func main() {
	br := bufio.NewReader(os.Stdin)

	line, _ := br.ReadString(';')
	t, err := tree.ReadNewick(strings.NewReader(line))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	ctab := charTable(t)
	for _, array := range ctab.Data {
		for i := 0; i < ctab.NumTaxa; i++ {
			fmt.Print(array.Bit(i))
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
// This version is optimized with the assumption that taxa are only
// at the leaves and that they are distinct.
func charTable(t *tree.Node) CharacterTable {
	taxa := findTaxa(tree.Edge{t,0}, nil)
	sort.StringSlice(taxa).Sort()

	taxaInv := map[string]int{}
	for i, taxon := range taxa {
		taxaInv[taxon] = i
	}

	ctab, _ := subtreeCharTable(tree.Edge{t, 0}, taxaInv, CharacterTable{})
	// The last array just includes all the taxa.
	ctab.Data = ctab.Data[:len(ctab.Data)-1]
	ctab.NumTaxa = len(taxaInv)
	return ctab
}

func subtreeCharTable(t tree.Edge, taxaInv map[string]int, table CharacterTable) (CharacterTable, *big.Int) {
	s := big.NewInt(0)

	if t.Label != "" {
		s.SetInt64(1)
		s.Lsh(s, uint(taxaInv[t.Label]))
	}

	for _, u := range t.Children {
		var r *big.Int
		table, r = subtreeCharTable(u, taxaInv, table)
		s.Or(s, r)
	}
	if len(t.Children) > 0 {
		table.Data = append(table.Data, s)
	}

	return table, s
}
