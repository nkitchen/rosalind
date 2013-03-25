package main

import "bufio"
import "flag"
import "fmt"
import "os"
import "rosalind/tree"
import "strconv"
import "strings"

func main() {
	flag.Parse()
	size, _ := strconv.Atoi(flag.Arg(0))

	br := bufio.NewReader(os.Stdin)

	line, _ := br.ReadString('\n')
	taxa := strings.Fields(line)

	taxaToKeep := taxa[:size]
	taxaMap := map[string]bool{}
	for _, taxon := range taxaToKeep {
		taxaMap[taxon] = true
	}

	line, _ = br.ReadString('\n')
	t1, _ := tree.ReadNewick(strings.NewReader(line))
	s1 := tree.CollapseUnrootedBinary(t1, taxaMap)

	line, _ = br.ReadString('\n')
	t2, _ := tree.ReadNewick(strings.NewReader(line))
	s2 := tree.CollapseUnrootedBinary(t2, taxaMap)

	for _, taxon := range taxaToKeep {
		fmt.Printf("%v ", taxon)
	}
	fmt.Println()
	s1.WriteNewick(os.Stdout)
	fmt.Println()
	s2.WriteNewick(os.Stdout)
	fmt.Println()
}
