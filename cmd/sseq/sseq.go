package main

import "fmt"
import "log"
import "os"
import "rosalind/gene"

func main() {
	fasta, err := gene.ReadAllFasta(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	s := fasta[0].Data
	t := fasta[1].Data
	indices := []int{}

	j := 0
	for i := range s {
		if j >= len(t) {
			break
		}

		if s[i] == t[j] {
			indices = append(indices, i)
			j++
		}
	}

	if j < len(t) {
		log.Fatal("Subsequence not found")
	}

	for _, i := range indices {
		fmt.Printf("%v ", i + 1)
	}
	fmt.Println()
}
