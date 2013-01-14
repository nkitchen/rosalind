package main

import "fmt"
import "log"
import "os"
import "rosalind/gene"
import "strconv"

func main() {
	const k = 4
	const n = 1 << (2 * 4)

	trMap := [256]byte{'A': '0', 'C': '1', 'G': '2', 'T': '3'}
	comp := make([]int, n)
	dna, err := gene.ReadFasta(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(dna.Data) - k + 1; i++ {
		kmer := dna.Data[i:i+k]
		mapped := make([]byte, k)
		for j := 0; j < k; j++ {
			mapped[j] = trMap[kmer[j]]
		}
		p, err := strconv.ParseInt(string(mapped), 4, 32)
		if err != nil {
			log.Fatal(err)
		}
		comp[int(p)] += 1
	}
	fmt.Println(comp)
}
