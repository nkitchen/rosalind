package main

import "fmt"
import "log"
import "os"
import "rosalind/gene"
import "rosalind/strings"

func main() {
	fasta, err := gene.ReadAllFasta(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	s := fasta[0].Data
	t := fasta[1].Data
    m := (1 << 27) - 1
    n := strings.NumOptimalAlignments(s, t, m)
	fmt.Println(n)
}
