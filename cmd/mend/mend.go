package main

import "bytes"
import "fmt"
import "rosalind/tree"
import "os"

// Genotype distribution
type genoDist map[string]float64

func (d genoDist) String() string {
	b := &bytes.Buffer{}
	for _, g := range ([]string{"AA", "Aa", "aa"}) {
		fmt.Fprintf(b, "%.3f ", d[g])
	}
	return b.String()
}

func main() {
	pedigree, _ := tree.ReadNewick(os.Stdin)

	fmt.Println(genotypeDistribution(pedigree))
}

func genotypeDistribution(pedigree *tree.Node) genoDist {
	d := make(genoDist)
	if len(pedigree.Children) == 0 {
		genotype := pedigree.Label
		d[genotype] = 1.0
		return d
	}

	if len(pedigree.Children) != 2 {
		panic("Unexpected number of parents: " + pedigree.String())
	}

	parDist0 := genotypeDistribution(pedigree.Children[0].Node)
	parDist1 := genotypeDistribution(pedigree.Children[1].Node)
	for g0, p0 := range parDist0 {
		for _, allele0 := range g0 {
			for g1, p1 := range parDist1 {
				for _, allele1 := range g1 {
					crossed := []rune{allele0, allele1}
					if allele1 < allele0 {
						crossed[0] = allele1
						crossed[1] = allele0
					}
					g := string(crossed)
					d[g] += 0.5 * 0.5 * p0 * p1
				}
			}
		}
	}
	return d
}
