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
	fmt.Println(splits1)
	quartets1 := phylo.Quartets(t1, taxaInv)
	quartets2 := phylo.Quartets(t2, taxaInv)
	fmt.Println(quartets1)
	fmt.Println(quartets2)
	shared := 0
	for q := range quartets2 {
		if quartets2[q] {
			shared++
		}
	}

	d := len(quartets1) + len(quartets2) - 2 * shared
	fmt.Println(d)
} 
