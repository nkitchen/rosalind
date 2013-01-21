package main

import "fmt"
import "log"
import "os"
import "rosalind/strings"
import "rosalind/gene"

func main() {
	fasta, err := gene.ReadAllFasta(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	a := make([]string, len(fasta))
	for i, f := range fasta {
		a[i] = f.Data
	}
	d := strings.DistanceMatrix(a, strings.PDistance)
	for _, row := range d {
		for _, x := range row {
			fmt.Printf("%7.5f ", x)
		}
		fmt.Println()
	}
}
