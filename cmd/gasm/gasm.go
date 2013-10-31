package main

import "bufio"
import "errors"
import "fmt"
import "os"
import "rosalind/gene"

func main() {
	reads := []string{}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		reads = append(reads, scanner.Text())
	}

	for k := len(reads[0]) - 1; k > 0; k-- {
		fmt.Println(k)
		g, err := cyclicGenome(reads, k)
		if err == nil {
			fmt.Println(g)
			return
		}
	}

	fmt.Println("No genome found")
}

func cyclicGenome(reads []string, k int) (string, error) {
	adj := map[string][]string{}

	for _, r := range reads {
		for _, s := range ([]string{r, gene.ReverseComplement(r)}) {
			for i := 0; i < len(s) - k; i++ {
				u := s[i:i + k]
				v := s[i+1:i + 1 + k]
				adj[u] = add(adj[u], v)
			}
		}
	}

    cycle := func (start string) (string, error) {
		c := []byte{}
		visited := map[string]bool{}
		v := start
		for {
			if len(adj[v]) != 1 {
				e := fmt.Sprintf("Node %s has %d successors.", v, len(adj[v]))
				return "", errors.New(e)
			}
			v = adj[v][0]
			if visited[v] {
				return string(c), nil
			}
			visited[v] = true
			c = append(c, v[len(v) - 1])
		}
	}
		
	var v string
	for v = range adj {
		break
	}
	g, err := cycle(v)
	if err != nil {
		return "", err
	}

	gComp, err := cycle(gene.ReverseComplement(g[:k]))
	if err != nil {
		return "", err
	}

	if gComp != gene.ReverseComplement(g) {
		return "", errors.New(gComp + " is not the reverse complement of " + g + ".")
	}

	return g, nil
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
