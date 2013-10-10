package main

import "fmt"
import "math/big"
import "os"
import "rosalind/gene"

func main() {
	s, _ := gene.ReadFasta(os.Stdin)

	nA := 0
	nC := 0
	for _, b := range s.Data {
		switch b {
		case 'A': nA++
		case 'C': nC++
		}
	}

	matchings := big.NewInt(1)
	for i := 1; i <= nA; i++ {
		matchings.Mul(matchings, big.NewInt(int64(i)))
	}
	for i := 1; i <= nC; i++ {
		matchings.Mul(matchings, big.NewInt(int64(i)))
	}
	fmt.Println(matchings)
}
