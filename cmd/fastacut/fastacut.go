package main

import "flag"
import "log"
import "os"
import "rosalind/gene"
import "strconv"

func main() {
	flag.Parse()

	var m, n int

	m, _ = strconv.Atoi(flag.Arg(0))
	if flag.NArg() == 1 {
		n = m
	} else {
		n, _ = strconv.Atoi(flag.Arg(1))
	}
		
	fasta, err := gene.ReadAllFasta(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	s := fasta[0]
	t := fasta[1]

	if len(s.Data) > m {
		s.Data = s.Data[:m]
	}
	if len(t.Data) > n {
		t.Data = t.Data[:n]
	}

	gene.WriteFasta(os.Stdout, s)
	gene.WriteFasta(os.Stdout, t)
}
