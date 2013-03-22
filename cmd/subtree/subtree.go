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
	taxaInv := map[string]int{}
	for i, taxon := range taxa {
		taxaInv[taxon] = i
	}

	line, _ = br.ReadString('\n')
	t, _ := tree.ReadNewick(strings.NewReader(line))

	m := t.SubtreeLeaves(taxaInv)
	s := subtree(t, size, m)
	if s == nil {
		fmt.Println("No subtree of size", size, "found")
	} else {
		a := m[s]
		for _, i := range a {
			fmt.Printf("%v ", taxa[i])
		}
		fmt.Println()
		s.WriteNewick(os.Stdout)
		fmt.Println()
	}
}

func subtree(t *tree.Node, size int, leaves map[*tree.Node][]int) *tree.Node {
	a, ok := leaves[t]
	if ok && len(a) == size {
		return t
	}

	for _, child := range t.Children {
		s := subtree(child.Node, size, leaves)
		if s != nil {
			return s
		}
	}
	return nil
}
