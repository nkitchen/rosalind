package main

import "bufio"
import "flag"
import "fmt"
import "os"
import "rosalind/gene"

func main() {
	k := flag.Int("k", 3, "Node substring length")
	flag.Parse()

	adj := map[string][]string{}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		read := scanner.Text()
		if len(read) < *k + 1 {
			fmt.Fprintf(os.Stderr, "Error: Reads are too short for k.")
			os.Exit(1)
		}

		for _, s := range ([]string{read, gene.ReverseComplement(read)}) {
			for i := 0; i < len(s) - *k; i++ {
				u := s[i:i + *k]
				v := s[i+1:i + 1 + *k]
				adj[u] = add(adj[u], v)
			}
		}
	}

	nonce := 1
	fmt.Println("digraph debruijn {")
	for u := range adj {
		for _, v := range adj[u] {
			if u == v {
				fmt.Printf("  %s -> nonce%d;\n", u, nonce)
				fmt.Printf("  nonce%d -> %s;\n", nonce, u)
				nonce++
			} else {
				fmt.Printf("  %s -> %s;\n", u, v)
			}
		}
	}
	fmt.Println("}")
}

// Appends a string only if not already present.
func add(a []string, s string) []string {
	for _, r := range a {
		if r == s {
			return a
		}
	}

	return append(a, s)
}
