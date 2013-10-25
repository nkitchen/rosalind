package main

import "bufio"
import "fmt"
import "os"
import "rosalind/gene"

func main() {
	adj := map[string]map[string]bool{}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		fwd := scanner.Text()
		rev := gene.ReverseComplement(fwd)
		for _, s := range ([]string{fwd, rev}) {
			n := len(s)
			
			u := s[:n-1]
			v := s[1:]

			if adj[u] == nil {
				adj[u] = map[string]bool{}
			}
			adj[u][v] = true
		}
	}

	for u := range adj {
		for v := range adj[u] {
			fmt.Printf("(%s, %s)\n", u, v)
		}
	}
}
