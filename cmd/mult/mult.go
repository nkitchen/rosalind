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

	s := []string{}
	for _, f := range fasta {
		s = append(s, f.Data)
	}

	scoreFunc := func(c, d int16) int32 {
		if c == d {
			return 0
		}
		return -1
	}
	score, a := strings.MultipleAlignment(s, scoreFunc, '-')
	fmt.Println(score)
	for _, t := range a {
		fmt.Println(t)
	}
}
