package main

import "bufio"
import "fmt"
import "rosalind/gene"
import "log"
import "os"
import "strings"

func main() {
	br := bufio.NewReader(os.Stdin)

	line, err := br.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	dna := strings.TrimSpace(line)
	rnaBytes := []byte(dna)
	for i, b := range rnaBytes {
		if b == 'T' {
			rnaBytes[i] = 'U'
		}
	}
	rna := string(rnaBytes)
	revRna := strings.Replace(gene.DnaReverseComplement(dna), "T", "U", -1)

	c := gene.CodonsOf('M')
	if len(c) != 1 {
		panic(fmt.Sprintf("Unexpected start codons: %v", c))
	}
	startCodon := c[0]

	seen := make(map[string]bool)
	for _, s := range []string{rna, revRna} {
		for offset := 0; offset < 3; offset++ {
			for i := offset; i < len(s) - 2; i += 3 {
				if s[i:i+3] == startCodon {
					protein, err := gene.Translate(s[i:])
					if err == nil && !seen[protein] {
						fmt.Println(protein)
						seen[protein] = true
					}
				}
			}
		}
	}
} 
