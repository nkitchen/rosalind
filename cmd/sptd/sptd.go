package main

import "bufio"
import "fmt"
import "os"
import "rosalind/phylo"
import "rosalind/tree"
import "strings"

func main() {
	br := bufio.NewReader(os.Stdin)
	
	line, _ := br.ReadString('\n')
	taxa := strings.Fields(line)
	taxaInv := map[string]int{}
	for i, taxon := range taxa {
		taxaInv[taxon] = i
	}

	line, _ = br.ReadString('\n')
	t1, _ := tree.ReadNewick(strings.NewReader(line))
	line, _ = br.ReadString('\n')
	t2, _ := tree.ReadNewick(strings.NewReader(line))

	splits1 := phylo.Splits(t1, taxaInv)
	splits2 := phylo.Splits(t2, taxaInv)
	normalize(splits1)
	normalize(splits2)

	m1 := map[string]bool{}
	for _, a := range splits1 {
		m1[string(a)] = true
	}

	shared := 0
	for _, a := range splits2 {
		if m1[string(a)] {
			shared++
		}
	}

	d := 2 * (len(taxa) - 3 - shared)
	fmt.Println(d)
}

func normalize(s []phylo.CharArray) {
	for _, a := range s {
		if a[0] != 1 {
			a.Not(a)
		}
	}
}
