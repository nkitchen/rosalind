package main

import "bufio"
import "fmt"
import "log"
import "math"
import "os"
import "rosalind/spectrum"
import "sort"
import "strconv"

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)

	spec := []float64{}

	for scanner.Scan() {
		t := scanner.Text()
		w, err := strconv.ParseFloat(t, 64)
		if err != nil {
			log.Fatal(err)
		}

		spec = append(spec, w)
	}

	sort.Float64Slice(spec).Sort()

	fmt.Println(ProteinFromSpectrum(spec))
}

func ProteinFromSpectrum(spec []float64) string {
	// Amino acid by mass
	a := map[float64]byte{}
	for aa, w := range spectrum.MonoisotopicMass {
		a[snap(w)] = aa
	}
	fmt.Println(a)

	p := []byte{}
	for i := 1; i < len(spec); i++ {
		d := snap(spec[i]) - snap(spec[i - 1])

		aa, ok := a[snap(d)]
		if !ok {
			log.Fatalf("Cannot find residue for %v, %v: %v",
			           spec[i - 1], spec[i], snap(d))
		}
		p = append(p, aa)
	}

	return string(p)
}

const precision = 1e-3

func snap(x float64) float64 {
	y := x * (1 / precision)
	z := math.Floor(y + 0.5) / (1 / precision)
	return z
}
