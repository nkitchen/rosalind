package main

import "bufio"
import "fmt"
import "log"
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
	p := []byte{}
	for i := 1; i < len(spec); i++ {
		d := spec[i] - spec[i - 1]
		aa, ok := spectrum.ResidueByMass(d)
		if !ok {
			log.Fatalf("Cannot find residue for %v, %v: %v",
			           spec[i - 1], spec[i], d)
		}
		p = append(p, aa)
	}

	return string(p)
}
