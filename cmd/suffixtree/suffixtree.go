package main

import "bufio"
import "fmt"
import "os"
import "rosalind/strings/suffix"

var s string

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	s = scanner.Text() + "$"

	t := suffix.NewTree(s)

	fmt.Println("digraph suffix {")
	for tail, a := range t.Edges {
		for _, e := range a {
			u := s[e.Loc:e.Loc + e.Len]
			fmt.Printf("  node%d -> node%d [label=\"%s * %d\"];\n", tail, e.Head, u, e.LeafCount)
		}
	}
	fmt.Println("}")
}

