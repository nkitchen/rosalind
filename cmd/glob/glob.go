package main

import "flag"
import "fmt"
import "log"
import "os"
import "rosalind/gene"
import "rosalind/strings"

func main() {
	debug := flag.Bool("debug", false, "Show debug output")
	flag.Parse()

	strings.DebugEditDistance = *debug

	fasta, err := gene.ReadAllFasta(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	s := fasta[0].Data
	t := fasta[1].Data
	x := strings.MaxAlignmentScore(s, t, gene.Blosum62Scores, -5)
	fmt.Println(x)
}
