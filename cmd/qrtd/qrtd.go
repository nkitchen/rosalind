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

	d := phylo.QuartetDistance(t1, t2, taxaInv)
	fmt.Println(d)
} 
