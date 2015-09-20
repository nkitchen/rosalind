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

	s := fasta[0]
	t := fasta[1]

	if len(s.Data) != len(t.Data) {
		log.Fatalf("Length mismatch: %d vs. %d", len(s.Data), len(t.Data))
	}

	nTransit := float64(0)
	nTransver := float64(0)
	for i := range s.Data {
		switch pair(s.Data[i], t.Data[i]) {
		case pair('A', 'G'),
			pair('G', 'A'),
			pair('C', 'T'),
			pair('T', 'C'):
			nTransit += 1
		case pair('A', 'C'),
			pair('C', 'A'),
			pair('A', 'T'),
			pair('T', 'A'),
			pair('C', 'G'),
			pair('G', 'C'),
			pair('G', 'T'),
			pair('T', 'G'):
			nTransver += 1
		}
	}

	r := nTransit / nTransver
	fmt.Println(r)
}

func pair(a, b byte) [2]byte {
	return [2]byte{a, b}
}
