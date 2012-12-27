package main

import "bufio"
import "fmt"
import "log"
import "math"
import "strconv"
import "strings"
import "os"

func main() {
	br := bufio.NewReader(os.Stdin)

	line, err := br.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	fields := strings.Fields(line)

	k64, err := strconv.ParseInt(fields[0], 0, 32)
	if err != nil {
		log.Fatal(err)
	}
	k := int(k64)

	n64, err := strconv.ParseInt(fields[1], 0, 32)
	if err != nil {
		log.Fatal(err)
	}
	n := int(n64)

    // Genotype distributions for a single gene
	// Element i of a distribution is the probability for i recessive alleles.
	pGen := [][3]float64{[3]float64{0, 1, 0}}

	alleleDist := [][]float64{
		0: {0: 1.0},
		1: {0: 0.5, 1: 0.5},
		2: {1: 1.0},
	}

	for i := 1; i <= k; i++ {
		var p [3]float64
		//p[0] += pGen[i-1][0] * pGen[0][0] +
		//p[0] += pGen[i-1][0] * pGen[0][1] * 0.5 +
		//p[1] += pGen[i-1][0] * pGen[0][1] * 0.5 +
		//p[1] += pGen[i-1][0] * pGen[0][2] +
		//p[0] += pGen[i-1][1] * pGen[0][0] * 0.5 +
		//p[1] += pGen[i-1][1] * pGen[0][0] * 0.5 +
		//p[0] += pGen[i-1][1] * pGen[0][1] * 0.25
		//p[1] += pGen[i-1][1] * pGen[0][1] * 0.5 +
		//p[2] += pGen[i-1][1] * pGen[1][1] * 0.25 +
		//p[1] += pGen[i-1][1] * pGen[0][2] * 0.5 +
		//p[2] += pGen[i-1][1] * pGen[1][2] * 0.5 +
		//p[1] += pGen[i-1][2] * pGen[0][0] +
		//p[1] += pGen[i-1][2] * pGen[0][1] * 0.5
		//p[2] += pGen[i-1][2] * pGen[1][1] * 0.5 +
		//p[2] += pGen[i-1][2] * pGen[1][2]

		//p[0] = pGen[i-1][0] * pGen[0][0] +
		//       pGen[i-1][0] * pGen[0][1] * 0.5 +
		//	   pGen[i-1][1] * pGen[0][0] * 0.5 +
		//	   pGen[i-1][1] * pGen[0][1] * 0.25
		//p[1] = pGen[i-1][0] * pGen[0][1] * 0.5 +
		//       pGen[i-1][0] * pGen[0][2] +
		//	   pGen[i-1][1] * pGen[0][0] * 0.5 +
		//	   pGen[i-1][1] * pGen[0][1] * 0.25 +
		//	   pGen[i-1][1] * pGen[0][2] * 0.5 +
		//	   pGen[i-1][2] * pGen[0][0] +
		//	   pGen[i-1][2] * pGen[0][1] * 0.5
		//p[2] = pGen[i-1][1] * pGen[1][1] * 0.25 +
		//       pGen[i-1][1] * pGen[1][2] * 0.5 +
		//       pGen[i-1][2] * pGen[1][1] * 0.5 +
		//       pGen[i-1][2] * pGen[1][2]

        for par1 := 0; par1 < 3; par1++ {
			for par2 := 0; par2 < 3; par2++ {
				for u, q1 := range alleleDist[par1] {
					for v, q2 := range alleleDist[par2] {
						p[u+v] += pGen[i-1][par1] * q1 *
						          pGen[0][par2] * q2
					}
				}
			}
		}
		pGen = append(pGen, p)
	}

	p := 0.0
	pow2k := math.Pow(2, float64(k))
	for m := float64(n); m <= pow2k; m++ {
		q := pGen[k][1]
		p += math.Pow(q*q, m) *
		     math.Pow(1 - q*q, pow2k - m) *
			 binom(pow2k, m)
	}
	fmt.Println(p)
} 

func binom(n, k float64) float64 {
	if k > n - k {
		k = n - k
	}
	c := float64(1)
	for i := k; i >= 1; i-- {
		j := k - i
		a := (n - j) / i
		c *= a
	}
	return c
}
